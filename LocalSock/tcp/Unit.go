package tcp

import (
	"io"
	"log"
	"net"
	"strconv"
	"sync"
	"time"
)

type Error byte

const (
	ErrGeneralFailure       = Error(1)
	ErrConnectionNotAllowed = Error(2)
	ErrNetworkUnreachable   = Error(3)
	ErrHostUnreachable      = Error(4)
	ErrConnectionRefused    = Error(5)
	ErrTTLExpired           = Error(6)
	ErrCommandNotSupported  = Error(7)

	ErrAddressNotSupported = Error(8)
	InfoUDPAssociate       = Error(9)
)

const (
	CmdConnect      = 1
	CmdBind         = 2
	CmdUDPAssociate = 3
)

const (
	AtypIPv4       = 1
	AtypDomainName = 3
	AtypIPv6       = 4
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

func Handshake(rw io.ReadWriter) (Addr, error) {
	// Read RFC 1928 for request and reply structure and sizes.
	buf := make([]byte, 259)
	// read VER, NMETHODS, METHODS
	if _, err := io.ReadFull(rw, buf[:2]); err != nil {
		return nil, err
	}
	nmethods := buf[1]
	if _, err := io.ReadFull(rw, buf[:nmethods]); err != nil {
		return nil, err
	}
	// write VER METHOD
	log.Println("给客户端返回验证信息")
	if _, err := rw.Write([]byte{5, 0}); err != nil {
		return nil, err
	}
	// read VER CMD RSV ATYP DST.ADDR DST.PORT
	if _, err := io.ReadFull(rw, buf[:3]); err != nil {
		return nil, err
	}
	cmd := buf[1]
	addr, err := readAddr(rw, buf)
	if err != nil {
		return nil, err
	}
	log.Println("建立TCP连接成功", addr)
	switch cmd {

	case CmdConnect:
		_, err = rw.Write([]byte{5, 0, 0, 1, 0, 0, 0, 0, 0, 0}) // SOCKS v5, reply succeeded
	}

	return addr, err // skip VER, CMD, RSV fields
}

func readAddr(r io.Reader, b []byte) (Addr, error) {
	if len(b) < 259 {
		return nil, io.ErrShortBuffer
	}
	_, err := io.ReadFull(r, b[:1]) // read 1st byte for address type
	if err != nil {
		return nil, err
	}

	switch b[0] {
	case AtypDomainName:
		_, err = io.ReadFull(r, b[1:2]) // read 2nd byte for domain length
		if err != nil {
			return nil, err
		}
		_, err = io.ReadFull(r, b[2:2+int(b[1])+2])
		return b[:1+1+int(b[1])+2], err
	case AtypIPv4:
		_, err = io.ReadFull(r, b[1:1+net.IPv4len+2])
		return b[:1+net.IPv4len+2], err
	case AtypIPv6:
		_, err = io.ReadFull(r, b[1:1+net.IPv6len+2])
		return b[:1+net.IPv6len+2], err
	}

	return nil, nil
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
