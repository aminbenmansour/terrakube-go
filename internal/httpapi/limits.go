package httpapi

import (
	"bytes"
	"io"
	"net/http"
	"strconv"

	"github.com/aminbenmansour/terrakube-go/internal/jsonapi"
)

// LimitBody bounds the request body of a JSON surface. A declared length over
// the limit is refused before the handler runs; an undeclared (chunked) body is
// read up to the limit and refused if there is more. Streaming uploads (the
// state backend, configuration tarballs) have their own limits and are mounted
// outside this middleware.
func LimitBody(limit int64) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			if r.Body == nil || r.Body == http.NoBody {
				next.ServeHTTP(w, r)
				return
			}
			if r.ContentLength > limit {
				writeTooLarge(w, limit)
				return
			}
			if r.ContentLength >= 0 {
				// The JSON decoders cap their reads at the same limit.
				next.ServeHTTP(w, r)
				return
			}

			body, err := io.ReadAll(io.LimitReader(r.Body, limit+1))
			if err != nil {
				jsonapi.WriteError(w, http.StatusBadRequest, "Invalid Request", "The request body could not be read.")
				return
			}
			if int64(len(body)) > limit {
				writeTooLarge(w, limit)
				return
			}
			r.Body = io.NopCloser(bytes.NewReader(body))
			r.ContentLength = int64(len(body))
			next.ServeHTTP(w, r)
		})
	}
}

func writeTooLarge(w http.ResponseWriter, limit int64) {
	jsonapi.WriteError(w, http.StatusRequestEntityTooLarge, "Payload Too Large",
		"The request body exceeds the "+size(limit)+" limit.")
}

func size(bytes int64) string {
	if bytes%(1<<20) == 0 {
		return strconv.FormatInt(bytes>>20, 10) + " MiB"
	}
	return strconv.FormatInt(bytes, 10) + " bytes"
}
