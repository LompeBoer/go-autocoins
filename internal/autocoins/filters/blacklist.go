package filters

import "github.com/LompeBoer/go-autocoins/internal/exchange/binance"

type BlackListFilter struct {
	BlackList []string
}

func (f *BlackListFilter) KeepSymbol(symbol binance.Symbol) bool {
	return !blackListContainsSymbol(f.BlackList, symbol.Name)
}

func blackListContainsSymbol(a []string, x string) bool {
	for _, v := range a {
		if v == x {
			return true
		}
	}
	return false
}
