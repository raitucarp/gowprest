package gowprest

import (
	"encoding/json"
	"strings"
)

// Plugin status constants.
const (
	PluginStatusActive   = "active"
	PluginStatusInactive = "inactive"
)

// Plugin represents an installed WordPress plugin.
type Plugin struct {
	Plugin       string         `json:"plugin,omitempty"`
	Status       string         `json:"status,omitempty"`
	Name         string         `json:"name,omitempty"`
	PluginURI    string         `json:"plugin_uri,omitempty"`
	Author       string         `json:"author,omitempty"`
	AuthorURI    string         `json:"author_uri,omitempty"`
	Description  *Object        `json:"description,omitempty"`
	Version      string         `json:"version,omitempty"`
	NetworkOnly  bool           `json:"network_only,omitempty"`
	RequiresWP   string         `json:"requires_wp,omitempty"`
	RequiresPHP  string         `json:"requires_php,omitempty"`
	Textdomain   string         `json:"textdomain,omitempty"`
	Links        map[string]any `json:"_links,omitempty"`
}

// IsActive reports whether the plugin is currently active.
func (p *Plugin) IsActive() bool {
	return p.Status == PluginStatusActive
}

// IsInactive reports whether the plugin is currently inactive.
func (p *Plugin) IsInactive() bool {
	return p.Status == PluginStatusInactive
}

// Plugins handles requests to the WordPress Plugins API (/wp/v2/plugins).
type Plugins struct {
	client *RestClient
}

// Plugins returns a Plugins service instance.
func (c *RestClient) Plugins() *Plugins {
	return &Plugins{client: c}
}

// List returns a ListPlugins builder to query installed plugins.
func (api *Plugins) List() *ListPlugins {
	return &ListPlugins{
		endpoint:  "/wp/v2/plugins",
		client:    api.client,
		arguments: make(map[string]string),
	}
}

// Retrieve returns a RetrievePlugin builder to get a specific plugin.
// Accepts either a combined identifier ("akismet/akismet") or separate directory and file ("akismet", "akismet").
func (api *Plugins) Retrieve(plugin string, file ...string) *RetrievePlugin {
	path := "/wp/v2/plugins"
	if len(file) > 0 {
		path += "/" + strings.Trim(plugin, "/") + "/" + strings.Trim(file[0], "/")
	} else {
		path += "/" + strings.Trim(plugin, "/")
	}
	return &RetrievePlugin{
		endpoint:  path,
		client:    api.client,
		arguments: make(map[string]string),
	}
}

// Create returns a CreatePlugin builder to install a plugin by slug from WordPress.org.
func (api *Plugins) Create() *CreatePlugin {
	return &CreatePlugin{
		endpoint: "/wp/v2/plugins",
		client:   api.client,
		body:     make(map[string]any),
	}
}

// Update returns an UpdatePlugin builder to activate or deactivate a plugin.
func (api *Plugins) Update(plugin string, file ...string) *UpdatePlugin {
	path := "/wp/v2/plugins"
	if len(file) > 0 {
		path += "/" + strings.Trim(plugin, "/") + "/" + strings.Trim(file[0], "/")
	} else {
		path += "/" + strings.Trim(plugin, "/")
	}
	return &UpdatePlugin{
		endpoint:  path,
		client:    api.client,
		arguments: make(map[string]string),
		body:      make(map[string]any),
	}
}

// Delete returns a DeletePlugin builder to delete an installed plugin.
func (api *Plugins) Delete(plugin string, file ...string) *DeletePlugin {
	path := "/wp/v2/plugins"
	if len(file) > 0 {
		path += "/" + strings.Trim(plugin, "/") + "/" + strings.Trim(file[0], "/")
	} else {
		path += "/" + strings.Trim(plugin, "/")
	}
	return &DeletePlugin{
		endpoint:  path,
		client:    api.client,
		arguments: make(map[string]string),
	}
}

// ListPlugins handles querying plugin collections.
type ListPlugins struct {
	endpoint  string
	client    *RestClient
	arguments map[string]string
}

func (api *ListPlugins) Search(query string) *ListPlugins {
	api.arguments["search"] = query
	return api
}

func (api *ListPlugins) Status(statuses ...string) *ListPlugins {
	api.arguments["status"] = strings.Join(statuses, ",")
	return api
}

func (api *ListPlugins) StatusActive() *ListPlugins {
	return api.Status(PluginStatusActive)
}

func (api *ListPlugins) StatusInactive() *ListPlugins {
	return api.Status(PluginStatusInactive)
}

func (api *ListPlugins) Context(ctx string) *ListPlugins {
	api.arguments["context"] = ctx
	return api
}

func (api *ListPlugins) ContextView() *ListPlugins {
	return api.Context("view")
}

func (api *ListPlugins) ContextEdit() *ListPlugins {
	return api.Context("edit")
}

func (api *ListPlugins) ContextEmbed() *ListPlugins {
	return api.Context("embed")
}

func (api *ListPlugins) Fields(fields ...string) *ListPlugins {
	api.arguments["_fields"] = strings.Join(fields, ",")
	return api
}

func (api *ListPlugins) Do() (plugins []Plugin, err error) {
	endpoint := api.client.endpoint + api.endpoint

	restyClient := api.client.httpClient.R()
	if api.client.auth.Username != "" && api.client.auth.Password != "" {
		restyClient.SetBasicAuth(api.client.auth.Username, api.client.auth.Password)
	}

	resp, err := restyClient.
		SetHeader("Accept", "application/json").
		SetResult(&plugins).
		SetQueryParams(api.arguments).
		Get(endpoint)

	if err != nil {
		return nil, err
	}

	if resp.IsError() {
		var wpError WPRestError
		err = json.Unmarshal(resp.Bytes(), &wpError)
		if err != nil {
			return nil, err
		}
		return nil, &wpError
	}

	return plugins, nil
}

// RetrievePlugin handles getting a single plugin.
type RetrievePlugin struct {
	endpoint  string
	client    *RestClient
	arguments map[string]string
}

func (api *RetrievePlugin) Context(ctx string) *RetrievePlugin {
	api.arguments["context"] = ctx
	return api
}

func (api *RetrievePlugin) ContextView() *RetrievePlugin {
	return api.Context("view")
}

func (api *RetrievePlugin) ContextEdit() *RetrievePlugin {
	return api.Context("edit")
}

func (api *RetrievePlugin) ContextEmbed() *RetrievePlugin {
	return api.Context("embed")
}

func (api *RetrievePlugin) Fields(fields ...string) *RetrievePlugin {
	api.arguments["_fields"] = strings.Join(fields, ",")
	return api
}

func (api *RetrievePlugin) Do() (plugin *Plugin, err error) {
	endpoint := api.client.endpoint + api.endpoint

	restyClient := api.client.httpClient.R()
	if api.client.auth.Username != "" && api.client.auth.Password != "" {
		restyClient.SetBasicAuth(api.client.auth.Username, api.client.auth.Password)
	}

	resp, err := restyClient.
		SetHeader("Accept", "application/json").
		SetResult(&plugin).
		SetQueryParams(api.arguments).
		Get(endpoint)

	if err != nil {
		return nil, err
	}

	if resp.IsError() {
		var wpError WPRestError
		err = json.Unmarshal(resp.Bytes(), &wpError)
		if err != nil {
			return nil, err
		}
		return nil, &wpError
	}

	return plugin, nil
}

// CreatePlugin handles installing a plugin from WordPress.org.
type CreatePlugin struct {
	endpoint string
	client   *RestClient
	body     map[string]any
}

func (api *CreatePlugin) Slug(slug string) *CreatePlugin {
	api.body["slug"] = slug
	return api
}

func (api *CreatePlugin) Status(status string) *CreatePlugin {
	api.body["status"] = status
	return api
}

func (api *CreatePlugin) StatusActive() *CreatePlugin {
	return api.Status(PluginStatusActive)
}

func (api *CreatePlugin) StatusInactive() *CreatePlugin {
	return api.Status(PluginStatusInactive)
}

func (api *CreatePlugin) Do() (plugin *Plugin, err error) {
	endpoint := api.client.endpoint + api.endpoint

	restyClient := api.client.httpClient.R()
	if api.client.auth.Username != "" && api.client.auth.Password != "" {
		restyClient.SetBasicAuth(api.client.auth.Username, api.client.auth.Password)
	}

	resp, err := restyClient.
		SetHeader("Content-Type", "application/json").
		SetHeader("Accept", "application/json").
		SetBody(api.body).
		SetResult(&plugin).
		Post(endpoint)

	if err != nil {
		return nil, err
	}

	if resp.IsError() {
		var wpError WPRestError
		err = json.Unmarshal(resp.Bytes(), &wpError)
		if err != nil {
			return nil, err
		}
		return nil, &wpError
	}

	return plugin, nil
}

// UpdatePlugin handles updating a plugin's status (activation/deactivation).
type UpdatePlugin struct {
	endpoint  string
	client    *RestClient
	arguments map[string]string
	body      map[string]any
}

func (api *UpdatePlugin) Status(status string) *UpdatePlugin {
	api.body["status"] = status
	return api
}

func (api *UpdatePlugin) StatusActive() *UpdatePlugin {
	return api.Status(PluginStatusActive)
}

func (api *UpdatePlugin) StatusInactive() *UpdatePlugin {
	return api.Status(PluginStatusInactive)
}

func (api *UpdatePlugin) Activate() *UpdatePlugin {
	return api.StatusActive()
}

func (api *UpdatePlugin) Deactivate() *UpdatePlugin {
	return api.StatusInactive()
}

func (api *UpdatePlugin) Context(ctx string) *UpdatePlugin {
	api.arguments["context"] = ctx
	return api
}

func (api *UpdatePlugin) ContextEdit() *UpdatePlugin {
	return api.Context("edit")
}

func (api *UpdatePlugin) Do() (plugin *Plugin, err error) {
	endpoint := api.client.endpoint + api.endpoint

	restyClient := api.client.httpClient.R()
	if api.client.auth.Username != "" && api.client.auth.Password != "" {
		restyClient.SetBasicAuth(api.client.auth.Username, api.client.auth.Password)
	}

	resp, err := restyClient.
		SetHeader("Content-Type", "application/json").
		SetHeader("Accept", "application/json").
		SetBody(api.body).
		SetQueryParams(api.arguments).
		SetResult(&plugin).
		Post(endpoint)

	if err != nil {
		return nil, err
	}

	if resp.IsError() {
		var wpError WPRestError
		err = json.Unmarshal(resp.Bytes(), &wpError)
		if err != nil {
			return nil, err
		}
		return nil, &wpError
	}

	return plugin, nil
}

type deletePluginEnvelope struct {
	Plugin
	Deleted  bool    `json:"deleted"`
	Previous *Plugin `json:"previous"`
}

// DeletePlugin handles deleting an installed plugin.
type DeletePlugin struct {
	endpoint  string
	client    *RestClient
	arguments map[string]string
}

func (api *DeletePlugin) Context(ctx string) *DeletePlugin {
	api.arguments["context"] = ctx
	return api
}

func (api *DeletePlugin) Do() (plugin *Plugin, err error) {
	endpoint := api.client.endpoint + api.endpoint

	restyClient := api.client.httpClient.R()
	if api.client.auth.Username != "" && api.client.auth.Password != "" {
		restyClient.SetBasicAuth(api.client.auth.Username, api.client.auth.Password)
	}

	var env deletePluginEnvelope
	resp, err := restyClient.
		SetHeader("Accept", "application/json").
		SetQueryParams(api.arguments).
		SetResult(&env).
		Delete(endpoint)

	if err != nil {
		return nil, err
	}

	if resp.IsError() {
		var wpError WPRestError
		err = json.Unmarshal(resp.Bytes(), &wpError)
		if err != nil {
			return nil, err
		}
		return nil, &wpError
	}

	if env.Deleted && env.Previous != nil && env.Previous.Plugin != "" {
		return env.Previous, nil
	}

	return &env.Plugin, nil
}
