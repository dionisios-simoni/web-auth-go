package main

import (
	"crypto/hmac"
	"crypto/sha512"
	"fmt"
	"golang.org/x/crypto/bcrypt"
	"log"
)

// key is a secret between the two communicating services
var key []byte

func main() {

	for i := 0; i < 65; i++ {
		key = append(key, byte(i))
	}

	// a very strong password is encrypted with bcrypt hash
	pass := []byte("12345")
	hashedPass, err := hashPassword(pass)
	if err != nil {
		log.Fatal(err)
	}

	// checking if a given password results to the real password
	err = comparePassword(pass, hashedPass)
	if err != nil {
		log.Fatal(err)
	}

	fmt.Println("logged in, validating message authenticity")

	// A very secure message is signed with the secret key
	msg := []byte("this is a very secret message.")
	b, err := sign(msg)
	if err != nil {
		log.Fatal(err)
	}

	// A safe message should have a valid signature based on the secret key
	safe, err := checkSignature(msg, b)
	if err != nil {
		log.Fatal()
	}

	if safe {
		fmt.Println("message was safe!")
	} else {
		fmt.Println("something is fishy!!!")
	}
}

func hashPassword(password []byte) ([]byte, error) {
	bs, err := bcrypt.GenerateFromPassword(password, bcrypt.DefaultCost)
	if err != nil {
		return nil, err
	}

	return bs, nil
}

func comparePassword(password, hashedPassword []byte) error {
	return bcrypt.CompareHashAndPassword(hashedPassword, password)
}

func sign(msg []byte) ([]byte, error) {
	h := hmac.New(sha512.New, key)
	_, err := h.Write(msg)
	if err != nil {
		return nil, err
	}

	return h.Sum(nil), nil
}

func checkSignature(msg, sig []byte) (bool, error) {
	s, err := sign(msg)
	if err != nil {
		return false, err
	}

	same := hmac.Equal(s, sig)
	return same, nil
}
