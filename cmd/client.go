package cmd

import (
	"github.com/param108/grpc-chat/chat/chatclient"
	"github.com/spf13/cobra"
)

var ClientCmd = &cobra.Command{
	Use:   "client",
	Short: "grpc client code",
	RunE:  clientCmdF,
}

func clientCmdF(cmd *cobra.Command, args []string) error {
	ch := chatclient.NewMemoryChatClient(9090)
	ch.Start(args)
	return nil
}
