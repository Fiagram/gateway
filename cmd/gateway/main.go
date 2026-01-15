package main

import (
	"context"
	"fmt"
	"log"

	"github.com/Fiagram/gateway/internal/configs"
	handler "github.com/Fiagram/gateway/internal/handler/http"
	logic "github.com/Fiagram/gateway/internal/logic/http"
	"github.com/Fiagram/gateway/internal/utils"
	"github.com/spf13/cobra"
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

			cfg, err := configs.NewConfig("")
			if err != nil {
				log.Panic("Failed to read the config file")
			}

			logger, cleanup, err := utils.InitializeLogger(cfg.Log)
			defer cleanup()

			oapiLogic := logic.NewOapiLogic(logger)

			httpHandler := handler.NewHttpServer(
				cfg.Http,
				oapiLogic,
				logger,
			)
			httpHandler.Start(context.Background())

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
