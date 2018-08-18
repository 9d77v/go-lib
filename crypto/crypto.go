package crypto

import (
	"bytes"
	"crypto/aes"
	"crypto/cipher"
	"crypto/des"
	"crypto/hmac"
	"crypto/md5"
	"crypto/rand"
	"crypto/rsa"
	"crypto/sha1"
	"crypto/sha512"
	"crypto/x509"
	"encoding/base64"
	"encoding/hex"
	"encoding/pem"
	"errors"
	"fmt"
	"io"
	"io/ioutil"
	"log"
	"sort"
	"strconv"
	"strings"
)

//MD5 获取md5字符串
func MD5(s string, isUpper bool) string {
	h := md5.New()
	h.Write([]byte(s))
	rs := hex.EncodeToString(h.Sum(nil))
	if isUpper {
		return strings.ToUpper(rs)
	}
	return rs
}

//SHA1 获取sha1字符串
func SHA1(s string, isUpper bool) string {
	h := sha1.New()
	h.Write([]byte(s))
	rs := hex.EncodeToString(h.Sum(nil))
	if isUpper {
		return strings.ToUpper(rs)
	}
	return rs
}

//SHA512 获取sha512字符串
func SHA512(s string) string {
	h := sha512.New()
	h.Write([]byte(s))
	rs := hex.EncodeToString(h.Sum(nil))
	return strings.ToUpper(rs)
}

//Param 参数对象
type Param struct {
	Key   string
	Value string
}

//SignHMAC 生成HMAC签名
func SignHMAC(params []*Param, secret string) string {
	// 第一步：检查参数是否已经排序
	sort.Slice(params, func(i, j int) bool {
		return params[i].Key <= params[j].Key
	})
	// 第二步：把所有参数名和参数值串在一起
	queryString := ""
	for _, v := range params {
		queryString += v.Key + v.Value
	}
	// 第三步：使用HMAC加密
	key := []byte(secret)
	mac := hmac.New(md5.New, key)
	mac.Write([]byte(queryString))
	// 第四步：把二进制转化为大写的十六进制
	return strings.ToUpper(hex.EncodeToString(mac.Sum(nil)))
}

//Base64Encode base64加密
func Base64Encode(src []byte) string {
	return base64.StdEncoding.EncodeToString(src)
}

//Base64Decode base64解密
func Base64Decode(src string) string {
	code, err := base64.StdEncoding.DecodeString(src)
	if err != nil {
		fmt.Println("Base64解码失败!" + err.Error())
	}
	return string(code)
}

//PKCS5Padding ...
func PKCS5Padding(ciphertext []byte, blockSize int) []byte {
	padding := blockSize - len(ciphertext)%blockSize
	padtext := bytes.Repeat([]byte{byte(padding)}, padding)
	return append(ciphertext, padtext...)
}

//PKCS5UnPadding ...
func PKCS5UnPadding(origData []byte) []byte {
	length := len(origData)
	unpadding := int(origData[length-1])
	return origData[:(length - unpadding)]
}

//DesEncrypt Des加密
func DesEncrypt(data, key []byte) ([]byte, error) {
	block, err := des.NewCipher(key)
	if err != nil {
		return nil, err
	}
	bs := block.BlockSize()
	data = PKCS5Padding(data, bs)
	if len(data)%bs != 0 {
		return nil, errors.New("Need a multiple of the blocksize")
	}
	out := make([]byte, len(data))
	dst := out
	for len(data) > 0 {
		block.Encrypt(dst, data[:bs])
		data = data[bs:]
		dst = dst[bs:]
	}
	return out, err
}

//DesDecrypt Des解密
func DesDecrypt(data, key []byte) ([]byte, error) {
	block, err := des.NewCipher(key)
	if err != nil {
		return nil, err
	}
	bs := block.BlockSize()
	if len(data)%bs != 0 {
		return nil, errors.New("crypto/cipher: input not full blocks")
	}
	out := make([]byte, len(data))
	dst := out
	for len(data) > 0 {
		block.Decrypt(dst, data[:bs])
		data = data[bs:]
		dst = dst[bs:]
	}
	out = PKCS5UnPadding(out)
	return out, err
}

//RSAEncrypt RSA加密
func RSAEncrypt(data, publicKey []byte) ([]byte, error) {
	block, _ := pem.Decode(publicKey)
	if block == nil {
		return nil, errors.New("public key error")
	}
	pubInterface, err := x509.ParsePKIXPublicKey(block.Bytes)
	if err != nil {
		return nil, err
	}
	return rsa.EncryptPKCS1v15(rand.Reader, pubInterface.(*rsa.PublicKey), data)
}

//RSADecrypt RSA解密
func RSADecrypt(data, privateKey []byte) ([]byte, error) {
	block, _ := pem.Decode(privateKey)
	if block == nil {
		log.Println(block)
		return nil, errors.New("private key error")
	}
	priv, err := x509.ParsePKCS1PrivateKey(block.Bytes)
	if err != nil {
		log.Println(err)
		return nil, err
	}
	return rsa.DecryptPKCS1v15(rand.Reader, priv, data)
}

//RSAPublicKey ...
func RSAPublicKey(path string) (string, error) {
	publicKey, err := ioutil.ReadFile(path)
	if err != nil {
		return "", err
	}
	block, _ := pem.Decode(publicKey)
	if block == nil {
		return "", errors.New("public key error")
	}
	pubInterface, err := x509.ParsePKIXPublicKey(block.Bytes)
	if err != nil {
		return "", err
	}
	n := Base64Encode(pubInterface.(*rsa.PublicKey).N.Bytes())
	e := strconv.FormatInt(int64(pubInterface.(*rsa.PublicKey).E), 16)
	return n + "|" + e, err
}

//RSAPublicKeyOrigin ...
func RSAPublicKeyOrigin(path string) (string, error) {
	publicKey, err := ioutil.ReadFile(path)
	if err != nil {
		return "", err
	}
	block, _ := pem.Decode(publicKey)
	if block == nil {
		return "", errors.New("public key error")
	}
	return Base64Encode(block.Bytes), err
}

//AesDecrypt AES-CBC解密,PKCS#7,传入密文和密钥，[]byte
func AesDecrypt(src, key []byte) (dst []byte, err error) {
	block, err := aes.NewCipher(key)
	if err != nil {
		return nil, err
	}
	iv := make([]byte, aes.BlockSize)
	if _, err := io.ReadFull(rand.Reader, iv); err != nil {
		return nil, err
	}
	dst = make([]byte, len(src))
	cipher.NewCBCDecrypter(block, iv).CryptBlocks(dst, src)

	return PKCS7UnPad(dst), nil
}

//PKCS7UnPad PKSC#7解包
func PKCS7UnPad(msg []byte) []byte {
	length := len(msg)
	padlen := int(msg[length-1])
	return msg[:length-padlen]
}

//AesEncrypt AES-CBC加密+PKCS#7打包，传入明文和密钥
func AesEncrypt(src []byte, key []byte) ([]byte, error) {
	k := len(key)
	if len(src)%k != 0 {
		src = PKCS7Pad(src, k)
	}

	block, err := aes.NewCipher(key)
	if err != nil {
		return nil, err
	}

	iv := make([]byte, aes.BlockSize)
	if _, err := io.ReadFull(rand.Reader, iv); err != nil {
		return nil, err
	}

	dst := make([]byte, len(src))
	cipher.NewCBCEncrypter(block, iv).CryptBlocks(dst, src)

	return dst, nil
}

// PKCS7Pad PKCS#7打包
func PKCS7Pad(msg []byte, blockSize int) []byte {
	if blockSize < 1<<1 || blockSize >= 1<<8 {
		panic("unsupported block size")
	}
	padlen := blockSize - len(msg)%blockSize
	padding := bytes.Repeat([]byte{byte(padlen)}, padlen)
	return append(msg, padding...)
}
