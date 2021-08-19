package whdbv1

import "time"

type PositionState struct {
	LaunchID             string    `json:"LaunchId"`
	DateTime             time.Time `json:"Datetime"`
	Symbol               string    `json:"Symbol"`
	Status               string    `json:"Status"`
	Side                 string    `json:"Side"`
	BuyCount             int64     `json:"BuyCount"`
	Quantity             float64   `json:"Quantity"`
	AveragePrice         float64   `json:"AveragePrice"`
	TakeProfitPrice      float64   `json:"TakeProfitPrice"`
	StopLossPrice        float64   `json:"StopLossPrice"`
	TakeProfitLimitPrice string    `json:"TakeProfitLimitPrice"`
	Reason               string    `json:"Reason"`
}

func (d *Database) SelectOpenOrders() ([]string, error) {
	query := `
		SELECT m1.Symbol
		FROM PositionState m1 LEFT JOIN PositionState m2
		ON (m1.Symbol = m2.Symbol AND m1.Datetime < m2.Datetime)
		WHERE m2.Datetime IS NULL
		AND (m1.Status = 'Open' OR m1.Status = 'InitOpening' OR m1.Status = 'TPLimitPlacing' OR m1.Status = 'DCAOpening');
	`
	rows, err := d.db.Query(query)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var items []string
	for rows.Next() {
		var i string
		if err := rows.Scan(
			&i,
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

func (d *Database) SelectPositionStates() ([]PositionState, error) {
	rows, err := d.db.Query("SELECT LaunchId,Datetime,Symbol,Status,Side,BuyCount,Quantity,AveragePrice,TakeProfitPrice,StopLossPrice,TakeProfitLimitPrice,Reason FROM PositionState ORDER BY Datetime ASC")
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var items []PositionState
	for rows.Next() {
		var i PositionState
		if err := rows.Scan(
			&i.LaunchID,
			&i.DateTime,
			&i.Symbol,
			&i.Status,
			&i.Side,
			&i.BuyCount,
			&i.Quantity,
			&i.AveragePrice,
			&i.TakeProfitPrice,
			&i.StopLossPrice,
			&i.TakeProfitLimitPrice,
			&i.Reason,
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
