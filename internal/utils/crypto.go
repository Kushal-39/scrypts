package utils

import (
	"crypto/aes"
	"crypto/cipher"
	"crypto/rand"
)

func GenerateNonce(size int) ([]byte, error) {
	nonce := make([]byte, size)
	_, err := rand.Read(nonce)
	return nonce, err
}

func EncryptAESGCM(key, nonce, plaintext []byte) ([]byte, error) {
	block, err := aes.NewCipher(key)
	if err != nil {
		return nil, err
	}
	gcm, err := cipher.NewGCM(block)
	if err != nil {
		return nil, err
	}
	ciphertext := gcm.Seal(nil, nonce, plaintext, nil)
	return ciphertext, nil
}

func WrapKey(masterKey, keytoWrap []byte) ([]byte,[]byte,error){
	block,err:= aes.NewCipher(masterKey)
	if err!=nil{
		return nil,nil,err
	}
	gcm,err:=cipher.NewGCM(block)
	if err!=nil{
		return nil,nil,err
	}
	nonce:=make([]byte,gcm.NonceSize())
	if _,err:=rand.Read(nonce);err!=nil{
		return nil,nil,err
	}
	ct:= gcm.Seal(nil,nonce,keytoWrap,nil)
	return ct,nonce,nil
}

func UnwrapKey(masterKey,nonce,wrapped []byte)([]byte,error){
	block,err:=aes.NewCipher(masterKey)
	if err!=nil{
		return nil,err
	}
	gcm,err:=cipher.NewGCM(block)
	if err!=nil{
		return nil,err
	}
	pt,err:=gcm.Open(nil,nonce,wrapped,nil)
	if err!=nil{
		return nil,err
	}
	return pt,nil
}