package gowprest

import (
	"encoding/json"
	"strings"
)

// Block pattern source constants.
const (
	BlockPatternSourceCore                     = "core"
	BlockPatternSourcePlugin                   = "plugin"
	BlockPatternSourceTheme                    = "theme"
	BlockPatternSourcePatternDirectoryCore     = "pattern-directory/core"
	BlockPatternSourcePatternDirectoryTheme    = "pattern-directory/theme"
	BlockPatternSourcePatternDirectoryFeatured = "pattern-directory/featured"
)

// BlockPattern represents a registered WordPress block pattern.
type BlockPattern struct {
	Name          string   `json:"name,omitempty"`
	Title         string   `json:"title,omitempty"`
	Content       string   `json:"content,omitempty"`
	Description   string   `json:"description,omitempty"`
	ViewportWidth int      `json:"viewport_width,omitempty"`
	Inserter      *bool    `json:"inserter,omitempty"`
	Categories    []string `json:"categories,omitempty"`
	Keywords      []string `json:"keywords,omitempty"`
	BlockTypes    []string `json:"block_types,omitempty"`
	PostTypes     []string `json:"post_types,omitempty"`
	TemplateTypes []string `json:"template_types,omitempty"`
	Source        string   `json:"source,omitempty"`
}

// Patterns is an alias returning the BlockPatterns service instance.
func (api *BlockPatterns) Patterns() *BlockPatterns {
	return api
}

// List returns a ListBlockPatterns builder.
func (api *BlockPatterns) List() *ListBlockPatterns {
	return &ListBlockPatterns{
		endpoint:  "/wp/v2/block-patterns/patterns",
		client:    api.client,
		arguments: make(map[string]string),
	}
}

// Retrieve returns a RetrieveBlockPattern builder to find a specific pattern by name.
func (api *BlockPatterns) Retrieve(name string) *RetrieveBlockPattern {
	return &RetrieveBlockPattern{
		endpoint:    "/wp/v2/block-patterns/patterns",
		patternName: name,
		client:      api.client,
		arguments:   make(map[string]string),
	}
}

// ListBlockPatterns handles querying block patterns.
type ListBlockPatterns struct {
	endpoint  string
	client    *RestClient
	arguments map[string]string
}

// Context sets the context parameter (e.g., "view", "edit", "embed").
func (api *ListBlockPatterns) Context(ctx string) *ListBlockPatterns {
	api.arguments["context"] = ctx
	return api
}

// ContextView sets the context parameter to "view".
func (api *ListBlockPatterns) ContextView() *ListBlockPatterns {
	return api.Context("view")
}

// ContextEdit sets the context parameter to "edit".
func (api *ListBlockPatterns) ContextEdit() *ListBlockPatterns {
	return api.Context("edit")
}

// ContextEmbed sets the context parameter to "embed".
func (api *ListBlockPatterns) ContextEmbed() *ListBlockPatterns {
	return api.Context("embed")
}

// Fields limits the response to specific fields.
func (api *ListBlockPatterns) Fields(fields ...string) *ListBlockPatterns {
	api.arguments["_fields"] = strings.Join(fields, ",")
	return api
}

// Do executes the request and returns the list of block patterns.
func (api *ListBlockPatterns) Do() (patterns []BlockPattern, err error) {
	endpoint := api.client.endpoint + api.endpoint

	restyClient := api.client.httpClient.R()
	if api.client.auth.Username != "" && api.client.auth.Password != "" {
		restyClient = restyClient.SetBasicAuth(api.client.auth.Username, api.client.auth.Password)
	}

	resp, err := restyClient.
		SetHeader("Accept", "application/json").
		SetResult(&patterns).
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

	return patterns, nil
}

// RetrieveBlockPattern handles retrieving a specific block pattern by name.
type RetrieveBlockPattern struct {
	endpoint    string
	patternName string
	client      *RestClient
	arguments   map[string]string
}

// Context sets the context parameter (e.g., "view", "edit", "embed").
func (api *RetrieveBlockPattern) Context(ctx string) *RetrieveBlockPattern {
	api.arguments["context"] = ctx
	return api
}

// ContextView sets the context parameter to "view".
func (api *RetrieveBlockPattern) ContextView() *RetrieveBlockPattern {
	return api.Context("view")
}

// ContextEdit sets the context parameter to "edit".
func (api *RetrieveBlockPattern) ContextEdit() *RetrieveBlockPattern {
	return api.Context("edit")
}

// ContextEmbed sets the context parameter to "embed".
func (api *RetrieveBlockPattern) ContextEmbed() *RetrieveBlockPattern {
	return api.Context("embed")
}

// Fields limits the response to specific fields.
func (api *RetrieveBlockPattern) Fields(fields ...string) *RetrieveBlockPattern {
	api.arguments["_fields"] = strings.Join(fields, ",")
	return api
}

// Do executes the request and returns the matched block pattern, or a not-found error.
func (api *RetrieveBlockPattern) Do() (*BlockPattern, error) {
	listBuilder := &ListBlockPatterns{
		endpoint:  api.endpoint,
		client:    api.client,
		arguments: api.arguments,
	}

	patterns, err := listBuilder.Do()
	if err != nil {
		return nil, err
	}

	for _, p := range patterns {
		if p.Name == api.patternName {
			return &p, nil
		}
	}

	return nil, &WPRestError{
		Code:    "rest_block_pattern_not_found",
		Message: "Block pattern not found.",
		Data: struct {
			Status int `json:"status"`
		}{
			Status: 404,
		},
	}
}
