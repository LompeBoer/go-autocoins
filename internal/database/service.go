package database

import "database/sql"

type DatabaseService interface {
	SelectInstruments() ([]Instrument, error)
	InsertInstrument(item Instrument) error
	DeleteInstrument(symbol string) error
	UpdateInstrument(symbol string, item Instrument) error
	CreateInstrumentTable() error
	BulkInsertInstruments(items []Instrument) error
	TruncateAndBulkInsertInstruments(items []Instrument) error

	SelectPermittedInstruments() ([]Instrument, error)
	SelectNonPermittedInstruments() ([]Instrument, error)
	UpdatePermittedList(permitted []string) error
	SelectOpenOrders() ([]string, error)
	Close() error
}

type Instrument struct {
	Symbol           sql.NullString `json:"Symbol"`
	IsPermitted      bool           `json:"IsPermitted"`
	IsDefaultSetting bool           `json:"IsDefaultSettings"`
}
