package main

import (
	"fmt"
	"log"

	"github.com/darwinia-network/postman/src"
)

// main fn
func main() {
	config, err := postman.Config("./postman.toml")
	if err != nil {
		log.Fatalln(err)
	}

	for _, mail := range config.Get("mailing.list").([]interface{}) {
		addr := fmt.Sprintf("%s", mail)
		err := postman.Mail(config, addr, "This is from yandex")
		if err != nil {
			fmt.Println(err)
		}
	}

	fmt.Println("(Wait.) Oh yes, wait a minute, Mr. Postman!")
}
