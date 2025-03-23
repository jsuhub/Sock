package tcp

import (
	"ServerSock/utils"
	"log"
	"net"
	"strconv"
)

func ServerTCP(port string, cipher utils.Cipher) {
	listener, err := net.Listen("tcp", "127.0.0.1:"+port)
	defer listener.Close()
	if err != nil {
		log.Fatal("Server_listing_fail:%v", err)
	}
	log.Println("Server_listing_success")
	for {
		l, err := listener.Accept()
		lc := cipher.StreamConn(l)
		if err != nil {
			log.Fatal("Server_conn_fail:%v", err)
		}
		traget := make([]byte, 64)
		n, err := lc.Read(traget)
		tragets, port, err := parseAddr(traget[:n])
		if err != nil {
			log.Fatal("Server_traget_fail:%v", err)
		}
		log.Println("Server_traget_success", port, tragets)
		r, err := net.Dial("tcp", tragets+":"+strconv.Itoa(port))
		log.Println("建立连接成功")
		if err != nil {
			log.Fatal("Server_traget_fail:%v", err)
		}
		defer r.Close()
		if err = relay(lc, r); err != nil {
			log.Fatal("Server_relay_fail:%v", err)
		}

	}
}
