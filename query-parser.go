package main

import (
	"errors"
	"net/url"
	"strconv"
	"strings"
)

const (
	defaultPage = 0
	defaultLimit = 20
	defaultOrder = "desc"
)

type QueryParamsParser struct {
	values url.Values
	Page int
	Limit int
	Errors []error
	Q string
	Order string
}

func (q *QueryParamsParser) CheckCorrectness() bool {
	q.parsePage()
	q.parseLimit()
	q.parseQ()
	q.parseOrder()

	return len(q.Errors) == 0
}

func (q *QueryParamsParser) parsePage() {
	pageString := q.values["page"]
	if pageString == nil {
		q.Page = defaultPage
		return
	}

	page, err := strconv.Atoi(pageString[0])
	if err != nil {
			q.Errors = append(q.Errors, errors.New("page must be number"))
			return
	}
	if page < 0 {
		q.Errors = append(q.Errors, errors.New("page must be positive number"))
		return
	}
	q.Page = page
}

func (q *QueryParamsParser) parseLimit() {
	limitString := q.values["limit"]
	if limitString == nil {
		q.Limit = defaultLimit
		return
	}

	limit, err := strconv.Atoi(limitString[0])
	if err != nil {
			q.Errors = append(q.Errors, errors.New("limit must be number"))
			return
	}
	if limit < 0 {
			q.Errors = append(q.Errors, errors.New("limit must be positive number"))
			return
	}
	q.Limit = limit
}

func (q *QueryParamsParser) parseQ() {
	qString := q.values["q"]
	if qString == nil {
		q.Q = ""
		return
	}

	q.Q = qString[0]
}

func (q *QueryParamsParser) parseOrder() {
	orderString := q.values["order"]
	if orderString == nil {
		q.Order = defaultOrder
		return
	}

	if strings.ToLower(orderString[0]) == defaultOrder {
		q.Order = defaultOrder
		return
	}

	q.Order = "asc"
}


