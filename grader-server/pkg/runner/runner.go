package runner

import (
	"context"
	"errors"
	"log"
	"os"
	"os/exec"
	"path/filepath"
	"strings"

	"github.com/DeepAung/gradient/grader-server/graderconfig"
	"github.com/DeepAung/gradient/grader-server/proto"
)

var ErrInvalidLanguage = errors.New("invalid language")

type CodeRunner interface {
	Build(
		ctx context.Context,
		language proto.LanguageType,
		codeFilename string,
	) (bool, proto.StatusType)
	Run(
		ctx context.Context,
		language proto.LanguageType,
		codeFilename, inputFilename string,
	) (bool, proto.Result)
}

type codeRunner struct {
	graderCfg *graderconfig.Config
}

func NewCodeRunner(graderCfg *graderconfig.Config) CodeRunner {
	return &codeRunner{
		graderCfg: graderCfg,
	}
}

func (r *codeRunner) Build(
	ctx context.Context,
	language proto.LanguageType,
	codeFilename string,
) (bool, proto.StatusType) {
	languageInfo, ok := r.graderCfg.GetLanguageInfoFromProto(language)
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
) (bool, proto.StatusType) {
	languageInfo, ok := r.graderCfg.GetLanguageInfoFromProto(language)
	if !ok {
		log.Printf("error: %v", ErrInvalidLanguage.Error())
		return false, proto.StatusType_COMPILATION_ERROR
	}

	if languageInfo.RunCommand == "" {
		return true, 0
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
		return false, proto.StatusType_RUNTIME_ERROR
	}
	defer inputFile.Close()

	// Create result file
	resultFile, err := os.Create(resultFilename)
	if err != nil {
		return false, proto.StatusType_RUNTIME_ERROR
	}
	defer resultFile.Close()

	// Run run command
	runCommand := parseCommand(languageInfo.RunCommand, codeFilename, codeName)
	cmd := exec.CommandContext(ctx, runCommand[0], runCommand[1:]...)
	cmd.Stdin = inputFile
	cmd.Stdout = resultFile
	if err := cmd.Run(); err != nil {
		return false, proto.StatusType_RUNTIME_ERROR
	}
	return true, 0
}

func getNameAndExt(filename string) (name string, ext string) {
	ext = filepath.Ext(filename)
	name = filename[0 : len(filename)-len(ext)]
	return
}

func parseCommand(cmd, codeFilename, codeName string) []string {
	return strings.Split(
		strings.NewReplacer("{filename}", codeFilename, "{name}", codeName).Replace(cmd),
		" ",
	)
}
