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
	Links             map[string]any   `json:"_links,omitempty"`
	Embedded          map[string]any   `json:"_embedded,omitempty"`
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

func (api *Pages) Revisions(parentID int) *PageRevisions {
	return &PageRevisions{
		client:   api.client,
		parentID: parentID,
	}
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

func (api *ListPages) Context(ctx string) *ListPages {
	api.arguments["context"] = ctx
	return api
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

func (api *ListPages) Author(authorIDs ...int) *ListPages {
	authors := make([]string, 0, len(authorIDs))
	for _, id := range authorIDs {
		authors = append(authors, strconv.Itoa(id))
	}
	api.arguments["author"] = strings.Join(authors, ",")
	return api
}

func (api *ListPages) AuthorExclude(authorIDs ...int) *ListPages {
	authors := make([]string, 0, len(authorIDs))
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
	excludes := make([]string, 0, len(excludeIDs))
	for _, excludeId := range excludeIDs {
		excludes = append(excludes, strconv.Itoa(excludeId))
	}
	api.arguments["exclude"] = strings.Join(excludes, ",")
	return api
}

func (api *ListPages) Include(includeIDs ...int) *ListPages {
	includes := make([]string, 0, len(includeIDs))
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

func (api *ListPages) Order(order string) *ListPages {
	api.arguments["order"] = order
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

func (api *ListPages) OrderBy(orderBy string) *ListPages {
	api.arguments["orderby"] = orderBy
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

func (api *ListPages) OrderByMenuOrder() *ListPages {
	api.arguments["orderby"] = "menu_order"
	return api
}

func (api *ListPages) MenuOrder(menuOrder int) *ListPages {
	api.arguments["menu_order"] = strconv.Itoa(menuOrder)
	return api
}

func (api *ListPages) SearchColumns(columns ...string) *ListPages {
	api.arguments["search_columns"] = strings.Join(columns, ",")
	return api
}

func (api *ListPages) Slug(slugs ...string) *ListPages {
	api.arguments["slug"] = strings.Join(slugs, ",")
	return api
}

func (api *ListPages) Status(statuses ...string) *ListPages {
	api.arguments["status"] = strings.Join(statuses, ",")
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
	parents := make([]string, 0, len(parentIDs))
	for _, parentId := range parentIDs {
		parents = append(parents, strconv.Itoa(parentId))
	}
	api.arguments["parent_exclude"] = strings.Join(parents, ",")
	return api
}

func (api *ListPages) Embed() *ListPages {
	api.arguments["_embed"] = "true"
	return api
}

func (api *ListPages) Fields(fields ...string) *ListPages {
	api.arguments["_fields"] = strings.Join(fields, ",")
	return api
}

func (api *ListPages) Do() (pages []Page, err error) {
	restyClient := api.client.httpClient.R()
	if api.client.auth.Username != "" && api.client.auth.Password != "" {
		restyClient.SetBasicAuth(api.client.auth.Username, api.client.auth.Password)
	}

	resp, err := restyClient.
		SetHeader("Accept", "application/json").
		SetResult(&pages).
		SetQueryParams(api.arguments).
		Get(api.client.endpoint + api.endpoint)

	if resp.IsError() {
		var wpError WPRestError
		err = json.Unmarshal(resp.Bytes(), &wpError)
		if err != nil {
			return
		}
		return pages, &wpError
	}

	return
}

type CreatePage struct {
	endpoint string
	client   *RestClient
	page     PageData
}

// Create returns a CreatePage builder.
func (api *Pages) Create(page ...PageData) *CreatePage {
	builder := &CreatePage{
		endpoint: "/wp/v2/pages",
		client:   api.client,
	}
	if len(page) > 0 {
		builder.page = page[0]
	}
	return builder
}

func (api *CreatePage) Title(title string) *CreatePage {
	api.page.Title = title
	return api
}

func (api *CreatePage) Content(content string) *CreatePage {
	api.page.Content = content
	return api
}

func (api *CreatePage) Excerpt(excerpt string) *CreatePage {
	api.page.Excerpt = excerpt
	return api
}

func (api *CreatePage) Slug(slug string) *CreatePage {
	api.page.Slug = slug
	return api
}

func (api *CreatePage) Status(status PostStatus) *CreatePage {
	api.page.Status = status
	return api
}

func (api *CreatePage) StatusPublish() *CreatePage {
	api.page.Status = StatusPublished
	return api
}

func (api *CreatePage) StatusDraft() *CreatePage {
	api.page.Status = StatusDraft
	return api
}

func (api *CreatePage) StatusPending() *CreatePage {
	api.page.Status = StatusPending
	return api
}

func (api *CreatePage) StatusPrivate() *CreatePage {
	api.page.Status = StatusPrivate
	return api
}

func (api *CreatePage) StatusFuture() *CreatePage {
	api.page.Status = StatusFuture
	return api
}

func (api *CreatePage) Password(password string) *CreatePage {
	api.page.Password = password
	return api
}

func (api *CreatePage) Author(authorID int) *CreatePage {
	api.page.Author = authorID
	return api
}

func (api *CreatePage) FeaturedMedia(mediaID int) *CreatePage {
	api.page.FeaturedMedia = mediaID
	return api
}

func (api *CreatePage) CommentStatus(status OpenClosedStatus) *CreatePage {
	api.page.CommentStatus = status
	return api
}

func (api *CreatePage) PingStatus(status OpenClosedStatus) *CreatePage {
	api.page.PingStatus = status
	return api
}

func (api *CreatePage) MenuOrder(menuOrder int) *CreatePage {
	api.page.MenuOrder = menuOrder
	return api
}

func (api *CreatePage) Parent(parentID int) *CreatePage {
	api.page.Parent = parentID
	return api
}

func (api *CreatePage) Template(template string) *CreatePage {
	api.page.Template = template
	return api
}

func (api *CreatePage) Meta(meta map[string]any) *CreatePage {
	api.page.Meta = meta
	return api
}

func (api *CreatePage) Date(date time.Time) *CreatePage {
	api.page.Date = &Date{Time: date}
	return api
}

func (api *CreatePage) DateGMT(date time.Time) *CreatePage {
	api.page.DateGMT = &Date{Time: date}
	return api
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

func (api *RetrievePage) Context(ctx string) *RetrievePage {
	api.arguments["context"] = ctx
	return api
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

func (api *RetrievePage) Embed() *RetrievePage {
	api.arguments["_embed"] = "true"
	return api
}

func (api *RetrievePage) Fields(fields ...string) *RetrievePage {
	api.arguments["_fields"] = strings.Join(fields, ",")
	return api
}

func (api *RetrievePage) Do() (page *Page, err error) {
	endpoint := api.client.endpoint + api.endpoint

	restyClient := api.client.httpClient.R()
	if api.client.auth.Username != "" && api.client.auth.Password != "" {
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

	return
}

type UpdatePage struct {
	endpoint string
	client   *RestClient
	page     PageData
}

// Update returns an UpdatePage builder.
func (api *Pages) Update(page ...PageData) *UpdatePage {
	builder := &UpdatePage{
		endpoint: "/wp/v2/pages",
		client:   api.client,
	}
	if len(page) > 0 {
		builder.page = page[0]
		if builder.page.ID != 0 {
			builder.endpoint = "/wp/v2/pages/" + strconv.Itoa(builder.page.ID)
		}
	}
	return builder
}

func (api *UpdatePage) ID(id int) *UpdatePage {
	api.page.ID = id
	api.endpoint = "/wp/v2/pages/" + strconv.Itoa(id)
	return api
}

func (api *UpdatePage) Title(title string) *UpdatePage {
	api.page.Title = title
	return api
}

func (api *UpdatePage) Content(content string) *UpdatePage {
	api.page.Content = content
	return api
}

func (api *UpdatePage) Excerpt(excerpt string) *UpdatePage {
	api.page.Excerpt = excerpt
	return api
}

func (api *UpdatePage) Slug(slug string) *UpdatePage {
	api.page.Slug = slug
	return api
}

func (api *UpdatePage) Status(status PostStatus) *UpdatePage {
	api.page.Status = status
	return api
}

func (api *UpdatePage) StatusPublish() *UpdatePage {
	api.page.Status = StatusPublished
	return api
}

func (api *UpdatePage) StatusDraft() *UpdatePage {
	api.page.Status = StatusDraft
	return api
}

func (api *UpdatePage) StatusPending() *UpdatePage {
	api.page.Status = StatusPending
	return api
}

func (api *UpdatePage) StatusPrivate() *UpdatePage {
	api.page.Status = StatusPrivate
	return api
}

func (api *UpdatePage) StatusFuture() *UpdatePage {
	api.page.Status = StatusFuture
	return api
}

func (api *UpdatePage) Password(password string) *UpdatePage {
	api.page.Password = password
	return api
}

func (api *UpdatePage) Author(authorID int) *UpdatePage {
	api.page.Author = authorID
	return api
}

func (api *UpdatePage) FeaturedMedia(mediaID int) *UpdatePage {
	api.page.FeaturedMedia = mediaID
	return api
}

func (api *UpdatePage) CommentStatus(status OpenClosedStatus) *UpdatePage {
	api.page.CommentStatus = status
	return api
}

func (api *UpdatePage) PingStatus(status OpenClosedStatus) *UpdatePage {
	api.page.PingStatus = status
	return api
}

func (api *UpdatePage) MenuOrder(menuOrder int) *UpdatePage {
	api.page.MenuOrder = menuOrder
	return api
}

func (api *UpdatePage) Parent(parentID int) *UpdatePage {
	api.page.Parent = parentID
	return api
}

func (api *UpdatePage) Template(template string) *UpdatePage {
	api.page.Template = template
	return api
}

func (api *UpdatePage) Meta(meta map[string]any) *UpdatePage {
	api.page.Meta = meta
	return api
}

func (api *UpdatePage) Date(date time.Time) *UpdatePage {
	api.page.Date = &Date{Time: date}
	return api
}

func (api *UpdatePage) DateGMT(date time.Time) *UpdatePage {
	api.page.DateGMT = &Date{Time: date}
	return api
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

	return
}

type deletePageEnvelope struct {
	Page
	Deleted  bool  `json:"deleted"`
	Previous *Page `json:"previous"`
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
	var env deletePageEnvelope
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
		return page, &wpError
	}

	if err != nil {
		return
	}

	if env.Previous != nil && env.Previous.ID != 0 {
		return *env.Previous, nil
	}

	return env.Page, nil
}
