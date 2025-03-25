package tcp

import (
	"LocalSock/utils"
	"log"
	"net"
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

		traget, err := Handshake(conn)
		if err != nil {
			log.Fatal(err)
		}
		//traget = parse(string(traget))
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
			log.Println("proxy %s <-> %s <-> %s", rc.RemoteAddr(), server, traget)
			if err = relay(conn, rc); err != nil {
				log.Fatal("relay error: %v", err)
			}

			//time.Sleep(5 * time.Second)
			//read := make([]byte, 1024)
			//n, err := rc.Read(read)
			//if err != nil {
			//	log.Fatal("tcp conn read err: %v", err)
			//}
			//read = read[:n]
			//log.Println("client read ok")
			//conn.Write(read)

		}()

	}
}
