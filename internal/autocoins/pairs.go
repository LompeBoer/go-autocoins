package autocoins

import (
	"log"

	"github.com/LompeBoer/go-autocoins/internal/pairslist"
)

func (a *AutoCoins) SetPairs(useSafe bool) {
	list, err := pairslist.Read()
	if err != nil {
		log.Fatal(err)
	}

	positions, err := a.BotAPI.GetPositions()
	if err != nil {
		log.Fatal(err)
	}

	permittedCoins := []string{}
	quarantinedCoins := []string{}
	for _, p := range positions {
		pv := false
		for _, c := range list {
			if p.Symbol != c.Pair {
				continue
			}
			if (useSafe && c.IsPermitted && c.IsSafeAccount) || (!useSafe && c.IsPermitted) {
				permittedCoins = append(permittedCoins, c.Pair)
				pv = true
			}
		}
		if !pv {
			quarantinedCoins = append(quarantinedCoins, p.Symbol)
		}
	}

	a.BackupDatabase()
	err = a.BotAPI.UpdatePermittedList(permittedCoins, quarantinedCoins)
	if err != nil {
		log.Fatal(err)
	}

	log.Printf("Set %d pairs to permitted\n", len(permittedCoins))
}
