package gowprest

import (
	"encoding/json"
	"strconv"
	"strings"
	"time"
)

type Page struct {
	Date              *Date            `json:"date,omitempty"`
	DateGMT           *Date            `json:"date_gmt,omitempty"`
	GUID              *Object          `json:"guid,omitempty"`
	ID                int              `json:"id,omitempty"`
	Link              string           `json:"link,omitempty"`
	Modified          *Date            `json:"modified,omitempty"`
	ModifiedGMT       *Date            `json:"modified_gmt,omitempty"`
	Slug              string           `json:"slug,omitempty"`
	Status            PostStatus       `json:"status,omitempty"`
	Type              string           `json:"type,omitempty"`
	Password          string           `json:"password,omitempty"`
	PermalinkTemplate string           `json:"permalink_template,omitempty"`
	GeneratedSlug     string           `json:"generated_slug,omitempty"`
	Title             *Object          `json:"title,omitempty"`
	Content           *Object          `json:"content,omitempty"`
	Author            int              `json:"author,omitempty"`
	Excerpt           *Object          `json:"excerpt,omitempty"`
	FeaturedMedia     int              `json:"featured_media,omitempty"`
	CommentStatus     OpenClosedStatus `json:"comment_status,omitempty"`
	PingStatus        OpenClosedStatus `json:"ping_status,omitempty"`
	MenuOrder         int              `json:"menu_order,omitempty"`
	Meta              map[string]any   `json:"meta,omitempty"`
	Template          string           `json:"template,omitempty"`
	Parent            int              `json:"parent,omitempty"`
}

type PageData struct {
	ID            int              `json:"id,omitempty"`
	Date          *Date            `json:"date,omitempty"`
	DateGMT       *Date            `json:"date_gmt,omitempty"`
	Slug          string           `json:"slug,omitempty"`
	Status        PostStatus       `json:"status,omitempty"`
	Password      string           `json:"password,omitempty"`
	Title         string           `json:"title,omitempty"`
	Content       string           `json:"content,omitempty"`
	Author        int              `json:"author,omitempty"`
	Excerpt       string           `json:"excerpt,omitempty"`
	FeaturedMedia int              `json:"featured_media,omitempty"`
	CommentStatus OpenClosedStatus `json:"comment_status,omitempty"`
	PingStatus    OpenClosedStatus `json:"ping_status,omitempty"`
	MenuOrder     int              `json:"menu_order,omitempty"`
	Meta          map[string]any   `json:"meta,omitempty"`
	Template      string           `json:"template,omitempty"`
	Parent        int              `json:"parent,omitempty"`
}

type Pages struct {
	client *RestClient
}

func (c *RestClient) Pages() *Pages {
	return &Pages{client: c}
}

type ListPages struct {
	endpoint  string
	client    *RestClient
	arguments map[string]string
}

func (api *Pages) List() *ListPages {
	return &ListPages{
		endpoint:  "/wp/v2/pages",
		client:    api.client,
		arguments: make(map[string]string),
	}
}

func (api *ListPages) ContextView() *ListPages {
	api.arguments["context"] = "view"
	return api
}

func (api *ListPages) ContextEdit() *ListPages {
	api.arguments["context"] = "edit"
	return api
}

func (api *ListPages) ContextEmbed() *ListPages {
	api.arguments["context"] = "embed"
	return api
}

func (api *ListPages) Page(page int) *ListPages {
	api.arguments["page"] = strconv.Itoa(page)
	return api
}

func (api *ListPages) PerPage(perPage int) *ListPages {
	api.arguments["per_page"] = strconv.Itoa(perPage)
	return api
}

func (api *ListPages) Search(query string) *ListPages {
	api.arguments["search"] = query
	return api
}

func (api *ListPages) After(after time.Time) *ListPages {
	api.arguments["after"] = after.Format(time.RFC3339)
	return api
}

func (api *ListPages) ModifiedAfter(modifiedAfter time.Time) *ListPages {
	api.arguments["modified_after"] = modifiedAfter.Format(time.RFC3339)
	return api
}

func (api *ListPages) Author(authorID int) *ListPages {
	api.arguments["author"] = strconv.Itoa(authorID)
	return api
}

func (api *ListPages) AuthorExclude(authorIDs ...int) *ListPages {
	authors := []string{}

	for _, authorId := range authorIDs {
		authors = append(authors, strconv.Itoa(authorId))
	}

	api.arguments["author_exclude"] = strings.Join(authors, ",")
	return api
}

func (api *ListPages) Before(before time.Time) *ListPages {
	api.arguments["before"] = before.Format(time.RFC3339)
	return api
}

func (api *ListPages) ModifiedBefore(modifiedBefore time.Time) *ListPages {
	api.arguments["modified_before"] = modifiedBefore.Format(time.RFC3339)
	return api
}

func (api *ListPages) Exclude(excludeIDs ...int) *ListPages {
	excludes := []string{}
	for _, excludeId := range excludeIDs {
		excludes = append(excludes, strconv.Itoa(excludeId))
	}
	api.arguments["exclude"] = strings.Join(excludes, ",")
	return api
}

func (api *ListPages) Include(includeIDs ...int) *ListPages {
	includes := []string{}
	for _, includeId := range includeIDs {
		includes = append(includes, strconv.Itoa(includeId))
	}
	api.arguments["include"] = strings.Join(includes, ",")
	return api
}

func (api *ListPages) Offset(offset int) *ListPages {
	api.arguments["offset"] = strconv.Itoa(offset)
	return api
}

func (api *ListPages) OrderAsc() *ListPages {
	api.arguments["order"] = "asc"
	return api
}

func (api *ListPages) OrderDesc() *ListPages {
	api.arguments["order"] = "desc"
	return api
}

func (api *ListPages) OrderByAuthor() *ListPages {
	api.arguments["orderby"] = "author"
	return api
}

func (api *ListPages) OrderByDate() *ListPages {
	api.arguments["orderby"] = "date"
	return api
}

func (api *ListPages) OrderById() *ListPages {
	api.arguments["orderby"] = "id"
	return api
}

func (api *ListPages) OrderByInclude() *ListPages {
	api.arguments["orderby"] = "include"
	return api
}

func (api *ListPages) OrderByModified() *ListPages {
	api.arguments["orderby"] = "modified"
	return api
}

func (api *ListPages) OrderByParent() *ListPages {
	api.arguments["orderby"] = "parent"
	return api
}

func (api *ListPages) OrderByRelevance() *ListPages {
	api.arguments["orderby"] = "relevance"
	return api
}

func (api *ListPages) OrderBySlug() *ListPages {
	api.arguments["orderby"] = "slug"
	return api
}

func (api *ListPages) OrderByIncludeSlug() *ListPages {
	api.arguments["orderby"] = "include_slugs"
	return api
}

func (api *ListPages) OrderByTitle() *ListPages {
	api.arguments["orderby"] = "title"
	return api
}

func (api *ListPages) SearchColumns(columns ...string) *ListPages {
	api.arguments["search_columns"] = strings.Join(columns, ",")
	return api
}

func (api *ListPages) Slug(slug string) *ListPages {
	api.arguments["slug"] = slug
	return api
}

func (api *ListPages) StatusPublish() *ListPages {
	api.arguments["status"] = "publish"
	return api
}

func (api *ListPages) StatusDraft() *ListPages {
	api.arguments["status"] = "draft"
	return api
}

func (api *ListPages) StatusPending() *ListPages {
	api.arguments["status"] = "pending"
	return api
}

func (api *ListPages) StatusPrivate() *ListPages {
	api.arguments["status"] = "private"
	return api
}

func (api *ListPages) StatusFuture() *ListPages {
	api.arguments["status"] = "future"
	return api
}

func (api *ListPages) StatusTrash() *ListPages {
	api.arguments["status"] = "trash"
	return api
}

func (api *ListPages) StatusAny() *ListPages {
	api.arguments["status"] = "any"
	return api
}

func (api *ListPages) Parent(parentID int) *ListPages {
	api.arguments["parent"] = strconv.Itoa(parentID)
	return api
}

func (api *ListPages) ParentExclude(parentIDs ...int) *ListPages {
	parents := []string{}
	for _, parentId := range parentIDs {
		parents = append(parents, strconv.Itoa(parentId))
	}
	api.arguments["parent_exclude"] = strings.Join(parents, ",")
	return api
}

func (api *ListPages) Do() (pages []Page, err error) {
	_, err = api.client.httpClient.R().
		SetHeader("Accept", "application/json").
		SetResult(&pages).
		SetQueryParams(api.arguments).
		Get(api.client.endpoint + api.endpoint)

	if err != nil {
		return
	}

	return
}

type CreatePage struct {
	endpoint string
	client   *RestClient
	page     PageData
}

func (api *Pages) Create(page PageData) *CreatePage {
	return &CreatePage{
		endpoint: "/wp/v2/pages",
		client:   api.client,
		page:     page,
	}
}

func (api *CreatePage) Do() (page Page, err error) {
	resp, err := api.client.httpClient.R().
		SetHeader("Content-Type", "application/json").
		SetBasicAuth(api.client.auth.Username, api.client.auth.Password).
		SetResult(&page).
		SetBody(api.page).
		Post(api.client.endpoint + api.endpoint)

	if resp.IsError() {
		var wpError WPRestError
		err = json.Unmarshal(resp.Bytes(), &wpError)

		if err != nil {
			return
		}

		return page, &wpError
	}

	if err != nil {
		return
	}

	return
}

type RetrievePage struct {
	endpoint  string
	client    *RestClient
	arguments map[string]string
}

func (api *Pages) Retrieve(pageID int) *RetrievePage {
	return &RetrievePage{
		endpoint:  "/wp/v2/pages/" + strconv.Itoa(pageID),
		client:    api.client,
		arguments: make(map[string]string),
	}
}

func (api *RetrievePage) ContextView() *RetrievePage {
	api.arguments["context"] = "view"
	return api
}

func (api *RetrievePage) ContextEdit() *RetrievePage {
	api.arguments["context"] = "edit"
	return api
}

func (api *RetrievePage) ContextEmbed() *RetrievePage {
	api.arguments["context"] = "embed"
	return api
}

func (api *RetrievePage) Password(password string) *RetrievePage {
	api.arguments["password"] = password
	return api
}

func (api *RetrievePage) Do() (page *Page, err error) {
	endpoint := api.client.endpoint + api.endpoint

	restyClient := api.client.httpClient.R()
	if api.client.auth.Username != "" && api.client.auth.Password != "" && api.arguments["context"] == "edit" {
		restyClient.SetBasicAuth(api.client.auth.Username, api.client.auth.Password)
	}

	resp, err := restyClient.
		SetHeader("Accept", "application/json").
		SetResult(&page).
		SetQueryParams(api.arguments).
		Get(endpoint)

	if resp.IsError() {
		var wpError WPRestError
		err = json.Unmarshal(resp.Bytes(), &wpError)

		if err != nil {
			return
		}

		return page, &wpError
	}

	// TODO: need fixing of message = invalid suit value: trash
	if err != nil && err.Error() == "invalid suit value: trash" {
		err = nil
	}

	if err != nil {
		return
	}

	return
}

type UpdatePage struct {
	endpoint string
	client   *RestClient
	page     PageData
}

func (api *Pages) Update(page PageData) *UpdatePage {
	return &UpdatePage{
		endpoint: "/wp/v2/pages/" + strconv.Itoa(page.ID),
		client:   api.client,
		page:     page,
	}
}

func (api *UpdatePage) Do() (page Page, err error) {
	resp, err := api.client.httpClient.R().
		SetHeader("Content-Type", "application/json").
		SetBasicAuth(api.client.auth.Username, api.client.auth.Password).
		SetResult(&page).
		SetBody(api.page).
		Post(api.client.endpoint + api.endpoint)

	if resp.IsError() {
		var wpError WPRestError
		err = json.Unmarshal(resp.Bytes(), &wpError)

		if err != nil {
			return
		}

		return page, &wpError
	}

	if err != nil {
		return
	}

	return
}

type DeletePage struct {
	endpoint string
	client   *RestClient
	pageID   int
	force    bool
}

func (api *Pages) Delete(pageID int) *DeletePage {
	return &DeletePage{
		endpoint: "/wp/v2/pages",
		client:   api.client,
		pageID:   pageID,
	}
}

func (api *DeletePage) Force() *DeletePage {
	api.force = true
	return api
}

func (api *DeletePage) Do() (page Page, err error) {
	endpoint := api.client.endpoint + api.endpoint + "/" + strconv.Itoa(api.pageID)
	resp, err :=
		api.client.httpClient.R().
			SetHeader("Content-Type", "application/json").
			SetBasicAuth(api.client.auth.Username, api.client.auth.Password).
			SetResult(&page).
			SetPathParam("force", strconv.FormatBool(api.force)).
			Delete(endpoint)

	if resp.IsError() {
		var wpError WPRestError
		err = json.Unmarshal(resp.Bytes(), &wpError)

		if err != nil {
			return
		}

		return page, &wpError
	}

	// TODO: need fixing of message = invalid suit value: trash
	if err != nil && err.Error() == "invalid suit value: trash" {
		err = nil
	}

	return
}
