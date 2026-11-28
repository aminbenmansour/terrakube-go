package api

import (
	"errors"
	"net/http"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5/pgconn"

	"github.com/aminbenmansour/terrakube-go/internal/jsonapi"
	"github.com/aminbenmansour/terrakube-go/internal/store"
)

type federatedAttrs struct {
	Name      *string `json:"name"`
	IssuerURL *string `json:"issuerUrl"`
	Audience  *string `json:"audience"`
}

type federatedClaimAttrs struct {
	ClaimKey   *string `json:"claimKey"`
	ClaimValue *string `json:"claimValue"`
}

func federatedResource(federated store.Federated, claims []store.FederatedClaim) jsonapi.Resource {
	claimIDs := make([]string, 0, len(claims))
	for _, claim := range claims {
		claimIDs = append(claimIDs, claim.ID)
	}
	return jsonapi.Resource{
		Type: "federated",
		ID:   federated.ID,
		Attributes: map[string]any{
			"audience":  federated.Audience,
			"issuerUrl": federated.IssuerUrl,
			"name":      federated.Name,
		},
		Relationships: map[string]jsonapi.Relationship{
			"claims": jsonapi.LinkList("federated_claim", claimIDs...),
		},
	}
}

func federatedClaimResource(claim store.FederatedClaim) jsonapi.Resource {
	return jsonapi.Resource{
		Type: "federated_claim",
		ID:   claim.ID,
		Attributes: map[string]any{
			"claimKey":   claim.ClaimKey,
			"claimValue": claim.ClaimValue,
		},
		Relationships: map[string]jsonapi.Relationship{
			"federated": jsonapi.Link("federated", claim.FederatedID),
		},
	}
}

// federatedNotBlank is the constraint violation Elide answers for the federated
// resources, unlike the generic 409 the policy resources get.
func federatedNotBlank(w http.ResponseWriter, attribute string) {
	jsonapi.Write(w, http.StatusBadRequest, jsonapi.Document{Errors: []jsonapi.Error{{
		Code:   "NotBlank",
		Source: &jsonapi.ErrorSource{Pointer: "/data/attributes/" + attribute},
		Detail: attribute + " is required",
		Meta:   map[string]any{"type": "ConstraintViolation", "property": attribute},
	}}})
}

func (a *API) listFederated(w http.ResponseWriter, r *http.Request) {
	federated, err := a.deps.Store.ListFederated(r.Context())
	if err != nil {
		if !a.writeStoreError(w, "list federated credentials", err) {
			a.internalError(w, "list federated credentials", err)
		}
		return
	}
	claims, err := a.deps.Store.ListAllFederatedClaims(r.Context())
	if err != nil {
		if !a.writeStoreError(w, "list federated claims", err) {
			a.internalError(w, "list federated claims", err)
		}
		return
	}
	byFederated := make(map[string][]store.FederatedClaim, len(federated))
	for _, claim := range claims {
		byFederated[claim.FederatedID] = append(byFederated[claim.FederatedID], claim)
	}
	resources := make([]jsonapi.Resource, 0, len(federated))
	for _, credential := range federated {
		resources = append(resources, federatedResource(credential, byFederated[credential.ID]))
	}
	jsonapi.WriteCollection(w, http.StatusOK, resources, len(resources), jsonapi.Page{}, jsonapi.ParseFieldsets(r))
}

func (a *API) federatedFromPath(w http.ResponseWriter, r *http.Request) (store.Federated, bool) {
	id := param(r, "federatedID")
	federated, err := a.deps.Store.GetFederated(r.Context(), id)
	if err != nil {
		if a.writeStoreError(w, "get federated credential", err) {
			return store.Federated{}, false
		}
		a.internalError(w, "get federated credential", err)
		return store.Federated{}, false
	}
	return federated, true
}

func (a *API) getFederated(w http.ResponseWriter, r *http.Request) {
	federated, ok := a.federatedFromPath(w, r)
	if !ok {
		return
	}
	claims, err := a.deps.Store.ListFederatedClaims(r.Context(), federated.ID)
	if err != nil {
		if !a.writeStoreError(w, "list federated claims", err) {
			a.internalError(w, "list federated claims", err)
		}
		return
	}
	jsonapi.WriteResource(w, http.StatusOK, federatedResource(federated, claims), jsonapi.ParseFieldsets(r))
}

func (a *API) createFederated(w http.ResponseWriter, r *http.Request) {
	attrs, data, ok := decodeFederatedAttrs(w, r)
	if !ok {
		return
	}
	if !federatedRequired(w, data, attrs, true) {
		return
	}
	email := callerEmail(r)
	federated, err := a.deps.Store.CreateFederated(r.Context(), store.CreateFederatedParams{
		ID:        uuid.NewString(),
		Name:      *attrs.Name,
		IssuerUrl: attrs.IssuerURL,
		Audience:  attrs.Audience,
		CreatedBy: email,
	})
	if err != nil {
		if !a.writeStoreError(w, "create federated credential", err) {
			a.internalError(w, "create federated credential", err)
		}
		return
	}
	jsonapi.WriteResource(w, http.StatusCreated, federatedResource(federated, nil), jsonapi.ParseFieldsets(r))
}

// federatedRequired reports the first missing attribute the way Elide does: a
// 400 constraint violation on that attribute. Create requires all three; update
// checks only the ones the document names.
func federatedRequired(w http.ResponseWriter, data *jsonapi.RequestData, attrs federatedAttrs, create bool) bool {
	fields := []struct {
		name  string
		value *string
	}{{"name", attrs.Name}, {"issuerUrl", attrs.IssuerURL}, {"audience", attrs.Audience}}
	for _, field := range fields {
		named := jsonapi.HasAttribute(data, field.name)
		if (create && field.value == nil) || (named && (field.value == nil || *field.value == "")) {
			federatedNotBlank(w, field.name)
			return false
		}
	}
	return true
}

func (a *API) updateFederated(w http.ResponseWriter, r *http.Request) {
	federated, ok := a.federatedFromPath(w, r)
	if !ok {
		return
	}
	attrs, data, ok := decodeFederatedAttrs(w, r)
	if !ok {
		return
	}
	if !federatedRequired(w, data, attrs, false) {
		return
	}
	_, err := a.deps.Store.UpdateFederated(r.Context(), store.UpdateFederatedParams{
		ID:        federated.ID,
		Name:      stringValue(pickString(data, "name", attrs.Name, &federated.Name)),
		IssuerUrl: pickString(data, "issuerUrl", attrs.IssuerURL, federated.IssuerUrl),
		Audience:  pickString(data, "audience", attrs.Audience, federated.Audience),
		UpdatedBy: callerEmail(r),
	})
	if err != nil {
		if !a.writeStoreError(w, "update federated credential", err) {
			a.internalError(w, "update federated credential", err)
		}
		return
	}
	w.WriteHeader(http.StatusNoContent)
}

func (a *API) deleteFederated(w http.ResponseWriter, r *http.Request) {
	federated, ok := a.federatedFromPath(w, r)
	if !ok {
		return
	}
	if err := a.deps.Store.DeleteFederated(r.Context(), federated.ID); err != nil {
		var pgErr *pgconn.PgError
		if errors.As(err, &pgErr) && pgErr.Code == "23503" {
			jsonapi.WriteError(w, http.StatusConflict, "Conflict", "This operation conflicts with data that already exists.")
			return
		}
		if !a.writeStoreError(w, "delete federated credential", err) {
			a.internalError(w, "delete federated credential", err)
		}
		return
	}
	w.WriteHeader(http.StatusNoContent)
}

func (a *API) listFederatedClaims(w http.ResponseWriter, r *http.Request) {
	federated, ok := a.federatedFromPath(w, r)
	if !ok {
		return
	}
	claims, err := a.deps.Store.ListFederatedClaims(r.Context(), federated.ID)
	if err != nil {
		if !a.writeStoreError(w, "list federated claims", err) {
			a.internalError(w, "list federated claims", err)
		}
		return
	}
	resources := make([]jsonapi.Resource, 0, len(claims))
	for _, claim := range claims {
		resources = append(resources, federatedClaimResource(claim))
	}
	jsonapi.WriteCollection(w, http.StatusOK, resources, len(resources), jsonapi.Page{}, jsonapi.ParseFieldsets(r))
}

func (a *API) createFederatedClaim(w http.ResponseWriter, r *http.Request) {
	federated, ok := a.federatedFromPath(w, r)
	if !ok {
		return
	}
	attrs, _, ok := decodeFederatedClaimAttrs(w, r)
	if !ok {
		return
	}
	for _, field := range []struct {
		name  string
		value *string
	}{{"claimKey", attrs.ClaimKey}, {"claimValue", attrs.ClaimValue}} {
		if field.value == nil || *field.value == "" {
			federatedNotBlank(w, field.name)
			return
		}
	}
	claim, err := a.deps.Store.CreateFederatedClaim(r.Context(), store.CreateFederatedClaimParams{
		ID:          uuid.NewString(),
		ClaimKey:    *attrs.ClaimKey,
		ClaimValue:  *attrs.ClaimValue,
		FederatedID: federated.ID,
	})
	if err != nil {
		if !a.writeStoreError(w, "create federated claim", err) {
			a.internalError(w, "create federated claim", err)
		}
		return
	}
	jsonapi.WriteResource(w, http.StatusCreated, federatedClaimResource(claim), jsonapi.ParseFieldsets(r))
}

func (a *API) federatedClaimFromPath(w http.ResponseWriter, r *http.Request) (store.FederatedClaim, bool) {
	federated, ok := a.federatedFromPath(w, r)
	if !ok {
		return store.FederatedClaim{}, false
	}
	id := param(r, "claimID")
	claim, err := a.deps.Store.GetFederatedClaim(r.Context(), id)
	if err != nil {
		if a.writeStoreError(w, "get federated claim", err) {
			return store.FederatedClaim{}, false
		}
		a.internalError(w, "get federated claim", err)
		return store.FederatedClaim{}, false
	}
	if claim.FederatedID != federated.ID {
		jsonapi.WriteError(w, http.StatusNotFound, "Not Found", "Unknown identifier ["+claim.ID+"] for claims")
		return store.FederatedClaim{}, false
	}
	return claim, true
}

func (a *API) getFederatedClaim(w http.ResponseWriter, r *http.Request) {
	claim, ok := a.federatedClaimFromPath(w, r)
	if !ok {
		return
	}
	jsonapi.WriteResource(w, http.StatusOK, federatedClaimResource(claim), jsonapi.ParseFieldsets(r))
}

func (a *API) updateFederatedClaim(w http.ResponseWriter, r *http.Request) {
	claim, ok := a.federatedClaimFromPath(w, r)
	if !ok {
		return
	}
	attrs, data, ok := decodeFederatedClaimAttrs(w, r)
	if !ok {
		return
	}
	for _, field := range []struct {
		name  string
		value *string
	}{{"claimKey", pickString(data, "claimKey", attrs.ClaimKey, &claim.ClaimKey)},
		{"claimValue", pickString(data, "claimValue", attrs.ClaimValue, &claim.ClaimValue)}} {
		if jsonapi.HasAttribute(data, field.name) && (field.value == nil || *field.value == "") {
			federatedNotBlank(w, field.name)
			return
		}
	}
	if _, err := a.deps.Store.UpdateFederatedClaim(r.Context(), store.UpdateFederatedClaimParams{
		ID:         claim.ID,
		ClaimKey:   stringValue(pickString(data, "claimKey", attrs.ClaimKey, &claim.ClaimKey)),
		ClaimValue: stringValue(pickString(data, "claimValue", attrs.ClaimValue, &claim.ClaimValue)),
	}); err != nil {
		if !a.writeStoreError(w, "update federated claim", err) {
			a.internalError(w, "update federated claim", err)
		}
		return
	}
	w.WriteHeader(http.StatusNoContent)
}

func (a *API) deleteFederatedClaim(w http.ResponseWriter, r *http.Request) {
	claim, ok := a.federatedClaimFromPath(w, r)
	if !ok {
		return
	}
	if err := a.deps.Store.DeleteFederatedClaim(r.Context(), claim.ID); err != nil {
		if !a.writeStoreError(w, "delete federated claim", err) {
			a.internalError(w, "delete federated claim", err)
		}
		return
	}
	w.WriteHeader(http.StatusNoContent)
}

func decodeFederatedAttrs(w http.ResponseWriter, r *http.Request) (federatedAttrs, *jsonapi.RequestData, bool) {
	data, err := jsonapi.DecodeRequest(r)
	if err != nil {
		jsonapi.WriteError(w, http.StatusBadRequest, "Invalid Request", err.Error())
		return federatedAttrs{}, nil, false
	}
	if !requireType(w, data, "federated") {
		return federatedAttrs{}, nil, false
	}
	attrs, err := jsonapi.DecodeAttributes[federatedAttrs](data)
	if err != nil {
		jsonapi.WriteError(w, http.StatusBadRequest, "Invalid Attributes", err.Error())
		return federatedAttrs{}, nil, false
	}
	return attrs, data, true
}

func decodeFederatedClaimAttrs(w http.ResponseWriter, r *http.Request) (federatedClaimAttrs, *jsonapi.RequestData, bool) {
	data, err := jsonapi.DecodeRequest(r)
	if err != nil {
		jsonapi.WriteError(w, http.StatusBadRequest, "Invalid Request", err.Error())
		return federatedClaimAttrs{}, nil, false
	}
	if !requireType(w, data, "federated_claim") {
		return federatedClaimAttrs{}, nil, false
	}
	attrs, err := jsonapi.DecodeAttributes[federatedClaimAttrs](data)
	if err != nil {
		jsonapi.WriteError(w, http.StatusBadRequest, "Invalid Attributes", err.Error())
		return federatedClaimAttrs{}, nil, false
	}
	return attrs, data, true
}
