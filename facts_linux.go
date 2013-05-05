package facts

import (
	"bufio"
	"bytes"
	"io"
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
	if file, err := os.Open("/etc/redhat-release"); err == nil {
		defer file.Close()
		reader := bufio.NewReader(file)
		line, _ := reader.ReadString('\n')
		fields := strings.Split(line, " ")

		// Grab codename and clean it up. Usually looks like (Codename)
		codename := fields[len(fields)-1]
		codename = strings.TrimSpace(codename)
		codename = strings.TrimLeft(codename, "(")
		codename = strings.TrimRight(codename, ")")
		f.PlatformCodename = codename

		f.PlatformVersion = fields[len(fields)-2]
		f.Platform = strings.Join(fields[:len(fields)-3], " ")

		return nil
	}

	// Look for lsb-release file. Known to work on Ubuntu.
	if file, err := os.Open("/etc/lsb-release"); err == nil {
		defer file.Close()
		reader := bufio.NewReader(file)
		for {
			line, err := reader.ReadString('\n')
			fields := strings.SplitN(line, "=", 2)
			if len(fields) == 2 {
				k, v := fields[0], fields[1]
				v = strings.TrimSpace(v)
				switch k {
				case "DISTRIB_ID":
					f.Platform = v
				case "DISTRIB_RELEASE":
					f.PlatformVersion = v
				case "DISTRIB_CODENAME":
					f.PlatformCodename = v
				}
			}
			if err == io.EOF {
				break
			}
		}
	}
	return nil
}
