package main

import (
	"fmt"
	"log"

	"github.com/Fiagram/gateway/internal/app"
	"github.com/Fiagram/gateway/internal/configs"
	"github.com/spf13/cobra"
	"go.uber.org/fx"
)

var (
	version    string
	commitHash string
)

func main() {
	var configFilePath string

	rootCommand := &cobra.Command{
		Use:     "gateway",
		Short:   "Starts the gateway in standalone server mode.",
		Long:    "Gateway is a microservice for managing accounts belongs to Fiagram project.",
		Version: fmt.Sprintf("%s \ncommit: %s", version, commitHash),
		RunE: func(cmd *cobra.Command, _ []string) error {
			fx.New(
				fx.Supply(configs.ConfigFilePath(configFilePath)),
				app.Module,
			).Run()

			return nil
		},
	}

	rootCommand.Flags().StringVarP(&configFilePath,
		"config-file-path", "c", "",
		"Use the provided config file, otherwise the default embedded config applied.")

	if err := rootCommand.Execute(); err != nil {
		log.Panic(err)
	}
}
