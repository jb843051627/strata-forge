package model

type Page[T any] struct {
	Items  []T `json:"items"`
	Offset int `json:"offset"`
	Limit  int `json:"limit"`
	Total  int `json:"total"`
}

func NewPage[T any](items []T, offset, limit, total int) Page[T] {
	if offset < 0 {
		offset = 0
	}
	if limit <= 0 {
		limit = len(items)
	}
	return Page[T]{Items: items, Offset: offset, Limit: limit, Total: total}
}

func (p Page[T]) HasNext() bool {
	return p.Offset+p.Limit < p.Total
}
