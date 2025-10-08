package config

import (
	"crypto/sha256"
	"os"
)

var JwtSecret []byte
var MasterKey []byte

func Init() {
	s := os.Getenv("JWT_SECRET")
	JwtSecret = []byte(s)

	mk := os.Getenv("MASTER_KEY")
	if mk != "" {
		MasterKey = []byte(mk)
	} else {
		sum := sha256.Sum256(JwtSecret)
		MasterKey = sum[:]
	}
}
