package tfstate

import (
	"fmt"

	"github.com/aminbenmansour/terrakube-go/internal/config"
)

// Layout builds object-store keys using the same paths as the upstream
// storage implementations. Local and S3 use different prefixes.
type Layout struct {
	s3 bool
}

// NewLayout returns the key layout for a configured storage backend.
func NewLayout(storageType string) Layout {
	return Layout{s3: storageType == config.StorageTypeS3}
}

// CurrentState is the object holding the latest Terraform state.
func (l Layout) CurrentState(organizationID, workspaceID string) string {
	prefix := "backend"
	if l.s3 {
		prefix = "tfstate"
	}
	return fmt.Sprintf("%s/%s/%s/terraform.tfstate", prefix, organizationID, workspaceID)
}

// RawState is the raw Terraform state archived for a history row.
func (l Layout) RawState(organizationID, workspaceID, historyID string) string {
	return fmt.Sprintf("%s/%s/%s/state/%s.raw.json", l.statePrefix(), organizationID, workspaceID, historyID)
}

// JSONState is the JSON state archived for a history row.
func (l Layout) JSONState(organizationID, workspaceID, historyID string) string {
	return fmt.Sprintf("%s/%s/%s/state/%s.json", l.statePrefix(), organizationID, workspaceID, historyID)
}

// Plan is the saved Terraform plan for a job step.
func (l Layout) Plan(organizationID, workspaceID, jobID, stepID string) string {
	return fmt.Sprintf("%s/%s/%s/%s/%s/terraformLibrary.tfPlan", l.statePrefix(), organizationID, workspaceID, jobID, stepID)
}

// Output is the archived console log for a job step. Upstream uses the same
// S3 key; its local backend writes to a home-directory path instead.
func (l Layout) Output(organizationID string, jobID int64, stepID string) string {
	return fmt.Sprintf("tfoutput/%s/%d/%s.tfoutput", organizationID, jobID, stepID)
}

// Content is the uploaded Terraform configuration tarball for a
// configuration version. Upstream uses the same key for local and S3 storage.
func (l Layout) Content(contentID string) string {
	return fmt.Sprintf("content/%s/terraformContent.tar.gz", contentID)
}

func (l Layout) statePrefix() string {
	if l.s3 {
		return "tfstate"
	}
	return "state"
}
