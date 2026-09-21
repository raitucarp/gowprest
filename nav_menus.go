package gowprest

import (
	"encoding/json"
	"strconv"
	"strings"
)

// NavMenu represents a WordPress navigation menu term (nav_menu taxonomy).
type NavMenu struct {
	ID          int            `json:"id,omitempty"`
	Description string         `json:"description,omitempty"`
	Name        string         `json:"name,omitempty"`
	Slug        string         `json:"slug,omitempty"`
	Meta        any            `json:"meta,omitempty"`
	Locations   []string       `json:"locations,omitempty"`
	AutoAdd     bool           `json:"auto_add,omitempty"`
	Links       map[string]any `json:"_links,omitempty"`
}

// NavMenus provides access to navigation menus API (/wp/v2/menus).
type NavMenus struct {
	client *RestClient
}

// NavMenus returns a NavMenus service instance.
func (c *RestClient) NavMenus() *NavMenus {
	return &NavMenus{client: c}
}

// NavMenu is an ergonomic alias for NavMenus.
func (c *RestClient) NavMenu() *NavMenus {
	return c.NavMenus()
}

// Menus is an alias for NavMenus matching the /wp/v2/menus endpoint.
func (c *RestClient) Menus() *NavMenus {
	return c.NavMenus()
}

// Menu is an ergonomic alias for Menus.
func (c *RestClient) Menu() *NavMenus {
	return c.NavMenus()
}

// List returns a ListNavMenus builder.
func (api *NavMenus) List() *ListNavMenus {
	return &ListNavMenus{
		endpoint:  "/wp/v2/menus",
		client:    api.client,
		arguments: make(map[string]string),
	}
}

// Retrieve returns a RetrieveNavMenu builder for a specific menu ID.
func (api *NavMenus) Retrieve(id int) *RetrieveNavMenu {
	return &RetrieveNavMenu{
		endpoint:  "/wp/v2/menus/" + strconv.Itoa(id),
		client:    api.client,
		arguments: make(map[string]string),
	}
}

// Create returns a CreateNavMenu builder.
func (api *NavMenus) Create() *CreateNavMenu {
	return &CreateNavMenu{
		endpoint: "/wp/v2/menus",
		client:   api.client,
		body:     make(map[string]any),
	}
}

// Update returns an UpdateNavMenu builder for a specific menu ID.
func (api *NavMenus) Update(id int) *UpdateNavMenu {
	return &UpdateNavMenu{
		endpoint: "/wp/v2/menus/" + strconv.Itoa(id),
		client:   api.client,
		body:     make(map[string]any),
	}
}

// Delete returns a DeleteNavMenu builder for a specific menu ID.
func (api *NavMenus) Delete(id int) *DeleteNavMenu {
	return &DeleteNavMenu{
		endpoint: "/wp/v2/menus/" + strconv.Itoa(id),
		client:   api.client,
		force:    true,
	}
}

// ListNavMenus handles querying navigation menus.
type ListNavMenus struct {
	endpoint  string
	client    *RestClient
	arguments map[string]string
}

// Context sets the context parameter (e.g. "view", "edit", "embed").
func (api *ListNavMenus) Context(ctx string) *ListNavMenus {
	api.arguments["context"] = ctx
	return api
}

// ContextView sets the context parameter to "view".
func (api *ListNavMenus) ContextView() *ListNavMenus {
	return api.Context("view")
}

// ContextEdit sets the context parameter to "edit".
func (api *ListNavMenus) ContextEdit() *ListNavMenus {
	return api.Context("edit")
}

// ContextEmbed sets the context parameter to "embed".
func (api *ListNavMenus) ContextEmbed() *ListNavMenus {
	return api.Context("embed")
}

// Page sets the page number of the collection.
func (api *ListNavMenus) Page(page int) *ListNavMenus {
	api.arguments["page"] = strconv.Itoa(page)
	return api
}

// PerPage sets the number of items per page.
func (api *ListNavMenus) PerPage(perPage int) *ListNavMenus {
	api.arguments["per_page"] = strconv.Itoa(perPage)
	return api
}

// Search limits results to those matching a string.
func (api *ListNavMenus) Search(query string) *ListNavMenus {
	api.arguments["search"] = query
	return api
}

// Exclude ensures the result set excludes specific IDs.
func (api *ListNavMenus) Exclude(ids ...int) *ListNavMenus {
	strIDs := make([]string, len(ids))
	for i, id := range ids {
		strIDs[i] = strconv.Itoa(id)
	}
	api.arguments["exclude"] = strings.Join(strIDs, ",")
	return api
}

// Include limits the result set to specific IDs.
func (api *ListNavMenus) Include(ids ...int) *ListNavMenus {
	strIDs := make([]string, len(ids))
	for i, id := range ids {
		strIDs[i] = strconv.Itoa(id)
	}
	api.arguments["include"] = strings.Join(strIDs, ",")
	return api
}

// Offset offsets the result set by a specific number of items.
func (api *ListNavMenus) Offset(offset int) *ListNavMenus {
	api.arguments["offset"] = strconv.Itoa(offset)
	return api
}

// Order sets the sort order ("asc" or "desc").
func (api *ListNavMenus) Order(order string) *ListNavMenus {
	api.arguments["order"] = order
	return api
}

// OrderAsc sets the sort order to ascending.
func (api *ListNavMenus) OrderAsc() *ListNavMenus {
	return api.Order("asc")
}

// OrderDesc sets the sort order to descending.
func (api *ListNavMenus) OrderDesc() *ListNavMenus {
	return api.Order("desc")
}

// OrderBy sorts the collection by a term attribute.
func (api *ListNavMenus) OrderBy(field string) *ListNavMenus {
	api.arguments["orderby"] = field
	return api
}

// HideEmpty sets whether to hide terms not assigned to any posts.
func (api *ListNavMenus) HideEmpty(hideEmpty ...bool) *ListNavMenus {
	if len(hideEmpty) > 0 {
		api.arguments["hide_empty"] = strconv.FormatBool(hideEmpty[0])
	} else {
		api.arguments["hide_empty"] = "true"
	}
	return api
}

// Post limits results to terms assigned to a specific post ID.
func (api *ListNavMenus) Post(postID int) *ListNavMenus {
	api.arguments["post"] = strconv.Itoa(postID)
	return api
}

// Slug limits results to terms with one or more specific slugs.
func (api *ListNavMenus) Slug(slugs ...string) *ListNavMenus {
	api.arguments["slug"] = strings.Join(slugs, ",")
	return api
}

// Fields limits the response to specific fields.
func (api *ListNavMenus) Fields(fields ...string) *ListNavMenus {
	api.arguments["_fields"] = strings.Join(fields, ",")
	return api
}

// Do executes the list request.
func (api *ListNavMenus) Do() (menus []NavMenu, err error) {
	endpoint := api.client.endpoint + api.endpoint

	restyClient := api.client.httpClient.R()
	if api.client.auth.Username != "" && api.client.auth.Password != "" {
		restyClient = restyClient.SetBasicAuth(api.client.auth.Username, api.client.auth.Password)
	}

	resp, err := restyClient.
		SetHeader("Accept", "application/json").
		SetResult(&menus).
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

	return menus, nil
}

// RetrieveNavMenu handles retrieving a single navigation menu.
type RetrieveNavMenu struct {
	endpoint  string
	client    *RestClient
	arguments map[string]string
}

// Context sets the context parameter (e.g. "view", "edit", "embed").
func (api *RetrieveNavMenu) Context(ctx string) *RetrieveNavMenu {
	api.arguments["context"] = ctx
	return api
}

// ContextView sets the context parameter to "view".
func (api *RetrieveNavMenu) ContextView() *RetrieveNavMenu {
	return api.Context("view")
}

// ContextEdit sets the context parameter to "edit".
func (api *RetrieveNavMenu) ContextEdit() *RetrieveNavMenu {
	return api.Context("edit")
}

// ContextEmbed sets the context parameter to "embed".
func (api *RetrieveNavMenu) ContextEmbed() *RetrieveNavMenu {
	return api.Context("embed")
}

// Fields limits the response to specific fields.
func (api *RetrieveNavMenu) Fields(fields ...string) *RetrieveNavMenu {
	api.arguments["_fields"] = strings.Join(fields, ",")
	return api
}

// Do executes the retrieve request.
func (api *RetrieveNavMenu) Do() (menu *NavMenu, err error) {
	endpoint := api.client.endpoint + api.endpoint

	restyClient := api.client.httpClient.R()
	if api.client.auth.Username != "" && api.client.auth.Password != "" {
		restyClient = restyClient.SetBasicAuth(api.client.auth.Username, api.client.auth.Password)
	}

	resp, err := restyClient.
		SetHeader("Accept", "application/json").
		SetResult(&menu).
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

	return menu, nil
}

// CreateNavMenu handles creating a new navigation menu.
type CreateNavMenu struct {
	endpoint string
	client   *RestClient
	body     map[string]any
}

// Name sets the HTML title for the menu (required).
func (api *CreateNavMenu) Name(name string) *CreateNavMenu {
	api.body["name"] = name
	return api
}

// Description sets the HTML description of the menu.
func (api *CreateNavMenu) Description(description string) *CreateNavMenu {
	api.body["description"] = description
	return api
}

// Slug sets the alphanumeric identifier for the menu unique to its type.
func (api *CreateNavMenu) Slug(slug string) *CreateNavMenu {
	api.body["slug"] = slug
	return api
}

// Locations sets the locations assigned to the menu.
func (api *CreateNavMenu) Locations(locations ...string) *CreateNavMenu {
	api.body["locations"] = locations
	return api
}

// AutoAdd sets whether to automatically add top level pages to this menu.
func (api *CreateNavMenu) AutoAdd(autoAdd bool) *CreateNavMenu {
	api.body["auto_add"] = autoAdd
	return api
}

// Meta sets meta fields for the menu.
func (api *CreateNavMenu) Meta(meta any) *CreateNavMenu {
	api.body["meta"] = meta
	return api
}

// Do executes the create request.
func (api *CreateNavMenu) Do() (menu *NavMenu, err error) {
	endpoint := api.client.endpoint + api.endpoint

	restyClient := api.client.httpClient.R()
	if api.client.auth.Username != "" && api.client.auth.Password != "" {
		restyClient = restyClient.SetBasicAuth(api.client.auth.Username, api.client.auth.Password)
	}

	resp, err := restyClient.
		SetHeader("Content-Type", "application/json").
		SetHeader("Accept", "application/json").
		SetBody(api.body).
		SetResult(&menu).
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

	return menu, nil
}

// UpdateNavMenu handles modifying an existing navigation menu.
type UpdateNavMenu struct {
	endpoint string
	client   *RestClient
	body     map[string]any
}

// Name sets the HTML title for the menu.
func (api *UpdateNavMenu) Name(name string) *UpdateNavMenu {
	api.body["name"] = name
	return api
}

// Description sets the HTML description of the menu.
func (api *UpdateNavMenu) Description(description string) *UpdateNavMenu {
	api.body["description"] = description
	return api
}

// Slug sets the alphanumeric identifier for the menu.
func (api *UpdateNavMenu) Slug(slug string) *UpdateNavMenu {
	api.body["slug"] = slug
	return api
}

// Locations sets the locations assigned to the menu.
func (api *UpdateNavMenu) Locations(locations ...string) *UpdateNavMenu {
	api.body["locations"] = locations
	return api
}

// AutoAdd sets whether to automatically add top level pages to this menu.
func (api *UpdateNavMenu) AutoAdd(autoAdd bool) *UpdateNavMenu {
	api.body["auto_add"] = autoAdd
	return api
}

// Meta sets meta fields for the menu.
func (api *UpdateNavMenu) Meta(meta any) *UpdateNavMenu {
	api.body["meta"] = meta
	return api
}

// Do executes the update request.
func (api *UpdateNavMenu) Do() (menu *NavMenu, err error) {
	endpoint := api.client.endpoint + api.endpoint

	restyClient := api.client.httpClient.R()
	if api.client.auth.Username != "" && api.client.auth.Password != "" {
		restyClient = restyClient.SetBasicAuth(api.client.auth.Username, api.client.auth.Password)
	}

	resp, err := restyClient.
		SetHeader("Content-Type", "application/json").
		SetHeader("Accept", "application/json").
		SetBody(api.body).
		SetResult(&menu).
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

	return menu, nil
}

type deleteNavMenuEnvelope struct {
	NavMenu
	Deleted  bool     `json:"deleted"`
	Previous *NavMenu `json:"previous"`
}

// DeleteNavMenu handles deleting a navigation menu.
type DeleteNavMenu struct {
	endpoint string
	client   *RestClient
	force    bool
}

// Force sets whether to force deletion. Nav menu terms require force=true.
func (api *DeleteNavMenu) Force(force ...bool) *DeleteNavMenu {
	if len(force) > 0 {
		api.force = force[0]
	} else {
		api.force = true
	}
	return api
}

// Do executes the delete request.
func (api *DeleteNavMenu) Do() (menu *NavMenu, err error) {
	endpoint := api.client.endpoint + api.endpoint

	restyClient := api.client.httpClient.R()
	if api.client.auth.Username != "" && api.client.auth.Password != "" {
		restyClient = restyClient.SetBasicAuth(api.client.auth.Username, api.client.auth.Password)
	}

	if api.force {
		restyClient = restyClient.SetQueryParam("force", "true")
	}

	var env deleteNavMenuEnvelope
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

	return &env.NavMenu, nil
}
