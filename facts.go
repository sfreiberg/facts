package facts

import (
	"encoding/json"
	"net"
	"os"
	"runtime"
	"strings"
)

type Facts struct {
	Hostname   string
	Domain     string
	Fqdn       string
	Cpus       int
	Os         string
	Interfaces map[string]Interface
}

type Interface struct {
	Index        int
	MTU          int
	Name         string
	HardwareAddr string
	Flags        []string
	Addresses    []Address
}

type Address struct {
	Network   string
	Address   string
	IPNetwork string
}

// Gather all of the system facts available
func FindFacts() *Facts {
	f := &Facts{Interfaces: map[string]Interface{}}

	// get the domain info
	if fqdn, err := os.Hostname(); err == nil {
		f.Fqdn = fqdn
		a := strings.SplitN(fqdn, ".", 2)
		if len(a) == 2 {
			f.Hostname = a[0]
			f.Domain = a[1]
		} else if len(a) == 1 {
			f.Hostname = a[0]
		}
	}

	f.Cpus = runtime.NumCPU()

	f.Os = getOs()

	ifaces, err := net.Interfaces()
	if err == nil {
		for _, iface := range ifaces {
			i := Interface{
				Index:        iface.Index,
				MTU:          iface.MTU,
				Name:         iface.Name,
				HardwareAddr: iface.HardwareAddr.String(),
				Flags:        strings.Split(iface.Flags.String(), "|"),
			}

			addrs, err := iface.Addrs()
			if err == nil {
				for _, addr := range addrs {
					ip, ipnet, err := net.ParseCIDR(addr.String())
					if err == nil {
						a := Address{
							Network:   addr.Network(),
							Address:   ip.String(),
							IPNetwork: ipnet.String(),
						}
						i.Addresses = append(i.Addresses, a)
					}
				}

			}

			f.Interfaces[iface.Name] = i
		}
	}

	return f
}

// Return facts as a JSON document
func (f *Facts) ToJson() ([]byte, error) {
	return json.Marshal(f)
}

// Return facts as a JSON document with newlines and indentation added.
func (f *Facts) ToPrettyJson() ([]byte, error) {
	return json.MarshalIndent(f, "", "  ")
}
