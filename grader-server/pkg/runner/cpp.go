package runner

import (
	"context"
	"os"
	"os/exec"
	"path/filepath"

	"github.com/DeepAung/gradient/grader-server/proto"
)

var (
	cppBuildCommand = func(filename, name string) []string {
		return []string{"g++", filename, "-o", name}
	}
	cppRunCommand = func(name string) []string {
		return []string{name}
	}
)

type cppRunner struct{}

func NewCppRunner() CodeRunner {
	return cppRunner{}
}

func (r cppRunner) Build(ctx context.Context, codeFilename string) (bool, proto.ResultType) {
	codeExt := filepath.Ext(codeFilename)
	codeName := codeFilename[0 : len(codeFilename)-len(codeExt)]

	buildCommand := cppBuildCommand(codeFilename, codeName)

	cmd := exec.CommandContext(ctx, buildCommand[0], buildCommand[1:]...)
	if err := cmd.Run(); err != nil {
		return false, proto.ResultType_COMPILATION_ERROR
	}
	return true, 0
}

func (r cppRunner) Run(
	ctx context.Context,
	codeFilename, inputFilename string,
) (bool, proto.ResultType) {
	// codeFilename = path1/path2/code.cpp
	// inputFilename = path3/patth4/01.in

	codeExt := filepath.Ext(codeFilename)                        // ".cpp"
	codeName := codeFilename[0 : len(codeFilename)-len(codeExt)] // "path1/path2/code"
	codeDir := filepath.Dir(codeFilename)                        // "path1/path2"

	inputExt := filepath.Ext(inputFilename)                          // ".in"
	inputName := inputFilename[0 : len(inputFilename)-len(inputExt)] // path1/path2/01
	resultFilename := codeDir +
		"/" +
		filepath.Base(inputName) +
		".result" // "path1/path2" + "/" + "01" + ".result"

	runCommand := cppRunCommand(codeName)

	// Open input file
	inputFile, err := os.Open(inputFilename)
	if err != nil {
		return false, proto.ResultType_RUNTIME_ERROR
	}
	defer inputFile.Close()

	// Create result file
	tmp, err := os.Create(resultFilename)
	tmp.Close()
	resultFile, err := os.Create(resultFilename)
	if err != nil {
		return false, proto.ResultType_RUNTIME_ERROR
	}
	defer resultFile.Close()

	// Run run command
	cmd := exec.CommandContext(ctx, runCommand[0], runCommand[1:]...)
	cmd.Stdin = inputFile
	cmd.Stdout = resultFile
	if err := cmd.Run(); err != nil {
		return false, proto.ResultType_RUNTIME_ERROR
	}
	return true, 0
}
