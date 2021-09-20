package filters

import "github.com/LompeBoer/go-autocoins/internal/binance"

type MarginAssetsFilter struct {
	MarginAssets []string
}

func (f *MarginAssetsFilter) KeepSymbol(symbol binance.Symbol) bool {
	for _, asset := range f.MarginAssets {
		if symbol.MarginAsset == asset {
			return true
		}
	}
	return false
}
