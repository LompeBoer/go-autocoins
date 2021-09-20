package filters

import (
	"github.com/LompeBoer/go-autocoins/internal/binance"
	"github.com/LompeBoer/go-autocoins/internal/wickhunter"
)

type WickHunterDBFilter struct {
	Positions []wickhunter.Position
}

func (f *WickHunterDBFilter) KeepSymbol(symbol binance.Symbol) bool {
	for _, u := range f.Positions {
		if symbol.Name == u.Symbol {
			return true
		}
	}
	return false
}
