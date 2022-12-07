package pagination

import (
	"net/url"
	"strconv"

	"bank-service/src/libs/dto"
	"bank-service/src/libs/errors"
)

/*
GetPaginationFromQuery receives a queryString from request, extracts page
and page_size, then returns pagination object
*/
func GetPaginationFromQuery(queryString url.Values) (*dto.Pagination, error) {
	pageStr := queryString.Get("page")
	pageSizeStr := queryString.Get("page_size")

	var page, pageSize int
	var err error
	if pageStr != "" {
		page, err = strconv.Atoi(pageStr)
		if err != nil {
			return nil, errors.ErrFieldValidation("page", "not_number", "")
		}
	}
	if pageSizeStr != "" {
		pageSize, err = strconv.Atoi(pageSizeStr)
		if err != nil {
			return nil, errors.ErrFieldValidation("page_size", "not_number", "")
		}
	}
	if page > 100 {
		return nil, errors.ErrPageTooHigh
	}
	if page < 1 {
		page = 1
	}
	if pageSize > 100 {
		return nil, errors.ErrPageSizeTooHigh
	}
	if pageSize < 1 {
		pageSize = 20
	}
	return dto.NewPagination(page, pageSize, 0), nil
}
