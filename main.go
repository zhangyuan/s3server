package main

import (
	"fmt"
	"log"
)

func main() {
	if err := invoke(); err != nil {
		log.Fatalln(err)
	}
}

func invoke() error {
	fmt.Println("hello world")
	return nil
}
