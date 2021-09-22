package autocoins

import (
	"context"
	"fmt"
	"io"
	"log"
	"os"
	"time"

	"github.com/LompeBoer/go-autocoins/internal/pairslist"
)

// RunLoop the main loop that will be run by `Run`
func (a *AutoCoins) RunLoop() {
	log.Println("Calculating coin list ...")
	startTime := time.Now()

	// Download the permitted and safe pairs list from Google Docs.
	var pairsList []pairslist.Pair
	if a.Settings.Filters.GoogleSheet.Enabled {
		list, err := pairslist.ReadWithKey(a.Settings.Filters.GoogleSheet.APIKey)
		if err != nil {
			a.OutputWriter.WriteError(fmt.Sprintf("Unable to retrieve Google Doc WickHunter Pairs List: %s", err.Error()))
		}
		pairsList = list
	}

	// Process all the symbols.
	objects, lists, err := a.GetInfo(pairsList)
	if err != nil {
		a.OutputWriter.WriteError(err.Error())
	} else if a.DisableWrite {
		log.Println("READ ONLY not updating WickHunter")
	} else if len(lists.Permitted) == 0 {
		a.OutputWriter.WriteError("ERROR: No permitted coins (no action performed)")
	} else {
		a.BackupDatabase()
		a.BotAPI.UpdatePermittedList(lists.Permitted, lists.NotTrading)
	}

	a.outputRun(objects, lists, startTime)
}

func (a *AutoCoins) outputRun(objects []SymbolDataObject, lists SymbolLists, startTime time.Time) {
	if len(objects) > 0 {
		a.OutputWriter.WriteResult(objects, lists)

		p := len(lists.Permitted)
		q := len(lists.Quarantined)
		log.Printf("Permitted: %d Quarantined: %d Total: %d\n", p, q, p+q)
	}

	log.Printf("Elapsed: %s\n", time.Since(startTime))
	log.Printf("API Weight used: %d/%d\n", a.ExchangeAPI.UsedWeight, a.ExchangeAPI.WeightLimit)
}

// Start running the loop with a wait interval defined in settings.
func (a *AutoCoins) Run() {
	ctx, cancel := context.WithCancel(context.Background())
	a.ctx = ctx
	a.cancel = cancel
	a.wg.Add(1)
	defer a.wg.Done()

	a.IsRunning = true
	for {
		a.RunLoop()

		select {
		case <-ctx.Done():
			return
		case <-time.After(time.Duration(a.Settings.AutoCoins.Refresh) * time.Minute):
			a.ReloadConfig()
		}
	}
}

// Stop the run loop.
func (a *AutoCoins) Stop() {
	if !a.IsRunning {
		return
	}
	a.IsRunning = false
	a.ExchangeAPI.Cancel()
	a.cancel()
	a.wg.Wait()
}

// Reload the settings (from disk.)
func (a *AutoCoins) ReloadConfig() {
	a.Settings = *a.Settings.ReloadConfig()
}

func (a *AutoCoins) BackupDatabase() {
	original, err := os.Open(a.StorageFilename)
	if err != nil {
		log.Fatal(err)
	}
	defer original.Close()

	new, err := os.Create(a.StorageFilename + ".bak")
	if err != nil {
		log.Fatal(err)
	}
	defer new.Close()

	_, err = io.Copy(new, original)
	if err != nil {
		log.Fatal(err)
	}
}
