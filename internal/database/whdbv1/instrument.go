package whdbv1

import (
	"log"

	"github.com/LompeBoer/go-autocoins/internal/database"
)

func (d *Database) SelectInstruments() ([]database.Instrument, error) {
	rows, err := d.db.Query("SELECT Symbol, IsPermitted, IsDefaultSettings FROM Instrument")
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var items []database.Instrument
	for rows.Next() {
		var i database.Instrument
		if err := rows.Scan(
			&i.Symbol,
			&i.IsPermitted,
			&i.IsDefaultSetting,
		); err != nil {
			return nil, err
		}
		items = append(items, i)
	}
	if err := rows.Close(); err != nil {
		return nil, err
	}
	if err := rows.Err(); err != nil {
		return nil, err
	}
	return items, nil
}

func (d *Database) InsertInstrument(item database.Instrument) error {
	stmt, err := d.db.Prepare("INSERT INTO Instrument(Symbol, IsPermitted, IsDefaultSettings) values(?, ?, ?)")
	if err != nil {
		return err
	}
	defer stmt.Close()
	_, err = stmt.Exec(item.Symbol, item.IsPermitted, item.IsDefaultSetting)
	if err != nil {
		return err
	}
	return nil
}

func (d *Database) DeleteInstrument(symbol string) error {
	stmt, err := d.db.Prepare("DELETE FROM Instrument WHERE Symbol = ?")
	if err != nil {
		return err
	}
	defer stmt.Close()
	_, err = stmt.Exec(symbol)
	if err != nil {
		return err
	}
	return nil
}

func (d *Database) UpdateInstrument(symbol string, item database.Instrument) error {
	stmt, err := d.db.Prepare("UPDATE Instrument SET Symbol = ?, IsPermitted = ?, IsDefaultSettings = ? WHERE Symbol = ?")
	if err != nil {
		return err
	}
	defer stmt.Close()
	_, err = stmt.Exec(item.Symbol, item.IsPermitted, item.IsDefaultSetting, symbol)
	if err != nil {
		return err
	}
	return nil
}

func (d *Database) CreateInstrumentTable() error {
	query := "DROP TABLE IF EXISTS \"Instrument\";CREATE TABLE IF NOT EXISTS [Instrument] (Symbol TEXT NOT NULL,IsPermitted INTEGER NOT NULL,IsDefaultSettings INTEGER NOT NULL);"
	_, err := d.db.Exec(query)
	if err != nil {
		return err
	}
	return nil
}

func (d *Database) BulkInsertInstruments(items []database.Instrument) error {
	tx, err := d.db.Begin()
	if err != nil {
		return err
	}
	stmt, err := tx.Prepare("INSERT INTO Instrument(Symbol, IsPermitted, IsDefaultSettings) values(?, ?, ?)")
	if err != nil {
		return err
	}
	defer stmt.Close()
	for _, item := range items {
		_, err = stmt.Exec(item.Symbol.String, item.IsPermitted, item.IsDefaultSetting)
		if err != nil {
			return err
		}
	}
	tx.Commit()
	return nil
}

func (d *Database) TruncateAndBulkInsertInstruments(items []database.Instrument) error {
	tx, err := d.db.Begin()
	if err != nil {
		log.Printf("ERROR: TruncateAndBulkInsertInstruments:Begin: %s\n", err.Error())
		return err
	}
	query := "DROP TABLE IF EXISTS \"Instrument\";CREATE TABLE IF NOT EXISTS [Instrument] (Symbol TEXT NOT NULL,IsPermitted INTEGER NOT NULL,IsDefaultSettings INTEGER NOT NULL);"
	_, err = tx.Exec(query)
	if err != nil {
		log.Printf("ERROR: TruncateAndBulkInsertInstruments:Exec: %s\n", err.Error())
		return err
	}

	stmt, err := tx.Prepare("INSERT INTO Instrument(Symbol, IsPermitted, IsDefaultSettings) values(?, ?, ?)")
	if err != nil {
		log.Printf("ERROR: TruncateAndBulkInsertInstruments:Prepare: %s\n", err.Error())
		return err
	}
	defer stmt.Close()
	for _, item := range items {
		_, err = stmt.Exec(item.Symbol.String, item.IsPermitted, item.IsDefaultSetting)
		if err != nil {
			log.Printf("ERROR: TruncateAndBulkInsertInstruments:Exec: %s\n", err.Error())
			return err
		}
	}
	err = tx.Commit()
	if err != nil {
		log.Printf("ERROR: TruncateAndBulkInsertInstruments:Commit: %s\n", err.Error())
		return err
	}
	return nil

}

func (d *Database) UpdatePermittedList(permitted []string) error {
	tx, err := d.db.Begin()
	if err != nil {
		log.Printf("ERROR: UpdatePermittedList:Begin: %s\n", err.Error())
		return err
	}

	_, err = tx.Exec("UPDATE Instrument SET IsPermitted = 0")
	if err != nil {
		log.Printf("ERROR: UpdatePermittedList:Exec(P=0): %s\n", err.Error())
		return err
	}

	stmt, err := tx.Prepare("UPDATE Instrument SET IsPermitted = ? WHERE Symbol = ?")
	if err != nil {
		log.Printf("ERROR: UpdatePermittedList:Prepare: %s\n", err.Error())
		return err
	}
	defer stmt.Close()
	for _, item := range permitted {
		_, err = stmt.Exec(true, item)
		if err != nil {
			log.Printf("ERROR: UpdatePermittedList:Exec: %s\n", err.Error())
			return err
		}
	}
	err = tx.Commit()
	if err != nil {
		log.Printf("ERROR: UpdatePermittedList:Commit: %s\n", err.Error())
		return err
	}
	return nil
}

func (d *Database) SelectInstrumentsForPermitted(permitted bool) ([]database.Instrument, error) {
	stmt, err := d.db.Prepare("SELECT Symbol, IsPermitted, IsDefaultSettings FROM Instrument WHERE IsPermitted=?")
	if err != nil {
		return nil, err
	}
	defer stmt.Close()
	rows, err := stmt.Query(permitted)
	defer rows.Close()
	var items []database.Instrument
	for rows.Next() {
		var i database.Instrument
		if err := rows.Scan(
			&i.Symbol,
			&i.IsPermitted,
			&i.IsDefaultSetting,
		); err != nil {
			return nil, err
		}
		items = append(items, i)
	}
	if err := rows.Close(); err != nil {
		return nil, err
	}
	if err := rows.Err(); err != nil {
		return nil, err
	}
	return items, nil
}

func (d *Database) SelectPermittedInstruments() ([]database.Instrument, error) {
	return d.SelectInstrumentsForPermitted(true)
}

func (d *Database) SelectNonPermittedInstruments() ([]database.Instrument, error) {
	return d.SelectInstrumentsForPermitted(false)
}
