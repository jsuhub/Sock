package main

import (
	"LocalSock/core"
	"LocalSock/tcp"
	"flag"
	"log"
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
	flag.StringVar(&flags.Port, "p", "7890", "Port")
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

	if flags.TCP {
		go tcp.LocalTCP(flags.SerAddr, flags.Port, ciph)
	}
}
