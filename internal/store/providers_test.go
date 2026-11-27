package store

import (
	"testing"

	"github.com/google/uuid"

	"github.com/aminbenmansour/terrakube-go/internal/dbtest"
)

// TestProviderModuleQueries covers the provider/version/implementation and
// module/module_version round trips the registry screens need, including the
// organization-scoped joins used by include=.
func TestProviderModuleQueries(t *testing.T) {
	q := New(dbtest.Pool(t))
	ctx := t.Context()

	org, err := q.CreateOrganization(ctx, CreateOrganizationParams{
		ID: uuid.NewString(), Name: "registry-" + uuid.NewString()[:8], ExecutionMode: "remote",
	})
	if err != nil {
		t.Fatalf("create organization: %v", err)
	}

	imported := true
	namespace := "acme"
	provider, err := q.CreateProvider(ctx, CreateProviderParams{
		ID: uuid.NewString(), Name: "capture", Description: strPtr("capture"),
		OrganizationID: org.ID, Imported: &imported, RegistryNamespace: &namespace,
	})
	if err != nil {
		t.Fatalf("create provider: %v", err)
	}
	list, err := q.ListProvidersByOrganization(ctx, org.ID)
	if err != nil || len(list) != 1 || list[0].ID != provider.ID {
		t.Fatalf("list providers = %+v/%v", list, err)
	}

	version, err := q.CreateProviderVersion(ctx, CreateProviderVersionParams{
		ID: uuid.NewString(), VersionNumber: "1.2.3", Protocols: "5.0",
		ProviderID: provider.ID, Status: "active",
	})
	if err != nil {
		t.Fatalf("create version: %v", err)
	}
	implementations, err := q.ListImplementationsByOrganization(ctx, org.ID)
	if err != nil || len(implementations) != 0 {
		t.Fatalf("implementations before create = %+v/%v", implementations, err)
	}
	implementation, err := q.CreateImplementation(ctx, CreateImplementationParams{
		ID: uuid.NewString(), Os: "linux", Filename: "p.zip", DownloadUrl: "https://example.com/p.zip",
		ShasumsUrl: "https://example.com/SHA256SUMS", ShasumsSignatureUrl: "https://example.com/SHA256SUMS.sig",
		Shasum: "abc", KeyID: "key", TrustSignature: "trust", Source: "https://github.com/acme/p",
		SourceUrl: "https://github.com/acme/p", VersionID: version.ID,
	})
	if err != nil {
		t.Fatalf("create implementation: %v", err)
	}
	implementations, err = q.ListImplementationsByOrganization(ctx, org.ID)
	if err != nil || len(implementations) != 1 || implementations[0].ID != implementation.ID {
		t.Fatalf("implementations by organization = %+v/%v", implementations, err)
	}

	module, err := q.CreateModule(ctx, CreateModuleParams{
		ID: uuid.NewString(), Name: "capture", Provider: "aws", Source: "https://github.com/acme/m",
		OrganizationID: org.ID, DownloadQuantity: int64Ptr(0),
	})
	if err != nil {
		t.Fatalf("create module: %v", err)
	}
	moduleVersion, err := q.CreateModuleVersion(ctx, CreateModuleVersionParams{
		ID: uuid.NewString(), ModuleID: module.ID, Version: "1.0.0", CommitInfo: "abc", Status: "active",
	})
	if err != nil {
		t.Fatalf("create module version: %v", err)
	}
	if updated, err := q.UpdateModuleVersion(ctx, UpdateModuleVersionParams{
		ID: moduleVersion.ID, Version: "1.0.0", CommitInfo: "abc", Status: "deprecated",
		DeprecationMessage: strPtr("use 2.x"),
	}); err != nil || updated.Status != "deprecated" || updated.DeprecationMessage == nil {
		t.Fatalf("update module version = %+v/%v", updated, err)
	}

	if err := q.DeleteModule(ctx, module.ID); err != nil {
		t.Fatalf("delete module: %v", err)
	}
	if _, err := q.GetModuleVersion(ctx, moduleVersion.ID); err == nil {
		t.Fatal("module version survived the module delete")
	}

	if err := q.DeleteImplementation(ctx, implementation.ID); err != nil {
		t.Fatalf("delete implementation: %v", err)
	}
	if err := q.DeleteProviderVersion(ctx, version.ID); err != nil {
		t.Fatalf("delete version: %v", err)
	}
	if err := q.DeleteProvider(ctx, provider.ID); err != nil {
		t.Fatalf("delete provider: %v", err)
	}
}

func strPtr(value string) *string { return &value }

func int64Ptr(value int64) *int64 { return &value }
