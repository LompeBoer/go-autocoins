package whdbv0

type Order struct {
	TradeID          string  `json:"TradeId"`
	IsOpen           bool    `json:"IsOpen"`
	Range            int     `json:"Range"`
	CreationDateTime string  `json:"CreationDatetime"`
	ClientOrderID    string  `json:"ClientOrderId"`
	ID               string  `json:"Id"`
	Symbol           string  `json:"Symbol"`
	Type             string  `json:"type"`
	Direction        string  `json:"direction"`
	Quantity         float64 `json:"Quantity"`
	FilledQuantity   float64 `json:"FilledQuantity"`
	Price            float64 `json:"Price"`
	State            string  `json:"State"`
	ExecutionPrice   float64 `json:"ExecutionPrice"`
	Commision        float64 `json:"Commision"`
	Message          string  `json:"Message"`
	UpdateDateTime   string  `json:"UpdateDatetime"`
}

func (d *Database) SelectOpenOrders() ([]string, error) {
	rows, err := d.db.Query("SELECT [Symbol] FROM [Order] WHERE [State] = 'New'")
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

func (d *Database) SelectOrders() ([]Order, error) {
	rows, err := d.db.Query("SELECT * FROM [Order]")
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var items []Order
	for rows.Next() {
		var i Order
		if err := rows.Scan(
			&i.TradeID,
			&i.IsOpen,
			&i.Range,
			&i.CreationDateTime,
			&i.ClientOrderID,
			&i.ID,
			&i.Symbol,
			&i.Type,
			&i.Direction,
			&i.Quantity,
			&i.FilledQuantity,
			&i.Price,
			&i.State,
			&i.ExecutionPrice,
			&i.Commision,
			&i.Message,
			&i.UpdateDateTime,
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
