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

type NewPost struct {
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

type PostsAPI struct {
	client *RestClient
}

func (c *RestClient) Posts() *PostsAPI {
	return &PostsAPI{client: c}
}

type ListPostsAPI struct {
	endpoint  string
	client    *RestClient
	arguments map[string]string
}

func (api *PostsAPI) List() *ListPostsAPI {
	return &ListPostsAPI{
		endpoint:  "/wp/v2/posts",
		client:    api.client,
		arguments: make(map[string]string),
	}
}

func (api *ListPostsAPI) ContextView() *ListPostsAPI {
	api.arguments["context"] = "view"
	return api
}

func (api *ListPostsAPI) ContextEdit() *ListPostsAPI {
	api.arguments["context"] = "edit"
	return api
}

func (api *ListPostsAPI) ContextEmbed() *ListPostsAPI {
	api.arguments["context"] = "embed"
	return api
}

func (api *ListPostsAPI) Page(page int) *ListPostsAPI {
	api.arguments["page"] = strconv.Itoa(page)
	return api
}

func (api *ListPostsAPI) PerPage(perPage int) *ListPostsAPI {
	api.arguments["per_page"] = strconv.Itoa(perPage)
	return api
}

func (api *ListPostsAPI) Search(query string) *ListPostsAPI {
	api.arguments["search"] = query
	return api
}

func (api *ListPostsAPI) After(after time.Time) *ListPostsAPI {
	api.arguments["after"] = after.Format(time.RFC3339)
	return api
}

func (api *ListPostsAPI) ModifiedAfter(modifiedAfter time.Time) *ListPostsAPI {
	api.arguments["modified_after"] = modifiedAfter.Format(time.RFC3339)
	return api
}

func (api *ListPostsAPI) Author(authorID int) *ListPostsAPI {
	api.arguments["author"] = strconv.Itoa(authorID)
	return api
}

func (api *ListPostsAPI) AuthorExclude(authorIDs ...int) *ListPostsAPI {
	authors := []string{}

	for _, authorId := range authorIDs {
		authors = append(authors, strconv.Itoa(authorId))
	}

	api.arguments["author_exclude"] = strings.Join(authors, ",")
	return api
}

func (api *ListPostsAPI) Before(before time.Time) *ListPostsAPI {
	api.arguments["before"] = before.Format(time.RFC3339)
	return api
}

func (api *ListPostsAPI) ModifiedBefore(modifiedBefore time.Time) *ListPostsAPI {
	api.arguments["modified_before"] = modifiedBefore.Format(time.RFC3339)
	return api
}

func (api *ListPostsAPI) Exclude(excludeIDs ...int) *ListPostsAPI {
	excludes := []string{}
	for _, excludeId := range excludeIDs {
		excludes = append(excludes, strconv.Itoa(excludeId))
	}
	api.arguments["exclude"] = strings.Join(excludes, ",")
	return api
}

func (api *ListPostsAPI) Include(includeIDs ...int) *ListPostsAPI {
	includes := []string{}
	for _, includeId := range includeIDs {
		includes = append(includes, strconv.Itoa(includeId))
	}
	api.arguments["include"] = strings.Join(includes, ",")
	return api
}

func (api *ListPostsAPI) Offset(offset int) *ListPostsAPI {
	api.arguments["offset"] = strconv.Itoa(offset)
	return api
}

func (api *ListPostsAPI) OrderAsc() *ListPostsAPI {
	api.arguments["order"] = "asc"
	return api
}

func (api *ListPostsAPI) OrderDesc() *ListPostsAPI {
	api.arguments["order"] = "desc"
	return api
}

func (api *ListPostsAPI) OrderByAuthor() *ListPostsAPI {
	api.arguments["orderby"] = "author"
	return api
}

func (api *ListPostsAPI) OrderByDate() *ListPostsAPI {
	api.arguments["orderby"] = "date"
	return api
}

func (api *ListPostsAPI) OrderById() *ListPostsAPI {
	api.arguments["orderby"] = "id"
	return api
}

func (api *ListPostsAPI) OrderByInclude() *ListPostsAPI {
	api.arguments["orderby"] = "include"
	return api
}

func (api *ListPostsAPI) OrderByModified() *ListPostsAPI {
	api.arguments["orderby"] = "modified"
	return api
}

func (api *ListPostsAPI) OrderByParent() *ListPostsAPI {
	api.arguments["orderby"] = "parent"
	return api
}

func (api *ListPostsAPI) OrderByRelevance() *ListPostsAPI {
	api.arguments["orderby"] = "relevance"
	return api
}

func (api *ListPostsAPI) OrderBySlug() *ListPostsAPI {
	api.arguments["orderby"] = "slug"
	return api
}

func (api *ListPostsAPI) OrderByIncludeSlug() *ListPostsAPI {
	api.arguments["orderby"] = "include_slugs"
	return api
}

func (api *ListPostsAPI) OrderByTitle() *ListPostsAPI {
	api.arguments["orderby"] = "title"
	return api
}

func (api *ListPostsAPI) SearchColumns(columns ...string) *ListPostsAPI {
	api.arguments["search_columns"] = strings.Join(columns, ",")
	return api
}

func (api *ListPostsAPI) Slug(slug string) *ListPostsAPI {
	api.arguments["slug"] = slug
	return api
}

func (api *ListPostsAPI) StatusPublish() *ListPostsAPI {
	api.arguments["status"] = "publish"
	return api
}

func (api *ListPostsAPI) StatusDraft() *ListPostsAPI {
	api.arguments["status"] = "draft"
	return api
}

func (api *ListPostsAPI) StatusPending() *ListPostsAPI {
	api.arguments["status"] = "pending"
	return api
}

func (api *ListPostsAPI) StatusPrivate() *ListPostsAPI {
	api.arguments["status"] = "private"
	return api
}

func (api *ListPostsAPI) StatusFuture() *ListPostsAPI {
	api.arguments["status"] = "future"
	return api
}

func (api *ListPostsAPI) StatusTrash() *ListPostsAPI {
	api.arguments["status"] = "trash"
	return api
}

func (api *ListPostsAPI) StatusAny() *ListPostsAPI {
	api.arguments["status"] = "any"
	return api
}

func (api *ListPostsAPI) TaxAnd() *ListPostsAPI {
	api.arguments["tax_relation"] = "AND"
	return api
}

func (api *ListPostsAPI) TaxOr() *ListPostsAPI {
	api.arguments["tax_relation"] = "OR"
	return api
}

func (api *ListPostsAPI) Categories(categories ...string) *ListPostsAPI {
	api.arguments["category"] = strings.Join(categories, ",")
	return api
}

func (api *ListPostsAPI) CategoriesExclude(categories ...string) *ListPostsAPI {
	api.arguments["category_exclude"] = strings.Join(categories, ",")
	return api
}

func (api *ListPostsAPI) Tags(tags ...string) *ListPostsAPI {
	api.arguments["tags"] = strings.Join(tags, ",")
	return api
}

func (api *ListPostsAPI) TagsExclude(tags ...string) *ListPostsAPI {
	api.arguments["tags_exclude"] = strings.Join(tags, ",")
	return api
}

func (api *ListPostsAPI) Sticky(sticky bool) *ListPostsAPI {
	api.arguments["sticky"] = strconv.FormatBool(sticky)
	return api
}

func (api *ListPostsAPI) Do() (posts []Post, err error) {
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

type CreatePostAPI struct {
	endpoint string
	client   *RestClient
	post     NewPost
}

func (api *PostsAPI) Create(post NewPost) *CreatePostAPI {
	return &CreatePostAPI{
		endpoint: "/wp/v2/posts",
		client:   api.client,
		post:     post,
	}
}

func (api *CreatePostAPI) Do() (post Post, err error) {
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

type PostAPI struct {
	endpoint  string
	client    *RestClient
	postId    int
	arguments map[string]string
}

func (api *PostsAPI) Retrieve(postId int) *PostAPI {
	return &PostAPI{
		endpoint:  "/wp/v2/posts",
		client:    api.client,
		postId:    postId,
		arguments: make(map[string]string),
	}
}

func (api *PostAPI) ContextView() *PostAPI {
	api.arguments["context"] = "view"
	return api
}

func (api *PostAPI) ContextEdit() *PostAPI {
	api.arguments["context"] = "edit"
	return api
}

func (api *PostAPI) ContextEmbed() *PostAPI {
	api.arguments["context"] = "embed"
	return api
}

func (api *PostAPI) Password(password string) *PostAPI {
	api.arguments["password"] = password
	return api
}

func (api *PostAPI) Do() (post *Post, err error) {
	endpoint := api.client.endpoint + api.endpoint + "/" + strconv.Itoa(api.postId)

	_, err = api.client.httpClient.R().
		SetHeader("Accept", "application/json").
		SetResult(&post).
		SetQueryParams(api.arguments).
		Get(endpoint)

	if err != nil {
		return
	}

	return
}
