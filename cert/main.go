package main

import (
	"log"
	"os"

	"github.com/creativeprojects/http-stream-server/server"
)

func main() {
	pemCert, pemKey, pin, err := server.KeyPairWithPin()
	if err != nil {
		log.Fatal(err)
	}

	err = save("certificate.pem", pemCert)
	if err != nil {
		log.Fatal(err)
	}

	err = save("key.pem", pemKey)
	if err != nil {
		log.Fatal(err)
	}

	err = save("pin.txt", pin)
	if err != nil {
		log.Fatal(err)
	}
}

func save(filename string, content []byte) error {
	file, err := os.OpenFile(filename, os.O_WRONLY|os.O_CREATE|os.O_TRUNC, 0600)
	if err != nil {
		return err
	}
	defer file.Close()

	_, err = file.Write(content)
	if err != nil {
		return err
	}
	return nil
}
