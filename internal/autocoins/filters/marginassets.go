package filters

import "github.com/LompeBoer/go-autocoins/internal/binance"

type MarginAssetsFilter struct {
	MarginAssets []string
}

func (f *MarginAssetsFilter) KeepSymbol(symbol binance.Symbol) bool {
	return marginAssetsContainsSymbol(f.MarginAssets, symbol.MarginAsset)
}

func marginAssetsContainsSymbol(a []string, x string) bool {
	for _, v := range a {
		if v == x {
			return true
		}
	}
	return false
}
