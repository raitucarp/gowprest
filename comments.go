package gowprest

import (
	"encoding/json"
	"strconv"
	"strings"
	"time"
)

// CommentStatus defines the moderation state of a comment.
type CommentStatus string

const (
	CommentStatusApprove CommentStatus = "approve"
	CommentStatusHold    CommentStatus = "hold"
	CommentStatusSpam    CommentStatus = "spam"
	CommentStatusTrash   CommentStatus = "trash"
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
	Links            map[string]any    `json:"_links,omitempty"`
	Embedded         map[string]any    `json:"_embedded,omitempty"`
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
	Password        string         `json:"password,omitempty"`
}

type DeletedComment struct {
	Comment
	Previous *Comment `json:"previous,omitempty"`
	Deleted  bool     `json:"deleted"`
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

func (api *ListComments) Context(ctx string) *ListComments {
	api.arguments["context"] = ctx
	return api
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
	authors := make([]string, 0, len(authorIDs))
	for _, id := range authorIDs {
		authors = append(authors, strconv.Itoa(id))
	}
	api.arguments["author"] = strings.Join(authors, ",")
	return api
}

func (api *ListComments) AuthorExclude(authorIDs ...int) *ListComments {
	authors := make([]string, 0, len(authorIDs))
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
	excludes := make([]string, 0, len(excludeIDs))
	for _, id := range excludeIDs {
		excludes = append(excludes, strconv.Itoa(id))
	}
	api.arguments["exclude"] = strings.Join(excludes, ",")
	return api
}

func (api *ListComments) Include(includeIDs ...int) *ListComments {
	includes := make([]string, 0, len(includeIDs))
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

func (api *ListComments) Order(order string) *ListComments {
	api.arguments["order"] = order
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

func (api *ListComments) OrderBy(orderBy string) *ListComments {
	api.arguments["orderby"] = orderBy
	return api
}

func (api *ListComments) OrderByDate() *ListComments {
	api.arguments["orderby"] = "date"
	return api
}

func (api *ListComments) OrderByDateGMT() *ListComments {
	api.arguments["orderby"] = "date_gmt"
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

func (api *ListComments) OrderByType() *ListComments {
	api.arguments["orderby"] = "type"
	return api
}

func (api *ListComments) Parent(parentIDs ...int) *ListComments {
	parents := make([]string, 0, len(parentIDs))
	for _, id := range parentIDs {
		parents = append(parents, strconv.Itoa(id))
	}
	api.arguments["parent"] = strings.Join(parents, ",")
	return api
}

func (api *ListComments) ParentExclude(parentIDs ...int) *ListComments {
	parents := make([]string, 0, len(parentIDs))
	for _, id := range parentIDs {
		parents = append(parents, strconv.Itoa(id))
	}
	api.arguments["parent_exclude"] = strings.Join(parents, ",")
	return api
}

func (api *ListComments) Post(postIDs ...int) *ListComments {
	posts := make([]string, 0, len(postIDs))
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

func (api *ListComments) StatusApprove() *ListComments {
	api.arguments["status"] = string(CommentStatusApprove)
	return api
}

func (api *ListComments) StatusHold() *ListComments {
	api.arguments["status"] = string(CommentStatusHold)
	return api
}

func (api *ListComments) StatusSpam() *ListComments {
	api.arguments["status"] = string(CommentStatusSpam)
	return api
}

func (api *ListComments) StatusTrash() *ListComments {
	api.arguments["status"] = string(CommentStatusTrash)
	return api
}

func (api *ListComments) Type(commentType string) *ListComments {
	api.arguments["type"] = commentType
	return api
}

func (api *ListComments) Password(password string) *ListComments {
	api.arguments["password"] = password
	return api
}

func (api *ListComments) Embed() *ListComments {
	api.arguments["_embed"] = "true"
	return api
}

func (api *ListComments) Fields(fields ...string) *ListComments {
	api.arguments["_fields"] = strings.Join(fields, ",")
	return api
}

func (api *ListComments) Do() (comments []Comment, err error) {
	restyClient := api.client.httpClient.R()
	if api.client.auth.Username != "" && api.client.auth.Password != "" {
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

func (api *Comments) Create(comment ...CommentData) *CreateComment {
	req := &CreateComment{
		endpoint: "/wp/v2/comments",
		client:   api.client,
	}
	if len(comment) > 0 {
		req.comment = comment[0]
	}
	return req
}

func (api *CreateComment) Post(postID int) *CreateComment {
	api.comment.Post = postID
	return api
}

func (api *CreateComment) Parent(parentID int) *CreateComment {
	api.comment.Parent = parentID
	return api
}

func (api *CreateComment) Author(authorID int) *CreateComment {
	api.comment.Author = authorID
	return api
}

func (api *CreateComment) AuthorName(name string) *CreateComment {
	api.comment.AuthorName = name
	return api
}

func (api *CreateComment) AuthorEmail(email string) *CreateComment {
	api.comment.AuthorEmail = email
	return api
}

func (api *CreateComment) AuthorURL(url string) *CreateComment {
	api.comment.AuthorURL = url
	return api
}

func (api *CreateComment) AuthorIP(ip string) *CreateComment {
	api.comment.AuthorIP = ip
	return api
}

func (api *CreateComment) AuthorUserAgent(ua string) *CreateComment {
	api.comment.AuthorUserAgent = ua
	return api
}

func (api *CreateComment) Date(t time.Time) *CreateComment {
	api.comment.Date = &Date{Time: t}
	return api
}

func (api *CreateComment) DateGMT(t time.Time) *CreateComment {
	api.comment.DateGMT = &Date{Time: t}
	return api
}

func (api *CreateComment) Content(content string) *CreateComment {
	api.comment.Content = content
	return api
}

func (api *CreateComment) Status(status string) *CreateComment {
	api.comment.Status = status
	return api
}

func (api *CreateComment) StatusApprove() *CreateComment {
	api.comment.Status = string(CommentStatusApprove)
	return api
}

func (api *CreateComment) StatusHold() *CreateComment {
	api.comment.Status = string(CommentStatusHold)
	return api
}

func (api *CreateComment) StatusSpam() *CreateComment {
	api.comment.Status = string(CommentStatusSpam)
	return api
}

func (api *CreateComment) StatusTrash() *CreateComment {
	api.comment.Status = string(CommentStatusTrash)
	return api
}

func (api *CreateComment) Meta(meta map[string]any) *CreateComment {
	api.comment.Meta = meta
	return api
}

func (api *CreateComment) Password(password string) *CreateComment {
	api.comment.Password = password
	return api
}

func (api *CreateComment) Do() (comment Comment, err error) {
	restyClient := api.client.httpClient.R()
	if api.client.auth.Username != "" && api.client.auth.Password != "" {
		restyClient.SetBasicAuth(api.client.auth.Username, api.client.auth.Password)
	}

	resp, err := restyClient.
		SetHeader("Content-Type", "application/json").
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

func (api *RetrieveComment) Context(ctx string) *RetrieveComment {
	api.arguments["context"] = ctx
	return api
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

func (api *RetrieveComment) Embed() *RetrieveComment {
	api.arguments["_embed"] = "true"
	return api
}

func (api *RetrieveComment) Fields(fields ...string) *RetrieveComment {
	api.arguments["_fields"] = strings.Join(fields, ",")
	return api
}

func (api *RetrieveComment) Do() (comment *Comment, err error) {
	endpoint := api.client.endpoint + api.endpoint

	restyClient := api.client.httpClient.R()
	if api.client.auth.Username != "" && api.client.auth.Password != "" {
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

func (api *Comments) Update(comment ...CommentData) *UpdateComment {
	endpoint := "/wp/v2/comments"
	var c CommentData
	if len(comment) > 0 {
		c = comment[0]
		if c.ID != 0 {
			endpoint = "/wp/v2/comments/" + strconv.Itoa(c.ID)
		}
	}
	return &UpdateComment{
		endpoint: endpoint,
		client:   api.client,
		comment:  c,
	}
}

func (api *UpdateComment) ID(commentID int) *UpdateComment {
	api.comment.ID = commentID
	api.endpoint = "/wp/v2/comments/" + strconv.Itoa(commentID)
	return api
}

func (api *UpdateComment) Post(postID int) *UpdateComment {
	api.comment.Post = postID
	return api
}

func (api *UpdateComment) Parent(parentID int) *UpdateComment {
	api.comment.Parent = parentID
	return api
}

func (api *UpdateComment) Author(authorID int) *UpdateComment {
	api.comment.Author = authorID
	return api
}

func (api *UpdateComment) AuthorName(name string) *UpdateComment {
	api.comment.AuthorName = name
	return api
}

func (api *UpdateComment) AuthorEmail(email string) *UpdateComment {
	api.comment.AuthorEmail = email
	return api
}

func (api *UpdateComment) AuthorURL(url string) *UpdateComment {
	api.comment.AuthorURL = url
	return api
}

func (api *UpdateComment) AuthorIP(ip string) *UpdateComment {
	api.comment.AuthorIP = ip
	return api
}

func (api *UpdateComment) AuthorUserAgent(ua string) *UpdateComment {
	api.comment.AuthorUserAgent = ua
	return api
}

func (api *UpdateComment) Date(t time.Time) *UpdateComment {
	api.comment.Date = &Date{Time: t}
	return api
}

func (api *UpdateComment) DateGMT(t time.Time) *UpdateComment {
	api.comment.DateGMT = &Date{Time: t}
	return api
}

func (api *UpdateComment) Content(content string) *UpdateComment {
	api.comment.Content = content
	return api
}

func (api *UpdateComment) Status(status string) *UpdateComment {
	api.comment.Status = status
	return api
}

func (api *UpdateComment) StatusApprove() *UpdateComment {
	api.comment.Status = string(CommentStatusApprove)
	return api
}

func (api *UpdateComment) StatusHold() *UpdateComment {
	api.comment.Status = string(CommentStatusHold)
	return api
}

func (api *UpdateComment) StatusSpam() *UpdateComment {
	api.comment.Status = string(CommentStatusSpam)
	return api
}

func (api *UpdateComment) StatusTrash() *UpdateComment {
	api.comment.Status = string(CommentStatusTrash)
	return api
}

func (api *UpdateComment) Meta(meta map[string]any) *UpdateComment {
	api.comment.Meta = meta
	return api
}

func (api *UpdateComment) Password(password string) *UpdateComment {
	api.comment.Password = password
	return api
}

func (api *UpdateComment) Do() (comment Comment, err error) {
	restyClient := api.client.httpClient.R()
	if api.client.auth.Username != "" && api.client.auth.Password != "" {
		restyClient.SetBasicAuth(api.client.auth.Username, api.client.auth.Password)
	}

	resp, err := restyClient.
		SetHeader("Content-Type", "application/json").
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
	password  string
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

func (api *DeleteComment) Password(password string) *DeleteComment {
	api.password = password
	return api
}

func (api *DeleteComment) Do() (deletedComment DeletedComment, err error) {
	endpoint := api.client.endpoint + api.endpoint + "/" + strconv.Itoa(api.commentID)
	restyClient := api.client.httpClient.R()
	if api.client.auth.Username != "" && api.client.auth.Password != "" {
		restyClient.SetBasicAuth(api.client.auth.Username, api.client.auth.Password)
	}

	if api.password != "" {
		restyClient.SetQueryParam("password", api.password)
	}

	resp, err := restyClient.
		SetHeader("Content-Type", "application/json").
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

	if deletedComment.Previous != nil && deletedComment.Previous.ID != 0 {
		deletedComment.Comment = *deletedComment.Previous
	}

	return
}
