package tfe

import (
	"context"
	"errors"
	"fmt"
	"net/url"
)

// Compile-time proof of interface implementation.
var _ RegistryModules = (*registryModules)(nil)

// RegistryModules describes all the registry module related methods that the Terraform
// Enterprise API supports.
//
// TFE API docs: https://www.terraform.io/docs/cloud/api/modules.html
type RegistryModules interface {
	// Create a registry module without VCS
	Create(ctx context.Context, options RegistryModuleCreateOptions) (*RegistryModule, error)

	// Create and publish a VCS backed registry module
	CreateWithVCSConnection(ctx context.Context, options RegistryModuleCreateFromVCSConnectionOptions) (*RegistryModule, error)

	// Delete a registry module
	Delete(ctx context.Context, organization string, name string) error

	// Delete a specific registry module provider
	DeleteProvider(ctx context.Context, organization string, name string, provider string) error

	// Delete a specific registry module version
	DeleteVersion(ctx context.Context, organization string, name string, provider string, version string) error
}

// registryModules implements RegistryModules.
type registryModules struct {
	client *Client
}

// RegistryModuleStatus represents the status of the registry module
type RegistryModuleStatus string

// List of available registry module statuses
const (
	RegistryModuleStatusPending       RegistryModuleStatus = "pending"
	RegistryModuleStatusNoVersionTags RegistryModuleStatus = "no_version_tags"
	RegistryModuleStatusSetupFailed   RegistryModuleStatus = "setup_failed"
	RegistryModuleStatusSetupComplete RegistryModuleStatus = "setup_complete"
)

// RegistryModuleVersionStatus represents the status of a specific version of a registry module
type RegistryModuleVersionStatus string

// List of available registry module version statuses
const (
	RegistryModuleVersionStatusPending             RegistryModuleVersionStatus = "pending"
	RegistryModuleVersionStatusCloning             RegistryModuleVersionStatus = "cloning"
	RegistryModuleVersionStatusCloneFailed         RegistryModuleVersionStatus = "clone_failed"
	RegistryModuleVersionStatusRegIngressReqFailed RegistryModuleVersionStatus = "reg_ingress_req_failed"
	RegistryModuleVersionStatusRegIngressing       RegistryModuleVersionStatus = "reg_ingressing"
	RegistryModuleVersionStatusRegIngressFailed    RegistryModuleVersionStatus = "reg_ingress_failed"
	RegistryModuleVersionStatusOk                  RegistryModuleVersionStatus = "ok"
)

// RegistryModule represents a registry module
type RegistryModule struct {
	ID              string                          `jsonapi:"primary,registry-modules"`
	Name            string                          `jsonapi:"attr,name"`
	Provider        string                          `jsonapi:"attr,provider"`
	Status          RegistryModuleStatus            `jsonapi:"attr,status"`
	VersionStatuses []RegistryModuleVersionStatuses `jsonapi:"attr,version-statuses"`
	Permissions     *RegistryModulePermissions      `jsonapi:"attr,permissions"`
	CreatedAt       string                          `jsonapi:"attr,created-at"`
	UpdatedAt       string                          `jsonapi:"attr,updated-at"`

	// Relations
	Organization *Organization `jsonapi:"relation,organization"`
}

type RegistryModulePermissions struct {
	CanDelete bool `json:"can-delete"`
	CanResync bool `json:"can-resync"`
	CanRetry  bool `json:"can-retry"`
}

type RegistryModuleVersionStatuses struct {
	Version string                      `json:"version"`
	Status  RegistryModuleVersionStatus `json:"status"`
	Error   string                      `json:"error"`
}

// RegistryModuleCreateOptions is used when creating a registry module
type RegistryModuleCreateOptions struct {
	VCSRepo *VCSRepoOptions `jsonapi:"attr,vcs-repo,omitempty"`
}

// Create a new registry module to the TFE private registry
func (r *registryModules) Create(ctx context.Context, options RegistryModuleCreateOptions) (*RegistryModule, error) {
	req, err := r.client.newRequest("POST", "registry-modules", &options)
	if err != nil {
		return nil, err
	}

	m := &RegistryModule{}
	err = r.client.do(ctx, req, m)
	if err != nil {
		return nil, err
	}

	return m, nil
}

// RegistryModuleCreateOptions is used when creating a registry module
type RegistryModuleCreateFromVCSConnectionOptions struct {
	// VCS repository information
	VCSRepo *VCSRepoOptions `jsonapi:"attr,vcs-repo,omitempty"`
}

// CreateWithVCSConnection is used to create abd publish a new registry module from a VCS repo to the TFE private registry
func (r *registryModules) CreateWithVCSConnection(ctx context.Context, options RegistryModuleCreateFromVCSConnectionOptions) (*RegistryModule, error) {
	req, err := r.client.newRequest("POST", "registry-modules", &options)
	if err != nil {
		return nil, err
	}

	m := &RegistryModule{}
	err = r.client.do(ctx, req, m)
	if err != nil {
		return nil, err
	}

	return m, nil
}

// Delete is used to delete the entire module on the TFE private registry
func (r *registryModules) Delete(ctx context.Context, organization string, name string) error {
	if !validStringID(&organization) {
		return errors.New("invalid value for organization")
	}
	if !validString(&name) {
		return errors.New("name is required")
	}
	if !validStringID(&name) {
		return errors.New("invalid value for name")
	}

	u := fmt.Sprintf(
		"registry-modules/actions/delete/%s/%s",
		url.QueryEscape(organization),
		url.QueryEscape(name),
	)
	req, err := r.client.newRequest("POST", u, nil)
	if err != nil {
		return err
	}

	return r.client.do(ctx, req, nil)
}

// DeleteProvider is used to delete the specific module provider on the TFE private registry
func (r *registryModules) DeleteProvider(ctx context.Context, organization string, name string, provider string) error {
	if !validStringID(&organization) {
		return errors.New("invalid value for organization")
	}
	if !validString(&name) {
		return errors.New("name is required")
	}
	if !validStringID(&name) {
		return errors.New("invalid value for name")
	}
	if !validString(&provider) {
		return errors.New("provider is required")
	}
	if !validStringID(&provider) {
		return errors.New("invalid value for provider")
	}

	u := fmt.Sprintf(
		"registry-modules/actions/delete/%s/%s/%s",
		url.QueryEscape(organization),
		url.QueryEscape(name),
		url.QueryEscape(provider),
	)
	req, err := r.client.newRequest("POST", u, nil)
	if err != nil {
		return err
	}

	return r.client.do(ctx, req, nil)
}

// DeleteVersion is used to delete the specific module version on the TFE private registry
func (r *registryModules) DeleteVersion(ctx context.Context, organization string, name string, provider string, version string) error {
	if !validStringID(&organization) {
		return errors.New("invalid value for organization")
	}
	if !validString(&name) {
		return errors.New("name is required")
	}
	if !validStringID(&name) {
		return errors.New("invalid value for name")
	}
	if !validString(&provider) {
		return errors.New("provider is required")
	}
	if !validStringID(&provider) {
		return errors.New("invalid value for provider")
	}
	if !validString(&version) {
		return errors.New("version is required")
	}
	if !validStringID(&version) {
		return errors.New("invalid value for version")
	}

	u := fmt.Sprintf(
		"registry-modules/actions/delete/%s/%s/%s/%s",
		url.QueryEscape(organization),
		url.QueryEscape(name),
		url.QueryEscape(provider),
		url.QueryEscape(version),
	)
	req, err := r.client.newRequest("POST", u, nil)
	if err != nil {
		return err
	}

	return r.client.do(ctx, req, nil)
}
