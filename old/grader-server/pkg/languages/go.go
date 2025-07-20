package languages

type Go struct{}

func (l Go) Name() string                    { return "Go / go" }
func (l Go) DbName() string                  { return "go" }
func (l Go) SourceFileExtension() string     { return ".go" }
func (l Go) ExecutableFileExtension() string { return "" }
func (l Go) GetCompilationCommand(sourceFilename, executableFilename string) []string {
	return []string{
		"/usr/bin/go",
		"-o",
		executableFilename,
		sourceFilename,
	}
}

func (l Go) GetExecutionCommand(executableFilename string) []string {
	return []string{executableFilename}
}
