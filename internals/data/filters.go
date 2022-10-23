// Filename : internal/data/filters.go

package data

import (
	"math"
	"strings"

	"todoapi.miguelavila.net/internals/validator"
)

type Filters struct {
	Page     int
	PageSize int
	Sort     string
	SortList []string
}

func ValidateFilters(v *validator.Validator, f Filters) {
	// check page and page size parameters
	v.Check(f.Page > 0, "page", "must be greater than zero")
	v.Check(f.Page <= 1000, "page", "must be maximum of 1000")
	v.Check(f.PageSize > 0, "page_size", "must be greater than zero")
	v.Check(f.PageSize <= 100, "page_size", "must be maximum of 100")
	// check that the sort parameter matches a value the acceptable sort list
	v.Check(validator.In(f.Sort, f.SortList...), "sort", "invalid sort value")
}

// sortColumn() methods safety extracts the sort field query parameters
func (f Filters) sortColumn() string {

	for _, safeValue := range f.SortList {
		if f.Sort == safeValue {
			return strings.TrimPrefix(f.Sort, "-")
		}
	}
	panic("unsafe sort parameters")
}

// sortOrder() methods determines where we should sort by ASC/DESC
func (f Filters) sortOrder() string {

	if strings.HasPrefix(f.Sort, "-") {
		return "DESC"
	}
	return "ASC"
}

// limit() methods determines the LIMIT
func (f Filters) limit() int {
	return f.PageSize
}

// offset() methods calculates the OFFSET
func (f Filters) offset() int {
	return (f.Page - 1) * f.PageSize
}

// type Metadata contains the metadata to help with pagination
type Metadata struct {
	CurrentPage  int `json:"current_page,omitempty"`
	PageSize     int `json:"page_size,omitempty"`
	FirstPage    int `json:"first_page,omitempty"`
	LastPage     int `json:"last_page,omitempty"`
	TotalRecords int `json:"total_records,omitempty"`
}

// calculatesMetadata() methods computes the values for the metadata fields
func calculatesMetadata(totalRecords int, page int, pageSize int) Metadata {
	if totalRecords == 0 {
		return Metadata{}

	}
	return Metadata{
		CurrentPage:  page,
		PageSize:     pageSize,
		FirstPage:    1,
		LastPage:     int(math.Ceil(float64(totalRecords) / float64(pageSize))),
		TotalRecords: totalRecords,
	}
}
