package api

import "strings"

// filterSortPage is the in-memory implementation every list endpoint used
// before it moved to SQL. The database tests keep it as the reference for
// parity: a query that filters, sorts or pages differently fails one of them.
func filterSortPage[T any](items []T, q listQuery, name func(T) string) ([]T, int) {
	filtered := make([]T, 0, len(items))
	for _, item := range items {
		n := strings.ToLower(name(item))
		filter := strings.ToLower(q.FilterValue)
		switch {
		case q.FilterValue == "":
			filtered = append(filtered, item)
		case q.FilterPrefix && strings.HasPrefix(n, filter):
			filtered = append(filtered, item)
		case !q.FilterPrefix && n == filter:
			filtered = append(filtered, item)
		}
	}

	// Insertion sort keeps the helper dependency-free and the lists are small.
	for i := 1; i < len(filtered); i++ {
		for j := i; j > 0; j-- {
			less := strings.ToLower(name(filtered[j])) < strings.ToLower(name(filtered[j-1]))
			if q.SortDesc {
				less = !less
			}
			if !less {
				break
			}
			filtered[j], filtered[j-1] = filtered[j-1], filtered[j]
		}
	}

	total := len(filtered)
	if q.PageSize == 0 {
		return filtered, total
	}
	start := (q.PageNumber - 1) * q.PageSize
	if start >= total {
		return []T{}, total
	}
	end := start + q.PageSize
	if end > total {
		end = total
	}
	return filtered[start:end], total
}
