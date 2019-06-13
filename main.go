package main

import (
	"github.com/param108/grpc-chat-server/cmd"
	"log"
	"os"
)

func main() {
	if err := cmd.Run(os.Args[1:]); err != nil {
		log.Fatalf("Unable to run the command %s ", err.Error())
	}
}
