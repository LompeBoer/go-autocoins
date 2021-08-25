package wickhunter

import (
	"context"
	"encoding/csv"
	"errors"
	"fmt"
	"io"
	"log"
	"net/http"
	"net/url"
	"strings"

	"google.golang.org/api/option"
	"google.golang.org/api/sheets/v4"
)

const (
	DocumentID = "1XWadBbVkbdi5Ub7bFhCcAhqpHiQXBETbeTg644pkTdI"
	SheetName  = "Pairs list"
	ReadRange  = "A6:D"
)

type Pair struct {
	Pair          string
	IsPermitted   bool
	IsAvailable   bool
	IsSafeAccount bool
}

// ReadPairsList retrieves the pairs list by WickHunter.
// If no key is provided it will use the CSV export method.
func ReadPairsList(key string) ([]Pair, error) {
	if key == "" {
		return ReadPairsListWithoutKey()
	}
	ctx := context.Background()

	srv, err := sheets.NewService(ctx, option.WithAPIKey(key))
	if err != nil {
		log.Printf("Unable to retrieve Sheets client: %v\n", err)
		return nil, err
	}

	readRange := SheetName + "!" + ReadRange
	resp, err := srv.Spreadsheets.Values.Get(DocumentID, readRange).Do()
	if err != nil {
		log.Printf("Unable to retrieve data from sheet: %v\n", err)
		return nil, err
	}

	if len(resp.Values) == 0 {
		return []Pair{}, errors.New("no data found in pairs list")
	} else {
		list := []Pair{}
		for _, row := range resp.Values {
			v := []string{
				fmt.Sprint(row[0]),
				fmt.Sprint(row[1]),
				fmt.Sprint(row[2]),
				fmt.Sprint(row[3]),
			}
			list = append(list, convertRow(v))
		}

		return list, nil
	}
}

func ReadPairsListWithoutKey() ([]Pair, error) {
	sheet := url.PathEscape(SheetName)
	url := fmt.Sprintf("https://docs.google.com/spreadsheets/d/%s/gviz/tq?tqx=out:csv&sheet=%s&range=%s", DocumentID, sheet, ReadRange)
	resp, err := http.Get(url)
	if err != nil {
		return nil, err
	}

	list := []Pair{}
	r := csv.NewReader(resp.Body)
	for {
		row, err := r.Read()
		if err == io.EOF {
			break
		}
		if err != nil {
			log.Fatal(err)
		}

		list = append(list, convertRow(row))
	}

	return list, nil
}

func convertRow(row []string) Pair {
	coin := strings.TrimSuffix(fmt.Sprint(row[0]), "PERP")
	pair := Pair{
		Pair:          coin,
		IsPermitted:   row[1] == "TRUE",
		IsAvailable:   row[2] == "TRUE",
		IsSafeAccount: row[3] == "TRUE",
	}
	return pair
}
