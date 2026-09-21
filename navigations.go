package gowprest

import (
	"encoding/json"
	"strconv"
	"strings"
)

// Navigation represents a WordPress navigation post (wp_navigation).
type Navigation struct {
	ID          int            `json:"id,omitempty"`
	Date        *Date          `json:"date,omitempty"`
	DateGMT     *Date          `json:"date_gmt,omitempty"`
	GUID        *Object        `json:"guid,omitempty"`
	Link        string         `json:"link,omitempty"`
	Modified    *Date          `json:"modified,omitempty"`
	ModifiedGMT *Date          `json:"modified_gmt,omitempty"`
	Slug        string         `json:"slug,omitempty"`
	Status      PostStatus     `json:"status,omitempty"`
	Type        string         `json:"type,omitempty"`
	Password    string         `json:"password,omitempty"`
	Title       *Object        `json:"title,omitempty"`
	Content     *Object        `json:"content,omitempty"`
	Template    string         `json:"template,omitempty"`
	Links       map[string]any `json:"_links,omitempty"`
}

// TitleString returns the title string from either Raw or Rendered.
func (n *Navigation) TitleString() string {
	if n.Title == nil {
		return ""
	}
	if n.Title.Raw != "" {
		return n.Title.Raw
	}
	return n.Title.Rendered
}

// ContentString returns the content string from either Raw or Rendered.
func (n *Navigation) ContentString() string {
	if n.Content == nil {
		return ""
	}
	if n.Content.Raw != "" {
		return n.Content.Raw
	}
	return n.Content.Rendered
}

// Navigations provides access to navigation related APIs (/wp/v2/navigation).
type Navigations struct {
	client *RestClient
}

// Navigations returns a Navigations service instance.
func (c *RestClient) Navigations() *Navigations {
	return &Navigations{client: c}
}

// Navigation is an ergonomic convenience alias for Navigations.
func (c *RestClient) Navigation() *Navigations {
	return c.Navigations()
}

// List returns a ListNavigations builder.
func (api *Navigations) List() *ListNavigations {
	return &ListNavigations{
		endpoint:  "/wp/v2/navigation",
		client:    api.client,
		arguments: make(map[string]string),
	}
}

// Retrieve returns a RetrieveNavigation builder for a specific navigation ID.
func (api *Navigations) Retrieve(id int) *RetrieveNavigation {
	return &RetrieveNavigation{
		endpoint:  "/wp/v2/navigation/" + strconv.Itoa(id),
		client:    api.client,
		arguments: make(map[string]string),
	}
}

// Create returns a CreateNavigation builder to create a new navigation post.
func (api *Navigations) Create() *CreateNavigation {
	return &CreateNavigation{
		endpoint: "/wp/v2/navigation",
		client:   api.client,
		body:     make(map[string]any),
	}
}

// Update returns an UpdateNavigation builder for the given navigation ID.
func (api *Navigations) Update(id int) *UpdateNavigation {
	return &UpdateNavigation{
		endpoint: "/wp/v2/navigation/" + strconv.Itoa(id),
		client:   api.client,
		body:     make(map[string]any),
	}
}

// Delete returns a DeleteNavigation builder for the given navigation ID.
func (api *Navigations) Delete(id int) *DeleteNavigation {
	return &DeleteNavigation{
		endpoint: "/wp/v2/navigation/" + strconv.Itoa(id),
		client:   api.client,
		force:    true,
	}
}

// ListNavigations handles querying navigation posts.
type ListNavigations struct {
	endpoint  string
	client    *RestClient
	arguments map[string]string
}

// Context sets the context parameter (e.g. "view", "edit", "embed").
func (api *ListNavigations) Context(ctx string) *ListNavigations {
	api.arguments["context"] = ctx
	return api
}

// ContextView sets the context parameter to "view".
func (api *ListNavigations) ContextView() *ListNavigations {
	return api.Context("view")
}

// ContextEdit sets the context parameter to "edit".
func (api *ListNavigations) ContextEdit() *ListNavigations {
	return api.Context("edit")
}

// ContextEmbed sets the context parameter to "embed".
func (api *ListNavigations) ContextEmbed() *ListNavigations {
	return api.Context("embed")
}

// Page sets the page number of the collection.
func (api *ListNavigations) Page(page int) *ListNavigations {
	api.arguments["page"] = strconv.Itoa(page)
	return api
}

// PerPage sets the number of items per page.
func (api *ListNavigations) PerPage(perPage int) *ListNavigations {
	api.arguments["per_page"] = strconv.Itoa(perPage)
	return api
}

// Search limits the response to items matching a search string.
func (api *ListNavigations) Search(query string) *ListNavigations {
	api.arguments["search"] = query
	return api
}

// After limits response to posts published after a given ISO8601 date.
func (api *ListNavigations) After(date string) *ListNavigations {
	api.arguments["after"] = date
	return api
}

// ModifiedAfter limits response to posts modified after a given ISO8601 date.
func (api *ListNavigations) ModifiedAfter(date string) *ListNavigations {
	api.arguments["modified_after"] = date
	return api
}

// Before limits response to posts published before a given ISO8601 date.
func (api *ListNavigations) Before(date string) *ListNavigations {
	api.arguments["before"] = date
	return api
}

// ModifiedBefore limits response to posts modified before a given ISO8601 date.
func (api *ListNavigations) ModifiedBefore(date string) *ListNavigations {
	api.arguments["modified_before"] = date
	return api
}

// Exclude ensures the response excludes specific IDs.
func (api *ListNavigations) Exclude(ids ...int) *ListNavigations {
	strIDs := make([]string, len(ids))
	for i, id := range ids {
		strIDs[i] = strconv.Itoa(id)
	}
	api.arguments["exclude"] = strings.Join(strIDs, ",")
	return api
}

// Include limits the response to specific IDs.
func (api *ListNavigations) Include(ids ...int) *ListNavigations {
	strIDs := make([]string, len(ids))
	for i, id := range ids {
		strIDs[i] = strconv.Itoa(id)
	}
	api.arguments["include"] = strings.Join(strIDs, ",")
	return api
}

// Offset offsets the result set by a specific number of items.
func (api *ListNavigations) Offset(offset int) *ListNavigations {
	api.arguments["offset"] = strconv.Itoa(offset)
	return api
}

// Order sets the sort order ("asc" or "desc").
func (api *ListNavigations) Order(order string) *ListNavigations {
	api.arguments["order"] = order
	return api
}

// OrderAsc sets the sort order to ascending.
func (api *ListNavigations) OrderAsc() *ListNavigations {
	return api.Order("asc")
}

// OrderDesc sets the sort order to descending.
func (api *ListNavigations) OrderDesc() *ListNavigations {
	return api.Order("desc")
}

// OrderBy sorts the collection by a specific field.
func (api *ListNavigations) OrderBy(field string) *ListNavigations {
	api.arguments["orderby"] = field
	return api
}

// SearchColumns limits the columns searched.
func (api *ListNavigations) SearchColumns(columns ...string) *ListNavigations {
	api.arguments["search_columns"] = strings.Join(columns, ",")
	return api
}

// Slug limits the response to posts with specific slugs.
func (api *ListNavigations) Slug(slugs ...string) *ListNavigations {
	api.arguments["slug"] = strings.Join(slugs, ",")
	return api
}

// Status limits the response to posts assigned specific statuses.
func (api *ListNavigations) Status(statuses ...PostStatus) *ListNavigations {
	strStatuses := make([]string, len(statuses))
	for i, s := range statuses {
		strStatuses[i] = string(s)
	}
	api.arguments["status"] = strings.Join(strStatuses, ",")
	return api
}

// Fields limits the response to specific fields.
func (api *ListNavigations) Fields(fields ...string) *ListNavigations {
	api.arguments["_fields"] = strings.Join(fields, ",")
	return api
}

// Do executes the request and returns the list of navigation posts.
func (api *ListNavigations) Do() (navigations []Navigation, err error) {
	endpoint := api.client.endpoint + api.endpoint

	restyClient := api.client.httpClient.R()
	if api.client.auth.Username != "" && api.client.auth.Password != "" {
		restyClient = restyClient.SetBasicAuth(api.client.auth.Username, api.client.auth.Password)
	}

	resp, err := restyClient.
		SetHeader("Accept", "application/json").
		SetResult(&navigations).
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

	return navigations, nil
}

// RetrieveNavigation handles retrieving a single navigation post.
type RetrieveNavigation struct {
	endpoint  string
	client    *RestClient
	arguments map[string]string
}

// Context sets the context parameter (e.g. "view", "edit", "embed").
func (api *RetrieveNavigation) Context(ctx string) *RetrieveNavigation {
	api.arguments["context"] = ctx
	return api
}

// ContextView sets the context parameter to "view".
func (api *RetrieveNavigation) ContextView() *RetrieveNavigation {
	return api.Context("view")
}

// ContextEdit sets the context parameter to "edit".
func (api *RetrieveNavigation) ContextEdit() *RetrieveNavigation {
	return api.Context("edit")
}

// ContextEmbed sets the context parameter to "embed".
func (api *RetrieveNavigation) ContextEmbed() *RetrieveNavigation {
	return api.Context("embed")
}

// Password sets the password if password protected.
func (api *RetrieveNavigation) Password(password string) *RetrieveNavigation {
	api.arguments["password"] = password
	return api
}

// Fields limits the response to specific fields.
func (api *RetrieveNavigation) Fields(fields ...string) *RetrieveNavigation {
	api.arguments["_fields"] = strings.Join(fields, ",")
	return api
}

// Do executes the retrieve request.
func (api *RetrieveNavigation) Do() (nav *Navigation, err error) {
	endpoint := api.client.endpoint + api.endpoint

	restyClient := api.client.httpClient.R()
	if api.client.auth.Username != "" && api.client.auth.Password != "" {
		restyClient = restyClient.SetBasicAuth(api.client.auth.Username, api.client.auth.Password)
	}

	resp, err := restyClient.
		SetHeader("Accept", "application/json").
		SetResult(&nav).
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

	return nav, nil
}

// CreateNavigation handles creating a new navigation post.
type CreateNavigation struct {
	endpoint string
	client   *RestClient
	body     map[string]any
}

func (api *CreateNavigation) Title(title string) *CreateNavigation {
	api.body["title"] = title
	return api
}

func (api *CreateNavigation) Content(content string) *CreateNavigation {
	api.body["content"] = content
	return api
}

func (api *CreateNavigation) Slug(slug string) *CreateNavigation {
	api.body["slug"] = slug
	return api
}

func (api *CreateNavigation) Status(status PostStatus) *CreateNavigation {
	api.body["status"] = string(status)
	return api
}

func (api *CreateNavigation) StatusPublish() *CreateNavigation {
	return api.Status(StatusPublished)
}

func (api *CreateNavigation) StatusDraft() *CreateNavigation {
	return api.Status(StatusDraft)
}

func (api *CreateNavigation) StatusPrivate() *CreateNavigation {
	return api.Status(StatusPrivate)
}

func (api *CreateNavigation) StatusPending() *CreateNavigation {
	return api.Status(StatusPending)
}

func (api *CreateNavigation) Password(password string) *CreateNavigation {
	api.body["password"] = password
	return api
}

func (api *CreateNavigation) Template(template string) *CreateNavigation {
	api.body["template"] = template
	return api
}

func (api *CreateNavigation) Date(date string) *CreateNavigation {
	api.body["date"] = date
	return api
}

func (api *CreateNavigation) DateGMT(dateGMT string) *CreateNavigation {
	api.body["date_gmt"] = dateGMT
	return api
}

func (api *CreateNavigation) Do() (nav *Navigation, err error) {
	endpoint := api.client.endpoint + api.endpoint

	restyClient := api.client.httpClient.R()
	if api.client.auth.Username != "" && api.client.auth.Password != "" {
		restyClient = restyClient.SetBasicAuth(api.client.auth.Username, api.client.auth.Password)
	}

	resp, err := restyClient.
		SetHeader("Content-Type", "application/json").
		SetHeader("Accept", "application/json").
		SetBody(api.body).
		SetResult(&nav).
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

	return nav, nil
}

// UpdateNavigation handles modifying an existing navigation post.
type UpdateNavigation struct {
	endpoint string
	client   *RestClient
	body     map[string]any
}

func (api *UpdateNavigation) Title(title string) *UpdateNavigation {
	api.body["title"] = title
	return api
}

func (api *UpdateNavigation) Content(content string) *UpdateNavigation {
	api.body["content"] = content
	return api
}

func (api *UpdateNavigation) Slug(slug string) *UpdateNavigation {
	api.body["slug"] = slug
	return api
}

func (api *UpdateNavigation) Status(status PostStatus) *UpdateNavigation {
	api.body["status"] = string(status)
	return api
}

func (api *UpdateNavigation) StatusPublish() *UpdateNavigation {
	return api.Status(StatusPublished)
}

func (api *UpdateNavigation) StatusDraft() *UpdateNavigation {
	return api.Status(StatusDraft)
}

func (api *UpdateNavigation) StatusPrivate() *UpdateNavigation {
	return api.Status(StatusPrivate)
}

func (api *UpdateNavigation) StatusPending() *UpdateNavigation {
	return api.Status(StatusPending)
}

func (api *UpdateNavigation) Password(password string) *UpdateNavigation {
	api.body["password"] = password
	return api
}

func (api *UpdateNavigation) Template(template string) *UpdateNavigation {
	api.body["template"] = template
	return api
}

func (api *UpdateNavigation) Do() (nav *Navigation, err error) {
	endpoint := api.client.endpoint + api.endpoint

	restyClient := api.client.httpClient.R()
	if api.client.auth.Username != "" && api.client.auth.Password != "" {
		restyClient = restyClient.SetBasicAuth(api.client.auth.Username, api.client.auth.Password)
	}

	resp, err := restyClient.
		SetHeader("Content-Type", "application/json").
		SetHeader("Accept", "application/json").
		SetBody(api.body).
		SetResult(&nav).
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

	return nav, nil
}

type deleteNavigationEnvelope struct {
	Navigation
	Deleted  bool        `json:"deleted"`
	Previous *Navigation `json:"previous"`
}

// DeleteNavigation handles deleting a navigation post.
type DeleteNavigation struct {
	endpoint string
	client   *RestClient
	force    bool
}

func (api *DeleteNavigation) Force(force ...bool) *DeleteNavigation {
	if len(force) > 0 {
		api.force = force[0]
	} else {
		api.force = true
	}
	return api
}

func (api *DeleteNavigation) Do() (nav *Navigation, err error) {
	endpoint := api.client.endpoint + api.endpoint

	restyClient := api.client.httpClient.R()
	if api.client.auth.Username != "" && api.client.auth.Password != "" {
		restyClient = restyClient.SetBasicAuth(api.client.auth.Username, api.client.auth.Password)
	}

	if api.force {
		restyClient = restyClient.SetQueryParam("force", "true")
	}

	var env deleteNavigationEnvelope
	resp, err := restyClient.
		SetHeader("Accept", "application/json").
		SetResult(&env).
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

	return &env.Navigation, nil
}
