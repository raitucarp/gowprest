package gowprest

import (
	"encoding/json"
	"strconv"
	"strings"
)

// TemplateObject represents content or title objects in templates,
// handling both structured object {"raw": "...", "rendered": "..."} and string forms.
type TemplateObject struct {
	Raw          string `json:"raw,omitempty"`
	Rendered     string `json:"rendered,omitempty"`
	BlockVersion int    `json:"block_version,omitempty"`
}

// UnmarshalJSON implements custom unmarshaling for TemplateObject.
func (to *TemplateObject) UnmarshalJSON(b []byte) error {
	trimmed := strings.TrimSpace(string(b))
	if trimmed == "" || trimmed == "null" || trimmed == "[]" {
		return nil
	}
	if len(trimmed) > 0 && trimmed[0] == '[' {
		return nil
	}
	if len(trimmed) > 0 && trimmed[0] == '"' {
		var s string
		if err := json.Unmarshal(b, &s); err != nil {
			return err
		}
		to.Raw = s
		to.Rendered = s
		return nil
	}
	type Alias TemplateObject
	var aux Alias
	if err := json.Unmarshal(b, &aux); err != nil {
		return err
	}
	*to = TemplateObject(aux)
	return nil
}

// Template represents a WordPress block theme template record (/wp/v2/templates).
type Template struct {
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
	IsCustom       bool            `json:"is_custom,omitempty"`
	Author         int             `json:"author,omitempty"`
	AuthorText     string          `json:"author_text,omitempty"`
	Modified       *string         `json:"modified,omitempty"`
	Date           *string         `json:"date,omitempty"`
	OriginalSource string          `json:"original_source,omitempty"`
	Links          map[string]any  `json:"_links,omitempty"`
}

// TemplateData represents the request body for creating or updating a template.
type TemplateData struct {
	Slug        string `json:"slug,omitempty"`
	Theme       string `json:"theme,omitempty"`
	Type        string `json:"type,omitempty"`
	Content     any    `json:"content,omitempty"`
	Title       any    `json:"title,omitempty"`
	Description string `json:"description,omitempty"`
	Status      string `json:"status,omitempty"`
	Author      int    `json:"author,omitempty"`
}

// Templates provides access to WordPress template endpoints (/wp/v2/templates).
type Templates struct {
	client *RestClient
}

// Templates returns a Templates service instance.
func (c *RestClient) Templates() *Templates {
	return &Templates{client: c}
}

// Template is an ergonomic alias for Templates.
func (c *RestClient) Template() *Templates {
	return c.Templates()
}

// List returns a ListTemplates builder to query registered templates.
func (api *Templates) List() *ListTemplates {
	return &ListTemplates{
		endpoint:  "/wp/v2/templates",
		client:    api.client,
		arguments: make(map[string]string),
	}
}

// Retrieve returns a RetrieveTemplate builder for a specific template ID.
func (api *Templates) Retrieve(id string) *RetrieveTemplate {
	return &RetrieveTemplate{
		endpoint:  "/wp/v2/templates/" + id,
		client:    api.client,
		arguments: make(map[string]string),
	}
}

// Create returns a CreateTemplate builder to create a new template.
func (api *Templates) Create() *CreateTemplate {
	return &CreateTemplate{
		endpoint: "/wp/v2/templates",
		client:   api.client,
		data:     TemplateData{},
	}
}

// Update returns an UpdateTemplate builder to modify an existing template.
func (api *Templates) Update(id string) *UpdateTemplate {
	return &UpdateTemplate{
		endpoint: "/wp/v2/templates/" + id,
		client:   api.client,
		data:     TemplateData{},
	}
}

// Delete returns a DeleteTemplate builder to remove a template.
func (api *Templates) Delete(id string) *DeleteTemplate {
	return &DeleteTemplate{
		endpoint: "/wp/v2/templates/" + id,
		client:   api.client,
	}
}

// ListTemplates handles querying templates.
type ListTemplates struct {
	endpoint  string
	client    *RestClient
	arguments map[string]string
}

// Context sets the context parameter (e.g. "view", "edit", "embed").
func (api *ListTemplates) Context(ctx string) *ListTemplates {
	api.arguments["context"] = ctx
	return api
}

// ContextView sets the context parameter to "view".
func (api *ListTemplates) ContextView() *ListTemplates {
	return api.Context("view")
}

// ContextEdit sets the context parameter to "edit".
func (api *ListTemplates) ContextEdit() *ListTemplates {
	return api.Context("edit")
}

// ContextEmbed sets the context parameter to "embed".
func (api *ListTemplates) ContextEmbed() *ListTemplates {
	return api.Context("embed")
}

// WPID filters templates by post ID.
func (api *ListTemplates) WPID(wpid int) *ListTemplates {
	api.arguments["wp_id"] = strconv.Itoa(wpid)
	return api
}

// Area filters templates by template part area.
func (api *ListTemplates) Area(area string) *ListTemplates {
	api.arguments["area"] = area
	return api
}

// PostType filters templates by post type.
func (api *ListTemplates) PostType(postType string) *ListTemplates {
	api.arguments["post_type"] = postType
	return api
}

// Fields limits the response to specific fields.
func (api *ListTemplates) Fields(fields ...string) *ListTemplates {
	api.arguments["_fields"] = strings.Join(fields, ",")
	return api
}

// Do executes the list request.
func (api *ListTemplates) Do() (templates []Template, err error) {
	endpoint := api.client.endpoint + api.endpoint

	restyClient := api.client.httpClient.R()
	if api.client.auth.Username != "" && api.client.auth.Password != "" {
		restyClient = restyClient.SetBasicAuth(api.client.auth.Username, api.client.auth.Password)
	}

	resp, err := restyClient.
		SetHeader("Accept", "application/json").
		SetResult(&templates).
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

	return templates, nil
}

// RetrieveTemplate handles retrieving a single template.
type RetrieveTemplate struct {
	endpoint  string
	client    *RestClient
	arguments map[string]string
}

// Context sets the context parameter (e.g. "view", "edit", "embed").
func (api *RetrieveTemplate) Context(ctx string) *RetrieveTemplate {
	api.arguments["context"] = ctx
	return api
}

// ContextView sets the context parameter to "view".
func (api *RetrieveTemplate) ContextView() *RetrieveTemplate {
	return api.Context("view")
}

// ContextEdit sets the context parameter to "edit".
func (api *RetrieveTemplate) ContextEdit() *RetrieveTemplate {
	return api.Context("edit")
}

// ContextEmbed sets the context parameter to "embed".
func (api *RetrieveTemplate) ContextEmbed() *RetrieveTemplate {
	return api.Context("embed")
}

// Fields limits the response to specific fields.
func (api *RetrieveTemplate) Fields(fields ...string) *RetrieveTemplate {
	api.arguments["_fields"] = strings.Join(fields, ",")
	return api
}

// Do executes the retrieve request.
func (api *RetrieveTemplate) Do() (template *Template, err error) {
	endpoint := api.client.endpoint + api.endpoint

	restyClient := api.client.httpClient.R()
	if api.client.auth.Username != "" && api.client.auth.Password != "" {
		restyClient = restyClient.SetBasicAuth(api.client.auth.Username, api.client.auth.Password)
	}

	resp, err := restyClient.
		SetHeader("Accept", "application/json").
		SetResult(&template).
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

	return template, nil
}

// CreateTemplate handles creating a new template.
type CreateTemplate struct {
	endpoint string
	client   *RestClient
	data     TemplateData
}

// Slug sets the unique slug identifying the template (required).
func (api *CreateTemplate) Slug(slug string) *CreateTemplate {
	api.data.Slug = slug
	return api
}

// Theme sets the theme identifier for the template.
func (api *CreateTemplate) Theme(theme string) *CreateTemplate {
	api.data.Theme = theme
	return api
}

// Type sets the type of template.
func (api *CreateTemplate) Type(typ string) *CreateTemplate {
	api.data.Type = typ
	return api
}

// Content sets the content of template (string or object).
func (api *CreateTemplate) Content(content any) *CreateTemplate {
	api.data.Content = content
	return api
}

// Title sets the title of template (string or object).
func (api *CreateTemplate) Title(title any) *CreateTemplate {
	api.data.Title = title
	return api
}

// Description sets the description of template.
func (api *CreateTemplate) Description(desc string) *CreateTemplate {
	api.data.Description = desc
	return api
}

// Status sets the status of template ("publish", "future", "draft", "pending", "private").
func (api *CreateTemplate) Status(status string) *CreateTemplate {
	api.data.Status = status
	return api
}

// StatusPublish sets status to "publish".
func (api *CreateTemplate) StatusPublish() *CreateTemplate {
	return api.Status("publish")
}

// StatusDraft sets status to "draft".
func (api *CreateTemplate) StatusDraft() *CreateTemplate {
	return api.Status("draft")
}

// Author sets the ID for the author of the template.
func (api *CreateTemplate) Author(author int) *CreateTemplate {
	api.data.Author = author
	return api
}

// WithData populates the request body from a TemplateData struct.
func (api *CreateTemplate) WithData(data TemplateData) *CreateTemplate {
	api.data = data
	return api
}

// Do executes the create request.
func (api *CreateTemplate) Do() (template *Template, err error) {
	endpoint := api.client.endpoint + api.endpoint

	restyClient := api.client.httpClient.R()
	if api.client.auth.Username != "" && api.client.auth.Password != "" {
		restyClient = restyClient.SetBasicAuth(api.client.auth.Username, api.client.auth.Password)
	}

	resp, err := restyClient.
		SetHeader("Content-Type", "application/json").
		SetHeader("Accept", "application/json").
		SetBody(api.data).
		SetResult(&template).
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

	return template, nil
}

// UpdateTemplate handles modifying an existing template.
type UpdateTemplate struct {
	endpoint string
	client   *RestClient
	data     TemplateData
}

// Slug sets the unique slug identifying the template.
func (api *UpdateTemplate) Slug(slug string) *UpdateTemplate {
	api.data.Slug = slug
	return api
}

// Theme sets the theme identifier for the template.
func (api *UpdateTemplate) Theme(theme string) *UpdateTemplate {
	api.data.Theme = theme
	return api
}

// Type sets the type of template.
func (api *UpdateTemplate) Type(typ string) *UpdateTemplate {
	api.data.Type = typ
	return api
}

// Content sets the content of template (string or object).
func (api *UpdateTemplate) Content(content any) *UpdateTemplate {
	api.data.Content = content
	return api
}

// Title sets the title of template (string or object).
func (api *UpdateTemplate) Title(title any) *UpdateTemplate {
	api.data.Title = title
	return api
}

// Description sets the description of template.
func (api *UpdateTemplate) Description(desc string) *UpdateTemplate {
	api.data.Description = desc
	return api
}

// Status sets the status of template ("publish", "future", "draft", "pending", "private").
func (api *UpdateTemplate) Status(status string) *UpdateTemplate {
	api.data.Status = status
	return api
}

// StatusPublish sets status to "publish".
func (api *UpdateTemplate) StatusPublish() *UpdateTemplate {
	return api.Status("publish")
}

// StatusDraft sets status to "draft".
func (api *UpdateTemplate) StatusDraft() *UpdateTemplate {
	return api.Status("draft")
}

// Author sets the ID for the author of the template.
func (api *UpdateTemplate) Author(author int) *UpdateTemplate {
	api.data.Author = author
	return api
}

// WithData populates the request body from a TemplateData struct.
func (api *UpdateTemplate) WithData(data TemplateData) *UpdateTemplate {
	api.data = data
	return api
}

// Do executes the update request.
func (api *UpdateTemplate) Do() (template *Template, err error) {
	endpoint := api.client.endpoint + api.endpoint

	restyClient := api.client.httpClient.R()
	if api.client.auth.Username != "" && api.client.auth.Password != "" {
		restyClient = restyClient.SetBasicAuth(api.client.auth.Username, api.client.auth.Password)
	}

	resp, err := restyClient.
		SetHeader("Content-Type", "application/json").
		SetHeader("Accept", "application/json").
		SetBody(api.data).
		SetResult(&template).
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

	return template, nil
}

type deleteTemplateEnvelope struct {
	Template
	Deleted  bool      `json:"deleted"`
	Previous *Template `json:"previous"`
}

// DeleteTemplate handles deleting a template.
type DeleteTemplate struct {
	endpoint string
	client   *RestClient
	force    bool
}

// Force sets whether to bypass Trash and force permanent deletion.
func (api *DeleteTemplate) Force() *DeleteTemplate {
	api.force = true
	return api
}

// Do executes the delete request.
func (api *DeleteTemplate) Do() (template *Template, err error) {
	endpoint := api.client.endpoint + api.endpoint
	var env deleteTemplateEnvelope

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

	return &env.Template, nil
}
