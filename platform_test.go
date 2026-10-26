package facts

import (
	"strings"
	"testing"
)

func TestGetOSXCodename(t *testing.T) {
	tests := map[string]string{
		"10.0":    "Cheetah",
		"10.1.5":  "Puma",
		"10.6.8":  "Snow Leopard",
		"10.9.5":  "Mavericks",
		"10.10":   "Yosemite",
		"10.11.6": "El Capitan",
		"10.12":   "Sierra",
		"10.15.7": "Catalina",
		"11.7.10": "Big Sur",
		"14.2.1":  "Sonoma",
		"15.0":    "Sequoia",
		"26.1":    "Tahoe",
		"99.0":    "Unknown",
		"10":      "Unknown",
		"":        "Unknown",
		"abc":     "Unknown",
	}
	for version, want := range tests {
		if got := getOSXCodename(version); got != want {
			t.Errorf("getOSXCodename(%q) = %q, want %q", version, got, want)
		}
	}
}

func TestParseRedhatRelease(t *testing.T) {
	tests := []struct {
		line                        string
		platform, version, codename string
		ok                          bool
	}{
		{"CentOS Linux release 7.9.2009 (Core)\n", "CentOS Linux", "7.9.2009", "Core", true},
		{"Fedora release 38 (Thirty Eight)", "Fedora", "38", "Thirty Eight", true},
		{"Rocky Linux release 9.3 (Blue Onyx)", "Rocky Linux", "9.3", "Blue Onyx", true},
		{"Red Hat Enterprise Linux release 9.4 (Plow)", "Red Hat Enterprise Linux", "9.4", "Plow", true},
		{"CentOS Stream release 9", "CentOS Stream", "9", "", true},
		{"", "", "", "", false},
		{"garbage", "", "", "", false},
	}
	for _, tt := range tests {
		p, v, c, ok := parseRedhatRelease(tt.line)
		if p != tt.platform || v != tt.version || c != tt.codename || ok != tt.ok {
			t.Errorf("parseRedhatRelease(%q) = %q, %q, %q, %v; want %q, %q, %q, %v",
				tt.line, p, v, c, ok, tt.platform, tt.version, tt.codename, tt.ok)
		}
	}
}

func TestParseLSBRelease(t *testing.T) {
	in := `DISTRIB_ID=Ubuntu
DISTRIB_RELEASE=24.04
DISTRIB_CODENAME=noble
DISTRIB_DESCRIPTION="Ubuntu 24.04.1 LTS"`
	p, v, c := parseLSBRelease(strings.NewReader(in))
	if p != "Ubuntu" || v != "24.04" || c != "noble" {
		t.Fatalf("parseLSBRelease = %q, %q, %q", p, v, c)
	}
}
