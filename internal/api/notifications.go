package api

import (
	"encoding/json"
	"log/slog"
	"net/http"

	"github.com/go-chi/chi/v5"
	"github.com/google/uuid"

	"github.com/aminbenmansour/terrakube-go/internal/httpapi"
	"github.com/aminbenmansour/terrakube-go/internal/jsonapi"
	"github.com/aminbenmansour/terrakube-go/internal/notification"
	"github.com/aminbenmansour/terrakube-go/internal/store"
)

// Upstream's NotificationChannelType and NotificationMessageStyle enums.
var (
	notificationChannels = map[string]bool{"SLACK": true, "TEAMS": true, "WEBHOOK": true}
	notificationStyles   = map[string]bool{"DETAILED": true, "SIMPLE": true}
	// NotificationTrigger.jobStatus is upstream's JobStatus enum, including the
	// values the run pipeline never writes.
	notificationTriggerStatuses = map[string]bool{
		"pending": true, "waitingApproval": true, "approved": true, "queue": true,
		"running": true, "completed": true, "noChanges": true, "notExecuted": true,
		"rejected": true, "cancelled": true, "failed": true, "unknown": true,
		"NeverExecuted": true,
	}
)

type notificationConfigurationAttrs struct {
	Name           *string `json:"name"`
	Description    *string `json:"description"`
	ChannelType    *string `json:"channelType"`
	DestinationUrl *string `json:"destinationUrl"`
	SigningSecret  *string `json:"signingSecret"`
	Active         *bool   `json:"active"`
	MessageStyle   *string `json:"messageStyle"`
}

type notificationTriggerAttrs struct {
	JobStatus *string `json:"jobStatus"`
}

func notificationConfigurationResource(config store.NotificationConfiguration) jsonapi.Resource {
	workspace := jsonapi.NullLink()
	if config.WorkspaceID != nil {
		workspace = jsonapi.Link("workspace", *config.WorkspaceID)
	}
	return jsonapi.Resource{
		Type: "notification_configuration",
		ID:   config.ID,
		Attributes: map[string]any{
			"active":         config.Active,
			"channelType":    config.ChannelType,
			"createdBy":      config.CreatedBy,
			"createdDate":    timestamp(config.CreatedDate),
			"description":    config.Description,
			"destinationUrl": config.DestinationUrl,
			"messageStyle":   config.MessageStyle,
			"name":           config.Name,
			"updatedBy":      config.UpdatedBy,
			"updatedDate":    timestamp(config.UpdatedDate),
		},
		Relationships: map[string]jsonapi.Relationship{
			"organization": jsonapi.Link("organization", config.OrganizationID),
			"workspace":    workspace,
		},
	}
}

func notificationTriggerResource(trigger store.NotificationTrigger) jsonapi.Resource {
	return jsonapi.Resource{
		Type: "notification_trigger",
		ID:   trigger.ID,
		Attributes: map[string]any{
			"createdBy":   trigger.CreatedBy,
			"createdDate": timestamp(trigger.CreatedDate),
			"jobStatus":   trigger.JobStatus,
			"updatedBy":   trigger.UpdatedBy,
			"updatedDate": timestamp(trigger.UpdatedDate),
		},
		Relationships: map[string]jsonapi.Relationship{
			"configuration": jsonapi.Link("notification_configuration", trigger.ConfigurationID),
		},
	}
}

func (a *API) listNotificationConfigurations(w http.ResponseWriter, r *http.Request) {
	orgID := param(r, "organizationID")
	if _, err := a.deps.Store.GetOrganization(r.Context(), orgID); err != nil {
		if !a.writeStoreError(w, "get organization", err) {
			a.internalError(w, "get organization", err)
		}
		return
	}
	configs, err := a.deps.Store.ListNotificationConfigurationsByOrganization(r.Context(), orgID)
	if err != nil {
		if !a.writeStoreError(w, "list notification configurations", err) {
			a.internalError(w, "list notification configurations", err)
		}
		return
	}
	a.writeNotificationConfigurations(w, r, configs)
}

func (a *API) listWorkspaceNotificationConfigurations(w http.ResponseWriter, r *http.Request) {
	workspace, ok := a.notificationWorkspace(w, r)
	if !ok {
		return
	}
	configs, err := a.deps.Store.ListNotificationConfigurationsByWorkspace(r.Context(), &workspace.ID)
	if err != nil {
		if !a.writeStoreError(w, "list notification configurations", err) {
			a.internalError(w, "list notification configurations", err)
		}
		return
	}
	a.writeNotificationConfigurations(w, r, configs)
}

// notificationWorkspace resolves the workspace in the nested organization route
// and refuses one that belongs somewhere else.
func (a *API) notificationWorkspace(w http.ResponseWriter, r *http.Request) (store.Workspace, bool) {
	workspace, err := a.deps.Store.GetWorkspace(r.Context(), param(r, "workspaceID"))
	if err != nil {
		if !a.writeStoreError(w, "get workspace", err) {
			a.internalError(w, "get workspace", err)
		}
		return store.Workspace{}, false
	}
	if workspace.OrganizationID != param(r, "organizationID") {
		jsonapi.WriteError(w, http.StatusNotFound, "Not Found", "The requested resource does not exist.")
		return store.Workspace{}, false
	}
	return workspace, true
}

func (a *API) writeNotificationConfigurations(w http.ResponseWriter, r *http.Request, configs []store.NotificationConfiguration) {
	resources := make([]jsonapi.Resource, 0, len(configs))
	for _, config := range configs {
		resources = append(resources, notificationConfigurationResource(config))
	}
	jsonapi.WriteCollection(w, http.StatusOK, resources, len(resources), jsonapi.Page{}, jsonapi.ParseFieldsets(r))
}

func (a *API) getNotificationConfiguration(w http.ResponseWriter, r *http.Request) {
	config, ok := a.notificationConfigurationByID(w, r)
	if !ok {
		return
	}
	jsonapi.WriteResource(w, http.StatusOK, notificationConfigurationResource(config), jsonapi.ParseFieldsets(r))
}

func (a *API) createNotificationConfiguration(w http.ResponseWriter, r *http.Request) {
	orgID := param(r, "organizationID")
	if _, err := a.deps.Store.GetOrganization(r.Context(), orgID); err != nil {
		if !a.writeStoreError(w, "get organization", err) {
			a.internalError(w, "get organization", err)
		}
		return
	}
	a.createNotificationConfigurationIn(w, r, orgID, nil)
}

func (a *API) createWorkspaceNotificationConfiguration(w http.ResponseWriter, r *http.Request) {
	workspace, ok := a.notificationWorkspace(w, r)
	if !ok {
		return
	}
	a.createNotificationConfigurationIn(w, r, workspace.OrganizationID, &workspace.ID)
}

func (a *API) createNotificationConfigurationIn(w http.ResponseWriter, r *http.Request, orgID string, workspaceID *string) {
	attrs, _, ok := decodeNotificationConfigurationAttrs(w, r)
	if !ok {
		return
	}
	// Upstream's columns are NOT NULL, so Elide answers its data-integrity 409
	// rather than a validation error.
	if attrs.Name == nil || attrs.ChannelType == nil || attrs.DestinationUrl == nil {
		writeInUseConflict(w)
		return
	}
	if !notificationChannels[*attrs.ChannelType] {
		jsonapi.WriteError(w, http.StatusBadRequest, "Invalid Value", "Invalid value: "+*attrs.ChannelType)
		return
	}
	style := "DETAILED"
	if attrs.MessageStyle != nil {
		if !notificationStyles[*attrs.MessageStyle] {
			jsonapi.WriteError(w, http.StatusBadRequest, "Invalid Value", "Invalid value: "+*attrs.MessageStyle)
			return
		}
		style = *attrs.MessageStyle
	}
	config, err := a.deps.Store.CreateNotificationConfiguration(r.Context(), store.CreateNotificationConfigurationParams{
		ID:             uuid.NewString(),
		Name:           *attrs.Name,
		Description:    attrs.Description,
		OrganizationID: orgID,
		WorkspaceID:    workspaceID,
		ChannelType:    *attrs.ChannelType,
		DestinationUrl: *attrs.DestinationUrl,
		SigningSecret:  attrs.SigningSecret,
		Active:         boolValue(attrs.Active, true),
		MessageStyle:   style,
	})
	if err != nil {
		if !a.writeStoreError(w, "create notification configuration", err) {
			a.internalError(w, "create notification configuration", err)
		}
		return
	}
	jsonapi.WriteResource(w, http.StatusCreated, notificationConfigurationResource(config), jsonapi.ParseFieldsets(r))
}

func (a *API) updateNotificationConfiguration(w http.ResponseWriter, r *http.Request) {
	config, ok := a.notificationConfigurationByID(w, r)
	if !ok {
		return
	}
	attrs, data, ok := decodeNotificationConfigurationAttrs(w, r)
	if !ok {
		return
	}
	name := pickString(data, "name", attrs.Name, &config.Name)
	channel := pickString(data, "channelType", attrs.ChannelType, &config.ChannelType)
	destination := pickString(data, "destinationUrl", attrs.DestinationUrl, &config.DestinationUrl)
	if name == nil || channel == nil || destination == nil {
		writeInUseConflict(w)
		return
	}
	if !notificationChannels[*channel] {
		jsonapi.WriteError(w, http.StatusBadRequest, "Invalid Value", "Invalid value: "+*channel)
		return
	}
	style := pickString(data, "messageStyle", attrs.MessageStyle, &config.MessageStyle)
	if style != nil && !notificationStyles[*style] {
		jsonapi.WriteError(w, http.StatusBadRequest, "Invalid Value", "Invalid value: "+*style)
		return
	}
	if style == nil {
		style = &config.MessageStyle
	}
	if _, err := a.deps.Store.UpdateNotificationConfiguration(r.Context(), store.UpdateNotificationConfigurationParams{
		ID:             config.ID,
		Name:           *name,
		Description:    pickString(data, "description", attrs.Description, config.Description),
		ChannelType:    *channel,
		DestinationUrl: *destination,
		SigningSecret:  pickString(data, "signingSecret", attrs.SigningSecret, config.SigningSecret),
		Active:         pickBool(data, "active", attrs.Active, config.Active),
		MessageStyle:   *style,
	}); err != nil {
		if !a.writeStoreError(w, "update notification configuration", err) {
			a.internalError(w, "update notification configuration", err)
		}
		return
	}
	w.WriteHeader(http.StatusNoContent)
}

func (a *API) deleteNotificationConfiguration(w http.ResponseWriter, r *http.Request) {
	config, ok := a.notificationConfigurationByID(w, r)
	if !ok {
		return
	}
	if err := a.deps.Store.DeleteNotificationConfiguration(r.Context(), config.ID); err != nil {
		if !a.writeStoreError(w, "delete notification configuration", err) {
			a.internalError(w, "delete notification configuration", err)
		}
		return
	}
	w.WriteHeader(http.StatusNoContent)
}

func (a *API) createNotificationTrigger(w http.ResponseWriter, r *http.Request) {
	config, ok := a.notificationConfigurationByID(w, r)
	if !ok {
		return
	}
	data, err := jsonapi.DecodeRequest(r)
	if err != nil {
		jsonapi.WriteError(w, http.StatusBadRequest, "Invalid Request", err.Error())
		return
	}
	if !requireType(w, data, "notification_trigger") {
		return
	}
	attrs, err := jsonapi.DecodeAttributes[notificationTriggerAttrs](data)
	if err != nil {
		jsonapi.WriteError(w, http.StatusBadRequest, "Invalid Attributes", err.Error())
		return
	}
	if attrs.JobStatus == nil {
		writeInUseConflict(w)
		return
	}
	if !notificationTriggerStatuses[*attrs.JobStatus] {
		jsonapi.WriteError(w, http.StatusBadRequest, "Invalid Value", "Invalid value: "+*attrs.JobStatus)
		return
	}
	trigger, err := a.deps.Store.CreateNotificationTrigger(r.Context(), store.CreateNotificationTriggerParams{
		ID:              uuid.NewString(),
		JobStatus:       *attrs.JobStatus,
		ConfigurationID: config.ID,
	})
	if err != nil {
		if !a.writeStoreError(w, "create notification trigger", err) {
			a.internalError(w, "create notification trigger", err)
		}
		return
	}
	jsonapi.WriteResource(w, http.StatusCreated, notificationTriggerResource(trigger), jsonapi.ParseFieldsets(r))
}

func (a *API) deleteNotificationTrigger(w http.ResponseWriter, r *http.Request) {
	config, ok := a.notificationConfigurationByID(w, r)
	if !ok {
		return
	}
	triggerID := param(r, "triggerID")
	if _, err := uuid.Parse(triggerID); err != nil {
		jsonapi.WriteError(w, http.StatusBadRequest, "Invalid Value", "Invalid value: "+triggerID)
		return
	}
	trigger, err := a.deps.Store.GetNotificationTrigger(r.Context(), triggerID)
	if err != nil || trigger.ConfigurationID != config.ID {
		jsonapi.WriteError(w, http.StatusNotFound, "Not Found", "The requested resource does not exist.")
		return
	}
	if err := a.deps.Store.DeleteNotificationTrigger(r.Context(), trigger.ID); err != nil {
		if !a.writeStoreError(w, "delete notification trigger", err) {
			a.internalError(w, "delete notification trigger", err)
		}
		return
	}
	w.WriteHeader(http.StatusNoContent)
}

// replaceNotificationTemplates rewrites the "applies to these templates" set,
// which is how the UI saves the picker. An empty body or empty list means every
// template, matching upstream (there is no independent add/remove).
func (a *API) replaceNotificationTemplates(w http.ResponseWriter, r *http.Request) {
	config, ok := a.notificationConfigurationByID(w, r)
	if !ok {
		return
	}
	ids, err := jsonapi.RelationshipList(r)
	if err != nil {
		jsonapi.WriteError(w, http.StatusBadRequest, "Invalid Request", err.Error())
		return
	}
	for _, id := range ids {
		if id.Type != "template" {
			writeInUseConflict(w)
			return
		}
		if _, err := a.deps.Store.GetTemplate(r.Context(), id.ID); err != nil {
			jsonapi.WriteError(w, http.StatusNotFound, "Not Found", "Unknown identifier "+id.ID+" for template")
			return
		}
	}
	if err := a.deps.Store.DeleteNotificationConfigurationTemplates(r.Context(), config.ID); err != nil {
		if !a.writeStoreError(w, "clear notification templates", err) {
			a.internalError(w, "clear notification templates", err)
		}
		return
	}
	for _, id := range ids {
		if err := a.deps.Store.AddNotificationConfigurationTemplate(r.Context(), store.AddNotificationConfigurationTemplateParams{
			ConfigurationID: config.ID,
			TemplateID:      id.ID,
		}); err != nil {
			if !a.writeStoreError(w, "add notification template", err) {
				a.internalError(w, "add notification template", err)
			}
			return
		}
	}
	w.WriteHeader(http.StatusNoContent)
}

func (a *API) notificationConfigurationByID(w http.ResponseWriter, r *http.Request) (store.NotificationConfiguration, bool) {
	id := param(r, "configurationID")
	if _, err := uuid.Parse(id); err != nil {
		jsonapi.WriteError(w, http.StatusBadRequest, "Invalid Value", "Invalid value: "+id)
		return store.NotificationConfiguration{}, false
	}
	config, err := a.deps.Store.GetNotificationConfiguration(r.Context(), id)
	if err != nil {
		if !a.writeStoreError(w, "get notification configuration", err) {
			a.internalError(w, "get notification configuration", err)
		}
		return store.NotificationConfiguration{}, false
	}
	return config, true
}

func decodeNotificationConfigurationAttrs(w http.ResponseWriter, r *http.Request) (notificationConfigurationAttrs, *jsonapi.RequestData, bool) {
	data, err := jsonapi.DecodeRequest(r)
	if err != nil {
		jsonapi.WriteError(w, http.StatusBadRequest, "Invalid Request", err.Error())
		return notificationConfigurationAttrs{}, nil, false
	}
	if !requireType(w, data, "notification_configuration") {
		return notificationConfigurationAttrs{}, nil, false
	}
	attrs, err := jsonapi.DecodeAttributes[notificationConfigurationAttrs](data)
	if err != nil {
		jsonapi.WriteError(w, http.StatusBadRequest, "Invalid Attributes", err.Error())
		return notificationConfigurationAttrs{}, nil, false
	}
	return attrs, data, true
}

// NewNotification builds the root /notification/v1 router. Upstream mounts the
// test-send endpoints outside /api/v1, and the UI calls them at the API origin.
func NewNotification(deps Deps) http.Handler {
	if deps.Logger == nil {
		deps.Logger = slog.Default()
	}
	if deps.Deliverer == nil {
		deps.Deliverer = notification.NewDeliverer(true)
	}
	a := &API{deps: deps}
	r := chi.NewRouter()
	r.Use(httpapi.CORS(deps.CORSOrigins))
	r.Use(httpapi.LimitBody(jsonapi.MaxBodyBytes))
	r.Use(a.authenticate)
	r.Route("/configuration/{configurationID}", func(r chi.Router) {
		r.Use(a.notificationScope)
		r.With(a.requireWrite(canManageWorkspace)).Post("/test", a.sendNotificationTest)
	})
	r.Route("/organization/{organizationID}", func(r chi.Router) {
		r.Use(a.notificationTestOrganization)
		r.With(a.requireWrite(canManageWorkspace)).Post("/configuration/test", a.sendAdHocNotificationTest)
	})
	return r
}

// notificationTestOrganization answers an unknown organization with 404 before
// the access middleware can turn it into a non-member 403.
func (a *API) notificationTestOrganization(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		id := param(r, "organizationID")
		if _, err := uuid.Parse(id); err != nil {
			jsonapi.WriteError(w, http.StatusBadRequest, "Invalid Value", "Invalid value: "+id)
			return
		}
		if _, err := a.deps.Store.GetOrganization(r.Context(), id); err != nil {
			if !a.writeStoreError(w, "get organization", err) {
				a.internalError(w, "get organization", err)
			}
			return
		}
		a.organizationAccess(next).ServeHTTP(w, r)
	})
}

type notificationAdHocTestAttrs struct {
	ChannelType    string  `json:"channelType"`
	DestinationUrl string  `json:"destinationUrl"`
	SigningSecret  *string `json:"signingSecret"`
}

// sendNotificationTest renders and sends the same test payload as a real
// delivery, without persisting a configuration or a delivery row.
func (a *API) sendNotificationTest(w http.ResponseWriter, r *http.Request) {
	config, ok := a.notificationConfigurationByID(w, r)
	if !ok {
		return
	}
	organization, err := a.deps.Store.GetOrganization(r.Context(), config.OrganizationID)
	if err != nil {
		if !a.writeStoreError(w, "get organization", err) {
			a.internalError(w, "get organization", err)
		}
		return
	}
	workspaceName := "(test notification)"
	if config.WorkspaceID != nil && *config.WorkspaceID != "" {
		workspace, err := a.deps.Store.GetWorkspace(r.Context(), *config.WorkspaceID)
		if err != nil {
			if !a.writeStoreError(w, "get workspace", err) {
				a.internalError(w, "get workspace", err)
			}
			return
		}
		workspaceName = workspace.Name
	}
	configurationName := config.Name
	if configurationName == "" {
		configurationName = "Test Notification"
	}
	a.deliverNotificationTest(w, r, config.ChannelType, config.DestinationUrl, config.SigningSecret, notification.Context{
		OrganizationName:  organization.Name,
		WorkspaceName:     workspaceName,
		Status:            "completed",
		ConfigurationName: configurationName,
		MessageStyle:      config.MessageStyle,
	})
}

// sendAdHocNotificationTest lets the create form verify a destination before
// the configuration has ever been saved.
func (a *API) sendAdHocNotificationTest(w http.ResponseWriter, r *http.Request) {
	organization, err := a.deps.Store.GetOrganization(r.Context(), param(r, "organizationID"))
	if err != nil {
		if !a.writeStoreError(w, "get organization", err) {
			a.internalError(w, "get organization", err)
		}
		return
	}
	attrs := notificationAdHocTestAttrs{}
	if err := json.NewDecoder(r.Body).Decode(&attrs); err != nil {
		jsonapi.WriteError(w, http.StatusBadRequest, "Bad Request", "The request body is not valid JSON.")
		return
	}
	if !notificationChannels[attrs.ChannelType] {
		jsonapi.WriteError(w, http.StatusBadRequest, "Invalid Value", "Invalid value: "+attrs.ChannelType)
		return
	}
	a.deliverNotificationTest(w, r, attrs.ChannelType, attrs.DestinationUrl, attrs.SigningSecret, notification.Context{
		OrganizationName:  organization.Name,
		WorkspaceName:     "(test notification)",
		Status:            "completed",
		ConfigurationName: "Test Notification",
		MessageStyle:      notification.StyleDetailed,
	})
}

func (a *API) deliverNotificationTest(w http.ResponseWriter, r *http.Request, channelType, destinationURL string, signingSecret *string, context notification.Context) {
	payload, err := notification.BuildPayload(channelType, context)
	if err != nil {
		a.internalError(w, "render notification payload", err)
		return
	}
	err = a.deps.Deliverer.Deliver(r.Context(), notification.Configuration{
		ChannelType:    channelType,
		DestinationURL: destinationURL,
		SigningSecret:  signingSecret,
	}, payload)
	if err != nil {
		a.deps.Logger.Warn("test notification failed", "channel", channelType, "error", err)
		w.WriteHeader(http.StatusBadGateway)
		return
	}
	w.WriteHeader(http.StatusOK)
}
