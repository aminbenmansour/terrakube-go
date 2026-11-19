package tfstate

import (
	"bytes"
	"context"
	"errors"
	"net/http"
	"testing"

	"github.com/aminbenmansour/terrakube-go/internal/store"
)

func TestArchiveUploadRejectsAnOversizedBody(t *testing.T) {
	f := newFixture(t, 1024)
	f.store.archives["archive-big"] = store.GetArchiveRow{
		ID: "archive-big", HistoryID: "history-big", WorkspaceID: testWS1, OrganizationID: testOrg,
	}

	body := bytes.Repeat([]byte("x"), 2048)
	rec := f.request(t, http.MethodPut, "/archive/archive-big/terraform.tfstate",
		bytes.NewReader(body), map[string]string{"Authorization": ""})

	if rec.Code != http.StatusRequestEntityTooLarge {
		t.Fatalf("oversized archive upload = %d, want 413", rec.Code)
	}
	if _, ok := f.store.archives["archive-big"]; !ok {
		t.Fatal("archive row was deleted after a rejected upload")
	}
	if len(f.storage.objects) != 0 {
		t.Fatalf("storage wrote %d objects for a rejected upload", len(f.storage.objects))
	}
}

func TestRollbackKeepsNoHistoryRowWhenTheStateWriteFails(t *testing.T) {
	f := newFixture(t, 1024)
	jsonState := []byte(`{"version":4,"serial":2}`)
	rawState := []byte(`{"version":4,"serial":2,"raw":true}`)
	for key, data := range map[string][]byte{
		"state/" + testOrg + "/" + testWS1 + "/state/history-1.json":     jsonState,
		"state/" + testOrg + "/" + testWS1 + "/state/history-1.raw.json": rawState,
	} {
		if err := f.storage.Put(context.Background(), key, bytes.NewReader(data), int64(len(data)), "application/json"); err != nil {
			t.Fatalf("seed %s: %v", key, err)
		}
	}

	f.storage.failPut = func(string) error { return errors.New("storage unavailable") }
	rec := f.request(t, http.MethodPut, "/organization/"+testOrg+"/workspace/"+testWS1+"/rollback/history-1.json", nil, nil)

	if rec.Code != http.StatusInternalServerError {
		t.Fatalf("rollback with a failing write = %d, want 500", rec.Code)
	}
	if len(f.store.history) != 0 {
		t.Fatalf("history rows = %d, want 0; a row must not exist without its state objects", len(f.store.history))
	}
}
