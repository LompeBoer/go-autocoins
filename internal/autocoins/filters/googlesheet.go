package filters

import (
	"github.com/LompeBoer/go-autocoins/internal/exchange/binance"
	"github.com/LompeBoer/go-autocoins/internal/pairslist"
)

type GoogleSheetFilter struct {
	PairsList   []pairslist.Pair
	WhiteList   []string
	UseSafeList bool
}

func (f *GoogleSheetFilter) KeepSymbol(symbol binance.Symbol) bool {
	if whiteListContainsSymbol(f.WhiteList, symbol.Name) {
		return true
	}

	if n := pairsListContainsSymbol(f.PairsList, symbol.Name); n >= 0 {
		p := f.PairsList[n]
		if f.UseSafeList {
			if p.IsPermitted && p.IsSafeAccount {
				return true
			}
		} else if p.IsPermitted {
			return true
		}
	}

	return false
}

func pairsListContainsSymbol(a []pairslist.Pair, x string) int {
	for i, v := range a {
		if v.Pair == x {
			return i
		}
	}
	return -1
}

func whiteListContainsSymbol(a []string, x string) bool {
	for _, v := range a {
		if v == x {
			return true
		}
	}
	return false
}
