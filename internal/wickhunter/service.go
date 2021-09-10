package wickhunter

import "log"

func (a *API) UpdatePermittedList(permitted []string, quarantined []string) error {
	for _, symbol := range permitted {
		err := a.SetSymbolTrading(symbol, true)
		if err != nil {
			log.Printf("Error updating permitted symbol: %s\n", err.Error())
		}
	}
	for _, symbol := range quarantined {
		err := a.SetSymbolTrading(symbol, false)
		if err != nil {
			log.Printf("Error updating quarantined symbol: %s\n", err.Error())
		}
	}
	return nil
}

func (p *Position) IsOpen() bool {
	return (p.State != "Neutral")
}
