package runner

import (
	"context"
	"fmt"

	"github.com/DeepAung/gradient/grader-server/proto"
)

type CodeRunner interface {
	Build(ctx context.Context, codeFilename string) (bool, proto.StatusType)
	Run(ctx context.Context, codeFilename, inputFilename string) (bool, proto.StatusType)
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
