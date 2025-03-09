package languages

type C11 struct{}

func (l C11) Name() string                    { return "C11 / gcc" }
func (l C11) DbName() string                  { return "c11" }
func (l C11) SourceFileExtension() string     { return ".c" }
func (l C11) ExecutableFileExtension() string { return "" }
func (l C11) GetCompilationCommand(sourceFilename, executableFilename string) []string {
	return []string{
		"/usr/bin/gcc",
		"-std=gnu11",
		"-O2",
		"-pipe",
		"-static",
		"-s",
		"-o",
		executableFilename,
		sourceFilename,
		"-lm",
	}
}

func (l C11) GetExecutionCommand(executableFilename string) []string {
	return []string{executableFilename}
}
