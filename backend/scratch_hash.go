package main

import (
	"fmt"
	"golang.org/x/crypto/bcrypt"
)

func main() {
	bytes, _ := bcrypt.GenerateFromPassword([]byte("password"), 10)
	fmt.Println("Correct hash for 'password':")
	fmt.Println(string(bytes))
}
