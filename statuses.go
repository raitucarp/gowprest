package gowprest

import (
	"encoding/json"
	"strings"
)

// Status represents a WordPress post status.
type Status struct {
	Name         string         `json:"name,omitempty"`
	Private      bool           `json:"private,omitempty"`
	Protected    bool           `json:"protected,omitempty"`
	Public       bool           `json:"public,omitempty"`
	Queryable    bool           `json:"queryable,omitempty"`
	ShowInList   bool           `json:"show_in_list,omitempty"`
	Slug         string         `json:"slug,omitempty"`
	DateFloating bool           `json:"date_floating,omitempty"`
	Links        map[string]any `json:"_links,omitempty"`
	Embedded     map[string]any `json:"_embedded,omitempty"`
}

// Statuses handles requests to the WordPress post statuses API.
type Statuses struct {
	client *RestClient
}

// Statuses returns a Statuses service instance.
func (c *RestClient) Statuses() *Statuses {
	return &Statuses{client: c}
}

// PostStatuses is an alias for Statuses.
func (c *RestClient) PostStatuses() *Statuses {
	return &Statuses{client: c}
}

// List returns a ListStatuses builder to query post statuses.
func (api *Statuses) List() *ListStatuses {
	return &ListStatuses{
		endpoint:  "/wp/v2/statuses",
		client:    api.client,
		arguments: make(map[string]string),
	}
}

// Retrieve returns a RetrieveStatus builder to get a specific post status by slug.
func (api *Statuses) Retrieve(status string) *RetrieveStatus {
	return &RetrieveStatus{
		endpoint:  "/wp/v2/statuses/" + status,
		client:    api.client,
		arguments: make(map[string]string),
	}
}

// ListStatuses handles querying post status collections.
type ListStatuses struct {
	endpoint  string
	client    *RestClient
	arguments map[string]string
}

func (api *ListStatuses) Context(ctx string) *ListStatuses {
	api.arguments["context"] = ctx
	return api
}

func (api *ListStatuses) ContextView() *ListStatuses {
	api.arguments["context"] = "view"
	return api
}

func (api *ListStatuses) ContextEdit() *ListStatuses {
	api.arguments["context"] = "edit"
	return api
}

func (api *ListStatuses) ContextEmbed() *ListStatuses {
	api.arguments["context"] = "embed"
	return api
}

func (api *ListStatuses) Embed() *ListStatuses {
	api.arguments["_embed"] = "true"
	return api
}

func (api *ListStatuses) Fields(fields ...string) *ListStatuses {
	api.arguments["_fields"] = strings.Join(fields, ",")
	return api
}

func (api *ListStatuses) Do() (statuses map[string]Status, err error) {
	endpoint := api.client.endpoint + api.endpoint

	restyClient := api.client.httpClient.R()
	if api.client.auth.Username != "" && api.client.auth.Password != "" {
		restyClient.SetBasicAuth(api.client.auth.Username, api.client.auth.Password)
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
		return make(map[string]Status), nil
	}

	err = json.Unmarshal(raw, &statuses)
	if err != nil {
		return nil, err
	}
	return statuses, nil
}

// RetrieveStatus handles querying a single post status.
type RetrieveStatus struct {
	endpoint  string
	client    *RestClient
	arguments map[string]string
}

func (api *RetrieveStatus) Context(ctx string) *RetrieveStatus {
	api.arguments["context"] = ctx
	return api
}

func (api *RetrieveStatus) ContextView() *RetrieveStatus {
	api.arguments["context"] = "view"
	return api
}

func (api *RetrieveStatus) ContextEdit() *RetrieveStatus {
	api.arguments["context"] = "edit"
	return api
}

func (api *RetrieveStatus) ContextEmbed() *RetrieveStatus {
	api.arguments["context"] = "embed"
	return api
}

func (api *RetrieveStatus) Embed() *RetrieveStatus {
	api.arguments["_embed"] = "true"
	return api
}

func (api *RetrieveStatus) Fields(fields ...string) *RetrieveStatus {
	api.arguments["_fields"] = strings.Join(fields, ",")
	return api
}

func (api *RetrieveStatus) Do() (status *Status, err error) {
	endpoint := api.client.endpoint + api.endpoint

	restyClient := api.client.httpClient.R()
	if api.client.auth.Username != "" && api.client.auth.Password != "" {
		restyClient.SetBasicAuth(api.client.auth.Username, api.client.auth.Password)
	}

	resp, err := restyClient.
		SetHeader("Accept", "application/json").
		SetResult(&status).
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

	return status, nil
}
