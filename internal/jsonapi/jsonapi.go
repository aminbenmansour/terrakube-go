// Package jsonapi implements the small, typed subset of JSON:API 1.0 that
// Terrakube needs. Upstream uses Elide, so payload shapes and content types must
// match https://jsonapi.org/format/ rather than a bespoke envelope.
package jsonapi

import (
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/http"
)

// MediaType is the JSON:API content type.
const MediaType = "application/vnd.api+json"

const maxBodyBytes = 8 << 20 // 8 MiB

// Resource is a JSON:API resource object.
type Resource struct {
	Type          string                  `json:"type"`
	ID            string                  `json:"id,omitempty"`
	Attributes    map[string]any          `json:"attributes,omitempty"`
	Relationships map[string]Relationship `json:"relationships,omitempty"`
	Links         map[string]any          `json:"links,omitempty"`
	Meta          map[string]any          `json:"meta,omitempty"`
}

// Relationship is a to-one relationship in a resource object.
type Relationship struct {
	Data  *ResourceIdentifier `json:"data,omitempty"`
	Links map[string]any      `json:"links,omitempty"`
	Meta  map[string]any      `json:"meta,omitempty"`
}

// ResourceIdentifier identifies a related resource.
type ResourceIdentifier struct {
	Type string `json:"type"`
	ID   string `json:"id"`
}

// Document is a JSON:API top-level document.
type Document struct {
	Data     any            `json:"data,omitempty"`
	Meta     map[string]any `json:"meta,omitempty"`
	Links    map[string]any `json:"links,omitempty"`
	Included []Resource     `json:"included,omitempty"`
	Errors   []Error        `json:"errors,omitempty"`
}

// Error is a JSON:API error object.
type Error struct {
	Status string       `json:"status,omitempty"`
	Title  string       `json:"title,omitempty"`
	Detail string       `json:"detail,omitempty"`
	Source *ErrorSource `json:"source,omitempty"`
}

// ErrorSource points at the offending document location.
type ErrorSource struct {
	Pointer   string `json:"pointer,omitempty"`
	Parameter string `json:"parameter,omitempty"`
}

// RequestData is the data member of an incoming JSON:API document.
type RequestData struct {
	Type          string                         `json:"type"`
	ID            string                         `json:"id,omitempty"`
	Attributes    map[string]json.RawMessage     `json:"attributes,omitempty"`
	Relationships map[string]RequestRelationship `json:"relationships,omitempty"`
}

// RequestRelationship is a to-one relationship in an incoming document.
type RequestRelationship struct {
	Data *ResourceIdentifier `json:"data"`
}

// Write encodes a document with the correct content type.
func Write(w http.ResponseWriter, status int, doc Document) {
	w.Header().Set("Content-Type", MediaType+"; charset=utf-8")
	w.WriteHeader(status)
	_ = json.NewEncoder(w).Encode(doc)
}

// WriteResource writes a single resource document.
func WriteResource(w http.ResponseWriter, status int, resource Resource) {
	Write(w, status, Document{Data: resource})
}

// WriteCollection writes a resource collection. When totals is true a
// meta.totals member is included, matching Elide's page[totals] behaviour.
func WriteCollection(w http.ResponseWriter, status int, resources []Resource, total int, totals bool) {
	if resources == nil {
		resources = []Resource{}
	}
	doc := Document{Data: resources}
	if totals {
		doc.Meta = map[string]any{"totals": total}
	}
	Write(w, status, doc)
}

// WriteError writes a single JSON:API error document.
func WriteError(w http.ResponseWriter, status int, title, detail string) {
	Write(w, status, Document{
		Errors: []Error{{
			Status: fmt.Sprintf("%d", status),
			Title:  title,
			Detail: detail,
		}},
	})
}

// WriteUnprocessable writes an error that points at a request parameter.
func WriteUnprocessable(w http.ResponseWriter, status int, title, detail, parameter string) {
	Write(w, status, Document{
		Errors: []Error{{
			Status: fmt.Sprintf("%d", status),
			Title:  title,
			Detail: detail,
			Source: &ErrorSource{Parameter: parameter},
		}},
	})
}

// DecodeRequest reads an incoming JSON:API document.
func DecodeRequest(r *http.Request) (*RequestData, error) {
	if r.Body == nil {
		return nil, errors.New("request body is empty")
	}
	defer r.Body.Close()

	var doc struct {
		Data RequestData `json:"data"`
	}
	dec := json.NewDecoder(io.LimitReader(r.Body, maxBodyBytes))
	if err := dec.Decode(&doc); err != nil {
		return nil, fmt.Errorf("invalid JSON: %w", err)
	}
	if doc.Data.Type == "" {
		return nil, errors.New("data.type is required")
	}
	return &doc.Data, nil
}

// DecodeAttributes unmarshals data.attributes into T.
func DecodeAttributes[T any](data *RequestData) (T, error) {
	var out T
	if len(data.Attributes) == 0 {
		return out, nil
	}
	raw, err := json.Marshal(data.Attributes)
	if err != nil {
		return out, fmt.Errorf("encode attributes: %w", err)
	}
	if err := json.Unmarshal(raw, &out); err != nil {
		return out, fmt.Errorf("decode attributes: %w", err)
	}
	return out, nil
}

// HasAttribute reports whether the attribute was present in the request. This
// is how PATCH handlers distinguish "not sent" from "sent as null".
func HasAttribute(data *RequestData, name string) bool {
	if data == nil || data.Attributes == nil {
		return false
	}
	_, ok := data.Attributes[name]
	return ok
}

// RelationshipID returns the related resource ID. present is true when the
// relationship was sent; id is empty when it was sent as null.
func RelationshipID(data *RequestData, name string) (id string, present bool, err error) {
	if data == nil || data.Relationships == nil {
		return "", false, nil
	}
	rel, ok := data.Relationships[name]
	if !ok {
		return "", false, nil
	}
	if rel.Data == nil {
		return "", true, nil
	}
	if rel.Data.ID == "" {
		return "", true, errors.New("relationship data.id is required")
	}
	return rel.Data.ID, true, nil
}

// Link builds a to-one relationship.
func Link(resourceType, id string) Relationship {
	return Relationship{Data: &ResourceIdentifier{Type: resourceType, ID: id}}
}
