package gowprest

import (
	"encoding/json"
	"strings"
)

// WidgetType represents a WordPress widget type record (/wp/v2/widget-types).
type WidgetType struct {
	ID          string         `json:"id,omitempty"`
	Name        string         `json:"name,omitempty"`
	Description string         `json:"description,omitempty"`
	IsMulti     bool           `json:"is_multi,omitempty"`
	Classname   string         `json:"classname,omitempty"`
	Links       map[string]any `json:"_links,omitempty"`
}

// WidgetTypes provides access to WordPress widget types endpoints (/wp/v2/widget-types).
type WidgetTypes struct {
	client *RestClient
}

// WidgetTypes returns a WidgetTypes service instance.
func (c *RestClient) WidgetTypes() *WidgetTypes {
	return &WidgetTypes{client: c}
}

// WidgetType is an ergonomic alias for WidgetTypes.
func (c *RestClient) WidgetType() *WidgetTypes {
	return c.WidgetTypes()
}

// List returns a ListWidgetTypes builder to query registered widget types.
func (api *WidgetTypes) List() *ListWidgetTypes {
	return &ListWidgetTypes{
		endpoint:  "/wp/v2/widget-types",
		client:    api.client,
		arguments: make(map[string]string),
	}
}

// Retrieve returns a RetrieveWidgetType builder for a specific widget type ID.
func (api *WidgetTypes) Retrieve(id string) *RetrieveWidgetType {
	return &RetrieveWidgetType{
		endpoint:  "/wp/v2/widget-types/" + id,
		client:    api.client,
		arguments: make(map[string]string),
	}
}

// ListWidgetTypes handles querying widget types.
type ListWidgetTypes struct {
	endpoint  string
	client    *RestClient
	arguments map[string]string
}

// Context sets the context parameter (e.g. "view", "edit", "embed").
func (api *ListWidgetTypes) Context(ctx string) *ListWidgetTypes {
	api.arguments["context"] = ctx
	return api
}

// ContextView sets the context parameter to "view".
func (api *ListWidgetTypes) ContextView() *ListWidgetTypes {
	return api.Context("view")
}

// ContextEdit sets the context parameter to "edit".
func (api *ListWidgetTypes) ContextEdit() *ListWidgetTypes {
	return api.Context("edit")
}

// ContextEmbed sets the context parameter to "embed".
func (api *ListWidgetTypes) ContextEmbed() *ListWidgetTypes {
	return api.Context("embed")
}

// Fields limits the response to specific fields.
func (api *ListWidgetTypes) Fields(fields ...string) *ListWidgetTypes {
	api.arguments["_fields"] = strings.Join(fields, ",")
	return api
}

// Do executes the list request.
func (api *ListWidgetTypes) Do() (widgetTypes []WidgetType, err error) {
	endpoint := api.client.endpoint + api.endpoint

	restyClient := api.client.httpClient.R()
	if api.client.auth.Username != "" && api.client.auth.Password != "" {
		restyClient = restyClient.SetBasicAuth(api.client.auth.Username, api.client.auth.Password)
	}

	resp, err := restyClient.
		SetHeader("Accept", "application/json").
		SetResult(&widgetTypes).
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

	return widgetTypes, nil
}

// RetrieveWidgetType handles retrieving a single widget type.
type RetrieveWidgetType struct {
	endpoint  string
	client    *RestClient
	arguments map[string]string
}

// Context sets the context parameter (e.g. "view", "edit", "embed").
func (api *RetrieveWidgetType) Context(ctx string) *RetrieveWidgetType {
	api.arguments["context"] = ctx
	return api
}

// ContextView sets the context parameter to "view".
func (api *RetrieveWidgetType) ContextView() *RetrieveWidgetType {
	return api.Context("view")
}

// ContextEdit sets the context parameter to "edit".
func (api *RetrieveWidgetType) ContextEdit() *RetrieveWidgetType {
	return api.Context("edit")
}

// ContextEmbed sets the context parameter to "embed".
func (api *RetrieveWidgetType) ContextEmbed() *RetrieveWidgetType {
	return api.Context("embed")
}

// Fields limits the response to specific fields.
func (api *RetrieveWidgetType) Fields(fields ...string) *RetrieveWidgetType {
	api.arguments["_fields"] = strings.Join(fields, ",")
	return api
}

// Do executes the retrieve request.
func (api *RetrieveWidgetType) Do() (widgetType *WidgetType, err error) {
	endpoint := api.client.endpoint + api.endpoint

	restyClient := api.client.httpClient.R()
	if api.client.auth.Username != "" && api.client.auth.Password != "" {
		restyClient = restyClient.SetBasicAuth(api.client.auth.Username, api.client.auth.Password)
	}

	resp, err := restyClient.
		SetHeader("Accept", "application/json").
		SetResult(&widgetType).
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

	return widgetType, nil
}
