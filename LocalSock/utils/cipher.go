package utils

import (
	"crypto/md5"
	"net"

	"github.com/shadowsocks/go-shadowsocks2/shadowaead"
)

/*
实现了KDF对密码生成密钥+AES对称加密通过调用其他函数  不过只实现了AES-256-GCM的加密
*/

type Cipher interface {
	StreamConnCipher
	PacketConnCipher
}

type StreamConnCipher interface {
	StreamConn(net.Conn) net.Conn
}

type PacketConnCipher interface {
	PacketConn(net.PacketConn) net.PacketConn
}

// aed 表示加密算法需要的结构体
var aeadList = map[string]struct {
	KeySize int
	New     func([]byte) (shadowaead.Cipher, error)
	//相当于在函数体里面定义了一个类方法，指定了返回值和输入
}{
	"AEAD_AES_256_GCM": {32, shadowaead.AESGCM},
	//关键shadowaead.AESGCM是一个可以利用32位加密返回一个shadowaead.Cipher
}

type aeadCipher struct{ shadowaead.Cipher }

func (aead *aeadCipher) StreamConn(c net.Conn) net.Conn { return shadowaead.NewConn(c, aead) }

//返回一个加密过后的通道StreamConnCipher实现StreamConn,类方法

func (aead *aeadCipher) PacketConn(c net.PacketConn) net.PacketConn {
	return shadowaead.NewPacketConn(c, aead)
}

func NewCipher(method string, password string) (cipher Cipher, err error) {
	var choice = aeadList[method]
	var key = make([]byte, choice.KeySize)
	key = kdf(password, choice.KeySize)
	aead, err := choice.New(key)
	return &aeadCipher{aead}, err
	//&相当于NEW
}

func kdf(password string, keyLen int) []byte {
	var b, prev []byte
	h := md5.New()
	//生成md5加密器
	for len(b) < keyLen {
		h.Write(prev)
		//把prev写入
		h.Write([]byte(password))
		//形成prev+password的格式
		b = h.Sum(b)
		//表示使用这两个加密并且添加进去appen[b,preb...]
		prev = b[len(b)-h.Size():]
		//prev每次都是加密后添加的部分
		h.Reset()
	}
	return b[:keyLen]
	//表示返回的长度符合
}
