package gowprest

import (
	"encoding/json"
	"strconv"
	"strings"
)

type Category struct {
	ID          int    `json:"id,omitempty"`
	Count       int    `json:"count,omitempty"`
	Description string `json:"description,omitempty"`
	Link        string `json:"link,omitempty"`
	Name        string `json:"name,omitempty"`
	Slug        string `json:"slug,omitempty"`
	Taxonomy    string `json:"taxonomy,omitempty"`
	Parent      int    `json:"parent,omitempty"`
	Meta        any    `json:"meta,omitempty"`
}

type CategoryData struct {
	ID          int    `json:"id,omitempty"`
	Description string `json:"description,omitempty"`
	Name        string `json:"name,omitempty"`
	Slug        string `json:"slug,omitempty"`
	Parent      int    `json:"parent,omitempty"`
	Meta        any    `json:"meta,omitempty"`
}

type Categories struct {
	client *RestClient
}

func (c *RestClient) Categories() *Categories {
	return &Categories{client: c}
}

type ListCategories struct {
	endpoint  string
	client    *RestClient
	arguments map[string]string
}

func (api *Categories) List() *ListCategories {
	return &ListCategories{
		endpoint:  "/wp/v2/categories",
		client:    api.client,
		arguments: make(map[string]string),
	}
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
	excludes := []string{}
	for _, excludeId := range excludeIDs {
		excludes = append(excludes, strconv.Itoa(excludeId))
	}
	api.arguments["exclude"] = strings.Join(excludes, ",")
	return api
}

func (api *ListCategories) Include(includeIDs ...int) *ListCategories {
	includes := []string{}
	for _, includeId := range includeIDs {
		includes = append(includes, strconv.Itoa(includeId))
	}
	api.arguments["include"] = strings.Join(includes, ",")
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

func (api *ListCategories) Do() (categories []Category, err error) {
	_, err = api.client.httpClient.R().
		SetHeader("Accept", "application/json").
		SetResult(&categories).
		SetQueryParams(api.arguments).
		Get(api.client.endpoint + api.endpoint)

	if err != nil {
		return
	}

	return
}

type CreateCategory struct {
	endpoint string
	client   *RestClient
	category CategoryData
}

func (api *Categories) Create(category CategoryData) *CreateCategory {
	return &CreateCategory{
		endpoint: "/wp/v2/categories",
		client:   api.client,
		category: category,
	}
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

	if err != nil {
		return
	}

	return
}

type RetrieveCategory struct {
	endpoint  string
	client    *RestClient
	arguments map[string]string
}

func (api *Categories) Retrieve(categoryId int) *RetrieveCategory {
	return &RetrieveCategory{
		endpoint:  "/wp/v2/categories/" + strconv.Itoa(categoryId),
		client:    api.client,
		arguments: make(map[string]string),
	}
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

func (api *RetrieveCategory) Do() (category *Category, err error) {
	endpoint := api.client.endpoint + api.endpoint

	restyClient := api.client.httpClient.R()
	if api.client.auth.Username != "" && api.client.auth.Password != "" && api.arguments["context"] == "edit" {
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

	if err != nil {
		return
	}

	return
}

type UpdateCategory struct {
	endpoint string
	client   *RestClient
	category CategoryData
}

func (api *Categories) Update(category CategoryData) *UpdateCategory {
	return &UpdateCategory{
		endpoint: "/wp/v2/categories/" + strconv.Itoa(category.ID),
		client:   api.client,
		category: category,
	}
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

	if err != nil {
		return
	}

	return
}

type DeleteCategory struct {
	endpoint   string
	client     *RestClient
	categoryId int
	force      bool
}

func (api *Categories) Delete(categoryId int) *DeleteCategory {
	return &DeleteCategory{
		endpoint:   "/wp/v2/categories",
		client:     api.client,
		categoryId: categoryId,
	}
}

func (api *DeleteCategory) Force() *DeleteCategory {
	api.force = true
	return api
}

func (api *DeleteCategory) Do() (category Category, err error) {
	endpoint := api.client.endpoint + api.endpoint + "/" + strconv.Itoa(api.categoryId)
	resp, err :=
		api.client.httpClient.R().
			SetHeader("Content-Type", "application/json").
			SetBasicAuth(api.client.auth.Username, api.client.auth.Password).
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

	if api.force {
		var nested struct {
			Deleted  bool     `json:"deleted"`
			Previous Category `json:"previous"`
		}
		err = json.Unmarshal(resp.Bytes(), &nested)
		if err != nil {
			return
		}
		return nested.Previous, nil
	}

	err = json.Unmarshal(resp.Bytes(), &category)
	return
}
