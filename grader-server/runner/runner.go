package runner

import "context"

type CodeRunner interface {
	Run(ctx context.Context, codeFilename, inputFilename, outputFilename string) error
}
