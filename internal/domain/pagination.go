package domain

type Page struct {
	Number int
	Size   int
	Sort   string
	Desc   bool
}
type PageResult[T any] struct {
	Items []T
	Total int
	Page  int
	Size  int
}

func (p Page) Normalize() Page {
	if p.Number < 1 {
		p.Number = 1
	}
	if p.Size < 1 {
		p.Size = 20
	}
	if p.Size > 200 {
		p.Size = 200
	}
	return p
}
func (p Page) Offset() int { q := p.Normalize(); return (q.Number - 1) * q.Size }
func (p Page) Allowed(column string) bool {
	switch column {
	case "code", "name", "year", "created_at", "status":
		return true
	default:
		return false
	}
}
func SlicePage[T any](items []T, p Page) PageResult[T] {
	q := p.Normalize()
	start := q.Offset()
	if start > len(items) {
		start = len(items)
	}
	end := start + q.Size
	if end > len(items) {
		end = len(items)
	}
	return PageResult[T]{Items: items[start:end], Total: len(items), Page: q.Number, Size: q.Size}
}
