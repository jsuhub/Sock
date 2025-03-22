package main

import (
	"LocalSock/core"
	"LocalSock/tcp"
	"flag"
	"fmt"
	"log"
	"os"
	"os/signal"
	"syscall"
)

func main() {

	var flags struct {
		Cipher  string
		Port    string
		Key     string
		Keygen  int
		SerAddr string
		TCP     bool
		UDP     bool
		Socks   bool
	}

	flag.StringVar(&flags.Key, "key", "", "Password")
	flag.StringVar(&flags.Port, "p", "8888", "Port")
	flag.StringVar(&flags.Cipher, "cipher", "AEAD_AES_256_GCM", "cipher")
	flag.BoolVar(&flags.TCP, "tcp", false, "tcp")
	flag.BoolVar(&flags.UDP, "udp", false, "udp")
	flag.BoolVar(&flags.Socks, "socks", false, "socks")

	flag.Parse()
	if flags.TCP == false && flags.UDP == false && flags.Socks == false {
		flag.Usage()
	}

	//开始加密，先生成密钥 ——利用密码
	ciph, err := core.NewCipher(flags.Cipher, flags.Key)
	if err != nil {
		log.Fatal("生成密钥错误%v", err)
	}
	log.Println("密钥生成成功")
	if flags.TCP {

		//SerAddr 表示服务端目标 Port 表示客户端监听的端口 ciph表示加密
		go tcp.LocalTCP(flags.SerAddr, flags.Port, ciph)

	}
	sigCh := make(chan os.Signal, 1)
	signal.Notify(sigCh, syscall.SIGINT, syscall.SIGTERM)

	fmt.Println("服务已启动，按 Ctrl+C 退出...")
	<-sigCh

	fmt.Println("收到退出信号，开始清理")

}
