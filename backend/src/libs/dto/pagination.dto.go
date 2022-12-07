package dto

import "math"

/*
Pagination struct useful for: repository methods (offset and limit)
and controllers (set response headers)
*/
type Pagination struct {
	Page       int
	PageSize   int
	TotalCount int64
}

/*
NewPagination is a constructor for Pagination struct
*/
func NewPagination(page, pageSize int, totalCount int64) *Pagination {
	return &Pagination{
		Page:       page,
		PageSize:   pageSize,
		TotalCount: totalCount,
	}
}

// Offset calculates the offset with page and page_size
func (p *Pagination) Offset() int {
	return (p.Page - 1) * p.PageSize
}

// PageCount calculates the page count (quantity of available pages) with total count and page size
func (p *Pagination) PageCount() int {
	return int(math.Ceil(float64(p.TotalCount) / float64(p.PageSize)))
}
