package gowprest

import (
	"encoding/json"
	"strconv"
	"strings"
	"time"
)

type Comment struct {
	ID               int               `json:"id,omitempty"`
	Post             int               `json:"post,omitempty"`
	Parent           int               `json:"parent,omitempty"`
	Author           int               `json:"author,omitempty"`
	AuthorName       string            `json:"author_name,omitempty"`
	AuthorEmail      string            `json:"author_email,omitempty"`
	AuthorURL        string            `json:"author_url,omitempty"`
	AuthorIP         string            `json:"author_ip,omitempty"`
	AuthorUserAgent  string            `json:"author_user_agent,omitempty"`
	Date             *Date             `json:"date,omitempty"`
	DateGMT          *Date             `json:"date_gmt,omitempty"`
	Content          *Object           `json:"content,omitempty"`
	Link             string            `json:"link,omitempty"`
	Status           string            `json:"status,omitempty"`
	Type             string            `json:"type,omitempty"`
	AuthorAvatarURLs map[string]string `json:"author_avatar_urls,omitempty"`
	Meta             map[string]any    `json:"meta,omitempty"`
}

type CommentData struct {
	ID              int            `json:"id,omitempty"`
	Post            int            `json:"post,omitempty"`
	Parent          int            `json:"parent,omitempty"`
	Author          int            `json:"author,omitempty"`
	AuthorName      string         `json:"author_name,omitempty"`
	AuthorEmail     string         `json:"author_email,omitempty"`
	AuthorURL       string         `json:"author_url,omitempty"`
	AuthorIP        string         `json:"author_ip,omitempty"`
	AuthorUserAgent string         `json:"author_user_agent,omitempty"`
	Date            *Date          `json:"date,omitempty"`
	DateGMT         *Date          `json:"date_gmt,omitempty"`
	Content         string         `json:"content,omitempty"`
	Status          string         `json:"status,omitempty"`
	Meta            map[string]any `json:"meta,omitempty"`
}

type DeletedComment struct {
	Previous Comment `json:"previous"`
	Deleted  bool    `json:"deleted"`
}

type Comments struct {
	client *RestClient
}

func (c *RestClient) Comments() *Comments {
	return &Comments{client: c}
}

type ListComments struct {
	endpoint  string
	client    *RestClient
	arguments map[string]string
}

func (api *Comments) List() *ListComments {
	return &ListComments{
		endpoint:  "/wp/v2/comments",
		client:    api.client,
		arguments: make(map[string]string),
	}
}

func (api *ListComments) ContextView() *ListComments {
	api.arguments["context"] = "view"
	return api
}

func (api *ListComments) ContextEdit() *ListComments {
	api.arguments["context"] = "edit"
	return api
}

func (api *ListComments) ContextEmbed() *ListComments {
	api.arguments["context"] = "embed"
	return api
}

func (api *ListComments) Page(page int) *ListComments {
	api.arguments["page"] = strconv.Itoa(page)
	return api
}

func (api *ListComments) PerPage(perPage int) *ListComments {
	api.arguments["per_page"] = strconv.Itoa(perPage)
	return api
}

func (api *ListComments) Search(query string) *ListComments {
	api.arguments["search"] = query
	return api
}

func (api *ListComments) After(after time.Time) *ListComments {
	api.arguments["after"] = after.Format(time.RFC3339)
	return api
}

func (api *ListComments) Before(before time.Time) *ListComments {
	api.arguments["before"] = before.Format(time.RFC3339)
	return api
}

func (api *ListComments) Author(authorIDs ...int) *ListComments {
	authors := []string{}
	for _, id := range authorIDs {
		authors = append(authors, strconv.Itoa(id))
	}
	api.arguments["author"] = strings.Join(authors, ",")
	return api
}

func (api *ListComments) AuthorExclude(authorIDs ...int) *ListComments {
	authors := []string{}
	for _, id := range authorIDs {
		authors = append(authors, strconv.Itoa(id))
	}
	api.arguments["author_exclude"] = strings.Join(authors, ",")
	return api
}

func (api *ListComments) AuthorEmail(email string) *ListComments {
	api.arguments["author_email"] = email
	return api
}

func (api *ListComments) Exclude(excludeIDs ...int) *ListComments {
	excludes := []string{}
	for _, id := range excludeIDs {
		excludes = append(excludes, strconv.Itoa(id))
	}
	api.arguments["exclude"] = strings.Join(excludes, ",")
	return api
}

func (api *ListComments) Include(includeIDs ...int) *ListComments {
	includes := []string{}
	for _, id := range includeIDs {
		includes = append(includes, strconv.Itoa(id))
	}
	api.arguments["include"] = strings.Join(includes, ",")
	return api
}

func (api *ListComments) Offset(offset int) *ListComments {
	api.arguments["offset"] = strconv.Itoa(offset)
	return api
}

func (api *ListComments) OrderAsc() *ListComments {
	api.arguments["order"] = "asc"
	return api
}

func (api *ListComments) OrderDesc() *ListComments {
	api.arguments["order"] = "desc"
	return api
}

func (api *ListComments) OrderByDate() *ListComments {
	api.arguments["orderby"] = "date"
	return api
}

func (api *ListComments) OrderByID() *ListComments {
	api.arguments["orderby"] = "id"
	return api
}

func (api *ListComments) OrderByInclude() *ListComments {
	api.arguments["orderby"] = "include"
	return api
}

func (api *ListComments) OrderByPost() *ListComments {
	api.arguments["orderby"] = "post"
	return api
}

func (api *ListComments) OrderByParent() *ListComments {
	api.arguments["orderby"] = "parent"
	return api
}

func (api *ListComments) OrderByCommentType() *ListComments {
	api.arguments["orderby"] = "type"
	return api
}

func (api *ListComments) Parent(parentIDs ...int) *ListComments {
	parents := []string{}
	for _, id := range parentIDs {
		parents = append(parents, strconv.Itoa(id))
	}
	api.arguments["parent"] = strings.Join(parents, ",")
	return api
}

func (api *ListComments) ParentExclude(parentIDs ...int) *ListComments {
	parents := []string{}
	for _, id := range parentIDs {
		parents = append(parents, strconv.Itoa(id))
	}
	api.arguments["parent_exclude"] = strings.Join(parents, ",")
	return api
}

func (api *ListComments) Post(postIDs ...int) *ListComments {
	posts := []string{}
	for _, id := range postIDs {
		posts = append(posts, strconv.Itoa(id))
	}
	api.arguments["post"] = strings.Join(posts, ",")
	return api
}

func (api *ListComments) Status(status string) *ListComments {
	api.arguments["status"] = status
	return api
}

func (api *ListComments) Type(commentType string) *ListComments {
	api.arguments["type"] = commentType
	return api
}

func (api *ListComments) Do() (comments []Comment, err error) {
	restyClient := api.client.httpClient.R()
	if api.client.auth.Username != "" && api.client.auth.Password != "" && api.arguments["context"] == "edit" {
		restyClient.SetBasicAuth(api.client.auth.Username, api.client.auth.Password)
	}

	resp, err := restyClient.
		SetHeader("Accept", "application/json").
		SetResult(&comments).
		SetQueryParams(api.arguments).
		Get(api.client.endpoint + api.endpoint)

	if err != nil {
		return
	}

	if resp.IsError() {
		var wpError WPRestError
		err = json.Unmarshal(resp.Bytes(), &wpError)
		if err != nil {
			return
		}
		return comments, &wpError
	}

	return
}

type CreateComment struct {
	endpoint string
	client   *RestClient
	comment  CommentData
}

func (api *Comments) Create(comment CommentData) *CreateComment {
	return &CreateComment{
		endpoint: "/wp/v2/comments",
		client:   api.client,
		comment:  comment,
	}
}

func (api *CreateComment) Do() (comment Comment, err error) {
	resp, err := api.client.httpClient.R().
		SetHeader("Content-Type", "application/json").
		SetBasicAuth(api.client.auth.Username, api.client.auth.Password).
		SetResult(&comment).
		SetBody(api.comment).
		Post(api.client.endpoint + api.endpoint)

	if err != nil {
		return
	}

	if resp.IsError() {
		var wpError WPRestError
		err = json.Unmarshal(resp.Bytes(), &wpError)
		if err != nil {
			return
		}
		return comment, &wpError
	}

	return
}

type RetrieveComment struct {
	endpoint  string
	client    *RestClient
	arguments map[string]string
}

func (api *Comments) Retrieve(commentID int) *RetrieveComment {
	return &RetrieveComment{
		endpoint:  "/wp/v2/comments/" + strconv.Itoa(commentID),
		client:    api.client,
		arguments: make(map[string]string),
	}
}

func (api *RetrieveComment) ContextView() *RetrieveComment {
	api.arguments["context"] = "view"
	return api
}

func (api *RetrieveComment) ContextEdit() *RetrieveComment {
	api.arguments["context"] = "edit"
	return api
}

func (api *RetrieveComment) ContextEmbed() *RetrieveComment {
	api.arguments["context"] = "embed"
	return api
}

func (api *RetrieveComment) Password(password string) *RetrieveComment {
	api.arguments["password"] = password
	return api
}

func (api *RetrieveComment) Do() (comment *Comment, err error) {
	endpoint := api.client.endpoint + api.endpoint

	restyClient := api.client.httpClient.R()
	if api.client.auth.Username != "" && api.client.auth.Password != "" && api.arguments["context"] == "edit" {
		restyClient.SetBasicAuth(api.client.auth.Username, api.client.auth.Password)
	}

	resp, err := restyClient.
		SetHeader("Accept", "application/json").
		SetResult(&comment).
		SetQueryParams(api.arguments).
		Get(endpoint)

	if err != nil {
		return
	}

	if resp.IsError() {
		var wpError WPRestError
		err = json.Unmarshal(resp.Bytes(), &wpError)
		if err != nil {
			return
		}
		return comment, &wpError
	}

	return
}

type UpdateComment struct {
	endpoint string
	client   *RestClient
	comment  CommentData
}

func (api *Comments) Update(comment CommentData) *UpdateComment {
	return &UpdateComment{
		endpoint: "/wp/v2/comments/" + strconv.Itoa(comment.ID),
		client:   api.client,
		comment:  comment,
	}
}

func (api *UpdateComment) Do() (comment Comment, err error) {
	resp, err := api.client.httpClient.R().
		SetHeader("Content-Type", "application/json").
		SetBasicAuth(api.client.auth.Username, api.client.auth.Password).
		SetResult(&comment).
		SetBody(api.comment).
		Post(api.client.endpoint + api.endpoint)

	if err != nil {
		return
	}

	if resp.IsError() {
		var wpError WPRestError
		err = json.Unmarshal(resp.Bytes(), &wpError)
		if err != nil {
			return
		}
		return comment, &wpError
	}

	return
}

type DeleteComment struct {
	endpoint  string
	client    *RestClient
	commentID int
	force     bool
}

func (api *Comments) Delete(commentID int) *DeleteComment {
	return &DeleteComment{
		endpoint:  "/wp/v2/comments",
		client:    api.client,
		commentID: commentID,
	}
}

func (api *DeleteComment) Force() *DeleteComment {
	api.force = true
	return api
}

func (api *DeleteComment) Do() (deletedComment DeletedComment, err error) {
	endpoint := api.client.endpoint + api.endpoint + "/" + strconv.Itoa(api.commentID)
	resp, err := api.client.httpClient.R().
		SetHeader("Content-Type", "application/json").
		SetBasicAuth(api.client.auth.Username, api.client.auth.Password).
		SetResult(&deletedComment).
		SetQueryParam("force", strconv.FormatBool(api.force)).
		Delete(endpoint)

	if err != nil {
		return
	}

	if resp.IsError() {
		var wpError WPRestError
		err = json.Unmarshal(resp.Bytes(), &wpError)
		if err != nil {
			return
		}
		return deletedComment, &wpError
	}

	return
}
