// Gather facts (memory, cpus, filesystems, etc...) about the current machine.
package facts

import (
	"github.com/cloudfoundry/gosigar"

	"encoding/json"
	"net"
	"os"
	"runtime"
	"strings"
)

type Facts struct {
	Hostname    string
	Domain      string
	Fqdn        string
	Cpus        int
	Os          string
	Uptime      float64
	Memory      uint64                // In megabytes
	Swap        uint64                // In megabytes
	Interfaces  map[string]Interface  // key is nic name ex: en0, eth0
	FileSystems map[string]FileSystem // key is DeviceName
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

type FileSystem struct {
	Directory  string
	DeviceType string
	SysType    string
	Options    string
	Size       uint64 // in megabytes
}

// Gather all of the system facts available
func FindFacts() *Facts {
	f := &Facts{
		Interfaces:  map[string]Interface{},
		FileSystems: map[string]FileSystem{},
	}

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

	f.loadInterfaces()

	// sigar related items
	f.loadUptime()
	f.loadMemory()
	f.loadSwap()
	f.loadFileSystems()

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

func (f *Facts) loadInterfaces() {
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
}

func (f *Facts) loadUptime() {
	uptime := sigar.Uptime{}
	if err := uptime.Get(); err == nil {
		f.Uptime = uptime.Length
	}
}

func (f *Facts) loadMemory() {
	mem := sigar.Mem{}
	if err := mem.Get(); err == nil {
		f.Memory = mem.Total / 1024 / 1024
	}
}

func (f *Facts) loadSwap() {
	swap := sigar.Swap{}
	if err := swap.Get(); err == nil {
		f.Swap = swap.Total / 1024 / 1024
	}
}

func (f *Facts) loadFileSystems() {
	fileSystems := sigar.FileSystemList{}
	if err := fileSystems.Get(); err != nil {
		return
	}
	for _, fs := range fileSystems.List {
		filesystem := FileSystem{
			Directory:  fs.DirName,
			DeviceType: fs.TypeName,
			SysType:    fs.SysTypeName,
			Options:    fs.Options,
		}
		space := sigar.FileSystemUsage{}
		if err := space.Get(fs.DirName); err == nil {
			filesystem.Size = space.Total / 1024
		}
		f.FileSystems[fs.DevName] = filesystem
	}
}
