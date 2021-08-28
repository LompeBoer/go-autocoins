package main

import (
	"flag"
	"fmt"
	"log"
	"os"
	"os/signal"

	_ "embed"

	"github.com/LompeBoer/go-autocoins/internal/autocoins"
	"github.com/LompeBoer/go-autocoins/internal/binance"
	"github.com/LompeBoer/go-autocoins/internal/database"
	"github.com/LompeBoer/go-autocoins/internal/database/whdbv0"
	"github.com/LompeBoer/go-autocoins/internal/database/whdbv1"
	"github.com/LompeBoer/go-autocoins/internal/discord"
)

const VersionNumber = "0.9.12"

func main() {
	flags := initFlags()

	log.Printf("Starting autocoins (v%s)\n\n", VersionNumber)

	settings := autocoins.LoadConfig(flags.ConfigFilename)

	db := initDatabase(flags.StorageFilename, settings.Version)
	defer db.Close()

	autoCoins := initAutoCoins(db, settings, flags.StorageFilename)
	if flags.SetPairs || flags.SetSafePairs {
		autoCoins.SetPairs(flags.SetSafePairs)
	} else {
		go autoCoins.Run()

		stop := make(chan os.Signal)
		signal.Notify(stop, os.Interrupt)
		<-stop
		autoCoins.Stop()
	}
	log.Printf("Exiting autocoins")
}

func initAutoCoins(db database.DatabaseService, settings *autocoins.Settings, storageFilename string) autocoins.AutoCoins {
	discordHook := discord.DiscordWebHook{
		Enabled: true,
		URL:     settings.Discord,
	}
	autoCoins := autocoins.AutoCoins{
		Settings: *settings,
		API: binance.NewAPI(binance.BinanceAPIParams{
			BaseURL:            "https://fapi.binance.com",
			ProxyURL:           settings.Proxy,
			ProxyUser:          settings.ProxyUser,
			ProxyPassword:      settings.ProxyPass,
			DebugSaveResponses: false,
			DebugReadResponses: false,
		}),
		DB:                         db,
		MaxFailedSymbolsPercentage: 0.1,
		StorageFilename:            storageFilename,
		DisableWrite:               settings.Version == 1, // Because WickHunter does not yet pickup changes to the storage.db file disable writing for v1.0.
		OutputWriter: autocoins.OutputWriter{
			Writers: []autocoins.Writer{
				&autocoins.ConsoleOutputWriter{},
				&autocoins.DiscordOutputWriter{
					WebHook:        discordHook,
					Version:        VersionNumber,
					MentionOnError: settings.MentionOnError,
				},
			},
		},
	}
	return autoCoins
}

func initDatabase(storageFilename string, version int) database.DatabaseService {
	if version == 1 {
		return whdbv1.New(storageFilename)
	}

	return whdbv0.New(storageFilename)
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
