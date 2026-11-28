package api

import (
	"context"
	"net/http"
	"time"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5/pgtype"

	"github.com/aminbenmansour/terrakube-go/internal/auth"
	"github.com/aminbenmansour/terrakube-go/internal/jsonapi"
	"github.com/aminbenmansour/terrakube-go/internal/store"
)

// policyEnforcementLevels is upstream's PolicyEnforcementLevel enum.
var policyEnforcementLevels = map[string]bool{
	"ADVISORY": true, "SOFT_MANDATORY": true, "HARD_MANDATORY": true,
}

type policySetAttrs struct {
	Name                   *string `json:"name"`
	Description            *string `json:"description"`
	EnforcementLevel       *string `json:"enforcementLevel"`
	ShadowEnforcementLevel *string `json:"shadowEnforcementLevel"`
	OverrideTeam           *string `json:"overrideTeam"`
	Global                 *bool   `json:"global"`
	Repository             *string `json:"repository"`
	Branch                 *string `json:"branch"`
	Folder                 *string `json:"folder"`
	OpaVersion             *string `json:"opaVersion"`
}

type policySetParameterAttrs struct {
	Key         *string `json:"key"`
	Value       *string `json:"value"`
	Description *string `json:"description"`
}

type policyExemptionAttrs struct {
	RuleID          *string    `json:"ruleId"`
	TicketReference *string    `json:"ticketReference"`
	Justification   *string    `json:"justification"`
	ExpiresAt       *time.Time `json:"expiresAt"`
}

// policySetChildren are the rows upstream serialises as relationships on a
// policy set, plus the rows an include asks for.
type policySetChildren struct {
	attachments []store.PolicyAttachment
	parameters  []store.PolicySetParameter
	exemptions  []store.PolicyExemption
}

func policySetResource(set store.PolicySet, children policySetChildren) jsonapi.Resource {
	attachmentIDs := make([]string, 0, len(children.attachments))
	for _, attachment := range children.attachments {
		attachmentIDs = append(attachmentIDs, attachment.ID)
	}
	parameterIDs := make([]string, 0, len(children.parameters))
	for _, parameter := range children.parameters {
		parameterIDs = append(parameterIDs, parameter.ID)
	}
	exemptionIDs := make([]string, 0, len(children.exemptions))
	for _, exemption := range children.exemptions {
		exemptionIDs = append(exemptionIDs, exemption.ID)
	}
	return jsonapi.Resource{
		Type: "policy_set",
		ID:   set.ID,
		Attributes: map[string]any{
			"branch":                 set.Branch,
			"createdBy":              set.CreatedBy,
			"createdDate":            timestamp(set.CreatedDate),
			"description":            set.Description,
			"enforcementLevel":       set.EnforcementLevel,
			"folder":                 set.Folder,
			"global":                 set.Global,
			"name":                   set.Name,
			"opaVersion":             set.OpaVersion,
			"overrideTeam":           set.OverrideTeam,
			"repository":             set.Repository,
			"shadowEnforcementLevel": set.ShadowEnforcementLevel,
			"updatedBy":              set.UpdatedBy,
			"updatedDate":            timestamp(set.UpdatedDate),
		},
		Relationships: map[string]jsonapi.Relationship{
			"attachments":               jsonapi.LinkList("policy_attachment", attachmentIDs...),
			"exemptions":                jsonapi.LinkList("policy_exemption", exemptionIDs...),
			"notificationConfiguration": optionalLink("notification_configuration", set.NotificationConfigurationID),
			"organization":              jsonapi.Link("organization", set.OrganizationID),
			"parameters":                jsonapi.LinkList("policy_set_parameter", parameterIDs...),
			"vcs":                       optionalLink("vcs", set.VcsID),
		},
	}
}

func policyAttachmentResource(attachment store.PolicyAttachment) jsonapi.Resource {
	return jsonapi.Resource{
		Type: "policy_attachment",
		ID:   attachment.ID,
		Attributes: map[string]any{
			"createdBy":   attachment.CreatedBy,
			"createdDate": timestamp(attachment.CreatedDate),
			"updatedBy":   attachment.UpdatedBy,
			"updatedDate": timestamp(attachment.UpdatedDate),
		},
		Relationships: map[string]jsonapi.Relationship{
			"policySet": jsonapi.Link("policy_set", attachment.PolicySetID),
			"project":   optionalLink("project", attachment.ProjectID),
			"tag":       optionalLink("tag", attachment.TagID),
			"workspace": optionalLink("workspace", attachment.WorkspaceID),
		},
	}
}

func policySetParameterResource(parameter store.PolicySetParameter) jsonapi.Resource {
	return jsonapi.Resource{
		Type: "policy_set_parameter",
		ID:   parameter.ID,
		Attributes: map[string]any{
			"createdBy":   parameter.CreatedBy,
			"createdDate": timestamp(parameter.CreatedDate),
			"description": parameter.Description,
			"key":         parameter.ParameterKey,
			"updatedBy":   parameter.UpdatedBy,
			"updatedDate": timestamp(parameter.UpdatedDate),
			"value":       parameter.ParameterValue,
		},
		Relationships: map[string]jsonapi.Relationship{
			"policySet": jsonapi.Link("policy_set", parameter.PolicySetID),
		},
	}
}

func policyExemptionResource(exemption store.PolicyExemption) jsonapi.Resource {
	return jsonapi.Resource{
		Type: "policy_exemption",
		ID:   exemption.ID,
		Attributes: map[string]any{
			"createdBy":       exemption.CreatedBy,
			"createdDate":     timestamp(exemption.CreatedDate),
			"expiresAt":       timestamp(exemption.ExpiresAt),
			"justification":   exemption.Justification,
			"ruleId":          exemption.RuleID,
			"ticketReference": exemption.TicketReference,
			"updatedBy":       exemption.UpdatedBy,
			"updatedDate":     timestamp(exemption.UpdatedDate),
		},
		Relationships: map[string]jsonapi.Relationship{
			"organization": jsonapi.Link("organization", exemption.OrganizationID),
			"policySet":    jsonapi.Link("policy_set", exemption.PolicySetID),
			"project":      optionalLink("project", exemption.ProjectID),
			"workspace":    optionalLink("workspace", exemption.WorkspaceID),
		},
	}
}

// listPolicySets is the root collection the UI reads. A non-owner sees only the
// organizations their groups belong to, like the organization list.
func (a *API) listPolicySets(w http.ResponseWriter, r *http.Request) {
	var (
		sets []store.PolicySet
		err  error
	)
	if a.callerIsOwner(r) {
		sets, err = a.deps.Store.ListPolicySets(r.Context())
	} else {
		identity, _ := auth.IdentityFromContext(r.Context())
		ids, listErr := a.deps.Store.ListOrganizationIDsByGroupNames(r.Context(), identity.Groups)
		if listErr != nil {
			a.internalError(w, "list policy sets", listErr)
			return
		}
		if len(ids) > 0 {
			sets, err = a.deps.Store.ListPolicySetsByOrganizations(r.Context(), ids)
		}
	}
	if err != nil {
		if !a.writeStoreError(w, "list policy sets", err) {
			a.internalError(w, "list policy sets", err)
		}
		return
	}
	a.writePolicySets(w, r, sets)
}

func (a *API) listOrganizationPolicySets(w http.ResponseWriter, r *http.Request) {
	orgID := param(r, "organizationID")
	if _, err := a.deps.Store.GetOrganization(r.Context(), orgID); err != nil {
		if !a.writeStoreError(w, "get organization", err) {
			a.internalError(w, "get organization", err)
		}
		return
	}
	sets, err := a.deps.Store.ListPolicySetsByOrganization(r.Context(), orgID)
	if err != nil {
		if !a.writeStoreError(w, "list policy sets", err) {
			a.internalError(w, "list policy sets", err)
		}
		return
	}
	a.writePolicySets(w, r, sets)
}

func (a *API) writePolicySets(w http.ResponseWriter, r *http.Request, sets []store.PolicySet) {
	resources := make([]jsonapi.Resource, 0, len(sets))
	included := []jsonapi.Resource{}
	for _, set := range sets {
		resource, extras, ok := a.policySetDocument(w, r, set)
		if !ok {
			return
		}
		resources = append(resources, resource)
		included = append(included, extras...)
	}
	jsonapi.Write(w, http.StatusOK, jsonapi.Document{Data: resources, Included: included}, jsonapi.ParseFieldsets(r))
}

func (a *API) getPolicySet(w http.ResponseWriter, r *http.Request) {
	set, ok := a.policySetFromScope(w, r)
	if !ok {
		return
	}
	resource, included, ok := a.policySetDocument(w, r, set)
	if !ok {
		return
	}
	jsonapi.Write(w, http.StatusOK, jsonapi.Document{Data: resource, Included: included}, jsonapi.ParseFieldsets(r))
}

// policySetDocument builds one policy set and the resources its include names.
func (a *API) policySetDocument(w http.ResponseWriter, r *http.Request, set store.PolicySet) (jsonapi.Resource, []jsonapi.Resource, bool) {
	children, err := a.policySetChildren(r.Context(), set.ID)
	if err != nil {
		if !a.writeStoreError(w, "load policy set children", err) {
			a.internalError(w, "load policy set children", err)
		}
		return jsonapi.Resource{}, nil, false
	}
	want := parseInclude(r)
	included := []jsonapi.Resource{}
	if want["attachments"] {
		for _, attachment := range children.attachments {
			included = append(included, policyAttachmentResource(attachment))
		}
	}
	if want["parameters"] {
		for _, parameter := range children.parameters {
			included = append(included, policySetParameterResource(parameter))
		}
	}
	if want["exemptions"] {
		for _, exemption := range children.exemptions {
			included = append(included, policyExemptionResource(exemption))
		}
	}
	if want["vcs"] && set.VcsID != nil && *set.VcsID != "" {
		provider, err := a.deps.Store.GetVCS(r.Context(), *set.VcsID)
		if err != nil {
			if !a.writeStoreError(w, "get vcs provider", err) {
				a.internalError(w, "get vcs provider", err)
			}
			return jsonapi.Resource{}, nil, false
		}
		resource, err := a.vcsResource(r.Context(), provider)
		if err != nil {
			a.internalError(w, "get vcs workspaces", err)
			return jsonapi.Resource{}, nil, false
		}
		included = append(included, resource)
	}
	if want["notificationConfiguration"] && set.NotificationConfigurationID != nil && *set.NotificationConfigurationID != "" {
		config, err := a.deps.Store.GetNotificationConfiguration(r.Context(), *set.NotificationConfigurationID)
		if err != nil {
			if !a.writeStoreError(w, "get notification configuration", err) {
				a.internalError(w, "get notification configuration", err)
			}
			return jsonapi.Resource{}, nil, false
		}
		included = append(included, notificationConfigurationResource(config))
	}
	return policySetResource(set, children), included, true
}

func (a *API) policySetChildren(ctx context.Context, id string) (policySetChildren, error) {
	attachments, err := a.deps.Store.ListPolicyAttachments(ctx, id)
	if err != nil {
		return policySetChildren{}, err
	}
	parameters, err := a.deps.Store.ListPolicySetParameters(ctx, id)
	if err != nil {
		return policySetChildren{}, err
	}
	exemptions, err := a.deps.Store.ListPolicyExemptionsByPolicySet(ctx, id)
	if err != nil {
		return policySetChildren{}, err
	}
	return policySetChildren{attachments: attachments, parameters: parameters, exemptions: exemptions}, nil
}

func (a *API) createPolicySet(w http.ResponseWriter, r *http.Request) {
	attrs, data, ok := decodePolicySetAttrs(w, r)
	if !ok {
		return
	}
	if attrs.Name == nil {
		dataIntegrityConflict(w)
		return
	}
	enforcement, ok := policyEnforcement(w, data, attrs.EnforcementLevel, "HARD_MANDATORY")
	if !ok {
		return
	}
	if !validShadowEnforcement(w, data, attrs.ShadowEnforcementLevel) {
		return
	}
	orgID, ok := requiredRelationshipID(w, data, "organization")
	if !ok {
		return
	}
	if !a.requirePolicyWrite(w, r, orgID) {
		return
	}
	links, ok := policyLinks(w, data, map[string]*string{"vcs": nil, "notificationConfiguration": nil})
	if !ok {
		return
	}
	if !a.validatePolicySetLinks(w, r, orgID, links) {
		return
	}
	email := callerEmail(r)
	set, err := a.deps.Store.CreatePolicySet(r.Context(), store.CreatePolicySetParams{
		ID:                          uuid.NewString(),
		Name:                        *attrs.Name,
		Description:                 attrs.Description,
		EnforcementLevel:            enforcement,
		ShadowEnforcementLevel:      attrs.ShadowEnforcementLevel,
		OverrideTeam:                attrs.OverrideTeam,
		Global:                      boolValue(attrs.Global, false),
		Repository:                  attrs.Repository,
		Branch:                      strPtr(stringOr(attrs.Branch, "main")),
		Folder:                      strPtr(stringOr(attrs.Folder, "/")),
		VcsID:                       links["vcs"],
		OrganizationID:              orgID,
		NotificationConfigurationID: links["notificationConfiguration"],
		CreatedBy:                   email,
		OpaVersion:                  attrs.OpaVersion,
	})
	if err != nil {
		if !a.writeStoreError(w, "create policy set", err) {
			a.internalError(w, "create policy set", err)
		}
		return
	}
	jsonapi.WriteResource(w, http.StatusCreated, policySetResource(set, policySetChildren{}), jsonapi.ParseFieldsets(r))
}

func (a *API) updatePolicySet(w http.ResponseWriter, r *http.Request) {
	set, ok := a.policySetFromScope(w, r)
	if !ok {
		return
	}
	attrs, data, ok := decodePolicySetAttrs(w, r)
	if !ok {
		return
	}
	if jsonapi.HasAttribute(data, "name") && attrs.Name == nil {
		dataIntegrityConflict(w)
		return
	}
	if jsonapi.HasAttribute(data, "global") && attrs.Global == nil {
		dataIntegrityConflict(w)
		return
	}
	enforcement, ok := policyEnforcement(w, data, attrs.EnforcementLevel, set.EnforcementLevel)
	if !ok {
		return
	}
	if !validShadowEnforcement(w, data, attrs.ShadowEnforcementLevel) {
		return
	}
	links, ok := policyLinks(w, data, map[string]*string{
		"vcs": set.VcsID, "notificationConfiguration": set.NotificationConfigurationID,
	})
	if !ok {
		return
	}
	if !a.validatePolicySetLinks(w, r, set.OrganizationID, links) {
		return
	}
	if _, err := a.deps.Store.UpdatePolicySet(r.Context(), store.UpdatePolicySetParams{
		ID:                          set.ID,
		Name:                        pickStringValue(data, "name", attrs.Name, set.Name),
		Description:                 pickString(data, "description", attrs.Description, set.Description),
		EnforcementLevel:            enforcement,
		ShadowEnforcementLevel:      pickString(data, "shadowEnforcementLevel", attrs.ShadowEnforcementLevel, set.ShadowEnforcementLevel),
		OverrideTeam:                pickString(data, "overrideTeam", attrs.OverrideTeam, set.OverrideTeam),
		Global:                      pickBool(data, "global", attrs.Global, set.Global),
		Repository:                  pickString(data, "repository", attrs.Repository, set.Repository),
		Branch:                      pickString(data, "branch", attrs.Branch, set.Branch),
		Folder:                      pickString(data, "folder", attrs.Folder, set.Folder),
		VcsID:                       links["vcs"],
		NotificationConfigurationID: links["notificationConfiguration"],
		OpaVersion:                  pickString(data, "opaVersion", attrs.OpaVersion, set.OpaVersion),
		UpdatedBy:                   callerEmail(r),
	}); err != nil {
		if !a.writeStoreError(w, "update policy set", err) {
			a.internalError(w, "update policy set", err)
		}
		return
	}
	w.WriteHeader(http.StatusNoContent)
}

func (a *API) deletePolicySet(w http.ResponseWriter, r *http.Request) {
	set, ok := a.policySetFromScope(w, r)
	if !ok {
		return
	}
	if err := a.deps.Store.DeletePolicySet(r.Context(), set.ID); err != nil {
		if !a.writeStoreError(w, "delete policy set", err) {
			a.internalError(w, "delete policy set", err)
		}
		return
	}
	w.WriteHeader(http.StatusNoContent)
}

func (a *API) listPolicyAttachments(w http.ResponseWriter, r *http.Request) {
	set, ok := a.policySetFromScope(w, r)
	if !ok {
		return
	}
	attachments, err := a.deps.Store.ListPolicyAttachments(r.Context(), set.ID)
	if err != nil {
		if !a.writeStoreError(w, "list policy attachments", err) {
			a.internalError(w, "list policy attachments", err)
		}
		return
	}
	resources := make([]jsonapi.Resource, 0, len(attachments))
	for _, attachment := range attachments {
		resources = append(resources, policyAttachmentResource(attachment))
	}
	jsonapi.WriteCollection(w, http.StatusOK, resources, len(resources), jsonapi.Page{}, jsonapi.ParseFieldsets(r))
}

func (a *API) createPolicyAttachment(w http.ResponseWriter, r *http.Request) {
	set, ok := a.policySetFromScope(w, r)
	if !ok {
		return
	}
	data, err := jsonapi.DecodeRequest(r)
	if err != nil {
		jsonapi.WriteError(w, http.StatusBadRequest, "Invalid Request", err.Error())
		return
	}
	if !requireType(w, data, "policy_attachment") {
		return
	}
	links, ok := policyLinks(w, data, map[string]*string{"workspace": nil, "project": nil, "tag": nil})
	if !ok {
		return
	}
	if !a.validatePolicyAttachmentLinks(w, r, set.OrganizationID, links) {
		return
	}
	email := callerEmail(r)
	attachment, err := a.deps.Store.CreatePolicyAttachment(r.Context(), store.CreatePolicyAttachmentParams{
		ID:          uuid.NewString(),
		PolicySetID: set.ID,
		WorkspaceID: links["workspace"],
		ProjectID:   links["project"],
		TagID:       links["tag"],
		CreatedBy:   email,
	})
	if err != nil {
		if !a.writeStoreError(w, "create policy attachment", err) {
			a.internalError(w, "create policy attachment", err)
		}
		return
	}
	jsonapi.WriteResource(w, http.StatusCreated, policyAttachmentResource(attachment), jsonapi.ParseFieldsets(r))
}

func (a *API) getPolicyAttachment(w http.ResponseWriter, r *http.Request) {
	attachment, ok := a.policyAttachmentInSet(w, r)
	if !ok {
		return
	}
	jsonapi.WriteResource(w, http.StatusOK, policyAttachmentResource(attachment), jsonapi.ParseFieldsets(r))
}

func (a *API) updatePolicyAttachment(w http.ResponseWriter, r *http.Request) {
	attachment, ok := a.policyAttachmentInSet(w, r)
	if !ok {
		return
	}
	set, ok := a.policySetFromScope(w, r)
	if !ok {
		return
	}
	data, err := jsonapi.DecodeRequest(r)
	if err != nil {
		jsonapi.WriteError(w, http.StatusBadRequest, "Invalid Request", err.Error())
		return
	}
	if !requireType(w, data, "policy_attachment") {
		return
	}
	links, ok := policyLinks(w, data, map[string]*string{
		"workspace": attachment.WorkspaceID, "project": attachment.ProjectID, "tag": attachment.TagID,
	})
	if !ok {
		return
	}
	if !a.validatePolicyAttachmentLinks(w, r, set.OrganizationID, links) {
		return
	}
	if _, err := a.deps.Store.UpdatePolicyAttachment(r.Context(), store.UpdatePolicyAttachmentParams{
		ID:          attachment.ID,
		WorkspaceID: links["workspace"],
		ProjectID:   links["project"],
		TagID:       links["tag"],
		UpdatedBy:   callerEmail(r),
	}); err != nil {
		if !a.writeStoreError(w, "update policy attachment", err) {
			a.internalError(w, "update policy attachment", err)
		}
		return
	}
	w.WriteHeader(http.StatusNoContent)
}

func (a *API) deletePolicyAttachment(w http.ResponseWriter, r *http.Request) {
	attachment, ok := a.policyAttachmentInSet(w, r)
	if !ok {
		return
	}
	if err := a.deps.Store.DeletePolicyAttachment(r.Context(), attachment.ID); err != nil {
		if !a.writeStoreError(w, "delete policy attachment", err) {
			a.internalError(w, "delete policy attachment", err)
		}
		return
	}
	w.WriteHeader(http.StatusNoContent)
}

func (a *API) listPolicySetParameters(w http.ResponseWriter, r *http.Request) {
	set, ok := a.policySetFromScope(w, r)
	if !ok {
		return
	}
	parameters, err := a.deps.Store.ListPolicySetParameters(r.Context(), set.ID)
	if err != nil {
		if !a.writeStoreError(w, "list policy set parameters", err) {
			a.internalError(w, "list policy set parameters", err)
		}
		return
	}
	resources := make([]jsonapi.Resource, 0, len(parameters))
	for _, parameter := range parameters {
		resources = append(resources, policySetParameterResource(parameter))
	}
	jsonapi.WriteCollection(w, http.StatusOK, resources, len(resources), jsonapi.Page{}, jsonapi.ParseFieldsets(r))
}

func (a *API) createPolicySetParameter(w http.ResponseWriter, r *http.Request) {
	set, ok := a.policySetFromScope(w, r)
	if !ok {
		return
	}
	attrs, _, ok := decodePolicySetParameterAttrs(w, r)
	if !ok {
		return
	}
	if attrs.Key == nil || attrs.Value == nil {
		dataIntegrityConflict(w)
		return
	}
	email := callerEmail(r)
	parameter, err := a.deps.Store.CreatePolicySetParameter(r.Context(), store.CreatePolicySetParameterParams{
		ID:             uuid.NewString(),
		PolicySetID:    set.ID,
		ParameterKey:   *attrs.Key,
		ParameterValue: *attrs.Value,
		Description:    attrs.Description,
		CreatedBy:      email,
	})
	if err != nil {
		if !a.writeStoreError(w, "create policy set parameter", err) {
			a.internalError(w, "create policy set parameter", err)
		}
		return
	}
	jsonapi.WriteResource(w, http.StatusCreated, policySetParameterResource(parameter), jsonapi.ParseFieldsets(r))
}

func (a *API) updatePolicySetParameter(w http.ResponseWriter, r *http.Request) {
	parameter, ok := a.policySetParameterInSet(w, r)
	if !ok {
		return
	}
	attrs, data, ok := decodePolicySetParameterAttrs(w, r)
	if !ok {
		return
	}
	if (jsonapi.HasAttribute(data, "key") && attrs.Key == nil) ||
		(jsonapi.HasAttribute(data, "value") && attrs.Value == nil) {
		dataIntegrityConflict(w)
		return
	}
	if _, err := a.deps.Store.UpdatePolicySetParameter(r.Context(), store.UpdatePolicySetParameterParams{
		ID:             parameter.ID,
		ParameterKey:   pickStringValue(data, "key", attrs.Key, parameter.ParameterKey),
		ParameterValue: pickStringValue(data, "value", attrs.Value, parameter.ParameterValue),
		Description:    pickString(data, "description", attrs.Description, parameter.Description),
		UpdatedBy:      callerEmail(r),
	}); err != nil {
		if !a.writeStoreError(w, "update policy set parameter", err) {
			a.internalError(w, "update policy set parameter", err)
		}
		return
	}
	w.WriteHeader(http.StatusNoContent)
}

func (a *API) deletePolicySetParameter(w http.ResponseWriter, r *http.Request) {
	parameter, ok := a.policySetParameterInSet(w, r)
	if !ok {
		return
	}
	if err := a.deps.Store.DeletePolicySetParameter(r.Context(), parameter.ID); err != nil {
		if !a.writeStoreError(w, "delete policy set parameter", err) {
			a.internalError(w, "delete policy set parameter", err)
		}
		return
	}
	w.WriteHeader(http.StatusNoContent)
}

func (a *API) listPolicyExemptions(w http.ResponseWriter, r *http.Request) {
	orgID := param(r, "organizationID")
	if _, err := a.deps.Store.GetOrganization(r.Context(), orgID); err != nil {
		if !a.writeStoreError(w, "get organization", err) {
			a.internalError(w, "get organization", err)
		}
		return
	}
	exemptions, err := a.deps.Store.ListPolicyExemptionsByOrganization(r.Context(), orgID)
	if err != nil {
		if !a.writeStoreError(w, "list policy exemptions", err) {
			a.internalError(w, "list policy exemptions", err)
		}
		return
	}
	a.writePolicyExemptions(w, r, exemptions)
}

func (a *API) writePolicyExemptions(w http.ResponseWriter, r *http.Request, exemptions []store.PolicyExemption) {
	want := parseInclude(r)
	resources := make([]jsonapi.Resource, 0, len(exemptions))
	included := []jsonapi.Resource{}
	for _, exemption := range exemptions {
		if want["policySet"] {
			set, err := a.deps.Store.GetPolicySet(r.Context(), exemption.PolicySetID)
			if err != nil {
				if !a.writeStoreError(w, "get policy set", err) {
					a.internalError(w, "get policy set", err)
				}
				return
			}
			included = append(included, policySetResource(set, policySetChildren{}))
		}
		if want["workspace"] && exemption.WorkspaceID != nil && *exemption.WorkspaceID != "" {
			ws, err := a.deps.Store.GetWorkspace(r.Context(), *exemption.WorkspaceID)
			if err != nil {
				if !a.writeStoreError(w, "get workspace", err) {
					a.internalError(w, "get workspace", err)
				}
				return
			}
			included = append(included, workspaceResource(ws))
		}
		if want["project"] && exemption.ProjectID != nil && *exemption.ProjectID != "" {
			project, err := a.deps.Store.GetProject(r.Context(), *exemption.ProjectID)
			if err != nil {
				if !a.writeStoreError(w, "get project", err) {
					a.internalError(w, "get project", err)
				}
				return
			}
			included = append(included, projectResource(project))
		}
		resources = append(resources, policyExemptionResource(exemption))
	}
	jsonapi.Write(w, http.StatusOK, jsonapi.Document{Data: resources, Included: included}, jsonapi.ParseFieldsets(r))
}

func (a *API) createPolicyExemption(w http.ResponseWriter, r *http.Request) {
	attrs, data, ok := decodePolicyExemptionAttrs(w, r)
	if !ok {
		return
	}
	orgID, ok := requiredRelationshipID(w, data, "organization")
	if !ok {
		return
	}
	if !a.requirePolicyWrite(w, r, orgID) {
		return
	}
	policySetID, ok := requiredRelationshipID(w, data, "policySet")
	if !ok {
		return
	}
	set, err := a.deps.Store.GetPolicySet(r.Context(), policySetID)
	if err != nil {
		if !a.writeStoreError(w, "get policy set", err) {
			a.internalError(w, "get policy set", err)
		}
		return
	}
	if set.OrganizationID != orgID {
		jsonapi.WriteError(w, http.StatusForbidden, "Forbidden", "Policy set does not belong to the organization")
		return
	}
	if attrs.RuleID == nil || attrs.TicketReference == nil || attrs.Justification == nil {
		dataIntegrityConflict(w)
		return
	}
	links, ok := policyLinks(w, data, map[string]*string{"workspace": nil, "project": nil})
	if !ok {
		return
	}
	if !a.validatePolicyAttachmentLinks(w, r, orgID, links) {
		return
	}
	email := callerEmail(r)
	exemption, err := a.deps.Store.CreatePolicyExemption(r.Context(), store.CreatePolicyExemptionParams{
		ID:              uuid.NewString(),
		PolicySetID:     set.ID,
		RuleID:          *attrs.RuleID,
		WorkspaceID:     links["workspace"],
		ProjectID:       links["project"],
		TicketReference: *attrs.TicketReference,
		Justification:   *attrs.Justification,
		ExpiresAt:       timestamptz(attrs.ExpiresAt),
		OrganizationID:  orgID,
		CreatedBy:       email,
	})
	if err != nil {
		if !a.writeStoreError(w, "create policy exemption", err) {
			a.internalError(w, "create policy exemption", err)
		}
		return
	}
	jsonapi.WriteResource(w, http.StatusCreated, policyExemptionResource(exemption), jsonapi.ParseFieldsets(r))
}

func (a *API) getPolicyExemption(w http.ResponseWriter, r *http.Request) {
	exemption, ok := a.policyExemptionFromScope(w, r)
	if !ok {
		return
	}
	jsonapi.WriteResource(w, http.StatusOK, policyExemptionResource(exemption), jsonapi.ParseFieldsets(r))
}

func (a *API) updatePolicyExemption(w http.ResponseWriter, r *http.Request) {
	exemption, ok := a.policyExemptionFromScope(w, r)
	if !ok {
		return
	}
	attrs, data, ok := decodePolicyExemptionAttrs(w, r)
	if !ok {
		return
	}
	for name, value := range map[string]*string{
		"ruleId": attrs.RuleID, "ticketReference": attrs.TicketReference, "justification": attrs.Justification,
	} {
		if jsonapi.HasAttribute(data, name) && value == nil {
			dataIntegrityConflict(w)
			return
		}
	}
	links, ok := policyLinks(w, data, map[string]*string{
		"workspace": exemption.WorkspaceID, "project": exemption.ProjectID,
	})
	if !ok {
		return
	}
	if !a.validatePolicyAttachmentLinks(w, r, exemption.OrganizationID, links) {
		return
	}
	if _, err := a.deps.Store.UpdatePolicyExemption(r.Context(), store.UpdatePolicyExemptionParams{
		ID:              exemption.ID,
		RuleID:          pickStringValue(data, "ruleId", attrs.RuleID, exemption.RuleID),
		WorkspaceID:     links["workspace"],
		ProjectID:       links["project"],
		TicketReference: pickStringValue(data, "ticketReference", attrs.TicketReference, exemption.TicketReference),
		Justification:   pickStringValue(data, "justification", attrs.Justification, exemption.Justification),
		ExpiresAt:       pickTimestamp(data, "expiresAt", attrs.ExpiresAt, exemption.ExpiresAt),
		UpdatedBy:       callerEmail(r),
	}); err != nil {
		if !a.writeStoreError(w, "update policy exemption", err) {
			a.internalError(w, "update policy exemption", err)
		}
		return
	}
	w.WriteHeader(http.StatusNoContent)
}

func (a *API) deletePolicyExemption(w http.ResponseWriter, r *http.Request) {
	exemption, ok := a.policyExemptionFromScope(w, r)
	if !ok {
		return
	}
	if err := a.deps.Store.DeletePolicyExemption(r.Context(), exemption.ID); err != nil {
		if !a.writeStoreError(w, "delete policy exemption", err) {
			a.internalError(w, "delete policy exemption", err)
		}
		return
	}
	w.WriteHeader(http.StatusNoContent)
}

// policySetScope resolves a root-level policy set and its organization access.
func (a *API) policySetScope(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		id := param(r, "policySetID")
		if _, err := uuid.Parse(id); err != nil {
			jsonapi.WriteError(w, http.StatusBadRequest, "Invalid Value", "Invalid value: "+id)
			return
		}
		set, err := a.deps.Store.GetPolicySet(r.Context(), id)
		if err != nil {
			if !a.writeStoreError(w, "get policy set", err) {
				a.internalError(w, "get policy set", err)
			}
			return
		}
		access, err := a.resolveOrganizationAuthorization(r, set.OrganizationID)
		if err != nil {
			a.internalError(w, "resolve policy set access", err)
			return
		}
		if !access.Member {
			writePermissionDenied(w, r.Method)
			return
		}
		next.ServeHTTP(w, r.WithContext(withAuthorization(r.Context(), access)))
	})
}

// policyExemptionScope is the policy set counterpart for root exemptions.
func (a *API) policyExemptionScope(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		id := param(r, "exemptionID")
		if _, err := uuid.Parse(id); err != nil {
			jsonapi.WriteError(w, http.StatusBadRequest, "Invalid Value", "Invalid value: "+id)
			return
		}
		exemption, err := a.deps.Store.GetPolicyExemption(r.Context(), id)
		if err != nil {
			if !a.writeStoreError(w, "get policy exemption", err) {
				a.internalError(w, "get policy exemption", err)
			}
			return
		}
		access, err := a.resolveOrganizationAuthorization(r, exemption.OrganizationID)
		if err != nil {
			a.internalError(w, "resolve policy exemption access", err)
			return
		}
		if !access.Member {
			writePermissionDenied(w, r.Method)
			return
		}
		next.ServeHTTP(w, r.WithContext(withAuthorization(r.Context(), access)))
	})
}

func (a *API) policySetFromScope(w http.ResponseWriter, r *http.Request) (store.PolicySet, bool) {
	set, err := a.deps.Store.GetPolicySet(r.Context(), param(r, "policySetID"))
	if err != nil {
		if !a.writeStoreError(w, "get policy set", err) {
			a.internalError(w, "get policy set", err)
		}
		return store.PolicySet{}, false
	}
	return set, true
}

func (a *API) policyExemptionFromScope(w http.ResponseWriter, r *http.Request) (store.PolicyExemption, bool) {
	exemption, err := a.deps.Store.GetPolicyExemption(r.Context(), param(r, "exemptionID"))
	if err != nil {
		if !a.writeStoreError(w, "get policy exemption", err) {
			a.internalError(w, "get policy exemption", err)
		}
		return store.PolicyExemption{}, false
	}
	return exemption, true
}

func (a *API) policyAttachmentInSet(w http.ResponseWriter, r *http.Request) (store.PolicyAttachment, bool) {
	set, ok := a.policySetFromScope(w, r)
	if !ok {
		return store.PolicyAttachment{}, false
	}
	attachment, err := a.deps.Store.GetPolicyAttachment(r.Context(), param(r, "attachmentID"))
	if err != nil {
		if !a.writeStoreError(w, "get policy attachment", err) {
			a.internalError(w, "get policy attachment", err)
		}
		return store.PolicyAttachment{}, false
	}
	if attachment.PolicySetID != set.ID {
		jsonapi.WriteError(w, http.StatusNotFound, "Not Found", "Unknown identifier ["+attachment.ID+"] for policy_attachment")
		return store.PolicyAttachment{}, false
	}
	return attachment, true
}

func (a *API) policySetParameterInSet(w http.ResponseWriter, r *http.Request) (store.PolicySetParameter, bool) {
	set, ok := a.policySetFromScope(w, r)
	if !ok {
		return store.PolicySetParameter{}, false
	}
	parameter, err := a.deps.Store.GetPolicySetParameter(r.Context(), param(r, "parameterID"))
	if err != nil {
		if !a.writeStoreError(w, "get policy set parameter", err) {
			a.internalError(w, "get policy set parameter", err)
		}
		return store.PolicySetParameter{}, false
	}
	if parameter.PolicySetID != set.ID {
		jsonapi.WriteError(w, http.StatusNotFound, "Not Found", "Unknown identifier ["+parameter.ID+"] for policy_set_parameter")
		return store.PolicySetParameter{}, false
	}
	return parameter, true
}

// requirePolicyWrite is the root-route equivalent of requireWrite for the
// policy collections, which carry their organization in the body.
func (a *API) requirePolicyWrite(w http.ResponseWriter, r *http.Request, orgID string) bool {
	access, err := a.resolveOrganizationAuthorization(r, orgID)
	if err != nil {
		a.internalError(w, "resolve policy access", err)
		return false
	}
	if !access.Member {
		writePermissionDenied(w, r.Method)
		return false
	}
	if !access.Owner && !access.Permissions.ManagePolicies {
		writePermissionDenied(w, r.Method)
		return false
	}
	return true
}

func (a *API) validatePolicySetLinks(w http.ResponseWriter, r *http.Request, orgID string, links map[string]*string) bool {
	if links["vcs"] != nil && !a.checkVCSRelationship(w, r, *links["vcs"], orgID) {
		return false
	}
	if links["notificationConfiguration"] != nil {
		config, err := a.deps.Store.GetNotificationConfiguration(r.Context(), *links["notificationConfiguration"])
		if err != nil {
			return a.relationshipLoadFailed(w, "get notification configuration", err,
				*links["notificationConfiguration"], "notification_configuration")
		}
		if config.OrganizationID != orgID {
			jsonapi.WriteError(w, http.StatusForbidden, "Forbidden", "Notification configuration does not belong to the policy set organization")
			return false
		}
	}
	return true
}

// validatePolicyAttachmentLinks keeps attachments and exemptions inside the
// policy set's organization.
func (a *API) validatePolicyAttachmentLinks(w http.ResponseWriter, r *http.Request, orgID string, links map[string]*string) bool {
	if links["workspace"] != nil {
		ws, err := a.deps.Store.GetWorkspace(r.Context(), *links["workspace"])
		if err != nil {
			return a.relationshipLoadFailed(w, "get workspace", err, *links["workspace"], "workspace")
		}
		if ws.OrganizationID != orgID {
			jsonapi.WriteError(w, http.StatusForbidden, "Forbidden", "Workspace does not belong to the policy set organization")
			return false
		}
	}
	if links["project"] != nil {
		project, err := a.deps.Store.GetProject(r.Context(), *links["project"])
		if err != nil {
			return a.relationshipLoadFailed(w, "get project", err, *links["project"], "project")
		}
		if project.OrganizationID != orgID {
			jsonapi.WriteError(w, http.StatusForbidden, "Forbidden", "Project does not belong to the policy set organization")
			return false
		}
	}
	if links["tag"] != nil {
		tag, err := a.deps.Store.GetTag(r.Context(), *links["tag"])
		if err != nil {
			return a.relationshipLoadFailed(w, "get tag", err, *links["tag"], "tag")
		}
		if tag.OrganizationID != orgID {
			jsonapi.WriteError(w, http.StatusForbidden, "Forbidden", "Tag does not belong to the policy set organization")
			return false
		}
	}
	return true
}

func policyLinks(w http.ResponseWriter, data *jsonapi.RequestData, current map[string]*string) (map[string]*string, bool) {
	links := map[string]*string{}
	for name, id := range current {
		links[name] = id
	}
	for name := range current {
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

func requiredRelationshipID(w http.ResponseWriter, data *jsonapi.RequestData, name string) (string, bool) {
	id, present, err := jsonapi.RelationshipID(data, name)
	if err != nil {
		jsonapi.WriteError(w, http.StatusBadRequest, "Invalid Relationship", err.Error())
		return "", false
	}
	if !present || id == "" {
		dataIntegrityConflict(w)
		return "", false
	}
	return id, true
}

// policyEnforcement validates an enforcement level, keeping fallback when the
// attribute is absent.
func policyEnforcement(w http.ResponseWriter, data *jsonapi.RequestData, value *string, fallback string) (string, bool) {
	if !jsonapi.HasAttribute(data, "enforcementLevel") {
		return fallback, true
	}
	if value == nil {
		dataIntegrityConflict(w)
		return "", false
	}
	if !policyEnforcementLevels[*value] {
		jsonapi.WriteError(w, http.StatusBadRequest, "Invalid Value", "Invalid value: "+*value)
		return "", false
	}
	return *value, true
}

func validShadowEnforcement(w http.ResponseWriter, data *jsonapi.RequestData, value *string) bool {
	if value == nil {
		return true
	}
	if !policyEnforcementLevels[*value] {
		jsonapi.WriteError(w, http.StatusBadRequest, "Invalid Value", "Invalid value: "+*value)
		return false
	}
	return true
}

func decodePolicySetAttrs(w http.ResponseWriter, r *http.Request) (policySetAttrs, *jsonapi.RequestData, bool) {
	data, err := jsonapi.DecodeRequest(r)
	if err != nil {
		jsonapi.WriteError(w, http.StatusBadRequest, "Invalid Request", err.Error())
		return policySetAttrs{}, nil, false
	}
	if !requireType(w, data, "policy_set") {
		return policySetAttrs{}, nil, false
	}
	attrs, err := jsonapi.DecodeAttributes[policySetAttrs](data)
	if err != nil {
		jsonapi.WriteError(w, http.StatusBadRequest, "Invalid Attributes", err.Error())
		return policySetAttrs{}, nil, false
	}
	return attrs, data, true
}

func decodePolicySetParameterAttrs(w http.ResponseWriter, r *http.Request) (policySetParameterAttrs, *jsonapi.RequestData, bool) {
	data, err := jsonapi.DecodeRequest(r)
	if err != nil {
		jsonapi.WriteError(w, http.StatusBadRequest, "Invalid Request", err.Error())
		return policySetParameterAttrs{}, nil, false
	}
	if !requireType(w, data, "policy_set_parameter") {
		return policySetParameterAttrs{}, nil, false
	}
	attrs, err := jsonapi.DecodeAttributes[policySetParameterAttrs](data)
	if err != nil {
		jsonapi.WriteError(w, http.StatusBadRequest, "Invalid Attributes", err.Error())
		return policySetParameterAttrs{}, nil, false
	}
	return attrs, data, true
}

func decodePolicyExemptionAttrs(w http.ResponseWriter, r *http.Request) (policyExemptionAttrs, *jsonapi.RequestData, bool) {
	data, err := jsonapi.DecodeRequest(r)
	if err != nil {
		jsonapi.WriteError(w, http.StatusBadRequest, "Invalid Request", err.Error())
		return policyExemptionAttrs{}, nil, false
	}
	if !requireType(w, data, "policy_exemption") {
		return policyExemptionAttrs{}, nil, false
	}
	attrs, err := jsonapi.DecodeAttributes[policyExemptionAttrs](data)
	if err != nil {
		jsonapi.WriteError(w, http.StatusBadRequest, "Invalid Attributes", err.Error())
		return policyExemptionAttrs{}, nil, false
	}
	return attrs, data, true
}

func timestamptz(value *time.Time) pgtype.Timestamptz {
	if value == nil {
		return pgtype.Timestamptz{}
	}
	return pgtype.Timestamptz{Time: *value, Valid: true}
}

// pickTimestamp keeps the existing value when the attribute is absent and
// clears it when the request sends null.
func pickTimestamp(data *jsonapi.RequestData, name string, incoming *time.Time, existing pgtype.Timestamptz) pgtype.Timestamptz {
	if !jsonapi.HasAttribute(data, name) {
		return existing
	}
	return timestamptz(incoming)
}
