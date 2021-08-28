package autocoins

import (
	"context"
	"fmt"
	"io"
	"log"
	"os"
	"time"

	"github.com/LompeBoer/go-autocoins/internal/wickhunter"
)

// RunLoop the main loop that will be run by `Run`
func (a *AutoCoins) RunLoop() {
	log.Println("Calculating coin list ...")
	startTime := time.Now()

	// Download the permitted and safe pairs list from Google Docs.
	var pairsList []wickhunter.Pair
	if a.Settings.Filters.GoogleSheetPermitted || a.Settings.Filters.GoogleSheetSafe {
		list, err := wickhunter.ReadPairsList(a.Settings.GoogleApiKey)
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
		log.Println("READ ONLY not updating WickHunter DB")
	} else if len(lists.Permitted) == 0 {
		a.OutputWriter.WriteError("ERROR: No permitted coins (no action performed)")
	} else {
		a.BackupDatabase()
		a.DB.UpdatePermittedList(lists.Permitted)

		a.OutputWriter.WriteResult(objects, lists)
	}

	elapsed := time.Since(startTime)
	log.Printf("Elapsed: %s\n", elapsed)
	log.Printf("API Weight used: %d/%d\n", a.API.UsedWeight, a.API.WeightLimit)
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
		case <-time.After(time.Duration(a.Settings.Refresh) * time.Minute):
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
	a.API.Cancel()
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
