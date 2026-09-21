package gowprest

import (
	"encoding/json"
	"strings"
)

// Sidebar represents a WordPress sidebar (widget area).
type Sidebar struct {
	ID           string         `json:"id,omitempty"`
	Name         string         `json:"name,omitempty"`
	Description  string         `json:"description,omitempty"`
	Class        string         `json:"class,omitempty"`
	BeforeWidget string         `json:"before_widget,omitempty"`
	AfterWidget  string         `json:"after_widget,omitempty"`
	BeforeTitle  string         `json:"before_title,omitempty"`
	AfterTitle   string         `json:"after_title,omitempty"`
	Status       string         `json:"status,omitempty"`
	Widgets      []any          `json:"widgets,omitempty"`
	Links        map[string]any `json:"_links,omitempty"`
}

// Sidebars provides access to sidebars APIs (/wp/v2/sidebars).
type Sidebars struct {
	client *RestClient
}

// Sidebars returns a Sidebars service instance.
func (c *RestClient) Sidebars() *Sidebars {
	return &Sidebars{client: c}
}

// Sidebar is an ergonomic alias for Sidebars.
func (c *RestClient) Sidebar() *Sidebars {
	return c.Sidebars()
}

// List returns a ListSidebars builder to list registered sidebars.
func (api *Sidebars) List() *ListSidebars {
	return &ListSidebars{
		endpoint:  "/wp/v2/sidebars",
		client:    api.client,
		arguments: make(map[string]string),
	}
}

// Retrieve returns a RetrieveSidebar builder for a specific sidebar ID.
func (api *Sidebars) Retrieve(id string) *RetrieveSidebar {
	return &RetrieveSidebar{
		endpoint:  "/wp/v2/sidebars/" + id,
		client:    api.client,
		arguments: make(map[string]string),
	}
}

// Update returns an UpdateSidebar builder for a specific sidebar ID.
func (api *Sidebars) Update(id string) *UpdateSidebar {
	return &UpdateSidebar{
		endpoint: "/wp/v2/sidebars/" + id,
		client:   api.client,
		body:     make(map[string]any),
	}
}

// ListSidebars handles querying registered sidebars.
type ListSidebars struct {
	endpoint  string
	client    *RestClient
	arguments map[string]string
}

// Context sets the context parameter (e.g. "view", "edit", "embed").
func (api *ListSidebars) Context(ctx string) *ListSidebars {
	api.arguments["context"] = ctx
	return api
}

// ContextView sets the context parameter to "view".
func (api *ListSidebars) ContextView() *ListSidebars {
	return api.Context("view")
}

// ContextEdit sets the context parameter to "edit".
func (api *ListSidebars) ContextEdit() *ListSidebars {
	return api.Context("edit")
}

// ContextEmbed sets the context parameter to "embed".
func (api *ListSidebars) ContextEmbed() *ListSidebars {
	return api.Context("embed")
}

// Fields limits the response to specific fields.
func (api *ListSidebars) Fields(fields ...string) *ListSidebars {
	api.arguments["_fields"] = strings.Join(fields, ",")
	return api
}

// Do executes the list request.
func (api *ListSidebars) Do() (sidebars []Sidebar, err error) {
	endpoint := api.client.endpoint + api.endpoint

	restyClient := api.client.httpClient.R()
	if api.client.auth.Username != "" && api.client.auth.Password != "" {
		restyClient = restyClient.SetBasicAuth(api.client.auth.Username, api.client.auth.Password)
	}

	resp, err := restyClient.
		SetHeader("Accept", "application/json").
		SetResult(&sidebars).
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

	return sidebars, nil
}

// RetrieveSidebar handles retrieving a single sidebar.
type RetrieveSidebar struct {
	endpoint  string
	client    *RestClient
	arguments map[string]string
}

// Context sets the context parameter (e.g. "view", "edit", "embed").
func (api *RetrieveSidebar) Context(ctx string) *RetrieveSidebar {
	api.arguments["context"] = ctx
	return api
}

// ContextView sets the context parameter to "view".
func (api *RetrieveSidebar) ContextView() *RetrieveSidebar {
	return api.Context("view")
}

// ContextEdit sets the context parameter to "edit".
func (api *RetrieveSidebar) ContextEdit() *RetrieveSidebar {
	return api.Context("edit")
}

// ContextEmbed sets the context parameter to "embed".
func (api *RetrieveSidebar) ContextEmbed() *RetrieveSidebar {
	return api.Context("embed")
}

// Fields limits the response to specific fields.
func (api *RetrieveSidebar) Fields(fields ...string) *RetrieveSidebar {
	api.arguments["_fields"] = strings.Join(fields, ",")
	return api
}

// Do executes the retrieve request.
func (api *RetrieveSidebar) Do() (sidebar *Sidebar, err error) {
	endpoint := api.client.endpoint + api.endpoint

	restyClient := api.client.httpClient.R()
	if api.client.auth.Username != "" && api.client.auth.Password != "" {
		restyClient = restyClient.SetBasicAuth(api.client.auth.Username, api.client.auth.Password)
	}

	resp, err := restyClient.
		SetHeader("Accept", "application/json").
		SetResult(&sidebar).
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

	return sidebar, nil
}

// UpdateSidebar handles modifying an existing sidebar's assigned widgets.
type UpdateSidebar struct {
	endpoint string
	client   *RestClient
	body     map[string]any
}

// Widgets sets the nested widgets assigned to the sidebar.
func (api *UpdateSidebar) Widgets(widgets ...any) *UpdateSidebar {
	api.body["widgets"] = widgets
	return api
}

// Do executes the update request.
func (api *UpdateSidebar) Do() (sidebar *Sidebar, err error) {
	endpoint := api.client.endpoint + api.endpoint

	restyClient := api.client.httpClient.R()
	if api.client.auth.Username != "" && api.client.auth.Password != "" {
		restyClient = restyClient.SetBasicAuth(api.client.auth.Username, api.client.auth.Password)
	}

	resp, err := restyClient.
		SetHeader("Content-Type", "application/json").
		SetHeader("Accept", "application/json").
		SetBody(api.body).
		SetResult(&sidebar).
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

	return sidebar, nil
}
