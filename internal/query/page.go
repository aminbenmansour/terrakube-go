package query

import (
	"fmt"
	"strconv"
	"strings"
)

// DefaultPageSize is Elide's implicit ceiling for a collection the client does
// not page explicitly.
const DefaultPageSize int64 = 500

// MaxPageSize is the largest page an Elide entity with @Paginate(10000)
// accepts. Larger requests are rejected rather than silently clamped, matching
// upstream.
const MaxPageSize int64 = 10000

// ParsePage applies Elide's offset paging: first is the page size and after is
// the decimal offset to start at. A page larger than the maximum is rejected,
// as Elide does, rather than silently clamped.
func ParsePage(first *int, after *string) (limit, offset int64, err error) {
	limit = DefaultPageSize
	if first != nil && *first > 0 {
		limit = int64(*first)
	}
	if limit > MaxPageSize {
		return 0, 0, fmt.Errorf("query: page size %d exceeds the maximum of %d", limit, MaxPageSize)
	}

	if after == nil || *after == "" {
		return limit, 0, nil
	}
	offset, err = strconv.ParseInt(strings.TrimSpace(*after), 10, 64)
	if err != nil || offset < 0 {
		return 0, 0, fmt.Errorf("query: invalid page offset %q", *after)
	}
	return limit, offset, nil
}
