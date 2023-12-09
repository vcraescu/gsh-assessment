package domain

type OrderRow struct {
	Quantity int `json:"quantity,omitempty"`
	Pack     int `json:"pack,omitempty"`
}

type Order struct {
	Rows []OrderRow `json:"rows,omitempty"`
}

type Pack struct {
	Size int `json:"size,omitempty"`
}
