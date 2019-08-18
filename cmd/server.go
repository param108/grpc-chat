package cmd

import (
	"fmt"
	"github.com/param108/grpc-chat/chat/chatserver"
	"github.com/spf13/cobra"
	"io/ioutil"
	"os"
	"strconv"
)

var ServerCmd = &cobra.Command{
	Use:   "server",
	Short: "grpc server code",
	RunE:  serverCmdF,
}

func createPIDFile() error {
	pid := os.Getpid()
	err := ioutil.WriteFile("PID", []byte(strconv.Itoa(pid)), 0644)
	return err
}

func serverCmdF(cmd *cobra.Command, args []string) error {
	port, err := strconv.Atoi(os.Getenv("SERVER_PORT"))
	if err != nil {
		fmt.Printf("Failed to identify server port: %s\n", err.Error())
		return err
	}

	err = createPIDFile()
	if err != nil {
		fmt.Printf("Failed to open PID file %s", err.Error())
		return err
	}
	server, err := chatserver.NewChatServer(port)

	if err != nil {
		fmt.Printf("Failed to create new Server %s\n", err.Error())
		return err
	}
	server.Start()

	return nil
}
