package languages

type Python3 struct{}

func (l Python3) Name() string                    { return "Python3 / python3" }
func (l Python3) DbName() string                  { return "python3" }
func (l Python3) SourceFileExtension() string     { return ".py" }
func (l Python3) ExecutableFileExtension() string { return ".py" }
func (l Python3) GetCompilationCommand(sourceFilename, executableFilename string) []string {
	return []string{}
}

func (l Python3) GetExecutionCommand(executableFilename string) []string {
	return []string{
		"/usr/bin/python3",
		executableFilename,
	}
}
