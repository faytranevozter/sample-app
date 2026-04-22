package helpers

import (
	"net/url"
	"strconv"
)

// GetLimitOffset calculates pagination parameters based on query values.
// It extracts the "page" and "limit" values from the provided URL query parameters,
// and computes the corresponding page, limit, and offset values. If the "page" or "limit"
// values are not provided or are invalid, it defaults to page 1 and a limit of 10, or
// a user-specified default limit if provided.
//
// Parameters:
//
//	query: URL query parameters containing "page" and "limit" values.
//	defaultLimit: Optional parameter to specify a default limit if "limit" is not provided.
//
// Returns:
//
//	page: The current page number.
//	limit: The number of items per page.
//	offset: The offset for the database query, calculated as (page - 1) * limit.
func GetLimitOffset(query url.Values, defaultLimit ...int) (page, limit, offset int64) {
	pageQuery := query.Get("page")
	pageInt, _ := strconv.Atoi(pageQuery)
	limitQuery := query.Get("limit")
	limitInt, _ := strconv.Atoi(limitQuery)

	if pageInt == 0 {
		pageInt = 1
	}

	if limitInt == 0 {
		if len(defaultLimit) > 0 {
			limitInt = defaultLimit[0]
		} else {
			limitInt = 10
		}
	}

	page = int64(pageInt)
	limit = int64(limitInt)
	offset = (page - 1) * limit

	return
}
