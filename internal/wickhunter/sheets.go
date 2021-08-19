package wickhunter

import (
	"context"
	"errors"
	"fmt"
	"log"
	"strings"

	"google.golang.org/api/option"
	"google.golang.org/api/sheets/v4"
)

type Pair struct {
	Pair          string
	IsPermitted   bool
	IsAvailable   bool
	IsSafeAccount bool
}

func ReadPairsList(key string) ([]Pair, error) {
	ctx := context.Background()

	srv, err := sheets.NewService(ctx, option.WithAPIKey(key))
	if err != nil {
		log.Printf("Unable to retrieve Sheets client: %v\n", err)
		return nil, err
	}

	spreadsheetId := "1XWadBbVkbdi5Ub7bFhCcAhqpHiQXBETbeTg644pkTdI"
	readRange := "Pairs list!A5:D"
	resp, err := srv.Spreadsheets.Values.Get(spreadsheetId, readRange).Do()
	if err != nil {
		log.Printf("Unable to retrieve data from sheet: %v\n", err)
		return nil, err
	}

	if len(resp.Values) == 0 {
		return []Pair{}, errors.New("no data found in pairs list")
	} else {
		list := []Pair{}
		for _, row := range resp.Values {
			coin := strings.TrimSuffix(fmt.Sprint(row[0]), "PERP")
			pair := Pair{
				Pair:          coin,
				IsPermitted:   row[1] == "TRUE",
				IsAvailable:   row[2] == "TRUE",
				IsSafeAccount: row[3] == "TRUE",
			}
			list = append(list, pair)
		}

		return list, nil
	}
}
