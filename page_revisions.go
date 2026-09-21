package gowprest

import (
	"encoding/json"
	"strconv"
	"strings"
)

// PageRevisions anchors revision-related operations for a specific page.
type PageRevisions struct {
	client   *RestClient
	parentID int
}

// List returns a ListPageRevisions struct to list revisions for the parent page.
func (api *PageRevisions) List() *ListPageRevisions {
	return &ListPageRevisions{
		endpoint:  "/wp/v2/pages/" + strconv.Itoa(api.parentID) + "/revisions",
		client:    api.client,
		arguments: make(map[string]string),
	}
}

// Retrieve returns a RetrievePageRevision struct to get a specific revision.
func (api *PageRevisions) Retrieve(revisionID int) *RetrievePageRevision {
	return &RetrievePageRevision{
		endpoint:   "/wp/v2/pages/" + strconv.Itoa(api.parentID) + "/revisions",
		client:     api.client,
		revisionID: revisionID,
		arguments:  make(map[string]string),
	}
}

// Delete returns a DeletePageRevision struct to delete a specific revision.
func (api *PageRevisions) Delete(revisionID int) *DeletePageRevision {
	return &DeletePageRevision{
		endpoint:   "/wp/v2/pages/" + strconv.Itoa(api.parentID) + "/revisions",
		client:     api.client,
		revisionID: revisionID,
		force:      true, // WP page revisions always require force=true
	}
}

// Create returns a CreatePageRevision struct to create a new revision (autosave).
// It can be used directly or with variadic optional PageData for backward compatibility.
func (api *PageRevisions) Create(revision ...PageData) *CreatePageRevision {
	req := &CreatePageRevision{
		endpoint: "/wp/v2/pages/" + strconv.Itoa(api.parentID) + "/autosaves",
		client:   api.client,
	}
	if len(revision) > 0 {
		req.revision = revision[0]
	}
	return req
}

// Autosaves returns an Autosaves endpoint builder for this page.
func (api *PageRevisions) Autosaves() *PageAutosaves {
	return &PageAutosaves{
		client:   api.client,
		parentID: api.parentID,
	}
}

// ListPageRevisions handles listing revisions.
type ListPageRevisions struct {
	endpoint  string
	client    *RestClient
	arguments map[string]string
}

func (api *ListPageRevisions) Context(ctx string) *ListPageRevisions {
	api.arguments["context"] = ctx
	return api
}

func (api *ListPageRevisions) ContextView() *ListPageRevisions {
	api.arguments["context"] = "view"
	return api
}

func (api *ListPageRevisions) ContextEdit() *ListPageRevisions {
	api.arguments["context"] = "edit"
	return api
}

func (api *ListPageRevisions) ContextEmbed() *ListPageRevisions {
	api.arguments["context"] = "embed"
	return api
}

func (api *ListPageRevisions) Page(page int) *ListPageRevisions {
	api.arguments["page"] = strconv.Itoa(page)
	return api
}

func (api *ListPageRevisions) PerPage(perPage int) *ListPageRevisions {
	api.arguments["per_page"] = strconv.Itoa(perPage)
	return api
}

func (api *ListPageRevisions) Search(query string) *ListPageRevisions {
	api.arguments["search"] = query
	return api
}

func (api *ListPageRevisions) Exclude(excludeIDs ...int) *ListPageRevisions {
	excludes := make([]string, 0, len(excludeIDs))
	for _, id := range excludeIDs {
		excludes = append(excludes, strconv.Itoa(id))
	}
	api.arguments["exclude"] = strings.Join(excludes, ",")
	return api
}

func (api *ListPageRevisions) Include(includeIDs ...int) *ListPageRevisions {
	includes := make([]string, 0, len(includeIDs))
	for _, id := range includeIDs {
		includes = append(includes, strconv.Itoa(id))
	}
	api.arguments["include"] = strings.Join(includes, ",")
	return api
}

func (api *ListPageRevisions) Offset(offset int) *ListPageRevisions {
	api.arguments["offset"] = strconv.Itoa(offset)
	return api
}

func (api *ListPageRevisions) Order(order string) *ListPageRevisions {
	api.arguments["order"] = order
	return api
}

func (api *ListPageRevisions) OrderAsc() *ListPageRevisions {
	api.arguments["order"] = "asc"
	return api
}

func (api *ListPageRevisions) OrderDesc() *ListPageRevisions {
	api.arguments["order"] = "desc"
	return api
}

func (api *ListPageRevisions) OrderBy(orderBy string) *ListPageRevisions {
	api.arguments["orderby"] = orderBy
	return api
}

func (api *ListPageRevisions) OrderByDate() *ListPageRevisions {
	api.arguments["orderby"] = "date"
	return api
}

func (api *ListPageRevisions) OrderById() *ListPageRevisions {
	api.arguments["orderby"] = "id"
	return api
}

func (api *ListPageRevisions) OrderByInclude() *ListPageRevisions {
	api.arguments["orderby"] = "include"
	return api
}

func (api *ListPageRevisions) OrderByRelevance() *ListPageRevisions {
	api.arguments["orderby"] = "relevance"
	return api
}

func (api *ListPageRevisions) OrderBySlug() *ListPageRevisions {
	api.arguments["orderby"] = "slug"
	return api
}

func (api *ListPageRevisions) OrderByTitle() *ListPageRevisions {
	api.arguments["orderby"] = "title"
	return api
}

func (api *ListPageRevisions) Embed() *ListPageRevisions {
	api.arguments["_embed"] = "true"
	return api
}

func (api *ListPageRevisions) Fields(fields ...string) *ListPageRevisions {
	api.arguments["_fields"] = strings.Join(fields, ",")
	return api
}

func (api *ListPageRevisions) Do() (revisions []Revision, err error) {
	restyClient := api.client.httpClient.R()
	if api.client.auth.Username != "" && api.client.auth.Password != "" {
		restyClient = restyClient.SetBasicAuth(api.client.auth.Username, api.client.auth.Password)
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

// RetrievePageRevision handles retrieving a specific revision.
type RetrievePageRevision struct {
	endpoint   string
	client     *RestClient
	revisionID int
	arguments  map[string]string
}

func (api *RetrievePageRevision) Context(ctx string) *RetrievePageRevision {
	api.arguments["context"] = ctx
	return api
}

func (api *RetrievePageRevision) ContextView() *RetrievePageRevision {
	api.arguments["context"] = "view"
	return api
}

func (api *RetrievePageRevision) ContextEdit() *RetrievePageRevision {
	api.arguments["context"] = "edit"
	return api
}

func (api *RetrievePageRevision) ContextEmbed() *RetrievePageRevision {
	api.arguments["context"] = "embed"
	return api
}

func (api *RetrievePageRevision) Embed() *RetrievePageRevision {
	api.arguments["_embed"] = "true"
	return api
}

func (api *RetrievePageRevision) Fields(fields ...string) *RetrievePageRevision {
	api.arguments["_fields"] = strings.Join(fields, ",")
	return api
}

func (api *RetrievePageRevision) Do() (revision *Revision, err error) {
	endpoint := api.client.endpoint + api.endpoint + "/" + strconv.Itoa(api.revisionID)

	restyClient := api.client.httpClient.R()
	if api.client.auth.Username != "" && api.client.auth.Password != "" {
		restyClient = restyClient.SetBasicAuth(api.client.auth.Username, api.client.auth.Password)
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

// DeletePageRevision handles deleting a specific revision.
type DeletePageRevision struct {
	endpoint   string
	client     *RestClient
	revisionID int
	force      bool
}

func (api *DeletePageRevision) Force() *DeletePageRevision {
	api.force = true
	return api
}

func (api *DeletePageRevision) Do() (revision Revision, err error) {
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

// CreatePageRevision handles creating a new revision (autosave).
type CreatePageRevision struct {
	endpoint string
	client   *RestClient
	revision PageData
}

func (api *CreatePageRevision) Title(title string) *CreatePageRevision {
	api.revision.Title = title
	return api
}

func (api *CreatePageRevision) Content(content string) *CreatePageRevision {
	api.revision.Content = content
	return api
}

func (api *CreatePageRevision) Excerpt(excerpt string) *CreatePageRevision {
	api.revision.Excerpt = excerpt
	return api
}

func (api *CreatePageRevision) Slug(slug string) *CreatePageRevision {
	api.revision.Slug = slug
	return api
}

func (api *CreatePageRevision) Status(status PostStatus) *CreatePageRevision {
	api.revision.Status = status
	return api
}

func (api *CreatePageRevision) Author(author int) *CreatePageRevision {
	api.revision.Author = author
	return api
}

func (api *CreatePageRevision) Password(password string) *CreatePageRevision {
	api.revision.Password = password
	return api
}

func (api *CreatePageRevision) FeaturedMedia(mediaID int) *CreatePageRevision {
	api.revision.FeaturedMedia = mediaID
	return api
}

func (api *CreatePageRevision) CommentStatus(status OpenClosedStatus) *CreatePageRevision {
	api.revision.CommentStatus = status
	return api
}

func (api *CreatePageRevision) PingStatus(status OpenClosedStatus) *CreatePageRevision {
	api.revision.PingStatus = status
	return api
}

func (api *CreatePageRevision) MenuOrder(order int) *CreatePageRevision {
	api.revision.MenuOrder = order
	return api
}

func (api *CreatePageRevision) Meta(meta map[string]any) *CreatePageRevision {
	api.revision.Meta = meta
	return api
}

func (api *CreatePageRevision) Template(template string) *CreatePageRevision {
	api.revision.Template = template
	return api
}

func (api *CreatePageRevision) Parent(parentID int) *CreatePageRevision {
	api.revision.Parent = parentID
	return api
}

func (api *CreatePageRevision) Do() (revision Revision, err error) {
	restyClient := api.client.httpClient.R()
	if api.client.auth.Username != "" && api.client.auth.Password != "" {
		restyClient = restyClient.SetBasicAuth(api.client.auth.Username, api.client.auth.Password)
	}

	resp, err := restyClient.
		SetHeader("Content-Type", "application/json").
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

// PageAutosaves provides operations for page autosaves.
type PageAutosaves struct {
	client   *RestClient
	parentID int
}

// List returns a ListPageAutosaves builder to get all autosaves for this page.
func (api *PageAutosaves) List() *ListPageAutosaves {
	return &ListPageAutosaves{
		endpoint:  "/wp/v2/pages/" + strconv.Itoa(api.parentID) + "/autosaves",
		client:    api.client,
		arguments: make(map[string]string),
	}
}

// Retrieve returns a RetrievePageAutosave builder to get a specific autosave by id.
func (api *PageAutosaves) Retrieve(autosaveID int) *RetrievePageAutosave {
	return &RetrievePageAutosave{
		endpoint:   "/wp/v2/pages/" + strconv.Itoa(api.parentID) + "/autosaves/" + strconv.Itoa(autosaveID),
		client:     api.client,
		arguments:  make(map[string]string),
	}
}

// Create returns a CreatePageRevision builder to create a new autosave for this page.
func (api *PageAutosaves) Create(revision ...PageData) *CreatePageRevision {
	req := &CreatePageRevision{
		endpoint: "/wp/v2/pages/" + strconv.Itoa(api.parentID) + "/autosaves",
		client:   api.client,
	}
	if len(revision) > 0 {
		req.revision = revision[0]
	}
	return req
}

// ListPageAutosaves handles listing page autosaves.
type ListPageAutosaves struct {
	endpoint  string
	client    *RestClient
	arguments map[string]string
}

func (api *ListPageAutosaves) Context(ctx string) *ListPageAutosaves {
	api.arguments["context"] = ctx
	return api
}

func (api *ListPageAutosaves) ContextView() *ListPageAutosaves {
	api.arguments["context"] = "view"
	return api
}

func (api *ListPageAutosaves) ContextEdit() *ListPageAutosaves {
	api.arguments["context"] = "edit"
	return api
}

func (api *ListPageAutosaves) ContextEmbed() *ListPageAutosaves {
	api.arguments["context"] = "embed"
	return api
}

func (api *ListPageAutosaves) Embed() *ListPageAutosaves {
	api.arguments["_embed"] = "true"
	return api
}

func (api *ListPageAutosaves) Fields(fields ...string) *ListPageAutosaves {
	api.arguments["_fields"] = strings.Join(fields, ",")
	return api
}

func (api *ListPageAutosaves) Do() (revisions []Revision, err error) {
	restyClient := api.client.httpClient.R()
	if api.client.auth.Username != "" && api.client.auth.Password != "" {
		restyClient = restyClient.SetBasicAuth(api.client.auth.Username, api.client.auth.Password)
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

// RetrievePageAutosave handles retrieving a specific page autosave.
type RetrievePageAutosave struct {
	endpoint  string
	client    *RestClient
	arguments map[string]string
}

func (api *RetrievePageAutosave) Context(ctx string) *RetrievePageAutosave {
	api.arguments["context"] = ctx
	return api
}

func (api *RetrievePageAutosave) ContextView() *RetrievePageAutosave {
	api.arguments["context"] = "view"
	return api
}

func (api *RetrievePageAutosave) ContextEdit() *RetrievePageAutosave {
	api.arguments["context"] = "edit"
	return api
}

func (api *RetrievePageAutosave) ContextEmbed() *RetrievePageAutosave {
	api.arguments["context"] = "embed"
	return api
}

func (api *RetrievePageAutosave) Embed() *RetrievePageAutosave {
	api.arguments["_embed"] = "true"
	return api
}

func (api *RetrievePageAutosave) Fields(fields ...string) *RetrievePageAutosave {
	api.arguments["_fields"] = strings.Join(fields, ",")
	return api
}

func (api *RetrievePageAutosave) Do() (revision *Revision, err error) {
	restyClient := api.client.httpClient.R()
	if api.client.auth.Username != "" && api.client.auth.Password != "" {
		restyClient = restyClient.SetBasicAuth(api.client.auth.Username, api.client.auth.Password)
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
