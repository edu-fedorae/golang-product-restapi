package main

type Product struct {
	ID       string  `json:"id"`
	Name     string  `json:"name"`
	Quantity int     `json:"quantity"`
	Price    float32 `json:"price"`
	Store    *Store  `json:"store"`
}

type Store struct {
	ID   string `json:"id"`
	Name string `json:"name"`
}
