package gowprest

import (
	"encoding/json"
	"strconv"
	"strings"
)

// Search result object type constants.
const (
	SearchTypePost       = "post"
	SearchTypeTerm       = "term"
	SearchTypePostFormat = "post-format"
)

// Search result subtype constants.
const (
	SearchSubtypePost     = "post"
	SearchSubtypePage     = "page"
	SearchSubtypeCategory = "category"
	SearchSubtypePostTag  = "post_tag"
	SearchSubtypeAny      = "any"
)

// SearchResult represents a WordPress search result item.
type SearchResult struct {
	ID       any            `json:"id,omitempty"`
	Title    string         `json:"title,omitempty"`
	URL      string         `json:"url,omitempty"`
	Type     string         `json:"type,omitempty"`
	Subtype  string         `json:"subtype,omitempty"`
	Links    map[string]any `json:"_links,omitempty"`
	Embedded map[string]any `json:"_embedded,omitempty"`
}

// IntID returns the search result ID as an integer, converting from float64 or string if needed.
func (r *SearchResult) IntID() int {
	switch v := r.ID.(type) {
	case int:
		return v
	case float64:
		return int(v)
	case string:
		i, _ := strconv.Atoi(v)
		return i
	default:
		return 0
	}
}

// StringID returns the search result ID as a string.
func (r *SearchResult) StringID() string {
	switch v := r.ID.(type) {
	case string:
		return v
	case float64:
		return strconv.Itoa(int(v))
	case int:
		return strconv.Itoa(v)
	default:
		return ""
	}
}

// Search handles querying the WordPress search API.
type Search struct {
	endpoint  string
	client    *RestClient
	arguments map[string]string
}

// Search returns a Search builder to perform searches.
func (c *RestClient) Search() *Search {
	return &Search{
		endpoint:  "/wp/v2/search",
		client:    c,
		arguments: make(map[string]string),
	}
}

// SearchResults is an alias for Search.
func (c *RestClient) SearchResults() *Search {
	return c.Search()
}

// List returns the Search builder itself for consistent API ergonomics.
func (api *Search) List() *Search {
	return api
}

func (api *Search) Query(query string) *Search {
	api.arguments["search"] = query
	return api
}

func (api *Search) Search(query string) *Search {
	return api.Query(query)
}

func (api *Search) Page(page int) *Search {
	api.arguments["page"] = strconv.Itoa(page)
	return api
}

func (api *Search) PerPage(perPage int) *Search {
	api.arguments["per_page"] = strconv.Itoa(perPage)
	return api
}

func (api *Search) Type(objectType string) *Search {
	api.arguments["type"] = objectType
	return api
}

func (api *Search) TypePost() *Search {
	return api.Type(SearchTypePost)
}

func (api *Search) TypeTerm() *Search {
	return api.Type(SearchTypeTerm)
}

func (api *Search) TypePostFormat() *Search {
	return api.Type(SearchTypePostFormat)
}

func (api *Search) Subtype(subtypes ...string) *Search {
	api.arguments["subtype"] = strings.Join(subtypes, ",")
	return api
}

func (api *Search) SubtypePost() *Search {
	return api.Subtype(SearchSubtypePost)
}

func (api *Search) SubtypePage() *Search {
	return api.Subtype(SearchSubtypePage)
}

func (api *Search) SubtypeCategory() *Search {
	return api.Subtype(SearchSubtypeCategory)
}

func (api *Search) SubtypeTag() *Search {
	return api.Subtype(SearchSubtypePostTag)
}

func (api *Search) Exclude(ids ...int) *Search {
	strIDs := make([]string, len(ids))
	for i, id := range ids {
		strIDs[i] = strconv.Itoa(id)
	}
	api.arguments["exclude"] = strings.Join(strIDs, ",")
	return api
}

func (api *Search) Include(ids ...int) *Search {
	strIDs := make([]string, len(ids))
	for i, id := range ids {
		strIDs[i] = strconv.Itoa(id)
	}
	api.arguments["include"] = strings.Join(strIDs, ",")
	return api
}

func (api *Search) Context(ctx string) *Search {
	api.arguments["context"] = ctx
	return api
}

func (api *Search) ContextView() *Search {
	return api.Context("view")
}

func (api *Search) ContextEmbed() *Search {
	return api.Context("embed")
}

func (api *Search) Embed() *Search {
	api.arguments["_embed"] = "true"
	return api
}

func (api *Search) Fields(fields ...string) *Search {
	api.arguments["_fields"] = strings.Join(fields, ",")
	return api
}

func (api *Search) Do() (results []SearchResult, err error) {
	endpoint := api.client.endpoint + api.endpoint

	restyClient := api.client.httpClient.R()
	if api.client.auth.Username != "" && api.client.auth.Password != "" {
		restyClient.SetBasicAuth(api.client.auth.Username, api.client.auth.Password)
	}

	resp, err := restyClient.
		SetHeader("Accept", "application/json").
		SetResult(&results).
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

	return results, nil
}
