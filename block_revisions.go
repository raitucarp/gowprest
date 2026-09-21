package gowprest

import (
	"encoding/json"
	"strconv"
	"strings"
)

// BlockRevision represents a WordPress block (wp_block) revision.
type BlockRevision struct {
	Author      int            `json:"author,omitempty"`
	Date        *Date          `json:"date,omitempty"`
	DateGMT     *Date          `json:"date_gmt,omitempty"`
	GUID        *Object        `json:"guid,omitempty"`
	ID          int            `json:"id,omitempty"`
	Modified    *Date          `json:"modified,omitempty"`
	ModifiedGMT *Date          `json:"modified_gmt,omitempty"`
	Parent      int            `json:"parent,omitempty"`
	Slug        string         `json:"slug,omitempty"`
	Title       *Object        `json:"title,omitempty"`
	Content     *Object        `json:"content,omitempty"`
	Excerpt     *Object        `json:"excerpt,omitempty"`
	Links       map[string]any `json:"_links,omitempty"`
	Embedded    map[string]any `json:"_embedded,omitempty"`
}

// TitleString returns the revision's title string from either Raw or Rendered.
func (r *BlockRevision) TitleString() string {
	if r.Title == nil {
		return ""
	}
	if r.Title.Raw != "" {
		return r.Title.Raw
	}
	return r.Title.Rendered
}

// ContentString returns the revision's content string from either Raw or Rendered.
func (r *BlockRevision) ContentString() string {
	if r.Content == nil {
		return ""
	}
	if r.Content.Raw != "" {
		return r.Content.Raw
	}
	return r.Content.Rendered
}

// BlockRevisions anchors revision-related operations for a specific Editor Block.
type BlockRevisions struct {
	client   *RestClient
	parentID int
}

// List returns a ListBlockRevisions struct to list revisions for the parent block.
func (api *BlockRevisions) List() *ListBlockRevisions {
	return &ListBlockRevisions{
		endpoint:  "/wp/v2/blocks/" + strconv.Itoa(api.parentID) + "/revisions",
		client:    api.client,
		arguments: make(map[string]string),
	}
}

// Retrieve returns a RetrieveBlockRevision struct to get a specific revision.
func (api *BlockRevisions) Retrieve(revisionID int) *RetrieveBlockRevision {
	return &RetrieveBlockRevision{
		endpoint:   "/wp/v2/blocks/" + strconv.Itoa(api.parentID) + "/revisions",
		client:     api.client,
		revisionID: revisionID,
		arguments:  make(map[string]string),
	}
}

// Delete returns a DeleteBlockRevision struct to delete a specific revision.
func (api *BlockRevisions) Delete(revisionID int) *DeleteBlockRevision {
	return &DeleteBlockRevision{
		endpoint:   "/wp/v2/blocks/" + strconv.Itoa(api.parentID) + "/revisions",
		client:     api.client,
		revisionID: revisionID,
		force:      true, // Block revisions always require force=true
	}
}

// Create returns a CreateBlockRevision builder to create a new revision (autosave).
func (api *BlockRevisions) Create(revision ...BlockData) *CreateBlockRevision {
	req := &CreateBlockRevision{
		endpoint: "/wp/v2/blocks/" + strconv.Itoa(api.parentID) + "/autosaves",
		client:   api.client,
	}
	if len(revision) > 0 {
		req.revision = revision[0]
	}
	return req
}

// Autosaves returns an Autosaves endpoint builder for this block.
func (api *BlockRevisions) Autosaves() *BlockAutosaves {
	return &BlockAutosaves{
		client:   api.client,
		parentID: api.parentID,
	}
}

// ListBlockRevisions handles listing block revisions.
type ListBlockRevisions struct {
	endpoint  string
	client    *RestClient
	arguments map[string]string
}

func (api *ListBlockRevisions) Context(ctx string) *ListBlockRevisions {
	api.arguments["context"] = ctx
	return api
}

func (api *ListBlockRevisions) ContextView() *ListBlockRevisions {
	return api.Context("view")
}

func (api *ListBlockRevisions) ContextEdit() *ListBlockRevisions {
	return api.Context("edit")
}

func (api *ListBlockRevisions) ContextEmbed() *ListBlockRevisions {
	return api.Context("embed")
}

func (api *ListBlockRevisions) Page(page int) *ListBlockRevisions {
	api.arguments["page"] = strconv.Itoa(page)
	return api
}

func (api *ListBlockRevisions) PerPage(perPage int) *ListBlockRevisions {
	api.arguments["per_page"] = strconv.Itoa(perPage)
	return api
}

func (api *ListBlockRevisions) Search(query string) *ListBlockRevisions {
	api.arguments["search"] = query
	return api
}

func (api *ListBlockRevisions) Exclude(excludeIDs ...int) *ListBlockRevisions {
	ids := make([]string, 0, len(excludeIDs))
	for _, id := range excludeIDs {
		ids = append(ids, strconv.Itoa(id))
	}
	api.arguments["exclude"] = strings.Join(ids, ",")
	return api
}

func (api *ListBlockRevisions) Include(includeIDs ...int) *ListBlockRevisions {
	ids := make([]string, 0, len(includeIDs))
	for _, id := range includeIDs {
		ids = append(ids, strconv.Itoa(id))
	}
	api.arguments["include"] = strings.Join(ids, ",")
	return api
}

func (api *ListBlockRevisions) Offset(offset int) *ListBlockRevisions {
	api.arguments["offset"] = strconv.Itoa(offset)
	return api
}

func (api *ListBlockRevisions) Order(order string) *ListBlockRevisions {
	api.arguments["order"] = order
	return api
}

func (api *ListBlockRevisions) OrderAsc() *ListBlockRevisions {
	return api.Order("asc")
}

func (api *ListBlockRevisions) OrderDesc() *ListBlockRevisions {
	return api.Order("desc")
}

func (api *ListBlockRevisions) OrderBy(orderBy string) *ListBlockRevisions {
	api.arguments["orderby"] = orderBy
	return api
}

func (api *ListBlockRevisions) OrderByDate() *ListBlockRevisions {
	return api.OrderBy("date")
}

func (api *ListBlockRevisions) OrderByID() *ListBlockRevisions {
	return api.OrderBy("id")
}

func (api *ListBlockRevisions) OrderByInclude() *ListBlockRevisions {
	return api.OrderBy("include")
}

func (api *ListBlockRevisions) OrderByRelevance() *ListBlockRevisions {
	return api.OrderBy("relevance")
}

func (api *ListBlockRevisions) OrderBySlug() *ListBlockRevisions {
	return api.OrderBy("slug")
}

func (api *ListBlockRevisions) OrderByTitle() *ListBlockRevisions {
	return api.OrderBy("title")
}

func (api *ListBlockRevisions) Embed() *ListBlockRevisions {
	api.arguments["_embed"] = "true"
	return api
}

func (api *ListBlockRevisions) Fields(fields ...string) *ListBlockRevisions {
	api.arguments["_fields"] = strings.Join(fields, ",")
	return api
}

func (api *ListBlockRevisions) Do() (revisions []BlockRevision, err error) {
	endpoint := api.client.endpoint + api.endpoint

	restyClient := api.client.httpClient.R()
	if api.client.auth.Username != "" && api.client.auth.Password != "" {
		restyClient = restyClient.SetBasicAuth(api.client.auth.Username, api.client.auth.Password)
	}

	resp, err := restyClient.
		SetHeader("Accept", "application/json").
		SetResult(&revisions).
		SetQueryParams(api.arguments).
		Get(endpoint)

	if err != nil {
		return nil, err
	}

	if resp.IsError() {
		var wpError WPRestError
		err = json.Unmarshal(resp.Bytes(), &wpError)
		if err != nil {
			return nil, err
		}
		return nil, &wpError
	}

	return revisions, nil
}

// RetrieveBlockRevision handles getting a specific block revision.
type RetrieveBlockRevision struct {
	endpoint   string
	client     *RestClient
	revisionID int
	arguments  map[string]string
}

func (api *RetrieveBlockRevision) Context(ctx string) *RetrieveBlockRevision {
	api.arguments["context"] = ctx
	return api
}

func (api *RetrieveBlockRevision) ContextView() *RetrieveBlockRevision {
	return api.Context("view")
}

func (api *RetrieveBlockRevision) ContextEdit() *RetrieveBlockRevision {
	return api.Context("edit")
}

func (api *RetrieveBlockRevision) ContextEmbed() *RetrieveBlockRevision {
	return api.Context("embed")
}

func (api *RetrieveBlockRevision) Embed() *RetrieveBlockRevision {
	api.arguments["_embed"] = "true"
	return api
}

func (api *RetrieveBlockRevision) Fields(fields ...string) *RetrieveBlockRevision {
	api.arguments["_fields"] = strings.Join(fields, ",")
	return api
}

func (api *RetrieveBlockRevision) Do() (revision *BlockRevision, err error) {
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

	if err != nil {
		return nil, err
	}

	if resp.IsError() {
		var wpError WPRestError
		err = json.Unmarshal(resp.Bytes(), &wpError)
		if err != nil {
			return nil, err
		}
		return nil, &wpError
	}

	return revision, nil
}

type deleteBlockRevisionEnvelope struct {
	BlockRevision
	Deleted  bool           `json:"deleted"`
	Previous *BlockRevision `json:"previous"`
}

// DeleteBlockRevision handles deleting a specific revision.
type DeleteBlockRevision struct {
	endpoint   string
	client     *RestClient
	revisionID int
	force      bool
}

func (api *DeleteBlockRevision) Force() *DeleteBlockRevision {
	api.force = true
	return api
}

func (api *DeleteBlockRevision) Do() (revision *BlockRevision, err error) {
	endpoint := api.client.endpoint + api.endpoint + "/" + strconv.Itoa(api.revisionID)
	restyClient := api.client.httpClient.R()

	if api.client.auth.Username != "" && api.client.auth.Password != "" {
		restyClient = restyClient.SetBasicAuth(api.client.auth.Username, api.client.auth.Password)
	}

	var env deleteBlockRevisionEnvelope
	resp, err := restyClient.
		SetHeader("Accept", "application/json").
		SetResult(&env).
		SetQueryParam("force", strconv.FormatBool(api.force)).
		Delete(endpoint)

	if err != nil {
		return nil, err
	}

	if resp.IsError() {
		var wpError WPRestError
		err = json.Unmarshal(resp.Bytes(), &wpError)
		if err != nil {
			return nil, err
		}
		return nil, &wpError
	}

	if env.Previous != nil && env.Previous.ID != 0 {
		return env.Previous, nil
	}

	return &env.BlockRevision, nil
}

// CreateBlockRevision handles creating a new revision (autosave).
type CreateBlockRevision struct {
	endpoint string
	client   *RestClient
	revision BlockData
}

func (api *CreateBlockRevision) Title(title string) *CreateBlockRevision {
	api.revision.Title = title
	return api
}

func (api *CreateBlockRevision) Content(content string) *CreateBlockRevision {
	api.revision.Content = content
	return api
}

func (api *CreateBlockRevision) Excerpt(excerpt string) *CreateBlockRevision {
	api.revision.Excerpt = excerpt
	return api
}

func (api *CreateBlockRevision) Slug(slug string) *CreateBlockRevision {
	api.revision.Slug = slug
	return api
}

func (api *CreateBlockRevision) Status(status PostStatus) *CreateBlockRevision {
	api.revision.Status = status
	return api
}

func (api *CreateBlockRevision) StatusPublish() *CreateBlockRevision {
	return api.Status(StatusPublished)
}

func (api *CreateBlockRevision) StatusDraft() *CreateBlockRevision {
	return api.Status(StatusDraft)
}

func (api *CreateBlockRevision) Password(password string) *CreateBlockRevision {
	api.revision.Password = password
	return api
}

func (api *CreateBlockRevision) Date(date string) *CreateBlockRevision {
	api.revision.Date = &date
	return api
}

func (api *CreateBlockRevision) DateGMT(dateGMT string) *CreateBlockRevision {
	api.revision.DateGMT = &dateGMT
	return api
}

func (api *CreateBlockRevision) Template(template string) *CreateBlockRevision {
	api.revision.Template = template
	return api
}

func (api *CreateBlockRevision) Meta(meta map[string]any) *CreateBlockRevision {
	api.revision.Meta = meta
	return api
}

func (api *CreateBlockRevision) WPPatternCategory(categories ...int) *CreateBlockRevision {
	api.revision.WPPatternCategory = categories
	return api
}

func (api *CreateBlockRevision) Do() (revision *BlockRevision, err error) {
	endpoint := api.client.endpoint + api.endpoint

	restyClient := api.client.httpClient.R()
	if api.client.auth.Username != "" && api.client.auth.Password != "" {
		restyClient = restyClient.SetBasicAuth(api.client.auth.Username, api.client.auth.Password)
	}

	resp, err := restyClient.
		SetHeader("Content-Type", "application/json").
		SetHeader("Accept", "application/json").
		SetBody(api.revision).
		SetResult(&revision).
		Post(endpoint)

	if err != nil {
		return nil, err
	}

	if resp.IsError() {
		var wpError WPRestError
		err = json.Unmarshal(resp.Bytes(), &wpError)
		if err != nil {
			return nil, err
		}
		return nil, &wpError
	}

	return revision, nil
}

// BlockAutosaves provides operations for block autosaves.
type BlockAutosaves struct {
	client   *RestClient
	parentID int
}

// List returns a ListBlockAutosaves builder to get all autosaves for this block.
func (api *BlockAutosaves) List() *ListBlockAutosaves {
	return &ListBlockAutosaves{
		endpoint:  "/wp/v2/blocks/" + strconv.Itoa(api.parentID) + "/autosaves",
		client:    api.client,
		arguments: make(map[string]string),
	}
}

// Retrieve returns a RetrieveBlockAutosave builder to get a specific autosave by id.
func (api *BlockAutosaves) Retrieve(autosaveID int) *RetrieveBlockAutosave {
	return &RetrieveBlockAutosave{
		endpoint:   "/wp/v2/blocks/" + strconv.Itoa(api.parentID) + "/autosaves/" + strconv.Itoa(autosaveID),
		client:     api.client,
		arguments:  make(map[string]string),
	}
}

// Create returns a CreateBlockRevision builder to create a new autosave for this block.
func (api *BlockAutosaves) Create(revision ...BlockData) *CreateBlockRevision {
	req := &CreateBlockRevision{
		endpoint: "/wp/v2/blocks/" + strconv.Itoa(api.parentID) + "/autosaves",
		client:   api.client,
	}
	if len(revision) > 0 {
		req.revision = revision[0]
	}
	return req
}

// ListBlockAutosaves handles listing block autosaves.
type ListBlockAutosaves struct {
	endpoint  string
	client    *RestClient
	arguments map[string]string
}

func (api *ListBlockAutosaves) Context(ctx string) *ListBlockAutosaves {
	api.arguments["context"] = ctx
	return api
}

func (api *ListBlockAutosaves) ContextView() *ListBlockAutosaves {
	return api.Context("view")
}

func (api *ListBlockAutosaves) ContextEdit() *ListBlockAutosaves {
	return api.Context("edit")
}

func (api *ListBlockAutosaves) ContextEmbed() *ListBlockAutosaves {
	return api.Context("embed")
}

func (api *ListBlockAutosaves) Embed() *ListBlockAutosaves {
	api.arguments["_embed"] = "true"
	return api
}

func (api *ListBlockAutosaves) Fields(fields ...string) *ListBlockAutosaves {
	api.arguments["_fields"] = strings.Join(fields, ",")
	return api
}

func (api *ListBlockAutosaves) Do() (revisions []BlockRevision, err error) {
	endpoint := api.client.endpoint + api.endpoint

	restyClient := api.client.httpClient.R()
	if api.client.auth.Username != "" && api.client.auth.Password != "" {
		restyClient = restyClient.SetBasicAuth(api.client.auth.Username, api.client.auth.Password)
	}

	resp, err := restyClient.
		SetHeader("Accept", "application/json").
		SetResult(&revisions).
		SetQueryParams(api.arguments).
		Get(endpoint)

	if err != nil {
		return nil, err
	}

	if resp.IsError() {
		var wpError WPRestError
		err = json.Unmarshal(resp.Bytes(), &wpError)
		if err != nil {
			return nil, err
		}
		return nil, &wpError
	}

	return revisions, nil
}

// RetrieveBlockAutosave handles retrieving a specific block autosave.
type RetrieveBlockAutosave struct {
	endpoint  string
	client    *RestClient
	arguments map[string]string
}

func (api *RetrieveBlockAutosave) Context(ctx string) *RetrieveBlockAutosave {
	api.arguments["context"] = ctx
	return api
}

func (api *RetrieveBlockAutosave) ContextView() *RetrieveBlockAutosave {
	return api.Context("view")
}

func (api *RetrieveBlockAutosave) ContextEdit() *RetrieveBlockAutosave {
	return api.Context("edit")
}

func (api *RetrieveBlockAutosave) ContextEmbed() *RetrieveBlockAutosave {
	return api.Context("embed")
}

func (api *RetrieveBlockAutosave) Embed() *RetrieveBlockAutosave {
	api.arguments["_embed"] = "true"
	return api
}

func (api *RetrieveBlockAutosave) Fields(fields ...string) *RetrieveBlockAutosave {
	api.arguments["_fields"] = strings.Join(fields, ",")
	return api
}

func (api *RetrieveBlockAutosave) Do() (revision *BlockRevision, err error) {
	endpoint := api.client.endpoint + api.endpoint

	restyClient := api.client.httpClient.R()
	if api.client.auth.Username != "" && api.client.auth.Password != "" {
		restyClient = restyClient.SetBasicAuth(api.client.auth.Username, api.client.auth.Password)
	}

	resp, err := restyClient.
		SetHeader("Accept", "application/json").
		SetResult(&revision).
		SetQueryParams(api.arguments).
		Get(endpoint)

	if err != nil {
		return nil, err
	}

	if resp.IsError() {
		var wpError WPRestError
		err = json.Unmarshal(resp.Bytes(), &wpError)
		if err != nil {
			return nil, err
		}
		return nil, &wpError
	}

	return revision, nil
}
