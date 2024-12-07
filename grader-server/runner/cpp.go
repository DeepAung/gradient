package runner

import (
	"context"
	"os"
	"os/exec"
	"path/filepath"
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

func (r cppRunner) Run(
	ctx context.Context,
	codeFilename, inputFilename, outputFilename string,
) error {
	codeExt := filepath.Ext(codeFilename)
	codeName := codeFilename[0 : len(codeFilename)-len(codeExt)]

	buildCommand := cppBuildCommand(codeFilename, codeName)
	runCommand := cppRunCommand(codeName)

	// run build command
	cmd := exec.CommandContext(ctx, buildCommand[0], buildCommand[1:]...)
	if err := cmd.Run(); err != nil {
		return err
	}

	inputFile, err := os.Open(inputFilename)
	if err != nil {
		return err
	}
	defer inputFile.Close()

	outputFile, err := os.Create(outputFilename)
	if err != nil {
		return err
	}
	defer outputFile.Close()

	// run run command
	cmd = exec.CommandContext(ctx, runCommand[0], runCommand[1:]...)
	cmd.Stdin = inputFile
	cmd.Stdout = outputFile
	return cmd.Run()
}
