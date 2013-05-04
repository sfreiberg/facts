package facts

import (
	"bytes"
	"os/exec"
)

func getOs() string {
	return "linux"
}

func getArch() string {
	cmd := exec.Command("/bin/sh", "-c", "/bin/uname -m")
	out, err := cmd.Output()
	if err != nil {
		return ""
	}
	out = bytes.TrimSpace(out)
	return string(out)
}
