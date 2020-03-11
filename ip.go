package snow

import (
	"net"
	"strings"
)

// InferHostIPv4 infers the host IP address.
func InferHostIPv4(ifaceName string) string {
	if ifaceName != "" {
		ipList, _ := HostIPv4(IfaceNameMatchExactly, ifaceName)
		if len(ipList) > 0 {
			return ipList[0]
		}

		ipList, _ = HostIPv4(IfaceNameMatchPrefix, ifaceName)
		if len(ipList) > 0 {
			return ipList[0]
		}
	}

	ipList, _ := HostIPv4(IfaceNameMatchPrefix, "en", "eth")

	if len(ipList) == 0 {
		ipList, _ = HostIPv4(IfaceNameMatchPrefix)
	}

	if len(ipList) == 1 { // nolint gomnd
		return ipList[0]
	}

	ipList2, _ := HostIPv4(IfaceNameMatchPrefix, "en0", "eth0")
	if len(ipList2) > 0 {
		return ipList2[0]
	}

	return ipList[0]
}

// IfaceNameMatchMode defines the mode for IfaceName matching
type IfaceNameMatchMode int

const (
	// IfaceNameMatchPrefix matches iface name in prefix mode
	IfaceNameMatchPrefix IfaceNameMatchMode = iota
	// IfaceNameMatchExactly matches iface name in exactly same mode
	IfaceNameMatchExactly
)

// HostIPv4 根据 primaryIfaceName 确定的名字，返回主IP PrimaryIP，以及以空格分隔的本机IP列表 ipList
// PrimaryIfaceName 表示主网卡的名称，用于获取主IP(v4)，不设置时，从eth0(linux), en0(darwin)，或者第一个ip v4的地址
// eg.  HostIP("eth0", "en0") // nolint
func HostIPv4(mode IfaceNameMatchMode, primaryIfaceNames ...string) (ipList []string, err error) {
	ips, err := ListIfaces(IPv4)
	if err != nil {
		return
	}

	ipList = make([]string, 0, len(ips))

	for _, addr := range ips {
		if mode.match(primaryIfaceNames, addr.IfaceName) {
			ipList = append(ipList, addr.IP.String())
		}
	}

	return
}

func (mode IfaceNameMatchMode) match(ifaceNames []string, name string) bool {
	if len(ifaceNames) == 0 {
		return true
	}

	for _, ifaceName := range ifaceNames {
		switch mode {
		case IfaceNameMatchPrefix:
			if strings.HasPrefix(name, ifaceName) {
				return true
			}
		case IfaceNameMatchExactly:
			if name == ifaceName {
				return true
			}
		}
	}

	return false
}

// ListMode defines the mode for iface listing
type ListMode int

const (
	// IPv4v6 list all ipv4 and ipv6
	IPv4v6 ListMode = iota
	// IPv4 list only all ipv4
	IPv4
	// IPv6 list only all ipv6
	IPv6
)

// Iface 表示一个IP地址和网卡名称的结构
type Iface struct {
	IP        net.IP
	IfaceName string
}

// ListIfaces 根据mode 列出本机所有IP和网卡名称
func ListIfaces(mode ...ListMode) ([]Iface, error) {
	ret := make([]Iface, 0)
	list, err := net.Interfaces()

	if err != nil {
		return ret, err
	}

	modeMap := makeModeMap(mode)

	for _, iface := range list {
		addrs, err := iface.Addrs()
		if err != nil {
			return ret, err
		}

		for _, addr := range addrs {
			ipnet, ok := addr.(*net.IPNet)
			if !ok || ipnet.IP.IsLoopback() {
				continue
			}

			if matchesMode(modeMap, ipnet) {
				ret = append(ret, Iface{
					IP:        ipnet.IP,
					IfaceName: iface.Name,
				})
			}
		}
	}

	return ret, nil
}

func matchesMode(modeMap map[ListMode]bool, ipnet *net.IPNet) bool {
	if _, all := modeMap[IPv4v6]; all {
		return true
	}

	if _, v4 := modeMap[IPv4]; v4 {
		return ipnet.IP.To4() != nil
	}

	if _, v6 := modeMap[IPv6]; v6 {
		return ipnet.IP.To16() != nil
	}

	return false
}

func makeModeMap(mode []ListMode) map[ListMode]bool {
	modeMap := make(map[ListMode]bool)

	for _, m := range mode {
		modeMap[m] = true
	}

	if len(modeMap) == 0 {
		modeMap[IPv4v6] = true
	}

	return modeMap
}
