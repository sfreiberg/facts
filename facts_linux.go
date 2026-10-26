package facts

import (
	"bytes"
	"os"
	"os/exec"
	"strings"
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

func (f *Facts) loadPlatformInfo() error {
	// Check for Redhat derived OS
	if b, err := os.ReadFile("/etc/redhat-release"); err == nil {
		line, _, _ := strings.Cut(string(b), "\n")
		if platform, version, codename, ok := parseRedhatRelease(line); ok {
			f.Platform = platform
			f.PlatformVersion = version
			f.PlatformCodename = codename
			return nil
		}
	}

	// Look for lsb-release file. Known to work on Ubuntu.
	if file, err := os.Open("/etc/lsb-release"); err == nil {
		defer file.Close()
		f.Platform, f.PlatformVersion, f.PlatformCodename = parseLSBRelease(file)
	}
	return nil
}
