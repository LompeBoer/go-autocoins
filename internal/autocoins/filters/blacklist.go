package filters

import "github.com/LompeBoer/go-autocoins/internal/binance"

type BlackListFilter struct {
	BlackList []string
}

func (f *BlackListFilter) KeepSymbol(symbol binance.Symbol) bool {
	for _, b := range f.BlackList {
		if symbol.Name == b {
			return false
		}
	}

	return true
}
