package main

import (
	"flag"
	"fmt"
	"log"
	"os"
	"os/signal"

	"github.com/LompeBoer/go-autocoins/internal/autocoins"
	"github.com/LompeBoer/go-autocoins/internal/binance"
	"github.com/LompeBoer/go-autocoins/internal/discord"
	"github.com/LompeBoer/go-autocoins/internal/wickhunter"
)

const (
	VersionNumber = "0.10.0-pre"
)

func main() {
	flags := initFlags()

	log.Printf("Starting autocoins (v%s)\n\n", VersionNumber)

	checkLatestVersion()

	settings := autocoins.LoadConfig(flags.ConfigFilename)

	autoCoins := initAutoCoins(settings, flags.StorageFilename)
	if flags.SetPairs || flags.SetSafePairs {
		autoCoins.SetPairs(flags.SetSafePairs)
	} else {
		go autoCoins.Run()

		stop := make(chan os.Signal, 1)
		signal.Notify(stop, os.Interrupt)
		<-stop
		autoCoins.Stop()
	}
	log.Printf("Exiting autocoins")
}

func initAutoCoins(settings *autocoins.Settings, storageFilename string) autocoins.AutoCoins {
	discordHook := discord.DiscordWebHook{
		Enabled: true,
		URL:     settings.Discord.WebHook,
	}
	autoCoins := autocoins.AutoCoins{
		Settings: *settings,
		ExchangeAPI: binance.NewAPI(binance.APIParams{
			BaseURL:            "https://fapi.binance.com",
			ProxyURL:           settings.Proxy.Address,
			ProxyUser:          settings.Proxy.Username,
			ProxyPassword:      settings.Proxy.Password,
			DebugSaveResponses: false,
			DebugReadResponses: true,
		}),
		BotAPI:                     wickhunter.NewAPI(settings.API),
		MaxFailedSymbolsPercentage: 0.1,
		StorageFilename:            storageFilename,
		DisableWrite:               false,
		OutputWriter: autocoins.OutputWriter{
			Writers: []autocoins.Writer{
				&autocoins.ConsoleOutputWriter{},
				&autocoins.DiscordOutputWriter{
					WebHook:        discordHook,
					Version:        VersionNumber,
					MentionOnError: settings.Discord.MentionOnError,
				},
			},
		},
	}
	return autoCoins
}

type StartupFlags struct {
	NoConfig        bool
	ConfigFilename  string
	StorageFilename string
	SetPairs        bool
	SetSafePairs    bool
}

func initFlags() StartupFlags {
	version := flag.Bool("version", false, "prints current go-autocoins version")
	noConfig := flag.Bool("noconfig", false, "use default settings without a config file")
	configFilename := flag.String("config", "autoCoins.json", "path to the config file")
	storageFilename := flag.String("storage", "storage.db", "path to the storage file")
	setPairs := flag.Bool("pairs", false, "set pairs to permitted from the Google Sheet Pairs List and exits the program")
	setSafePairs := flag.Bool("safepairs", false, "set safe pairs to permitted from the Google Sheet Pairs List and exits the program")
	flag.Parse()

	if *version {
		fmt.Println(VersionNumber)
		os.Exit(0)
	}

	if _, err := os.Stat(*storageFilename); os.IsNotExist(err) {
		log.Fatalf("Storage file '%s' does not exist.\n", *storageFilename)
	}
	if *noConfig {
		*configFilename = ""
	} else if _, err := os.Stat(*configFilename); os.IsNotExist(err) {
		log.Fatalf("Config file '%s' does not exist.\n", *configFilename)
	}

	return StartupFlags{
		NoConfig:        *noConfig,
		ConfigFilename:  *configFilename,
		StorageFilename: *storageFilename,
		SetPairs:        *setPairs,
		SetSafePairs:    *setSafePairs,
	}
}
