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
	ch := chatclient.NewMemoryChatClient(args[1], args[2], args[3])
	ch.Start(args)
	return nil
}
