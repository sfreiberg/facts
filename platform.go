package facts

import (
	"bufio"
	"io"
	"regexp"
	"strconv"
	"strings"
)

// redhatRelease matches /etc/redhat-release style lines such as
// "CentOS Linux release 7.9.2009 (Core)" or "Fedora release 38 (Thirty Eight)".
var redhatRelease = regexp.MustCompile(`^(.*?)\s+release\s+(\S+)(?:\s+\((.*)\))?`)

// parseRedhatRelease extracts the platform, version and codename from the
// first line of /etc/redhat-release. ok is false if the line isn't in the
// expected format.
func parseRedhatRelease(line string) (platform, version, codename string, ok bool) {
	m := redhatRelease.FindStringSubmatch(strings.TrimSpace(line))
	if m == nil {
		return "", "", "", false
	}
	return m[1], m[2], m[3], true
}

// parseLSBRelease extracts the platform, version and codename from the
// contents of /etc/lsb-release.
func parseLSBRelease(r io.Reader) (platform, version, codename string) {
	scanner := bufio.NewScanner(r)
	for scanner.Scan() {
		k, v, ok := strings.Cut(scanner.Text(), "=")
		if !ok {
			continue
		}
		v = strings.Trim(strings.TrimSpace(v), `"'`)
		switch strings.TrimSpace(k) {
		case "DISTRIB_ID":
			platform = v
		case "DISTRIB_RELEASE":
			version = v
		case "DISTRIB_CODENAME":
			codename = v
		}
	}
	return platform, version, codename
}

var macOSCodenames = map[int]string{
	11: "Big Sur",
	12: "Monterey",
	13: "Ventura",
	14: "Sonoma",
	15: "Sequoia",
	26: "Tahoe",
}

var macOSXCodenames = map[int]string{
	0:  "Cheetah",
	1:  "Puma",
	2:  "Jaguar",
	3:  "Panther",
	4:  "Tiger",
	5:  "Leopard",
	6:  "Snow Leopard",
	7:  "Lion",
	8:  "Mountain Lion",
	9:  "Mavericks",
	10: "Yosemite",
	11: "El Capitan",
	12: "Sierra",
	13: "High Sierra",
	14: "Mojave",
	15: "Catalina",
}

// getOSXCodename returns the marketing name for a macOS version such as
// "10.15.7" or "14.2".
func getOSXCodename(version string) string {
	parts := strings.Split(strings.TrimSpace(version), ".")
	major, err := strconv.Atoi(parts[0])
	if err != nil {
		return "Unknown"
	}
	if major != 10 {
		if name, ok := macOSCodenames[major]; ok {
			return name
		}
		return "Unknown"
	}
	if len(parts) < 2 {
		return "Unknown"
	}
	minor, err := strconv.Atoi(parts[1])
	if err != nil {
		return "Unknown"
	}
	if name, ok := macOSXCodenames[minor]; ok {
		return name
	}
	return "Unknown"
}
