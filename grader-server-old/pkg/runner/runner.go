package runner

import (
	"context"
	"errors"
	"fmt"
	"log"
	"os"
	"os/exec"
	"path/filepath"
	"strings"

	"github.com/DeepAung/gradient/grader-server/graderconfig"
	"github.com/DeepAung/gradient/grader-server/proto"
)

var (
	ErrInvalidLanguage    = errors.New("invalid language")
	ErrRunCommandNotFound = errors.New("run command not found")
)

type CodeRunner interface {
	// Build() will return ok, status
	Build(
		ctx context.Context,
		language proto.LanguageType,
		codeFilename string,
	) (bool, proto.StatusType)

	// Run() will return ok, result, err
	// if err != nil then return err
	// if ok then proceed to check the output using checker
	// else return bad result (e.g. runtime error, time limit exceeded)
	Run(
		ctx context.Context,
		language proto.LanguageType,
		codeFilename, inputFilename string,
	) (bool, proto.Result, error)
}

type codeRunner struct {
	cfg *graderconfig.Config
}

func NewCodeRunner(cfg *graderconfig.Config) CodeRunner {
	return &codeRunner{
		cfg: cfg,
	}
}

func (r *codeRunner) Build(
	ctx context.Context,
	language proto.LanguageType,
	codeFilename string,
) (bool, proto.StatusType) {
	languageInfo, ok := r.cfg.GetLanguageInfoFromProto(language)
	if !ok {
		log.Printf("error: %v", ErrInvalidLanguage.Error())
		return false, proto.StatusType_COMPILATION_ERROR
	}

	if languageInfo.BuildCommand == "" {
		return true, 0
	}

	codeName, _ := getNameAndExt(codeFilename)
	buildCommand := parseCommand(languageInfo.BuildCommand, codeFilename, codeName)

	cmd := exec.CommandContext(ctx, buildCommand[0], buildCommand[1:]...)
	if err := cmd.Run(); err != nil {
		return false, proto.StatusType_COMPILATION_ERROR
	}
	return true, 0
}

// codeFilename = path1/path2/code.cpp
// inputFilename = path3/patth4/01.in
func (r *codeRunner) Run(
	ctx context.Context,
	language proto.LanguageType,
	codeFilename, inputFilename string,
) (bool, proto.Result, error) {
	languageInfo, ok := r.cfg.GetLanguageInfoFromProto(language)
	if !ok {
		return false, proto.Result{}, ErrInvalidLanguage
	}

	if languageInfo.RunCommand == "" {
		return false, proto.Result{}, ErrRunCommandNotFound
	}

	codeExt := filepath.Ext(codeFilename)                        // ".cpp"
	codeName := codeFilename[0 : len(codeFilename)-len(codeExt)] // "path1/path2/code"
	codeDir := filepath.Dir(codeFilename)                        // "path1/path2"

	inputExt := filepath.Ext(inputFilename)                          // ".in"
	inputName := inputFilename[0 : len(inputFilename)-len(inputExt)] // path1/path2/01
	resultFilename := codeDir +
		"/" +
		filepath.Base(inputName) +
		".result" // "path1/path2" + "/" + "01" + ".result"

	// Open input file
	inputFile, err := os.Open(inputFilename)
	if err != nil {
		return false, proto.Result{}, err
	}
	defer inputFile.Close()

	// Create result file
	resultFile, err := os.Create(resultFilename)
	if err != nil {
		return false, proto.Result{}, err
	}
	defer resultFile.Close()

	// Run run command
	runCommand := parseCommand(languageInfo.RunCommand, codeFilename, codeName)
	cmd := exec.CommandContext(ctx, runCommand[0], runCommand[1:]...)
	cmd.Stdin = inputFile
	cmd.Stdout = resultFile
	if err := cmd.Run(); err != nil {
		return false, proto.Result{Status: proto.StatusType_RUNTIME_ERROR}, nil
	}
	return true, proto.Result{}, nil
}

func getNameAndExt(filename string) (name string, ext string) {
	ext = filepath.Ext(filename)
	name = filename[0 : len(filename)-len(ext)]
	return
}

func parseCommand(cmd, codeFilename, codeName string) []string {
	cmd = fmt.Sprintf("isolate --run -- %s", cmd)
	return strings.Split(
		strings.NewReplacer("{filename}", codeFilename, "{name}", codeName).Replace(cmd),
		" ",
	)
}
