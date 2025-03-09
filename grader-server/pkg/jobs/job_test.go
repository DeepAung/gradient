package jobs

import (
	"errors"
	"log"
	"os"
	"os/exec"
	"testing"

	"github.com/DeepAung/gradient/grader-server/pkg/languages"
	"github.com/DeepAung/gradient/website-server/pkg/asserts"
)

var (
	boxDir = "/var/local/lib/isolate/0/box"

	pythonSourceFilename = "code.py"
	inputFilename        = "01.in"
	pythonLanguage       = languages.Python3{}
	timeLimitMs          = 1000 // 1 second
	memoryLimitKiB       = 20   // 20 KiB
	pythonCode           string
)

func init() {
	outerSourceFilename := "code.py"
	b, err := os.ReadFile(outerSourceFilename)
	if err != nil {
		log.Fatal(err)
	}
	pythonCode = string(b)
}

func TestCompilationJobSetupFilesPython(t *testing.T) {
	job := NewCompilationJob(
		pythonCode,
		pythonLanguage,
		timeLimitMs,
		memoryLimitKiB,
		pythonSourceFilename,
	)

	err := exec.Command("isolate", "--init").Run()
	asserts.EqualError(t, err, nil)

	sourceFilename, executableFilename, err := job.SetupFiles(boxDir)
	asserts.EqualError(t, err, nil)
	asserts.Equal(t, "sourceFilename", sourceFilename, boxDir+"/code.py")
	asserts.Equal(t, "executableFilename", executableFilename, boxDir+"/code.py")

	testExist(t, sourceFilename)

	err = exec.Command("isolate", "--cleanup").Run()
	asserts.EqualError(t, err, nil)
}

func TestExecutionJobSetupFilesPython(t *testing.T) {
	job := NewExecutionJob(
		pythonCode,
		pythonLanguage,
		timeLimitMs,
		memoryLimitKiB,
		pythonSourceFilename,
		inputFilename,
	)

	err := exec.Command("isolate", "--init").Run()
	asserts.EqualError(t, err, nil)

	executableFilename, innerInputFilename, resultFilename, err := job.SetupFiles(boxDir)
	asserts.EqualError(t, err, nil)
	asserts.Equal(t, "executableFilename", executableFilename, boxDir+"/code.py")
	asserts.Equal(t, "inputFilename", innerInputFilename, boxDir+"/01.in")
	asserts.Equal(t, "resultFilename", resultFilename, boxDir+"/01.result")

	testExist(t, executableFilename)
	testExist(t, innerInputFilename)
	testExist(t, resultFilename)

	err = exec.Command("isolate", "--cleanup").Run()
	asserts.EqualError(t, err, nil)
}

func testExist(t *testing.T, path string) {
	t.Helper()

	if _, err := os.Stat(path); errors.Is(err, os.ErrNotExist) {
		t.Fatal(err)
	}
}
