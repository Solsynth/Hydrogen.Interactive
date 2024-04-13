package main

import (
	"git.solsynth.dev/hydrogen/interactive/pkg/grpc"
	"git.solsynth.dev/hydrogen/interactive/pkg/server"
	"git.solsynth.dev/hydrogen/interactive/pkg/services"
	"github.com/robfig/cron/v3"
	"os"
	"os/signal"
	"syscall"

	interactive "git.solsynth.dev/hydrogen/interactive/pkg"
	"git.solsynth.dev/hydrogen/interactive/pkg/database"
	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
	"github.com/spf13/viper"
)

func init() {
	zerolog.TimeFieldFormat = zerolog.TimeFormatUnix
	log.Logger = log.Output(zerolog.ConsoleWriter{Out: os.Stdout})
}

func main() {
	// Configure settings
	viper.AddConfigPath(".")
	viper.AddConfigPath("..")
	viper.SetConfigName("settings")
	viper.SetConfigType("toml")

	// Load settings
	if err := viper.ReadInConfig(); err != nil {
		log.Panic().Err(err).Msg("An error occurred when loading settings.")
	}

	// Connect to database
	if err := database.NewSource(); err != nil {
		log.Fatal().Err(err).Msg("An error occurred when connect to database.")
	} else if err := database.RunMigration(database.C); err != nil {
		log.Fatal().Err(err).Msg("An error occurred when running database auto migration.")
	}

	// Connect other services
	go func() {
		if err := grpc.ConnectPassport(); err != nil {
			log.Fatal().Err(err).Msg("An error occurred when connecting to passport grpc endpoint...")
		}
	}()

	// Configure timed tasks
	quartz := cron.New(cron.WithLogger(cron.VerbosePrintfLogger(&log.Logger)))
	quartz.AddFunc("@every 60m", services.DoAutoDatabaseCleanup)
	quartz.Start()

	// Server
	server.NewServer()
	go server.Listen()

	// Messages
	log.Info().Msgf("Interactive v%s is started...", interactive.AppVersion)

	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit

	log.Info().Msgf("Interactive v%s is quitting...", interactive.AppVersion)

	quartz.Stop()
}
