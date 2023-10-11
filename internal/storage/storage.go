package storage

type Transaction struct {
	Operation int `json:"operation"`
	Sum       int `json:"sum"`
	To        int `json:"to"`
	From      int `json:"from"`
}
