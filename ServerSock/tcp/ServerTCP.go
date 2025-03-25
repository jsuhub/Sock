package tcp

import (
	"ServerSock/utils"
	"log"
	"net"
	"strconv"
)

func ServerTCP(port string, cipher utils.Cipher) {
	log.Println("0.0.0.0:" + port)
	listener, err := net.Listen("tcp", "0.0.0.0:"+port)
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
		if err != nil {
			log.Fatal("Server_traget_fail:%v", err)
		}
		tragets, port, err := parseAddr(traget[:n])

		log.Println("Server_traget_success", tragets+":"+strconv.Itoa(port))
		r, err := net.Dial("tcp", tragets+":"+strconv.Itoa(port))
		if err != nil {
			log.Fatal("Server_traget_fail:%v", err)
		}
		log.Println("建立连接成功")
		defer r.Close()
		if err = relay(lc, r); err != nil {
			log.Fatal("Server_relay_fail:%v", err)
		}

	}
}
