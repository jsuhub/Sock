package tcp

import (
	"LocalSock/core"
	"bufio"
	"bytes"
	"log"
	"net"
	"net/http"
)

func LocalTCP(server, localport string, ciph core.Cipher) {
	//建立本地监听
	l, err := net.Listen("tcp", "127.0.0.1:"+localport)
	if err != nil {
		log.Fatal(err)
	}
	log.Println("Local TCP server listening on " + localport)
	defer l.Close()
	for {
		conn, err := l.Accept()
		reader := bufio.NewReader(conn)
		peeked, err := reader.Peek(1024) // 看前 1024 字节
		if err != nil {
			log.Println("Peek failed:", err)
			return
		}
		read := bufio.NewReader(bytes.NewReader(peeked))
		//通过生成一个新的缓冲区并且利用peek只是观看缓冲区里面的数据但是并没有推进的效果为io.copy提供不变
		req, err := http.ReadRequest(read)
		if err != nil {
			log.Println("Read request failed:", err)
		}
		traget := parse(req.Host)
		if err != nil {
			continue
		}
		go func() {
			defer conn.Close()
			r, err := net.Dial("tcp", server)
			defer r.Close()
			rc := ciph.StreamConn(r)
			if _, err = rc.Write([]byte(traget)); err != nil {
				log.Println("tcp conn write err: %v", err)
				return
			}
			if err = relay(conn, rc); err != nil {
				log.Println("tcp relay err: %v", err)
			}
		}()

	}
}
 