package autocoins

import (
	"log"

	"github.com/LompeBoer/go-autocoins/internal/wickhunter"
)

func (a *AutoCoins) SetPairs(useSafe bool) {
	list, err := wickhunter.ReadPairsList("")
	if err != nil {
		log.Fatal(err)
	}

	permittedCoins := []string{}
	for _, c := range list {
		if (useSafe && c.IsPermitted && c.IsSafeAccount) || (!useSafe && c.IsPermitted) {
			permittedCoins = append(permittedCoins, c.Pair)
		}
	}

	a.BackupDatabase()
	err = a.DB.UpdatePermittedList(permittedCoins)
	if err != nil {
		log.Fatal(err)
	}

	log.Printf("Set %d pairs to permitted\n", len(permittedCoins))
}
