package gowprest

import (
	"encoding/json"
	"strconv"
	"strings"
)

// PatternDirectoryItem represents an item from the WordPress.org Pattern Directory.
type PatternDirectoryItem struct {
	ID            int      `json:"id,omitempty"`
	Title         string   `json:"title,omitempty"`
	Content       string   `json:"content,omitempty"`
	Categories    []string `json:"categories,omitempty"`
	Keywords      []string `json:"keywords,omitempty"`
	Description   string   `json:"description,omitempty"`
	ViewportWidth int      `json:"viewport_width,omitempty"`
	BlockTypes    []string `json:"block_types,omitempty"`
}

// PatternDirectory provides access to WordPress.org Pattern Directory APIs (/wp/v2/pattern-directory/patterns).
type PatternDirectory struct {
	client *RestClient
}

// PatternDirectory returns a PatternDirectory service instance.
func (c *RestClient) PatternDirectory() *PatternDirectory {
	return &PatternDirectory{client: c}
}

// PatternDirectoryItems is a convenience alias for PatternDirectory.
func (c *RestClient) PatternDirectoryItems() *PatternDirectory {
	return c.PatternDirectory()
}

// List returns a ListPatternDirectoryItems builder to query pattern directory items.
func (api *PatternDirectory) List() *ListPatternDirectoryItems {
	return &ListPatternDirectoryItems{
		endpoint:  "/wp/v2/pattern-directory/patterns",
		client:    api.client,
		arguments: make(map[string]string),
	}
}

// Search returns a ListPatternDirectoryItems builder with a search query preconfigured.
func (api *PatternDirectory) Search(query string) *ListPatternDirectoryItems {
	return api.List().Search(query)
}

// ListPatternDirectoryItems handles querying the pattern directory.
type ListPatternDirectoryItems struct {
	endpoint  string
	client    *RestClient
	arguments map[string]string
}

// Context sets the context parameter (e.g. "view", "edit", "embed").
func (api *ListPatternDirectoryItems) Context(ctx string) *ListPatternDirectoryItems {
	api.arguments["context"] = ctx
	return api
}

// ContextView sets the context parameter to "view".
func (api *ListPatternDirectoryItems) ContextView() *ListPatternDirectoryItems {
	return api.Context("view")
}

// ContextEdit sets the context parameter to "edit".
func (api *ListPatternDirectoryItems) ContextEdit() *ListPatternDirectoryItems {
	return api.Context("edit")
}

// ContextEmbed sets the context parameter to "embed".
func (api *ListPatternDirectoryItems) ContextEmbed() *ListPatternDirectoryItems {
	return api.Context("embed")
}

// Page sets the page number of the collection.
func (api *ListPatternDirectoryItems) Page(page int) *ListPatternDirectoryItems {
	api.arguments["page"] = strconv.Itoa(page)
	return api
}

// PerPage sets the maximum number of items to return in the result set.
func (api *ListPatternDirectoryItems) PerPage(perPage int) *ListPatternDirectoryItems {
	api.arguments["per_page"] = strconv.Itoa(perPage)
	return api
}

// Search limits results to those matching a string.
func (api *ListPatternDirectoryItems) Search(query string) *ListPatternDirectoryItems {
	api.arguments["search"] = query
	return api
}

// Category limits results to those matching a category ID.
func (api *ListPatternDirectoryItems) Category(categoryID int) *ListPatternDirectoryItems {
	api.arguments["category"] = strconv.Itoa(categoryID)
	return api
}

// Keyword limits results to those matching a keyword ID.
func (api *ListPatternDirectoryItems) Keyword(keywordID int) *ListPatternDirectoryItems {
	api.arguments["keyword"] = strconv.Itoa(keywordID)
	return api
}

// Slug limits results to those matching a pattern slug.
func (api *ListPatternDirectoryItems) Slug(slug string) *ListPatternDirectoryItems {
	api.arguments["slug"] = slug
	return api
}

// Offset offsets the result set by a specific number of items.
func (api *ListPatternDirectoryItems) Offset(offset int) *ListPatternDirectoryItems {
	api.arguments["offset"] = strconv.Itoa(offset)
	return api
}

// Order sets the sort order ("asc" or "desc").
func (api *ListPatternDirectoryItems) Order(order string) *ListPatternDirectoryItems {
	api.arguments["order"] = order
	return api
}

// OrderAsc sets the sort order to ascending.
func (api *ListPatternDirectoryItems) OrderAsc() *ListPatternDirectoryItems {
	return api.Order("asc")
}

// OrderDesc sets the sort order to descending.
func (api *ListPatternDirectoryItems) OrderDesc() *ListPatternDirectoryItems {
	return api.Order("desc")
}

// OrderBy sorts the collection by a pattern attribute.
func (api *ListPatternDirectoryItems) OrderBy(field string) *ListPatternDirectoryItems {
	api.arguments["orderby"] = field
	return api
}

// Fields limits the response to specific fields.
func (api *ListPatternDirectoryItems) Fields(fields ...string) *ListPatternDirectoryItems {
	api.arguments["_fields"] = strings.Join(fields, ",")
	return api
}

// Do executes the request and returns the pattern directory items.
func (api *ListPatternDirectoryItems) Do() (items []PatternDirectoryItem, err error) {
	endpoint := api.client.endpoint + api.endpoint

	restyClient := api.client.httpClient.R()
	if api.client.auth.Username != "" && api.client.auth.Password != "" {
		restyClient = restyClient.SetBasicAuth(api.client.auth.Username, api.client.auth.Password)
	}

	resp, err := restyClient.
		SetHeader("Accept", "application/json").
		SetResult(&items).
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

	return items, nil
}
