package filters

import (
	"testing"

	"github.com/LompeBoer/go-autocoins/internal/exchange/binance"
	"github.com/LompeBoer/go-autocoins/internal/pairslist"
)

func TestGoogleSheetFilterPermitted(t *testing.T) {
	symbolName := "TEST"

	filter := GoogleSheetFilter{
		UseSafeList: false,
		PairsList: []pairslist.Pair{
			{
				Pair:          symbolName,
				IsPermitted:   true,
				IsSafeAccount: false,
				IsAvailable:   false,
			},
		},
	}

	symbol := binance.Symbol{Name: symbolName}
	keep := filter.KeepSymbol(symbol)
	if !keep {
		t.Errorf("google sheet filter 'permitted' invalid result: expected %v got %v", false, keep)
	}
}

func TestGoogleSheetFilterSafe(t *testing.T) {
	symbolName := "TEST"

	filter := GoogleSheetFilter{
		UseSafeList: true,
		PairsList: []pairslist.Pair{
			{
				Pair:          symbolName,
				IsPermitted:   true,
				IsSafeAccount: true,
				IsAvailable:   false,
			},
		},
	}

	symbol := binance.Symbol{Name: symbolName}
	keep := filter.KeepSymbol(symbol)
	if !keep {
		t.Errorf("google sheet filter 'safe' invalid result: expected %v got %v", false, keep)
	}
}

func TestGoogleSheetFilterBlock(t *testing.T) {
	symbolName := "TEST"

	filter := GoogleSheetFilter{
		UseSafeList: false,
		PairsList: []pairslist.Pair{
			{
				Pair:          symbolName,
				IsPermitted:   false,
				IsSafeAccount: false,
				IsAvailable:   true,
			},
		},
	}

	symbol := binance.Symbol{Name: symbolName}
	keep := filter.KeepSymbol(symbol)
	if keep {
		t.Errorf("google sheet filter 'block' invalid result: expected %v got %v", true, keep)
	}
}

func TestGoogleSheetFilterOther(t *testing.T) {
	filter := GoogleSheetFilter{
		UseSafeList: false,
		PairsList: []pairslist.Pair{
			{
				Pair:          "INCLUDED",
				IsPermitted:   true,
				IsSafeAccount: true,
				IsAvailable:   false,
			},
		},
	}

	symbol := binance.Symbol{Name: "EXCLUDED"}
	keep := filter.KeepSymbol(symbol)
	if keep {
		t.Errorf("google sheet filter 'other' invalid result: expected %v got %v", true, keep)
	}
}
