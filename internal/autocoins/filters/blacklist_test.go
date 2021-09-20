package filters

import (
	"testing"

	"github.com/LompeBoer/go-autocoins/internal/binance"
)

func TestBlackListFilter(t *testing.T) {
	symbolName := "TEST"

	filter := BlackListFilter{
		BlackList: []string{symbolName},
	}

	symbol := binance.Symbol{Name: symbolName}
	keep := filter.KeepSymbol(symbol)
	if keep {
		t.Errorf("blacklist filter invalid result: expected %v got %v", false, keep)
	}
}

func TestBlackListFilterOther(t *testing.T) {
	filter := BlackListFilter{
		BlackList: []string{"INCLUDED"},
	}

	symbol := binance.Symbol{Name: "EXCLUDED"}
	keep := filter.KeepSymbol(symbol)
	if !keep {
		t.Errorf("blacklist filter 'other' invalid result: expected %v got %v", true, keep)
	}
}
