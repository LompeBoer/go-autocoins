package main

import (
	"log"

	"github.com/LompeBoer/go-autocoins/internal/pairslist"
)

func main() {
	list, err := pairslist.Read()
	if err != nil {
		log.Fatal(err)
	}

	permittedCoins := []string{}
	for _, c := range list {
		if c.IsPermitted {
			permittedCoins = append(permittedCoins, c.Pair)
		}
	}

	log.Println(permittedCoins)
}
