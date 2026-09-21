package gowprest

import (
	"encoding/json"
	"strconv"
	"strings"
)

// TemplatePart represents a WordPress block theme template part record (/wp/v2/template-parts).
type TemplatePart struct {
	ID             string          `json:"id,omitempty"`
	Slug           string          `json:"slug,omitempty"`
	Theme          string          `json:"theme,omitempty"`
	Type           string          `json:"type,omitempty"`
	Source         string          `json:"source,omitempty"`
	Origin         *string         `json:"origin,omitempty"`
	Content        *TemplateObject `json:"content,omitempty"`
	Title          *TemplateObject `json:"title,omitempty"`
	Description    string          `json:"description,omitempty"`
	Status         string          `json:"status,omitempty"`
	WPID           int             `json:"wp_id,omitempty"`
	HasThemeFile   bool            `json:"has_theme_file,omitempty"`
	Author         int             `json:"author,omitempty"`
	AuthorText     string          `json:"author_text,omitempty"`
	Area           string          `json:"area,omitempty"`
	Modified       *string         `json:"modified,omitempty"`
	Date           *string         `json:"date,omitempty"`
	OriginalSource string          `json:"original_source,omitempty"`
	Links          map[string]any  `json:"_links,omitempty"`
}

// TitleString returns the template part's title string from either Raw or Rendered.
func (tp *TemplatePart) TitleString() string {
	if tp.Title == nil {
		return ""
	}
	if tp.Title.Raw != "" {
		return tp.Title.Raw
	}
	return tp.Title.Rendered
}

// ContentString returns the template part's content string from either Raw or Rendered.
func (tp *TemplatePart) ContentString() string {
	if tp.Content == nil {
		return ""
	}
	if tp.Content.Raw != "" {
		return tp.Content.Raw
	}
	return tp.Content.Rendered
}

// TemplatePartData represents the request body for creating or updating a template part.
type TemplatePartData struct {
	Slug        string `json:"slug,omitempty"`
	Theme       string `json:"theme,omitempty"`
	Type        string `json:"type,omitempty"`
	Content     any    `json:"content,omitempty"`
	Title       any    `json:"title,omitempty"`
	Description string `json:"description,omitempty"`
	Status      string `json:"status,omitempty"`
	Author      int    `json:"author,omitempty"`
	Area        string `json:"area,omitempty"`
}

// TemplateParts provides access to WordPress template part endpoints (/wp/v2/template-parts).
type TemplateParts struct {
	client *RestClient
}

// TemplateParts returns a TemplateParts service instance.
func (c *RestClient) TemplateParts() *TemplateParts {
	return &TemplateParts{client: c}
}

// TemplatePart is an ergonomic alias for TemplateParts.
func (c *RestClient) TemplatePart() *TemplateParts {
	return c.TemplateParts()
}

// List returns a ListTemplateParts builder to query template parts.
func (api *TemplateParts) List() *ListTemplateParts {
	return &ListTemplateParts{
		endpoint:  "/wp/v2/template-parts",
		client:    api.client,
		arguments: make(map[string]string),
	}
}

// Retrieve returns a RetrieveTemplatePart builder for a specific template part ID.
func (api *TemplateParts) Retrieve(id string) *RetrieveTemplatePart {
	return &RetrieveTemplatePart{
		endpoint:  "/wp/v2/template-parts/" + id,
		client:    api.client,
		arguments: make(map[string]string),
	}
}

// Create returns a CreateTemplatePart builder to create a new template part.
func (api *TemplateParts) Create() *CreateTemplatePart {
	return &CreateTemplatePart{
		endpoint: "/wp/v2/template-parts",
		client:   api.client,
		data:     TemplatePartData{},
	}
}

// Update returns an UpdateTemplatePart builder to modify an existing template part.
func (api *TemplateParts) Update(id string) *UpdateTemplatePart {
	return &UpdateTemplatePart{
		endpoint: "/wp/v2/template-parts/" + id,
		client:   api.client,
		data:     TemplatePartData{},
	}
}

// Delete returns a DeleteTemplatePart builder to remove a template part.
func (api *TemplateParts) Delete(id string) *DeleteTemplatePart {
	return &DeleteTemplatePart{
		endpoint: "/wp/v2/template-parts/" + id,
		client:   api.client,
	}
}

// ListTemplateParts handles querying template parts.
type ListTemplateParts struct {
	endpoint  string
	client    *RestClient
	arguments map[string]string
}

// Context sets the context parameter (e.g. "view", "edit", "embed").
func (api *ListTemplateParts) Context(ctx string) *ListTemplateParts {
	api.arguments["context"] = ctx
	return api
}

// ContextView sets the context parameter to "view".
func (api *ListTemplateParts) ContextView() *ListTemplateParts {
	return api.Context("view")
}

// ContextEdit sets the context parameter to "edit".
func (api *ListTemplateParts) ContextEdit() *ListTemplateParts {
	return api.Context("edit")
}

// ContextEmbed sets the context parameter to "embed".
func (api *ListTemplateParts) ContextEmbed() *ListTemplateParts {
	return api.Context("embed")
}

// WPID filters template parts by post ID.
func (api *ListTemplateParts) WPID(wpid int) *ListTemplateParts {
	api.arguments["wp_id"] = strconv.Itoa(wpid)
	return api
}

// Area filters template parts by template part area (e.g. "header", "footer", "sidebar").
func (api *ListTemplateParts) Area(area string) *ListTemplateParts {
	api.arguments["area"] = area
	return api
}

// PostType filters template parts by post type.
func (api *ListTemplateParts) PostType(postType string) *ListTemplateParts {
	api.arguments["post_type"] = postType
	return api
}

// Fields limits the response to specific fields.
func (api *ListTemplateParts) Fields(fields ...string) *ListTemplateParts {
	api.arguments["_fields"] = strings.Join(fields, ",")
	return api
}

// Do executes the list request.
func (api *ListTemplateParts) Do() (parts []TemplatePart, err error) {
	endpoint := api.client.endpoint + api.endpoint

	restyClient := api.client.httpClient.R()
	if api.client.auth.Username != "" && api.client.auth.Password != "" {
		restyClient = restyClient.SetBasicAuth(api.client.auth.Username, api.client.auth.Password)
	}

	resp, err := restyClient.
		SetHeader("Accept", "application/json").
		SetResult(&parts).
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

	return parts, nil
}

// RetrieveTemplatePart handles retrieving a single template part.
type RetrieveTemplatePart struct {
	endpoint  string
	client    *RestClient
	arguments map[string]string
}

// Context sets the context parameter (e.g. "view", "edit", "embed").
func (api *RetrieveTemplatePart) Context(ctx string) *RetrieveTemplatePart {
	api.arguments["context"] = ctx
	return api
}

// ContextView sets the context parameter to "view".
func (api *RetrieveTemplatePart) ContextView() *RetrieveTemplatePart {
	return api.Context("view")
}

// ContextEdit sets the context parameter to "edit".
func (api *RetrieveTemplatePart) ContextEdit() *RetrieveTemplatePart {
	return api.Context("edit")
}

// ContextEmbed sets the context parameter to "embed".
func (api *RetrieveTemplatePart) ContextEmbed() *RetrieveTemplatePart {
	return api.Context("embed")
}

// Fields limits the response to specific fields.
func (api *RetrieveTemplatePart) Fields(fields ...string) *RetrieveTemplatePart {
	api.arguments["_fields"] = strings.Join(fields, ",")
	return api
}

// Do executes the retrieve request.
func (api *RetrieveTemplatePart) Do() (part *TemplatePart, err error) {
	endpoint := api.client.endpoint + api.endpoint

	restyClient := api.client.httpClient.R()
	if api.client.auth.Username != "" && api.client.auth.Password != "" {
		restyClient = restyClient.SetBasicAuth(api.client.auth.Username, api.client.auth.Password)
	}

	resp, err := restyClient.
		SetHeader("Accept", "application/json").
		SetResult(&part).
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

	return part, nil
}

// CreateTemplatePart handles creating a new template part.
type CreateTemplatePart struct {
	endpoint string
	client   *RestClient
	data     TemplatePartData
}

// Slug sets the unique slug identifying the template part (required).
func (api *CreateTemplatePart) Slug(slug string) *CreateTemplatePart {
	api.data.Slug = slug
	return api
}

// Theme sets the theme identifier for the template part.
func (api *CreateTemplatePart) Theme(theme string) *CreateTemplatePart {
	api.data.Theme = theme
	return api
}

// Type sets the type of template part.
func (api *CreateTemplatePart) Type(typ string) *CreateTemplatePart {
	api.data.Type = typ
	return api
}

// Content sets the content of template part (string or object).
func (api *CreateTemplatePart) Content(content any) *CreateTemplatePart {
	api.data.Content = content
	return api
}

// Title sets the title of template part (string or object).
func (api *CreateTemplatePart) Title(title any) *CreateTemplatePart {
	api.data.Title = title
	return api
}

// Description sets the description of template part.
func (api *CreateTemplatePart) Description(desc string) *CreateTemplatePart {
	api.data.Description = desc
	return api
}

// Status sets the status of template part ("publish", "future", "draft", "pending", "private").
func (api *CreateTemplatePart) Status(status string) *CreateTemplatePart {
	api.data.Status = status
	return api
}

// StatusPublish sets status to "publish".
func (api *CreateTemplatePart) StatusPublish() *CreateTemplatePart {
	return api.Status("publish")
}

// StatusDraft sets status to "draft".
func (api *CreateTemplatePart) StatusDraft() *CreateTemplatePart {
	return api.Status("draft")
}

// Author sets the ID for the author of the template part.
func (api *CreateTemplatePart) Author(author int) *CreateTemplatePart {
	api.data.Author = author
	return api
}

// Area sets where the template part is intended for use (header, footer, sidebar, etc.).
func (api *CreateTemplatePart) Area(area string) *CreateTemplatePart {
	api.data.Area = area
	return api
}

// WithData populates the request body from a TemplatePartData struct.
func (api *CreateTemplatePart) WithData(data TemplatePartData) *CreateTemplatePart {
	api.data = data
	return api
}

// Do executes the create request.
func (api *CreateTemplatePart) Do() (part *TemplatePart, err error) {
	endpoint := api.client.endpoint + api.endpoint

	restyClient := api.client.httpClient.R()
	if api.client.auth.Username != "" && api.client.auth.Password != "" {
		restyClient = restyClient.SetBasicAuth(api.client.auth.Username, api.client.auth.Password)
	}

	resp, err := restyClient.
		SetHeader("Content-Type", "application/json").
		SetHeader("Accept", "application/json").
		SetBody(api.data).
		SetResult(&part).
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

	return part, nil
}

// UpdateTemplatePart handles modifying an existing template part.
type UpdateTemplatePart struct {
	endpoint string
	client   *RestClient
	data     TemplatePartData
}

// Slug sets the unique slug identifying the template part.
func (api *UpdateTemplatePart) Slug(slug string) *UpdateTemplatePart {
	api.data.Slug = slug
	return api
}

// Theme sets the theme identifier for the template part.
func (api *UpdateTemplatePart) Theme(theme string) *UpdateTemplatePart {
	api.data.Theme = theme
	return api
}

// Type sets the type of template part.
func (api *UpdateTemplatePart) Type(typ string) *UpdateTemplatePart {
	api.data.Type = typ
	return api
}

// Content sets the content of template part (string or object).
func (api *UpdateTemplatePart) Content(content any) *UpdateTemplatePart {
	api.data.Content = content
	return api
}

// Title sets the title of template part (string or object).
func (api *UpdateTemplatePart) Title(title any) *UpdateTemplatePart {
	api.data.Title = title
	return api
}

// Description sets the description of template part.
func (api *UpdateTemplatePart) Description(desc string) *UpdateTemplatePart {
	api.data.Description = desc
	return api
}

// Status sets the status of template part ("publish", "future", "draft", "pending", "private").
func (api *UpdateTemplatePart) Status(status string) *UpdateTemplatePart {
	api.data.Status = status
	return api
}

// StatusPublish sets status to "publish".
func (api *UpdateTemplatePart) StatusPublish() *UpdateTemplatePart {
	return api.Status("publish")
}

// StatusDraft sets status to "draft".
func (api *UpdateTemplatePart) StatusDraft() *UpdateTemplatePart {
	return api.Status("draft")
}

// Author sets the ID for the author of the template part.
func (api *UpdateTemplatePart) Author(author int) *UpdateTemplatePart {
	api.data.Author = author
	return api
}

// Area sets where the template part is intended for use (header, footer, sidebar, etc.).
func (api *UpdateTemplatePart) Area(area string) *UpdateTemplatePart {
	api.data.Area = area
	return api
}

// WithData populates the request body from a TemplatePartData struct.
func (api *UpdateTemplatePart) WithData(data TemplatePartData) *UpdateTemplatePart {
	api.data = data
	return api
}

// Do executes the update request.
func (api *UpdateTemplatePart) Do() (part *TemplatePart, err error) {
	endpoint := api.client.endpoint + api.endpoint

	restyClient := api.client.httpClient.R()
	if api.client.auth.Username != "" && api.client.auth.Password != "" {
		restyClient = restyClient.SetBasicAuth(api.client.auth.Username, api.client.auth.Password)
	}

	resp, err := restyClient.
		SetHeader("Content-Type", "application/json").
		SetHeader("Accept", "application/json").
		SetBody(api.data).
		SetResult(&part).
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

	return part, nil
}

type deleteTemplatePartEnvelope struct {
	TemplatePart
	Deleted  bool          `json:"deleted"`
	Previous *TemplatePart `json:"previous"`
}

// DeleteTemplatePart handles deleting a template part.
type DeleteTemplatePart struct {
	endpoint string
	client   *RestClient
	force    bool
}

// Force sets whether to bypass Trash and force permanent deletion.
func (api *DeleteTemplatePart) Force() *DeleteTemplatePart {
	api.force = true
	return api
}

// Do executes the delete request.
func (api *DeleteTemplatePart) Do() (part *TemplatePart, err error) {
	endpoint := api.client.endpoint + api.endpoint
	var env deleteTemplatePartEnvelope

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

	if env.Deleted && env.Previous != nil && env.Previous.ID != "" {
		return env.Previous, nil
	}

	return &env.TemplatePart, nil
}
