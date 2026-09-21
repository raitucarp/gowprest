package gowprest

import (
	"encoding/json"
	"strconv"
	"strings"
)

// Category represents a WordPress post category taxonomy term.
type Category struct {
	ID          int            `json:"id,omitempty"`
	Count       int            `json:"count,omitempty"`
	Description string         `json:"description,omitempty"`
	Link        string         `json:"link,omitempty"`
	Name        string         `json:"name,omitempty"`
	Slug        string         `json:"slug,omitempty"`
	Taxonomy    string         `json:"taxonomy,omitempty"`
	Parent      int            `json:"parent,omitempty"`
	Meta        any            `json:"meta,omitempty"`
	Links       map[string]any `json:"_links,omitempty"`
	Embedded    map[string]any `json:"_embedded,omitempty"`
}

// CategoryData holds mutable payload data for creating or updating a category.
type CategoryData struct {
	ID          int    `json:"id,omitempty"`
	Description string `json:"description,omitempty"`
	Name        string `json:"name,omitempty"`
	Slug        string `json:"slug,omitempty"`
	Parent      int    `json:"parent,omitempty"`
	Meta        any    `json:"meta,omitempty"`
}

// Categories provides operations for managing WordPress post categories (/wp/v2/categories).
type Categories struct {
	client *RestClient
}

// Categories returns the Categories service for managing post categories.
func (c *RestClient) Categories() *Categories {
	return &Categories{client: c}
}

// ListCategories is a fluent builder for querying collections of categories.
type ListCategories struct {
	endpoint  string
	client    *RestClient
	arguments map[string]string
}

// List initiates a query builder to list and filter categories.
func (api *Categories) List() *ListCategories {
	return &ListCategories{
		endpoint:  "/wp/v2/categories",
		client:    api.client,
		arguments: make(map[string]string),
	}
}

func (api *ListCategories) Context(ctx string) *ListCategories {
	api.arguments["context"] = ctx
	return api
}

func (api *ListCategories) ContextView() *ListCategories {
	api.arguments["context"] = "view"
	return api
}

func (api *ListCategories) ContextEdit() *ListCategories {
	api.arguments["context"] = "edit"
	return api
}

func (api *ListCategories) ContextEmbed() *ListCategories {
	api.arguments["context"] = "embed"
	return api
}

func (api *ListCategories) Page(page int) *ListCategories {
	api.arguments["page"] = strconv.Itoa(page)
	return api
}

func (api *ListCategories) PerPage(perPage int) *ListCategories {
	api.arguments["per_page"] = strconv.Itoa(perPage)
	return api
}

func (api *ListCategories) Search(query string) *ListCategories {
	api.arguments["search"] = query
	return api
}

func (api *ListCategories) Exclude(excludeIDs ...int) *ListCategories {
	excludes := make([]string, 0, len(excludeIDs))
	for _, excludeId := range excludeIDs {
		excludes = append(excludes, strconv.Itoa(excludeId))
	}
	api.arguments["exclude"] = strings.Join(excludes, ",")
	return api
}

func (api *ListCategories) Include(includeIDs ...int) *ListCategories {
	includes := make([]string, 0, len(includeIDs))
	for _, includeId := range includeIDs {
		includes = append(includes, strconv.Itoa(includeId))
	}
	api.arguments["include"] = strings.Join(includes, ",")
	return api
}

func (api *ListCategories) Offset(offset int) *ListCategories {
	api.arguments["offset"] = strconv.Itoa(offset)
	return api
}

func (api *ListCategories) Order(order string) *ListCategories {
	api.arguments["order"] = order
	return api
}

func (api *ListCategories) OrderAsc() *ListCategories {
	api.arguments["order"] = "asc"
	return api
}

func (api *ListCategories) OrderDesc() *ListCategories {
	api.arguments["order"] = "desc"
	return api
}

func (api *ListCategories) OrderBy(orderBy string) *ListCategories {
	api.arguments["orderby"] = orderBy
	return api
}

func (api *ListCategories) OrderById() *ListCategories {
	api.arguments["orderby"] = "id"
	return api
}

func (api *ListCategories) OrderByInclude() *ListCategories {
	api.arguments["orderby"] = "include"
	return api
}

func (api *ListCategories) OrderByName() *ListCategories {
	api.arguments["orderby"] = "name"
	return api
}

func (api *ListCategories) OrderBySlug() *ListCategories {
	api.arguments["orderby"] = "slug"
	return api
}

func (api *ListCategories) OrderByIncludeSlug() *ListCategories {
	api.arguments["orderby"] = "include_slugs"
	return api
}

func (api *ListCategories) OrderByTermGroup() *ListCategories {
	api.arguments["orderby"] = "term_group"
	return api
}

func (api *ListCategories) OrderByDescription() *ListCategories {
	api.arguments["orderby"] = "description"
	return api
}

func (api *ListCategories) OrderByCount() *ListCategories {
	api.arguments["orderby"] = "count"
	return api
}

func (api *ListCategories) HideEmpty(hide bool) *ListCategories {
	api.arguments["hide_empty"] = strconv.FormatBool(hide)
	return api
}

func (api *ListCategories) Parent(parentID int) *ListCategories {
	api.arguments["parent"] = strconv.Itoa(parentID)
	return api
}

func (api *ListCategories) Post(postID int) *ListCategories {
	api.arguments["post"] = strconv.Itoa(postID)
	return api
}

func (api *ListCategories) Slug(slugs ...string) *ListCategories {
	api.arguments["slug"] = strings.Join(slugs, ",")
	return api
}

func (api *ListCategories) Embed() *ListCategories {
	api.arguments["_embed"] = "true"
	return api
}

func (api *ListCategories) Fields(fields ...string) *ListCategories {
	api.arguments["_fields"] = strings.Join(fields, ",")
	return api
}

func (api *ListCategories) Do() (categories []Category, err error) {
	restyClient := api.client.httpClient.R()
	if api.client.auth.Username != "" && api.client.auth.Password != "" {
		restyClient.SetBasicAuth(api.client.auth.Username, api.client.auth.Password)
	}

	resp, err := restyClient.
		SetHeader("Accept", "application/json").
		SetResult(&categories).
		SetQueryParams(api.arguments).
		Get(api.client.endpoint + api.endpoint)

	if resp.IsError() {
		var wpError WPRestError
		err = json.Unmarshal(resp.Bytes(), &wpError)
		if err != nil {
			return
		}
		return categories, &wpError
	}

	return
}

// CreateCategory is a fluent builder for creating a new post category.
type CreateCategory struct {
	endpoint string
	client   *RestClient
	category CategoryData
}

// Create returns a CreateCategory builder. Can be called without arguments for fluent chaining,
// or with CategoryData for backward compatibility.
func (api *Categories) Create(category ...CategoryData) *CreateCategory {
	builder := &CreateCategory{
		endpoint: "/wp/v2/categories",
		client:   api.client,
	}
	if len(category) > 0 {
		builder.category = category[0]
	}
	return builder
}

func (api *CreateCategory) Name(name string) *CreateCategory {
	api.category.Name = name
	return api
}

func (api *CreateCategory) Description(description string) *CreateCategory {
	api.category.Description = description
	return api
}

func (api *CreateCategory) Slug(slug string) *CreateCategory {
	api.category.Slug = slug
	return api
}

func (api *CreateCategory) Parent(parentID int) *CreateCategory {
	api.category.Parent = parentID
	return api
}

func (api *CreateCategory) Meta(meta any) *CreateCategory {
	api.category.Meta = meta
	return api
}

func (api *CreateCategory) Do() (category Category, err error) {
	resp, err := api.client.httpClient.R().
		SetHeader("Content-Type", "application/json").
		SetBasicAuth(api.client.auth.Username, api.client.auth.Password).
		SetResult(&category).
		SetBody(api.category).
		Post(api.client.endpoint + api.endpoint)

	if resp.IsError() {
		var wpError WPRestError
		err = json.Unmarshal(resp.Bytes(), &wpError)
		if err != nil {
			return
		}
		return category, &wpError
	}

	return
}

// RetrieveCategory is a fluent builder for retrieving a single category.
type RetrieveCategory struct {
	endpoint  string
	client    *RestClient
	arguments map[string]string
}

// Retrieve returns a builder to fetch the category with the given ID.
func (api *Categories) Retrieve(categoryId int) *RetrieveCategory {
	return &RetrieveCategory{
		endpoint:  "/wp/v2/categories/" + strconv.Itoa(categoryId),
		client:    api.client,
		arguments: make(map[string]string),
	}
}

func (api *RetrieveCategory) Context(ctx string) *RetrieveCategory {
	api.arguments["context"] = ctx
	return api
}

func (api *RetrieveCategory) ContextView() *RetrieveCategory {
	api.arguments["context"] = "view"
	return api
}

func (api *RetrieveCategory) ContextEdit() *RetrieveCategory {
	api.arguments["context"] = "edit"
	return api
}

func (api *RetrieveCategory) ContextEmbed() *RetrieveCategory {
	api.arguments["context"] = "embed"
	return api
}

func (api *RetrieveCategory) Embed() *RetrieveCategory {
	api.arguments["_embed"] = "true"
	return api
}

func (api *RetrieveCategory) Fields(fields ...string) *RetrieveCategory {
	api.arguments["_fields"] = strings.Join(fields, ",")
	return api
}

func (api *RetrieveCategory) Do() (category *Category, err error) {
	endpoint := api.client.endpoint + api.endpoint

	restyClient := api.client.httpClient.R()
	if api.client.auth.Username != "" && api.client.auth.Password != "" {
		restyClient.SetBasicAuth(api.client.auth.Username, api.client.auth.Password)
	}

	resp, err := restyClient.
		SetHeader("Accept", "application/json").
		SetResult(&category).
		SetQueryParams(api.arguments).
		Get(endpoint)

	if resp.IsError() {
		var wpError WPRestError
		err = json.Unmarshal(resp.Bytes(), &wpError)
		if err != nil {
			return
		}
		return category, &wpError
	}

	return
}

// UpdateCategory is a fluent builder for modifying an existing post category.
type UpdateCategory struct {
	endpoint string
	client   *RestClient
	category CategoryData
}

// Update returns an UpdateCategory builder. Can be called without arguments or with CategoryData.
func (api *Categories) Update(category ...CategoryData) *UpdateCategory {
	builder := &UpdateCategory{
		endpoint: "/wp/v2/categories",
		client:   api.client,
	}
	if len(category) > 0 {
		builder.category = category[0]
		if builder.category.ID != 0 {
			builder.endpoint = "/wp/v2/categories/" + strconv.Itoa(builder.category.ID)
		}
	}
	return builder
}

func (api *UpdateCategory) ID(id int) *UpdateCategory {
	api.category.ID = id
	api.endpoint = "/wp/v2/categories/" + strconv.Itoa(id)
	return api
}

func (api *UpdateCategory) Name(name string) *UpdateCategory {
	api.category.Name = name
	return api
}

func (api *UpdateCategory) Description(description string) *UpdateCategory {
	api.category.Description = description
	return api
}

func (api *UpdateCategory) Slug(slug string) *UpdateCategory {
	api.category.Slug = slug
	return api
}

func (api *UpdateCategory) Parent(parentID int) *UpdateCategory {
	api.category.Parent = parentID
	return api
}

func (api *UpdateCategory) Meta(meta any) *UpdateCategory {
	api.category.Meta = meta
	return api
}

func (api *UpdateCategory) Do() (category Category, err error) {
	resp, err := api.client.httpClient.R().
		SetHeader("Content-Type", "application/json").
		SetBasicAuth(api.client.auth.Username, api.client.auth.Password).
		SetResult(&category).
		SetBody(api.category).
		Post(api.client.endpoint + api.endpoint)

	if resp.IsError() {
		var wpError WPRestError
		err = json.Unmarshal(resp.Bytes(), &wpError)
		if err != nil {
			return
		}
		return category, &wpError
	}

	return
}

type deleteCategoryEnvelope struct {
	Category
	Deleted  bool      `json:"deleted"`
	Previous *Category `json:"previous"`
}

// DeleteCategory is a fluent builder for deleting a post category.
type DeleteCategory struct {
	endpoint   string
	client     *RestClient
	categoryId int
	force      bool
}

// Delete returns a builder to delete the category with the given ID.
func (api *Categories) Delete(categoryId int) *DeleteCategory {
	return &DeleteCategory{
		endpoint:   "/wp/v2/categories",
		client:     api.client,
		categoryId: categoryId,
		force:      true, // Terms do not support trashing in WP REST API
	}
}

func (api *DeleteCategory) Force() *DeleteCategory {
	api.force = true
	return api
}

func (api *DeleteCategory) Do() (category Category, err error) {
	endpoint := api.client.endpoint + api.endpoint + "/" + strconv.Itoa(api.categoryId)
	var env deleteCategoryEnvelope
	resp, err :=
		api.client.httpClient.R().
			SetHeader("Content-Type", "application/json").
			SetBasicAuth(api.client.auth.Username, api.client.auth.Password).
			SetResult(&env).
			SetQueryParam("force", strconv.FormatBool(api.force)).
			Delete(endpoint)

	if resp.IsError() {
		var wpError WPRestError
		err = json.Unmarshal(resp.Bytes(), &wpError)
		if err != nil {
			return
		}
		return category, &wpError
	}

	if err != nil {
		return
	}

	if env.Previous != nil && env.Previous.ID != 0 {
		return *env.Previous, nil
	}

	return env.Category, nil
}
