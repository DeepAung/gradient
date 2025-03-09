package sandboxes

import (
	"bytes"
	"context"
	"fmt"
	"log/slog"
	"os"
	"os/exec"
	"strconv"
	"strings"
	"time"

	"github.com/DeepAung/gradient/grader-server/pkg/jobs"
)

const channelLen = 10

var nextId = 0

type IsolateSandBox struct {
	boxId            int
	baseBoxDir       string
	sandBoxDirectory string
	maxTimeMs        int // in millisecond
	maxMemoryKiB     int // in KiB
	logger           *slog.Logger

	payloadCh chan payload
}

type payload struct {
	jobGroup      *jobs.JobGroup
	compilationCh chan<- jobs.CompilationResult
	executionCh   chan<- jobs.ExecutionResult
}

// One SandBox represent one isolated container that have one thread.
// Added jobGroup needs to wait until to all prior jobGroups is finished to be able to processed and get the result.
func NewIsolateSandBox(
	baseBoxDir string, // default = "/var/local/lib/isolate"
	sandBoxDirectory string, // default = "/tmp"
	maxTimeMs int, // default = 20000 (20 second)
	maxMemoryKiB int, // default = 512 * 1024 (512 MiB)
	logger *slog.Logger,
) *IsolateSandBox {
	sandbox := &IsolateSandBox{
		boxId:            nextId,
		baseBoxDir:       baseBoxDir,
		sandBoxDirectory: sandBoxDirectory,
		maxTimeMs:        maxTimeMs,
		maxMemoryKiB:     maxMemoryKiB,
		logger:           logger,

		payloadCh: make(chan payload, channelLen),
	}
	nextId++

	// run loop infinitely
	go func() {
		for {
			payload := <-sandbox.payloadCh
			sandbox.runJobGroup(payload.jobGroup, payload.compilationCh, payload.executionCh)
		}
	}()

	return sandbox
}

func (s *IsolateSandBox) RunJobGroup(
	jobGroup *jobs.JobGroup,
	compilationCh chan jobs.CompilationResult,
	executionCh chan jobs.ExecutionResult,
) {
	s.payloadCh <- payload{
		jobGroup:      jobGroup,
		compilationCh: compilationCh,
		executionCh:   executionCh,
	}
}

func (s *IsolateSandBox) runJobGroup(
	jobGroup *jobs.JobGroup,
	compilationCh chan<- jobs.CompilationResult,
	executionCh chan<- jobs.ExecutionResult,
) {
	s.initBox()

	compilationCh <- s.runCompilationJob(&jobGroup.CompilationJob)
	close(compilationCh)
	for _, executableJob := range jobGroup.ExecutionJobs {
		executionCh <- s.runExecutionJob(&executableJob)
	}
	close(executionCh)

	s.cleanupBox()
}

func (s *IsolateSandBox) runCompilationJob(job *jobs.CompilationJob) jobs.CompilationResult {
	if job.TimeLimitMs() > s.maxTimeMs {
		s.logger.Error(
			fmt.Sprintf(
				"job's time limit (%d ms) exceed sandbox's time limit (%d ms)",
				job.TimeLimitMs(),
				s.maxTimeMs,
			),
		)
	}

	if job.MemoryLimitKiB() > s.maxMemoryKiB {
		s.logger.Error(
			fmt.Sprintf(
				"job's memory limit (%d KiB) exceed sandbox's memory limit (%d KiB)",
				job.MemoryLimitKiB(),
				s.maxMemoryKiB,
			),
		)
	}

	sourceFilename, executableFilename, err := job.SetupFiles(s.getBoxDir())
	if err != nil {
		s.logger.Error(err.Error())
		return jobs.CompilationResult{
			Status:    jobs.CompilationError,
			TimeMs:    0,
			MemoryKiB: 0,
		}
	}
	cmd := job.GetCommand(sourceFilename, executableFilename)

	prefix := []string{
		"isolate",
		"-b",
		fmt.Sprint(s.boxId),
		"-M",
		s.getLogFilename(),
		"-t",
		job.TimeLimitFormatSeconds(),
		"-m",
		job.MemoryLimitFormatBytes(),
		"--wait",
		"--run",
		"--",
	}
	cmd = append(prefix, cmd...)

	ctx, cancel := context.WithTimeout(
		context.Background(),
		time.Millisecond*time.Duration(s.maxTimeMs),
	)
	defer cancel()

	execCmd := exec.CommandContext(ctx, cmd[0], cmd[1:]...)
	errBuffer := bytes.NewBufferString("")
	execCmd.Stderr = errBuffer

	if err := execCmd.Run(); err != nil {
		s.logger.Error(err.Error())
	}
	if errMsg := errBuffer.String(); errMsg != "" {
		s.logger.Error(errMsg)
	}

	log, err := s.readLog()
	if err != nil {
		s.logger.Error(err.Error())
	}
	status := jobs.CompilationPass
	if log.Status != "" {
		status = jobs.CompilationError
	}

	return jobs.CompilationResult{
		Status:    status,
		TimeMs:    log.TimeMs,
		MemoryKiB: log.MemoryKiB,
	}
}

func (s *IsolateSandBox) runExecutionJob(job *jobs.ExecutionJob) jobs.ExecutionResult {
	if job.TimeLimitMs() > s.maxTimeMs {
		s.logger.Error(
			fmt.Sprintf(
				"job's time limit (%d ms) exceed sandbox's time limit (%d ms)",
				job.TimeLimitMs(),
				s.maxTimeMs,
			),
		)
	}

	if job.MemoryLimitKiB() > s.maxMemoryKiB {
		s.logger.Error(
			fmt.Sprintf(
				"job's memory limit (%d KiB) exceed sandbox's memory limit (%d KiB)",
				job.MemoryLimitKiB(),
				s.maxMemoryKiB,
			),
		)
	}

	executableFilename, inputFilename, resultFilename, err := job.SetupFiles(s.getBoxDir())
	if err != nil {
		s.logger.Error(err.Error())
		return jobs.ExecutionResult{
			Status:    jobs.RuntimeError,
			TimeMs:    0,
			MemoryKiB: 0,
		}
	}
	cmd := job.GetCommand(executableFilename)

	prefix := []string{
		"isolate",
		"-b",
		fmt.Sprint(s.boxId),
		"-i",
		inputFilename,
		"-o",
		resultFilename,
		"-M",
		s.getLogFilename(),
		"-t",
		job.TimeLimitFormatSeconds(),
		"-m",
		job.MemoryLimitFormatBytes(),
		"--wait",
		"--run",
		"--",
	}
	cmd = append(prefix, cmd...)

	ctx, cancel := context.WithTimeout(
		context.Background(),
		time.Millisecond*time.Duration(s.maxTimeMs),
	)
	defer cancel()

	execCmd := exec.CommandContext(ctx, cmd[0], cmd[1:]...)
	errBuffer := bytes.NewBufferString("")
	execCmd.Stderr = errBuffer

	if err := execCmd.Run(); err != nil {
		s.logger.Error(err.Error())
	}
	if errMsg := errBuffer.String(); errMsg != "" {
		s.logger.Error(errMsg)
	}

	log, err := s.readLog()
	if err != nil {
		s.logger.Error(err.Error())
	}
	var status jobs.ExecutionStatus
	switch log.Status {
	case "RE":
		status = jobs.RuntimeError
	case "SG":
		status = jobs.ExitSignalError
	case "TO":
		status = jobs.TimeLimitExceeded
	case "XX":
		status = jobs.SandBoxError
	default:
		status = jobs.ExecutionPass
	}

	return jobs.ExecutionResult{
		Status:    status,
		TimeMs:    log.TimeMs,
		MemoryKiB: log.MemoryKiB,
	}
}

func (s *IsolateSandBox) initBox() {
	exec.Command("isolate", "-b", fmt.Sprint(s.boxId), "--init")
}

func (s *IsolateSandBox) cleanupBox() {
	exec.Command("isolate", "-b", fmt.Sprint(s.boxId), "--cleanup")
}

func (s *IsolateSandBox) getBoxDir() string {
	return fmt.Sprintf("%s/%d/box", s.baseBoxDir, s.boxId)
}

func (s *IsolateSandBox) getLogFilename() string {
	return fmt.Sprintf("%s/run.log.%d", s.sandBoxDirectory, s.boxId)
}

func (s *IsolateSandBox) readLog() (IsolateLog, error) {
	log := IsolateLog{}

	b, err := os.ReadFile(s.getLogFilename())
	if err != nil {
		return IsolateLog{}, err
	}

	iter := strings.SplitSeq(string(b), "\n")
	for line := range iter {
		tmp := strings.Split(line, ":")
		key, value := tmp[0], tmp[1]

		switch key {
		case "time":
			timeSeconds, err := strconv.ParseFloat(value, 32)
			if err == nil {
				log.TimeMs = int(timeSeconds * 1000)
			}
		case "max-rss":
			memoryBytes, err := strconv.Atoi(value)
			if err == nil {
				log.MemoryKiB = memoryBytes / 1024
			}
		case "status":
			log.Status = value
		case "message":
			log.Message = value

		}
	}

	return log, nil
}

type IsolateLog struct {
	TimeMs    int
	MemoryKiB int
	Status    string
	Message   string
}
