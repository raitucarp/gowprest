package gowprest

import (
	"encoding/json"
	"strconv"
	"strings"
	"time"
)

// Block represents a WordPress Editor Block (reusable block / wp_block post type).
type Block struct {
	Date              *Date          `json:"date,omitempty"`
	DateGMT           *Date          `json:"date_gmt,omitempty"`
	GUID              *Object        `json:"guid,omitempty"`
	ID                int            `json:"id,omitempty"`
	Link              string         `json:"link,omitempty"`
	Modified          *Date          `json:"modified,omitempty"`
	ModifiedGMT       *Date          `json:"modified_gmt,omitempty"`
	Slug              string         `json:"slug,omitempty"`
	Status            PostStatus     `json:"status,omitempty"`
	Type              string         `json:"type,omitempty"`
	Password          string         `json:"password,omitempty"`
	Title             *Object        `json:"title,omitempty"`
	Content           *Object        `json:"content,omitempty"`
	Excerpt           *Object        `json:"excerpt,omitempty"`
	Template          string         `json:"template,omitempty"`
	Meta              map[string]any `json:"meta,omitempty"`
	WPPatternCategory []int          `json:"wp_pattern_category,omitempty"`
	Links             map[string]any `json:"_links,omitempty"`
	Embedded          map[string]any `json:"_embedded,omitempty"`
}

// EditorBlock is an alias for Block.
type EditorBlock = Block

// TitleString returns the title string from either Raw or Rendered.
func (b *Block) TitleString() string {
	if b.Title == nil {
		return ""
	}
	if b.Title.Raw != "" {
		return b.Title.Raw
	}
	return b.Title.Rendered
}

// ContentString returns the content string from either Raw or Rendered.
func (b *Block) ContentString() string {
	if b.Content == nil {
		return ""
	}
	if b.Content.Raw != "" {
		return b.Content.Raw
	}
	return b.Content.Rendered
}

// BlockData represents data used to create or update an Editor Block.
type BlockData struct {
	ID                int            `json:"id,omitempty"`
	Date              *string        `json:"date,omitempty"`
	DateGMT           *string        `json:"date_gmt,omitempty"`
	Slug              string         `json:"slug,omitempty"`
	Status            PostStatus     `json:"status,omitempty"`
	Password          string         `json:"password,omitempty"`
	Title             string         `json:"title,omitempty"`
	Content           string         `json:"content,omitempty"`
	Excerpt           string         `json:"excerpt,omitempty"`
	Template          string         `json:"template,omitempty"`
	Meta              map[string]any `json:"meta,omitempty"`
	WPPatternCategory []int          `json:"wp_pattern_category,omitempty"`
}

// Blocks handles requests to the WordPress Editor Blocks API (/wp/v2/blocks).
type Blocks struct {
	client *RestClient
}

// Blocks returns a Blocks service instance.
func (c *RestClient) Blocks() *Blocks {
	return &Blocks{client: c}
}

// EditorBlocks is a convenience alias for Blocks.
func (c *RestClient) EditorBlocks() *Blocks {
	return c.Blocks()
}

// BlockRevisions is a convenience method to access revisions of a specific block.
func (c *RestClient) BlockRevisions(parentID int) *BlockRevisions {
	return c.Blocks().Revisions(parentID)
}

// Revisions returns a BlockRevisions service instance for the specified block.
func (api *Blocks) Revisions(parentID int) *BlockRevisions {
	return &BlockRevisions{
		client:   api.client,
		parentID: parentID,
	}
}

// List returns a ListBlocks builder to query Editor Blocks collections.
func (api *Blocks) List() *ListBlocks {
	return &ListBlocks{
		endpoint:  "/wp/v2/blocks",
		client:    api.client,
		arguments: make(map[string]string),
	}
}

// Create returns a CreateBlock builder to create a new Editor Block.
func (api *Blocks) Create(block ...BlockData) *CreateBlock {
	builder := &CreateBlock{
		endpoint: "/wp/v2/blocks",
		client:   api.client,
	}
	if len(block) > 0 {
		builder.block = block[0]
	}
	return builder
}

// Retrieve returns a RetrieveBlock builder to fetch an Editor Block by ID.
func (api *Blocks) Retrieve(id int) *RetrieveBlock {
	return &RetrieveBlock{
		endpoint:  "/wp/v2/blocks/" + strconv.Itoa(id),
		client:    api.client,
		arguments: make(map[string]string),
	}
}

// Update returns an UpdateBlock builder to modify an existing Editor Block.
func (api *Blocks) Update(id int, block ...BlockData) *UpdateBlock {
	builder := &UpdateBlock{
		endpoint: "/wp/v2/blocks/" + strconv.Itoa(id),
		client:   api.client,
		block: BlockData{
			ID: id,
		},
	}
	if len(block) > 0 {
		builder.block = block[0]
		builder.block.ID = id
	}
	return builder
}

// Delete returns a DeleteBlock builder to remove an Editor Block.
func (api *Blocks) Delete(id int) *DeleteBlock {
	return &DeleteBlock{
		endpoint: "/wp/v2/blocks",
		client:   api.client,
		id:       id,
	}
}

// ListBlocks handles querying block collections.
type ListBlocks struct {
	endpoint  string
	client    *RestClient
	arguments map[string]string
}

func (api *ListBlocks) Context(ctx string) *ListBlocks {
	api.arguments["context"] = ctx
	return api
}

func (api *ListBlocks) ContextView() *ListBlocks {
	return api.Context("view")
}

func (api *ListBlocks) ContextEdit() *ListBlocks {
	return api.Context("edit")
}

func (api *ListBlocks) ContextEmbed() *ListBlocks {
	return api.Context("embed")
}

func (api *ListBlocks) Page(page int) *ListBlocks {
	api.arguments["page"] = strconv.Itoa(page)
	return api
}

func (api *ListBlocks) PerPage(perPage int) *ListBlocks {
	api.arguments["per_page"] = strconv.Itoa(perPage)
	return api
}

func (api *ListBlocks) Search(query string) *ListBlocks {
	api.arguments["search"] = query
	return api
}

func (api *ListBlocks) After(after any) *ListBlocks {
	switch v := after.(type) {
	case time.Time:
		api.arguments["after"] = v.Format(time.RFC3339)
	case string:
		api.arguments["after"] = v
	}
	return api
}

func (api *ListBlocks) Before(before any) *ListBlocks {
	switch v := before.(type) {
	case time.Time:
		api.arguments["before"] = v.Format(time.RFC3339)
	case string:
		api.arguments["before"] = v
	}
	return api
}

func (api *ListBlocks) ModifiedAfter(modifiedAfter any) *ListBlocks {
	switch v := modifiedAfter.(type) {
	case time.Time:
		api.arguments["modified_after"] = v.Format(time.RFC3339)
	case string:
		api.arguments["modified_after"] = v
	}
	return api
}

func (api *ListBlocks) ModifiedBefore(modifiedBefore any) *ListBlocks {
	switch v := modifiedBefore.(type) {
	case time.Time:
		api.arguments["modified_before"] = v.Format(time.RFC3339)
	case string:
		api.arguments["modified_before"] = v
	}
	return api
}

func (api *ListBlocks) Exclude(ids ...int) *ListBlocks {
	strIDs := make([]string, 0, len(ids))
	for _, id := range ids {
		strIDs = append(strIDs, strconv.Itoa(id))
	}
	api.arguments["exclude"] = strings.Join(strIDs, ",")
	return api
}

func (api *ListBlocks) Include(ids ...int) *ListBlocks {
	strIDs := make([]string, 0, len(ids))
	for _, id := range ids {
		strIDs = append(strIDs, strconv.Itoa(id))
	}
	api.arguments["include"] = strings.Join(strIDs, ",")
	return api
}

func (api *ListBlocks) Offset(offset int) *ListBlocks {
	api.arguments["offset"] = strconv.Itoa(offset)
	return api
}

func (api *ListBlocks) Order(order string) *ListBlocks {
	api.arguments["order"] = order
	return api
}

func (api *ListBlocks) OrderAsc() *ListBlocks {
	return api.Order("asc")
}

func (api *ListBlocks) OrderDesc() *ListBlocks {
	return api.Order("desc")
}

func (api *ListBlocks) OrderBy(orderby string) *ListBlocks {
	api.arguments["orderby"] = orderby
	return api
}

func (api *ListBlocks) OrderByDate() *ListBlocks {
	return api.OrderBy("date")
}

func (api *ListBlocks) OrderByID() *ListBlocks {
	return api.OrderBy("id")
}

func (api *ListBlocks) OrderByTitle() *ListBlocks {
	return api.OrderBy("title")
}

func (api *ListBlocks) OrderBySlug() *ListBlocks {
	return api.OrderBy("slug")
}

func (api *ListBlocks) OrderByModified() *ListBlocks {
	return api.OrderBy("modified")
}

func (api *ListBlocks) OrderByRelevance() *ListBlocks {
	return api.OrderBy("relevance")
}

func (api *ListBlocks) SearchColumns(cols ...string) *ListBlocks {
	api.arguments["search_columns"] = strings.Join(cols, ",")
	return api
}

func (api *ListBlocks) Slug(slugs ...string) *ListBlocks {
	api.arguments["slug"] = strings.Join(slugs, ",")
	return api
}

func (api *ListBlocks) Status(statuses ...string) *ListBlocks {
	api.arguments["status"] = strings.Join(statuses, ",")
	return api
}

func (api *ListBlocks) StatusPublish() *ListBlocks {
	return api.Status(string(StatusPublished))
}

func (api *ListBlocks) StatusDraft() *ListBlocks {
	return api.Status(string(StatusDraft))
}

func (api *ListBlocks) StatusPending() *ListBlocks {
	return api.Status(string(StatusPending))
}

func (api *ListBlocks) StatusPrivate() *ListBlocks {
	return api.Status(string(StatusPrivate))
}

func (api *ListBlocks) StatusFuture() *ListBlocks {
	return api.Status(string(StatusFuture))
}

func (api *ListBlocks) Fields(fields ...string) *ListBlocks {
	api.arguments["_fields"] = strings.Join(fields, ",")
	return api
}

func (api *ListBlocks) Embed() *ListBlocks {
	api.arguments["_embed"] = "true"
	return api
}

func (api *ListBlocks) Do() (blocks []Block, err error) {
	endpoint := api.client.endpoint + api.endpoint

	restyClient := api.client.httpClient.R()
	if api.client.auth.Username != "" && api.client.auth.Password != "" {
		restyClient.SetBasicAuth(api.client.auth.Username, api.client.auth.Password)
	}

	resp, err := restyClient.
		SetHeader("Accept", "application/json").
		SetResult(&blocks).
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

	return blocks, nil
}

// CreateBlock handles creating an Editor Block.
type CreateBlock struct {
	endpoint string
	client   *RestClient
	block    BlockData
}

func (api *CreateBlock) Title(title string) *CreateBlock {
	api.block.Title = title
	return api
}

func (api *CreateBlock) Content(content string) *CreateBlock {
	api.block.Content = content
	return api
}

func (api *CreateBlock) Excerpt(excerpt string) *CreateBlock {
	api.block.Excerpt = excerpt
	return api
}

func (api *CreateBlock) Slug(slug string) *CreateBlock {
	api.block.Slug = slug
	return api
}

func (api *CreateBlock) Status(status PostStatus) *CreateBlock {
	api.block.Status = status
	return api
}

func (api *CreateBlock) StatusPublish() *CreateBlock {
	return api.Status(StatusPublished)
}

func (api *CreateBlock) StatusDraft() *CreateBlock {
	return api.Status(StatusDraft)
}

func (api *CreateBlock) StatusPending() *CreateBlock {
	return api.Status(StatusPending)
}

func (api *CreateBlock) StatusPrivate() *CreateBlock {
	return api.Status(StatusPrivate)
}

func (api *CreateBlock) StatusFuture() *CreateBlock {
	return api.Status(StatusFuture)
}

func (api *CreateBlock) Password(password string) *CreateBlock {
	api.block.Password = password
	return api
}

func (api *CreateBlock) Date(date string) *CreateBlock {
	api.block.Date = &date
	return api
}

func (api *CreateBlock) DateGMT(dateGMT string) *CreateBlock {
	api.block.DateGMT = &dateGMT
	return api
}

func (api *CreateBlock) Template(template string) *CreateBlock {
	api.block.Template = template
	return api
}

func (api *CreateBlock) Meta(meta map[string]any) *CreateBlock {
	api.block.Meta = meta
	return api
}

func (api *CreateBlock) WPPatternCategory(categories ...int) *CreateBlock {
	api.block.WPPatternCategory = categories
	return api
}

func (api *CreateBlock) SetPayload(block BlockData) *CreateBlock {
	api.block = block
	return api
}

func (api *CreateBlock) Do() (block *Block, err error) {
	endpoint := api.client.endpoint + api.endpoint

	restyClient := api.client.httpClient.R()
	if api.client.auth.Username != "" && api.client.auth.Password != "" {
		restyClient.SetBasicAuth(api.client.auth.Username, api.client.auth.Password)
	}

	resp, err := restyClient.
		SetHeader("Content-Type", "application/json").
		SetHeader("Accept", "application/json").
		SetBody(api.block).
		SetResult(&block).
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

	return block, nil
}

// RetrieveBlock handles fetching a single Editor Block.
type RetrieveBlock struct {
	endpoint  string
	client    *RestClient
	arguments map[string]string
}

func (api *RetrieveBlock) Context(ctx string) *RetrieveBlock {
	api.arguments["context"] = ctx
	return api
}

func (api *RetrieveBlock) ContextView() *RetrieveBlock {
	return api.Context("view")
}

func (api *RetrieveBlock) ContextEdit() *RetrieveBlock {
	return api.Context("edit")
}

func (api *RetrieveBlock) ContextEmbed() *RetrieveBlock {
	return api.Context("embed")
}

func (api *RetrieveBlock) Password(password string) *RetrieveBlock {
	api.arguments["password"] = password
	return api
}

func (api *RetrieveBlock) Fields(fields ...string) *RetrieveBlock {
	api.arguments["_fields"] = strings.Join(fields, ",")
	return api
}

func (api *RetrieveBlock) Embed() *RetrieveBlock {
	api.arguments["_embed"] = "true"
	return api
}

func (api *RetrieveBlock) Do() (block *Block, err error) {
	endpoint := api.client.endpoint + api.endpoint

	restyClient := api.client.httpClient.R()
	if api.client.auth.Username != "" && api.client.auth.Password != "" {
		restyClient.SetBasicAuth(api.client.auth.Username, api.client.auth.Password)
	}

	resp, err := restyClient.
		SetHeader("Accept", "application/json").
		SetResult(&block).
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

	return block, nil
}

// UpdateBlock handles modifying an Editor Block.
type UpdateBlock struct {
	endpoint string
	client   *RestClient
	block    BlockData
}

func (api *UpdateBlock) Title(title string) *UpdateBlock {
	api.block.Title = title
	return api
}

func (api *UpdateBlock) Content(content string) *UpdateBlock {
	api.block.Content = content
	return api
}

func (api *UpdateBlock) Excerpt(excerpt string) *UpdateBlock {
	api.block.Excerpt = excerpt
	return api
}

func (api *UpdateBlock) Slug(slug string) *UpdateBlock {
	api.block.Slug = slug
	return api
}

func (api *UpdateBlock) Status(status PostStatus) *UpdateBlock {
	api.block.Status = status
	return api
}

func (api *UpdateBlock) StatusPublish() *UpdateBlock {
	return api.Status(StatusPublished)
}

func (api *UpdateBlock) StatusDraft() *UpdateBlock {
	return api.Status(StatusDraft)
}

func (api *UpdateBlock) StatusPending() *UpdateBlock {
	return api.Status(StatusPending)
}

func (api *UpdateBlock) StatusPrivate() *UpdateBlock {
	return api.Status(StatusPrivate)
}

func (api *UpdateBlock) StatusFuture() *UpdateBlock {
	return api.Status(StatusFuture)
}

func (api *UpdateBlock) Password(password string) *UpdateBlock {
	api.block.Password = password
	return api
}

func (api *UpdateBlock) Date(date string) *UpdateBlock {
	api.block.Date = &date
	return api
}

func (api *UpdateBlock) DateGMT(dateGMT string) *UpdateBlock {
	api.block.DateGMT = &dateGMT
	return api
}

func (api *UpdateBlock) Template(template string) *UpdateBlock {
	api.block.Template = template
	return api
}

func (api *UpdateBlock) Meta(meta map[string]any) *UpdateBlock {
	api.block.Meta = meta
	return api
}

func (api *UpdateBlock) WPPatternCategory(categories ...int) *UpdateBlock {
	api.block.WPPatternCategory = categories
	return api
}

func (api *UpdateBlock) SetPayload(block BlockData) *UpdateBlock {
	api.block = block
	return api
}

func (api *UpdateBlock) Do() (block *Block, err error) {
	endpoint := api.client.endpoint + api.endpoint

	restyClient := api.client.httpClient.R()
	if api.client.auth.Username != "" && api.client.auth.Password != "" {
		restyClient.SetBasicAuth(api.client.auth.Username, api.client.auth.Password)
	}

	resp, err := restyClient.
		SetHeader("Content-Type", "application/json").
		SetHeader("Accept", "application/json").
		SetBody(api.block).
		SetResult(&block).
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

	return block, nil
}

type deleteBlockEnvelope struct {
	Block
	Deleted  bool   `json:"deleted"`
	Previous *Block `json:"previous"`
}

// DeleteBlock handles deleting an Editor Block.
type DeleteBlock struct {
	endpoint string
	client   *RestClient
	id       int
	force    bool
}

// Force sets whether to bypass Trash and force permanent deletion.
func (api *DeleteBlock) Force() *DeleteBlock {
	api.force = true
	return api
}

func (api *DeleteBlock) Do() (block *Block, err error) {
	endpoint := api.client.endpoint + api.endpoint + "/" + strconv.Itoa(api.id)
	var env deleteBlockEnvelope

	restyClient := api.client.httpClient.R()
	if api.client.auth.Username != "" && api.client.auth.Password != "" {
		restyClient.SetBasicAuth(api.client.auth.Username, api.client.auth.Password)
	}

	if api.force {
		restyClient.SetQueryParam("force", "true")
	}

	resp, err := restyClient.
		SetHeader("Accept", "application/json").
		SetResult(&env).
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

	if env.Deleted && env.Previous != nil && env.Previous.ID != 0 {
		return env.Previous, nil
	}

	return &env.Block, nil
}
