package api

import (
	"context"
	"errors"
	"fmt"
	"net/http"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5"

	"github.com/aminbenmansour/terrakube-go/internal/auth"
	"github.com/aminbenmansour/terrakube-go/internal/jsonapi"
	"github.com/aminbenmansour/terrakube-go/internal/rbac"
	"github.com/aminbenmansour/terrakube-go/internal/store"
)

// maxOutboundRunTriggers mirrors io.terrakube.run-trigger.max-outbound-triggers-per-workspace's
// default. Upstream applies it when an edge is created or enabled, not at dispatch.
const maxOutboundRunTriggers = 20

type runTriggerAttrs struct {
	Enabled *bool `json:"enabled"`
}

func runTriggerResource(trigger store.WorkspaceRunTrigger) jsonapi.Resource {
	template := jsonapi.NullLink()
	if trigger.TemplateID != nil && *trigger.TemplateID != "" {
		template = jsonapi.Link("template", *trigger.TemplateID)
	}
	return jsonapi.Resource{
		Type: "runTrigger",
		ID:   trigger.ID,
		Attributes: map[string]any{
			"createdBy":   trigger.CreatedBy,
			"createdDate": timestamp(trigger.CreatedDate),
			"enabled":     trigger.Enabled,
			"updatedBy":   trigger.UpdatedBy,
			"updatedDate": timestamp(trigger.UpdatedDate),
		},
		Relationships: map[string]jsonapi.Relationship{
			"destinationWorkspace": jsonapi.Link("workspace", trigger.DestinationWorkspaceID),
			"sourceWorkspace":      jsonapi.Link("workspace", trigger.SourceWorkspaceID),
			"template":             template,
		},
	}
}

func (a *API) listRunTriggers(w http.ResponseWriter, r *http.Request) {
	q, err := parseList(r, "runTrigger")
	if err != nil {
		jsonapi.WriteError(w, http.StatusBadRequest, "Invalid Filter", err.Error())
		return
	}

	// Reads are open to every member of the edge's organization; a caller
	// outside it sees an empty collection rather than a denial.
	var allowed []string
	if !a.callerIsOwner(r) {
		identity, _ := auth.IdentityFromContext(r.Context())
		ids, err := a.deps.Store.ListOrganizationIDsByGroupNames(r.Context(), identity.Groups)
		if err != nil {
			a.internalError(w, "list run triggers", err)
			return
		}
		allowed = ids
		if allowed == nil {
			allowed = []string{}
		}
	}

	limit, offset := q.pageWindow()
	items, err := a.deps.Lists.ListRunTriggers(r.Context(), []any{allowed}, q.Filter, q.Sorts, limit, offset)
	if err != nil {
		if !a.writeStoreError(w, "list run triggers", err) {
			a.internalError(w, "list run triggers", err)
		}
		return
	}

	var total int64
	if q.Totals {
		total, err = a.deps.Lists.CountRunTriggers(r.Context(), []any{allowed}, q.Filter)
		if err != nil {
			if !a.writeStoreError(w, "count run triggers", err) {
				a.internalError(w, "count run triggers", err)
			}
			return
		}
	}

	resources := make([]jsonapi.Resource, 0, len(items))
	for _, trigger := range items {
		resources = append(resources, runTriggerResource(trigger))
	}
	included, err := a.runTriggerIncluded(r.Context(), r, items)
	if err != nil {
		if !a.writeStoreError(w, "include run trigger relationships", err) {
			a.internalError(w, "include run trigger relationships", err)
		}
		return
	}

	doc := jsonapi.Document{Data: resources, Included: included}
	if q.Paged {
		limit := max(q.PageSize, 1)
		page := map[string]any{"number": max(q.PageNumber, 1), "limit": limit}
		if q.Totals {
			page["totalPages"] = (total + int64(limit) - 1) / int64(limit)
			page["totalRecords"] = total
		}
		doc.Meta = map[string]any{"page": page}
	}
	jsonapi.Write(w, http.StatusOK, doc, jsonapi.ParseFieldsets(r))
}

func (a *API) getRunTrigger(w http.ResponseWriter, r *http.Request) {
	trigger, ok := a.runTriggerByID(w, r)
	if !ok {
		return
	}
	included, err := a.runTriggerIncluded(r.Context(), r, []store.WorkspaceRunTrigger{trigger})
	if err != nil {
		if !a.writeStoreError(w, "include run trigger relationships", err) {
			a.internalError(w, "include run trigger relationships", err)
		}
		return
	}
	jsonapi.Write(w, http.StatusOK, jsonapi.Document{Data: runTriggerResource(trigger), Included: included}, jsonapi.ParseFieldsets(r))
}

func (a *API) createRunTrigger(w http.ResponseWriter, r *http.Request) {
	data, err := jsonapi.DecodeRequest(r)
	if err != nil {
		jsonapi.WriteError(w, http.StatusBadRequest, "Invalid Request", err.Error())
		return
	}
	if !requireType(w, data, "runTrigger") {
		return
	}
	attrs, err := jsonapi.DecodeAttributes[runTriggerAttrs](data)
	if err != nil {
		jsonapi.WriteError(w, http.StatusBadRequest, "Invalid Attributes", err.Error())
		return
	}

	sourceID, sourcePresent, err := jsonapi.RelationshipID(data, "sourceWorkspace")
	if err != nil {
		jsonapi.WriteError(w, http.StatusBadRequest, "Invalid Request", err.Error())
		return
	}
	destinationID, destinationPresent, err := jsonapi.RelationshipID(data, "destinationWorkspace")
	if err != nil {
		jsonapi.WriteError(w, http.StatusBadRequest, "Invalid Request", err.Error())
		return
	}
	if !sourcePresent || !destinationPresent || sourceID == "" || destinationID == "" {
		writePermissionDenied(w, r.Method)
		return
	}

	source, destination, ok := a.validTriggerEnds(w, r, sourceID, destinationID)
	if !ok {
		return
	}
	if !a.canManageTriggerDestination(w, r, destination) || !a.canViewTriggerSource(w, r, source) {
		return
	}

	templateID, ok := a.triggerTemplate(w, r, data, destination.OrganizationID)
	if !ok {
		return
	}
	enabled := boolValue(attrs.Enabled, true)

	if err := a.validateRunTriggerGraph(r.Context(), destination.OrganizationID, "", sourceID, destinationID, enabled); err != nil {
		jsonapi.WriteError(w, http.StatusBadRequest, "Bad Request", err.Error())
		return
	}

	trigger, err := a.deps.Store.CreateRunTrigger(r.Context(), store.CreateRunTriggerParams{
		ID:                     uuid.NewString(),
		SourceWorkspaceID:      sourceID,
		DestinationWorkspaceID: destinationID,
		TemplateID:             templateID,
		OrganizationID:         destination.OrganizationID,
		Enabled:                enabled,
	})
	if err != nil {
		if !a.writeStoreError(w, "create run trigger", err) {
			a.internalError(w, "create run trigger", err)
		}
		return
	}
	jsonapi.WriteResource(w, http.StatusCreated, runTriggerResource(trigger), jsonapi.ParseFieldsets(r))
}

func (a *API) updateRunTrigger(w http.ResponseWriter, r *http.Request) {
	trigger, ok := a.runTriggerByID(w, r)
	if !ok {
		return
	}
	destination, ok := a.triggerDestination(w, r, trigger, http.MethodPatch)
	if !ok {
		return
	}

	data, err := jsonapi.DecodeRequest(r)
	if err != nil {
		jsonapi.WriteError(w, http.StatusBadRequest, "Invalid Request", err.Error())
		return
	}
	if !requireType(w, data, "runTrigger") {
		return
	}
	attrs, err := jsonapi.DecodeAttributes[runTriggerAttrs](data)
	if err != nil {
		jsonapi.WriteError(w, http.StatusBadRequest, "Invalid Attributes", err.Error())
		return
	}

	sourceID, destinationID := trigger.SourceWorkspaceID, trigger.DestinationWorkspaceID
	topologyChanged := false
	for name, target := range map[string]*string{
		"sourceWorkspace":      &sourceID,
		"destinationWorkspace": &destinationID,
	} {
		id, present, err := jsonapi.RelationshipID(data, name)
		if err != nil {
			jsonapi.WriteError(w, http.StatusBadRequest, "Invalid Request", err.Error())
			return
		}
		if !present {
			continue
		}
		// Repointing an edge is deleting one dependency and declaring
		// another, so only the instance owner may do it.
		if !a.callerIsOwner(r) {
			writePermissionDenied(w, r.Method)
			return
		}
		if id == "" {
			writePermissionDenied(w, r.Method)
			return
		}
		if *target != id {
			*target = id
			topologyChanged = true
		}
	}

	source, destination, ok := a.validTriggerEnds(w, r, sourceID, destinationID)
	if !ok {
		return
	}
	if !a.canManageTriggerDestination(w, r, destination) || !a.canViewTriggerSource(w, r, source) {
		return
	}

	templateID := trigger.TemplateID
	template, templatePresent, err := jsonapi.RelationshipID(data, "template")
	if err != nil {
		jsonapi.WriteError(w, http.StatusBadRequest, "Invalid Request", err.Error())
		return
	}
	if templatePresent {
		templateID = nil
		if template != "" {
			if _, ok := a.triggerTemplate(w, r, data, destination.OrganizationID); !ok {
				return
			}
			templateID = &template
		}
	}

	enabled := pickBool(data, "enabled", attrs.Enabled, trigger.Enabled)
	if topologyChanged || (enabled && !trigger.Enabled) {
		if err := a.validateRunTriggerGraph(r.Context(), destination.OrganizationID, trigger.ID, sourceID, destinationID, enabled); err != nil {
			jsonapi.WriteError(w, http.StatusBadRequest, "Bad Request", err.Error())
			return
		}
	}

	if _, err := a.deps.Store.UpdateRunTrigger(r.Context(), store.UpdateRunTriggerParams{
		ID:                     trigger.ID,
		SourceWorkspaceID:      sourceID,
		DestinationWorkspaceID: destinationID,
		TemplateID:             templateID,
		OrganizationID:         destination.OrganizationID,
		Enabled:                enabled,
	}); err != nil {
		if !a.writeStoreError(w, "update run trigger", err) {
			a.internalError(w, "update run trigger", err)
		}
		return
	}
	w.WriteHeader(http.StatusNoContent)
}

func (a *API) deleteRunTrigger(w http.ResponseWriter, r *http.Request) {
	// Upstream's Elide refuses a DELETE that declares a content type, because
	// the request carries no document; the Run Triggers page omits the header.
	if r.Header.Get("Content-Type") != "" {
		jsonapi.WriteError(w, http.StatusBadRequest, "Invalid Request", "A DELETE request must not declare a content type.")
		return
	}
	trigger, ok := a.runTriggerByID(w, r)
	if !ok {
		return
	}
	if _, ok := a.triggerDestination(w, r, trigger, http.MethodDelete); !ok {
		return
	}
	if err := a.deps.Store.DeleteRunTrigger(r.Context(), trigger.ID); err != nil {
		if !a.writeStoreError(w, "delete run trigger", err) {
			a.internalError(w, "delete run trigger", err)
		}
		return
	}
	w.WriteHeader(http.StatusNoContent)
}

// runTriggerByID loads the edge and answers a caller outside its organization
// with the same denial Elide's read check raises.
func (a *API) runTriggerByID(w http.ResponseWriter, r *http.Request) (store.WorkspaceRunTrigger, bool) {
	trigger, err := a.deps.Store.GetRunTrigger(r.Context(), param(r, "runTriggerID"))
	if err != nil {
		if !a.writeStoreError(w, "get run trigger", err) {
			a.internalError(w, "get run trigger", err)
		}
		return store.WorkspaceRunTrigger{}, false
	}
	if !a.callerIsOwner(r) {
		identity, _ := auth.IdentityFromContext(r.Context())
		access, err := rbac.ResolveOrganization(r.Context(), a.deps.Store, trigger.OrganizationID, identity.Groups)
		if err != nil {
			a.internalError(w, "resolve run trigger access", err)
			return store.WorkspaceRunTrigger{}, false
		}
		if !access.Member {
			writePermissionDenied(w, r.Method)
			return store.WorkspaceRunTrigger{}, false
		}
	}
	return trigger, true
}

// triggerDestination resolves manage rights on the destination. A missing
// workspace is treated as a denial, as upstream's structural check does.
func (a *API) triggerDestination(w http.ResponseWriter, r *http.Request, trigger store.WorkspaceRunTrigger, method string) (store.Workspace, bool) {
	destination, err := a.deps.Store.GetWorkspace(r.Context(), trigger.DestinationWorkspaceID)
	if errors.Is(err, pgx.ErrNoRows) {
		writePermissionDenied(w, method)
		return store.Workspace{}, false
	}
	if err != nil {
		a.internalError(w, "get trigger destination", err)
		return store.Workspace{}, false
	}
	access, err := a.resolveWorkspaceAuthorization(r, destination)
	if err != nil {
		a.internalError(w, "resolve trigger destination access", err)
		return store.Workspace{}, false
	}
	if !access.Member || !(access.Owner || access.Permissions.ManageWorkspace) {
		writePermissionDenied(w, method)
		return store.Workspace{}, false
	}
	return destination, true
}

// validTriggerEnds loads both workspaces and enforces the invariants that hold
// regardless of the caller: distinct ends, one organization.
func (a *API) validTriggerEnds(w http.ResponseWriter, r *http.Request, sourceID, destinationID string) (store.Workspace, store.Workspace, bool) {
	if sourceID == destinationID {
		writePermissionDenied(w, r.Method)
		return store.Workspace{}, store.Workspace{}, false
	}
	source, err := a.deps.Store.GetWorkspace(r.Context(), sourceID)
	if errors.Is(err, pgx.ErrNoRows) {
		writePermissionDenied(w, r.Method)
		return store.Workspace{}, store.Workspace{}, false
	}
	if err != nil {
		a.internalError(w, "get trigger source", err)
		return store.Workspace{}, store.Workspace{}, false
	}
	destination, err := a.deps.Store.GetWorkspace(r.Context(), destinationID)
	if errors.Is(err, pgx.ErrNoRows) {
		writePermissionDenied(w, r.Method)
		return store.Workspace{}, store.Workspace{}, false
	}
	if err != nil {
		a.internalError(w, "get trigger destination", err)
		return store.Workspace{}, store.Workspace{}, false
	}
	if source.OrganizationID != destination.OrganizationID {
		writePermissionDenied(w, r.Method)
		return store.Workspace{}, store.Workspace{}, false
	}
	return source, destination, true
}

func (a *API) canManageTriggerDestination(w http.ResponseWriter, r *http.Request, destination store.Workspace) bool {
	access, err := a.resolveWorkspaceAuthorization(r, destination)
	if err != nil {
		a.internalError(w, "resolve trigger destination access", err)
		return false
	}
	if !access.Member || !(access.Owner || access.Permissions.ManageWorkspace) {
		writePermissionDenied(w, r.Method)
		return false
	}
	return true
}

func (a *API) canViewTriggerSource(w http.ResponseWriter, r *http.Request, source store.Workspace) bool {
	access, err := a.resolveWorkspaceAuthorization(r, source)
	if err != nil {
		a.internalError(w, "resolve trigger source access", err)
		return false
	}
	if !access.Member {
		writePermissionDenied(w, r.Method)
		return false
	}
	return true
}

// triggerTemplate validates the optional template override. A template from
// another organization would run its TCL against these workspaces, so it is
// refused the same way an unknown one is.
func (a *API) triggerTemplate(w http.ResponseWriter, r *http.Request, data *jsonapi.RequestData, organizationID string) (*string, bool) {
	id, present, err := jsonapi.RelationshipID(data, "template")
	if err != nil {
		jsonapi.WriteError(w, http.StatusBadRequest, "Invalid Request", err.Error())
		return nil, false
	}
	if !present || id == "" {
		return nil, true
	}
	template, err := a.deps.Store.GetTemplate(r.Context(), id)
	if errors.Is(err, pgx.ErrNoRows) {
		writePermissionDenied(w, r.Method)
		return nil, false
	}
	if err != nil {
		a.internalError(w, "get trigger template", err)
		return nil, false
	}
	if template.OrganizationID != organizationID {
		writePermissionDenied(w, r.Method)
		return nil, false
	}
	return &id, true
}

// validateRunTriggerGraph keeps the organization's edges acyclic and bounds a
// source workspace's enabled fan-out. Cycle checks walk the whole graph, so a
// loop is refused before it can re-trigger itself.
//
// shortcut: the checks are not serialized in one transaction, so two
// concurrent creations can each pass and together close a loop. Dispatch's
// cascade depth is the runtime net; upgrade if triggers become hot.
func (a *API) validateRunTriggerGraph(ctx context.Context, organizationID, triggerID, sourceID, destinationID string, enabled bool) error {
	edges, err := a.deps.Store.ListRunTriggersByOrganization(ctx, organizationID)
	if err != nil {
		return fmt.Errorf("load run trigger graph: %w", err)
	}
	adjacency := make(map[string][]string, len(edges))
	for _, edge := range edges {
		if edge.ID == triggerID {
			continue
		}
		adjacency[edge.SourceWorkspaceID] = append(adjacency[edge.SourceWorkspaceID], edge.DestinationWorkspaceID)
	}
	if path := findTriggerPath(adjacency, destinationID, sourceID); path != nil {
		return fmt.Errorf("This run trigger would create a circular dependency: %s", formatTriggerPath(path, destinationID))
	}

	if !enabled {
		return nil
	}
	count, err := a.deps.Store.CountEnabledRunTriggersBySource(ctx, store.CountEnabledRunTriggersBySourceParams{
		SourceWorkspaceID: sourceID,
		ID:                triggerID,
	})
	if err != nil {
		return fmt.Errorf("count outbound run triggers: %w", err)
	}
	if count >= maxOutboundRunTriggers {
		return fmt.Errorf("This workspace would have %d enabled run trigger(s), above the configured limit of %d. Disable or delete an existing one before adding another.",
			count+1, maxOutboundRunTriggers)
	}
	return nil
}

// findTriggerPath is an iterative DFS from one workspace to another, so a deep
// chain cannot overflow the stack.
func findTriggerPath(adjacency map[string][]string, from, target string) []string {
	cameFrom := map[string]string{from: ""}
	stack := []string{from}
	for len(stack) > 0 {
		current := stack[len(stack)-1]
		stack = stack[:len(stack)-1]
		if current == target {
			path := []string{target}
			for node := target; node != from; {
				node = cameFrom[node]
				path = append(path, node)
			}
			for i, j := 0, len(path)-1; i < j; i, j = i+1, j-1 {
				path[i], path[j] = path[j], path[i]
			}
			return path
		}
		for _, next := range adjacency[current] {
			if _, seen := cameFrom[next]; seen {
				continue
			}
			cameFrom[next] = current
			stack = append(stack, next)
		}
	}
	return nil
}

func formatTriggerPath(path []string, destinationID string) string {
	rendered := ""
	for _, node := range path {
		rendered += node + " -> "
	}
	return rendered + destinationID
}

// runTriggerIncluded resolves the ends and the optional template the Run
// Triggers page asks for in one include block.
func (a *API) runTriggerIncluded(ctx context.Context, r *http.Request, triggers []store.WorkspaceRunTrigger) ([]jsonapi.Resource, error) {
	want := parseInclude(r)
	if !want["sourceWorkspace"] && !want["destinationWorkspace"] && !want["template"] {
		return nil, nil
	}

	workspaces := map[string]bool{}
	templates := map[string]bool{}
	for _, trigger := range triggers {
		if want["sourceWorkspace"] {
			workspaces[trigger.SourceWorkspaceID] = true
		}
		if want["destinationWorkspace"] {
			workspaces[trigger.DestinationWorkspaceID] = true
		}
		if want["template"] && trigger.TemplateID != nil && *trigger.TemplateID != "" {
			templates[*trigger.TemplateID] = true
		}
	}

	included := make([]jsonapi.Resource, 0, len(workspaces)+len(templates))
	for id := range workspaces {
		workspace, err := a.deps.Store.GetWorkspace(ctx, id)
		if errors.Is(err, pgx.ErrNoRows) {
			// A soft-deleted end is hidden by GetWorkspace but still referenced
			// by the edge; leave it out rather than failing the whole page.
			continue
		}
		if err != nil {
			return nil, err
		}
		included = append(included, workspaceResource(workspace))
	}
	for id := range templates {
		template, err := a.deps.Store.GetTemplate(ctx, id)
		if err != nil {
			return nil, err
		}
		included = append(included, templateResource(template))
	}
	return included, nil
}
