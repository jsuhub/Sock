package tcp

import (
	"LocalSock/utils"
	"bufio"
	"bytes"
	"fmt"
	"log"
	"net"
	"net/http"
	"time"
)

func LocalTCP(server, localport string, ciph utils.Cipher) {
	//建立本地监听
	l, err := net.Listen("tcp", "127.0.0.1:"+localport)
	if err != nil {
		log.Fatal(err)
	}
	log.Println("Local tcp server listening on " + localport)
	defer l.Close()
	for {
		conn, err := l.Accept()
		read := make([]byte, 1024)
		n, err := conn.Read(read)
		if err != nil {
			log.Fatal(err)
		}
		read = read[:n]
		reader := bufio.NewReader(bytes.NewReader(read))
		req, err := http.ReadRequest(reader)
		if err != nil {
			fmt.Println("Error reading request:", err)
			return
		}
		traget_s := req.Host
		traget := parse(traget_s)
		if err != nil {
			continue
		}

		go func() {
			defer conn.Close()
			log.Println(traget)
			log.Println("Local tcp connect to", server)
			r, err := net.Dial("tcp", server)

			defer r.Close()
			log.Println("Dial ok")
			rc := ciph.StreamConn(r)
			if _, err = rc.Write([]byte(traget)); err != nil {
				log.Println("tcp conn write err: %v", err)
				return
			}
			rc.Write(read)
			
			log.Println("client write ok")
			time.Sleep(5 * time.Second)
			read = make([]byte, 1024)
			n, err := rc.Read(read)
			if err != nil {
				log.Fatal("tcp conn read err: %v", err)
			}
			read = read[:n]
			log.Println("client read ok")
			conn.Write(read)

		}()

	}
}
