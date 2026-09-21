package gowprest

import (
	"encoding/json"
	"strings"
)

// WidgetInstance represents the instance settings of a WordPress widget.
type WidgetInstance struct {
	Encoded string         `json:"encoded,omitempty"`
	Hash    string         `json:"hash,omitempty"`
	Raw     map[string]any `json:"raw,omitempty"`
}

// Widget represents a WordPress widget record (/wp/v2/widgets).
type Widget struct {
	ID           string          `json:"id,omitempty"`
	IDBase       string          `json:"id_base,omitempty"`
	Sidebar      string          `json:"sidebar,omitempty"`
	Rendered     string          `json:"rendered,omitempty"`
	RenderedForm string          `json:"rendered_form,omitempty"`
	Instance     *WidgetInstance `json:"instance,omitempty"`
	FormData     string          `json:"form_data,omitempty"`
	Links        map[string]any  `json:"_links,omitempty"`
}

// WidgetData represents the request body for creating or updating a widget.
type WidgetData struct {
	ID       string `json:"id,omitempty"`
	IDBase   string `json:"id_base,omitempty"`
	Sidebar  string `json:"sidebar,omitempty"`
	Instance any    `json:"instance,omitempty"`
	FormData string `json:"form_data,omitempty"`
}

// Widgets provides access to WordPress widgets endpoints (/wp/v2/widgets).
type Widgets struct {
	client *RestClient
}

// Widgets returns a Widgets service instance.
func (c *RestClient) Widgets() *Widgets {
	return &Widgets{client: c}
}

// Widget is an ergonomic alias for Widgets.
func (c *RestClient) Widget() *Widgets {
	return c.Widgets()
}

// List returns a ListWidgets builder to query registered widgets.
func (api *Widgets) List() *ListWidgets {
	return &ListWidgets{
		endpoint:  "/wp/v2/widgets",
		client:    api.client,
		arguments: make(map[string]string),
	}
}

// Retrieve returns a RetrieveWidget builder for a specific widget ID.
func (api *Widgets) Retrieve(id string) *RetrieveWidget {
	return &RetrieveWidget{
		endpoint:  "/wp/v2/widgets/" + id,
		client:    api.client,
		arguments: make(map[string]string),
	}
}

// Create returns a CreateWidget builder to create a new widget.
func (api *Widgets) Create() *CreateWidget {
	return &CreateWidget{
		endpoint: "/wp/v2/widgets",
		client:   api.client,
		data:     WidgetData{},
	}
}

// Update returns an UpdateWidget builder to modify an existing widget.
func (api *Widgets) Update(id string) *UpdateWidget {
	return &UpdateWidget{
		endpoint: "/wp/v2/widgets/" + id,
		client:   api.client,
		data:     WidgetData{},
	}
}

// Delete returns a DeleteWidget builder to remove or deactivate a widget.
func (api *Widgets) Delete(id string) *DeleteWidget {
	return &DeleteWidget{
		endpoint: "/wp/v2/widgets/" + id,
		client:   api.client,
	}
}

// ListWidgets handles querying widgets.
type ListWidgets struct {
	endpoint  string
	client    *RestClient
	arguments map[string]string
}

// Context sets the context parameter (e.g. "view", "edit", "embed").
func (api *ListWidgets) Context(ctx string) *ListWidgets {
	api.arguments["context"] = ctx
	return api
}

// ContextView sets the context parameter to "view".
func (api *ListWidgets) ContextView() *ListWidgets {
	return api.Context("view")
}

// ContextEdit sets the context parameter to "edit".
func (api *ListWidgets) ContextEdit() *ListWidgets {
	return api.Context("edit")
}

// ContextEmbed sets the context parameter to "embed".
func (api *ListWidgets) ContextEmbed() *ListWidgets {
	return api.Context("embed")
}

// Sidebar filters widgets by sidebar ID.
func (api *ListWidgets) Sidebar(sidebar string) *ListWidgets {
	api.arguments["sidebar"] = sidebar
	return api
}

// Fields limits the response to specific fields.
func (api *ListWidgets) Fields(fields ...string) *ListWidgets {
	api.arguments["_fields"] = strings.Join(fields, ",")
	return api
}

// Do executes the list request.
func (api *ListWidgets) Do() (widgets []Widget, err error) {
	endpoint := api.client.endpoint + api.endpoint

	restyClient := api.client.httpClient.R()
	if api.client.auth.Username != "" && api.client.auth.Password != "" {
		restyClient = restyClient.SetBasicAuth(api.client.auth.Username, api.client.auth.Password)
	}

	resp, err := restyClient.
		SetHeader("Accept", "application/json").
		SetResult(&widgets).
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

	return widgets, nil
}

// RetrieveWidget handles retrieving a single widget.
type RetrieveWidget struct {
	endpoint  string
	client    *RestClient
	arguments map[string]string
}

// Context sets the context parameter (e.g. "view", "edit", "embed").
func (api *RetrieveWidget) Context(ctx string) *RetrieveWidget {
	api.arguments["context"] = ctx
	return api
}

// ContextView sets the context parameter to "view".
func (api *RetrieveWidget) ContextView() *RetrieveWidget {
	return api.Context("view")
}

// ContextEdit sets the context parameter to "edit".
func (api *RetrieveWidget) ContextEdit() *RetrieveWidget {
	return api.Context("edit")
}

// ContextEmbed sets the context parameter to "embed".
func (api *RetrieveWidget) ContextEmbed() *RetrieveWidget {
	return api.Context("embed")
}

// Fields limits the response to specific fields.
func (api *RetrieveWidget) Fields(fields ...string) *RetrieveWidget {
	api.arguments["_fields"] = strings.Join(fields, ",")
	return api
}

// Do executes the retrieve request.
func (api *RetrieveWidget) Do() (widget *Widget, err error) {
	endpoint := api.client.endpoint + api.endpoint

	restyClient := api.client.httpClient.R()
	if api.client.auth.Username != "" && api.client.auth.Password != "" {
		restyClient = restyClient.SetBasicAuth(api.client.auth.Username, api.client.auth.Password)
	}

	resp, err := restyClient.
		SetHeader("Accept", "application/json").
		SetResult(&widget).
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

	return widget, nil
}

// CreateWidget handles creating a new widget.
type CreateWidget struct {
	endpoint string
	client   *RestClient
	data     WidgetData
}

// ID sets the unique identifier for the widget.
func (api *CreateWidget) ID(id string) *CreateWidget {
	api.data.ID = id
	return api
}

// IDBase sets the type of the widget (corresponds to ID in widget-types endpoint).
func (api *CreateWidget) IDBase(idBase string) *CreateWidget {
	api.data.IDBase = idBase
	return api
}

// Sidebar sets the sidebar the widget belongs to (defaults to "wp_inactive_widgets").
func (api *CreateWidget) Sidebar(sidebar string) *CreateWidget {
	api.data.Sidebar = sidebar
	return api
}

// Instance sets instance settings of the widget (map, struct, or raw data).
func (api *CreateWidget) Instance(instance any) *CreateWidget {
	api.data.Instance = instance
	return api
}

// FormData sets URL-encoded form data for widgets that do not support instance.
func (api *CreateWidget) FormData(formData string) *CreateWidget {
	api.data.FormData = formData
	return api
}

// WithData populates the request body from a WidgetData struct.
func (api *CreateWidget) WithData(data WidgetData) *CreateWidget {
	api.data = data
	return api
}

// Do executes the create request.
func (api *CreateWidget) Do() (widget *Widget, err error) {
	endpoint := api.client.endpoint + api.endpoint

	restyClient := api.client.httpClient.R()
	if api.client.auth.Username != "" && api.client.auth.Password != "" {
		restyClient = restyClient.SetBasicAuth(api.client.auth.Username, api.client.auth.Password)
	}

	resp, err := restyClient.
		SetHeader("Content-Type", "application/json").
		SetHeader("Accept", "application/json").
		SetBody(api.data).
		SetResult(&widget).
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

	return widget, nil
}

// UpdateWidget handles modifying an existing widget.
type UpdateWidget struct {
	endpoint string
	client   *RestClient
	data     WidgetData
}

// IDBase sets the type of the widget.
func (api *UpdateWidget) IDBase(idBase string) *UpdateWidget {
	api.data.IDBase = idBase
	return api
}

// Sidebar sets the sidebar the widget belongs to.
func (api *UpdateWidget) Sidebar(sidebar string) *UpdateWidget {
	api.data.Sidebar = sidebar
	return api
}

// Instance sets instance settings of the widget.
func (api *UpdateWidget) Instance(instance any) *UpdateWidget {
	api.data.Instance = instance
	return api
}

// FormData sets URL-encoded form data.
func (api *UpdateWidget) FormData(formData string) *UpdateWidget {
	api.data.FormData = formData
	return api
}

// WithData populates the request body from a WidgetData struct.
func (api *UpdateWidget) WithData(data WidgetData) *UpdateWidget {
	api.data = data
	return api
}

// Do executes the update request.
func (api *UpdateWidget) Do() (widget *Widget, err error) {
	endpoint := api.client.endpoint + api.endpoint

	restyClient := api.client.httpClient.R()
	if api.client.auth.Username != "" && api.client.auth.Password != "" {
		restyClient = restyClient.SetBasicAuth(api.client.auth.Username, api.client.auth.Password)
	}

	resp, err := restyClient.
		SetHeader("Content-Type", "application/json").
		SetHeader("Accept", "application/json").
		SetBody(api.data).
		SetResult(&widget).
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

	return widget, nil
}

type deleteWidgetEnvelope struct {
	Widget
	Deleted  bool    `json:"deleted"`
	Previous *Widget `json:"previous"`
}

// DeleteWidget handles deleting or deactivating a widget.
type DeleteWidget struct {
	endpoint string
	client   *RestClient
	force    bool
}

// Force sets whether to permanently remove the widget (or move to inactive sidebar if false).
func (api *DeleteWidget) Force() *DeleteWidget {
	api.force = true
	return api
}

// Do executes the delete request.
func (api *DeleteWidget) Do() (widget *Widget, err error) {
	endpoint := api.client.endpoint + api.endpoint
	var env deleteWidgetEnvelope

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

	return &env.Widget, nil
}
