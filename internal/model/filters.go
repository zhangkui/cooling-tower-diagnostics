package model

import "time"

type Filter struct {
	TowerID string
	Sensor  string
	Quality string
	State   string
	From    time.Time
	To      time.Time
	Limit   int
	Offset  int
}

func (f Filter) WithTower(id string) Filter          { f.TowerID = id; return f }
func (f Filter) WithSensor(sensor string) Filter     { f.Sensor = sensor; return f }
func (f Filter) WithQuality(quality string) Filter   { f.Quality = quality; return f }
func (f Filter) WithState(state string) Filter       { f.State = state; return f }
func (f Filter) WithRange(from, to time.Time) Filter { f.From = from; f.To = to; return f }
func (f Filter) WithLimit(limit int) Filter          { f.Limit = limit; return f }
func (f Filter) WithOffset(offset int) Filter        { f.Offset = offset; return f }
func (f Filter) IsBounded() bool                     { return !f.From.IsZero() && !f.To.IsZero() }
func (f Filter) HasTower() bool                      { return f.TowerID != "" }
func (f Filter) HasSensor() bool                     { return f.Sensor != "" }
func (f Filter) Valid() bool {
	if f.Limit < 0 || f.Offset < 0 {
		return false
	}
	if f.IsBounded() && f.To.Before(f.From) {
		return false
	}
	return true
}

type Page struct {
	Number int
	Size   int
	Total  int
	Items  any
}

func NewPage(number, size, total int, items any) Page {
	if number < 1 {
		number = 1
	}
	if size < 1 {
		size = 20
	}
	return Page{Number: number, Size: size, Total: total, Items: items}
}
func (p Page) HasNext() bool     { return p.Number*p.Size < p.Total }
func (p Page) HasPrevious() bool { return p.Number > 1 }
func (p Page) Next() Page        { p.Number++; return p }
func (p Page) Previous() Page {
	if p.Number > 1 {
		p.Number--
	}
	return p
}
func (p Page) LastNumber() int {
	if p.Size <= 0 {
		return 1
	}
	n := (p.Total + p.Size - 1) / p.Size
	if n < 1 {
		n = 1
	}
	return n
}

type Sort struct {
	Field      string
	Descending bool
}

func (s Sort) Reverse() Sort           { s.Descending = !s.Descending; return s }
func (s Sort) Empty() bool             { return s.Field == "" }
func SortBy(field string) Sort         { return Sort{Field: field} }
func SortDescending(field string) Sort { return Sort{Field: field, Descending: true} }

type QueryOptions struct {
	Filter Filter
	Sort   Sort
	Page   Page
}

func DefaultQuery() QueryOptions {
	return QueryOptions{Filter: Filter{Limit: 50}, Sort: SortBy("created_at"), Page: NewPage(1, 50, 0, nil)}
}
func (q QueryOptions) Valid() bool                      { return q.Filter.Valid() && !q.Sort.Empty() }
func (q QueryOptions) WithFilter(f Filter) QueryOptions { q.Filter = f; return q }
func (q QueryOptions) WithSort(s Sort) QueryOptions     { q.Sort = s; return q }
func (q QueryOptions) WithPage(p Page) QueryOptions     { q.Page = p; return q }
