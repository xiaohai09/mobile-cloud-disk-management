package main

import (
	"fmt"
	"golang.org/x/crypto/bcrypt"
)

func main() {
	hash := "$2b$12$og/0Ll/lXUsCip/g6PcNy.vEZBHAUz8r0W3whcXCRYSGLj1Al4uj6"
	fmt.Println("Verify admin123:", bcrypt.CompareHashAndPassword([]byte(hash), []byte("admin123")) == nil)
	fmt.Println("Verify wrong:", bcrypt.CompareHashAndPassword([]byte(hash), []byte("wrong")) == nil)
}
