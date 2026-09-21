package gowprest

import (
	"encoding/json"
	"strconv"
	"strings"
)

// NavigationRevision represents a WordPress navigation (wp_navigation) revision.
type NavigationRevision struct {
	Author      int            `json:"author,omitempty"`
	Date        *Date          `json:"date,omitempty"`
	DateGMT     *Date          `json:"date_gmt,omitempty"`
	GUID        *Object        `json:"guid,omitempty"`
	ID          int            `json:"id,omitempty"`
	Modified    *Date          `json:"modified,omitempty"`
	ModifiedGMT *Date          `json:"modified_gmt,omitempty"`
	Parent      int            `json:"parent,omitempty"`
	Slug        string         `json:"slug,omitempty"`
	Title       *Object        `json:"title,omitempty"`
	Content     *Object        `json:"content,omitempty"`
	Links       map[string]any `json:"_links,omitempty"`
	Embedded    map[string]any `json:"_embedded,omitempty"`
}

// TitleString returns the revision's title string from either Raw or Rendered.
func (r *NavigationRevision) TitleString() string {
	if r.Title == nil {
		return ""
	}
	if r.Title.Raw != "" {
		return r.Title.Raw
	}
	return r.Title.Rendered
}

// ContentString returns the revision's content string from either Raw or Rendered.
func (r *NavigationRevision) ContentString() string {
	if r.Content == nil {
		return ""
	}
	if r.Content.Raw != "" {
		return r.Content.Raw
	}
	return r.Content.Rendered
}

// NavigationRevisions anchors revision-related operations for a specific Navigation parent.
type NavigationRevisions struct {
	client   *RestClient
	parentID int
}

// NavigationRevisions returns a NavigationRevisions service instance for the given parent ID.
func (c *RestClient) NavigationRevisions(parentID int) *NavigationRevisions {
	return &NavigationRevisions{
		client:   c,
		parentID: parentID,
	}
}

// Revisions returns a NavigationRevisions service instance for the given parent navigation ID.
func (api *Navigations) Revisions(parentID int) *NavigationRevisions {
	return &NavigationRevisions{
		client:   api.client,
		parentID: parentID,
	}
}

// List returns a ListNavigationRevisions builder to list revisions for the parent navigation.
func (api *NavigationRevisions) List() *ListNavigationRevisions {
	return &ListNavigationRevisions{
		endpoint:  "/wp/v2/navigation/" + strconv.Itoa(api.parentID) + "/revisions",
		client:    api.client,
		arguments: make(map[string]string),
	}
}

// Retrieve returns a RetrieveNavigationRevision builder to get a specific revision by ID.
func (api *NavigationRevisions) Retrieve(id int) *RetrieveNavigationRevision {
	return &RetrieveNavigationRevision{
		endpoint:  "/wp/v2/navigation/" + strconv.Itoa(api.parentID) + "/revisions/" + strconv.Itoa(id),
		client:    api.client,
		arguments: make(map[string]string),
	}
}

// Delete returns a DeleteNavigationRevision builder to delete a revision.
func (api *NavigationRevisions) Delete(id int) *DeleteNavigationRevision {
	return &DeleteNavigationRevision{
		endpoint: "/wp/v2/navigation/" + strconv.Itoa(api.parentID) + "/revisions/" + strconv.Itoa(id),
		client:   api.client,
		force:    true,
	}
}

// ListNavigationRevisions handles querying navigation revisions.
type ListNavigationRevisions struct {
	endpoint  string
	client    *RestClient
	arguments map[string]string
}

// Context sets the context parameter (e.g. "view", "edit", "embed").
func (api *ListNavigationRevisions) Context(ctx string) *ListNavigationRevisions {
	api.arguments["context"] = ctx
	return api
}

// ContextView sets the context parameter to "view".
func (api *ListNavigationRevisions) ContextView() *ListNavigationRevisions {
	return api.Context("view")
}

// ContextEdit sets the context parameter to "edit".
func (api *ListNavigationRevisions) ContextEdit() *ListNavigationRevisions {
	return api.Context("edit")
}

// ContextEmbed sets the context parameter to "embed".
func (api *ListNavigationRevisions) ContextEmbed() *ListNavigationRevisions {
	return api.Context("embed")
}

// Page sets the page number of the collection.
func (api *ListNavigationRevisions) Page(page int) *ListNavigationRevisions {
	api.arguments["page"] = strconv.Itoa(page)
	return api
}

// PerPage sets the number of items per page.
func (api *ListNavigationRevisions) PerPage(perPage int) *ListNavigationRevisions {
	api.arguments["per_page"] = strconv.Itoa(perPage)
	return api
}

// Search limits the response to revisions matching a search string.
func (api *ListNavigationRevisions) Search(query string) *ListNavigationRevisions {
	api.arguments["search"] = query
	return api
}

// Exclude ensures the response excludes specific revision IDs.
func (api *ListNavigationRevisions) Exclude(ids ...int) *ListNavigationRevisions {
	strIDs := make([]string, len(ids))
	for i, id := range ids {
		strIDs[i] = strconv.Itoa(id)
	}
	api.arguments["exclude"] = strings.Join(strIDs, ",")
	return api
}

// Include limits the response to specific revision IDs.
func (api *ListNavigationRevisions) Include(ids ...int) *ListNavigationRevisions {
	strIDs := make([]string, len(ids))
	for i, id := range ids {
		strIDs[i] = strconv.Itoa(id)
	}
	api.arguments["include"] = strings.Join(strIDs, ",")
	return api
}

// Offset offsets the result set by a specific number of items.
func (api *ListNavigationRevisions) Offset(offset int) *ListNavigationRevisions {
	api.arguments["offset"] = strconv.Itoa(offset)
	return api
}

// Order sets the sort order ("asc" or "desc").
func (api *ListNavigationRevisions) Order(order string) *ListNavigationRevisions {
	api.arguments["order"] = order
	return api
}

// OrderAsc sets the sort order to ascending.
func (api *ListNavigationRevisions) OrderAsc() *ListNavigationRevisions {
	return api.Order("asc")
}

// OrderDesc sets the sort order to descending.
func (api *ListNavigationRevisions) OrderDesc() *ListNavigationRevisions {
	return api.Order("desc")
}

// OrderBy sorts the collection by a specific field.
func (api *ListNavigationRevisions) OrderBy(field string) *ListNavigationRevisions {
	api.arguments["orderby"] = field
	return api
}

// Fields limits the response to specific fields.
func (api *ListNavigationRevisions) Fields(fields ...string) *ListNavigationRevisions {
	api.arguments["_fields"] = strings.Join(fields, ",")
	return api
}

// Do executes the request and returns the collection of navigation revisions.
func (api *ListNavigationRevisions) Do() (revisions []NavigationRevision, err error) {
	endpoint := api.client.endpoint + api.endpoint

	restyClient := api.client.httpClient.R()
	if api.client.auth.Username != "" && api.client.auth.Password != "" {
		restyClient = restyClient.SetBasicAuth(api.client.auth.Username, api.client.auth.Password)
	}

	resp, err := restyClient.
		SetHeader("Accept", "application/json").
		SetResult(&revisions).
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

	return revisions, nil
}

// RetrieveNavigationRevision handles retrieving a single navigation revision.
type RetrieveNavigationRevision struct {
	endpoint  string
	client    *RestClient
	arguments map[string]string
}

// Context sets the context parameter (e.g. "view", "edit", "embed").
func (api *RetrieveNavigationRevision) Context(ctx string) *RetrieveNavigationRevision {
	api.arguments["context"] = ctx
	return api
}

// ContextView sets the context parameter to "view".
func (api *RetrieveNavigationRevision) ContextView() *RetrieveNavigationRevision {
	return api.Context("view")
}

// ContextEdit sets the context parameter to "edit".
func (api *RetrieveNavigationRevision) ContextEdit() *RetrieveNavigationRevision {
	return api.Context("edit")
}

// ContextEmbed sets the context parameter to "embed".
func (api *RetrieveNavigationRevision) ContextEmbed() *RetrieveNavigationRevision {
	return api.Context("embed")
}

// Fields limits the response to specific fields.
func (api *RetrieveNavigationRevision) Fields(fields ...string) *RetrieveNavigationRevision {
	api.arguments["_fields"] = strings.Join(fields, ",")
	return api
}

// Do executes the request and returns the navigation revision.
func (api *RetrieveNavigationRevision) Do() (revision *NavigationRevision, err error) {
	endpoint := api.client.endpoint + api.endpoint

	restyClient := api.client.httpClient.R()
	if api.client.auth.Username != "" && api.client.auth.Password != "" {
		restyClient = restyClient.SetBasicAuth(api.client.auth.Username, api.client.auth.Password)
	}

	resp, err := restyClient.
		SetHeader("Accept", "application/json").
		SetResult(&revision).
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

	return revision, nil
}

type deleteNavigationRevisionEnvelope struct {
	NavigationRevision
	Deleted  bool                `json:"deleted"`
	Previous *NavigationRevision `json:"previous"`
}

// DeleteNavigationRevision handles deleting a navigation revision.
type DeleteNavigationRevision struct {
	endpoint string
	client   *RestClient
	force    bool
}

// Force sets whether to force delete the revision. Revisions do not support trashing.
func (api *DeleteNavigationRevision) Force(force bool) *DeleteNavigationRevision {
	api.force = force
	return api
}

// Do executes the delete request and returns the deleted navigation revision.
func (api *DeleteNavigationRevision) Do() (revision *NavigationRevision, err error) {
	endpoint := api.client.endpoint + api.endpoint

	restyClient := api.client.httpClient.R()
	if api.client.auth.Username != "" && api.client.auth.Password != "" {
		restyClient = restyClient.SetBasicAuth(api.client.auth.Username, api.client.auth.Password)
	}

	var env deleteNavigationRevisionEnvelope
	resp, err := restyClient.
		SetHeader("Accept", "application/json").
		SetResult(&env).
		SetQueryParam("force", strconv.FormatBool(api.force)).
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

	if env.Deleted && env.Previous != nil && env.Previous.ID != 0 {
		return env.Previous, nil
	}

	return &env.NavigationRevision, nil
}
