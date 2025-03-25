package tcp

import (
	"log"
	"net"
	"strconv"
)

const (
	AtypIPv4 = 1
	AtypIPv6 = 4
)

type Addr []byte

func parse(host string) Addr {
	var addr Addr

	if ip := net.ParseIP(host); ip != nil {

		if ip4 := ip.To4(); ip4 != nil {
			addr = make([]byte, 1+net.IPv4len+2)
			addr[0] = AtypIPv4
			copy(addr[1:], ip4)
		} else {

			addr = make([]byte, 1+net.IPv6len+2)
			addr[0] = AtypIPv6
			copy(addr[1:], ip)
		}
	} else {

		ips, err := net.LookupIP(host)
		if err != nil {
			log.Println("域名解析失败:", err)
			return nil
		}
		// 找出 IPv4 地址
		for _, ip := range ips {
			if ip.To4() != nil {
				addr = make([]byte, 1+net.IPv4len+2)
				addr[0] = AtypIPv4
				copy(addr[1:], []byte(ip.To4()))
				break
			}
		}
	}
	if len(host) > 255 {
		return nil
	}
	portnum, err := strconv.ParseUint("0080", 10, 16)
	if err != nil {
		return nil
	}

	addr[len(addr)-2], addr[len(addr)-1] = byte(portnum>>8), byte(portnum)

	return addr
}
