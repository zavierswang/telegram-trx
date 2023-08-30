package utils

import (
	"crypto/md5"
	"crypto/rand"
	"crypto/rsa"
	"crypto/x509"
	"encoding/hex"
	"encoding/pem"
	"github.com/pkg/errors"
	"os"
	"telegram-trx/pkg/core/cst"
)

func MD5(raw []byte) string {
	h := md5.New()
	h.Write(raw)
	x := h.Sum(nil)
	y := make([]byte, 32)
	hex.Encode(y, x)
	return string(y)
}

func RSAEncrypt(origData []byte) ([]byte, error) {
	publicKey, err := os.ReadFile(cst.PublicKey)
	if err != nil {
		return nil, err
	}

	block, _ := pem.Decode(publicKey) //将密钥解析成公钥实例
	if block == nil {
		return nil, errors.New("public key error")
	}
	pubInterface, err := x509.ParsePKIXPublicKey(block.Bytes)
	if err != nil {
		return nil, err
	}
	pub := pubInterface.(*rsa.PublicKey)
	return rsa.EncryptPKCS1v15(rand.Reader, pub, origData) //RSA算法加密
}

func RSADecrypt(ciphertext []byte) ([]byte, error) {
	privateKey, err := os.ReadFile(cst.PrivateKey)
	if err != nil {
		return nil, err
	}
	block, _ := pem.Decode(privateKey)
	if block == nil {
		return nil, errors.New("private key error!")
	}
	priv, err := x509.ParsePKCS1PrivateKey(block.Bytes)
	if err != nil {
		return nil, err
	}
	return rsa.DecryptPKCS1v15(rand.Reader, priv, ciphertext) //RSA算法解密
}
