package gowprest

import (
	"encoding/json"
	"strings"
)

// ApplicationPassword represents an application password record for a WordPress user.
type ApplicationPassword struct {
	UUID     string         `json:"uuid,omitempty"`
	AppID    string         `json:"app_id,omitempty"`
	Name     string         `json:"name,omitempty"`
	Password string         `json:"password,omitempty"`
	Created  *Date          `json:"created,omitempty"`
	LastUsed *Date          `json:"last_used,omitempty"`
	LastIP   *string        `json:"last_ip,omitempty"`
	Links    map[string]any `json:"_links,omitempty"`
}

// ApplicationPasswords handles requests to the WordPress Application Passwords API (/wp/v2/users/<id>/application-passwords).
type ApplicationPasswords struct {
	client *RestClient
	userID string
}

// List returns a ListApplicationPasswords builder to list application passwords.
func (api *ApplicationPasswords) List() *ListApplicationPasswords {
	return &ListApplicationPasswords{
		endpoint:  "/wp/v2/users/" + api.userID + "/application-passwords",
		client:    api.client,
		arguments: make(map[string]string),
	}
}

// Retrieve returns a RetrieveApplicationPassword builder to get an application password by UUID.
func (api *ApplicationPasswords) Retrieve(uuid string) *RetrieveApplicationPassword {
	return &RetrieveApplicationPassword{
		endpoint:  "/wp/v2/users/" + api.userID + "/application-passwords/" + uuid,
		client:    api.client,
		arguments: make(map[string]string),
	}
}

// Introspect returns an IntrospectApplicationPassword builder to check the current credentials.
func (api *ApplicationPasswords) Introspect() *IntrospectApplicationPassword {
	return &IntrospectApplicationPassword{
		endpoint:  "/wp/v2/users/" + api.userID + "/application-passwords/introspect",
		client:    api.client,
		arguments: make(map[string]string),
	}
}

// Create returns a CreateApplicationPassword builder to generate a new application password.
func (api *ApplicationPasswords) Create(name ...string) *CreateApplicationPassword {
	builder := &CreateApplicationPassword{
		endpoint: "/wp/v2/users/" + api.userID + "/application-passwords",
		client:   api.client,
		body:     make(map[string]any),
	}
	if len(name) > 0 {
		builder.Name(name[0])
	}
	return builder
}

// Update returns an UpdateApplicationPassword builder to update an existing application password.
func (api *ApplicationPasswords) Update(uuid string) *UpdateApplicationPassword {
	return &UpdateApplicationPassword{
		endpoint: "/wp/v2/users/" + api.userID + "/application-passwords/" + uuid,
		client:   api.client,
		body:     make(map[string]any),
	}
}

// Delete returns a DeleteApplicationPassword builder to delete a specific application password by UUID.
func (api *ApplicationPasswords) Delete(uuid string) *DeleteApplicationPassword {
	return &DeleteApplicationPassword{
		endpoint: "/wp/v2/users/" + api.userID + "/application-passwords/" + uuid,
		client:   api.client,
	}
}

// DeleteAll returns a DeleteAllApplicationPasswords builder to delete all application passwords for the user.
func (api *ApplicationPasswords) DeleteAll() *DeleteAllApplicationPasswords {
	return &DeleteAllApplicationPasswords{
		endpoint: "/wp/v2/users/" + api.userID + "/application-passwords",
		client:   api.client,
	}
}

// ListApplicationPasswords handles querying application password collections.
type ListApplicationPasswords struct {
	endpoint  string
	client    *RestClient
	arguments map[string]string
}

func (api *ListApplicationPasswords) Context(ctx string) *ListApplicationPasswords {
	api.arguments["context"] = ctx
	return api
}

func (api *ListApplicationPasswords) ContextView() *ListApplicationPasswords {
	return api.Context("view")
}

func (api *ListApplicationPasswords) ContextEdit() *ListApplicationPasswords {
	return api.Context("edit")
}

func (api *ListApplicationPasswords) ContextEmbed() *ListApplicationPasswords {
	return api.Context("embed")
}

func (api *ListApplicationPasswords) Fields(fields ...string) *ListApplicationPasswords {
	api.arguments["_fields"] = strings.Join(fields, ",")
	return api
}

func (api *ListApplicationPasswords) Do() (passwords []ApplicationPassword, err error) {
	endpoint := api.client.endpoint + api.endpoint

	restyClient := api.client.httpClient.R()
	if api.client.auth.Username != "" && api.client.auth.Password != "" {
		restyClient.SetBasicAuth(api.client.auth.Username, api.client.auth.Password)
	}

	resp, err := restyClient.
		SetHeader("Accept", "application/json").
		SetResult(&passwords).
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

	return passwords, nil
}

// RetrieveApplicationPassword handles getting a single application password.
type RetrieveApplicationPassword struct {
	endpoint  string
	client    *RestClient
	arguments map[string]string
}

func (api *RetrieveApplicationPassword) Context(ctx string) *RetrieveApplicationPassword {
	api.arguments["context"] = ctx
	return api
}

func (api *RetrieveApplicationPassword) ContextView() *RetrieveApplicationPassword {
	return api.Context("view")
}

func (api *RetrieveApplicationPassword) ContextEdit() *RetrieveApplicationPassword {
	return api.Context("edit")
}

func (api *RetrieveApplicationPassword) ContextEmbed() *RetrieveApplicationPassword {
	return api.Context("embed")
}

func (api *RetrieveApplicationPassword) Fields(fields ...string) *RetrieveApplicationPassword {
	api.arguments["_fields"] = strings.Join(fields, ",")
	return api
}

func (api *RetrieveApplicationPassword) Do() (password *ApplicationPassword, err error) {
	endpoint := api.client.endpoint + api.endpoint

	restyClient := api.client.httpClient.R()
	if api.client.auth.Username != "" && api.client.auth.Password != "" {
		restyClient.SetBasicAuth(api.client.auth.Username, api.client.auth.Password)
	}

	resp, err := restyClient.
		SetHeader("Accept", "application/json").
		SetResult(&password).
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

	return password, nil
}

// IntrospectApplicationPassword handles introspecting the currently authenticated application password.
type IntrospectApplicationPassword struct {
	endpoint  string
	client    *RestClient
	arguments map[string]string
}

func (api *IntrospectApplicationPassword) Context(ctx string) *IntrospectApplicationPassword {
	api.arguments["context"] = ctx
	return api
}

func (api *IntrospectApplicationPassword) ContextView() *IntrospectApplicationPassword {
	return api.Context("view")
}

func (api *IntrospectApplicationPassword) ContextEdit() *IntrospectApplicationPassword {
	return api.Context("edit")
}

func (api *IntrospectApplicationPassword) ContextEmbed() *IntrospectApplicationPassword {
	return api.Context("embed")
}

func (api *IntrospectApplicationPassword) Fields(fields ...string) *IntrospectApplicationPassword {
	api.arguments["_fields"] = strings.Join(fields, ",")
	return api
}

func (api *IntrospectApplicationPassword) Do() (password *ApplicationPassword, err error) {
	endpoint := api.client.endpoint + api.endpoint

	restyClient := api.client.httpClient.R()
	if api.client.auth.Username != "" && api.client.auth.Password != "" {
		restyClient.SetBasicAuth(api.client.auth.Username, api.client.auth.Password)
	}

	resp, err := restyClient.
		SetHeader("Accept", "application/json").
		SetResult(&password).
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

	return password, nil
}

// CreateApplicationPassword handles generating a new application password.
type CreateApplicationPassword struct {
	endpoint string
	client   *RestClient
	body     map[string]any
}

func (api *CreateApplicationPassword) Name(name string) *CreateApplicationPassword {
	api.body["name"] = name
	return api
}

func (api *CreateApplicationPassword) AppID(appID string) *CreateApplicationPassword {
	api.body["app_id"] = appID
	return api
}

func (api *CreateApplicationPassword) Do() (password *ApplicationPassword, err error) {
	endpoint := api.client.endpoint + api.endpoint

	restyClient := api.client.httpClient.R()
	if api.client.auth.Username != "" && api.client.auth.Password != "" {
		restyClient.SetBasicAuth(api.client.auth.Username, api.client.auth.Password)
	}

	resp, err := restyClient.
		SetHeader("Content-Type", "application/json").
		SetHeader("Accept", "application/json").
		SetBody(api.body).
		SetResult(&password).
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

	return password, nil
}

// UpdateApplicationPassword handles modifying an application password.
type UpdateApplicationPassword struct {
	endpoint string
	client   *RestClient
	body     map[string]any
}

func (api *UpdateApplicationPassword) Name(name string) *UpdateApplicationPassword {
	api.body["name"] = name
	return api
}

func (api *UpdateApplicationPassword) AppID(appID string) *UpdateApplicationPassword {
	api.body["app_id"] = appID
	return api
}

func (api *UpdateApplicationPassword) Do() (password *ApplicationPassword, err error) {
	endpoint := api.client.endpoint + api.endpoint

	restyClient := api.client.httpClient.R()
	if api.client.auth.Username != "" && api.client.auth.Password != "" {
		restyClient.SetBasicAuth(api.client.auth.Username, api.client.auth.Password)
	}

	resp, err := restyClient.
		SetHeader("Content-Type", "application/json").
		SetHeader("Accept", "application/json").
		SetBody(api.body).
		SetResult(&password).
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

	return password, nil
}

type deleteAppPasswordEnvelope struct {
	ApplicationPassword
	Deleted  bool                 `json:"deleted"`
	Previous *ApplicationPassword `json:"previous"`
}

// DeleteApplicationPassword handles deleting a single application password.
type DeleteApplicationPassword struct {
	endpoint string
	client   *RestClient
}

func (api *DeleteApplicationPassword) Do() (password *ApplicationPassword, err error) {
	endpoint := api.client.endpoint + api.endpoint

	restyClient := api.client.httpClient.R()
	if api.client.auth.Username != "" && api.client.auth.Password != "" {
		restyClient.SetBasicAuth(api.client.auth.Username, api.client.auth.Password)
	}

	var env deleteAppPasswordEnvelope
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

	if env.Deleted && env.Previous != nil && env.Previous.UUID != "" {
		return env.Previous, nil
	}

	return &env.ApplicationPassword, nil
}

type deleteAllAppPasswordsEnvelope struct {
	Deleted bool `json:"deleted"`
}

// DeleteAllApplicationPasswords handles deleting all application passwords for a user.
type DeleteAllApplicationPasswords struct {
	endpoint string
	client   *RestClient
}

func (api *DeleteAllApplicationPasswords) Do() (deleted bool, err error) {
	endpoint := api.client.endpoint + api.endpoint

	restyClient := api.client.httpClient.R()
	if api.client.auth.Username != "" && api.client.auth.Password != "" {
		restyClient.SetBasicAuth(api.client.auth.Username, api.client.auth.Password)
	}

	var env deleteAllAppPasswordsEnvelope
	resp, err := restyClient.
		SetHeader("Accept", "application/json").
		SetResult(&env).
		Delete(endpoint)

	if err != nil {
		return false, err
	}

	if resp.IsError() {
		var wpError WPRestError
		err = json.Unmarshal(resp.Bytes(), &wpError)
		if err != nil {
			return false, err
		}
		return false, &wpError
	}

	return env.Deleted, nil
}
