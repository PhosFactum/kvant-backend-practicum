// internal/utils/utils.go
package utils

import (
    "errors"
    "strconv"

    "github.com/gin-gonic/gin"
)

// ParsePagination extracts page and limit from query parameters, with defaults.
func ParsePagination(c *gin.Context) (page, limit int, err error) {
    pageStr := c.DefaultQuery("page", "1")
    limitStr := c.DefaultQuery("limit", "10")

    page, err = strconv.Atoi(pageStr)
    if err != nil || page < 1 {
        return 0, 0, errors.New("invalid page parameter")
    }

    limit, err = strconv.Atoi(limitStr)
    if err != nil || limit < 1 {
        return 0, 0, errors.New("invalid limit parameter")
    }

    return page, limit, nil
}

// ParseAgeFilters extracts min_age and max_age from query parameters.
func ParseAgeFilters(c *gin.Context) (minAge, maxAge int, err error) {
    if s := c.Query("min_age"); s != "" {
        minAge, err = strconv.Atoi(s)
        if err != nil || minAge < 0 {
            return 0, 0, errors.New("invalid min_age parameter")
        }
    }
    if s := c.Query("max_age"); s != "" {
        maxAge, err = strconv.Atoi(s)
        if err != nil || maxAge < 0 {
            return 0, 0, errors.New("invalid max_age parameter")
        }
    }
    return minAge, maxAge, nil
}

// ParseIDParam extracts a uint ID from path parameters.
func ParseIDParam(c *gin.Context, param string) (uint, error) {
    s := c.Param(param)
    id, err := strconv.Atoi(s)
    if err != nil || id < 0 {
        return 0, errors.New("invalid " + param + " parameter")
    }
    return uint(id), nil
}

// Async runs a function in a new goroutine, recovering from any panics.
func Async(f func()) {
    go func() {
        defer func() {
            if r := recover(); r != nil {
                // optionally log r
            }
        }()
        f()
    }()
}

