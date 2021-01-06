package main

import (
	"bytes"
	"crypto/des"
	"encoding/base64"
	"encoding/json"
	"errors"
	"fmt"
	"log"
	"strconv"
	"strings"
)

type verificationMailInfo struct {
	Uid   string `json:"u"`
	Email string `json:"e"`
	Time  int64  `json:"t"`
}

func main() {
	key, err := base64.StdEncoding.DecodeString("秘钥马赛克")
	if err != nil {
		log.Fatalln(err)
	}
	fmt.Println(string(key))

	jsonString := "{\"u\":\"%s\",\"e\":\"%s\",\"t\":%d}"
	jsonString = fmt.Sprintf(jsonString, "4315108321", "402374739@qq.com", 1608194547)
	mIn := []byte(jsonString)
	log.Println(string(mIn))
	/*mIn = []byte("{\"u\":\"4315108321\",\"e\":\"402374739@qq.com\",\"t\":1608192946}")
	log.Println(string(mIn), len(mIn))*/
	log.Println("------------------ 加密 --------------------")
	s, e := TripleEcbDesEncrypt(mIn, key)
	fmt.Println(s)
	fmt.Println(string(s))
	data := base64UrlSafeEncode(s)
	fmt.Println(data)

	log.Println("------------------ 解密 --------------------")
	data = "Ita9d24gj8818bwbH4Pp44WzqP7TXUGoSYiD-Ga6Bdg2-vAbphC1E-u4mP6PXxy9KVrcxVMa6RQex0RJOkJclQ2"
	safeurlB, err := base64URLDecode(data)
	fmt.Println(safeurlB)
	s, e = TripleEcbDesDecrypt(safeurlB, key)
	fmt.Println(s, e)
	m := &verificationMailInfo{}
	_ = json.Unmarshal(s, &m)
	fmt.Println(m)
	fmt.Println(m.Uid)
}
func paddingPKCS7(res []byte) []byte {
	paddingChar := 8 - (len(res) % 8)
	ch := rune(paddingChar)
	return []byte(string(res) + strings.Repeat(string(ch), paddingChar))
}

//[golang ECB 3DES Encrypt]
func TripleEcbDesEncrypt(origData, key []byte) ([]byte, error) {
	tkey := make([]byte, 24, 24)
	copy(tkey, key)
	k1 := tkey[:8]
	k2 := tkey[8:16]
	k3 := tkey[16:]

	block, err := des.NewCipher(k1)
	if err != nil {
		return nil, err
	}
	bs := block.BlockSize()
	origData = PKCS5Padding(origData, bs)

	buf1, err := encrypt(origData, k1)
	if err != nil {
		return nil, err
	}
	buf2, err := decrypt(buf1, k2)
	if err != nil {
		return nil, err
	}
	out, err := encrypt(buf2, k3)
	if err != nil {
		return nil, err
	}
	return out, nil
}

//[golang ECB 3DES Decrypt]
func TripleEcbDesDecrypt(crypted, key []byte) ([]byte, error) {
	tkey := make([]byte, 24, 24)
	copy(tkey, key)
	k1 := tkey[:8]
	k2 := tkey[8:16]
	k3 := tkey[16:]
	buf1, err := decrypt(crypted, k3)
	if err != nil {
		return nil, err
	}
	buf2, err := encrypt(buf1, k2)
	if err != nil {
		return nil, err
	}
	out, err := decrypt(buf2, k1)
	if err != nil {
		return nil, err
	}
	out = PKCS5Unpadding(out)
	return out, nil
}

//ECB PKCS5Padding
func PKCS5Padding(ciphertext []byte, blockSize int) []byte {
	padding := blockSize - len(ciphertext)%blockSize
	padtext := bytes.Repeat([]byte{byte(padding)}, padding)
	return append(ciphertext, padtext...)
}

//ECB PKCS5Unpadding
func PKCS5Unpadding(origData []byte) []byte {
	length := len(origData)
	unpadding := int(origData[length-1])
	return origData[:(length - unpadding)]
}

//Des加密
func encrypt(origData, key []byte) ([]byte, error) {
	if len(origData) < 1 || len(key) < 1 {
		return nil, errors.New("wrong data or key")
	}
	block, err := des.NewCipher(key)
	if err != nil {
		return nil, err
	}
	bs := block.BlockSize()
	if len(origData)%bs != 0 {
		return nil, errors.New("wrong padding")
	}
	out := make([]byte, len(origData))
	dst := out
	for len(origData) > 0 {
		block.Encrypt(dst, origData[:bs])
		origData = origData[bs:]
		dst = dst[bs:]
	}
	return out, nil
}

//Des解密
func decrypt(crypted, key []byte) ([]byte, error) {
	if len(crypted) < 1 || len(key) < 1 {
		return nil, errors.New("wrong data or key")
	}
	block, err := des.NewCipher(key)
	if err != nil {
		return nil, err
	}
	out := make([]byte, len(crypted))
	dst := out
	bs := block.BlockSize()
	if len(crypted)%bs != 0 {
		return nil, errors.New("wrong crypted size")
	}

	for len(crypted) > 0 {
		block.Decrypt(dst, crypted[:bs])
		crypted = crypted[bs:]
		dst = dst[bs:]
	}

	return out, nil
}

func base64UrlSafeEncode(source []byte) string {
	bytearr := base64.StdEncoding.EncodeToString(source)
	noOfEq := strings.Count(bytearr, "=")
	safeurl := strings.Replace(bytearr, "=", "", -1)
	fmt.Println(bytearr, noOfEq)
	safeurl += strconv.Itoa(noOfEq)
	safeurl = strings.Replace(safeurl, "+", "-", -1)
	safeurl = strings.Replace(safeurl, "/", "_", -1)
	return safeurl
}

func base64URLDecode(data string) ([]byte, error) {
	data = strings.Replace(data, "_", "/", -1)
	data = strings.Replace(data, "-", "+", -1)
	noOfEqInt, _ := strconv.Atoi(data[len(data)-1:])
	data = data[:len(data)-1]
	if noOfEqInt > 0 {
		data += strings.Repeat("=", noOfEqInt)
	}
	fmt.Println(noOfEqInt, data)
	decodeBytes, err := base64.StdEncoding.DecodeString(data)
	if err != nil {
		log.Fatalln(err, "asdfafdasfdadsf")
	}
	return decodeBytes, nil
}
