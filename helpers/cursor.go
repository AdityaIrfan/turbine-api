package helpers

import (
	"encoding/base64"
	"encoding/json"
	"errors"
	"net/http"
	"strconv"
	"strings"

	"github.com/labstack/echo/v4"
	"github.com/phuslu/log"
)

func GenerateCursorPaginationByEcho(c echo.Context, sortMap map[string]string) (*Cursor, error) {
	cursorNextParam := c.QueryParam("Next")
	cursorPrevParam := c.QueryParam("Prev")

	if cursorNextParam != "" && cursorPrevParam != "" {
		return nil, errors.New("cannot use next and prev query params at the same time")
	}

	if cursorNextParam != "" {
		cursor, err := decodeCursor(cursorNextParam)
		if err != nil {
			log.Error().Str("ERROR ENCODE CURSOR NEXT", err.Error()).Msg("")
			return nil, errors.New(http.StatusText(http.StatusUnprocessableEntity))
		}
		return cursor, nil
	}

	if cursorPrevParam != "" {
		cursor, err := decodeCursor(cursorPrevParam)
		if err != nil {
			log.Error().Str("ERROR ENCODE CURSOR NEXT", err.Error()).Msg("")
			return nil, errors.New(http.StatusText(http.StatusUnprocessableEntity))
		}
		return cursor, nil
	}

	limitParam := c.QueryParam("PerPage")
	sortOrderParam := strings.ToLower(c.QueryParam("SortOrder"))
	sortByParams := strings.ToLower(c.QueryParam("SortBy"))
	searchParam := c.QueryParam("Search")

	limit, _ := strconv.Atoi(limitParam)
	if limit <= 0 {
		limit = 10
	}

	if sortByParams != "" {
		value, ok := sortMap[sortByParams]
		if !ok {
			value = "createdat"
		}
		sortByParams = value
	} else {
		sortByParams = "CreatedAt"
	}

	if sortOrderParam != "asc" && sortOrderParam != "desc" {
		sortOrderParam = "asc"
	}

	return &Cursor{
		Action:      NEXT,
		PerPage:     limit,
		CurrentPage: 1,
		SortOrder:   CursorSortOrder(sortOrderParam),
		SortBy:      sortByParams,
		Search:      searchParam,
	}, nil
}

func Reverse[T any](s []T) []T {
	for i, j := 0, len(s)-1; i < j; i, j = i+1, j-1 {
		s[i], s[j] = s[j], s[i]
	}
	return s
}

type CursorPagination struct {
	NextCursor string `json:"Next"`
	PrevCursor string `json:"Prev"`
}

type CursorSortOrder string
type CursorAction string

const (
	ASC  CursorSortOrder = "asc"
	DESC CursorSortOrder = "desc"
	NEXT CursorAction    = "next"
	PREV CursorAction    = "prev"
)

type Cursor struct {
	PerPage     int             `json:"PerPage"`
	CurrentPage int             `json:"CurrentPage"`
	SortOrder   CursorSortOrder `json:"SortOrder"`
	SortBy      string          `json:"SortBy"`
	Search      string          `json:"Search"`
	Action      CursorAction    `json:"Action"`
}

func (c *Cursor) GeneratePager(totalData int64) *CursorPagination {
	if totalData < int64(c.PerPage) {
		return &CursorPagination{
			NextCursor: "",
			PrevCursor: "",
		}
	}

	if c.Action == NEXT {
		totalPage := totalData / int64(c.PerPage)
		if totalPage%int64(c.PerPage) > 0 {
			totalPage++
		}

		if c.CurrentPage == 1 {
			nextCursor := c
			nextCursor.Action = NEXT
			nextCursor.CurrentPage++
			return &CursorPagination{
				NextCursor: encodeCursor(nextCursor),
				PrevCursor: "",
			}
		}

		if totalPage <= int64(c.CurrentPage) {
			c.Action = PREV
			c.CurrentPage--
			return &CursorPagination{
				NextCursor: "",
				PrevCursor: encodeCursor(c),
			}
		}

		nextCursor := c
		prevCursor := c
		nextCursor.Action = NEXT
		prevCursor.Action = PREV
		nextCursor.CurrentPage++
		prevCursor.CurrentPage--

		return &CursorPagination{
			NextCursor: encodeCursor(nextCursor),
			PrevCursor: encodeCursor(prevCursor),
		}
	} else if c.Action == PREV {
		if c.CurrentPage == 1 {
			c.CurrentPage++
			return &CursorPagination{
				NextCursor: encodeCursor(c),
				PrevCursor: "",
			}
		}

		nextCursor := c
		prevCursor := c
		nextCursor.Action = NEXT
		prevCursor.Action = PREV
		nextCursor.CurrentPage++
		prevCursor.CurrentPage--

		return &CursorPagination{
			NextCursor: encodeCursor(nextCursor),
			PrevCursor: encodeCursor(prevCursor),
		}
	}

	c.Action = NEXT
	c.CurrentPage++
	return &CursorPagination{
		NextCursor: encodeCursor(c),
		PrevCursor: "",
	}
}

func encodeCursor(cursor *Cursor) string {
	serializedCursor, err := json.Marshal(cursor)
	if err != nil {
		return ""
	}
	encodedCursor := base64.StdEncoding.EncodeToString(serializedCursor)
	return encodedCursor
}

func decodeCursor(cursor string) (*Cursor, error) {
	decodedCursor, err := base64.StdEncoding.DecodeString(cursor)
	if err != nil {
		return nil, err
	}
	var cur *Cursor
	if err := json.Unmarshal(decodedCursor, &cur); err != nil {
		return nil, err
	}
	return cur, nil
}
