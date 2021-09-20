package filters

import (
	"github.com/LompeBoer/go-autocoins/internal/binance"
	"github.com/LompeBoer/go-autocoins/internal/pairslist"
)

type GoogleSheetFilter struct {
	PairsList   []pairslist.Pair
	UseSafeList bool
}

func (f *GoogleSheetFilter) KeepSymbol(symbol binance.Symbol) bool {
	found := false
	for _, p := range f.PairsList {
		if symbol.Name == p.Pair {
			if f.UseSafeList {
				if p.IsPermitted && p.IsSafeAccount {
					found = true
					break
				}
			} else if p.IsPermitted {
				found = true
				break
			}
		}
	}
	return found
}
