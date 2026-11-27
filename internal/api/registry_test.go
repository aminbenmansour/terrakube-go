package api

import (
	"context"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/jackc/pgx/v5"

	"github.com/aminbenmansour/terrakube-go/internal/store"
)

const (
	registryOrgID       = "41414141-4141-4141-8141-414141414141"
	registryProviderID  = "42424242-4242-4242-8242-424242424242"
	registryVersionID   = "43434343-4343-4343-8343-434343434343"
	registryImplementID = "44444444-4444-4444-8444-444444444444"
	registryModuleID    = "45454545-4545-4545-8545-454545454545"
	registryModuleVrsID = "46464646-4646-4646-8646-464646464646"
	registryVCSID       = "47474747-4747-4747-8747-474747474747"
	registrySSHID       = "48484848-4848-4848-8848-484848484848"
	registryForeignID   = "49494949-4949-4949-8949-494949494949"
)

type registryStub struct {
	store.Querier
	org             store.Organization
	providers       []store.Provider
	versions        []store.Version
	implementations []store.Implementation
	modules         []store.Module
	moduleVersions  []store.ModuleVersion

	createdProvider   *store.CreateProviderParams
	updatedProvider   *store.UpdateProviderParams
	createdVersion    *store.CreateProviderVersionParams
	updatedVersion    *store.UpdateProviderVersionParams
	createdImplement  *store.CreateImplementationParams
	createdModule     *store.CreateModuleParams
	updatedModule     *store.UpdateModuleParams
	updatedModuleVers *store.UpdateModuleVersionParams
	deleted           []string
}

func (s *registryStub) GetOrganization(_ context.Context, id string) (store.Organization, error) {
	if id != s.org.ID {
		return store.Organization{}, pgx.ErrNoRows
	}
	return s.org, nil
}

func (s *registryStub) ListProvidersByOrganization(context.Context, string) ([]store.Provider, error) {
	return s.providers, nil
}

func (s *registryStub) GetProvider(_ context.Context, id string) (store.Provider, error) {
	for _, provider := range s.providers {
		if provider.ID == id {
			return provider, nil
		}
	}
	return store.Provider{}, pgx.ErrNoRows
}

func (s *registryStub) CreateProvider(_ context.Context, arg store.CreateProviderParams) (store.Provider, error) {
	s.createdProvider = &arg
	return store.Provider{
		ID: arg.ID, Name: arg.Name, Description: arg.Description, OrganizationID: arg.OrganizationID,
		Imported: arg.Imported, RegistryNamespace: arg.RegistryNamespace,
	}, nil
}

func (s *registryStub) UpdateProvider(_ context.Context, arg store.UpdateProviderParams) (store.Provider, error) {
	s.updatedProvider = &arg
	return store.Provider{ID: arg.ID, Name: arg.Name}, nil
}

func (s *registryStub) DeleteProvider(_ context.Context, id string) error {
	s.deleted = append(s.deleted, id)
	return nil
}

func (s *registryStub) ListProviderVersionsByOrganization(context.Context, string) ([]store.Version, error) {
	return s.versions, nil
}

func (s *registryStub) ListProviderVersions(context.Context, string) ([]store.Version, error) {
	return s.versions, nil
}

func (s *registryStub) GetProviderVersion(_ context.Context, id string) (store.Version, error) {
	for _, version := range s.versions {
		if version.ID == id {
			return version, nil
		}
	}
	return store.Version{}, pgx.ErrNoRows
}

func (s *registryStub) CreateProviderVersion(_ context.Context, arg store.CreateProviderVersionParams) (store.Version, error) {
	s.createdVersion = &arg
	return store.Version{
		ID: arg.ID, VersionNumber: arg.VersionNumber, Protocols: arg.Protocols,
		ProviderID: arg.ProviderID, Status: arg.Status, DeprecationMessage: arg.DeprecationMessage,
	}, nil
}

func (s *registryStub) UpdateProviderVersion(_ context.Context, arg store.UpdateProviderVersionParams) (store.Version, error) {
	s.updatedVersion = &arg
	return store.Version{ID: arg.ID, Status: arg.Status}, nil
}

func (s *registryStub) DeleteProviderVersion(_ context.Context, id string) error {
	s.deleted = append(s.deleted, id)
	return nil
}

func (s *registryStub) ListImplementations(context.Context, string) ([]store.Implementation, error) {
	return s.implementations, nil
}

func (s *registryStub) ListImplementationsByOrganization(context.Context, string) ([]store.Implementation, error) {
	return s.implementations, nil
}

func (s *registryStub) GetImplementation(_ context.Context, id string) (store.Implementation, error) {
	for _, implementation := range s.implementations {
		if implementation.ID == id {
			return implementation, nil
		}
	}
	return store.Implementation{}, pgx.ErrNoRows
}

func (s *registryStub) CreateImplementation(_ context.Context, arg store.CreateImplementationParams) (store.Implementation, error) {
	s.createdImplement = &arg
	return store.Implementation{
		ID: arg.ID, Os: arg.Os, Filename: arg.Filename, DownloadUrl: arg.DownloadUrl,
		ShasumsUrl: arg.ShasumsUrl, ShasumsSignatureUrl: arg.ShasumsSignatureUrl, Shasum: arg.Shasum,
		KeyID: arg.KeyID, TrustSignature: arg.TrustSignature, Source: arg.Source,
		SourceUrl: arg.SourceUrl, VersionID: arg.VersionID,
	}, nil
}

func (s *registryStub) UpdateImplementation(_ context.Context, arg store.UpdateImplementationParams) (store.Implementation, error) {
	return store.Implementation{ID: arg.ID, Os: arg.Os}, nil
}

func (s *registryStub) DeleteImplementation(_ context.Context, id string) error {
	s.deleted = append(s.deleted, id)
	return nil
}

func (s *registryStub) ListModulesByOrganization(context.Context, string) ([]store.Module, error) {
	return s.modules, nil
}

func (s *registryStub) GetModule(_ context.Context, id string) (store.Module, error) {
	for _, module := range s.modules {
		if module.ID == id {
			return module, nil
		}
	}
	return store.Module{}, pgx.ErrNoRows
}

func (s *registryStub) CreateModule(_ context.Context, arg store.CreateModuleParams) (store.Module, error) {
	s.createdModule = &arg
	return store.Module{
		ID: arg.ID, Name: arg.Name, Description: arg.Description, Provider: arg.Provider,
		Source: arg.Source, OrganizationID: arg.OrganizationID, VcsID: arg.VcsID, SshID: arg.SshID,
		CreatedBy: arg.CreatedBy, UpdatedBy: arg.UpdatedBy, TagPrefix: arg.TagPrefix, Folder: arg.Folder,
	}, nil
}

func (s *registryStub) UpdateModule(_ context.Context, arg store.UpdateModuleParams) (store.Module, error) {
	s.updatedModule = &arg
	return store.Module{ID: arg.ID, Name: arg.Name}, nil
}

func (s *registryStub) DeleteModule(_ context.Context, id string) error {
	s.deleted = append(s.deleted, id)
	return nil
}

func (s *registryStub) ListModuleVersions(context.Context, string) ([]store.ModuleVersion, error) {
	return s.moduleVersions, nil
}

func (s *registryStub) GetModuleVersion(_ context.Context, id string) (store.ModuleVersion, error) {
	for _, version := range s.moduleVersions {
		if version.ID == id {
			return version, nil
		}
	}
	return store.ModuleVersion{}, pgx.ErrNoRows
}

func (s *registryStub) UpdateModuleVersion(_ context.Context, arg store.UpdateModuleVersionParams) (store.ModuleVersion, error) {
	s.updatedModuleVers = &arg
	return store.ModuleVersion{ID: arg.ID, Status: arg.Status}, nil
}

func (s *registryStub) GetVCS(_ context.Context, id string) (store.Vc, error) {
	if id != registryVCSID {
		return store.Vc{}, pgx.ErrNoRows
	}
	return store.Vc{ID: registryVCSID, Name: "capture", OrganizationID: registryOrgID}, nil
}

func (s *registryStub) ListWorkspaceIDsByVCS(context.Context, *string) ([]string, error) {
	return nil, nil
}

func (s *registryStub) GetSSHKey(_ context.Context, id string) (store.Ssh, error) {
	if id != registrySSHID {
		return store.Ssh{}, pgx.ErrNoRows
	}
	return store.Ssh{ID: registrySSHID, Name: strPtr("capture"), OrganizationID: registryOrgID}, nil
}

func newRegistryStub() *registryStub {
	imported := false
	return &registryStub{
		org: store.Organization{ID: registryOrgID, Name: "simple-governance"},
		providers: []store.Provider{{
			ID: registryProviderID, Name: "capture", Description: strPtr("capture"),
			OrganizationID: registryOrgID, Imported: &imported,
		}},
		versions: []store.Version{{
			ID: registryVersionID, VersionNumber: "1.2.3", Protocols: "5.0",
			ProviderID: registryProviderID, Status: "active",
		}},
		implementations: []store.Implementation{{
			ID: registryImplementID, Os: "linux", Filename: "p.zip", VersionID: registryVersionID,
		}},
		modules: []store.Module{{
			ID: registryModuleID, Name: "capture", Provider: "aws", Source: "https://github.com/acme/m",
			OrganizationID: registryOrgID, VcsID: strPtr(registryVCSID), SshID: strPtr(registrySSHID),
		}},
		moduleVersions: []store.ModuleVersion{{
			ID: registryModuleVrsID, ModuleID: registryModuleID, Version: "1.0.0",
			CommitInfo: "abc", Status: "active",
		}},
	}
}

func registryServe(t *testing.T, q store.Querier, method, path, body string) *httptest.ResponseRecorder {
	t.Helper()
	req := httptest.NewRequest(method, path, strings.NewReader(body))
	if body != "" {
		req.Header.Set("Content-Type", "application/vnd.api+json")
	}
	rec := httptest.NewRecorder()
	New(Deps{Store: q}).ServeHTTP(rec, req)
	return rec
}

func decodeResource(t *testing.T, rec *httptest.ResponseRecorder) map[string]any {
	t.Helper()
	var doc map[string]any
	if err := json.Unmarshal(rec.Body.Bytes(), &doc); err != nil {
		t.Fatalf("decode: %v; body = %s", err, rec.Body.String())
	}
	if resources, ok := doc["data"].([]any); ok && len(resources) > 0 {
		resource, _ := resources[0].(map[string]any)
		return resource
	}
	resource, _ := doc["data"].(map[string]any)
	return resource
}

func TestProviderRoutes(t *testing.T) {
	base := "/organization/" + registryOrgID
	q := newRegistryStub()

	rec := registryServe(t, q, http.MethodGet, base+"/provider?include=version,version.implementation", "")
	if rec.Code != http.StatusOK {
		t.Fatalf("list status = %d: %s", rec.Code, rec.Body.String())
	}
	resource := decodeResource(t, rec)
	attrs, _ := resource["attributes"].(map[string]any)
	for _, name := range []string{"description", "imported", "name", "registryNamespace"} {
		if _, ok := attrs[name]; !ok {
			t.Errorf("provider attributes missing %q: %v", name, attrs)
		}
	}
	var doc struct {
		Included []struct {
			Type string `json:"type"`
		} `json:"included"`
	}
	if err := json.Unmarshal(rec.Body.Bytes(), &doc); err != nil {
		t.Fatalf("decode included: %v", err)
	}
	if len(doc.Included) != 2 || doc.Included[0].Type != "version" || doc.Included[1].Type != "implementation" {
		t.Fatalf("included = %+v, want a version and an implementation", doc.Included)
	}

	rec = registryServe(t, q, http.MethodGet, base+"/provider/"+registryProviderID, "")
	if rec.Code != http.StatusOK {
		t.Fatalf("get status = %d: %s", rec.Code, rec.Body.String())
	}

	rec = registryServe(t, q, http.MethodPost, base+"/provider",
		`{"data":{"type":"provider","attributes":{"name":"new","description":"d","imported":true,"registryNamespace":"acme"}}}`)
	if rec.Code != http.StatusCreated || q.createdProvider == nil || q.createdProvider.Name != "new" {
		t.Fatalf("create = %d/%+v: %s", rec.Code, q.createdProvider, rec.Body.String())
	}

	rec = registryServe(t, q, http.MethodPost, base+"/provider", `{"data":{"type":"provider","attributes":{}}}`)
	if rec.Code != http.StatusConflict {
		t.Fatalf("create without a name = %d, want 409", rec.Code)
	}
	rec = registryServe(t, q, http.MethodPost, base+"/provider", `{"data":{"type":"provider","attributes":{"name":""}}}`)
	if rec.Code != http.StatusCreated || q.createdProvider.Imported == nil || *q.createdProvider.Imported {
		t.Fatalf("create with an empty name = %d/%+v, want 201 with imported false", rec.Code, q.createdProvider)
	}

	rec = registryServe(t, q, http.MethodPatch, base+"/provider/"+registryProviderID,
		`{"data":{"type":"provider","id":"`+registryProviderID+`","attributes":{"description":"after"}}}`)
	if rec.Code != http.StatusNoContent || q.updatedProvider == nil || *q.updatedProvider.Description != "after" {
		t.Fatalf("patch = %d/%+v: %s", rec.Code, q.updatedProvider, rec.Body.String())
	}

	rec = registryServe(t, q, http.MethodDelete, base+"/provider/"+registryProviderID, "")
	if rec.Code != http.StatusNoContent || len(q.deleted) != 1 || q.deleted[0] != registryProviderID {
		t.Fatalf("delete = %d/%v", rec.Code, q.deleted)
	}

	if rec := registryServe(t, q, http.MethodGet, base+"/provider/not-a-uuid", ""); rec.Code != http.StatusBadRequest {
		t.Fatalf("bad id = %d, want 400", rec.Code)
	}
	if rec := registryServe(t, q, http.MethodGet, base+"/provider/"+registryForeignID, ""); rec.Code != http.StatusNotFound {
		t.Fatalf("unknown provider = %d, want 404", rec.Code)
	}
}

func TestProviderVersionAndImplementationRoutes(t *testing.T) {
	base := "/organization/" + registryOrgID + "/provider/" + registryProviderID
	q := newRegistryStub()

	rec := registryServe(t, q, http.MethodPost, base+"/version",
		`{"data":{"type":"version","attributes":{"versionNumber":"2.0.0","protocols":"5.0"}}}`)
	if rec.Code != http.StatusCreated || q.createdVersion == nil || q.createdVersion.Status != "active" {
		t.Fatalf("create version = %d/%+v: %s", rec.Code, q.createdVersion, rec.Body.String())
	}
	rec = registryServe(t, q, http.MethodPost, base+"/version", `{"data":{"type":"version","attributes":{"protocols":"5.0"}}}`)
	if rec.Code != http.StatusConflict {
		t.Fatalf("version without a number = %d, want 409", rec.Code)
	}

	rec = registryServe(t, q, http.MethodPatch, base+"/version/"+registryVersionID,
		`{"data":{"type":"version","id":"`+registryVersionID+`","attributes":{"status":"deprecated","deprecationMessage":"use 2.x"}}}`)
	if rec.Code != http.StatusNoContent || q.updatedVersion == nil || q.updatedVersion.Status != "deprecated" {
		t.Fatalf("patch version = %d/%+v: %s", rec.Code, q.updatedVersion, rec.Body.String())
	}
	rec = registryServe(t, q, http.MethodPatch, base+"/version/"+registryVersionID,
		`{"data":{"type":"version","id":"`+registryVersionID+`","attributes":{"status":"nonsense"}}}`)
	if rec.Code != http.StatusBadRequest || !strings.Contains(rec.Body.String(), "Invalid value: nonsense") {
		t.Fatalf("invalid status = %d/%s, want 400 Invalid value", rec.Code, rec.Body.String())
	}

	rec = registryServe(t, q, http.MethodPost, base+"/version/"+registryVersionID+"/implementation",
		`{"data":{"type":"implementation","attributes":{"os":"linux","arch":"amd64","filename":"p.zip","downloadUrl":"https://example.com/p.zip","shasumsUrl":"https://example.com/SHA256SUMS","shasumsSignatureUrl":"https://example.com/SHA256SUMS.sig","shasum":"abc","keyId":"key","trustSignature":"trust","source":"https://github.com/acme/p","sourceUrl":"https://github.com/acme/p"}}}`)
	if rec.Code != http.StatusCreated || q.createdImplement == nil || q.createdImplement.VersionID != registryVersionID {
		t.Fatalf("create implementation = %d/%+v: %s", rec.Code, q.createdImplement, rec.Body.String())
	}
	rec = registryServe(t, q, http.MethodPost, base+"/version/"+registryVersionID+"/implementation",
		`{"data":{"type":"implementation","attributes":{"os":"linux"}}}`)
	if rec.Code != http.StatusConflict {
		t.Fatalf("implementation without required attributes = %d, want 409", rec.Code)
	}

	rec = registryServe(t, q, http.MethodDelete, base+"/version/"+registryVersionID+"/implementation/"+registryImplementID, "")
	if rec.Code != http.StatusNoContent || len(q.deleted) != 1 || q.deleted[0] != registryImplementID {
		t.Fatalf("delete implementation = %d/%v", rec.Code, q.deleted)
	}
}

func TestModuleRoutes(t *testing.T) {
	base := "/organization/" + registryOrgID
	q := newRegistryStub()

	rec := registryServe(t, q, http.MethodGet, base+"/module?fields[module]=name,provider", "")
	if rec.Code != http.StatusOK {
		t.Fatalf("list status = %d: %s", rec.Code, rec.Body.String())
	}
	attrs, _ := decodeResource(t, rec)["attributes"].(map[string]any)
	if len(attrs) != 2 {
		t.Fatalf("attributes = %v, want only the requested fields", attrs)
	}

	rec = registryServe(t, q, http.MethodGet, base+"/module/"+registryModuleID+"?include=vcs,version", "")
	if rec.Code != http.StatusOK {
		t.Fatalf("get status = %d: %s", rec.Code, rec.Body.String())
	}
	attrs, _ = decodeResource(t, rec)["attributes"].(map[string]any)
	if attrs["registryPath"] != "simple-governance/capture/aws" {
		t.Fatalf("registryPath = %v", attrs["registryPath"])
	}

	rec = registryServe(t, q, http.MethodPost, base+"/module",
		`{"data":{"type":"module","attributes":{"name":"new","provider":"aws","source":"https://github.com/acme/m"},"relationships":{"vcs":{"data":{"type":"vcs","id":"`+registryVCSID+`"}}}}}`)
	if rec.Code != http.StatusCreated || q.createdModule == nil || q.createdModule.VcsID == nil {
		t.Fatalf("create module = %d/%+v: %s", rec.Code, q.createdModule, rec.Body.String())
	}
	rec = registryServe(t, q, http.MethodPost, base+"/module",
		`{"data":{"type":"module","attributes":{"name":"new"}}}`)
	if rec.Code != http.StatusConflict {
		t.Fatalf("module without provider/source = %d, want 409", rec.Code)
	}
	rec = registryServe(t, q, http.MethodPost, base+"/module",
		`{"data":{"type":"module","attributes":{"name":"new","provider":"a-b","source":"https://github.com/acme/m"}}}`)
	if rec.Code != http.StatusBadRequest || !strings.Contains(rec.Body.String(), "is not a valid provider") {
		t.Fatalf("module with an invalid provider = %d/%s, want 400", rec.Code, rec.Body.String())
	}
	rec = registryServe(t, q, http.MethodPost, base+"/module",
		`{"data":{"type":"module","attributes":{"name":"new","provider":"aws","source":"https://github.com/acme/m"},"relationships":{"vcs":{"data":{"type":"vcs","id":"`+registryForeignID+`"}}}}}`)
	if rec.Code != http.StatusNotFound {
		t.Fatalf("foreign vcs = %d, want 404", rec.Code)
	}

	rec = registryServe(t, q, http.MethodPatch, base+"/module/"+registryModuleID,
		`{"data":{"type":"module","id":"`+registryModuleID+`","attributes":{"description":"after"}}}`)
	if rec.Code != http.StatusNoContent || q.updatedModule == nil || *q.updatedModule.Description != "after" {
		t.Fatalf("patch module = %d/%+v: %s", rec.Code, q.updatedModule, rec.Body.String())
	}

	rec = registryServe(t, q, http.MethodPatch, base+"/module/"+registryModuleID+"/version/"+registryModuleVrsID,
		`{"data":{"type":"module_version","id":"`+registryModuleVrsID+`","attributes":{"status":"removed","deprecationMessage":"gone"}}}`)
	if rec.Code != http.StatusNoContent || q.updatedModuleVers == nil || q.updatedModuleVers.Status != "removed" {
		t.Fatalf("patch module version = %d/%+v: %s", rec.Code, q.updatedModuleVers, rec.Body.String())
	}
	rec = registryServe(t, q, http.MethodPatch, base+"/module/"+registryModuleID+"/version/"+registryModuleVrsID,
		`{"data":{"type":"module_version","id":"`+registryModuleVrsID+`","attributes":{"status":"nonsense"}}}`)
	if rec.Code != http.StatusBadRequest || !strings.Contains(rec.Body.String(), "Invalid value: nonsense") {
		t.Fatalf("invalid module version status = %d/%s, want 400 Invalid value", rec.Code, rec.Body.String())
	}

	rec = registryServe(t, q, http.MethodDelete, base+"/module/"+registryModuleID, "")
	if rec.Code != http.StatusNoContent || len(q.deleted) != 1 || q.deleted[0] != registryModuleID {
		t.Fatalf("delete module = %d/%v", rec.Code, q.deleted)
	}
}

// registryCaptureStore stands in for the rows the provider/module captures were
// recorded against. The ids are read from the captures, so a re-recording does
// not need this table updated.
func registryCaptureStore(captures map[string]upstreamCapture) *registryStub {
	providerID := captureID(captures, "organization-provider-create")
	versionID := captureID(captures, "organization-provider-version-create")
	implementationID := captureID(captures, "organization-provider-implementation-create")
	moduleID := captureID(captures, "organization-module-create")
	imported := true
	description := "capture"
	namespace := "acme"
	return &registryStub{
		org: store.Organization{ID: capturedRelOrg, Name: "simple-governance"},
		providers: []store.Provider{{
			ID: providerID, Name: "capture-provider", Description: &description,
			OrganizationID: capturedRelOrg, Imported: &imported, RegistryNamespace: &namespace,
		}},
		versions: []store.Version{{
			ID: versionID, VersionNumber: "1.2.3", Protocols: "5.0",
			ProviderID: providerID, Status: "deprecated",
		}},
		implementations: []store.Implementation{{
			ID: implementationID, Os: "linux", Filename: "terraform-provider-capture_1.2.3_linux_amd64.zip",
			VersionID: versionID,
		}},
		modules: []store.Module{{
			ID: moduleID, Name: "capture-module", Provider: "aws",
			Source: "https://github.com/example/terraform-aws-capture", OrganizationID: capturedRelOrg,
		}},
		moduleVersions: []store.ModuleVersion{{
			ID: "bbbbbbbb-0000-4000-8000-000000000001", ModuleID: moduleID,
			Version: "1.0.0", CommitInfo: "abc123", Status: "deprecated",
		}},
	}
}

func registryProviderStore(t *testing.T, captures map[string]upstreamCapture) *registryStub {
	t.Helper()
	providerID := captures["organization-provider-create"].resourceID(t)
	imported := true
	description := "capture"
	namespace := "acme"
	return &registryStub{
		org: store.Organization{ID: capturedRelOrg, Name: "simple-governance"},
		providers: []store.Provider{{
			ID: providerID, Name: "capture-provider", Description: &description,
			OrganizationID: capturedRelOrg, Imported: &imported, RegistryNamespace: &namespace,
		}},
	}
}

func emptyRegistryStore() *registryStub {
	return &registryStub{org: store.Organization{ID: capturedRelOrg, Name: "simple-governance"}}
}

// captureID reads a recorded resource id without a *testing.T so the replay
// table can be built while the package initialises.
func captureID(captures map[string]upstreamCapture, name string) string {
	var doc struct {
		Data struct {
			ID string `json:"id"`
		} `json:"data"`
	}
	_ = json.Unmarshal(captures[name].body, &doc)
	return doc.Data.ID
}
