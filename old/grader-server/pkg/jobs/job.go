package jobs

import (
	"bytes"
	"context"
	"errors"
	"fmt"
	"os/exec"
	"path/filepath"
	"strings"
	"time"

	"github.com/DeepAung/gradient/grader-server/pkg/languages"
)

type JobGroup struct {
	CompilationJob CompilationJob
	ExecutionJobs  []ExecutionJob
}

// -------------------------------------------------------------------------- //

type baseJob struct {
	code           string
	language       languages.Language
	timeLimitMs    int
	memoryLimitKiB int
}

func (j *baseJob) Code() string                 { return j.code }
func (j *baseJob) Language() languages.Language { return j.language }
func (j *baseJob) TimeLimitMs() int             { return j.timeLimitMs }
func (j *baseJob) MemoryLimitKiB() int          { return j.memoryLimitKiB }

func (j *baseJob) TimeLimitFormatSeconds() string {
	return fmt.Sprintf("%d.%03d", j.timeLimitMs/1000, j.timeLimitMs%1000)
}

func (j *baseJob) MemoryLimitFormatBytes() string {
	return fmt.Sprint(j.memoryLimitKiB * 1024)
}

// -------------------------------------------------------------------------- //

type CompilationJob struct {
	baseJob
	outerSourceFilename string
}

type CompilationResult struct {
	Status    CompilationStatus
	TimeMs    int // in milliseconds
	MemoryKiB int // in KiB
}

type CompilationStatus string

const (
	CompilationPass  CompilationStatus = "Compilation Pass"
	CompilationError CompilationStatus = "Compilation Error"
)

func NewCompilationJob(
	code string,
	language languages.Language,
	timeLimitMs int,
	memoryLimitKiB int,
	outerSourceFilename string,
) *CompilationJob {
	return &CompilationJob{
		baseJob: baseJob{
			code:           code,
			language:       language,
			timeLimitMs:    timeLimitMs,
			memoryLimitKiB: memoryLimitKiB,
		},
		outerSourceFilename: outerSourceFilename,
	}
}

// copy sourceFilename to box directory
// return
func (j *CompilationJob) SetupFiles(
	boxDir string,
) (innerSourceFilename, innerExecutableFilename string, err error) {
	ctx, cancel := context.WithTimeout(context.Background(), 20*time.Second)
	defer cancel()

	innerSourceFilename = fmt.Sprintf(
		"%s/%s",
		boxDir,
		filepath.Base(j.outerSourceFilename),
	)
	if err = runCommand(ctx, "cp", j.outerSourceFilename, innerSourceFilename); err != nil {
		return "", "", err
	}

	withoutExt := strings.TrimSuffix(innerSourceFilename, filepath.Ext(innerSourceFilename))
	ext := j.language.ExecutableFileExtension()
	innerExecutableFilename = fmt.Sprintf("%s%s", withoutExt, ext)

	return innerSourceFilename, innerExecutableFilename, nil
}

func (j *CompilationJob) GetCommand(sourceFilename, executableFilename string) []string {
	return j.language.GetCompilationCommand(sourceFilename, executableFilename)
}

// -------------------------------------------------------------------------- //

type ExecutionJob struct {
	baseJob
	outerExecutableFilename string
	outerInputFilename      string
}

type ExecutionResult struct {
	Status    ExecutionStatus
	TimeMs    int // in milliseconds
	MemoryKiB int // in KiB
}

type ExecutionStatus string

var (
	ExecutionPass     ExecutionStatus = "Execution Pass"
	RuntimeError      ExecutionStatus = "Runtime Error"
	TimeLimitExceeded ExecutionStatus = "Time limit Exceeded"
	SandBoxError      ExecutionStatus = "SandBox Error"
	ExitSignalError   ExecutionStatus = "Exit Signal"
	// Pass ExecutionStatus = "Pass"
	// Incorrect ExecutionStatus = "Incorrect"
	// MemoryLimitExceeded ExecutionStatus = "Memory limit Exceeded"
)

func NewExecutionJob(
	code string,
	language languages.Language,
	timeLimitMs int,
	memoryLimitKiB int,
	outerExecutableFilename string,
	outerInputFilename string,
) *ExecutionJob {
	return &ExecutionJob{
		baseJob: baseJob{
			code:           code,
			language:       language,
			timeLimitMs:    timeLimitMs,
			memoryLimitKiB: memoryLimitKiB,
		},
		outerExecutableFilename: outerExecutableFilename,
		outerInputFilename:      outerInputFilename,
	}
}

// copy executableFilename & inputFilename to box directory
// create resultFilename in box directory
// return
func (j *ExecutionJob) SetupFiles(
	boxDir string,
) (innerExecutableFilename, innerInputFilename, innerResultFilename string, err error) {
	ctx, cancel := context.WithTimeout(context.Background(), 20*time.Second)
	defer cancel()

	innerExecutableFilename = fmt.Sprintf(
		"%s/%s",
		boxDir,
		filepath.Base(j.outerExecutableFilename),
	)
	if err = runCommand(ctx, "cp", j.outerExecutableFilename, innerExecutableFilename); err != nil {
		return
	}

	innerInputFilename = fmt.Sprintf("%s/%s", boxDir, filepath.Base(j.outerInputFilename))
	if err = runCommand(ctx, "cp", j.outerInputFilename, innerInputFilename); err != nil {
		return
	}

	withoutExt := strings.TrimSuffix(j.outerInputFilename, filepath.Ext(j.outerInputFilename))
	innerResultFilename = fmt.Sprintf("%s/%s.result", boxDir, filepath.Base(withoutExt))
	if err = runCommand(ctx, "touch", innerResultFilename); err != nil {
		return
	}

	return
}

func (j *ExecutionJob) GetCommand(executableFilename string) []string {
	return j.language.GetExecutionCommand(executableFilename)
}

// -------------------------------------------------------------------------- //

func runCommand(ctx context.Context, name string, args ...string) error {
	bufferStr := bytes.NewBufferString("")

	execCmd := exec.CommandContext(ctx, name, args...)
	execCmd.Stderr = bufferStr
	if err := execCmd.Run(); err != nil {
		return err
	}
	if errMsg := bufferStr.String(); errMsg != "" {
		return errors.New(errMsg)
	}
	return nil
}
