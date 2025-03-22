package tcp

import (
	"errors"
	"fmt"
	"io"
	"net"
	"os"
	"strconv"
	"sync"
	"time"
)

const (
	AtypIPv4 = 1
	AtypIPv6 = 4
)

type Addr []byte

func parse(s string) Addr {
	var addr Addr
	host, port, err := net.SplitHostPort(s)
	if err != nil {
		return nil
	}
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
			fmt.Println("域名解析失败:", err)
			return nil
		}
		// 找出 IPv4 地址
		for _, ip := range ips {
			if ip.To4() != nil {
				addr = make([]byte, 1+net.IPv4len+2)
				addr[0] = AtypIPv4
				copy(addr[1:], []byte(ip.String()))
				break
			}
		}
	}
	if len(host) > 255 {
		return nil
	}

	portnum, err := strconv.ParseUint(port, 10, 16)
	if err != nil {
		return nil
	}

	addr[len(addr)-2], addr[len(addr)-1] = byte(portnum>>8), byte(portnum)

	return addr
}

func relay(left, right net.Conn) error {
	var err, err1 error
	var wg sync.WaitGroup
	var wait = 5 * time.Second
	wg.Add(1)
	go func() {
		defer wg.Done()
		_, err1 = io.Copy(right, left)
		right.SetReadDeadline(time.Now().Add(wait))
	}()
	_, err = io.Copy(left, right)
	left.SetReadDeadline(time.Now().Add(wait))
	wg.Wait()
	if err1 != nil && !errors.Is(err1, os.ErrDeadlineExceeded) { // requires Go 1.15+
		return err1
	}
	if err != nil && !errors.Is(err, os.ErrDeadlineExceeded) {
		return err
	}
	return nil
}
