package executor

import (
	"bytes"
	"context"
	"crypto/md5"
	"encoding/hex"
	"encoding/json"
	"errors"
	"fmt"
	"io"

	"github.com/google/uuid"

	"github.com/aminbenmansour/terrakube-go/internal/storage"
	"github.com/aminbenmansour/terrakube-go/internal/store"
)

type stateMetadata struct {
	Serial  int32  `json:"serial"`
	Lineage string `json:"lineage"`
}

// archiveHistory records the state written by an apply or destroy. It mirrors
// the history rows the upstream API creates for state versions.
func (e *Executor) archiveHistory(ctx context.Context, job store.Job, ws store.Workspace) error {
	if e.deps.Storage == nil {
		return nil
	}
	key := e.deps.Layout.CurrentState(job.OrganizationID, job.WorkspaceID)
	rc, err := e.deps.Storage.Get(ctx, key)
	if errors.Is(err, storage.ErrNotFound) {
		return nil
	}
	if err != nil {
		return err
	}
	defer rc.Close()

	data, err := io.ReadAll(io.LimitReader(rc, 100<<20))
	if err != nil {
		return err
	}
	if len(data) == 0 {
		return nil
	}

	var metadata stateMetadata
	if err := json.Unmarshal(data, &metadata); err != nil {
		return fmt.Errorf("parse state metadata: %w", err)
	}
	lineage := metadata.Lineage
	if lineage == "" {
		lineage = "0"
	}
	sum := md5.Sum(data)
	historyID := uuid.NewString()
	jobReference := fmt.Sprintf("%d", job.ID)

	if _, err := e.deps.Store.CreateHistory(ctx, store.CreateHistoryParams{
		ID:           historyID,
		WorkspaceID:  ws.ID,
		JobReference: &jobReference,
		Output:       nil,
		Serial:       metadata.Serial,
		Md5:          hex.EncodeToString(sum[:]),
		Lineage:      lineage,
		CreatedBy:    stringPointer("executor"),
	}); err != nil {
		return fmt.Errorf("create history row: %w", err)
	}

	if err := e.deps.Storage.Put(ctx, e.deps.Layout.JSONState(job.OrganizationID, job.WorkspaceID, historyID),
		bytes.NewReader(data), int64(len(data)), "application/json"); err != nil {
		return fmt.Errorf("archive json state: %w", err)
	}
	if err := e.deps.Storage.Put(ctx, e.deps.Layout.RawState(job.OrganizationID, job.WorkspaceID, historyID),
		bytes.NewReader(data), int64(len(data)), "application/json"); err != nil {
		return fmt.Errorf("archive raw state: %w", err)
	}
	return nil
}
