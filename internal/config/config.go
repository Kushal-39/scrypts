package config

import (
	"log"
	"math"
	"os"
)

var JwtSecret []byte
var MasterKey []byte

// calculateEntropy measures the Shannon entropy of a byte slice
func calculateEntropy(data []byte) float64 {
	if len(data) == 0 {
		return 0.0
	}

	freq := make(map[byte]int)
	for _, b := range data {
		freq[b]++
	}

	var entropy float64
	length := float64(len(data))
	for _, count := range freq {
		p := float64(count) / length
		if p > 0 {
			entropy -= p * math.Log2(p)
		}
	}
	return entropy
}

func Init() {
	// Validate JWT_SECRET
	s := os.Getenv("JWT_SECRET")
	if s == "" || len(s) < 32 {
		log.Fatal("FATAL: JWT_SECRET must be set and at least 32 characters long")
	}
	JwtSecret = []byte(s)

	// Validate MASTER_KEY
	mk := os.Getenv("MASTER_KEY")
	if mk == "" || len(mk) < 32 {
		log.Fatal("FATAL: MASTER_KEY must be set and at least 32 characters long")
	}
	MasterKey = []byte(mk)

	// Check entropy of secrets (minimum 4.0 bits per byte is reasonable)
	jwtEntropy := calculateEntropy(JwtSecret)
	masterEntropy := calculateEntropy(MasterKey)

	if jwtEntropy < 4.0 {
		log.Printf("WARNING: JWT_SECRET has low entropy (%.2f bits/byte). Use a stronger secret.", jwtEntropy)
	}

	if masterEntropy < 4.0 {
		log.Printf("WARNING: MASTER_KEY has low entropy (%.2f bits/byte). Use a stronger secret.", masterEntropy)
	}

	log.Println("Configuration initialized successfully")
}
