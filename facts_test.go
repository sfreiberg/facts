package facts

import (
	"encoding/json"
	"os"
	"runtime"
	"strings"
	"testing"

	"go.yaml.in/yaml/v3"
)

func TestFindFacts(t *testing.T) {
	f := FindFacts()

	if f.Os != runtime.GOOS {
		t.Errorf("Os = %q, want %q", f.Os, runtime.GOOS)
	}
	if f.Cpus != runtime.NumCPU() {
		t.Errorf("Cpus = %d, want %d", f.Cpus, runtime.NumCPU())
	}
	if f.Arch == "" {
		t.Error("Arch is empty")
	}
	if f.Memory == 0 {
		t.Error("Memory is 0")
	}

	hostname, err := os.Hostname()
	if err != nil {
		t.Fatal(err)
	}
	if f.Fqdn != hostname {
		t.Errorf("Fqdn = %q, want %q", f.Fqdn, hostname)
	}
	if !strings.HasPrefix(hostname, f.Hostname) || f.Hostname == "" {
		t.Errorf("Hostname = %q is not a prefix of %q", f.Hostname, hostname)
	}
	if f.Domain != "" && f.Hostname+"."+f.Domain != hostname {
		t.Errorf("Hostname %q + Domain %q != %q", f.Hostname, f.Domain, hostname)
	}

	if len(f.Interfaces) == 0 {
		t.Error("no network interfaces found")
	}
	for name, iface := range f.Interfaces {
		if iface.Name != name {
			t.Errorf("interface key %q has Name %q", name, iface.Name)
		}
	}
}

func TestSerialization(t *testing.T) {
	f := &Facts{
		Hostname:   "web01",
		Domain:     "example.com",
		Cpus:       4,
		Memory:     2048,
		Interfaces: map[string]Interface{"eth0": {Name: "eth0", MTU: 1500}},
	}

	b, err := f.ToJson()
	if err != nil {
		t.Fatal(err)
	}
	var m map[string]any
	if err := json.Unmarshal(b, &m); err != nil {
		t.Fatal(err)
	}
	if m["hostname"] != "web01" || m["cpus"] != 4.0 || m["memory"] != 2048.0 {
		t.Errorf("unexpected JSON %s", b)
	}

	pretty, err := f.ToPrettyJson()
	if err != nil {
		t.Fatal(err)
	}
	if !strings.Contains(string(pretty), "\n  \"hostname\": \"web01\"") {
		t.Errorf("pretty JSON is not indented: %s", pretty)
	}

	y, err := f.ToYAML()
	if err != nil {
		t.Fatal(err)
	}
	var back Facts
	if err := yaml.Unmarshal(y, &back); err != nil {
		t.Fatal(err)
	}
	if back.Hostname != "web01" || back.Cpus != 4 || back.Interfaces["eth0"].MTU != 1500 {
		t.Errorf("YAML round trip mismatch: %+v", back)
	}
}
