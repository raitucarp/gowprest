package gowprest

import (
	"encoding/json"
	"strconv"
	"strings"
)

// TemplatePartRevision represents a WordPress template part revision or autosave record.
type TemplatePartRevision struct {
	ID             string          `json:"id,omitempty"`
	WPID           int             `json:"wp_id,omitempty"`
	Parent         int             `json:"parent,omitempty"`
	Author         int             `json:"author,omitempty"`
	AuthorText     string          `json:"author_text,omitempty"`
	Date           *string         `json:"date,omitempty"`
	DateGMT        *string         `json:"date_gmt,omitempty"`
	Modified       *string         `json:"modified,omitempty"`
	ModifiedGMT    *string         `json:"modified_gmt,omitempty"`
	Slug           string          `json:"slug,omitempty"`
	Title          *TemplateObject `json:"title,omitempty"`
	Content        *TemplateObject `json:"content,omitempty"`
	Description    string          `json:"description,omitempty"`
	Status         string          `json:"status,omitempty"`
	Theme          string          `json:"theme,omitempty"`
	Type           string          `json:"type,omitempty"`
	Source         string          `json:"source,omitempty"`
	Origin         *string         `json:"origin,omitempty"`
	HasThemeFile   bool            `json:"has_theme_file,omitempty"`
	Area           string          `json:"area,omitempty"`
	OriginalSource string          `json:"original_source,omitempty"`
	Plugin         string          `json:"plugin,omitempty"`
	Links          map[string]any  `json:"_links,omitempty"`
}

// TitleString returns the revision's title string from either Raw or Rendered.
func (r *TemplatePartRevision) TitleString() string {
	if r.Title == nil {
		return ""
	}
	if r.Title.Raw != "" {
		return r.Title.Raw
	}
	return r.Title.Rendered
}

// ContentString returns the revision's content string from either Raw or Rendered.
func (r *TemplatePartRevision) ContentString() string {
	if r.Content == nil {
		return ""
	}
	if r.Content.Raw != "" {
		return r.Content.Raw
	}
	return r.Content.Rendered
}

// TemplatePartRevisions anchors revision-related operations for a specific Template Part parent.
type TemplatePartRevisions struct {
	client *RestClient
	parent string
}

// TemplatePartRevisions returns a TemplatePartRevisions service instance for the given parent template part ID.
func (c *RestClient) TemplatePartRevisions(parent string) *TemplatePartRevisions {
	return &TemplatePartRevisions{
		client: c,
		parent: parent,
	}
}

// Revisions returns a TemplatePartRevisions service instance for the given parent template part ID.
func (api *TemplateParts) Revisions(parent string) *TemplatePartRevisions {
	return &TemplatePartRevisions{
		client: api.client,
		parent: parent,
	}
}

// Autosaves returns a TemplatePartAutosaves service instance for this template part parent.
func (api *TemplatePartRevisions) Autosaves() *TemplatePartAutosaves {
	return &TemplatePartAutosaves{
		client: api.client,
		parent: api.parent,
	}
}

// List returns a ListTemplatePartRevisions builder to list revisions for the parent template part.
func (api *TemplatePartRevisions) List() *ListTemplatePartRevisions {
	return &ListTemplatePartRevisions{
		endpoint:  "/wp/v2/template-parts/" + api.parent + "/revisions",
		client:    api.client,
		arguments: make(map[string]string),
	}
}

// Retrieve returns a RetrieveTemplatePartRevision builder for a specific revision ID.
func (api *TemplatePartRevisions) Retrieve(id int) *RetrieveTemplatePartRevision {
	return &RetrieveTemplatePartRevision{
		endpoint:  "/wp/v2/template-parts/" + api.parent + "/revisions/" + strconv.Itoa(id),
		client:    api.client,
		arguments: make(map[string]string),
	}
}

// Delete returns a DeleteTemplatePartRevision builder to delete a revision.
func (api *TemplatePartRevisions) Delete(id int) *DeleteTemplatePartRevision {
	return &DeleteTemplatePartRevision{
		endpoint: "/wp/v2/template-parts/" + api.parent + "/revisions/" + strconv.Itoa(id),
		client:   api.client,
	}
}

// ListTemplatePartRevisions handles querying template part revisions.
type ListTemplatePartRevisions struct {
	endpoint  string
	client    *RestClient
	arguments map[string]string
}

// Context sets the context parameter (e.g. "view", "edit", "embed").
func (api *ListTemplatePartRevisions) Context(ctx string) *ListTemplatePartRevisions {
	api.arguments["context"] = ctx
	return api
}

// ContextView sets the context parameter to "view".
func (api *ListTemplatePartRevisions) ContextView() *ListTemplatePartRevisions {
	return api.Context("view")
}

// ContextEdit sets the context parameter to "edit".
func (api *ListTemplatePartRevisions) ContextEdit() *ListTemplatePartRevisions {
	return api.Context("edit")
}

// ContextEmbed sets the context parameter to "embed".
func (api *ListTemplatePartRevisions) ContextEmbed() *ListTemplatePartRevisions {
	return api.Context("embed")
}

// Page sets the page number of the collection.
func (api *ListTemplatePartRevisions) Page(page int) *ListTemplatePartRevisions {
	api.arguments["page"] = strconv.Itoa(page)
	return api
}

// PerPage sets the number of items to return per page.
func (api *ListTemplatePartRevisions) PerPage(perPage int) *ListTemplatePartRevisions {
	api.arguments["per_page"] = strconv.Itoa(perPage)
	return api
}

// Search limits results to those matching a string.
func (api *ListTemplatePartRevisions) Search(query string) *ListTemplatePartRevisions {
	api.arguments["search"] = query
	return api
}

// Exclude ensures result set excludes specific revision IDs.
func (api *ListTemplatePartRevisions) Exclude(ids ...int) *ListTemplatePartRevisions {
	strIDs := make([]string, 0, len(ids))
	for _, id := range ids {
		strIDs = append(strIDs, strconv.Itoa(id))
	}
	api.arguments["exclude"] = strings.Join(strIDs, ",")
	return api
}

// Include limits result set to specific revision IDs.
func (api *ListTemplatePartRevisions) Include(ids ...int) *ListTemplatePartRevisions {
	strIDs := make([]string, 0, len(ids))
	for _, id := range ids {
		strIDs = append(strIDs, strconv.Itoa(id))
	}
	api.arguments["include"] = strings.Join(strIDs, ",")
	return api
}

// Offset offsets the result set by a specific number of items.
func (api *ListTemplatePartRevisions) Offset(offset int) *ListTemplatePartRevisions {
	api.arguments["offset"] = strconv.Itoa(offset)
	return api
}

// Order sets sort order ("asc" or "desc").
func (api *ListTemplatePartRevisions) Order(order string) *ListTemplatePartRevisions {
	api.arguments["order"] = order
	return api
}

// OrderAsc sets sort order to ascending.
func (api *ListTemplatePartRevisions) OrderAsc() *ListTemplatePartRevisions {
	return api.Order("asc")
}

// OrderDesc sets sort order to descending.
func (api *ListTemplatePartRevisions) OrderDesc() *ListTemplatePartRevisions {
	return api.Order("desc")
}

// OrderBy sets collection sort attribute.
func (api *ListTemplatePartRevisions) OrderBy(orderBy string) *ListTemplatePartRevisions {
	api.arguments["orderby"] = orderBy
	return api
}

// Fields limits the response to specific fields.
func (api *ListTemplatePartRevisions) Fields(fields ...string) *ListTemplatePartRevisions {
	api.arguments["_fields"] = strings.Join(fields, ",")
	return api
}

// Do executes the list request.
func (api *ListTemplatePartRevisions) Do() (revisions []TemplatePartRevision, err error) {
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

// RetrieveTemplatePartRevision handles retrieving a single template part revision.
type RetrieveTemplatePartRevision struct {
	endpoint  string
	client    *RestClient
	arguments map[string]string
}

// Context sets the context parameter (e.g. "view", "edit", "embed").
func (api *RetrieveTemplatePartRevision) Context(ctx string) *RetrieveTemplatePartRevision {
	api.arguments["context"] = ctx
	return api
}

// ContextView sets the context parameter to "view".
func (api *RetrieveTemplatePartRevision) ContextView() *RetrieveTemplatePartRevision {
	return api.Context("view")
}

// ContextEdit sets the context parameter to "edit".
func (api *RetrieveTemplatePartRevision) ContextEdit() *RetrieveTemplatePartRevision {
	return api.Context("edit")
}

// ContextEmbed sets the context parameter to "embed".
func (api *RetrieveTemplatePartRevision) ContextEmbed() *RetrieveTemplatePartRevision {
	return api.Context("embed")
}

// Fields limits the response to specific fields.
func (api *RetrieveTemplatePartRevision) Fields(fields ...string) *RetrieveTemplatePartRevision {
	api.arguments["_fields"] = strings.Join(fields, ",")
	return api
}

// Do executes the retrieve request.
func (api *RetrieveTemplatePartRevision) Do() (revision *TemplatePartRevision, err error) {
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

type deleteTemplatePartRevisionEnvelope struct {
	TemplatePartRevision
	Deleted  bool                  `json:"deleted"`
	Previous *TemplatePartRevision `json:"previous"`
}

// DeleteTemplatePartRevision handles deleting a template part revision.
type DeleteTemplatePartRevision struct {
	endpoint string
	client   *RestClient
	force    bool
}

// Force sets whether to bypass Trash and force permanent deletion.
func (api *DeleteTemplatePartRevision) Force() *DeleteTemplatePartRevision {
	api.force = true
	return api
}

// Do executes the delete request.
func (api *DeleteTemplatePartRevision) Do() (revision *TemplatePartRevision, err error) {
	endpoint := api.client.endpoint + api.endpoint
	var env deleteTemplatePartRevisionEnvelope

	restyClient := api.client.httpClient.R()
	if api.client.auth.Username != "" && api.client.auth.Password != "" {
		restyClient = restyClient.SetBasicAuth(api.client.auth.Username, api.client.auth.Password)
	}

	if api.force {
		restyClient.SetQueryParam("force", "true")
	}

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

	if env.Deleted && env.Previous != nil && (env.Previous.WPID != 0 || env.Previous.ID != "") {
		return env.Previous, nil
	}

	return &env.TemplatePartRevision, nil
}

// TemplatePartAutosaves anchors autosave operations for a specific template part parent.
type TemplatePartAutosaves struct {
	client *RestClient
	parent string
}

// TemplatePartAutosaves returns a TemplatePartAutosaves service instance for the given parent template part ID.
func (c *RestClient) TemplatePartAutosaves(parent string) *TemplatePartAutosaves {
	return &TemplatePartAutosaves{
		client: c,
		parent: parent,
	}
}

// Autosaves returns a TemplatePartAutosaves service instance for the given parent template part ID.
func (api *TemplateParts) Autosaves(parent string) *TemplatePartAutosaves {
	return &TemplatePartAutosaves{
		client: api.client,
		parent: parent,
	}
}

// List returns a ListTemplatePartAutosaves builder to query autosaves for the parent template part.
func (api *TemplatePartAutosaves) List() *ListTemplatePartAutosaves {
	return &ListTemplatePartAutosaves{
		endpoint:  "/wp/v2/template-parts/" + api.parent + "/autosaves",
		client:    api.client,
		arguments: make(map[string]string),
	}
}

// Retrieve returns a RetrieveTemplatePartAutosave builder for a specific autosave revision ID.
func (api *TemplatePartAutosaves) Retrieve(id int) *RetrieveTemplatePartAutosave {
	return &RetrieveTemplatePartAutosave{
		endpoint:  "/wp/v2/template-parts/" + api.parent + "/autosaves/" + strconv.Itoa(id),
		client:    api.client,
		arguments: make(map[string]string),
	}
}

// Create returns a CreateTemplatePartAutosave builder to create or update an autosave.
func (api *TemplatePartAutosaves) Create() *CreateTemplatePartAutosave {
	return &CreateTemplatePartAutosave{
		endpoint: "/wp/v2/template-parts/" + api.parent + "/autosaves",
		client:   api.client,
		data:     TemplatePartData{},
	}
}

// ListTemplatePartAutosaves handles querying template part autosaves.
type ListTemplatePartAutosaves struct {
	endpoint  string
	client    *RestClient
	arguments map[string]string
}

// Context sets the context parameter (e.g. "view", "edit", "embed").
func (api *ListTemplatePartAutosaves) Context(ctx string) *ListTemplatePartAutosaves {
	api.arguments["context"] = ctx
	return api
}

// ContextView sets the context parameter to "view".
func (api *ListTemplatePartAutosaves) ContextView() *ListTemplatePartAutosaves {
	return api.Context("view")
}

// ContextEdit sets the context parameter to "edit".
func (api *ListTemplatePartAutosaves) ContextEdit() *ListTemplatePartAutosaves {
	return api.Context("edit")
}

// ContextEmbed sets the context parameter to "embed".
func (api *ListTemplatePartAutosaves) ContextEmbed() *ListTemplatePartAutosaves {
	return api.Context("embed")
}

// Fields limits the response to specific fields.
func (api *ListTemplatePartAutosaves) Fields(fields ...string) *ListTemplatePartAutosaves {
	api.arguments["_fields"] = strings.Join(fields, ",")
	return api
}

// Do executes the list request.
func (api *ListTemplatePartAutosaves) Do() (autosaves []TemplatePartRevision, err error) {
	endpoint := api.client.endpoint + api.endpoint

	restyClient := api.client.httpClient.R()
	if api.client.auth.Username != "" && api.client.auth.Password != "" {
		restyClient = restyClient.SetBasicAuth(api.client.auth.Username, api.client.auth.Password)
	}

	resp, err := restyClient.
		SetHeader("Accept", "application/json").
		SetResult(&autosaves).
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

	return autosaves, nil
}

// RetrieveTemplatePartAutosave handles retrieving a single template part autosave.
type RetrieveTemplatePartAutosave struct {
	endpoint  string
	client    *RestClient
	arguments map[string]string
}

// Context sets the context parameter (e.g. "view", "edit", "embed").
func (api *RetrieveTemplatePartAutosave) Context(ctx string) *RetrieveTemplatePartAutosave {
	api.arguments["context"] = ctx
	return api
}

// ContextView sets the context parameter to "view".
func (api *RetrieveTemplatePartAutosave) ContextView() *RetrieveTemplatePartAutosave {
	return api.Context("view")
}

// ContextEdit sets the context parameter to "edit".
func (api *RetrieveTemplatePartAutosave) ContextEdit() *RetrieveTemplatePartAutosave {
	return api.Context("edit")
}

// ContextEmbed sets the context parameter to "embed".
func (api *RetrieveTemplatePartAutosave) ContextEmbed() *RetrieveTemplatePartAutosave {
	return api.Context("embed")
}

// Fields limits the response to specific fields.
func (api *RetrieveTemplatePartAutosave) Fields(fields ...string) *RetrieveTemplatePartAutosave {
	api.arguments["_fields"] = strings.Join(fields, ",")
	return api
}

// Do executes the retrieve request.
func (api *RetrieveTemplatePartAutosave) Do() (autosave *TemplatePartRevision, err error) {
	endpoint := api.client.endpoint + api.endpoint

	restyClient := api.client.httpClient.R()
	if api.client.auth.Username != "" && api.client.auth.Password != "" {
		restyClient = restyClient.SetBasicAuth(api.client.auth.Username, api.client.auth.Password)
	}

	resp, err := restyClient.
		SetHeader("Accept", "application/json").
		SetResult(&autosave).
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

	return autosave, nil
}

// CreateTemplatePartAutosave handles creating a template part autosave.
type CreateTemplatePartAutosave struct {
	endpoint string
	client   *RestClient
	data     TemplatePartData
}

// Slug sets the unique slug identifying the template part.
func (api *CreateTemplatePartAutosave) Slug(slug string) *CreateTemplatePartAutosave {
	api.data.Slug = slug
	return api
}

// Theme sets the theme identifier for the template part.
func (api *CreateTemplatePartAutosave) Theme(theme string) *CreateTemplatePartAutosave {
	api.data.Theme = theme
	return api
}

// Type sets the type of template part.
func (api *CreateTemplatePartAutosave) Type(typ string) *CreateTemplatePartAutosave {
	api.data.Type = typ
	return api
}

// Content sets the content of template part.
func (api *CreateTemplatePartAutosave) Content(content any) *CreateTemplatePartAutosave {
	api.data.Content = content
	return api
}

// Title sets the title of template part.
func (api *CreateTemplatePartAutosave) Title(title any) *CreateTemplatePartAutosave {
	api.data.Title = title
	return api
}

// Description sets the description of template part.
func (api *CreateTemplatePartAutosave) Description(desc string) *CreateTemplatePartAutosave {
	api.data.Description = desc
	return api
}

// Status sets the status of template part.
func (api *CreateTemplatePartAutosave) Status(status string) *CreateTemplatePartAutosave {
	api.data.Status = status
	return api
}

// Author sets the ID for the author of the template part.
func (api *CreateTemplatePartAutosave) Author(author int) *CreateTemplatePartAutosave {
	api.data.Author = author
	return api
}

// Area sets the area where the template part is intended for use.
func (api *CreateTemplatePartAutosave) Area(area string) *CreateTemplatePartAutosave {
	api.data.Area = area
	return api
}

// WithData populates the request body from a TemplatePartData struct.
func (api *CreateTemplatePartAutosave) WithData(data TemplatePartData) *CreateTemplatePartAutosave {
	api.data = data
	return api
}

// Do executes the create autosave request.
func (api *CreateTemplatePartAutosave) Do() (autosave *TemplatePartRevision, err error) {
	endpoint := api.client.endpoint + api.endpoint

	restyClient := api.client.httpClient.R()
	if api.client.auth.Username != "" && api.client.auth.Password != "" {
		restyClient = restyClient.SetBasicAuth(api.client.auth.Username, api.client.auth.Password)
	}

	resp, err := restyClient.
		SetHeader("Content-Type", "application/json").
		SetHeader("Accept", "application/json").
		SetBody(api.data).
		SetResult(&autosave).
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

	return autosave, nil
}
