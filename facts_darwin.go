package facts

import (
	"bytes"
	"os/exec"
)

func getOs() string {
	return "darwin"
}

func getArch() string {
	cmd := exec.Command("/bin/sh", "-c", "/usr/bin/uname -m")
	out, err := cmd.Output()
	if err != nil {
		return ""
	}
	out = bytes.TrimSpace(out)
	return string(out)
}

func (f *Facts) loadPlatformInfo() error {
	// Figure out what version
	cmd := exec.Command("/bin/sh", "-c", "/usr/bin/sw_vers -productVersion")
	out, err := cmd.Output()
	if err != nil {
		return err
	}
	f.PlatformVersion = string(bytes.TrimSpace(out))

	// Figure out the platform
	cmd = exec.Command("/bin/sh", "-c", "/usr/bin/sw_vers -productName")
	out, err = cmd.Output()
	if err != nil {
		return err
	}
	f.Platform = string(bytes.TrimSpace(out))

	// Set the codename
	f.PlatformCodename = getOSXCodename(f.PlatformVersion)

	return nil
}
