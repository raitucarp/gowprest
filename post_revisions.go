package gowprest

import (
	"encoding/json"
	"strconv"
	"strings"
)

// Revision represents a WordPress post revision.
type Revision struct {
	Author      int       `json:"author,omitempty"`
	Date        *Date     `json:"date,omitempty"`
	DateGMT     *Date     `json:"date_gmt,omitempty"`
	GUID        *Object   `json:"guid,omitempty"`
	ID          int       `json:"id,omitempty"`
	Modified    *Date     `json:"modified,omitempty"`
	ModifiedGMT *Date     `json:"modified_gmt,omitempty"`
	Parent      int       `json:"parent,omitempty"`
	Slug        string    `json:"slug,omitempty"`
	Title       *Object   `json:"title,omitempty"`
	Content     *Object   `json:"content,omitempty"`
	Excerpt     *Object   `json:"excerpt,omitempty"`
	Links       map[string]any `json:"_links,omitempty"`
	Embedded    map[string]any `json:"_embedded,omitempty"`
}

// PostRevisions anchors revision-related operations for a specific post.
type PostRevisions struct {
	client   *RestClient
	parentID int
}

// List returns a ListPostRevisions struct to list revisions for the parent post.
func (api *PostRevisions) List() *ListPostRevisions {
	return &ListPostRevisions{
		endpoint:  "/wp/v2/posts/" + strconv.Itoa(api.parentID) + "/revisions",
		client:    api.client,
		arguments: make(map[string]string),
	}
}

// Retrieve returns a RetrievePostRevision struct to get a specific revision.
func (api *PostRevisions) Retrieve(revisionID int) *RetrievePostRevision {
	return &RetrievePostRevision{
		endpoint:   "/wp/v2/posts/" + strconv.Itoa(api.parentID) + "/revisions",
		client:     api.client,
		revisionID: revisionID,
		arguments:  make(map[string]string),
	}
}

// Delete returns a DeletePostRevision struct to delete a specific revision.
func (api *PostRevisions) Delete(revisionID int) *DeletePostRevision {
	return &DeletePostRevision{
		endpoint:   "/wp/v2/posts/" + strconv.Itoa(api.parentID) + "/revisions",
		client:     api.client,
		revisionID: revisionID,
		force:      true, // WP post revisions always require force=true
	}
}

// Create returns a CreatePostRevision struct to create a new revision (autosave).
// It can be used directly or with variadic optional PostData for backward compatibility.
func (api *PostRevisions) Create(revision ...PostData) *CreatePostRevision {
	req := &CreatePostRevision{
		endpoint: "/wp/v2/posts/" + strconv.Itoa(api.parentID) + "/autosaves",
		client:   api.client,
	}
	if len(revision) > 0 {
		req.revision = revision[0]
	}
	return req
}

// Autosaves returns an Autosaves endpoint builder for this post.
func (api *PostRevisions) Autosaves() *PostAutosaves {
	return &PostAutosaves{
		client:   api.client,
		parentID: api.parentID,
	}
}

// ListPostRevisions handles listing revisions.
type ListPostRevisions struct {
	endpoint  string
	client    *RestClient
	arguments map[string]string
}

func (api *ListPostRevisions) Context(ctx string) *ListPostRevisions {
	api.arguments["context"] = ctx
	return api
}

func (api *ListPostRevisions) ContextView() *ListPostRevisions {
	api.arguments["context"] = "view"
	return api
}

func (api *ListPostRevisions) ContextEdit() *ListPostRevisions {
	api.arguments["context"] = "edit"
	return api
}

func (api *ListPostRevisions) ContextEmbed() *ListPostRevisions {
	api.arguments["context"] = "embed"
	return api
}

func (api *ListPostRevisions) Page(page int) *ListPostRevisions {
	api.arguments["page"] = strconv.Itoa(page)
	return api
}

func (api *ListPostRevisions) PerPage(perPage int) *ListPostRevisions {
	api.arguments["per_page"] = strconv.Itoa(perPage)
	return api
}

func (api *ListPostRevisions) Search(query string) *ListPostRevisions {
	api.arguments["search"] = query
	return api
}

func (api *ListPostRevisions) Exclude(excludeIDs ...int) *ListPostRevisions {
	excludes := make([]string, 0, len(excludeIDs))
	for _, id := range excludeIDs {
		excludes = append(excludes, strconv.Itoa(id))
	}
	api.arguments["exclude"] = strings.Join(excludes, ",")
	return api
}

func (api *ListPostRevisions) Include(includeIDs ...int) *ListPostRevisions {
	includes := make([]string, 0, len(includeIDs))
	for _, id := range includeIDs {
		includes = append(includes, strconv.Itoa(id))
	}
	api.arguments["include"] = strings.Join(includes, ",")
	return api
}

func (api *ListPostRevisions) Offset(offset int) *ListPostRevisions {
	api.arguments["offset"] = strconv.Itoa(offset)
	return api
}

func (api *ListPostRevisions) Order(order string) *ListPostRevisions {
	api.arguments["order"] = order
	return api
}

func (api *ListPostRevisions) OrderAsc() *ListPostRevisions {
	api.arguments["order"] = "asc"
	return api
}

func (api *ListPostRevisions) OrderDesc() *ListPostRevisions {
	api.arguments["order"] = "desc"
	return api
}

func (api *ListPostRevisions) OrderBy(orderBy string) *ListPostRevisions {
	api.arguments["orderby"] = orderBy
	return api
}

func (api *ListPostRevisions) OrderByDate() *ListPostRevisions {
	api.arguments["orderby"] = "date"
	return api
}

func (api *ListPostRevisions) OrderById() *ListPostRevisions {
	api.arguments["orderby"] = "id"
	return api
}

func (api *ListPostRevisions) OrderByInclude() *ListPostRevisions {
	api.arguments["orderby"] = "include"
	return api
}

func (api *ListPostRevisions) OrderByRelevance() *ListPostRevisions {
	api.arguments["orderby"] = "relevance"
	return api
}

func (api *ListPostRevisions) OrderBySlug() *ListPostRevisions {
	api.arguments["orderby"] = "slug"
	return api
}

func (api *ListPostRevisions) OrderByTitle() *ListPostRevisions {
	api.arguments["orderby"] = "title"
	return api
}

func (api *ListPostRevisions) Embed() *ListPostRevisions {
	api.arguments["_embed"] = "true"
	return api
}

func (api *ListPostRevisions) Fields(fields ...string) *ListPostRevisions {
	api.arguments["_fields"] = strings.Join(fields, ",")
	return api
}

func (api *ListPostRevisions) Do() (revisions []Revision, err error) {
	restyClient := api.client.httpClient.R()
	if api.client.auth.Username != "" && api.client.auth.Password != "" {
		restyClient.SetBasicAuth(api.client.auth.Username, api.client.auth.Password)
	}

	resp, err := restyClient.
		SetHeader("Accept", "application/json").
		SetResult(&revisions).
		SetQueryParams(api.arguments).
		Get(api.client.endpoint + api.endpoint)

	if resp.IsError() {
		var wpError WPRestError
		err = json.Unmarshal(resp.Bytes(), &wpError)
		if err != nil {
			return
		}
		return revisions, &wpError
	}

	return
}

// RetrievePostRevision handles retrieving a specific revision.
type RetrievePostRevision struct {
	endpoint   string
	client     *RestClient
	revisionID int
	arguments  map[string]string
}

func (api *RetrievePostRevision) Context(ctx string) *RetrievePostRevision {
	api.arguments["context"] = ctx
	return api
}

func (api *RetrievePostRevision) ContextView() *RetrievePostRevision {
	api.arguments["context"] = "view"
	return api
}

func (api *RetrievePostRevision) ContextEdit() *RetrievePostRevision {
	api.arguments["context"] = "edit"
	return api
}

func (api *RetrievePostRevision) ContextEmbed() *RetrievePostRevision {
	api.arguments["context"] = "embed"
	return api
}

func (api *RetrievePostRevision) Embed() *RetrievePostRevision {
	api.arguments["_embed"] = "true"
	return api
}

func (api *RetrievePostRevision) Fields(fields ...string) *RetrievePostRevision {
	api.arguments["_fields"] = strings.Join(fields, ",")
	return api
}

func (api *RetrievePostRevision) Do() (revision *Revision, err error) {
	endpoint := api.client.endpoint + api.endpoint + "/" + strconv.Itoa(api.revisionID)

	restyClient := api.client.httpClient.R()
	if api.client.auth.Username != "" && api.client.auth.Password != "" {
		restyClient.SetBasicAuth(api.client.auth.Username, api.client.auth.Password)
	}

	resp, err := restyClient.
		SetHeader("Accept", "application/json").
		SetResult(&revision).
		SetQueryParams(api.arguments).
		Get(endpoint)

	if resp.IsError() {
		var wpError WPRestError
		err = json.Unmarshal(resp.Bytes(), &wpError)
		if err != nil {
			return
		}
		return revision, &wpError
	}

	return
}

type deleteRevisionEnvelope struct {
	Revision
	Deleted  bool      `json:"deleted"`
	Previous *Revision `json:"previous"`
}

// DeletePostRevision handles deleting a specific revision.
type DeletePostRevision struct {
	endpoint   string
	client     *RestClient
	revisionID int
	force      bool
}

func (api *DeletePostRevision) Force() *DeletePostRevision {
	api.force = true
	return api
}

func (api *DeletePostRevision) Do() (revision Revision, err error) {
	endpoint := api.client.endpoint + api.endpoint + "/" + strconv.Itoa(api.revisionID)
	restyClient := api.client.httpClient.R()

	if api.client.auth.Username != "" && api.client.auth.Password != "" {
		restyClient = restyClient.SetBasicAuth(api.client.auth.Username, api.client.auth.Password)
	}

	var env deleteRevisionEnvelope
	resp, err := restyClient.
		SetHeader("Content-Type", "application/json").
		SetResult(&env).
		SetQueryParam("force", strconv.FormatBool(api.force)).
		Delete(endpoint)

	if resp.IsError() {
		var wpError WPRestError
		err = json.Unmarshal(resp.Bytes(), &wpError)
		if err != nil {
			return
		}
		return revision, &wpError
	}

	if err != nil {
		return
	}

	if env.Previous != nil && env.Previous.ID != 0 {
		return *env.Previous, nil
	}

	return env.Revision, nil
}

// CreatePostRevision handles creating a new revision (autosave).
type CreatePostRevision struct {
	endpoint string
	client   *RestClient
	revision PostData
}

func (api *CreatePostRevision) Title(title string) *CreatePostRevision {
	api.revision.Title = title
	return api
}

func (api *CreatePostRevision) Content(content string) *CreatePostRevision {
	api.revision.Content = content
	return api
}

func (api *CreatePostRevision) Excerpt(excerpt string) *CreatePostRevision {
	api.revision.Excerpt = excerpt
	return api
}

func (api *CreatePostRevision) Slug(slug string) *CreatePostRevision {
	api.revision.Slug = slug
	return api
}

func (api *CreatePostRevision) Status(status PostStatus) *CreatePostRevision {
	api.revision.Status = status
	return api
}

func (api *CreatePostRevision) Author(author int) *CreatePostRevision {
	api.revision.Author = author
	return api
}

func (api *CreatePostRevision) Password(password string) *CreatePostRevision {
	api.revision.Password = password
	return api
}

func (api *CreatePostRevision) FeaturedMedia(mediaID int) *CreatePostRevision {
	api.revision.FeaturedMedia = mediaID
	return api
}

func (api *CreatePostRevision) CommentStatus(status OpenClosedStatus) *CreatePostRevision {
	api.revision.CommentStatus = status
	return api
}

func (api *CreatePostRevision) PingStatus(status OpenClosedStatus) *CreatePostRevision {
	api.revision.PingStatus = status
	return api
}

func (api *CreatePostRevision) Format(format Format) *CreatePostRevision {
	api.revision.Format = format
	return api
}

func (api *CreatePostRevision) Meta(meta map[string]any) *CreatePostRevision {
	api.revision.Meta = meta
	return api
}

func (api *CreatePostRevision) Sticky(sticky bool) *CreatePostRevision {
	api.revision.Sticky = sticky
	return api
}

func (api *CreatePostRevision) Template(template string) *CreatePostRevision {
	api.revision.Template = template
	return api
}

func (api *CreatePostRevision) Categories(categoryIDs ...int) *CreatePostRevision {
	api.revision.Categories = categoryIDs
	return api
}

func (api *CreatePostRevision) Tags(tagIDs ...int) *CreatePostRevision {
	api.revision.Tags = tagIDs
	return api
}

func (api *CreatePostRevision) Do() (revision Revision, err error) {
	resp, err := api.client.httpClient.R().
		SetHeader("Content-Type", "application/json").
		SetBasicAuth(api.client.auth.Username, api.client.auth.Password).
		SetResult(&revision).
		SetBody(api.revision).
		Post(api.client.endpoint + api.endpoint)

	if resp.IsError() {
		var wpError WPRestError
		err = json.Unmarshal(resp.Bytes(), &wpError)
		if err != nil {
			return
		}
		return revision, &wpError
	}

	return
}

// PostAutosaves provides operations for post autosaves.
type PostAutosaves struct {
	client   *RestClient
	parentID int
}

// List returns a ListPostAutosaves builder to get all autosaves for this post.
func (api *PostAutosaves) List() *ListPostAutosaves {
	return &ListPostAutosaves{
		endpoint:  "/wp/v2/posts/" + strconv.Itoa(api.parentID) + "/autosaves",
		client:    api.client,
		arguments: make(map[string]string),
	}
}

// Retrieve returns a RetrievePostAutosave builder to get a specific autosave by id.
func (api *PostAutosaves) Retrieve(autosaveID int) *RetrievePostAutosave {
	return &RetrievePostAutosave{
		endpoint:   "/wp/v2/posts/" + strconv.Itoa(api.parentID) + "/autosaves/" + strconv.Itoa(autosaveID),
		client:     api.client,
		arguments:  make(map[string]string),
	}
}

// Create returns a CreatePostRevision builder to create a new autosave for this post.
func (api *PostAutosaves) Create(revision ...PostData) *CreatePostRevision {
	req := &CreatePostRevision{
		endpoint: "/wp/v2/posts/" + strconv.Itoa(api.parentID) + "/autosaves",
		client:   api.client,
	}
	if len(revision) > 0 {
		req.revision = revision[0]
	}
	return req
}

// ListPostAutosaves handles listing post autosaves.
type ListPostAutosaves struct {
	endpoint  string
	client    *RestClient
	arguments map[string]string
}

func (api *ListPostAutosaves) Context(ctx string) *ListPostAutosaves {
	api.arguments["context"] = ctx
	return api
}

func (api *ListPostAutosaves) ContextView() *ListPostAutosaves {
	api.arguments["context"] = "view"
	return api
}

func (api *ListPostAutosaves) ContextEdit() *ListPostAutosaves {
	api.arguments["context"] = "edit"
	return api
}

func (api *ListPostAutosaves) ContextEmbed() *ListPostAutosaves {
	api.arguments["context"] = "embed"
	return api
}

func (api *ListPostAutosaves) Embed() *ListPostAutosaves {
	api.arguments["_embed"] = "true"
	return api
}

func (api *ListPostAutosaves) Fields(fields ...string) *ListPostAutosaves {
	api.arguments["_fields"] = strings.Join(fields, ",")
	return api
}

func (api *ListPostAutosaves) Do() (revisions []Revision, err error) {
	restyClient := api.client.httpClient.R()
	if api.client.auth.Username != "" && api.client.auth.Password != "" {
		restyClient.SetBasicAuth(api.client.auth.Username, api.client.auth.Password)
	}

	resp, err := restyClient.
		SetHeader("Accept", "application/json").
		SetResult(&revisions).
		SetQueryParams(api.arguments).
		Get(api.client.endpoint + api.endpoint)

	if resp.IsError() {
		var wpError WPRestError
		err = json.Unmarshal(resp.Bytes(), &wpError)
		if err != nil {
			return
		}
		return revisions, &wpError
	}

	return
}

// RetrievePostAutosave handles retrieving a specific post autosave.
type RetrievePostAutosave struct {
	endpoint  string
	client    *RestClient
	arguments map[string]string
}

func (api *RetrievePostAutosave) Context(ctx string) *RetrievePostAutosave {
	api.arguments["context"] = ctx
	return api
}

func (api *RetrievePostAutosave) ContextView() *RetrievePostAutosave {
	api.arguments["context"] = "view"
	return api
}

func (api *RetrievePostAutosave) ContextEdit() *RetrievePostAutosave {
	api.arguments["context"] = "edit"
	return api
}

func (api *RetrievePostAutosave) ContextEmbed() *RetrievePostAutosave {
	api.arguments["context"] = "embed"
	return api
}

func (api *RetrievePostAutosave) Embed() *RetrievePostAutosave {
	api.arguments["_embed"] = "true"
	return api
}

func (api *RetrievePostAutosave) Fields(fields ...string) *RetrievePostAutosave {
	api.arguments["_fields"] = strings.Join(fields, ",")
	return api
}

func (api *RetrievePostAutosave) Do() (revision *Revision, err error) {
	restyClient := api.client.httpClient.R()
	if api.client.auth.Username != "" && api.client.auth.Password != "" {
		restyClient.SetBasicAuth(api.client.auth.Username, api.client.auth.Password)
	}

	resp, err := restyClient.
		SetHeader("Accept", "application/json").
		SetResult(&revision).
		SetQueryParams(api.arguments).
		Get(api.client.endpoint + api.endpoint)

	if resp.IsError() {
		var wpError WPRestError
		err = json.Unmarshal(resp.Bytes(), &wpError)
		if err != nil {
			return
		}
		return revision, &wpError
	}

	return
}
