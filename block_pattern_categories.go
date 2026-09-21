package gowprest

import (
	"encoding/json"
	"strings"
)

// BlockPatternCategory represents a registered WordPress block pattern category.
type BlockPatternCategory struct {
	Name        string `json:"name,omitempty"`
	Label       string `json:"label,omitempty"`
	Description string `json:"description,omitempty"`
}

// BlockPatternCategories handles requests to /wp/v2/block-patterns/categories.
type BlockPatternCategories struct {
	client *RestClient
}

// BlockPatterns provides access to block pattern related APIs.
type BlockPatterns struct {
	client *RestClient
}

// BlockPatternCategories returns a BlockPatternCategories service instance.
func (c *RestClient) BlockPatternCategories() *BlockPatternCategories {
	return &BlockPatternCategories{client: c}
}

// BlockPatterns returns a BlockPatterns service instance.
func (c *RestClient) BlockPatterns() *BlockPatterns {
	return &BlockPatterns{client: c}
}

// Categories returns a BlockPatternCategories service instance.
func (api *BlockPatterns) Categories() *BlockPatternCategories {
	return &BlockPatternCategories{client: api.client}
}

// List returns a ListBlockPatternCategories builder.
func (api *BlockPatternCategories) List() *ListBlockPatternCategories {
	return &ListBlockPatternCategories{
		endpoint:  "/wp/v2/block-patterns/categories",
		client:    api.client,
		arguments: make(map[string]string),
	}
}

// Retrieve returns a RetrieveBlockPatternCategory builder to find a specific category by name.
func (api *BlockPatternCategories) Retrieve(name string) *RetrieveBlockPatternCategory {
	return &RetrieveBlockPatternCategory{
		endpoint:     "/wp/v2/block-patterns/categories",
		categoryName: name,
		client:       api.client,
		arguments:    make(map[string]string),
	}
}

// ListBlockPatternCategories handles querying block pattern categories.
type ListBlockPatternCategories struct {
	endpoint  string
	client    *RestClient
	arguments map[string]string
}

// Context sets the context parameter (e.g., "view", "edit", "embed").
func (api *ListBlockPatternCategories) Context(ctx string) *ListBlockPatternCategories {
	api.arguments["context"] = ctx
	return api
}

// ContextView sets the context parameter to "view".
func (api *ListBlockPatternCategories) ContextView() *ListBlockPatternCategories {
	return api.Context("view")
}

// ContextEdit sets the context parameter to "edit".
func (api *ListBlockPatternCategories) ContextEdit() *ListBlockPatternCategories {
	return api.Context("edit")
}

// ContextEmbed sets the context parameter to "embed".
func (api *ListBlockPatternCategories) ContextEmbed() *ListBlockPatternCategories {
	return api.Context("embed")
}

// Fields limits the response to specific fields.
func (api *ListBlockPatternCategories) Fields(fields ...string) *ListBlockPatternCategories {
	api.arguments["_fields"] = strings.Join(fields, ",")
	return api
}

// Do executes the request and returns the list of block pattern categories.
func (api *ListBlockPatternCategories) Do() (categories []BlockPatternCategory, err error) {
	endpoint := api.client.endpoint + api.endpoint

	restyClient := api.client.httpClient.R()
	if api.client.auth.Username != "" && api.client.auth.Password != "" {
		restyClient = restyClient.SetBasicAuth(api.client.auth.Username, api.client.auth.Password)
	}

	resp, err := restyClient.
		SetHeader("Accept", "application/json").
		SetResult(&categories).
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

	return categories, nil
}

// RetrieveBlockPatternCategory handles retrieving a specific block pattern category by name.
type RetrieveBlockPatternCategory struct {
	endpoint     string
	categoryName string
	client       *RestClient
	arguments    map[string]string
}

// Context sets the context parameter (e.g., "view", "edit", "embed").
func (api *RetrieveBlockPatternCategory) Context(ctx string) *RetrieveBlockPatternCategory {
	api.arguments["context"] = ctx
	return api
}

// ContextView sets the context parameter to "view".
func (api *RetrieveBlockPatternCategory) ContextView() *RetrieveBlockPatternCategory {
	return api.Context("view")
}

// ContextEdit sets the context parameter to "edit".
func (api *RetrieveBlockPatternCategory) ContextEdit() *RetrieveBlockPatternCategory {
	return api.Context("edit")
}

// ContextEmbed sets the context parameter to "embed".
func (api *RetrieveBlockPatternCategory) ContextEmbed() *RetrieveBlockPatternCategory {
	return api.Context("embed")
}

// Fields limits the response to specific fields.
func (api *RetrieveBlockPatternCategory) Fields(fields ...string) *RetrieveBlockPatternCategory {
	api.arguments["_fields"] = strings.Join(fields, ",")
	return api
}

// Do executes the request and returns the matched block pattern category, or a not-found error.
func (api *RetrieveBlockPatternCategory) Do() (*BlockPatternCategory, error) {
	listBuilder := &ListBlockPatternCategories{
		endpoint:  api.endpoint,
		client:    api.client,
		arguments: api.arguments,
	}

	categories, err := listBuilder.Do()
	if err != nil {
		return nil, err
	}

	for _, cat := range categories {
		if cat.Name == api.categoryName {
			return &cat, nil
		}
	}

	return nil, &WPRestError{
		Code:    "rest_block_pattern_category_not_found",
		Message: "Block pattern category not found.",
		Data: struct {
			Status int `json:"status"`
		}{
			Status: 404,
		},
	}
}
