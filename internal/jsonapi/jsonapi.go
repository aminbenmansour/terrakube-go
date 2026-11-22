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
	"strings"
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

// Relationship is a to-one or to-many relationship in a resource object.
// Data is used for to-one links; DataList is used for to-many links. A non-nil
// empty DataList marshals as an empty array, which is how upstream exposes
// run-events.
type Relationship struct {
	Data     *ResourceIdentifier  `json:"data,omitempty"`
	DataList []ResourceIdentifier `json:"-"`
	Links    map[string]any       `json:"links,omitempty"`
	Meta     map[string]any       `json:"meta,omitempty"`
}

// MarshalJSON emits either a to-one data object or a to-many data array.
func (r Relationship) MarshalJSON() ([]byte, error) {
	if r.DataList != nil {
		return json.Marshal(struct {
			Data  []ResourceIdentifier `json:"data"`
			Links map[string]any       `json:"links,omitempty"`
			Meta  map[string]any       `json:"meta,omitempty"`
		}{Data: r.DataList, Links: r.Links, Meta: r.Meta})
	}
	type alias Relationship
	return json.Marshal(alias(r))
}

// ResourceIdentifier identifies a related resource. Attributes is only used
// when decoding the non-standard tag relationship payloads go-tfe sends on
// workspace creation; response links leave it empty.
type ResourceIdentifier struct {
	Type       string                     `json:"type"`
	ID         string                     `json:"id"`
	Attributes map[string]json.RawMessage `json:"attributes,omitempty"`
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
	LID           string                         `json:"lid,omitempty"`
	Attributes    map[string]json.RawMessage     `json:"attributes,omitempty"`
	Relationships map[string]RequestRelationship `json:"relationships,omitempty"`
}

// RequestRelationship is a relationship in an incoming document. Data holds a
// to-one relationship; DataList holds a to-many relationship.
type RequestRelationship struct {
	Data     *ResourceIdentifier
	DataList []ResourceIdentifier
}

// UnmarshalJSON accepts both JSON:API to-one objects and to-many arrays. The
// standard request model only needs to-one relationships, but workspace
// creation carries tag lists in relationship arrays.
func (r *RequestRelationship) UnmarshalJSON(body []byte) error {
	var raw struct {
		Data json.RawMessage `json:"data"`
	}
	if err := json.Unmarshal(body, &raw); err != nil {
		return err
	}
	if len(raw.Data) == 0 || string(raw.Data) == "null" {
		return nil
	}
	if raw.Data[0] == '[' {
		return json.Unmarshal(raw.Data, &r.DataList)
	}
	var identifier ResourceIdentifier
	if err := json.Unmarshal(raw.Data, &identifier); err != nil {
		return err
	}
	r.Data = &identifier
	return nil
}

// Write encodes a document with the correct content type.
func Write(w http.ResponseWriter, status int, doc Document, fields ...Fieldsets) {
	if len(fields) > 0 {
		doc = fields[0].Restrict(doc)
	}
	w.Header().Set("Content-Type", MediaType+"; charset=utf-8")
	w.WriteHeader(status)
	_ = json.NewEncoder(w).Encode(doc)
}

// WriteResource writes a single resource document.
func WriteResource(w http.ResponseWriter, status int, resource Resource, fields ...Fieldsets) {
	Write(w, status, Document{Data: resource}, fields...)
}

// Fieldsets holds the fields[type]=a,b constraints of a request. JSON:API
// requires the server to omit every attribute a parameter does not name.
type Fieldsets map[string]map[string]bool

// ParseFieldsets reads fields[type] query parameters. An empty value is a
// parameter that names no fields, which still counts as present.
func ParseFieldsets(r *http.Request) Fieldsets {
	var sets Fieldsets
	for key, values := range r.URL.Query() {
		if !strings.HasPrefix(key, "fields[") || !strings.HasSuffix(key, "]") {
			continue
		}
		resourceType := key[len("fields[") : len(key)-1]
		names := map[string]bool{}
		for _, value := range values {
			for _, name := range strings.Split(value, ",") {
				if name != "" {
					names[name] = true
				}
			}
		}
		if sets == nil {
			sets = Fieldsets{}
		}
		sets[resourceType] = names
	}
	return sets
}

// Restrict applies the fieldsets to a document. A resource type with no
// fields parameter loses its attributes entirely once any fields parameter is
// present, which is what upstream's Elide does. Names the resource does not
// carry are dropped: upstream answers 400 for those, but it knows attribute
// sets we do not model yet, and refusing a field the UI asked for would break
// a screen over an attribute we have not implemented.
func (f Fieldsets) Restrict(doc Document) Document {
	if len(f) == 0 {
		return doc
	}
	switch data := doc.Data.(type) {
	case Resource:
		doc.Data = f.restrict(data)
	case []Resource:
		restricted := make([]Resource, len(data))
		for i, resource := range data {
			restricted[i] = f.restrict(resource)
		}
		doc.Data = restricted
	}
	for i, resource := range doc.Included {
		doc.Included[i] = f.restrict(resource)
	}
	return doc
}

func (f Fieldsets) restrict(resource Resource) Resource {
	wanted, ok := f[resource.Type]
	if !ok {
		resource.Attributes = nil
		return resource
	}
	kept := make(map[string]any, len(wanted))
	for name := range wanted {
		if value, present := resource.Attributes[name]; present {
			kept[name] = value
		}
	}
	if len(kept) == 0 {
		kept = nil
	}
	resource.Attributes = kept
	return resource
}

// Page describes how a collection was cut. The zero value means no page
// parameter was sent, and no meta is emitted, which is what Elide does.
type Page struct {
	Number int
	Limit  int
	Totals bool
}

// WriteCollection writes a resource collection. When the request carried page
// parameters, a meta.page member is included in Elide's shape:
//
//	{"page":{"number":1,"limit":2,"totalPages":1,"totalRecords":2}}
//
// totalPages and totalRecords appear only for page[totals]. Verified against
// upstream 2.34.3 (docs/upstream/rest/organization-list.json).
func WriteCollection(w http.ResponseWriter, status int, resources []Resource, total int, page Page, fields ...Fieldsets) {
	if resources == nil {
		resources = []Resource{}
	}
	doc := Document{Data: resources}
	if page.Number > 0 || page.Limit > 0 || page.Totals {
		number := max(page.Number, 1)
		limit := max(page.Limit, 1)
		meta := map[string]any{"number": number, "limit": limit}
		if page.Totals {
			meta["totalPages"] = (total + limit - 1) / limit
			meta["totalRecords"] = total
		}
		doc.Meta = map[string]any{"page": meta}
	}
	Write(w, status, doc, fields...)
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

// RelationshipIdentifier reads a relationship document, {"data":{...}} or
// {"data":null}. id is empty when the document clears the relationship; sent is
// false when the body has no data member at all.
func RelationshipIdentifier(r *http.Request) (typ, id string, sent bool, err error) {
	if r.Body == nil {
		return "", "", false, errors.New("request body is empty")
	}
	defer r.Body.Close()

	var doc struct {
		Data json.RawMessage `json:"data"`
	}
	dec := json.NewDecoder(io.LimitReader(r.Body, maxBodyBytes))
	if err := dec.Decode(&doc); err != nil {
		return "", "", false, fmt.Errorf("invalid JSON: %w", err)
	}
	if len(doc.Data) == 0 {
		return "", "", false, errors.New("data is required")
	}
	if string(doc.Data) == "null" {
		return "", "", true, nil
	}
	var identifier ResourceIdentifier
	if err := json.Unmarshal(doc.Data, &identifier); err != nil {
		return "", "", true, errors.New("data must be a resource identifier")
	}
	if identifier.Type == "" || identifier.ID == "" {
		return "", "", true, errors.New("data.type and data.id are required")
	}
	return identifier.Type, identifier.ID, true, nil
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

// RelationshipIdentifiers returns the identifiers of a to-many relationship.
// present is true when the relationship was sent; an empty list means it was
// sent as null or an empty array.
func RelationshipIdentifiers(data *RequestData, name string) (ids []ResourceIdentifier, present bool, err error) {
	if data == nil || data.Relationships == nil {
		return nil, false, nil
	}
	rel, ok := data.Relationships[name]
	if !ok {
		return nil, false, nil
	}
	return rel.DataList, true, nil
}

// Link builds a to-one relationship.
func Link(resourceType, id string) Relationship {
	return Relationship{Data: &ResourceIdentifier{Type: resourceType, ID: id}}
}

// LinkList builds a to-many relationship. An empty ids slice is preserved as an
// empty JSON array rather than being omitted.
func LinkList(resourceType string, ids ...string) Relationship {
	data := make([]ResourceIdentifier, 0, len(ids))
	for _, id := range ids {
		data = append(data, ResourceIdentifier{Type: resourceType, ID: id})
	}
	return Relationship{DataList: data}
}
