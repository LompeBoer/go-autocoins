package main

import (
	"log"

	"github.com/LompeBoer/go-autocoins/internal/autocoins"
	"github.com/LompeBoer/go-autocoins/internal/database"
	"github.com/LompeBoer/go-autocoins/internal/database/whdbv0"
	"github.com/LompeBoer/go-autocoins/internal/database/whdbv1"
	"github.com/LompeBoer/go-autocoins/internal/wickhunter"
)

func main() {
	list, err := wickhunter.ReadPairsList("")
	if err != nil {
		log.Fatal(err)
	}

	permittedCoins := []string{}
	for _, c := range list {
		if c.IsPermitted {
			permittedCoins = append(permittedCoins, c.Pair)
		}
	}

	db := initDatabase("storage.db", 1)
	defer db.Close()
	a := autocoins.AutoCoins{
		DB: db,
	}
	a.BackupDatabase()
	a.DB.UpdatePermittedList(permittedCoins)
}

func initDatabase(storageFilename string, version int) database.DatabaseService {
	if version == 1 {
		return whdbv1.New(storageFilename)
	}

	return whdbv0.New(storageFilename)
}
