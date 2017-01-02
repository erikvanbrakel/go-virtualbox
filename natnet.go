package virtualbox

import (
	"bufio"
	"net"
	"strings"
)

// A NATNet defines a NAT network.
type NATNet struct {
	Name           string
	IPv4           net.IPNet
	IPv6           net.IPNet
	IPv6Enabled    bool
	DHCPEnabled    bool
	Enabled        bool
}

// NATNets gets all NAT networks in a  map keyed by NATNet.Name.
func NATNets() (map[string]NATNet, error) {
	out, err := vbmOut("list", "natnets")
	if err != nil {
		return nil, err
	}
	s := bufio.NewScanner(strings.NewReader(out))
	m := map[string]NATNet{}
	n := NATNet{}
	for s.Scan() {
		line := s.Text()
		if line == "" {
			m[n.Name] = n
			n = NATNet{}
			continue
		}
		res := reColonLine.FindStringSubmatch(line)
		if res == nil {
			continue
		}
		switch key, val := res[1], res[2]; key {
		case "NetworkName":
			n.Name = val
		case "IP":
			n.IPv4.IP = net.ParseIP(val)
		case "Network":
			_, ipnet, err := net.ParseCIDR(val)
			if err != nil {
				return nil, err
			}
			n.IPv4.Mask = ipnet.Mask
		case "IPv6 Prefix":
			if val == "" {
				continue
			}
			_, ipnet, err := net.ParseCIDR(val)
			if err != nil {
				return nil, err
			}
			n.IPv6.Mask = ipnet.Mask
		case "IPv6 Enabled":
			n.IPv6Enabled = (val == "Yes")
		case "DHCP Enabled":
			n.DHCPEnabled = (val == "Yes")
		case "Enabled":
			n.Enabled = (val == "Yes")
		}
	}
	if err := s.Err(); err != nil {
		return nil, err
	}
	return m, nil
}

func CreateNATNet(name, cidr_block string, supports_dhcp, supports_ipv6 bool) (*NATNet, error) {

	args := []string{
		"natnetwork","add",
		"--netname",name,
		"--network",cidr_block,
	}

	if supports_dhcp {
		args = append(args, "--dhcp", "on")
	} else {
		args = append(args, "--dhcp", "off")
	}

	if supports_ipv6 {
		args = append(args, "--ipv6", "on")
	} else {
		args = append(args, "--ipv6", "off")
	}

	if err := vbm(args...); err != nil {
		return nil, err
	}

	nats, err := NATNets()

	if err != nil {
		return nil, err
	}

	nat, _ := nats[name]

	return &nat, nil
}

func DeleteNATNet(name string) error {
	args := []string{
		"natnetwork", "remove",
		"--netname", name,
	}

	if err := vbm(args...); err != nil {
		return err
	}

	return nil
}

func (nat *NATNet) Update() {

	args := []string {
		"natnetwork", "modify",
		"--netname", nat.Name,
	}

	if nat.DHCPEnabled {
		args = append(args, "--dhcp", "on")
	} else {
		args = append(args, "--dhcp", "off")
	}

	if nat.IPv6Enabled {
		args = append(args, "--ipv6", "on")
	} else {
		args = append(args, "--ipv6", "off")
	}
}
