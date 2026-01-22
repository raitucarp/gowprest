package gowprest

import (
	"encoding/json"
	"strconv"
	"strings"
	"time"
)

type Object struct {
	Rendered  string `json:"rendered"`
	Protected bool   `json:"protected,omitempty,omitzero"`
}

type Date struct {
	time.Time
}

func (ct *Date) UnmarshalJSON(b []byte) error {
	s := strings.Trim(string(b), "\"")
	t, err := time.Parse("2006-01-02T15:04:05", s)
	if err != nil {
		return err
	}
	ct.Time = t
	return nil
}

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
}

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

type Posts struct {
	client *RestClient
}

func (c *RestClient) Posts() *Posts {
	return &Posts{client: c}
}

func (api *Posts) Revisions(parentID int) *PostRevisions {
	return &PostRevisions{
		client:   api.client,
		parentID: parentID,
	}
}

type ListPosts struct {
	endpoint  string
	client    *RestClient
	arguments map[string]string
}

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

func (api *ListPosts) Author(authorID int) *ListPosts {
	api.arguments["author"] = strconv.Itoa(authorID)
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

func (api *ListPosts) SearchColumns(columns ...string) *ListPosts {
	api.arguments["search_columns"] = strings.Join(columns, ",")
	return api
}

func (api *ListPosts) Slug(slug string) *ListPosts {
	api.arguments["slug"] = slug
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

func (api *ListPosts) TaxAnd() *ListPosts {
	api.arguments["tax_relation"] = "AND"
	return api
}

func (api *ListPosts) TaxOr() *ListPosts {
	api.arguments["tax_relation"] = "OR"
	return api
}

func (api *ListPosts) Categories(categories ...string) *ListPosts {
	api.arguments["category"] = strings.Join(categories, ",")
	return api
}

func (api *ListPosts) CategoriesExclude(categories ...string) *ListPosts {
	api.arguments["category_exclude"] = strings.Join(categories, ",")
	return api
}

func (api *ListPosts) Tags(tags ...string) *ListPosts {
	api.arguments["tags"] = strings.Join(tags, ",")
	return api
}

func (api *ListPosts) TagsExclude(tags ...string) *ListPosts {
	api.arguments["tags_exclude"] = strings.Join(tags, ",")
	return api
}

func (api *ListPosts) Sticky(sticky bool) *ListPosts {
	api.arguments["sticky"] = strconv.FormatBool(sticky)
	return api
}

func (api *ListPosts) Do() (posts []Post, err error) {
	_, err = api.client.httpClient.R().
		SetHeader("Accept", "application/json").
		SetResult(&posts).
		SetQueryParams(api.arguments).
		Get(api.client.endpoint + api.endpoint)

	if err != nil {
		return
	}

	return
}

type CreatePost struct {
	endpoint string
	client   *RestClient
	post     PostData
}

func (api *Posts) Create(post PostData) *CreatePost {
	return &CreatePost{
		endpoint: "/wp/v2/posts",
		client:   api.client,
		post:     post,
	}
}

func (api *CreatePost) Do() (post Post, err error) {
	resp, err := api.client.httpClient.R().
		SetHeader("Content-Type", "application/json").
		SetBasicAuth(api.client.auth.Username, api.client.auth.Password).
		SetResult(&post).
		SetBody(api.post).
		// SetError(err).
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

type RetrievePost struct {
	endpoint  string
	client    *RestClient
	arguments map[string]string
}

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

func (api *RetrievePost) Do() (post *Post, err error) {
	endpoint := api.client.endpoint + api.endpoint

	restyClient := api.client.httpClient.R()
	if api.client.auth.Username != "" && api.client.auth.Password != "" && api.arguments["context"] == "edit" {
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

	// TODO: need fixing of message = invalid suit value: trash
	if err != nil && err.Error() == "invalid suit value: trash" {
		err = nil
	}

	if err != nil {
		return
	}

	return
}

type UpdatePost struct {
	endpoint string
	client   *RestClient
	post     PostData
}

func (api *Posts) Update(post PostData) *UpdatePost {
	return &UpdatePost{
		endpoint: "/wp/v2/posts/" + strconv.Itoa(post.ID),
		client:   api.client,
		post:     post,
	}
}

func (api *UpdatePost) Do() (post Post, err error) {
	resp, err := api.client.httpClient.R().
		SetHeader("Content-Type", "application/json").
		SetBasicAuth(api.client.auth.Username, api.client.auth.Password).
		SetResult(&post).
		SetBody(api.post).
		// SetError(err).
		Post(api.client.endpoint + api.endpoint)

	if resp.IsError() {
		var wpError WPRestError
		err = json.Unmarshal(resp.Bytes(), &wpError)

		if err != nil {
			return
		}

		return post, &wpError
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

type DeletePost struct {
	endpoint string
	client   *RestClient
	postId   int
	force    bool
}

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

func (api *DeletePost) Do() (post Post, err error) {
	endpoint := api.client.endpoint + api.endpoint + "/" + strconv.Itoa(api.postId)
	resp, err :=
		api.client.httpClient.R().
			SetHeader("Content-Type", "application/json").
			SetBasicAuth(api.client.auth.Username, api.client.auth.Password).
			SetResult(&post).
			SetPathParam("force", strconv.FormatBool(api.force)).
			Delete(endpoint)

	if resp.IsError() {
		var wpError WPRestError
		err = json.Unmarshal(resp.Bytes(), &wpError)

		if err != nil {
			return
		}

		return post, &wpError
	}

	// TODO: need fixing of message = invalid suit value: trash
	if err != nil && err.Error() == "invalid suit value: trash" {
		err = nil
	}

	return
}
