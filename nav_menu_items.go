package gowprest

import (
	"encoding/json"
	"strconv"
	"strings"
)

// NavMenuItem represents a WordPress navigation menu item (nav_menu_item post type).
type NavMenuItem struct {
	Title       *Object        `json:"title,omitempty"`
	ID          int            `json:"id,omitempty"`
	TypeLabel   string         `json:"type_label,omitempty"`
	Type        string         `json:"type,omitempty"`
	Status      string         `json:"status,omitempty"`
	Parent      int            `json:"parent,omitempty"`
	AttrTitle   string         `json:"attr_title,omitempty"`
	Classes     []string       `json:"classes,omitempty"`
	Description string         `json:"description,omitempty"`
	MenuOrder   int            `json:"menu_order,omitempty"`
	Object      string         `json:"object,omitempty"`
	ObjectID    int            `json:"object_id,omitempty"`
	Target      string         `json:"target,omitempty"`
	URL         string         `json:"url,omitempty"`
	XFN         []string       `json:"xfn,omitempty"`
	Invalid     bool           `json:"invalid,omitempty"`
	Menus       int            `json:"menus,omitempty"`
	Meta        any            `json:"meta,omitempty"`
	Links       map[string]any `json:"_links,omitempty"`
}

// TitleString returns the title string from either Raw or Rendered.
func (n *NavMenuItem) TitleString() string {
	if n.Title == nil {
		return ""
	}
	if n.Title.Raw != "" {
		return n.Title.Raw
	}
	return n.Title.Rendered
}

// NavMenuItems provides access to navigation menu items APIs (/wp/v2/menu-items).
type NavMenuItems struct {
	client *RestClient
}

// NavMenuItems returns a NavMenuItems service instance.
func (c *RestClient) NavMenuItems() *NavMenuItems {
	return &NavMenuItems{client: c}
}

// NavMenuItem is an ergonomic alias for NavMenuItems.
func (c *RestClient) NavMenuItem() *NavMenuItems {
	return c.NavMenuItems()
}

// MenuItems is an alias for NavMenuItems matching the /wp/v2/menu-items endpoint path.
func (c *RestClient) MenuItems() *NavMenuItems {
	return c.NavMenuItems()
}

// MenuItem is an ergonomic alias for MenuItems.
func (c *RestClient) MenuItem() *NavMenuItems {
	return c.NavMenuItems()
}

// List returns a ListNavMenuItems builder.
func (api *NavMenuItems) List() *ListNavMenuItems {
	return &ListNavMenuItems{
		endpoint:  "/wp/v2/menu-items",
		client:    api.client,
		arguments: make(map[string]string),
	}
}

// Retrieve returns a RetrieveNavMenuItem builder for a specific menu item ID.
func (api *NavMenuItems) Retrieve(id int) *RetrieveNavMenuItem {
	return &RetrieveNavMenuItem{
		endpoint:  "/wp/v2/menu-items/" + strconv.Itoa(id),
		client:    api.client,
		arguments: make(map[string]string),
	}
}

// Create returns a CreateNavMenuItem builder.
func (api *NavMenuItems) Create() *CreateNavMenuItem {
	return &CreateNavMenuItem{
		endpoint: "/wp/v2/menu-items",
		client:   api.client,
		body:     make(map[string]any),
	}
}

// Update returns an UpdateNavMenuItem builder for a specific menu item ID.
func (api *NavMenuItems) Update(id int) *UpdateNavMenuItem {
	return &UpdateNavMenuItem{
		endpoint: "/wp/v2/menu-items/" + strconv.Itoa(id),
		client:   api.client,
		body:     make(map[string]any),
	}
}

// Delete returns a DeleteNavMenuItem builder for a specific menu item ID.
func (api *NavMenuItems) Delete(id int) *DeleteNavMenuItem {
	return &DeleteNavMenuItem{
		endpoint: "/wp/v2/menu-items/" + strconv.Itoa(id),
		client:   api.client,
		force:    true,
	}
}

// ListNavMenuItems handles querying navigation menu items.
type ListNavMenuItems struct {
	endpoint  string
	client    *RestClient
	arguments map[string]string
}

// Context sets the context parameter (e.g. "view", "edit", "embed").
func (api *ListNavMenuItems) Context(ctx string) *ListNavMenuItems {
	api.arguments["context"] = ctx
	return api
}

// ContextView sets the context parameter to "view".
func (api *ListNavMenuItems) ContextView() *ListNavMenuItems {
	return api.Context("view")
}

// ContextEdit sets the context parameter to "edit".
func (api *ListNavMenuItems) ContextEdit() *ListNavMenuItems {
	return api.Context("edit")
}

// ContextEmbed sets the context parameter to "embed".
func (api *ListNavMenuItems) ContextEmbed() *ListNavMenuItems {
	return api.Context("embed")
}

// Page sets the page number of the collection.
func (api *ListNavMenuItems) Page(page int) *ListNavMenuItems {
	api.arguments["page"] = strconv.Itoa(page)
	return api
}

// PerPage sets the maximum number of items to return.
func (api *ListNavMenuItems) PerPage(perPage int) *ListNavMenuItems {
	api.arguments["per_page"] = strconv.Itoa(perPage)
	return api
}

// Search limits results to those matching a string.
func (api *ListNavMenuItems) Search(query string) *ListNavMenuItems {
	api.arguments["search"] = query
	return api
}

// After limits response to posts published after a given ISO8601 date.
func (api *ListNavMenuItems) After(date string) *ListNavMenuItems {
	api.arguments["after"] = date
	return api
}

// ModifiedAfter limits response to posts modified after a given ISO8601 date.
func (api *ListNavMenuItems) ModifiedAfter(date string) *ListNavMenuItems {
	api.arguments["modified_after"] = date
	return api
}

// Before limits response to posts published before a given ISO8601 date.
func (api *ListNavMenuItems) Before(date string) *ListNavMenuItems {
	api.arguments["before"] = date
	return api
}

// ModifiedBefore limits response to posts modified before a given ISO8601 date.
func (api *ListNavMenuItems) ModifiedBefore(date string) *ListNavMenuItems {
	api.arguments["modified_before"] = date
	return api
}

// Exclude ensures the result set excludes specific IDs.
func (api *ListNavMenuItems) Exclude(ids ...int) *ListNavMenuItems {
	strIDs := make([]string, len(ids))
	for i, id := range ids {
		strIDs[i] = strconv.Itoa(id)
	}
	api.arguments["exclude"] = strings.Join(strIDs, ",")
	return api
}

// Include limits the result set to specific IDs.
func (api *ListNavMenuItems) Include(ids ...int) *ListNavMenuItems {
	strIDs := make([]string, len(ids))
	for i, id := range ids {
		strIDs[i] = strconv.Itoa(id)
	}
	api.arguments["include"] = strings.Join(strIDs, ",")
	return api
}

// Offset offsets the result set by a specific number of items.
func (api *ListNavMenuItems) Offset(offset int) *ListNavMenuItems {
	api.arguments["offset"] = strconv.Itoa(offset)
	return api
}

// Order sets the sort order ("asc" or "desc").
func (api *ListNavMenuItems) Order(order string) *ListNavMenuItems {
	api.arguments["order"] = order
	return api
}

// OrderAsc sets the sort order to ascending.
func (api *ListNavMenuItems) OrderAsc() *ListNavMenuItems {
	return api.Order("asc")
}

// OrderDesc sets the sort order to descending.
func (api *ListNavMenuItems) OrderDesc() *ListNavMenuItems {
	return api.Order("desc")
}

// OrderBy sorts collection by object attribute.
func (api *ListNavMenuItems) OrderBy(field string) *ListNavMenuItems {
	api.arguments["orderby"] = field
	return api
}

// SearchColumns limits search to specific column names.
func (api *ListNavMenuItems) SearchColumns(columns ...string) *ListNavMenuItems {
	api.arguments["search_columns"] = strings.Join(columns, ",")
	return api
}

// Slug limits response to posts with specific slugs.
func (api *ListNavMenuItems) Slug(slugs ...string) *ListNavMenuItems {
	api.arguments["slug"] = strings.Join(slugs, ",")
	return api
}

// Status limits response to posts assigned one or more statuses.
func (api *ListNavMenuItems) Status(statuses ...string) *ListNavMenuItems {
	api.arguments["status"] = strings.Join(statuses, ",")
	return api
}

// TaxRelation limits result set based on relationship between multiple taxonomies (AND, OR).
func (api *ListNavMenuItems) TaxRelation(relation string) *ListNavMenuItems {
	api.arguments["tax_relation"] = relation
	return api
}

// Menus limits result set to items with specific terms assigned in menus taxonomy.
func (api *ListNavMenuItems) Menus(menuIDs ...int) *ListNavMenuItems {
	strIDs := make([]string, len(menuIDs))
	for i, id := range menuIDs {
		strIDs[i] = strconv.Itoa(id)
	}
	api.arguments["menus"] = strings.Join(strIDs, ",")
	return api
}

// MenusExclude limits result set to items except those with specific terms assigned in menus taxonomy.
func (api *ListNavMenuItems) MenusExclude(menuIDs ...int) *ListNavMenuItems {
	strIDs := make([]string, len(menuIDs))
	for i, id := range menuIDs {
		strIDs[i] = strconv.Itoa(id)
	}
	api.arguments["menus_exclude"] = strings.Join(strIDs, ",")
	return api
}

// MenuOrder limits result set to posts with a specific menu_order value.
func (api *ListNavMenuItems) MenuOrder(order int) *ListNavMenuItems {
	api.arguments["menu_order"] = strconv.Itoa(order)
	return api
}

// Fields limits the response to specific fields.
func (api *ListNavMenuItems) Fields(fields ...string) *ListNavMenuItems {
	api.arguments["_fields"] = strings.Join(fields, ",")
	return api
}

// Do executes the list request.
func (api *ListNavMenuItems) Do() (items []NavMenuItem, err error) {
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

// RetrieveNavMenuItem handles retrieving a single navigation menu item.
type RetrieveNavMenuItem struct {
	endpoint  string
	client    *RestClient
	arguments map[string]string
}

// Context sets the context parameter (e.g. "view", "edit", "embed").
func (api *RetrieveNavMenuItem) Context(ctx string) *RetrieveNavMenuItem {
	api.arguments["context"] = ctx
	return api
}

// ContextView sets the context parameter to "view".
func (api *RetrieveNavMenuItem) ContextView() *RetrieveNavMenuItem {
	return api.Context("view")
}

// ContextEdit sets the context parameter to "edit".
func (api *RetrieveNavMenuItem) ContextEdit() *RetrieveNavMenuItem {
	return api.Context("edit")
}

// ContextEmbed sets the context parameter to "embed".
func (api *RetrieveNavMenuItem) ContextEmbed() *RetrieveNavMenuItem {
	return api.Context("embed")
}

// Fields limits the response to specific fields.
func (api *RetrieveNavMenuItem) Fields(fields ...string) *RetrieveNavMenuItem {
	api.arguments["_fields"] = strings.Join(fields, ",")
	return api
}

// Do executes the retrieve request.
func (api *RetrieveNavMenuItem) Do() (item *NavMenuItem, err error) {
	endpoint := api.client.endpoint + api.endpoint

	restyClient := api.client.httpClient.R()
	if api.client.auth.Username != "" && api.client.auth.Password != "" {
		restyClient = restyClient.SetBasicAuth(api.client.auth.Username, api.client.auth.Password)
	}

	resp, err := restyClient.
		SetHeader("Accept", "application/json").
		SetResult(&item).
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

	return item, nil
}

// CreateNavMenuItem handles creating a new navigation menu item.
type CreateNavMenuItem struct {
	endpoint string
	client   *RestClient
	body     map[string]any
}

// Title sets the title for the object.
func (api *CreateNavMenuItem) Title(title string) *CreateNavMenuItem {
	api.body["title"] = title
	return api
}

// Type sets the family of objects originally represented (taxonomy, post_type, post_type_archive, custom).
func (api *CreateNavMenuItem) Type(itemType string) *CreateNavMenuItem {
	api.body["type"] = itemType
	return api
}

// Status sets the named status for the object (publish, future, draft, pending, private).
func (api *CreateNavMenuItem) Status(status string) *CreateNavMenuItem {
	api.body["status"] = status
	return api
}

// Parent sets the ID for the parent of the object.
func (api *CreateNavMenuItem) Parent(parent int) *CreateNavMenuItem {
	api.body["parent"] = parent
	return api
}

// AttrTitle sets text for the title attribute of the link element.
func (api *CreateNavMenuItem) AttrTitle(attrTitle string) *CreateNavMenuItem {
	api.body["attr_title"] = attrTitle
	return api
}

// Classes sets class names for the link element.
func (api *CreateNavMenuItem) Classes(classes ...string) *CreateNavMenuItem {
	api.body["classes"] = classes
	return api
}

// Description sets the description of this menu item.
func (api *CreateNavMenuItem) Description(description string) *CreateNavMenuItem {
	api.body["description"] = description
	return api
}

// MenuOrder sets the DB ID of the nav_menu_item that is this item's menu parent, if any, otherwise 0.
func (api *CreateNavMenuItem) MenuOrder(menuOrder int) *CreateNavMenuItem {
	api.body["menu_order"] = menuOrder
	return api
}

// Object sets the type of object originally represented (category, post, attachment, custom).
func (api *CreateNavMenuItem) Object(object string) *CreateNavMenuItem {
	api.body["object"] = object
	return api
}

// ObjectID sets the database ID of the original object this menu item represents.
func (api *CreateNavMenuItem) ObjectID(objectID int) *CreateNavMenuItem {
	api.body["object_id"] = objectID
	return api
}

// Target sets the target attribute of the link element (_blank, etc.).
func (api *CreateNavMenuItem) Target(target string) *CreateNavMenuItem {
	api.body["target"] = target
	return api
}

// URL sets the URL to which this menu item points.
func (api *CreateNavMenuItem) URL(url string) *CreateNavMenuItem {
	api.body["url"] = url
	return api
}

// XFN sets the XFN relationship expressed in the link.
func (api *CreateNavMenuItem) XFN(xfn ...string) *CreateNavMenuItem {
	api.body["xfn"] = xfn
	return api
}

// Menus sets the terms assigned to the object in the nav_menu taxonomy.
func (api *CreateNavMenuItem) Menus(menuID int) *CreateNavMenuItem {
	api.body["menus"] = menuID
	return api
}

// Meta sets meta fields.
func (api *CreateNavMenuItem) Meta(meta any) *CreateNavMenuItem {
	api.body["meta"] = meta
	return api
}

// Do executes the create request.
func (api *CreateNavMenuItem) Do() (item *NavMenuItem, err error) {
	endpoint := api.client.endpoint + api.endpoint

	restyClient := api.client.httpClient.R()
	if api.client.auth.Username != "" && api.client.auth.Password != "" {
		restyClient = restyClient.SetBasicAuth(api.client.auth.Username, api.client.auth.Password)
	}

	resp, err := restyClient.
		SetHeader("Content-Type", "application/json").
		SetHeader("Accept", "application/json").
		SetBody(api.body).
		SetResult(&item).
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

	return item, nil
}

// UpdateNavMenuItem handles modifying an existing navigation menu item.
type UpdateNavMenuItem struct {
	endpoint string
	client   *RestClient
	body     map[string]any
}

// Title sets the title for the object.
func (api *UpdateNavMenuItem) Title(title string) *UpdateNavMenuItem {
	api.body["title"] = title
	return api
}

// Type sets the family of objects originally represented (taxonomy, post_type, post_type_archive, custom).
func (api *UpdateNavMenuItem) Type(itemType string) *UpdateNavMenuItem {
	api.body["type"] = itemType
	return api
}

// Status sets the named status for the object (publish, future, draft, pending, private).
func (api *UpdateNavMenuItem) Status(status string) *UpdateNavMenuItem {
	api.body["status"] = status
	return api
}

// Parent sets the ID for the parent of the object.
func (api *UpdateNavMenuItem) Parent(parent int) *UpdateNavMenuItem {
	api.body["parent"] = parent
	return api
}

// AttrTitle sets text for the title attribute of the link element.
func (api *UpdateNavMenuItem) AttrTitle(attrTitle string) *UpdateNavMenuItem {
	api.body["attr_title"] = attrTitle
	return api
}

// Classes sets class names for the link element.
func (api *UpdateNavMenuItem) Classes(classes ...string) *UpdateNavMenuItem {
	api.body["classes"] = classes
	return api
}

// Description sets the description of this menu item.
func (api *UpdateNavMenuItem) Description(description string) *UpdateNavMenuItem {
	api.body["description"] = description
	return api
}

// MenuOrder sets the DB ID of the nav_menu_item that is this item's menu parent, if any, otherwise 0.
func (api *UpdateNavMenuItem) MenuOrder(menuOrder int) *UpdateNavMenuItem {
	api.body["menu_order"] = menuOrder
	return api
}

// Object sets the type of object originally represented (category, post, attachment, custom).
func (api *UpdateNavMenuItem) Object(object string) *UpdateNavMenuItem {
	api.body["object"] = object
	return api
}

// ObjectID sets the database ID of the original object this menu item represents.
func (api *UpdateNavMenuItem) ObjectID(objectID int) *UpdateNavMenuItem {
	api.body["object_id"] = objectID
	return api
}

// Target sets the target attribute of the link element (_blank, etc.).
func (api *UpdateNavMenuItem) Target(target string) *UpdateNavMenuItem {
	api.body["target"] = target
	return api
}

// URL sets the URL to which this menu item points.
func (api *UpdateNavMenuItem) URL(url string) *UpdateNavMenuItem {
	api.body["url"] = url
	return api
}

// XFN sets the XFN relationship expressed in the link.
func (api *UpdateNavMenuItem) XFN(xfn ...string) *UpdateNavMenuItem {
	api.body["xfn"] = xfn
	return api
}

// Menus sets the terms assigned to the object in the nav_menu taxonomy.
func (api *UpdateNavMenuItem) Menus(menuID int) *UpdateNavMenuItem {
	api.body["menus"] = menuID
	return api
}

// Meta sets meta fields.
func (api *UpdateNavMenuItem) Meta(meta any) *UpdateNavMenuItem {
	api.body["meta"] = meta
	return api
}

// Do executes the update request.
func (api *UpdateNavMenuItem) Do() (item *NavMenuItem, err error) {
	endpoint := api.client.endpoint + api.endpoint

	restyClient := api.client.httpClient.R()
	if api.client.auth.Username != "" && api.client.auth.Password != "" {
		restyClient = restyClient.SetBasicAuth(api.client.auth.Username, api.client.auth.Password)
	}

	resp, err := restyClient.
		SetHeader("Content-Type", "application/json").
		SetHeader("Accept", "application/json").
		SetBody(api.body).
		SetResult(&item).
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

	return item, nil
}

type deleteNavMenuItemEnvelope struct {
	NavMenuItem
	Deleted  bool         `json:"deleted"`
	Previous *NavMenuItem `json:"previous"`
}

// DeleteNavMenuItem handles deleting a navigation menu item.
type DeleteNavMenuItem struct {
	endpoint string
	client   *RestClient
	force    bool
}

// Force sets whether to bypass Trash and force deletion.
func (api *DeleteNavMenuItem) Force(force ...bool) *DeleteNavMenuItem {
	if len(force) > 0 {
		api.force = force[0]
	} else {
		api.force = true
	}
	return api
}

// Do executes the delete request.
func (api *DeleteNavMenuItem) Do() (item *NavMenuItem, err error) {
	endpoint := api.client.endpoint + api.endpoint

	restyClient := api.client.httpClient.R()
	if api.client.auth.Username != "" && api.client.auth.Password != "" {
		restyClient = restyClient.SetBasicAuth(api.client.auth.Username, api.client.auth.Password)
	}

	if api.force {
		restyClient = restyClient.SetQueryParam("force", "true")
	}

	var env deleteNavMenuItemEnvelope
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

	return &env.NavMenuItem, nil
}
