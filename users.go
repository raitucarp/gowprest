package gowprest

import (
	"encoding/json"
	"strconv"
	"strings"
)

// User represents a WordPress user.
type User struct {
	ID                int               `json:"id,omitempty"`
	Username          string            `json:"username,omitempty"`
	Name              string            `json:"name,omitempty"`
	FirstName         string            `json:"first_name,omitempty"`
	LastName          string            `json:"last_name,omitempty"`
	Email             string            `json:"email,omitempty"`
	URL               string            `json:"url,omitempty"`
	Description       string            `json:"description,omitempty"`
	Link              string            `json:"link,omitempty"`
	Locale            string            `json:"locale,omitempty"`
	Nickname          string            `json:"nickname,omitempty"`
	Slug              string            `json:"slug,omitempty"`
	RegisteredDate    *Date             `json:"registered_date,omitempty"`
	Roles             []string          `json:"roles,omitempty"`
	Password          string            `json:"password,omitempty"`
	Capabilities      map[string]bool   `json:"capabilities,omitempty"`
	ExtraCapabilities map[string]bool   `json:"extra_capabilities,omitempty"`
	AvatarURLs        map[string]string `json:"avatar_urls,omitempty"`
	Meta              Meta              `json:"meta,omitempty"`
	Links             map[string]any    `json:"_links,omitempty"`
	Embedded          map[string]any    `json:"_embedded,omitempty"`
}

// UserData represents user payload for create and update operations.
type UserData struct {
	ID          int            `json:"id,omitempty"`
	Username    string         `json:"username,omitempty"`
	Name        string         `json:"name,omitempty"`
	FirstName   string         `json:"first_name,omitempty"`
	LastName    string         `json:"last_name,omitempty"`
	Email       string         `json:"email,omitempty"`
	URL         string         `json:"url,omitempty"`
	Description string         `json:"description,omitempty"`
	Locale      string         `json:"locale,omitempty"`
	Nickname    string         `json:"nickname,omitempty"`
	Slug        string         `json:"slug,omitempty"`
	Roles       []string       `json:"roles,omitempty"`
	Password    string         `json:"password,omitempty"`
	Meta        map[string]any `json:"meta,omitempty"`
}

// DeletedUser represents response when a user is deleted.
type DeletedUser struct {
	User
	Deleted  bool  `json:"deleted"`
	Previous *User `json:"previous,omitempty"`
}

// Users handles requests to the WordPress users API.
type Users struct {
	client *RestClient
}

// Users returns a Users service instance.
func (c *RestClient) Users() *Users {
	return &Users{client: c}
}

// UserMe handles requests to the /wp/v2/users/me endpoint.
type UserMe struct {
	client *RestClient
}

// Me returns a UserMe service instance for operating on the currently authenticated user.
func (api *Users) Me() *UserMe {
	return &UserMe{client: api.client}
}

// ApplicationPasswords returns an ApplicationPasswords service instance for the given user ID (defaults to "me").
func (api *Users) ApplicationPasswords(userID ...any) *ApplicationPasswords {
	target := "me"
	if len(userID) > 0 {
		switch v := userID[0].(type) {
		case int:
			target = strconv.Itoa(v)
		case string:
			target = v
		}
	}
	return &ApplicationPasswords{
		client: api.client,
		userID: target,
	}
}

// ApplicationPasswords returns an ApplicationPasswords service instance for the authenticated user.
func (m *UserMe) ApplicationPasswords() *ApplicationPasswords {
	return &ApplicationPasswords{
		client: m.client,
		userID: "me",
	}
}

// ApplicationPasswords returns an ApplicationPasswords service instance (defaults to "me").
func (c *RestClient) ApplicationPasswords(userID ...any) *ApplicationPasswords {
	return c.Users().ApplicationPasswords(userID...)
}

// Retrieve returns a RetrieveUser builder configured for /wp/v2/users/me.
func (m *UserMe) Retrieve() *RetrieveUser {
	return &RetrieveUser{
		endpoint:  "/wp/v2/users/me",
		client:    m.client,
		arguments: make(map[string]string),
	}
}

// Update returns an UpdateUser builder configured for /wp/v2/users/me.
func (m *UserMe) Update(user ...UserData) *UpdateUser {
	var u UserData
	if len(user) > 0 {
		u = user[0]
	}
	return &UpdateUser{
		endpoint: "/wp/v2/users/me",
		client:   m.client,
		user:     u,
	}
}

// Delete returns a DeleteUser builder configured for /wp/v2/users/me.
func (m *UserMe) Delete(reassignID ...int) *DeleteUser {
	del := &DeleteUser{
		endpoint: "/wp/v2/users/me",
		client:   m.client,
		force:    true,
	}
	if len(reassignID) > 0 {
		del.Reassign(reassignID[0])
	}
	return del
}

// List returns a ListUsers builder to query user collections.
func (api *Users) List() *ListUsers {
	return &ListUsers{
		endpoint:  "/wp/v2/users",
		client:    api.client,
		arguments: make(map[string]string),
	}
}

// Retrieve returns a RetrieveUser builder to get a specific user by ID.
func (api *Users) Retrieve(userID int) *RetrieveUser {
	return &RetrieveUser{
		endpoint:  "/wp/v2/users/" + strconv.Itoa(userID),
		client:    api.client,
		arguments: make(map[string]string),
	}
}

// Create returns a CreateUser builder to create a new user.
func (api *Users) Create(user ...UserData) *CreateUser {
	builder := &CreateUser{
		endpoint: "/wp/v2/users",
		client:   api.client,
	}
	if len(user) > 0 {
		builder.user = user[0]
	}
	return builder
}

// Update returns an UpdateUser builder to modify a user.
func (api *Users) Update(user ...UserData) *UpdateUser {
	endpoint := "/wp/v2/users"
	var u UserData
	if len(user) > 0 {
		u = user[0]
		if u.ID != 0 {
			endpoint = "/wp/v2/users/" + strconv.Itoa(u.ID)
		}
	}
	return &UpdateUser{
		endpoint: endpoint,
		client:   api.client,
		user:     u,
	}
}

// Delete returns a DeleteUser builder to delete a user.
func (api *Users) Delete(userID int) *DeleteUser {
	return &DeleteUser{
		endpoint: "/wp/v2/users/" + strconv.Itoa(userID),
		client:   api.client,
		userID:   userID,
		force:    true, // Users require force=true
	}
}

// ListUsers handles querying user collections.
type ListUsers struct {
	endpoint  string
	client    *RestClient
	arguments map[string]string
}

func (api *ListUsers) Context(ctx string) *ListUsers {
	api.arguments["context"] = ctx
	return api
}

func (api *ListUsers) ContextView() *ListUsers {
	api.arguments["context"] = "view"
	return api
}

func (api *ListUsers) ContextEdit() *ListUsers {
	api.arguments["context"] = "edit"
	return api
}

func (api *ListUsers) ContextEmbed() *ListUsers {
	api.arguments["context"] = "embed"
	return api
}

func (api *ListUsers) Page(page int) *ListUsers {
	api.arguments["page"] = strconv.Itoa(page)
	return api
}

func (api *ListUsers) PerPage(perPage int) *ListUsers {
	api.arguments["per_page"] = strconv.Itoa(perPage)
	return api
}

func (api *ListUsers) Offset(offset int) *ListUsers {
	api.arguments["offset"] = strconv.Itoa(offset)
	return api
}

func (api *ListUsers) Search(search string) *ListUsers {
	api.arguments["search"] = search
	return api
}

func (api *ListUsers) Exclude(ids ...int) *ListUsers {
	strIDs := make([]string, len(ids))
	for i, id := range ids {
		strIDs[i] = strconv.Itoa(id)
	}
	api.arguments["exclude"] = strings.Join(strIDs, ",")
	return api
}

func (api *ListUsers) Include(ids ...int) *ListUsers {
	strIDs := make([]string, len(ids))
	for i, id := range ids {
		strIDs[i] = strconv.Itoa(id)
	}
	api.arguments["include"] = strings.Join(strIDs, ",")
	return api
}

func (api *ListUsers) Order(order string) *ListUsers {
	api.arguments["order"] = order
	return api
}

func (api *ListUsers) OrderAsc() *ListUsers {
	api.arguments["order"] = "asc"
	return api
}

func (api *ListUsers) OrderDesc() *ListUsers {
	api.arguments["order"] = "desc"
	return api
}

func (api *ListUsers) OrderBy(orderBy string) *ListUsers {
	api.arguments["orderby"] = orderBy
	return api
}

func (api *ListUsers) OrderByID() *ListUsers {
	api.arguments["orderby"] = "id"
	return api
}

func (api *ListUsers) OrderById() *ListUsers {
	return api.OrderByID()
}

func (api *ListUsers) OrderByInclude() *ListUsers {
	api.arguments["orderby"] = "include"
	return api
}

func (api *ListUsers) OrderByName() *ListUsers {
	api.arguments["orderby"] = "name"
	return api
}

func (api *ListUsers) OrderByRegisteredDate() *ListUsers {
	api.arguments["orderby"] = "registered_date"
	return api
}

func (api *ListUsers) OrderBySlug() *ListUsers {
	api.arguments["orderby"] = "slug"
	return api
}

func (api *ListUsers) OrderByIncludeSlugs() *ListUsers {
	api.arguments["orderby"] = "include_slugs"
	return api
}

func (api *ListUsers) OrderByEmail() *ListUsers {
	api.arguments["orderby"] = "email"
	return api
}

func (api *ListUsers) OrderByURL() *ListUsers {
	api.arguments["orderby"] = "url"
	return api
}

func (api *ListUsers) OrderByUrl() *ListUsers {
	return api.OrderByURL()
}

func (api *ListUsers) Slug(slugs ...string) *ListUsers {
	api.arguments["slug"] = strings.Join(slugs, ",")
	return api
}

func (api *ListUsers) Roles(roles ...string) *ListUsers {
	api.arguments["roles"] = strings.Join(roles, ",")
	return api
}

func (api *ListUsers) Capabilities(caps ...string) *ListUsers {
	api.arguments["capabilities"] = strings.Join(caps, ",")
	return api
}

func (api *ListUsers) Who(who string) *ListUsers {
	api.arguments["who"] = who
	return api
}

func (api *ListUsers) WhoAuthors() *ListUsers {
	api.arguments["who"] = "authors"
	return api
}

func (api *ListUsers) HasPublishedPosts(has bool) *ListUsers {
	api.arguments["has_published_posts"] = strconv.FormatBool(has)
	return api
}

func (api *ListUsers) HasPublishedPostsFor(postTypes ...string) *ListUsers {
	api.arguments["has_published_posts"] = strings.Join(postTypes, ",")
	return api
}

func (api *ListUsers) Embed() *ListUsers {
	api.arguments["_embed"] = "true"
	return api
}

func (api *ListUsers) Fields(fields ...string) *ListUsers {
	api.arguments["_fields"] = strings.Join(fields, ",")
	return api
}

func (api *ListUsers) Do() (users []User, err error) {
	restyClient := api.client.httpClient.R()
	if api.client.auth.Username != "" && api.client.auth.Password != "" {
		restyClient.SetBasicAuth(api.client.auth.Username, api.client.auth.Password)
	}

	resp, err := restyClient.
		SetHeader("Accept", "application/json").
		SetResult(&users).
		SetQueryParams(api.arguments).
		Get(api.client.endpoint + api.endpoint)

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

	return users, nil
}

// CreateUser handles creating a new user.
type CreateUser struct {
	endpoint string
	client   *RestClient
	user     UserData
}

func (api *CreateUser) Username(username string) *CreateUser {
	api.user.Username = username
	return api
}

func (api *CreateUser) Name(name string) *CreateUser {
	api.user.Name = name
	return api
}

func (api *CreateUser) FirstName(firstName string) *CreateUser {
	api.user.FirstName = firstName
	return api
}

func (api *CreateUser) LastName(lastName string) *CreateUser {
	api.user.LastName = lastName
	return api
}

func (api *CreateUser) Email(email string) *CreateUser {
	api.user.Email = email
	return api
}

func (api *CreateUser) URL(url string) *CreateUser {
	api.user.URL = url
	return api
}

func (api *CreateUser) Description(description string) *CreateUser {
	api.user.Description = description
	return api
}

func (api *CreateUser) Locale(locale string) *CreateUser {
	api.user.Locale = locale
	return api
}

func (api *CreateUser) Nickname(nickname string) *CreateUser {
	api.user.Nickname = nickname
	return api
}

func (api *CreateUser) Slug(slug string) *CreateUser {
	api.user.Slug = slug
	return api
}

func (api *CreateUser) Roles(roles ...string) *CreateUser {
	api.user.Roles = roles
	return api
}

func (api *CreateUser) Password(password string) *CreateUser {
	api.user.Password = password
	return api
}

func (api *CreateUser) Meta(meta map[string]any) *CreateUser {
	api.user.Meta = meta
	return api
}

func (api *CreateUser) Do() (*User, error) {
	restyClient := api.client.httpClient.R()
	if api.client.auth.Username != "" && api.client.auth.Password != "" {
		restyClient.SetBasicAuth(api.client.auth.Username, api.client.auth.Password)
	}

	var user User
	resp, err := restyClient.
		SetHeader("Content-Type", "application/json").
		SetHeader("Accept", "application/json").
		SetBody(api.user).
		SetResult(&user).
		Post(api.client.endpoint + api.endpoint)

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

	return &user, nil
}

// RetrieveUser handles retrieving a single user.
type RetrieveUser struct {
	endpoint  string
	client    *RestClient
	arguments map[string]string
}

func (api *RetrieveUser) Context(ctx string) *RetrieveUser {
	api.arguments["context"] = ctx
	return api
}

func (api *RetrieveUser) ContextView() *RetrieveUser {
	api.arguments["context"] = "view"
	return api
}

func (api *RetrieveUser) ContextEdit() *RetrieveUser {
	api.arguments["context"] = "edit"
	return api
}

func (api *RetrieveUser) ContextEmbed() *RetrieveUser {
	api.arguments["context"] = "embed"
	return api
}

func (api *RetrieveUser) Embed() *RetrieveUser {
	api.arguments["_embed"] = "true"
	return api
}

func (api *RetrieveUser) Fields(fields ...string) *RetrieveUser {
	api.arguments["_fields"] = strings.Join(fields, ",")
	return api
}

func (api *RetrieveUser) Do() (*User, error) {
	restyClient := api.client.httpClient.R()
	if api.client.auth.Username != "" && api.client.auth.Password != "" {
		restyClient.SetBasicAuth(api.client.auth.Username, api.client.auth.Password)
	}

	var user User
	resp, err := restyClient.
		SetHeader("Accept", "application/json").
		SetResult(&user).
		SetQueryParams(api.arguments).
		Get(api.client.endpoint + api.endpoint)

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

	return &user, nil
}

// UpdateUser handles updating a user.
type UpdateUser struct {
	endpoint string
	client   *RestClient
	user     UserData
}

func (api *UpdateUser) ID(id int) *UpdateUser {
	api.user.ID = id
	if !strings.HasSuffix(api.endpoint, "/me") {
		api.endpoint = "/wp/v2/users/" + strconv.Itoa(id)
	}
	return api
}

func (api *UpdateUser) Username(username string) *UpdateUser {
	api.user.Username = username
	return api
}

func (api *UpdateUser) Name(name string) *UpdateUser {
	api.user.Name = name
	return api
}

func (api *UpdateUser) FirstName(firstName string) *UpdateUser {
	api.user.FirstName = firstName
	return api
}

func (api *UpdateUser) LastName(lastName string) *UpdateUser {
	api.user.LastName = lastName
	return api
}

func (api *UpdateUser) Email(email string) *UpdateUser {
	api.user.Email = email
	return api
}

func (api *UpdateUser) URL(url string) *UpdateUser {
	api.user.URL = url
	return api
}

func (api *UpdateUser) Description(description string) *UpdateUser {
	api.user.Description = description
	return api
}

func (api *UpdateUser) Locale(locale string) *UpdateUser {
	api.user.Locale = locale
	return api
}

func (api *UpdateUser) Nickname(nickname string) *UpdateUser {
	api.user.Nickname = nickname
	return api
}

func (api *UpdateUser) Slug(slug string) *UpdateUser {
	api.user.Slug = slug
	return api
}

func (api *UpdateUser) Roles(roles ...string) *UpdateUser {
	api.user.Roles = roles
	return api
}

func (api *UpdateUser) Password(password string) *UpdateUser {
	api.user.Password = password
	return api
}

func (api *UpdateUser) Meta(meta map[string]any) *UpdateUser {
	api.user.Meta = meta
	return api
}

func (api *UpdateUser) Do() (*User, error) {
	restyClient := api.client.httpClient.R()
	if api.client.auth.Username != "" && api.client.auth.Password != "" {
		restyClient.SetBasicAuth(api.client.auth.Username, api.client.auth.Password)
	}

	var user User
	resp, err := restyClient.
		SetHeader("Content-Type", "application/json").
		SetHeader("Accept", "application/json").
		SetBody(api.user).
		SetResult(&user).
		Post(api.client.endpoint + api.endpoint)

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

	return &user, nil
}

// DeleteUser handles deleting a user.
type DeleteUser struct {
	endpoint string
	client   *RestClient
	userID   int
	force    bool
	reassign *int
}

func (api *DeleteUser) Force(force ...bool) *DeleteUser {
	if len(force) > 0 {
		api.force = force[0]
	} else {
		api.force = true
	}
	return api
}

func (api *DeleteUser) Reassign(reassignID int) *DeleteUser {
	api.reassign = &reassignID
	return api
}

func (api *DeleteUser) Do() (*DeletedUser, error) {
	restyClient := api.client.httpClient.R()
	if api.client.auth.Username != "" && api.client.auth.Password != "" {
		restyClient.SetBasicAuth(api.client.auth.Username, api.client.auth.Password)
	}

	queryParams := make(map[string]string)
	queryParams["force"] = strconv.FormatBool(api.force)
	if api.reassign != nil {
		queryParams["reassign"] = strconv.Itoa(*api.reassign)
	}

	var res DeletedUser
	resp, err := restyClient.
		SetHeader("Accept", "application/json").
		SetQueryParams(queryParams).
		SetResult(&res).
		Delete(api.client.endpoint + api.endpoint)

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

	if res.Previous != nil {
		res.User = *res.Previous
	}

	return &res, nil
}
