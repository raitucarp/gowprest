package gowprest

import (
	"encoding/json"
	"strconv"
	"strings"
)

// BlockDirectoryItem represents an item from the WordPress.org Block Directory.
type BlockDirectoryItem struct {
	Name              string         `json:"name,omitempty"`
	Title             string         `json:"title,omitempty"`
	Description       string         `json:"description,omitempty"`
	ID                string         `json:"id,omitempty"`
	Rating            float64        `json:"rating,omitempty"`
	RatingCount       int            `json:"rating_count,omitempty"`
	ActiveInstalls    int            `json:"active_installs,omitempty"`
	AuthorBlockRating float64        `json:"author_block_rating,omitempty"`
	AuthorBlockCount  int            `json:"author_block_count,omitempty"`
	Author            string         `json:"author,omitempty"`
	Icon              string         `json:"icon,omitempty"`
	LastUpdated       string         `json:"last_updated,omitempty"`
	HumanizedUpdated  string         `json:"humanized_updated,omitempty"`
	Links             map[string]any `json:"_links,omitempty"`
}

// BlockDirectory handles searching the WordPress.org Block Directory (/wp/v2/block-directory/search).
type BlockDirectory struct {
	client *RestClient
}

// BlockDirectory returns a BlockDirectory service instance.
func (c *RestClient) BlockDirectory() *BlockDirectory {
	return &BlockDirectory{client: c}
}

// BlockDirectoryItems is a convenience alias for BlockDirectory.
func (c *RestClient) BlockDirectoryItems() *BlockDirectory {
	return c.BlockDirectory()
}

// DirectorySearch is a convenience alias for BlockDirectory.
func (c *RestClient) DirectorySearch() *BlockDirectory {
	return c.BlockDirectory()
}

// Search returns a SearchBlockDirectory builder for querying block directory items.
func (api *BlockDirectory) Search(term ...string) *SearchBlockDirectory {
	builder := &SearchBlockDirectory{
		endpoint:  "/wp/v2/block-directory/search",
		client:    api.client,
		arguments: make(map[string]string),
	}
	if len(term) > 0 {
		builder.Term(term[0])
	}
	return builder
}

// List returns a SearchBlockDirectory builder (alias for Search).
func (api *BlockDirectory) List(term ...string) *SearchBlockDirectory {
	return api.Search(term...)
}

// SearchBlockDirectory handles querying the block directory.
type SearchBlockDirectory struct {
	endpoint  string
	client    *RestClient
	arguments map[string]string
}

// Term sets the search term query parameter.
func (api *SearchBlockDirectory) Term(term string) *SearchBlockDirectory {
	api.arguments["term"] = term
	return api
}

// Page sets the page number of the results.
func (api *SearchBlockDirectory) Page(page int) *SearchBlockDirectory {
	api.arguments["page"] = strconv.Itoa(page)
	return api
}

// PerPage sets the number of results per page.
func (api *SearchBlockDirectory) PerPage(perPage int) *SearchBlockDirectory {
	api.arguments["per_page"] = strconv.Itoa(perPage)
	return api
}

// Context sets the context parameter (default is "view").
func (api *SearchBlockDirectory) Context(ctx string) *SearchBlockDirectory {
	api.arguments["context"] = ctx
	return api
}

// ContextView sets the context parameter to "view".
func (api *SearchBlockDirectory) ContextView() *SearchBlockDirectory {
	return api.Context("view")
}

// Fields limits the response to specific fields.
func (api *SearchBlockDirectory) Fields(fields ...string) *SearchBlockDirectory {
	api.arguments["_fields"] = strings.Join(fields, ",")
	return api
}

// Do executes the search request and returns the matching block directory items.
func (api *SearchBlockDirectory) Do() (items []BlockDirectoryItem, err error) {
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
