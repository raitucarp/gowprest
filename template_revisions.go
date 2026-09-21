package gowprest

import (
	"encoding/json"
	"strconv"
	"strings"
)

// TemplateRevision represents a WordPress template revision or autosave record.
type TemplateRevision struct {
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
	IsCustom       bool            `json:"is_custom,omitempty"`
	OriginalSource string          `json:"original_source,omitempty"`
	Plugin         string          `json:"plugin,omitempty"`
	Links          map[string]any  `json:"_links,omitempty"`
}

// TitleString returns the revision's title string from either Raw or Rendered.
func (r *TemplateRevision) TitleString() string {
	if r.Title == nil {
		return ""
	}
	if r.Title.Raw != "" {
		return r.Title.Raw
	}
	return r.Title.Rendered
}

// ContentString returns the revision's content string from either Raw or Rendered.
func (r *TemplateRevision) ContentString() string {
	if r.Content == nil {
		return ""
	}
	if r.Content.Raw != "" {
		return r.Content.Raw
	}
	return r.Content.Rendered
}

// TemplateRevisions anchors revision-related operations for a specific Template parent.
type TemplateRevisions struct {
	client *RestClient
	parent string
}

// TemplateRevisions returns a TemplateRevisions service instance for the given parent template ID.
func (c *RestClient) TemplateRevisions(parent string) *TemplateRevisions {
	return &TemplateRevisions{
		client: c,
		parent: parent,
	}
}

// Revisions returns a TemplateRevisions service instance for the given parent template ID.
func (api *Templates) Revisions(parent string) *TemplateRevisions {
	return &TemplateRevisions{
		client: api.client,
		parent: parent,
	}
}

// Autosaves returns a TemplateAutosaves service instance for this template parent.
func (api *TemplateRevisions) Autosaves() *TemplateAutosaves {
	return &TemplateAutosaves{
		client: api.client,
		parent: api.parent,
	}
}

// List returns a ListTemplateRevisions builder to list revisions for the parent template.
func (api *TemplateRevisions) List() *ListTemplateRevisions {
	return &ListTemplateRevisions{
		endpoint:  "/wp/v2/templates/" + api.parent + "/revisions",
		client:    api.client,
		arguments: make(map[string]string),
	}
}

// Retrieve returns a RetrieveTemplateRevision builder for a specific revision ID.
func (api *TemplateRevisions) Retrieve(id int) *RetrieveTemplateRevision {
	return &RetrieveTemplateRevision{
		endpoint:  "/wp/v2/templates/" + api.parent + "/revisions/" + strconv.Itoa(id),
		client:    api.client,
		arguments: make(map[string]string),
	}
}

// Delete returns a DeleteTemplateRevision builder to delete a revision.
func (api *TemplateRevisions) Delete(id int) *DeleteTemplateRevision {
	return &DeleteTemplateRevision{
		endpoint: "/wp/v2/templates/" + api.parent + "/revisions/" + strconv.Itoa(id),
		client:   api.client,
	}
}

// ListTemplateRevisions handles querying template revisions.
type ListTemplateRevisions struct {
	endpoint  string
	client    *RestClient
	arguments map[string]string
}

// Context sets the context parameter (e.g. "view", "edit", "embed").
func (api *ListTemplateRevisions) Context(ctx string) *ListTemplateRevisions {
	api.arguments["context"] = ctx
	return api
}

// ContextView sets the context parameter to "view".
func (api *ListTemplateRevisions) ContextView() *ListTemplateRevisions {
	return api.Context("view")
}

// ContextEdit sets the context parameter to "edit".
func (api *ListTemplateRevisions) ContextEdit() *ListTemplateRevisions {
	return api.Context("edit")
}

// ContextEmbed sets the context parameter to "embed".
func (api *ListTemplateRevisions) ContextEmbed() *ListTemplateRevisions {
	return api.Context("embed")
}

// Page sets the page number of the collection.
func (api *ListTemplateRevisions) Page(page int) *ListTemplateRevisions {
	api.arguments["page"] = strconv.Itoa(page)
	return api
}

// PerPage sets the number of items to return per page.
func (api *ListTemplateRevisions) PerPage(perPage int) *ListTemplateRevisions {
	api.arguments["per_page"] = strconv.Itoa(perPage)
	return api
}

// Search limits results to those matching a string.
func (api *ListTemplateRevisions) Search(query string) *ListTemplateRevisions {
	api.arguments["search"] = query
	return api
}

// Exclude ensures result set excludes specific revision IDs.
func (api *ListTemplateRevisions) Exclude(ids ...int) *ListTemplateRevisions {
	strIDs := make([]string, 0, len(ids))
	for _, id := range ids {
		strIDs = append(strIDs, strconv.Itoa(id))
	}
	api.arguments["exclude"] = strings.Join(strIDs, ",")
	return api
}

// Include limits result set to specific revision IDs.
func (api *ListTemplateRevisions) Include(ids ...int) *ListTemplateRevisions {
	strIDs := make([]string, 0, len(ids))
	for _, id := range ids {
		strIDs = append(strIDs, strconv.Itoa(id))
	}
	api.arguments["include"] = strings.Join(strIDs, ",")
	return api
}

// Offset offsets the result set by a specific number of items.
func (api *ListTemplateRevisions) Offset(offset int) *ListTemplateRevisions {
	api.arguments["offset"] = strconv.Itoa(offset)
	return api
}

// Order sets sort order ("asc" or "desc").
func (api *ListTemplateRevisions) Order(order string) *ListTemplateRevisions {
	api.arguments["order"] = order
	return api
}

// OrderAsc sets sort order to ascending.
func (api *ListTemplateRevisions) OrderAsc() *ListTemplateRevisions {
	return api.Order("asc")
}

// OrderDesc sets sort order to descending.
func (api *ListTemplateRevisions) OrderDesc() *ListTemplateRevisions {
	return api.Order("desc")
}

// OrderBy sets collection sort attribute.
func (api *ListTemplateRevisions) OrderBy(orderBy string) *ListTemplateRevisions {
	api.arguments["orderby"] = orderBy
	return api
}

// Fields limits the response to specific fields.
func (api *ListTemplateRevisions) Fields(fields ...string) *ListTemplateRevisions {
	api.arguments["_fields"] = strings.Join(fields, ",")
	return api
}

// Do executes the list request.
func (api *ListTemplateRevisions) Do() (revisions []TemplateRevision, err error) {
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

// RetrieveTemplateRevision handles retrieving a single template revision.
type RetrieveTemplateRevision struct {
	endpoint  string
	client    *RestClient
	arguments map[string]string
}

// Context sets the context parameter (e.g. "view", "edit", "embed").
func (api *RetrieveTemplateRevision) Context(ctx string) *RetrieveTemplateRevision {
	api.arguments["context"] = ctx
	return api
}

// ContextView sets the context parameter to "view".
func (api *RetrieveTemplateRevision) ContextView() *RetrieveTemplateRevision {
	return api.Context("view")
}

// ContextEdit sets the context parameter to "edit".
func (api *RetrieveTemplateRevision) ContextEdit() *RetrieveTemplateRevision {
	return api.Context("edit")
}

// ContextEmbed sets the context parameter to "embed".
func (api *RetrieveTemplateRevision) ContextEmbed() *RetrieveTemplateRevision {
	return api.Context("embed")
}

// Fields limits the response to specific fields.
func (api *RetrieveTemplateRevision) Fields(fields ...string) *RetrieveTemplateRevision {
	api.arguments["_fields"] = strings.Join(fields, ",")
	return api
}

// Do executes the retrieve request.
func (api *RetrieveTemplateRevision) Do() (revision *TemplateRevision, err error) {
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

type deleteTemplateRevisionEnvelope struct {
	TemplateRevision
	Deleted  bool              `json:"deleted"`
	Previous *TemplateRevision `json:"previous"`
}

// DeleteTemplateRevision handles deleting a template revision.
type DeleteTemplateRevision struct {
	endpoint string
	client   *RestClient
	force    bool
}

// Force sets whether to bypass Trash and force permanent deletion.
func (api *DeleteTemplateRevision) Force() *DeleteTemplateRevision {
	api.force = true
	return api
}

// Do executes the delete request.
func (api *DeleteTemplateRevision) Do() (revision *TemplateRevision, err error) {
	endpoint := api.client.endpoint + api.endpoint
	var env deleteTemplateRevisionEnvelope

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

	return &env.TemplateRevision, nil
}

// TemplateAutosaves anchors autosave operations for a specific template parent.
type TemplateAutosaves struct {
	client *RestClient
	parent string
}

// TemplateAutosaves returns a TemplateAutosaves service instance for the given parent template ID.
func (c *RestClient) TemplateAutosaves(parent string) *TemplateAutosaves {
	return &TemplateAutosaves{
		client: c,
		parent: parent,
	}
}

// Autosaves returns a TemplateAutosaves service instance for the given parent template ID.
func (api *Templates) Autosaves(parent string) *TemplateAutosaves {
	return &TemplateAutosaves{
		client: api.client,
		parent: parent,
	}
}

// List returns a ListTemplateAutosaves builder to query autosaves for the parent template.
func (api *TemplateAutosaves) List() *ListTemplateAutosaves {
	return &ListTemplateAutosaves{
		endpoint:  "/wp/v2/templates/" + api.parent + "/autosaves",
		client:    api.client,
		arguments: make(map[string]string),
	}
}

// Retrieve returns a RetrieveTemplateAutosave builder for a specific autosave revision ID.
func (api *TemplateAutosaves) Retrieve(id int) *RetrieveTemplateAutosave {
	return &RetrieveTemplateAutosave{
		endpoint:  "/wp/v2/templates/" + api.parent + "/autosaves/" + strconv.Itoa(id),
		client:    api.client,
		arguments: make(map[string]string),
	}
}

// Create returns a CreateTemplateAutosave builder to create or update an autosave.
func (api *TemplateAutosaves) Create() *CreateTemplateAutosave {
	return &CreateTemplateAutosave{
		endpoint: "/wp/v2/templates/" + api.parent + "/autosaves",
		client:   api.client,
		data:     TemplateData{},
	}
}

// ListTemplateAutosaves handles querying template autosaves.
type ListTemplateAutosaves struct {
	endpoint  string
	client    *RestClient
	arguments map[string]string
}

// Context sets the context parameter (e.g. "view", "edit", "embed").
func (api *ListTemplateAutosaves) Context(ctx string) *ListTemplateAutosaves {
	api.arguments["context"] = ctx
	return api
}

// ContextView sets the context parameter to "view".
func (api *ListTemplateAutosaves) ContextView() *ListTemplateAutosaves {
	return api.Context("view")
}

// ContextEdit sets the context parameter to "edit".
func (api *ListTemplateAutosaves) ContextEdit() *ListTemplateAutosaves {
	return api.Context("edit")
}

// ContextEmbed sets the context parameter to "embed".
func (api *ListTemplateAutosaves) ContextEmbed() *ListTemplateAutosaves {
	return api.Context("embed")
}

// Fields limits the response to specific fields.
func (api *ListTemplateAutosaves) Fields(fields ...string) *ListTemplateAutosaves {
	api.arguments["_fields"] = strings.Join(fields, ",")
	return api
}

// Do executes the list request.
func (api *ListTemplateAutosaves) Do() (autosaves []TemplateRevision, err error) {
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

// RetrieveTemplateAutosave handles retrieving a single template autosave.
type RetrieveTemplateAutosave struct {
	endpoint  string
	client    *RestClient
	arguments map[string]string
}

// Context sets the context parameter (e.g. "view", "edit", "embed").
func (api *RetrieveTemplateAutosave) Context(ctx string) *RetrieveTemplateAutosave {
	api.arguments["context"] = ctx
	return api
}

// ContextView sets the context parameter to "view".
func (api *RetrieveTemplateAutosave) ContextView() *RetrieveTemplateAutosave {
	return api.Context("view")
}

// ContextEdit sets the context parameter to "edit".
func (api *RetrieveTemplateAutosave) ContextEdit() *RetrieveTemplateAutosave {
	return api.Context("edit")
}

// ContextEmbed sets the context parameter to "embed".
func (api *RetrieveTemplateAutosave) ContextEmbed() *RetrieveTemplateAutosave {
	return api.Context("embed")
}

// Fields limits the response to specific fields.
func (api *RetrieveTemplateAutosave) Fields(fields ...string) *RetrieveTemplateAutosave {
	api.arguments["_fields"] = strings.Join(fields, ",")
	return api
}

// Do executes the retrieve request.
func (api *RetrieveTemplateAutosave) Do() (autosave *TemplateRevision, err error) {
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

// CreateTemplateAutosave handles creating a template autosave.
type CreateTemplateAutosave struct {
	endpoint string
	client   *RestClient
	data     TemplateData
}

// Slug sets the unique slug identifying the template.
func (api *CreateTemplateAutosave) Slug(slug string) *CreateTemplateAutosave {
	api.data.Slug = slug
	return api
}

// Theme sets the theme identifier for the template.
func (api *CreateTemplateAutosave) Theme(theme string) *CreateTemplateAutosave {
	api.data.Theme = theme
	return api
}

// Type sets the type of template.
func (api *CreateTemplateAutosave) Type(typ string) *CreateTemplateAutosave {
	api.data.Type = typ
	return api
}

// Content sets the content of template.
func (api *CreateTemplateAutosave) Content(content any) *CreateTemplateAutosave {
	api.data.Content = content
	return api
}

// Title sets the title of template.
func (api *CreateTemplateAutosave) Title(title any) *CreateTemplateAutosave {
	api.data.Title = title
	return api
}

// Description sets the description of template.
func (api *CreateTemplateAutosave) Description(desc string) *CreateTemplateAutosave {
	api.data.Description = desc
	return api
}

// Status sets the status of template.
func (api *CreateTemplateAutosave) Status(status string) *CreateTemplateAutosave {
	api.data.Status = status
	return api
}

// Author sets the ID for the author of the template.
func (api *CreateTemplateAutosave) Author(author int) *CreateTemplateAutosave {
	api.data.Author = author
	return api
}

// WithData populates the request body from a TemplateData struct.
func (api *CreateTemplateAutosave) WithData(data TemplateData) *CreateTemplateAutosave {
	api.data = data
	return api
}

// Do executes the create autosave request.
func (api *CreateTemplateAutosave) Do() (autosave *TemplateRevision, err error) {
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
