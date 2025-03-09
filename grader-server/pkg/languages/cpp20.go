package languages

type CPP20 struct{}

func (l CPP20) Name() string                    { return "C++20 / g++" }
func (l CPP20) DbName() string                  { return "cpp20" }
func (l CPP20) SourceFileExtension() string     { return ".cpp" }
func (l CPP20) ExecutableFileExtension() string { return "" }
func (l CPP20) GetCompilationCommand(sourceFilename, executableFilename string) []string {
	return []string{
		"/usr/bin/g++",
		"-std=gnu++20",
		"-O2",
		"-pipe",
		"-static",
		"-s",
		"-o",
		executableFilename,
		sourceFilename,
	}
}

func (l CPP20) GetExecutionCommand(executableFilename string) []string {
	return []string{executableFilename}
}
