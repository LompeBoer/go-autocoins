package main

import (
	"database/sql"
	"encoding/json"
	"log"
	"net/http"

	"github.com/LompeBoer/go-autocoins/internal/database/whdbv1"
	"github.com/LompeBoer/go-autocoins/internal/wickhunter"
)

type Controller struct {
	db *whdbv1.Database
}

func main() {
	c := Controller{
		db: whdbv1.New("storage.db"),
	}

	http.HandleFunc("/bot/positions", c.handlePositions)

	log.Fatal(http.ListenAndServe(":5001", nil))
}

func (c *Controller) handlePositions(w http.ResponseWriter, r *http.Request) {
	w.Header().Add("content-type", "application/json")
	w.Header().Add("expires", "0")
	w.Header().Add("pragma", "no-cache")

	instruments, err := c.db.SelectInstruments()
	if err != nil {
		log.Fatal(err)
	}

	retval := []wickhunter.Position{}
	for _, ins := range instruments {
		state, err := c.db.SelectPositionState(ins.Symbol.String)
		if err != nil && err != sql.ErrNoRows {
			log.Fatal(err)
		}

		status := state.Status
		if status == "" {
			status = "Neutral"
		}

		retval = append(retval, wickhunter.Position{
			Symbol:    ins.Symbol.String,
			Permitted: ins.IsPermitted,
			State:     status,
		})
	}

	b, err := json.Marshal(retval)
	if err != nil {
		log.Fatal(err)
	}

	w.Write(b)
}
