package api

import (
	"net/http"

	"github.com/google/uuid"

	"github.com/aminbenmansour/terrakube-go/internal/auth"
	"github.com/aminbenmansour/terrakube-go/internal/jsonapi"
	"github.com/aminbenmansour/terrakube-go/internal/store"
)

type moduleAttrs struct {
	Name        *string `json:"name"`
	Description *string `json:"description"`
	Provider    *string `json:"provider"`
	Source      *string `json:"source"`
	TagPrefix   *string `json:"tagPrefix"`
	Folder      *string `json:"folder"`
}

type moduleVersionAttrs struct {
	Version            *string `json:"version"`
	Commit             *string `json:"commit"`
	GitTag             *string `json:"gitTag"`
	Status             *string `json:"status"`
	DeprecationMessage *string `json:"deprecationMessage"`
}

func moduleResource(module store.Module, organizationName string, versions []store.ModuleVersion) jsonapi.Resource {
	ids := make([]string, 0, len(versions))
	for _, version := range versions {
		ids = append(ids, version.ID)
	}
	return jsonapi.Resource{
		Type: "module",
		ID:   module.ID,
		Attributes: map[string]any{
			"createdBy":        module.CreatedBy,
			"createdDate":      timestamp(module.CreatedDate),
			"description":      module.Description,
			"downloadQuantity": module.DownloadQuantity,
			"folder":           module.Folder,
			"latestVersion":    module.LatestVersion,
			"name":             module.Name,
			"provider":         module.Provider,
			"registryPath":     organizationName + "/" + module.Name + "/" + module.Provider,
			"source":           module.Source,
			"tagPrefix":        module.TagPrefix,
			"updatedBy":        module.UpdatedBy,
			"updatedDate":      timestamp(module.UpdatedDate),
		},
		Relationships: map[string]jsonapi.Relationship{
			"organization": jsonapi.Link("organization", module.OrganizationID),
			"ssh":          optionalLink("ssh", module.SshID),
			"vcs":          optionalLink("vcs", module.VcsID),
			"version":      jsonapi.LinkList("module_version", ids...),
		},
	}
}

func moduleVersionResource(version store.ModuleVersion) jsonapi.Resource {
	return jsonapi.Resource{
		Type: "module_version",
		ID:   version.ID,
		Attributes: map[string]any{
			"commit":             version.CommitInfo,
			"deprecationMessage": version.DeprecationMessage,
			"gitTag":             version.GitTag,
			"status":             version.Status,
			"version":            version.Version,
		},
		Relationships: map[string]jsonapi.Relationship{
			"module": jsonapi.Link("module", version.ModuleID),
		},
	}
}

func optionalLink(resourceType string, id *string) jsonapi.Relationship {
	if id == nil || *id == "" {
		return jsonapi.NullLink()
	}
	return jsonapi.Link(resourceType, *id)
}

func (a *API) listModules(w http.ResponseWriter, r *http.Request) {
	orgID := param(r, "organizationID")
	organization, err := a.deps.Store.GetOrganization(r.Context(), orgID)
	if err != nil {
		if !a.writeStoreError(w, "get organization", err) {
			a.internalError(w, "get organization", err)
		}
		return
	}
	modules, err := a.deps.Store.ListModulesByOrganization(r.Context(), orgID)
	if err != nil {
		if !a.writeStoreError(w, "list modules", err) {
			a.internalError(w, "list modules", err)
		}
		return
	}
	resources, included, ok := a.moduleResources(w, r, modules, organization.Name)
	if !ok {
		return
	}
	jsonapi.Write(w, http.StatusOK, jsonapi.Document{Data: resources, Included: included}, jsonapi.ParseFieldsets(r))
}

func (a *API) getModule(w http.ResponseWriter, r *http.Request) {
	module, ok := a.moduleInOrganization(w, r)
	if !ok {
		return
	}
	organization, err := a.deps.Store.GetOrganization(r.Context(), module.OrganizationID)
	if err != nil {
		if !a.writeStoreError(w, "get organization", err) {
			a.internalError(w, "get organization", err)
		}
		return
	}
	resources, included, ok := a.moduleResources(w, r, []store.Module{module}, organization.Name)
	if !ok {
		return
	}
	jsonapi.Write(w, http.StatusOK, jsonapi.Document{Data: resources[0], Included: included}, jsonapi.ParseFieldsets(r))
}

// moduleResources builds the module documents and the relationships a request
// asks to include: vcs, ssh and the module versions.
func (a *API) moduleResources(w http.ResponseWriter, r *http.Request, modules []store.Module, organizationName string) ([]jsonapi.Resource, []jsonapi.Resource, bool) {
	want := parseInclude(r)
	resources := make([]jsonapi.Resource, 0, len(modules))
	included := []jsonapi.Resource{}
	for _, module := range modules {
		versions, err := a.deps.Store.ListModuleVersions(r.Context(), module.ID)
		if err != nil {
			if !a.writeStoreError(w, "list module versions", err) {
				a.internalError(w, "list module versions", err)
			}
			return nil, nil, false
		}
		if want["version"] {
			for _, version := range versions {
				included = append(included, moduleVersionResource(version))
			}
		}
		if want["vcs"] && module.VcsID != nil && *module.VcsID != "" {
			provider, err := a.deps.Store.GetVCS(r.Context(), *module.VcsID)
			if err != nil {
				if !a.writeStoreError(w, "get vcs provider", err) {
					a.internalError(w, "get vcs provider", err)
				}
				return nil, nil, false
			}
			resource, err := a.vcsResource(r.Context(), provider)
			if err != nil {
				if !a.writeStoreError(w, "get vcs workspaces", err) {
					a.internalError(w, "get vcs workspaces", err)
				}
				return nil, nil, false
			}
			included = append(included, resource)
		}
		if want["ssh"] && module.SshID != nil && *module.SshID != "" {
			key, err := a.deps.Store.GetSSHKey(r.Context(), *module.SshID)
			if err != nil {
				if !a.writeStoreError(w, "get ssh key", err) {
					a.internalError(w, "get ssh key", err)
				}
				return nil, nil, false
			}
			included = append(included, sshResource(key))
		}
		resources = append(resources, moduleResource(module, organizationName, versions))
	}
	return resources, included, true
}

func (a *API) createModule(w http.ResponseWriter, r *http.Request) {
	orgID := param(r, "organizationID")
	if _, err := a.deps.Store.GetOrganization(r.Context(), orgID); err != nil {
		if !a.writeStoreError(w, "get organization", err) {
			a.internalError(w, "get organization", err)
		}
		return
	}
	attrs, data, ok := decodeModuleAttrs(w, r)
	if !ok {
		return
	}
	if attrs.Name == nil || attrs.Provider == nil || attrs.Source == nil {
		dataIntegrityConflict(w)
		return
	}
	if !validModuleProvider(*attrs.Provider) {
		moduleProviderError(w, *attrs.Provider)
		return
	}
	links, ok := moduleLinks(w, data, nil)
	if !ok {
		return
	}
	if !a.validateModuleLinks(w, r, orgID, links) {
		return
	}
	email := callerEmail(r)
	module, err := a.deps.Store.CreateModule(r.Context(), store.CreateModuleParams{
		ID:               uuid.NewString(),
		Name:             *attrs.Name,
		Description:      attrs.Description,
		Provider:         *attrs.Provider,
		Source:           *attrs.Source,
		OrganizationID:   orgID,
		VcsID:            links["vcs"],
		SshID:            links["ssh"],
		DownloadQuantity: int64Ptr(0),
		CreatedBy:        email,
		UpdatedBy:        email,
		TagPrefix:        attrs.TagPrefix,
		Folder:           attrs.Folder,
		LatestVersion:    nil,
	})
	if err != nil {
		if !a.writeStoreError(w, "create module", err) {
			a.internalError(w, "create module", err)
		}
		return
	}
	organization, err := a.deps.Store.GetOrganization(r.Context(), orgID)
	if err != nil {
		if !a.writeStoreError(w, "get organization", err) {
			a.internalError(w, "get organization", err)
		}
		return
	}
	jsonapi.WriteResource(w, http.StatusCreated, moduleResource(module, organization.Name, nil), jsonapi.ParseFieldsets(r))
}

func (a *API) updateModule(w http.ResponseWriter, r *http.Request) {
	module, ok := a.moduleInOrganization(w, r)
	if !ok {
		return
	}
	attrs, data, ok := decodeModuleAttrs(w, r)
	if !ok {
		return
	}
	for name, value := range map[string]*string{
		"name": attrs.Name, "provider": attrs.Provider, "source": attrs.Source,
	} {
		if jsonapi.HasAttribute(data, name) && value == nil {
			dataIntegrityConflict(w)
			return
		}
	}
	provider := pickStringValue(data, "provider", attrs.Provider, module.Provider)
	if !validModuleProvider(provider) {
		moduleProviderError(w, provider)
		return
	}
	links, ok := moduleLinks(w, data, map[string]*string{"vcs": module.VcsID, "ssh": module.SshID})
	if !ok {
		return
	}
	if !a.validateModuleLinks(w, r, module.OrganizationID, links) {
		return
	}
	if _, err := a.deps.Store.UpdateModule(r.Context(), store.UpdateModuleParams{
		ID:            module.ID,
		Name:          pickStringValue(data, "name", attrs.Name, module.Name),
		Description:   pickString(data, "description", attrs.Description, module.Description),
		Provider:      provider,
		Source:        pickStringValue(data, "source", attrs.Source, module.Source),
		VcsID:         links["vcs"],
		SshID:         links["ssh"],
		TagPrefix:     pickString(data, "tagPrefix", attrs.TagPrefix, module.TagPrefix),
		Folder:        pickString(data, "folder", attrs.Folder, module.Folder),
		LatestVersion: module.LatestVersion,
		UpdatedBy:     callerEmail(r),
	}); err != nil {
		if !a.writeStoreError(w, "update module", err) {
			a.internalError(w, "update module", err)
		}
		return
	}
	w.WriteHeader(http.StatusNoContent)
}

func (a *API) deleteModule(w http.ResponseWriter, r *http.Request) {
	module, ok := a.moduleInOrganization(w, r)
	if !ok {
		return
	}
	if err := a.deps.Store.DeleteModule(r.Context(), module.ID); err != nil {
		if !a.writeStoreError(w, "delete module", err) {
			a.internalError(w, "delete module", err)
		}
		return
	}
	w.WriteHeader(http.StatusNoContent)
}

func (a *API) listModuleVersions(w http.ResponseWriter, r *http.Request) {
	module, ok := a.moduleInOrganization(w, r)
	if !ok {
		return
	}
	versions, err := a.deps.Store.ListModuleVersions(r.Context(), module.ID)
	if err != nil {
		if !a.writeStoreError(w, "list module versions", err) {
			a.internalError(w, "list module versions", err)
		}
		return
	}
	resources := make([]jsonapi.Resource, 0, len(versions))
	for _, version := range versions {
		resources = append(resources, moduleVersionResource(version))
	}
	jsonapi.WriteCollection(w, http.StatusOK, resources, len(resources), jsonapi.Page{}, jsonapi.ParseFieldsets(r))
}

func (a *API) getModuleVersion(w http.ResponseWriter, r *http.Request) {
	version, ok := a.moduleVersionInModule(w, r)
	if !ok {
		return
	}
	jsonapi.WriteResource(w, http.StatusOK, moduleVersionResource(version), jsonapi.ParseFieldsets(r))
}

func (a *API) updateModuleVersion(w http.ResponseWriter, r *http.Request) {
	version, ok := a.moduleVersionInModule(w, r)
	if !ok {
		return
	}
	data, err := jsonapi.DecodeRequest(r)
	if err != nil {
		jsonapi.WriteError(w, http.StatusBadRequest, "Invalid Request", err.Error())
		return
	}
	if !requireType(w, data, "module_version") {
		return
	}
	attrs, err := jsonapi.DecodeAttributes[moduleVersionAttrs](data)
	if err != nil {
		jsonapi.WriteError(w, http.StatusBadRequest, "Invalid Attributes", err.Error())
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
	if _, err := a.deps.Store.UpdateModuleVersion(r.Context(), store.UpdateModuleVersionParams{
		ID:                 version.ID,
		Version:            pickStringValue(data, "version", attrs.Version, version.Version),
		CommitInfo:         pickStringValue(data, "commit", attrs.Commit, version.CommitInfo),
		GitTag:             pickString(data, "gitTag", attrs.GitTag, version.GitTag),
		Status:             status,
		DeprecationMessage: pickString(data, "deprecationMessage", attrs.DeprecationMessage, version.DeprecationMessage),
	}); err != nil {
		if !a.writeStoreError(w, "update module version", err) {
			a.internalError(w, "update module version", err)
		}
		return
	}
	w.WriteHeader(http.StatusNoContent)
}

// moduleInOrganization resolves :moduleID and refuses a module from another
// organization.
func (a *API) moduleInOrganization(w http.ResponseWriter, r *http.Request) (store.Module, bool) {
	id := param(r, "moduleID")
	if _, err := uuid.Parse(id); err != nil {
		jsonapi.WriteError(w, http.StatusBadRequest, "Invalid Value", "Invalid value: "+id)
		return store.Module{}, false
	}
	module, err := a.deps.Store.GetModule(r.Context(), id)
	if err != nil {
		if !a.writeStoreError(w, "get module", err) {
			a.internalError(w, "get module", err)
		}
		return store.Module{}, false
	}
	if module.OrganizationID != param(r, "organizationID") {
		jsonapi.WriteError(w, http.StatusNotFound, "Not Found", "Unknown identifier ["+id+"] for module")
		return store.Module{}, false
	}
	return module, true
}

func (a *API) moduleVersionInModule(w http.ResponseWriter, r *http.Request) (store.ModuleVersion, bool) {
	module, ok := a.moduleInOrganization(w, r)
	if !ok {
		return store.ModuleVersion{}, false
	}
	id := param(r, "versionID")
	if _, err := uuid.Parse(id); err != nil {
		jsonapi.WriteError(w, http.StatusBadRequest, "Invalid Value", "Invalid value: "+id)
		return store.ModuleVersion{}, false
	}
	version, err := a.deps.Store.GetModuleVersion(r.Context(), id)
	if err != nil {
		if !a.writeStoreError(w, "get module version", err) {
			a.internalError(w, "get module version", err)
		}
		return store.ModuleVersion{}, false
	}
	if version.ModuleID != module.ID {
		jsonapi.WriteError(w, http.StatusNotFound, "Not Found", "Unknown identifier ["+id+"] for module_version")
		return store.ModuleVersion{}, false
	}
	return version, true
}

// moduleLinks reads the vcs and ssh relationships: absent keeps the current id,
// null clears it. Module create and update share it with the workspace links.
func moduleLinks(w http.ResponseWriter, data *jsonapi.RequestData, current map[string]*string) (map[string]*string, bool) {
	links := map[string]*string{"vcs": nil, "ssh": nil}
	for name, id := range current {
		links[name] = id
	}
	for _, name := range []string{"vcs", "ssh"} {
		id, present, err := jsonapi.RelationshipID(data, name)
		if err != nil {
			jsonapi.WriteError(w, http.StatusBadRequest, "Invalid Relationship", err.Error())
			return nil, false
		}
		if !present {
			continue
		}
		if id == "" {
			links[name] = nil
			continue
		}
		links[name] = &id
	}
	return links, true
}

func (a *API) validateModuleLinks(w http.ResponseWriter, r *http.Request, orgID string, links map[string]*string) bool {
	if links["vcs"] != nil && !a.checkVCSRelationship(w, r, *links["vcs"], orgID) {
		return false
	}
	if links["ssh"] != nil && !a.checkSSHRelationship(w, r, *links["ssh"], orgID) {
		return false
	}
	return true
}

func callerEmail(r *http.Request) *string {
	identity, ok := auth.IdentityFromContext(r.Context())
	if !ok || identity.Email == "" {
		return nil
	}
	return &identity.Email
}

func int64Ptr(value int64) *int64 {
	return &value
}

func decodeModuleAttrs(w http.ResponseWriter, r *http.Request) (moduleAttrs, *jsonapi.RequestData, bool) {
	data, err := jsonapi.DecodeRequest(r)
	if err != nil {
		jsonapi.WriteError(w, http.StatusBadRequest, "Invalid Request", err.Error())
		return moduleAttrs{}, nil, false
	}
	if !requireType(w, data, "module") {
		return moduleAttrs{}, nil, false
	}
	attrs, err := jsonapi.DecodeAttributes[moduleAttrs](data)
	if err != nil {
		jsonapi.WriteError(w, http.StatusBadRequest, "Invalid Attributes", err.Error())
		return moduleAttrs{}, nil, false
	}
	return attrs, data, true
}

// validModuleProvider mirrors upstream's validator on the module provider
// attribute: the registry system name, not a host or organization.
func validModuleProvider(provider string) bool {
	if len(provider) == 0 || len(provider) > 64 {
		return false
	}
	for _, r := range provider {
		switch {
		case r >= 'a' && r <= 'z', r >= 'A' && r <= 'Z', r >= '0' && r <= '9':
		default:
			return false
		}
	}
	return true
}

func moduleProviderError(w http.ResponseWriter, provider string) {
	jsonapi.WriteError(w, http.StatusBadRequest, "Invalid Value", provider+
		" is not a valid provider. Provider is the OpenTofu/Terraform module registry system (the last path segment of the module address, e.g. aws), not the repository name, and must match ^[A-Za-z0-9]{1,64}$.")
}
