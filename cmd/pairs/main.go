package main

import (
	"log"

	"github.com/LompeBoer/go-autocoins/internal/autocoins"
	"github.com/LompeBoer/go-autocoins/internal/wickhunter"
)

func main() {
	settings := autocoins.LoadConfig("autoCoins.json")

	list, err := wickhunter.ReadPairsList(settings.GoogleApiKey)
	if err != nil {
		log.Fatal(err)
	}
	wickhunter.Calculate(list)
}
