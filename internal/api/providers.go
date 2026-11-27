package api

import (
	"net/http"

	"github.com/google/uuid"

	"github.com/aminbenmansour/terrakube-go/internal/jsonapi"
	"github.com/aminbenmansour/terrakube-go/internal/store"
)

type providerAttrs struct {
	Name              *string `json:"name"`
	Description       *string `json:"description"`
	Imported          *bool   `json:"imported"`
	RegistryNamespace *string `json:"registryNamespace"`
}

type providerVersionAttrs struct {
	VersionNumber      *string `json:"versionNumber"`
	Protocols          *string `json:"protocols"`
	Status             *string `json:"status"`
	DeprecationMessage *string `json:"deprecationMessage"`
}

type implementationAttrs struct {
	Os                  *string `json:"os"`
	Arch                *string `json:"arch"`
	Filename            *string `json:"filename"`
	DownloadURL         *string `json:"downloadUrl"`
	ShasumsURL          *string `json:"shasumsUrl"`
	ShasumsSignatureURL *string `json:"shasumsSignatureUrl"`
	Shasum              *string `json:"shasum"`
	KeyID               *string `json:"keyId"`
	AsciiArmor          *string `json:"asciiArmor"`
	TrustSignature      *string `json:"trustSignature"`
	Source              *string `json:"source"`
	SourceURL           *string `json:"sourceUrl"`
}

func providerResource(provider store.Provider, versions []store.Version) jsonapi.Resource {
	ids := make([]string, 0, len(versions))
	for _, version := range versions {
		ids = append(ids, version.ID)
	}
	return jsonapi.Resource{
		Type: "provider",
		ID:   provider.ID,
		Attributes: map[string]any{
			"description":       provider.Description,
			"imported":          provider.Imported,
			"name":              provider.Name,
			"registryNamespace": provider.RegistryNamespace,
		},
		Relationships: map[string]jsonapi.Relationship{
			"organization": jsonapi.Link("organization", provider.OrganizationID),
			"version":      jsonapi.LinkList("version", ids...),
		},
	}
}

func providerVersionResource(version store.Version, implementations []store.Implementation) jsonapi.Resource {
	ids := make([]string, 0, len(implementations))
	for _, implementation := range implementations {
		ids = append(ids, implementation.ID)
	}
	return jsonapi.Resource{
		Type: "version",
		ID:   version.ID,
		Attributes: map[string]any{
			"deprecationMessage": version.DeprecationMessage,
			"protocols":          version.Protocols,
			"status":             version.Status,
			"versionNumber":      version.VersionNumber,
		},
		Relationships: map[string]jsonapi.Relationship{
			"implementation": jsonapi.LinkList("implementation", ids...),
			"provider":       jsonapi.Link("provider", version.ProviderID),
		},
	}
}

func implementationResource(implementation store.Implementation) jsonapi.Resource {
	return jsonapi.Resource{
		Type: "implementation",
		ID:   implementation.ID,
		Attributes: map[string]any{
			"arch":                implementation.Arch,
			"asciiArmor":          implementation.AsciiArmor,
			"downloadUrl":         implementation.DownloadUrl,
			"filename":            implementation.Filename,
			"keyId":               implementation.KeyID,
			"os":                  implementation.Os,
			"shasum":              implementation.Shasum,
			"shasumsSignatureUrl": implementation.ShasumsSignatureUrl,
			"shasumsUrl":          implementation.ShasumsUrl,
			"source":              implementation.Source,
			"sourceUrl":           implementation.SourceUrl,
			"trustSignature":      implementation.TrustSignature,
		},
		Relationships: map[string]jsonapi.Relationship{
			"version": jsonapi.Link("version", implementation.VersionID),
		},
	}
}

func (a *API) listProviders(w http.ResponseWriter, r *http.Request) {
	orgID := param(r, "organizationID")
	if _, err := a.deps.Store.GetOrganization(r.Context(), orgID); err != nil {
		if !a.writeStoreError(w, "get organization", err) {
			a.internalError(w, "get organization", err)
		}
		return
	}
	providers, err := a.deps.Store.ListProvidersByOrganization(r.Context(), orgID)
	if err != nil {
		if !a.writeStoreError(w, "list providers", err) {
			a.internalError(w, "list providers", err)
		}
		return
	}

	want := parseInclude(r)
	versions, implementations, ok := a.providerIncludes(w, r, orgID, want)
	if !ok {
		return
	}
	byProvider := map[string][]store.Version{}
	for _, version := range versions {
		byProvider[version.ProviderID] = append(byProvider[version.ProviderID], version)
	}

	resources := make([]jsonapi.Resource, 0, len(providers))
	for _, provider := range providers {
		resources = append(resources, providerResource(provider, byProvider[provider.ID]))
	}
	jsonapi.Write(w, http.StatusOK, jsonapi.Document{
		Data:     resources,
		Included: providerIncluded(versions, implementations, want),
	}, jsonapi.ParseFieldsets(r))
}

func (a *API) getProvider(w http.ResponseWriter, r *http.Request) {
	provider, ok := a.providerInOrganization(w, r)
	if !ok {
		return
	}
	want := parseInclude(r)
	versions, implementations, ok := a.providerIncludes(w, r, provider.OrganizationID, want)
	if !ok {
		return
	}
	versions = filterVersions(versions, provider.ID)
	jsonapi.Write(w, http.StatusOK, jsonapi.Document{
		Data:     providerResource(provider, versions),
		Included: providerIncluded(versions, implementations, want),
	}, jsonapi.ParseFieldsets(r))
}

// providerIncludes loads the versions (and implementations) a request asks to
// include. With no include the provider relationship still lists its versions,
// which is what upstream serialises after a create, so they are loaded anyway.
func (a *API) providerIncludes(w http.ResponseWriter, r *http.Request, orgID string, want map[string]bool) ([]store.Version, []store.Implementation, bool) {
	versions, err := a.deps.Store.ListProviderVersionsByOrganization(r.Context(), orgID)
	if err != nil {
		if !a.writeStoreError(w, "list provider versions", err) {
			a.internalError(w, "list provider versions", err)
		}
		return nil, nil, false
	}
	if !want["version.implementation"] {
		return versions, nil, true
	}
	implementations, err := a.deps.Store.ListImplementationsByOrganization(r.Context(), orgID)
	if err != nil {
		if !a.writeStoreError(w, "list implementations", err) {
			a.internalError(w, "list implementations", err)
		}
		return nil, nil, false
	}
	return versions, implementations, true
}

func filterVersions(versions []store.Version, providerID string) []store.Version {
	filtered := make([]store.Version, 0, len(versions))
	for _, version := range versions {
		if version.ProviderID == providerID {
			filtered = append(filtered, version)
		}
	}
	return filtered
}

func providerIncluded(versions []store.Version, implementations []store.Implementation, want map[string]bool) []jsonapi.Resource {
	if !want["version"] {
		return nil
	}
	includedVersions := make(map[string]bool, len(versions))
	byVersion := map[string][]store.Implementation{}
	for _, version := range versions {
		includedVersions[version.ID] = true
	}
	for _, implementation := range implementations {
		if includedVersions[implementation.VersionID] {
			byVersion[implementation.VersionID] = append(byVersion[implementation.VersionID], implementation)
		}
	}
	included := make([]jsonapi.Resource, 0, len(versions)+len(implementations))
	for _, version := range versions {
		included = append(included, providerVersionResource(version, byVersion[version.ID]))
	}
	for _, implementation := range implementations {
		if includedVersions[implementation.VersionID] {
			included = append(included, implementationResource(implementation))
		}
	}
	return included
}

func (a *API) createProvider(w http.ResponseWriter, r *http.Request) {
	orgID := param(r, "organizationID")
	if _, err := a.deps.Store.GetOrganization(r.Context(), orgID); err != nil {
		if !a.writeStoreError(w, "get organization", err) {
			a.internalError(w, "get organization", err)
		}
		return
	}
	attrs, _, ok := decodeProviderAttrs(w, r)
	if !ok {
		return
	}
	if attrs.Name == nil {
		dataIntegrityConflict(w)
		return
	}
	imported := false
	if attrs.Imported != nil {
		imported = *attrs.Imported
	}
	provider, err := a.deps.Store.CreateProvider(r.Context(), store.CreateProviderParams{
		ID:                uuid.NewString(),
		Name:              *attrs.Name,
		Description:       attrs.Description,
		OrganizationID:    orgID,
		Imported:          &imported,
		RegistryNamespace: attrs.RegistryNamespace,
	})
	if err != nil {
		if !a.writeStoreError(w, "create provider", err) {
			a.internalError(w, "create provider", err)
		}
		return
	}
	jsonapi.WriteResource(w, http.StatusCreated, providerResource(provider, nil), jsonapi.ParseFieldsets(r))
}

func (a *API) updateProvider(w http.ResponseWriter, r *http.Request) {
	provider, ok := a.providerInOrganization(w, r)
	if !ok {
		return
	}
	attrs, data, ok := decodeProviderAttrs(w, r)
	if !ok {
		return
	}
	name := provider.Name
	if jsonapi.HasAttribute(data, "name") {
		if attrs.Name == nil {
			dataIntegrityConflict(w)
			return
		}
		name = *attrs.Name
	}
	imported := provider.Imported
	if jsonapi.HasAttribute(data, "imported") {
		imported = attrs.Imported
	}
	if _, err := a.deps.Store.UpdateProvider(r.Context(), store.UpdateProviderParams{
		ID:                provider.ID,
		Name:              name,
		Description:       pickString(data, "description", attrs.Description, provider.Description),
		Imported:          imported,
		RegistryNamespace: pickString(data, "registryNamespace", attrs.RegistryNamespace, provider.RegistryNamespace),
	}); err != nil {
		if !a.writeStoreError(w, "update provider", err) {
			a.internalError(w, "update provider", err)
		}
		return
	}
	w.WriteHeader(http.StatusNoContent)
}

func (a *API) deleteProvider(w http.ResponseWriter, r *http.Request) {
	provider, ok := a.providerInOrganization(w, r)
	if !ok {
		return
	}
	if err := a.deps.Store.DeleteProvider(r.Context(), provider.ID); err != nil {
		if !a.writeStoreError(w, "delete provider", err) {
			a.internalError(w, "delete provider", err)
		}
		return
	}
	w.WriteHeader(http.StatusNoContent)
}

func (a *API) listProviderVersions(w http.ResponseWriter, r *http.Request) {
	provider, ok := a.providerInOrganization(w, r)
	if !ok {
		return
	}
	versions, err := a.deps.Store.ListProviderVersions(r.Context(), provider.ID)
	if err != nil {
		if !a.writeStoreError(w, "list versions", err) {
			a.internalError(w, "list versions", err)
		}
		return
	}
	resources := make([]jsonapi.Resource, 0, len(versions))
	for _, version := range versions {
		resources = append(resources, providerVersionResource(version, nil))
	}
	jsonapi.WriteCollection(w, http.StatusOK, resources, len(resources), jsonapi.Page{}, jsonapi.ParseFieldsets(r))
}

func (a *API) getProviderVersion(w http.ResponseWriter, r *http.Request) {
	version, ok := a.versionInProvider(w, r)
	if !ok {
		return
	}
	implementations, err := a.deps.Store.ListImplementations(r.Context(), version.ID)
	if err != nil {
		if !a.writeStoreError(w, "list implementations", err) {
			a.internalError(w, "list implementations", err)
		}
		return
	}
	jsonapi.WriteResource(w, http.StatusOK, providerVersionResource(version, implementations), jsonapi.ParseFieldsets(r))
}

func (a *API) createProviderVersion(w http.ResponseWriter, r *http.Request) {
	provider, ok := a.providerInOrganization(w, r)
	if !ok {
		return
	}
	attrs, _, ok := decodeProviderVersionAttrs(w, r)
	if !ok {
		return
	}
	if attrs.VersionNumber == nil || attrs.Protocols == nil {
		dataIntegrityConflict(w)
		return
	}
	status := "active"
	if attrs.Status != nil {
		status = *attrs.Status
	}
	if !validVersionStatus(status) {
		versionStatusError(w, status)
		return
	}
	version, err := a.deps.Store.CreateProviderVersion(r.Context(), store.CreateProviderVersionParams{
		ID:                 uuid.NewString(),
		VersionNumber:      *attrs.VersionNumber,
		Protocols:          *attrs.Protocols,
		ProviderID:         provider.ID,
		Status:             status,
		DeprecationMessage: attrs.DeprecationMessage,
	})
	if err != nil {
		if !a.writeStoreError(w, "create version", err) {
			a.internalError(w, "create version", err)
		}
		return
	}
	jsonapi.WriteResource(w, http.StatusCreated, providerVersionResource(version, nil), jsonapi.ParseFieldsets(r))
}

func (a *API) updateProviderVersion(w http.ResponseWriter, r *http.Request) {
	version, ok := a.versionInProvider(w, r)
	if !ok {
		return
	}
	attrs, data, ok := decodeProviderVersionAttrs(w, r)
	if !ok {
		return
	}
	status := version.Status
	if jsonapi.HasAttribute(data, "status") {
		if attrs.Status == nil {
			jsonapi.WriteUnprocessable(w, http.StatusBadRequest, "NotNull", "must not be null.", "status")
			return
		}
		status = *attrs.Status
	}
	if !validVersionStatus(status) {
		versionStatusError(w, status)
		return
	}
	if _, err := a.deps.Store.UpdateProviderVersion(r.Context(), store.UpdateProviderVersionParams{
		ID:                 version.ID,
		VersionNumber:      pickStringValue(data, "versionNumber", attrs.VersionNumber, version.VersionNumber),
		Protocols:          pickStringValue(data, "protocols", attrs.Protocols, version.Protocols),
		Status:             status,
		DeprecationMessage: pickString(data, "deprecationMessage", attrs.DeprecationMessage, version.DeprecationMessage),
	}); err != nil {
		if !a.writeStoreError(w, "update version", err) {
			a.internalError(w, "update version", err)
		}
		return
	}
	w.WriteHeader(http.StatusNoContent)
}

func (a *API) deleteProviderVersion(w http.ResponseWriter, r *http.Request) {
	version, ok := a.versionInProvider(w, r)
	if !ok {
		return
	}
	if err := a.deps.Store.DeleteProviderVersion(r.Context(), version.ID); err != nil {
		if !a.writeStoreError(w, "delete version", err) {
			a.internalError(w, "delete version", err)
		}
		return
	}
	w.WriteHeader(http.StatusNoContent)
}

func (a *API) listImplementations(w http.ResponseWriter, r *http.Request) {
	version, ok := a.versionInProvider(w, r)
	if !ok {
		return
	}
	implementations, err := a.deps.Store.ListImplementations(r.Context(), version.ID)
	if err != nil {
		if !a.writeStoreError(w, "list implementations", err) {
			a.internalError(w, "list implementations", err)
		}
		return
	}
	resources := make([]jsonapi.Resource, 0, len(implementations))
	for _, implementation := range implementations {
		resources = append(resources, implementationResource(implementation))
	}
	jsonapi.WriteCollection(w, http.StatusOK, resources, len(resources), jsonapi.Page{}, jsonapi.ParseFieldsets(r))
}

func (a *API) getImplementation(w http.ResponseWriter, r *http.Request) {
	implementation, ok := a.implementationInVersion(w, r)
	if !ok {
		return
	}
	jsonapi.WriteResource(w, http.StatusOK, implementationResource(implementation), jsonapi.ParseFieldsets(r))
}

func (a *API) createImplementation(w http.ResponseWriter, r *http.Request) {
	version, ok := a.versionInProvider(w, r)
	if !ok {
		return
	}
	attrs, _, ok := decodeImplementationAttrs(w, r)
	if !ok {
		return
	}
	if !requiredImplementationAttrs(w, attrs) {
		return
	}
	implementation, err := a.deps.Store.CreateImplementation(r.Context(), store.CreateImplementationParams{
		ID:                  uuid.NewString(),
		Os:                  *attrs.Os,
		Arch:                attrs.Arch,
		Filename:            *attrs.Filename,
		DownloadUrl:         *attrs.DownloadURL,
		ShasumsUrl:          *attrs.ShasumsURL,
		ShasumsSignatureUrl: *attrs.ShasumsSignatureURL,
		Shasum:              *attrs.Shasum,
		KeyID:               *attrs.KeyID,
		AsciiArmor:          attrs.AsciiArmor,
		TrustSignature:      *attrs.TrustSignature,
		Source:              *attrs.Source,
		SourceUrl:           *attrs.SourceURL,
		VersionID:           version.ID,
	})
	if err != nil {
		if !a.writeStoreError(w, "create implementation", err) {
			a.internalError(w, "create implementation", err)
		}
		return
	}
	jsonapi.WriteResource(w, http.StatusCreated, implementationResource(implementation), jsonapi.ParseFieldsets(r))
}

func (a *API) updateImplementation(w http.ResponseWriter, r *http.Request) {
	implementation, ok := a.implementationInVersion(w, r)
	if !ok {
		return
	}
	attrs, data, ok := decodeImplementationAttrs(w, r)
	if !ok {
		return
	}
	if _, err := a.deps.Store.UpdateImplementation(r.Context(), store.UpdateImplementationParams{
		ID:                  implementation.ID,
		Os:                  pickStringValue(data, "os", attrs.Os, implementation.Os),
		Arch:                pickString(data, "arch", attrs.Arch, implementation.Arch),
		Filename:            pickStringValue(data, "filename", attrs.Filename, implementation.Filename),
		DownloadUrl:         pickStringValue(data, "downloadUrl", attrs.DownloadURL, implementation.DownloadUrl),
		ShasumsUrl:          pickStringValue(data, "shasumsUrl", attrs.ShasumsURL, implementation.ShasumsUrl),
		ShasumsSignatureUrl: pickStringValue(data, "shasumsSignatureUrl", attrs.ShasumsSignatureURL, implementation.ShasumsSignatureUrl),
		Shasum:              pickStringValue(data, "shasum", attrs.Shasum, implementation.Shasum),
		KeyID:               pickStringValue(data, "keyId", attrs.KeyID, implementation.KeyID),
		AsciiArmor:          pickString(data, "asciiArmor", attrs.AsciiArmor, implementation.AsciiArmor),
		TrustSignature:      pickStringValue(data, "trustSignature", attrs.TrustSignature, implementation.TrustSignature),
		Source:              pickStringValue(data, "source", attrs.Source, implementation.Source),
		SourceUrl:           pickStringValue(data, "sourceUrl", attrs.SourceURL, implementation.SourceUrl),
	}); err != nil {
		if !a.writeStoreError(w, "update implementation", err) {
			a.internalError(w, "update implementation", err)
		}
		return
	}
	w.WriteHeader(http.StatusNoContent)
}

func (a *API) deleteImplementation(w http.ResponseWriter, r *http.Request) {
	implementation, ok := a.implementationInVersion(w, r)
	if !ok {
		return
	}
	if err := a.deps.Store.DeleteImplementation(r.Context(), implementation.ID); err != nil {
		if !a.writeStoreError(w, "delete implementation", err) {
			a.internalError(w, "delete implementation", err)
		}
		return
	}
	w.WriteHeader(http.StatusNoContent)
}

// providerInOrganization resolves the :providerID in the URL and refuses a
// provider from another organization, as the template routes do.
func (a *API) providerInOrganization(w http.ResponseWriter, r *http.Request) (store.Provider, bool) {
	id := param(r, "providerID")
	if _, err := uuid.Parse(id); err != nil {
		jsonapi.WriteError(w, http.StatusBadRequest, "Invalid Value", "Invalid value: "+id)
		return store.Provider{}, false
	}
	provider, err := a.deps.Store.GetProvider(r.Context(), id)
	if err != nil {
		if !a.writeStoreError(w, "get provider", err) {
			a.internalError(w, "get provider", err)
		}
		return store.Provider{}, false
	}
	if provider.OrganizationID != param(r, "organizationID") {
		jsonapi.WriteError(w, http.StatusNotFound, "Not Found", "Unknown identifier ["+id+"] for provider")
		return store.Provider{}, false
	}
	return provider, true
}

func (a *API) versionInProvider(w http.ResponseWriter, r *http.Request) (store.Version, bool) {
	provider, ok := a.providerInOrganization(w, r)
	if !ok {
		return store.Version{}, false
	}
	id := param(r, "versionID")
	if _, err := uuid.Parse(id); err != nil {
		jsonapi.WriteError(w, http.StatusBadRequest, "Invalid Value", "Invalid value: "+id)
		return store.Version{}, false
	}
	version, err := a.deps.Store.GetProviderVersion(r.Context(), id)
	if err != nil {
		if !a.writeStoreError(w, "get version", err) {
			a.internalError(w, "get version", err)
		}
		return store.Version{}, false
	}
	if version.ProviderID != provider.ID {
		jsonapi.WriteError(w, http.StatusNotFound, "Not Found", "Unknown identifier ["+id+"] for version")
		return store.Version{}, false
	}
	return version, true
}

func (a *API) implementationInVersion(w http.ResponseWriter, r *http.Request) (store.Implementation, bool) {
	version, ok := a.versionInProvider(w, r)
	if !ok {
		return store.Implementation{}, false
	}
	id := param(r, "implementationID")
	if _, err := uuid.Parse(id); err != nil {
		jsonapi.WriteError(w, http.StatusBadRequest, "Invalid Value", "Invalid value: "+id)
		return store.Implementation{}, false
	}
	implementation, err := a.deps.Store.GetImplementation(r.Context(), id)
	if err != nil {
		if !a.writeStoreError(w, "get implementation", err) {
			a.internalError(w, "get implementation", err)
		}
		return store.Implementation{}, false
	}
	if implementation.VersionID != version.ID {
		jsonapi.WriteError(w, http.StatusNotFound, "Not Found", "Unknown identifier ["+id+"] for implementation")
		return store.Implementation{}, false
	}
	return implementation, true
}

func validVersionStatus(status string) bool {
	switch status {
	case "active", "deprecated", "removed":
		return true
	}
	return false
}

func requiredImplementationAttrs(w http.ResponseWriter, attrs implementationAttrs) bool {
	required := []struct {
		name  string
		value *string
	}{
		{"os", attrs.Os}, {"filename", attrs.Filename}, {"downloadUrl", attrs.DownloadURL},
		{"shasumsUrl", attrs.ShasumsURL}, {"shasumsSignatureUrl", attrs.ShasumsSignatureURL},
		{"shasum", attrs.Shasum}, {"keyId", attrs.KeyID}, {"trustSignature", attrs.TrustSignature},
		{"source", attrs.Source}, {"sourceUrl", attrs.SourceURL},
	}
	for _, field := range required {
		if field.value == nil {
			dataIntegrityConflict(w)
			return false
		}
	}
	return true
}

func decodeProviderAttrs(w http.ResponseWriter, r *http.Request) (providerAttrs, *jsonapi.RequestData, bool) {
	data, err := jsonapi.DecodeRequest(r)
	if err != nil {
		jsonapi.WriteError(w, http.StatusBadRequest, "Invalid Request", err.Error())
		return providerAttrs{}, nil, false
	}
	if !requireType(w, data, "provider") {
		return providerAttrs{}, nil, false
	}
	attrs, err := jsonapi.DecodeAttributes[providerAttrs](data)
	if err != nil {
		jsonapi.WriteError(w, http.StatusBadRequest, "Invalid Attributes", err.Error())
		return providerAttrs{}, nil, false
	}
	return attrs, data, true
}

func decodeProviderVersionAttrs(w http.ResponseWriter, r *http.Request) (providerVersionAttrs, *jsonapi.RequestData, bool) {
	data, err := jsonapi.DecodeRequest(r)
	if err != nil {
		jsonapi.WriteError(w, http.StatusBadRequest, "Invalid Request", err.Error())
		return providerVersionAttrs{}, nil, false
	}
	if !requireType(w, data, "version") {
		return providerVersionAttrs{}, nil, false
	}
	attrs, err := jsonapi.DecodeAttributes[providerVersionAttrs](data)
	if err != nil {
		jsonapi.WriteError(w, http.StatusBadRequest, "Invalid Attributes", err.Error())
		return providerVersionAttrs{}, nil, false
	}
	return attrs, data, true
}

func decodeImplementationAttrs(w http.ResponseWriter, r *http.Request) (implementationAttrs, *jsonapi.RequestData, bool) {
	data, err := jsonapi.DecodeRequest(r)
	if err != nil {
		jsonapi.WriteError(w, http.StatusBadRequest, "Invalid Request", err.Error())
		return implementationAttrs{}, nil, false
	}
	if !requireType(w, data, "implementation") {
		return implementationAttrs{}, nil, false
	}
	attrs, err := jsonapi.DecodeAttributes[implementationAttrs](data)
	if err != nil {
		jsonapi.WriteError(w, http.StatusBadRequest, "Invalid Attributes", err.Error())
		return implementationAttrs{}, nil, false
	}
	return attrs, data, true
}

// versionStatusError writes the invalid-enum error shared by provider and
// module version updates: Elide names the value it rejected.
func versionStatusError(w http.ResponseWriter, status string) {
	jsonapi.WriteError(w, http.StatusBadRequest, "Invalid Value", "Invalid value: "+status)
}

// dataIntegrityConflict is Elide's answer when a required column is missing:
// a generic 409 rather than a 4xx naming the attribute (see globalvars.go).
func dataIntegrityConflict(w http.ResponseWriter) {
	jsonapi.WriteError(w, http.StatusConflict, "Conflict", "This operation conflicts with data that already exists.")
}
