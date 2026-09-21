package gowprest

import (
	"encoding/json"
	"strconv"
	"strings"
)

// NavMenuItemRevision represents a WordPress navigation menu item revision (autosave).
type NavMenuItemRevision struct {
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
	PreviewLink string         `json:"preview_link,omitempty"`
	Meta        any            `json:"meta,omitempty"`
	Links       map[string]any `json:"_links,omitempty"`
}

// NavMenuItemAutosave is an alias for NavMenuItemRevision.
type NavMenuItemAutosave = NavMenuItemRevision

// TitleString returns the revision's title string from either Raw or Rendered.
func (r *NavMenuItemRevision) TitleString() string {
	if r.Title == nil {
		return ""
	}
	if r.Title.Raw != "" {
		return r.Title.Raw
	}
	return r.Title.Rendered
}

// NavMenuItemRevisions provides operations for nav menu item revisions / autosaves.
type NavMenuItemRevisions struct {
	client   *RestClient
	parentID int
}

// Revisions returns a NavMenuItemRevisions service instance for the parent menu item ID.
func (api *NavMenuItems) Revisions(parentID int) *NavMenuItemRevisions {
	return &NavMenuItemRevisions{
		client:   api.client,
		parentID: parentID,
	}
}

// Autosaves returns a NavMenuItemRevisions service instance for the parent menu item ID.
func (api *NavMenuItems) Autosaves(parentID int) *NavMenuItemRevisions {
	return api.Revisions(parentID)
}

// NavMenuItemRevisions returns a NavMenuItemRevisions service instance for the parent menu item ID.
func (c *RestClient) NavMenuItemRevisions(parentID int) *NavMenuItemRevisions {
	return &NavMenuItemRevisions{
		client:   c,
		parentID: parentID,
	}
}

// NavMenuItemAutosaves is an alias for NavMenuItemRevisions.
func (c *RestClient) NavMenuItemAutosaves(parentID int) *NavMenuItemRevisions {
	return c.NavMenuItemRevisions(parentID)
}

// MenuItemRevisions is an alias for NavMenuItemRevisions.
func (c *RestClient) MenuItemRevisions(parentID int) *NavMenuItemRevisions {
	return c.NavMenuItemRevisions(parentID)
}

// MenuItemAutosaves is an alias for NavMenuItemRevisions.
func (c *RestClient) MenuItemAutosaves(parentID int) *NavMenuItemRevisions {
	return c.NavMenuItemRevisions(parentID)
}

// List returns a ListNavMenuItemRevisions builder to list autosaves for the parent menu item.
func (api *NavMenuItemRevisions) List() *ListNavMenuItemRevisions {
	return &ListNavMenuItemRevisions{
		endpoint:  "/wp/v2/menu-items/" + strconv.Itoa(api.parentID) + "/autosaves",
		client:    api.client,
		arguments: make(map[string]string),
	}
}

// Retrieve returns a RetrieveNavMenuItemRevision builder to retrieve a specific autosave by ID.
func (api *NavMenuItemRevisions) Retrieve(id int) *RetrieveNavMenuItemRevision {
	return &RetrieveNavMenuItemRevision{
		endpoint:  "/wp/v2/menu-items/" + strconv.Itoa(api.parentID) + "/autosaves/" + strconv.Itoa(id),
		client:    api.client,
		arguments: make(map[string]string),
	}
}

// Create returns a CreateNavMenuItemRevision builder to create an autosave for the parent menu item.
func (api *NavMenuItemRevisions) Create() *CreateNavMenuItemRevision {
	body := make(map[string]any)
	body["parent"] = api.parentID
	return &CreateNavMenuItemRevision{
		endpoint: "/wp/v2/menu-items/" + strconv.Itoa(api.parentID) + "/autosaves",
		client:   api.client,
		body:     body,
	}
}

// ListNavMenuItemRevisions handles querying nav menu item autosaves.
type ListNavMenuItemRevisions struct {
	endpoint  string
	client    *RestClient
	arguments map[string]string
}

// Context sets the context parameter (e.g. "view", "edit", "embed").
func (api *ListNavMenuItemRevisions) Context(ctx string) *ListNavMenuItemRevisions {
	api.arguments["context"] = ctx
	return api
}

// ContextView sets the context parameter to "view".
func (api *ListNavMenuItemRevisions) ContextView() *ListNavMenuItemRevisions {
	return api.Context("view")
}

// ContextEdit sets the context parameter to "edit".
func (api *ListNavMenuItemRevisions) ContextEdit() *ListNavMenuItemRevisions {
	return api.Context("edit")
}

// ContextEmbed sets the context parameter to "embed".
func (api *ListNavMenuItemRevisions) ContextEmbed() *ListNavMenuItemRevisions {
	return api.Context("embed")
}

// Fields limits the response to specific fields.
func (api *ListNavMenuItemRevisions) Fields(fields ...string) *ListNavMenuItemRevisions {
	api.arguments["_fields"] = strings.Join(fields, ",")
	return api
}

// Do executes the request and returns the list of revisions.
func (api *ListNavMenuItemRevisions) Do() (revisions []NavMenuItemRevision, err error) {
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

// RetrieveNavMenuItemRevision handles retrieving a single nav menu item autosave.
type RetrieveNavMenuItemRevision struct {
	endpoint  string
	client    *RestClient
	arguments map[string]string
}

// Context sets the context parameter (e.g. "view", "edit", "embed").
func (api *RetrieveNavMenuItemRevision) Context(ctx string) *RetrieveNavMenuItemRevision {
	api.arguments["context"] = ctx
	return api
}

// ContextView sets the context parameter to "view".
func (api *RetrieveNavMenuItemRevision) ContextView() *RetrieveNavMenuItemRevision {
	return api.Context("view")
}

// ContextEdit sets the context parameter to "edit".
func (api *RetrieveNavMenuItemRevision) ContextEdit() *RetrieveNavMenuItemRevision {
	return api.Context("edit")
}

// ContextEmbed sets the context parameter to "embed".
func (api *RetrieveNavMenuItemRevision) ContextEmbed() *RetrieveNavMenuItemRevision {
	return api.Context("embed")
}

// Fields limits the response to specific fields.
func (api *RetrieveNavMenuItemRevision) Fields(fields ...string) *RetrieveNavMenuItemRevision {
	api.arguments["_fields"] = strings.Join(fields, ",")
	return api
}

// Do executes the request and returns the revision.
func (api *RetrieveNavMenuItemRevision) Do() (revision *NavMenuItemRevision, err error) {
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

// CreateNavMenuItemRevision handles creating a nav menu item autosave/revision.
type CreateNavMenuItemRevision struct {
	endpoint string
	client   *RestClient
	body     map[string]any
}

// Title sets the title for the object.
func (api *CreateNavMenuItemRevision) Title(title string) *CreateNavMenuItemRevision {
	api.body["title"] = title
	return api
}

// Type sets the family of objects originally represented (taxonomy, post_type, post_type_archive, custom).
func (api *CreateNavMenuItemRevision) Type(itemType string) *CreateNavMenuItemRevision {
	api.body["type"] = itemType
	return api
}

// Status sets the named status for the object (publish, future, draft, pending, private).
func (api *CreateNavMenuItemRevision) Status(status string) *CreateNavMenuItemRevision {
	api.body["status"] = status
	return api
}

// AttrTitle sets text for the title attribute of the link element.
func (api *CreateNavMenuItemRevision) AttrTitle(attrTitle string) *CreateNavMenuItemRevision {
	api.body["attr_title"] = attrTitle
	return api
}

// Classes sets class names for the link element.
func (api *CreateNavMenuItemRevision) Classes(classes ...string) *CreateNavMenuItemRevision {
	api.body["classes"] = classes
	return api
}

// Description sets the description of this menu item.
func (api *CreateNavMenuItemRevision) Description(description string) *CreateNavMenuItemRevision {
	api.body["description"] = description
	return api
}

// MenuOrder sets the DB ID of the nav_menu_item that is this item's menu parent, if any, otherwise 0.
func (api *CreateNavMenuItemRevision) MenuOrder(menuOrder int) *CreateNavMenuItemRevision {
	api.body["menu_order"] = menuOrder
	return api
}

// Object sets the type of object originally represented (category, post, attachment, custom).
func (api *CreateNavMenuItemRevision) Object(object string) *CreateNavMenuItemRevision {
	api.body["object"] = object
	return api
}

// ObjectID sets the database ID of the original object this menu item represents.
func (api *CreateNavMenuItemRevision) ObjectID(objectID int) *CreateNavMenuItemRevision {
	api.body["object_id"] = objectID
	return api
}

// Target sets the target attribute of the link element (_blank, etc.).
func (api *CreateNavMenuItemRevision) Target(target string) *CreateNavMenuItemRevision {
	api.body["target"] = target
	return api
}

// URL sets the URL to which this menu item points.
func (api *CreateNavMenuItemRevision) URL(url string) *CreateNavMenuItemRevision {
	api.body["url"] = url
	return api
}

// XFN sets the XFN relationship expressed in the link.
func (api *CreateNavMenuItemRevision) XFN(xfn ...string) *CreateNavMenuItemRevision {
	api.body["xfn"] = xfn
	return api
}

// Menus sets the terms assigned to the object in the nav_menu taxonomy.
func (api *CreateNavMenuItemRevision) Menus(menuID int) *CreateNavMenuItemRevision {
	api.body["menus"] = menuID
	return api
}

// Meta sets meta fields.
func (api *CreateNavMenuItemRevision) Meta(meta any) *CreateNavMenuItemRevision {
	api.body["meta"] = meta
	return api
}

// Do executes the create request.
func (api *CreateNavMenuItemRevision) Do() (revision *NavMenuItemRevision, err error) {
	endpoint := api.client.endpoint + api.endpoint

	restyClient := api.client.httpClient.R()
	if api.client.auth.Username != "" && api.client.auth.Password != "" {
		restyClient = restyClient.SetBasicAuth(api.client.auth.Username, api.client.auth.Password)
	}

	resp, err := restyClient.
		SetHeader("Content-Type", "application/json").
		SetHeader("Accept", "application/json").
		SetBody(api.body).
		SetResult(&revision).
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

	return revision, nil
}
