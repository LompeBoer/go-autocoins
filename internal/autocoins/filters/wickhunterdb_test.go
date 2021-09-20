package filters

import (
	"testing"

	"github.com/LompeBoer/go-autocoins/internal/binance"
	"github.com/LompeBoer/go-autocoins/internal/wickhunter"
)

func TestWickHunterDBFilter(t *testing.T) {
	symbolName := "TEST"

	filter := WickHunterDBFilter{
		Positions: []wickhunter.Position{{Symbol: symbolName}},
	}

	symbol := binance.Symbol{Name: symbolName}
	keep := filter.KeepSymbol(symbol)
	if !keep {
		t.Errorf("wickhunterdb filter invalid result: expected %v got %v", true, keep)
	}
}

func TestWickHunterDBFilterOther(t *testing.T) {
	filter := WickHunterDBFilter{
		Positions: []wickhunter.Position{{Symbol: "INCLUDED"}},
	}

	symbol := binance.Symbol{Name: "EXCLUDED"}
	keep := filter.KeepSymbol(symbol)
	if keep {
		t.Errorf("wickhunterdb filter invalid result: expected %v got %v", false, keep)
	}
}
