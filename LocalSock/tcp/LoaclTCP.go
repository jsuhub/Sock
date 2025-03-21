package tcp

import (
	"LocalSock/core"
	"log"
	"net"
)

func LocalTCP(server, loacalport string, ciph core.Cipher) {
	//建立本地监听
	l, err := net.Listen("tcp", "127.0.0.1:"+loacalport)
	defer l.Close()
	if err != nil {
		log.Fatal(err)
	}

}
