package cmd

import (
	"bode.fun/go/oga/server"
	"github.com/charmbracelet/log"
	"github.com/spf13/cobra"
)

func NewServeCommand(logger *log.Logger) *cobra.Command {
	command := &cobra.Command{
		Use:   "serve",
		Short: "Start the server",
		RunE: func(cmd *cobra.Command, args []string) error {
			host, err := cmd.Flags().GetString("host")
			if err != nil {
				return err
			}

			port, err := cmd.Flags().GetUint("port")
			if err != nil {
				return err
			}

			server := server.New(logger,
				server.WithHost(host),
				server.WithPort(port),
			)

			return server.Serve()
		},
	}

	command.Flags().String("host", "localhost", "the host to listen on")
	command.Flags().UintP("port", "p", 3080, "the port to listen on")

	return command
}
