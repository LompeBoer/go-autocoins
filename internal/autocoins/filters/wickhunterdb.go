package filters

import (
	"github.com/LompeBoer/go-autocoins/internal/exchange/binance"
	"github.com/LompeBoer/go-autocoins/internal/wickhunter"
)

type WickHunterDBFilter struct {
	Positions []wickhunter.Position
}

func (f *WickHunterDBFilter) KeepSymbol(symbol binance.Symbol) bool {
	return positionContainsSymbol(f.Positions, symbol.Name)
}

func positionContainsSymbol(a []wickhunter.Position, x string) bool {
	for _, u := range a {
		if x == u.Symbol {
			return true
		}
	}
	return false
}
