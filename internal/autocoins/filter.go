package autocoins

import (
	"github.com/LompeBoer/go-autocoins/internal/autocoins/filters"
	"github.com/LompeBoer/go-autocoins/internal/binance"
	"github.com/LompeBoer/go-autocoins/internal/pairslist"
	"github.com/LompeBoer/go-autocoins/internal/wickhunter"
)

type Filter interface {
	KeepSymbol(binance.Symbol) bool
}

func (a *AutoCoins) createFilters(usedSymbols []wickhunter.Position, symbols []binance.Symbol, pairsList []pairslist.Pair) []Filter {
	filterList := []Filter{}
	// Check if symbol is present in the WickHunter Bot Instrument table.
	if a.Settings.Filters.WickHunterDB {
		filterList = append(filterList, &filters.WickHunterDBFilter{Positions: usedSymbols})
	}
	// Check if symbol is on the blacklist in the settings file.
	if a.Settings.Filters.Blacklist {
		filterList = append(filterList, &filters.BlackListFilter{BlackList: a.Settings.BlackList})
	}
	// Check if the margin asset is permitted in the settings file.
	if a.Settings.Filters.MarginAssets && len(a.Settings.MarginAssets) > 0 {
		filterList = append(filterList, &filters.MarginAssetsFilter{MarginAssets: a.Settings.MarginAssets})
	}
	// Check if the symbol is permitted in the Google Doc file by STP Todd.
	if a.Settings.Filters.GoogleSheetPermitted || a.Settings.Filters.GoogleSheetSafe {
		filterList = append(filterList, &filters.GoogleSheetFilter{PairsList: pairsList, UseSafeList: a.Settings.Filters.GoogleSheetSafe})
	}

	return filterList
}

// filterSymbols filters out the symbols from the exchangeInfo that are not used in the local storage file.
// It also checks the MarginAssets setting and filters out any symbol which uses a margin asset not in this list.
func (a *AutoCoins) filterSymbols(usedSymbols []wickhunter.Position, symbols []binance.Symbol, pairsList []pairslist.Pair) ([]binance.Symbol, error) {
	filters := a.createFilters(usedSymbols, symbols, pairsList)

	keepSymbol := func(symbol binance.Symbol) bool {
		for _, filter := range filters {
			if !filter.KeepSymbol(symbol) {
				return false
			}
		}
		return true
	}

	i := 0
	for _, s := range symbols {
		if keepSymbol(s) {
			symbols[i] = s
			i++
		}
	}
	symbols = symbols[:i]
	return symbols, nil

}
