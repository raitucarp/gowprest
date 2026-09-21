package gowprest

import (
	"encoding/json"
	"strconv"
	"strings"
	"time"
)

type Object struct {
	Rendered     string `json:"rendered"`
	Raw          string `json:"raw,omitempty"`
	Protected    bool   `json:"protected,omitempty,omitzero"`
	BlockVersion int    `json:"block_version,omitempty"`
}

func (o *Object) UnmarshalJSON(b []byte) error {
	trimmed := strings.TrimSpace(string(b))
	if trimmed == "" || trimmed == "null" || trimmed == "[]" {
		return nil
	}
	type Alias Object
	var aux Alias
	if err := json.Unmarshal(b, &aux); err != nil {
		return err
	}
	*o = Object(aux)
	return nil
}

type Date struct {
	time.Time
}

func (ct *Date) UnmarshalJSON(b []byte) error {
	s := strings.Trim(string(b), "\"")
	if s == "" || s == "null" {
		return nil
	}
	t, err := time.Parse(time.RFC3339, s)
	if err == nil {
		ct.Time = t
		return nil
	}
	t, err = time.Parse("2006-01-02T15:04:05", s)
	if err != nil {
		return err
	}
	ct.Time = t
	return nil
}

// Post represents a WordPress post resource returned by the REST API.
type Post struct {
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
	Format            Format           `json:"format,omitempty"`
	Meta              map[string]any   `json:"meta,omitempty"`
	Sticky            bool             `json:"sticky,omitempty"`
	Template          string           `json:"template,omitempty"`
	Categories        []int            `json:"categories,omitempty"`
	Tags              []int            `json:"tags,omitempty"`
	Links             map[string]any   `json:"_links,omitempty"`
	Embedded          map[string]any   `json:"_embedded,omitempty"`
}

// PostData holds mutable payload data for creating or updating a WordPress post.
type PostData struct {
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
	Format        Format           `json:"format,omitempty"`
	Meta          map[string]any   `json:"meta,omitempty"`
	Sticky        bool             `json:"sticky,omitempty"`
	Template      string           `json:"template,omitempty"`
	Categories    []int            `json:"categories,omitempty"`
	Tags          []int            `json:"tags,omitempty"`
}

// Posts provides operations for managing WordPress posts (/wp/v2/posts).
type Posts struct {
	client *RestClient
}

// Posts returns the Posts service for querying, creating, retrieving, updating, and deleting posts.
func (c *RestClient) Posts() *Posts {
	return &Posts{client: c}
}

// Revisions returns the PostRevisions service for inspecting revisions of a specific post.
func (api *Posts) Revisions(parentID int) *PostRevisions {
	return &PostRevisions{
		client:   api.client,
		parentID: parentID,
	}
}

// ListPosts is a fluent builder for querying a collection of posts.
type ListPosts struct {
	endpoint  string
	client    *RestClient
	arguments map[string]string
}

// List initiates a query builder to list and filter posts.
func (api *Posts) List() *ListPosts {
	return &ListPosts{
		endpoint:  "/wp/v2/posts",
		client:    api.client,
		arguments: make(map[string]string),
	}
}

func (api *ListPosts) ContextView() *ListPosts {
	api.arguments["context"] = "view"
	return api
}

func (api *ListPosts) ContextEdit() *ListPosts {
	api.arguments["context"] = "edit"
	return api
}

func (api *ListPosts) ContextEmbed() *ListPosts {
	api.arguments["context"] = "embed"
	return api
}

func (api *ListPosts) Page(page int) *ListPosts {
	api.arguments["page"] = strconv.Itoa(page)
	return api
}

func (api *ListPosts) PerPage(perPage int) *ListPosts {
	api.arguments["per_page"] = strconv.Itoa(perPage)
	return api
}

func (api *ListPosts) Search(query string) *ListPosts {
	api.arguments["search"] = query
	return api
}

func (api *ListPosts) After(after time.Time) *ListPosts {
	api.arguments["after"] = after.Format(time.RFC3339)
	return api
}

func (api *ListPosts) ModifiedAfter(modifiedAfter time.Time) *ListPosts {
	api.arguments["modified_after"] = modifiedAfter.Format(time.RFC3339)
	return api
}

func (api *ListPosts) Author(authorIDs ...int) *ListPosts {
	authors := []string{}
	for _, authorID := range authorIDs {
		authors = append(authors, strconv.Itoa(authorID))
	}
	api.arguments["author"] = strings.Join(authors, ",")
	return api
}

func (api *ListPosts) AuthorExclude(authorIDs ...int) *ListPosts {
	authors := []string{}

	for _, authorId := range authorIDs {
		authors = append(authors, strconv.Itoa(authorId))
	}

	api.arguments["author_exclude"] = strings.Join(authors, ",")
	return api
}

func (api *ListPosts) Before(before time.Time) *ListPosts {
	api.arguments["before"] = before.Format(time.RFC3339)
	return api
}

func (api *ListPosts) ModifiedBefore(modifiedBefore time.Time) *ListPosts {
	api.arguments["modified_before"] = modifiedBefore.Format(time.RFC3339)
	return api
}

func (api *ListPosts) Exclude(excludeIDs ...int) *ListPosts {
	excludes := []string{}
	for _, excludeId := range excludeIDs {
		excludes = append(excludes, strconv.Itoa(excludeId))
	}
	api.arguments["exclude"] = strings.Join(excludes, ",")
	return api
}

func (api *ListPosts) Include(includeIDs ...int) *ListPosts {
	includes := []string{}
	for _, includeId := range includeIDs {
		includes = append(includes, strconv.Itoa(includeId))
	}
	api.arguments["include"] = strings.Join(includes, ",")
	return api
}

func (api *ListPosts) Offset(offset int) *ListPosts {
	api.arguments["offset"] = strconv.Itoa(offset)
	return api
}

func (api *ListPosts) OrderAsc() *ListPosts {
	api.arguments["order"] = "asc"
	return api
}

func (api *ListPosts) OrderDesc() *ListPosts {
	api.arguments["order"] = "desc"
	return api
}

func (api *ListPosts) Order(order string) *ListPosts {
	api.arguments["order"] = order
	return api
}

func (api *ListPosts) OrderByAuthor() *ListPosts {
	api.arguments["orderby"] = "author"
	return api
}

func (api *ListPosts) OrderByDate() *ListPosts {
	api.arguments["orderby"] = "date"
	return api
}

func (api *ListPosts) OrderById() *ListPosts {
	api.arguments["orderby"] = "id"
	return api
}

func (api *ListPosts) OrderByInclude() *ListPosts {
	api.arguments["orderby"] = "include"
	return api
}

func (api *ListPosts) OrderByModified() *ListPosts {
	api.arguments["orderby"] = "modified"
	return api
}

func (api *ListPosts) OrderByParent() *ListPosts {
	api.arguments["orderby"] = "parent"
	return api
}

func (api *ListPosts) OrderByRelevance() *ListPosts {
	api.arguments["orderby"] = "relevance"
	return api
}

func (api *ListPosts) OrderBySlug() *ListPosts {
	api.arguments["orderby"] = "slug"
	return api
}

func (api *ListPosts) OrderByIncludeSlug() *ListPosts {
	api.arguments["orderby"] = "include_slugs"
	return api
}

func (api *ListPosts) OrderByTitle() *ListPosts {
	api.arguments["orderby"] = "title"
	return api
}

func (api *ListPosts) OrderBy(orderby string) *ListPosts {
	api.arguments["orderby"] = orderby
	return api
}

func (api *ListPosts) SearchColumns(columns ...string) *ListPosts {
	api.arguments["search_columns"] = strings.Join(columns, ",")
	return api
}

func (api *ListPosts) Slug(slugs ...string) *ListPosts {
	api.arguments["slug"] = strings.Join(slugs, ",")
	return api
}

func (api *ListPosts) StatusPublish() *ListPosts {
	api.arguments["status"] = "publish"
	return api
}

func (api *ListPosts) StatusDraft() *ListPosts {
	api.arguments["status"] = "draft"
	return api
}

func (api *ListPosts) StatusPending() *ListPosts {
	api.arguments["status"] = "pending"
	return api
}

func (api *ListPosts) StatusPrivate() *ListPosts {
	api.arguments["status"] = "private"
	return api
}

func (api *ListPosts) StatusFuture() *ListPosts {
	api.arguments["status"] = "future"
	return api
}

func (api *ListPosts) StatusTrash() *ListPosts {
	api.arguments["status"] = "trash"
	return api
}

func (api *ListPosts) StatusAny() *ListPosts {
	api.arguments["status"] = "any"
	return api
}

func (api *ListPosts) Status(statuses ...string) *ListPosts {
	api.arguments["status"] = strings.Join(statuses, ",")
	return api
}

func (api *ListPosts) TaxAnd() *ListPosts {
	api.arguments["tax_relation"] = "AND"
	return api
}

func (api *ListPosts) TaxOr() *ListPosts {
	api.arguments["tax_relation"] = "OR"
	return api
}

func (api *ListPosts) Categories(categories ...string) *ListPosts {
	api.arguments["categories"] = strings.Join(categories, ",")
	return api
}

func (api *ListPosts) CategoryIDs(categoryIDs ...int) *ListPosts {
	ids := []string{}
	for _, id := range categoryIDs {
		ids = append(ids, strconv.Itoa(id))
	}
	api.arguments["categories"] = strings.Join(ids, ",")
	return api
}

func (api *ListPosts) CategoriesExclude(categories ...string) *ListPosts {
	api.arguments["categories_exclude"] = strings.Join(categories, ",")
	return api
}

func (api *ListPosts) CategoryIDsExclude(categoryIDs ...int) *ListPosts {
	ids := []string{}
	for _, id := range categoryIDs {
		ids = append(ids, strconv.Itoa(id))
	}
	api.arguments["categories_exclude"] = strings.Join(ids, ",")
	return api
}

func (api *ListPosts) Tags(tags ...string) *ListPosts {
	api.arguments["tags"] = strings.Join(tags, ",")
	return api
}

func (api *ListPosts) TagIDs(tagIDs ...int) *ListPosts {
	ids := []string{}
	for _, id := range tagIDs {
		ids = append(ids, strconv.Itoa(id))
	}
	api.arguments["tags"] = strings.Join(ids, ",")
	return api
}

func (api *ListPosts) TagsExclude(tags ...string) *ListPosts {
	api.arguments["tags_exclude"] = strings.Join(tags, ",")
	return api
}

func (api *ListPosts) TagIDsExclude(tagIDs ...int) *ListPosts {
	ids := []string{}
	for _, id := range tagIDs {
		ids = append(ids, strconv.Itoa(id))
	}
	api.arguments["tags_exclude"] = strings.Join(ids, ",")
	return api
}

func (api *ListPosts) Sticky(sticky bool) *ListPosts {
	api.arguments["sticky"] = strconv.FormatBool(sticky)
	return api
}

func (api *ListPosts) Embed() *ListPosts {
	api.arguments["_embed"] = "true"
	return api
}

func (api *ListPosts) Fields(fields ...string) *ListPosts {
	api.arguments["_fields"] = strings.Join(fields, ",")
	return api
}

func (api *ListPosts) Do() (posts []Post, err error) {
	restyClient := api.client.httpClient.R()
	if api.client.auth.Username != "" && api.client.auth.Password != "" && (api.arguments["context"] == "edit" || api.arguments["status"] != "") {
		restyClient.SetBasicAuth(api.client.auth.Username, api.client.auth.Password)
	}

	resp, err := restyClient.
		SetHeader("Accept", "application/json").
		SetResult(&posts).
		SetQueryParams(api.arguments).
		Get(api.client.endpoint + api.endpoint)

	if resp.IsError() {
		var wpError WPRestError
		err = json.Unmarshal(resp.Bytes(), &wpError)
		if err != nil {
			return
		}
		return posts, &wpError
	}

	if err != nil {
		return
	}

	return
}

// CreatePost is a fluent builder for creating a new post.
type CreatePost struct {
	endpoint string
	client   *RestClient
	post     PostData
}

// Create initializes a fluent builder to create a new WordPress post.
func (api *Posts) Create(post ...PostData) *CreatePost {
	var p PostData
	if len(post) > 0 {
		p = post[0]
	}
	return &CreatePost{
		endpoint: "/wp/v2/posts",
		client:   api.client,
		post:     p,
	}
}

func (api *CreatePost) Title(title string) *CreatePost                 { api.post.Title = title; return api }
func (api *CreatePost) Content(content string) *CreatePost             { api.post.Content = content; return api }
func (api *CreatePost) Excerpt(excerpt string) *CreatePost             { api.post.Excerpt = excerpt; return api }
func (api *CreatePost) Status(status PostStatus) *CreatePost           { api.post.Status = status; return api }
func (api *CreatePost) StatusPublish() *CreatePost                     { api.post.Status = StatusPublished; return api }
func (api *CreatePost) StatusDraft() *CreatePost                       { api.post.Status = StatusDraft; return api }
func (api *CreatePost) StatusPending() *CreatePost                     { api.post.Status = StatusPending; return api }
func (api *CreatePost) StatusPrivate() *CreatePost                     { api.post.Status = StatusPrivate; return api }
func (api *CreatePost) StatusFuture() *CreatePost                      { api.post.Status = StatusFuture; return api }
func (api *CreatePost) Slug(slug string) *CreatePost                   { api.post.Slug = slug; return api }
func (api *CreatePost) Password(password string) *CreatePost           { api.post.Password = password; return api }
func (api *CreatePost) Author(authorID int) *CreatePost                { api.post.Author = authorID; return api }
func (api *CreatePost) FeaturedMedia(mediaID int) *CreatePost          { api.post.FeaturedMedia = mediaID; return api }
func (api *CreatePost) CommentStatus(status OpenClosedStatus) *CreatePost {
	api.post.CommentStatus = status
	return api
}
func (api *CreatePost) PingStatus(status OpenClosedStatus) *CreatePost {
	api.post.PingStatus = status
	return api
}
func (api *CreatePost) Format(format Format) *CreatePost         { api.post.Format = format; return api }
func (api *CreatePost) Meta(meta map[string]any) *CreatePost     { api.post.Meta = meta; return api }
func (api *CreatePost) Sticky(sticky bool) *CreatePost           { api.post.Sticky = sticky; return api }
func (api *CreatePost) Template(template string) *CreatePost     { api.post.Template = template; return api }
func (api *CreatePost) Categories(categories ...int) *CreatePost {
	api.post.Categories = append(api.post.Categories, categories...)
	return api
}
func (api *CreatePost) Tags(tags ...int) *CreatePost {
	api.post.Tags = append(api.post.Tags, tags...)
	return api
}
func (api *CreatePost) Date(t time.Time) *CreatePost    { api.post.Date = &Date{Time: t}; return api }
func (api *CreatePost) DateGMT(t time.Time) *CreatePost { api.post.DateGMT = &Date{Time: t}; return api }

func (api *CreatePost) Do() (post Post, err error) {
	resp, err := api.client.httpClient.R().
		SetHeader("Content-Type", "application/json").
		SetBasicAuth(api.client.auth.Username, api.client.auth.Password).
		SetResult(&post).
		SetBody(api.post).
		Post(api.client.endpoint + api.endpoint)

	if resp.IsError() {
		var wpError WPRestError
		err = json.Unmarshal(resp.Bytes(), &wpError)

		if err != nil {
			return
		}

		return post, &wpError
	}

	if err != nil {
		return
	}

	return
}

// RetrievePost is a fluent builder for fetching a single post by ID.
type RetrievePost struct {
	endpoint  string
	client    *RestClient
	arguments map[string]string
}

// Retrieve returns a builder to fetch the post with the given ID.
func (api *Posts) Retrieve(postId int) *RetrievePost {
	return &RetrievePost{
		endpoint:  "/wp/v2/posts/" + strconv.Itoa(postId),
		client:    api.client,
		arguments: make(map[string]string),
	}
}

func (api *RetrievePost) ContextView() *RetrievePost {
	api.arguments["context"] = "view"
	return api
}

func (api *RetrievePost) ContextEdit() *RetrievePost {
	api.arguments["context"] = "edit"
	return api
}

func (api *RetrievePost) ContextEmbed() *RetrievePost {
	api.arguments["context"] = "embed"
	return api
}

func (api *RetrievePost) Password(password string) *RetrievePost {
	api.arguments["password"] = password
	return api
}

func (api *RetrievePost) Embed() *RetrievePost {
	api.arguments["_embed"] = "true"
	return api
}

func (api *RetrievePost) Fields(fields ...string) *RetrievePost {
	api.arguments["_fields"] = strings.Join(fields, ",")
	return api
}

func (api *RetrievePost) Do() (post *Post, err error) {
	endpoint := api.client.endpoint + api.endpoint

	restyClient := api.client.httpClient.R()
	if api.client.auth.Username != "" && api.client.auth.Password != "" && (api.arguments["context"] == "edit" || api.arguments["password"] != "") {
		restyClient.SetBasicAuth(api.client.auth.Username, api.client.auth.Password)
	}

	resp, err := restyClient.
		SetHeader("Accept", "application/json").
		SetResult(&post).
		SetQueryParams(api.arguments).
		Get(endpoint)

	if resp.IsError() {
		var wpError WPRestError
		err = json.Unmarshal(resp.Bytes(), &wpError)

		if err != nil {
			return
		}

		return post, &wpError
	}

	if err != nil {
		return
	}

	return
}

// UpdatePost is a fluent builder for modifying an existing post.
type UpdatePost struct {
	endpoint string
	client   *RestClient
	post     PostData
}

// Update returns a builder to modify an existing post.
func (api *Posts) Update(post ...PostData) *UpdatePost {
	var p PostData
	endpoint := "/wp/v2/posts"
	if len(post) > 0 {
		p = post[0]
		if p.ID != 0 {
			endpoint = "/wp/v2/posts/" + strconv.Itoa(p.ID)
		}
	}
	return &UpdatePost{
		endpoint: endpoint,
		client:   api.client,
		post:     p,
	}
}

func (api *UpdatePost) ID(id int) *UpdatePost {
	api.post.ID = id
	api.endpoint = "/wp/v2/posts/" + strconv.Itoa(id)
	return api
}

func (api *UpdatePost) Title(title string) *UpdatePost                 { api.post.Title = title; return api }
func (api *UpdatePost) Content(content string) *UpdatePost             { api.post.Content = content; return api }
func (api *UpdatePost) Excerpt(excerpt string) *UpdatePost             { api.post.Excerpt = excerpt; return api }
func (api *UpdatePost) Status(status PostStatus) *UpdatePost           { api.post.Status = status; return api }
func (api *UpdatePost) StatusPublish() *UpdatePost                     { api.post.Status = StatusPublished; return api }
func (api *UpdatePost) StatusDraft() *UpdatePost                       { api.post.Status = StatusDraft; return api }
func (api *UpdatePost) StatusPending() *UpdatePost                     { api.post.Status = StatusPending; return api }
func (api *UpdatePost) StatusPrivate() *UpdatePost                     { api.post.Status = StatusPrivate; return api }
func (api *UpdatePost) StatusFuture() *UpdatePost                      { api.post.Status = StatusFuture; return api }
func (api *UpdatePost) Slug(slug string) *UpdatePost                   { api.post.Slug = slug; return api }
func (api *UpdatePost) Password(password string) *UpdatePost           { api.post.Password = password; return api }
func (api *UpdatePost) Author(authorID int) *UpdatePost                { api.post.Author = authorID; return api }
func (api *UpdatePost) FeaturedMedia(mediaID int) *UpdatePost          { api.post.FeaturedMedia = mediaID; return api }
func (api *UpdatePost) CommentStatus(status OpenClosedStatus) *UpdatePost {
	api.post.CommentStatus = status
	return api
}
func (api *UpdatePost) PingStatus(status OpenClosedStatus) *UpdatePost {
	api.post.PingStatus = status
	return api
}
func (api *UpdatePost) Format(format Format) *UpdatePost         { api.post.Format = format; return api }
func (api *UpdatePost) Meta(meta map[string]any) *UpdatePost     { api.post.Meta = meta; return api }
func (api *UpdatePost) Sticky(sticky bool) *UpdatePost           { api.post.Sticky = sticky; return api }
func (api *UpdatePost) Template(template string) *UpdatePost     { api.post.Template = template; return api }
func (api *UpdatePost) Categories(categories ...int) *UpdatePost {
	api.post.Categories = append(api.post.Categories, categories...)
	return api
}
func (api *UpdatePost) Tags(tags ...int) *UpdatePost {
	api.post.Tags = append(api.post.Tags, tags...)
	return api
}
func (api *UpdatePost) Date(t time.Time) *UpdatePost    { api.post.Date = &Date{Time: t}; return api }
func (api *UpdatePost) DateGMT(t time.Time) *UpdatePost { api.post.DateGMT = &Date{Time: t}; return api }

func (api *UpdatePost) Do() (post Post, err error) {
	resp, err := api.client.httpClient.R().
		SetHeader("Content-Type", "application/json").
		SetBasicAuth(api.client.auth.Username, api.client.auth.Password).
		SetResult(&post).
		SetBody(api.post).
		Post(api.client.endpoint + api.endpoint)

	if resp.IsError() {
		var wpError WPRestError
		err = json.Unmarshal(resp.Bytes(), &wpError)

		if err != nil {
			return
		}

		return post, &wpError
	}

	if err != nil {
		return
	}

	return
}

// DeletePost is a fluent builder for deleting a post.
type DeletePost struct {
	endpoint string
	client   *RestClient
	postId   int
	force    bool
}

// Delete returns a builder to delete the post with the given ID.
func (api *Posts) Delete(postId int) *DeletePost {
	return &DeletePost{
		endpoint: "/wp/v2/posts",
		client:   api.client,
		postId:   postId,
	}
}

func (api *DeletePost) Force() *DeletePost {
	api.force = true
	return api
}

type deletePostEnvelope struct {
	Post
	Deleted  bool  `json:"deleted"`
	Previous *Post `json:"previous"`
}

func (api *DeletePost) Do() (post Post, err error) {
	endpoint := api.client.endpoint + api.endpoint + "/" + strconv.Itoa(api.postId)
	var env deletePostEnvelope
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

		return post, &wpError
	}

	if err != nil {
		return
	}

	if env.Previous != nil && env.Previous.ID != 0 {
		return *env.Previous, nil
	}

	return env.Post, nil
}
