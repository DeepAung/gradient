package runner_test

import (
	"os"
	"os/exec"
	"testing"
)

func TestRunPython(t *testing.T) {
	runCommand(t, "isolate --cleanup")

	runCommand(t, "isolate --init")
	// defer runCommand(t, "isolate --cleanup")

	runCommand(t, "cp code.py /var/local/lib/isolate/0/box/")
	runCommand(t,
		"cp 01.in /var/local/lib/isolate/0/box/",
	)
	runCommand(t, "touch /var/local/lib/isolate/0/box/01.result")
	// 20480 bytes = 20 * 1024 = 20 KiB
	runCommand(
		t,
		"isolate --wait -i 01.in -o 01.result -t 0.1 -m 10000 -M run.log --run -- /usr/bin/python3 code.py",
	)
}

// func TestBuildCPP(t *testing.T) {
// 	runCommand(t, "isolate --init")
// 	defer runCommand(t, "isolate --cleanup")
//
// 	runCommand(t, "cp ../../examples/code.cpp /var/local/lib/isolate/0/box/")
// 	runCommand(
// 		t,
// 		"isolate --wait --run -- /usr/bin/g++ -DEVAL -std=gnu++20 -O2 -pipe -static -s -o code code.cpp",
// 	)
// 	// runCommand(t,
// 	// 	"cp ../../../website-server/migrations/testcases/1/01.in /var/local/lib/isolate/0/box/",
// 	// )
// 	// runCommand(t, "touch /var/local/lib/isolate/0/box/01.result")
// 	// runCommand(t, "isolate --run -i 01.in -o 01.result -- /usr/bin/python3 code.py")
// }
//
// func TestRunCPP(t *testing.T) {
// 	runCommand(t, "isolate --init")
// 	defer runCommand(t, "isolate --cleanup")
//
// 	runCommand(t,
// 		"cp ../../../website-server/migrations/testcases/1/01.in /var/local/lib/isolate/0/box/",
// 	)
// 	runCommand(t, "touch /var/local/lib/isolate/0/box/01.result")
// 	runCommand(t, "isolate --wait --run -i 01.in -o 01.result -- code")
// }

func runCommand(t *testing.T, cmdStr string) {
	t.Helper()
	cmd := exec.Command("bash", "-c", cmdStr)
	cmd.Stdout = os.Stdout
	if err := cmd.Run(); err != nil {
		t.Fatal(err.Error())
	}
}
