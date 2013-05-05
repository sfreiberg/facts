package facts

import (
	"bytes"
	"os/exec"
	"strings"
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

func getOSXCodename(version string) string {
	if strings.HasPrefix(version, "10.9") {
		return "Cabernet"
	} else if strings.HasPrefix(version, "10.8") {
		return "Mountain Lion"
	} else if strings.HasPrefix(version, "10.7") {
		return "Lion"
	} else if strings.HasPrefix(version, "10.6") {
		return "Snow Leopard"
	} else if strings.HasPrefix(version, "10.5") {
		return "Leopard"
	} else if strings.HasPrefix(version, "10.4") {
		return "Tiger"
	} else if strings.HasPrefix(version, "10.3") {
		return "Panther"
	} else if strings.HasPrefix(version, "10.2") {
		return "Jaguar"
	} else if strings.HasPrefix(version, "10.1") {
		return "Puma"
	} else if strings.HasPrefix(version, "10.0") {
		return "Cheetah"
	} else {
		return "Unknown"
	}
}
