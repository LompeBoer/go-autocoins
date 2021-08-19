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

const VersionNumber = "0.9.8"

func main() {
	flags := initFlags()

	log.Printf("Starting autocoins (v%s)\n\n", VersionNumber)

	settings := autocoins.LoadConfig(flags.ConfigFilename)

	db := initDatabase(flags.StorageFilename, settings.Version)
	defer db.Close()

	autoCoins := initAutoCoins(db, settings, flags.StorageFilename)
	go autoCoins.Run()

	stop := make(chan os.Signal)
	signal.Notify(stop, os.Interrupt)
	<-stop
	autoCoins.Stop()
	log.Printf("Exiting autocoins")
}

func initAutoCoins(db database.DatabaseService, settings *autocoins.Settings, storageFilename string) autocoins.AutoCoins {
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
		Discord: discord.DiscordWebHook{
			Enabled: true,
			URL:     settings.Discord,
		},
		DB:                         db,
		MaxFailedSymbolsPercentage: 0.1,
		StorageFilename:            storageFilename,
		DisableWrite:               settings.Version == 1, // Because WickHunter does not yet pickup changes to the storage.db file disable writing for v1.0.
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
}

func initFlags() StartupFlags {
	version := flag.Bool("version", false, "prints current go-autocoins version")
	noConfig := flag.Bool("noconfig", false, "use default settings without a config file")
	configFilename := flag.String("config", "autoCoins.json", "path to the config file")
	storageFilename := flag.String("storage", "storage.db", "path to the storage file")
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
	}
}
