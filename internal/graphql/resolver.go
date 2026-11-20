package graphql

// This file will not be regenerated automatically.
//
// It serves as dependency injection for your app, add any dependencies you require here.

type Resolver struct {
	db DB
}

func (r *Resolver) NotificationConfiguration() NotificationConfigurationResolver {
	return &notificationConfigurationResolver{r}
}

func (r *Resolver) Organization() OrganizationResolver { return &organizationResolver{r} }

func (r *Resolver) Project() ProjectResolver { return &projectResolver{r} }

func (r *Resolver) Query() QueryResolver { return &queryResolver{r} }

func (r *Resolver) Workspace() WorkspaceResolver { return &workspaceResolver{r} }

type notificationConfigurationResolver struct{ *Resolver }
type organizationResolver struct{ *Resolver }
type projectResolver struct{ *Resolver }
type queryResolver struct{ *Resolver }
type workspaceResolver struct{ *Resolver }
