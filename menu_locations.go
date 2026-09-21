package gowprest

import (
	"encoding/json"
	"strings"
)

// MenuLocation represents a registered WordPress navigation menu location.
type MenuLocation struct {
	Name        string         `json:"name,omitempty"`
	Description string         `json:"description,omitempty"`
	Menu        int            `json:"menu,omitempty"`
	Links       map[string]any `json:"_links,omitempty"`
}

// MenuLocations handles requests to /wp/v2/menu-locations.
type MenuLocations struct {
	client *RestClient
}

// MenuLocations returns a MenuLocations service instance.
func (c *RestClient) MenuLocations() *MenuLocations {
	return &MenuLocations{client: c}
}

// NavMenuLocations is a convenience alias for MenuLocations.
func (c *RestClient) NavMenuLocations() *MenuLocations {
	return c.MenuLocations()
}

// List returns a ListMenuLocations builder.
func (api *MenuLocations) List() *ListMenuLocations {
	return &ListMenuLocations{
		endpoint:  "/wp/v2/menu-locations",
		client:    api.client,
		arguments: make(map[string]string),
	}
}

// Retrieve returns a RetrieveMenuLocation builder for a specific menu location name.
func (api *MenuLocations) Retrieve(location string) *RetrieveMenuLocation {
	return &RetrieveMenuLocation{
		endpoint:  "/wp/v2/menu-locations/" + location,
		client:    api.client,
		arguments: make(map[string]string),
	}
}

// ListMenuLocations handles querying registered menu locations.
type ListMenuLocations struct {
	endpoint  string
	client    *RestClient
	arguments map[string]string
}

// Context sets the context parameter (e.g. "view", "edit", "embed").
func (api *ListMenuLocations) Context(ctx string) *ListMenuLocations {
	api.arguments["context"] = ctx
	return api
}

// ContextView sets the context parameter to "view".
func (api *ListMenuLocations) ContextView() *ListMenuLocations {
	return api.Context("view")
}

// ContextEdit sets the context parameter to "edit".
func (api *ListMenuLocations) ContextEdit() *ListMenuLocations {
	return api.Context("edit")
}

// ContextEmbed sets the context parameter to "embed".
func (api *ListMenuLocations) ContextEmbed() *ListMenuLocations {
	return api.Context("embed")
}

// Fields limits the response to specific fields.
func (api *ListMenuLocations) Fields(fields ...string) *ListMenuLocations {
	api.arguments["_fields"] = strings.Join(fields, ",")
	return api
}

// Do executes the request and returns the map of menu locations.
// WordPress returns an object map when locations exist, or an empty JSON array `[]` when none exist.
func (api *ListMenuLocations) Do() (locations map[string]MenuLocation, err error) {
	endpoint := api.client.endpoint + api.endpoint

	restyClient := api.client.httpClient.R()
	if api.client.auth.Username != "" && api.client.auth.Password != "" {
		restyClient = restyClient.SetBasicAuth(api.client.auth.Username, api.client.auth.Password)
	}

	var raw json.RawMessage
	resp, err := restyClient.
		SetHeader("Accept", "application/json").
		SetResult(&raw).
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

	body := strings.TrimSpace(string(raw))
	if strings.HasPrefix(body, "[") {
		return make(map[string]MenuLocation), nil
	}

	err = json.Unmarshal(raw, &locations)
	if err != nil {
		return nil, err
	}

	return locations, nil
}

// RetrieveMenuLocation handles retrieving a single menu location.
type RetrieveMenuLocation struct {
	endpoint  string
	client    *RestClient
	arguments map[string]string
}

// Context sets the context parameter (e.g. "view", "edit", "embed").
func (api *RetrieveMenuLocation) Context(ctx string) *RetrieveMenuLocation {
	api.arguments["context"] = ctx
	return api
}

// ContextView sets the context parameter to "view".
func (api *RetrieveMenuLocation) ContextView() *RetrieveMenuLocation {
	return api.Context("view")
}

// ContextEdit sets the context parameter to "edit".
func (api *RetrieveMenuLocation) ContextEdit() *RetrieveMenuLocation {
	return api.Context("edit")
}

// ContextEmbed sets the context parameter to "embed".
func (api *RetrieveMenuLocation) ContextEmbed() *RetrieveMenuLocation {
	return api.Context("embed")
}

// Fields limits the response to specific fields.
func (api *RetrieveMenuLocation) Fields(fields ...string) *RetrieveMenuLocation {
	api.arguments["_fields"] = strings.Join(fields, ",")
	return api
}

// Do executes the request and returns the menu location.
func (api *RetrieveMenuLocation) Do() (location *MenuLocation, err error) {
	endpoint := api.client.endpoint + api.endpoint

	restyClient := api.client.httpClient.R()
	if api.client.auth.Username != "" && api.client.auth.Password != "" {
		restyClient = restyClient.SetBasicAuth(api.client.auth.Username, api.client.auth.Password)
	}

	resp, err := restyClient.
		SetHeader("Accept", "application/json").
		SetResult(&location).
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

	return location, nil
}
