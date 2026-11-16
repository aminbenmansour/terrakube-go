package tfe

import (
	"net/http"
	"regexp"
	"strconv"
	"strings"
)

const (
	defaultWorkspacePageSize = 20
	maxWorkspacePageSize     = 100
)

var taggedParameter = regexp.MustCompile(`^filter\[tagged]\[(\d+)]\[(key|value)]$`)

type workspaceQuery struct {
	searchTags []string
	tagged     []tagFilter
	searchName string
	projectID  string
	pageNumber int
	pageSize   int
}

type tagFilter struct {
	key   string
	value string
}

func parseWorkspaceQuery(r *http.Request) workspaceQuery {
	values := r.URL.Query()
	query := workspaceQuery{
		pageNumber: 1,
		pageSize:   defaultWorkspacePageSize,
	}

	if tags := values.Get("search[tags]"); tags != "" {
		for _, tag := range strings.Split(tags, ",") {
			if tag = strings.TrimSpace(tag); tag != "" {
				query.searchTags = append(query.searchTags, tag)
			}
		}
	}
	query.searchName = strings.TrimSpace(values.Get("search[name]"))
	query.projectID = strings.TrimSpace(values.Get("filter[project][id]"))

	if raw := values.Get("page[number]"); raw != "" {
		if n, err := strconv.Atoi(strings.TrimSpace(raw)); err == nil && n > 0 {
			query.pageNumber = n
		}
	}
	if raw := values.Get("page[size]"); raw != "" {
		if n, err := strconv.Atoi(strings.TrimSpace(raw)); err == nil && n > 0 {
			query.pageSize = n
		}
	}
	if query.pageSize > maxWorkspacePageSize {
		query.pageSize = maxWorkspacePageSize
	}

	// Key/value tag filters arrive as filter[tagged][N][key] and
	// filter[tagged][N][value]. The index only pairs the two; its order is
	// not meaningful.
	byIndex := map[string]tagFilter{}
	for key, rawValues := range values {
		matches := taggedParameter.FindStringSubmatch(key)
		if matches == nil || len(rawValues) == 0 {
			continue
		}
		filter := byIndex[matches[1]]
		if matches[2] == "key" {
			filter.key = strings.TrimSpace(rawValues[0])
		} else {
			filter.value = rawValues[0]
		}
		byIndex[matches[1]] = filter
	}
	for _, filter := range byIndex {
		if filter.key != "" {
			query.tagged = append(query.tagged, filter)
		}
	}

	return query
}

func (q workspaceQuery) hasAnyCondition() bool {
	return len(q.searchTags) > 0 || len(q.tagged) > 0 || q.searchName != "" || q.projectID != ""
}

func (q workspaceQuery) offset() int {
	return (q.pageNumber - 1) * q.pageSize
}
