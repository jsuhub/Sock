package tcp

import (
	"fmt"
	"log"
	"net"
	"sync"
	"time"
)

const (
	AtypIPv4 = 1
	AtypIPv6 = 4
)

type Addr []byte

func parseAddr(addr Addr) (string, int, error) {
	if len(addr) < 1 {
		return "", 0, fmt.Errorf("addr too short")
	}

	atyp := addr[0]
	switch atyp {
	case AtypIPv4:
		if len(addr) < 1+net.IPv4len+2 {
			return "", 0, fmt.Errorf("invalid IPv4 addr length")
		}
		ip := net.IP(addr[1 : 1+net.IPv4len])
		port := int(addr[len(addr)-2])<<8 | int(addr[len(addr)-1])
		return ip.String(), port, nil

	case AtypIPv6:
		if len(addr) < 1+net.IPv6len+2 {
			return "", 0, fmt.Errorf("invalid IPv6 addr length")
		}
		ip := net.IP(addr[1 : 1+net.IPv6len])
		port := int(addr[len(addr)-2])<<8 | int(addr[len(addr)-1])
		return ip.String(), port, nil

	default:
		return "", 0, fmt.Errorf("unknown address type: %d", atyp)
	}
}

func relay(left, right net.Conn) error {
	var err error
	var wg sync.WaitGroup
	var wait = 1 * time.Second
	wg.Add(1)
	read := make([]byte, 1024)
	n, err := left.Read(read)
	read = read[:n]
	_, err = right.Write(read)
	if err != nil {
		return err
	}
	log.Println("进行了一次读")
	time.Sleep(wait)
	read = make([]byte, 1024)
	n, err = right.Read(read)
	read = read[:n]
	_, err = left.Write(read)
	log.Println("进行了一次写")
	if err != nil {
		return err
	}
	return nil
}
