package runner

import (
	"context"
	"fmt"

	"github.com/DeepAung/gradient/grader-server/proto"
)

type CodeRunner interface {
	Build(ctx context.Context, codeFilename string) (bool, proto.ResultType)
	Run(ctx context.Context, codeFilename, inputFilename string) (bool, proto.ResultType)
}

var codeRunners = map[proto.LanguageType]func() CodeRunner{
	proto.LanguageType_CPP: NewCppRunner,
}

func NewCodeRunner(language proto.LanguageType) (CodeRunner, error) {
	runner, ok := codeRunners[language]
	if !ok {
		return nil, fmt.Errorf("no code runner for language %q", language)
	}

	return runner(), nil
}
