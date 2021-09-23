package filters

import (
	"testing"

	"github.com/LompeBoer/go-autocoins/internal/exchange/binance"
)

func TestMarginAssetsFilter(t *testing.T) {
	marginAsset := "TEST"

	filter := MarginAssetsFilter{
		MarginAssets: []string{marginAsset},
	}

	symbol := binance.Symbol{Name: "TESTNAME", MarginAsset: marginAsset}
	keep := filter.KeepSymbol(symbol)
	if !keep {
		t.Errorf("marginassets filter invalid result: expected %v got %v", true, keep)
	}
}

func TestMarginAssetsFilterOther(t *testing.T) {
	filter := MarginAssetsFilter{
		MarginAssets: []string{"INCLUDED"},
	}

	symbol := binance.Symbol{Name: "TESTNAME", MarginAsset: "EXCLUDED"}
	keep := filter.KeepSymbol(symbol)
	if keep {
		t.Errorf("marginassets filter 'other' invalid result: expected %v got %v", false, keep)
	}
}
