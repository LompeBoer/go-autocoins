package wickhunter

import (
	"fmt"
)

type PairsListValues struct {
}

func Calculate(list []Pair) {
	// cr := []string{
	// 	"BTCUSDT", "BTCBUSD", "BTCUSDT_210326", "BTCUSDT_210924", "BTCUSDT_210625", "BTCSTUSDT", "ETHUSDT_210326", "ETHUSDT_210625", "ETHUSDT_210924", "ETHUSD", "ETHUSDT", "BNBUSDT", "LRCUSDT", "QTUMUSDT", "BLZUSDT", "BELUSDT", "BANDUSDT", "TOMOUSDT", "KNCUSDT", "HNTUSDT", "ETHBUSD", "ERNUSDT", "COCOSUSDT", "TWTUSDT", "RUNEUSDT", "CRVUSDT", "ICPUSDT", "RSRUSDT",
	// }

	// blacklistNames := []string{}
	// blacklist := []Pair{}
	// currentBlacklist := []Pair{}
	// currentBlacklistNotFound := []Pair{}
	// for _, p := range list {
	// 	found := false
	// 	for _, c := range cr {
	// 		if c == p.Pair {
	// 			found = true
	// 			break
	// 		}
	// 	}
	// 	if found {
	// 		currentBlacklist = append(currentBlacklist, p)
	// 	} else {
	// 		currentBlacklistNotFound = append(currentBlacklistNotFound, p)
	// 	}

	// 	if !p.IsPermitted {
	// 		blacklistNames = append(blacklistNames, fmt.Sprintf("\"%s\"", p.Pair))
	// 		blacklist = append(blacklist, p)
	// 	}
	// }
	// blstr := strings.Join(blacklistNames, ",")

	// fmt.Println(blstr)

	fmt.Printf("List length:      %d\n", len(list))
	pc := 0
	ac := 0
	sc := 0
	for _, p := range list {
		if p.IsPermitted {
			pc++
		}
		if p.IsAvailable {
			ac++
		}
		if p.IsSafeAccount {
			sc++
		}
	}
	fmt.Printf("Permitted length: %d\n", pc)
	fmt.Printf("Available length: %d\n", ac)
	fmt.Printf("Safe length:      %d\n", sc)

	// for _, p := range list {
	// 	if !p.IsPermitted {
	// 		continue
	// 	}
	// 	fmt.Printf("%s p=%t a=%t s=%t\n", p.Pair, p.IsPermitted, p.IsAvailable, p.IsSafeAccount)
	// }

	// for _, p := range blO {
	// 	fmt.Printf("%s p=%t a=%t s=%t\n", p.Pair, p.IsPermitted, p.IsAvailable, p.IsSafeAccount)
	// }

}
