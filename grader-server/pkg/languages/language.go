package languages

type Language interface {
	Name() string
	DbName() string
	SourceFileExtension() string
	ExecutableFileExtension() string

	// sourceFilename (e.g. "/path/to/code.cpp")
	// executableFilename (e.g. "/path/to/code")
	// GetCompilationCommand can return empty array (e.g. python3)
	GetCompilationCommand(sourceFilename, executableFilename string) []string
	GetExecutionCommand(executableFilename string) []string
}
