package autocoins

import (
	"testing"
	"time"

	"github.com/LompeBoer/go-autocoins/internal/exchange/binance"
	"github.com/LompeBoer/go-autocoins/internal/pairslist"
	"github.com/LompeBoer/go-autocoins/internal/wickhunter"
)

func TestMakeLists(t *testing.T) {
	objects := []SymbolDataObject{
		{
			Symbol:    binance.Symbol{Name: "TEST"},
			Open:      false,
			Time:      time.Now(),
			APIFailed: false,
			Result: SymbolDataResult{
				Percent1Hour:  true,
				Percent4Hour:  true,
				Percent24Hour: true,
				AllTimeHigh:   true,
				Age:           true,
			},
			Values: SymbolDataValues{
				Percent1Hour:  []float64{0, 1, 2, 4},
				Percent4Hour:  0.0,
				Percent24Hour: 0.0,
				AllTimeHigh:   0.0,
				Age:           0.0,
			},
		},
	}
	positions := []wickhunter.Position{
		{
			Symbol:    "TEST",
			Permitted: true,
			State:     "Neutral",
		},
	}
	a := AutoCoins{}
	lists, err := a.makeLists(objects, positions)
	if err != nil {
		t.Errorf("error returned: %s", err.Error())
	}

	expectTrading := 1
	expectNotTrading := 0
	countTrading := len(lists.Permitted)
	countNotTrading := len(lists.NotTrading)
	if countTrading != expectTrading {
		t.Errorf("invalid count trading: expected %d got %d", expectTrading, countTrading)
	}
	if countNotTrading != expectNotTrading {
		t.Errorf("invalid count not trading: expected %d got %d", expectNotTrading, countNotTrading)
	}

	t.Logf("expected trading %d got %d; not trading %d got %d", expectTrading, countTrading, expectNotTrading, countNotTrading)
}

func TestFilterSymbols(t *testing.T) {
	a := AutoCoins{
		Settings: Settings{
			Filters: SettingsFilters{
				BlackList: []string{"TEST"},
			},
		},
	}
	positions := []wickhunter.Position{
		{
			Symbol:    "TEST",
			Permitted: true,
			State:     "Neutral",
		},
	}

	binanceSymbols := []binance.Symbol{{Name: "TEST"}}
	symbols, err := a.filterSymbols(positions, binanceSymbols, []pairslist.Pair{})
	if err != nil {
		t.Errorf("filterSymbols returned error: %s", err.Error())
	}

	expect := 0
	got := len(symbols)
	if got != expect {
		t.Errorf("invalid filter: got %d expect %d", got, expect)
	}
}
