package utils

import (
	"crypto/sha256"
	"net"

	"github.com/shadowsocks/go-shadowsocks2/shadowaead"
	"golang.org/x/crypto/pbkdf2"
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
	var chipher = aeadList[method]
	passByte := []byte(password)
	var key = make([]byte, chipher.KeySize)
	if len(password) < chipher.KeySize {
		//通过使用KDF把密码加密成密钥
		key = increasePasswd(passByte)
	} else {
		key = passByte[:chipher.KeySize]
	}
	aead, err := chipher.New(key)
	return &aeadCipher{aead}, err
	//&相当于NEW
}

func increasePasswd(password []byte) []byte {
	var salt = []byte{
		0x7a, 0x31, 0x9c, 0x8d, 0xaf, 0x43, 0x76, 0x55,
		0x88, 0xe2, 0x0f, 0x99, 0x34, 0xc1, 0x6a, 0x12,
		0x58, 0x6d, 0xa9, 0x7b, 0x40, 0xcd, 0x28, 0xee,
		0x91, 0x2a, 0x6b, 0xfa, 0x13, 0xf0, 0x2d, 0xc3,
	}
	iterations := 100000
	keyLen := 32 // 256 位
	key := pbkdf2.Key(password, salt, iterations, keyLen, sha256.New)
	return key
}
