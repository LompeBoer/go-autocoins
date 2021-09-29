package autocoins

import (
	"context"
	"fmt"
	"sort"
	"sync"

	"github.com/LompeBoer/go-autocoins/internal/exchange/binance"
	"github.com/LompeBoer/go-autocoins/internal/pairslist"
	"github.com/LompeBoer/go-autocoins/internal/wickhunter"
)

type AutoCoins struct {
	Settings                   Settings
	ExchangeAPI                *binance.API
	BotAPI                     *wickhunter.API
	ctx                        context.Context
	cancel                     context.CancelFunc
	wg                         sync.WaitGroup
	IsRunning                  bool
	MaxFailedSymbolsPercentage float64
	StorageFilename            string
	DisableWrite               bool
	OutputWriter               OutputWriter
}

// GetInfo retrieves all symbol data and calculates market swing.
// It returns a list of permitted coins to trade.
func (a *AutoCoins) GetInfo(pairsList []pairslist.Pair) ([]SymbolDataObject, SymbolLists, error) {
	exchangeInfo, err := a.ExchangeAPI.GetExchangeInfo()
	if err != nil {
		return nil, SymbolLists{}, err
	}

	usedSymbols, err := a.BotAPI.GetPositions()
	if err != nil {
		return nil, SymbolLists{}, fmt.Errorf("botapi:getpositions: %s", err.Error())
	}

	// Remove symbols from the list based on the enabled filters.
	symbols, err := a.filterSymbols(usedSymbols, exchangeInfo.Symbols, pairsList)
	if err != nil {
		return nil, SymbolLists{}, fmt.Errorf("unable to filter symbols: %s", err.Error())
	}
	sort.Sort(binance.BySymbolName(symbols))

	// Will pause execution when rate limit will be exceeded.
	a.ExchangeAPI.RateLimitChecks(len(symbols))

	prices24Hours, err := a.ExchangeAPI.GetTicker()
	if err != nil {
		return nil, SymbolLists{}, err
	}

	c := make(chan SymbolDataObject)
	count := a.RetrieveAllSymbolData(symbols, prices24Hours, c)

	objects := []SymbolDataObject{}
	for i := 0; i < count; i++ {
		object := <-c
		objects = append(objects, object)
	}

	positions, err := a.BotAPI.GetPositions()
	if err != nil {
		return nil, SymbolLists{}, fmt.Errorf("botapi:getpositions: %s", err.Error())
	}

	lists, err := a.makeLists(objects, positions)
	if err != nil {
		return nil, SymbolLists{}, fmt.Errorf("unable to make list: %s", err.Error())
	}

	// If not enough symbol data is retrieved from the API fail this run.
	percentageFailed := float64(len(lists.FailedToProcess)) / float64(len(symbols))
	if percentageFailed > a.MaxFailedSymbolsPercentage {
		return nil, SymbolLists{}, fmt.Errorf("unable to retrieve enough data from Binance API (%.0f%% failed)", percentageFailed*100)
	}

	return objects, lists, nil
}

// SymbolLists contains all the calculated lists.
type SymbolLists struct {
	Quarantined          []string // Quarantined symbols to quarantine.
	QuarantinedNew       []string // QuarantinedNew newly added symbols to the quarantine list.
	QuarantinedSkipped   []string // QuarantinedSkipped should be quarantined but have currently open trade.
	QuarantinedExcluded  []string // QuarantinedExcluded should be quarantined but have been excluded.
	QuarantinedCurrently []string // QuarantinedCurrently already quarantined.
	QuarantinedRemoved   []string // QuarantinedRemoved no longer quarantined.
	Permitted            []string // Permitted symbols allowed to trade.
	PermittedCurrently   []string // PermittedCurrently symbols that were already being traded.
	FailedToProcess      []string // FailedToProcess symbols that failed to retrieve enough data to make calculations.
	NotTrading           []string // NotTrading coins that are excluded from trading.
}

// makeLists makes the SymbolLists object, this groups all the symbols in a certain list.
func (a *AutoCoins) makeLists(objects []SymbolDataObject, positions []wickhunter.Position) (SymbolLists, error) {
	openPositions := []string{}
	permittedCurrently := []string{}
	quarantinedCurrently := []string{}
	for _, p := range positions {
		if p.IsOpen() {
			openPositions = append(openPositions, p.Symbol)
		}
		if p.Permitted {
			permittedCurrently = append(permittedCurrently, p.Symbol)
		} else {
			quarantinedCurrently = append(quarantinedCurrently, p.Symbol)
		}
	}

	sort.Strings(openPositions)
	sort.Strings(permittedCurrently)
	sort.Strings(quarantinedCurrently)

	// Quarantined / Permitted
	quarantined := []string{}
	permitted := []string{}
	quarantinedSkipped := []string{}  // Skipped because it is currently being traded.
	quarantinedExcluded := []string{} // Skipped because it is excluded.
	failed := []string{}
	for _, object := range objects {
		if ContainsString(openPositions, object.Symbol.Name) {
			object.Open = true
		}
		if ContainsString(a.Settings.Filters.ExcludeList, object.Symbol.Name) {
			object.Excluded = true
		}

		if object.Open || object.Excluded {
			permitted = append(permitted, object.Symbol.Name)
		}

		if object.APIFailed {
			failed = append(failed, object.Symbol.Name)
		} else if object.ShouldQuarantine() {
			if object.Open {
				quarantinedSkipped = append(quarantinedSkipped, object.Symbol.Name)
			} else if object.Excluded {
				quarantinedExcluded = append(quarantinedExcluded, object.Symbol.Name)
			} else {
				quarantined = append(quarantined, object.Symbol.Name)
			}
		} else if !object.Open && !object.Excluded {
			permitted = append(permitted, object.Symbol.Name)
		}
	}
	sort.Strings(quarantined)
	sort.Strings(permitted)
	sort.Strings(quarantinedSkipped)
	sort.Strings(quarantinedExcluded)
	sort.Strings(failed)

	quarantinedNew := []string{}
	for _, q := range quarantined {
		if !ContainsString(quarantinedCurrently, q) {
			quarantinedNew = append(quarantinedNew, q)
		}
	}
	sort.Strings(quarantinedNew)

	quarantinedRemoved := []string{}
	for _, p := range permitted {
		if !ContainsString(permittedCurrently, p) {
			quarantinedRemoved = append(quarantinedRemoved, p)
		}
	}
	sort.Strings(quarantinedRemoved)

	notTrading := []string{}
	for _, p := range positions {
		found := false
		for _, pem := range permitted {
			if p.Symbol == pem {
				found = true
				break
			}
		}

		if !found {
			notTrading = append(notTrading, p.Symbol)
		}
	}
	sort.Strings(notTrading)

	return SymbolLists{
		Quarantined:          quarantined,
		QuarantinedNew:       quarantinedNew,
		QuarantinedSkipped:   quarantinedSkipped,
		QuarantinedExcluded:  quarantinedExcluded,
		QuarantinedCurrently: quarantinedCurrently,
		QuarantinedRemoved:   quarantinedRemoved,
		Permitted:            permitted,
		PermittedCurrently:   permittedCurrently,
		FailedToProcess:      failed,
		NotTrading:           notTrading,
	}, nil
}

func (a *AutoCoins) RetrieveAllSymbolData(symbols []binance.Symbol, prices24Hours []binance.Ticker, c chan SymbolDataObject) int {
	count := 0
	for _, symbol := range symbols {
		go a.RetrieveSymbolData(symbol, &prices24Hours, c)
		count++
	}
	return count
}
