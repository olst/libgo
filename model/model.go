package model

// Book - represents a book
type Book struct {
	ID     string  `json:"id,omitempty"`
	Title  string  `json:"title,omitempty"`
	Genres string  `json:"genres,omitempty"`
	Pages  int     `json:"pages,omitempty"`
	Price  float32 `json:"price,omitempty"`
}
