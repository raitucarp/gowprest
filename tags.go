package gowprest

import (
	"encoding/json"
	"strconv"
	"strings"
)

// Tag represents a WordPress post_tag term.
type Tag struct {
	ID          int            `json:"id,omitempty"`
	Count       int            `json:"count,omitempty"`
	Description string         `json:"description,omitempty"`
	Link        string         `json:"link,omitempty"`
	Name        string         `json:"name,omitempty"`
	Slug        string         `json:"slug,omitempty"`
	Taxonomy    string         `json:"taxonomy,omitempty"`
	Meta        any            `json:"meta,omitempty"`
	Links       map[string]any `json:"_links,omitempty"`
	Embedded    map[string]any `json:"_embedded,omitempty"`
}

// TagData represents payload for creating or updating a tag.
type TagData struct {
	ID          int    `json:"id,omitempty"`
	Description string `json:"description,omitempty"`
	Name        string `json:"name,omitempty"`
	Slug        string `json:"slug,omitempty"`
	Meta        any    `json:"meta,omitempty"`
}

// Tags anchors tag-related operations.
type Tags struct {
	client *RestClient
}

// Tags returns a Tags API client.
func (c *RestClient) Tags() *Tags {
	return &Tags{client: c}
}

// ListTags handles listing tags.
type ListTags struct {
	endpoint  string
	client    *RestClient
	arguments map[string]string
}

// List returns a ListTags builder.
func (api *Tags) List() *ListTags {
	return &ListTags{
		endpoint:  "/wp/v2/tags",
		client:    api.client,
		arguments: make(map[string]string),
	}
}

func (api *ListTags) Context(ctx string) *ListTags {
	api.arguments["context"] = ctx
	return api
}

func (api *ListTags) ContextView() *ListTags {
	api.arguments["context"] = "view"
	return api
}

func (api *ListTags) ContextEdit() *ListTags {
	api.arguments["context"] = "edit"
	return api
}

func (api *ListTags) ContextEmbed() *ListTags {
	api.arguments["context"] = "embed"
	return api
}

func (api *ListTags) Page(page int) *ListTags {
	api.arguments["page"] = strconv.Itoa(page)
	return api
}

func (api *ListTags) PerPage(perPage int) *ListTags {
	api.arguments["per_page"] = strconv.Itoa(perPage)
	return api
}

func (api *ListTags) Search(query string) *ListTags {
	api.arguments["search"] = query
	return api
}

func (api *ListTags) Exclude(excludeIDs ...int) *ListTags {
	excludes := make([]string, 0, len(excludeIDs))
	for _, id := range excludeIDs {
		excludes = append(excludes, strconv.Itoa(id))
	}
	api.arguments["exclude"] = strings.Join(excludes, ",")
	return api
}

func (api *ListTags) Include(includeIDs ...int) *ListTags {
	includes := make([]string, 0, len(includeIDs))
	for _, id := range includeIDs {
		includes = append(includes, strconv.Itoa(id))
	}
	api.arguments["include"] = strings.Join(includes, ",")
	return api
}

func (api *ListTags) Offset(offset int) *ListTags {
	api.arguments["offset"] = strconv.Itoa(offset)
	return api
}

func (api *ListTags) Order(order string) *ListTags {
	api.arguments["order"] = order
	return api
}

func (api *ListTags) OrderAsc() *ListTags {
	api.arguments["order"] = "asc"
	return api
}

func (api *ListTags) OrderDesc() *ListTags {
	api.arguments["order"] = "desc"
	return api
}

func (api *ListTags) OrderBy(orderBy string) *ListTags {
	api.arguments["orderby"] = orderBy
	return api
}

func (api *ListTags) OrderById() *ListTags {
	api.arguments["orderby"] = "id"
	return api
}

func (api *ListTags) OrderByInclude() *ListTags {
	api.arguments["orderby"] = "include"
	return api
}

func (api *ListTags) OrderByName() *ListTags {
	api.arguments["orderby"] = "name"
	return api
}

func (api *ListTags) OrderBySlug() *ListTags {
	api.arguments["orderby"] = "slug"
	return api
}

func (api *ListTags) OrderByIncludeSlug() *ListTags {
	api.arguments["orderby"] = "include_slugs"
	return api
}

func (api *ListTags) OrderByTermGroup() *ListTags {
	api.arguments["orderby"] = "term_group"
	return api
}

func (api *ListTags) OrderByDescription() *ListTags {
	api.arguments["orderby"] = "description"
	return api
}

func (api *ListTags) OrderByCount() *ListTags {
	api.arguments["orderby"] = "count"
	return api
}

func (api *ListTags) HideEmpty(hide bool) *ListTags {
	api.arguments["hide_empty"] = strconv.FormatBool(hide)
	return api
}

func (api *ListTags) Post(postID int) *ListTags {
	api.arguments["post"] = strconv.Itoa(postID)
	return api
}

func (api *ListTags) Slug(slugs ...string) *ListTags {
	api.arguments["slug"] = strings.Join(slugs, ",")
	return api
}

func (api *ListTags) Embed() *ListTags {
	api.arguments["_embed"] = "true"
	return api
}

func (api *ListTags) Fields(fields ...string) *ListTags {
	api.arguments["_fields"] = strings.Join(fields, ",")
	return api
}

func (api *ListTags) Do() (tags []Tag, err error) {
	restyClient := api.client.httpClient.R()
	if api.client.auth.Username != "" && api.client.auth.Password != "" {
		restyClient.SetBasicAuth(api.client.auth.Username, api.client.auth.Password)
	}

	resp, err := restyClient.
		SetHeader("Accept", "application/json").
		SetResult(&tags).
		SetQueryParams(api.arguments).
		Get(api.client.endpoint + api.endpoint)

	if resp.IsError() {
		var wpError WPRestError
		err = json.Unmarshal(resp.Bytes(), &wpError)
		if err != nil {
			return
		}
		return tags, &wpError
	}

	return
}

// CreateTag handles tag creation.
type CreateTag struct {
	endpoint string
	client   *RestClient
	tag      TagData
}

// Create returns a CreateTag builder.
func (api *Tags) Create(tag ...TagData) *CreateTag {
	builder := &CreateTag{
		endpoint: "/wp/v2/tags",
		client:   api.client,
	}
	if len(tag) > 0 {
		builder.tag = tag[0]
	}
	return builder
}

func (api *CreateTag) Name(name string) *CreateTag {
	api.tag.Name = name
	return api
}

func (api *CreateTag) Description(description string) *CreateTag {
	api.tag.Description = description
	return api
}

func (api *CreateTag) Slug(slug string) *CreateTag {
	api.tag.Slug = slug
	return api
}

func (api *CreateTag) Meta(meta any) *CreateTag {
	api.tag.Meta = meta
	return api
}

func (api *CreateTag) Do() (tag Tag, err error) {
	resp, err := api.client.httpClient.R().
		SetHeader("Content-Type", "application/json").
		SetBasicAuth(api.client.auth.Username, api.client.auth.Password).
		SetResult(&tag).
		SetBody(api.tag).
		Post(api.client.endpoint + api.endpoint)

	if resp.IsError() {
		var wpError WPRestError
		err = json.Unmarshal(resp.Bytes(), &wpError)
		if err != nil {
			return
		}
		return tag, &wpError
	}

	return
}

// RetrieveTag handles retrieving a single tag.
type RetrieveTag struct {
	endpoint  string
	client    *RestClient
	arguments map[string]string
}

// Retrieve returns a RetrieveTag builder.
func (api *Tags) Retrieve(tagId int) *RetrieveTag {
	return &RetrieveTag{
		endpoint:  "/wp/v2/tags/" + strconv.Itoa(tagId),
		client:    api.client,
		arguments: make(map[string]string),
	}
}

func (api *RetrieveTag) Context(ctx string) *RetrieveTag {
	api.arguments["context"] = ctx
	return api
}

func (api *RetrieveTag) ContextView() *RetrieveTag {
	api.arguments["context"] = "view"
	return api
}

func (api *RetrieveTag) ContextEdit() *RetrieveTag {
	api.arguments["context"] = "edit"
	return api
}

func (api *RetrieveTag) ContextEmbed() *RetrieveTag {
	api.arguments["context"] = "embed"
	return api
}

func (api *RetrieveTag) Embed() *RetrieveTag {
	api.arguments["_embed"] = "true"
	return api
}

func (api *RetrieveTag) Fields(fields ...string) *RetrieveTag {
	api.arguments["_fields"] = strings.Join(fields, ",")
	return api
}

func (api *RetrieveTag) Do() (tag *Tag, err error) {
	endpoint := api.client.endpoint + api.endpoint

	restyClient := api.client.httpClient.R()
	if api.client.auth.Username != "" && api.client.auth.Password != "" {
		restyClient.SetBasicAuth(api.client.auth.Username, api.client.auth.Password)
	}

	resp, err := restyClient.
		SetHeader("Accept", "application/json").
		SetResult(&tag).
		SetQueryParams(api.arguments).
		Get(endpoint)

	if resp.IsError() {
		var wpError WPRestError
		err = json.Unmarshal(resp.Bytes(), &wpError)
		if err != nil {
			return
		}
		return tag, &wpError
	}

	return
}

// UpdateTag handles updating a tag.
type UpdateTag struct {
	endpoint string
	client   *RestClient
	tag      TagData
}

// Update returns an UpdateTag builder.
func (api *Tags) Update(tag ...TagData) *UpdateTag {
	builder := &UpdateTag{
		endpoint: "/wp/v2/tags",
		client:   api.client,
	}
	if len(tag) > 0 {
		builder.tag = tag[0]
		if builder.tag.ID != 0 {
			builder.endpoint = "/wp/v2/tags/" + strconv.Itoa(builder.tag.ID)
		}
	}
	return builder
}

func (api *UpdateTag) ID(id int) *UpdateTag {
	api.tag.ID = id
	api.endpoint = "/wp/v2/tags/" + strconv.Itoa(id)
	return api
}

func (api *UpdateTag) Name(name string) *UpdateTag {
	api.tag.Name = name
	return api
}

func (api *UpdateTag) Description(description string) *UpdateTag {
	api.tag.Description = description
	return api
}

func (api *UpdateTag) Slug(slug string) *UpdateTag {
	api.tag.Slug = slug
	return api
}

func (api *UpdateTag) Meta(meta any) *UpdateTag {
	api.tag.Meta = meta
	return api
}

func (api *UpdateTag) Do() (tag Tag, err error) {
	resp, err := api.client.httpClient.R().
		SetHeader("Content-Type", "application/json").
		SetBasicAuth(api.client.auth.Username, api.client.auth.Password).
		SetResult(&tag).
		SetBody(api.tag).
		Post(api.client.endpoint + api.endpoint)

	if resp.IsError() {
		var wpError WPRestError
		err = json.Unmarshal(resp.Bytes(), &wpError)
		if err != nil {
			return
		}
		return tag, &wpError
	}

	return
}

type deleteTagEnvelope struct {
	Tag
	Deleted  bool `json:"deleted"`
	Previous *Tag `json:"previous"`
}

// DeleteTag handles deleting a tag.
type DeleteTag struct {
	endpoint string
	client   *RestClient
	tagId    int
	force    bool
}

// Delete returns a DeleteTag builder.
func (api *Tags) Delete(tagId int) *DeleteTag {
	return &DeleteTag{
		endpoint: "/wp/v2/tags",
		client:   api.client,
		tagId:    tagId,
		force:    true, // Terms require force=true in WP REST API
	}
}

func (api *DeleteTag) Force() *DeleteTag {
	api.force = true
	return api
}

func (api *DeleteTag) Do() (tag Tag, err error) {
	endpoint := api.client.endpoint + api.endpoint + "/" + strconv.Itoa(api.tagId)
	var env deleteTagEnvelope
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
		return tag, &wpError
	}

	if err != nil {
		return
	}

	if env.Previous != nil && env.Previous.ID != 0 {
		return *env.Previous, nil
	}

	return env.Tag, nil
}
