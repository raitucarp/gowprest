package gowprest

import (
	"encoding/json"
	"io"
	"mime"
	"net/http"
	"os"
	"path/filepath"
	"strconv"
	"strings"
	"time"
)

// Meta represents metadata associated with a WordPress object, handling both {} and [] from PHP JSON.
type Meta map[string]any

func (m *Meta) UnmarshalJSON(data []byte) error {
	trimmed := strings.TrimSpace(string(data))
	if trimmed == "[]" || trimmed == "null" || trimmed == "" {
		*m = make(map[string]any)
		return nil
	}
	var res map[string]any
	if err := json.Unmarshal(data, &res); err != nil {
		*m = make(map[string]any)
		return nil
	}
	*m = res
	return nil
}

// MediaSizes describes size variants of a media item, handling empty arrays from PHP JSON.
type MediaSizes map[string]MediaSize

func (s *MediaSizes) UnmarshalJSON(data []byte) error {
	trimmed := strings.TrimSpace(string(data))
	if trimmed == "[]" || trimmed == "null" || trimmed == "" {
		*s = make(map[string]MediaSize)
		return nil
	}
	var res map[string]MediaSize
	if err := json.Unmarshal(data, &res); err != nil {
		*s = make(map[string]MediaSize)
		return nil
	}
	*s = res
	return nil
}

// ImageMeta describes image camera/exif metadata, handling empty arrays from PHP JSON.
type ImageMeta map[string]any

func (im *ImageMeta) UnmarshalJSON(data []byte) error {
	trimmed := strings.TrimSpace(string(data))
	if trimmed == "[]" || trimmed == "null" || trimmed == "" {
		*im = make(map[string]any)
		return nil
	}
	var res map[string]any
	if err := json.Unmarshal(data, &res); err != nil {
		*im = make(map[string]any)
		return nil
	}
	*im = res
	return nil
}

// MediaSize describes the dimensions and file details of an image size.
type MediaSize struct {
	File      string `json:"file,omitempty"`
	Width     int    `json:"width,omitempty"`
	Height    int    `json:"height,omitempty"`
	MimeType  string `json:"mime_type,omitempty"`
	SourceURL string `json:"source_url,omitempty"`
}

// MediaDetails contains metadata for an uploaded media attachment.
type MediaDetails struct {
	Width     int        `json:"width,omitempty"`
	Height    int        `json:"height,omitempty"`
	File      string     `json:"file,omitempty"`
	Filesize  int        `json:"filesize,omitempty"`
	Sizes     MediaSizes `json:"sizes,omitempty"`
	ImageMeta ImageMeta  `json:"image_meta,omitempty"`
}

// Media represents a WordPress media attachment.
type Media struct {
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
	PermalinkTemplate string           `json:"permalink_template,omitempty"`
	GeneratedSlug     string           `json:"generated_slug,omitempty"`
	Title             *Object          `json:"title,omitempty"`
	Author            int              `json:"author,omitempty"`
	CommentStatus     OpenClosedStatus `json:"comment_status,omitempty"`
	PingStatus        OpenClosedStatus `json:"ping_status,omitempty"`
	Meta              Meta             `json:"meta,omitempty"`
	Template          string           `json:"template,omitempty"`
	AltText           string           `json:"alt_text,omitempty"`
	Caption           *Object          `json:"caption,omitempty"`
	Description       *Object          `json:"description,omitempty"`
	MediaType         string           `json:"media_type,omitempty"`
	MimeType          string           `json:"mime_type,omitempty"`
	MediaDetails      MediaDetails     `json:"media_details,omitempty"`
	Post              *int             `json:"post,omitempty"`
	SourceURL         string           `json:"source_url,omitempty"`
	MissingImageSizes []string         `json:"missing_image_sizes,omitempty"`
	Links             map[string]any   `json:"_links,omitempty"`
	Embedded          map[string]any   `json:"_embedded,omitempty"`
}

// MediaData holds fields used to create or update a media attachment.
type MediaData struct {
	ID            int              `json:"id,omitempty"`
	Date          *Date            `json:"date,omitempty"`
	DateGMT       *Date            `json:"date_gmt,omitempty"`
	Slug          string           `json:"slug,omitempty"`
	Status        PostStatus       `json:"status,omitempty"`
	Title         string           `json:"title,omitempty"`
	Author        int              `json:"author,omitempty"`
	CommentStatus OpenClosedStatus `json:"comment_status,omitempty"`
	PingStatus    OpenClosedStatus `json:"ping_status,omitempty"`
	Meta          map[string]any   `json:"meta,omitempty"`
	Template      string           `json:"template,omitempty"`
	AltText       string           `json:"alt_text,omitempty"`
	Caption       string           `json:"caption,omitempty"`
	Description   string           `json:"description,omitempty"`
	Post          int              `json:"post,omitempty"`
}

// DeletedMedia represents the response returned when deleting a media item.
type DeletedMedia struct {
	Media
	Previous *Media `json:"previous,omitempty"`
	Deleted  bool   `json:"deleted"`
}

// MediaAPI provides methods to interact with WordPress media items.
type MediaAPI struct {
	client *RestClient
}

// Media returns a MediaAPI instance for the REST client.
func (c *RestClient) Media() *MediaAPI {
	return &MediaAPI{client: c}
}

// Medias is an alias for Media.
func (c *RestClient) Medias() *MediaAPI {
	return &MediaAPI{client: c}
}

// List returns a ListMedia builder to query media collections.
func (api *MediaAPI) List() *ListMedia {
	return &ListMedia{
		endpoint:  "/wp/v2/media",
		client:    api.client,
		arguments: make(map[string]string),
	}
}

// Retrieve returns a RetrieveMedia builder to get a specific media item.
func (api *MediaAPI) Retrieve(mediaID int) *RetrieveMedia {
	return &RetrieveMedia{
		endpoint:  "/wp/v2/media/" + strconv.Itoa(mediaID),
		client:    api.client,
		arguments: make(map[string]string),
	}
}

// Create returns a CreateMedia builder to upload and create a media item.
func (api *MediaAPI) Create(media ...MediaData) *CreateMedia {
	req := &CreateMedia{
		endpoint:    "/wp/v2/media",
		client:      api.client,
		queryParams: make(map[string]string),
	}
	if len(media) > 0 {
		req.media = media[0]
	}
	return req
}

// Update returns an UpdateMedia builder to modify a media item.
func (api *MediaAPI) Update(media ...MediaData) *UpdateMedia {
	endpoint := "/wp/v2/media"
	var m MediaData
	if len(media) > 0 {
		m = media[0]
		if m.ID != 0 {
			endpoint = "/wp/v2/media/" + strconv.Itoa(m.ID)
		}
	}
	return &UpdateMedia{
		endpoint: endpoint,
		client:   api.client,
		media:    m,
	}
}

// Delete returns a DeleteMedia builder to delete a media item.
func (api *MediaAPI) Delete(mediaID int) *DeleteMedia {
	return &DeleteMedia{
		endpoint: "/wp/v2/media/" + strconv.Itoa(mediaID),
		client:   api.client,
		mediaID:  mediaID,
		force:    true, // WordPress media requires force=true unless MEDIA_TRASH is enabled
	}
}

// ListMedia handles querying media collections.
type ListMedia struct {
	endpoint  string
	client    *RestClient
	arguments map[string]string
}

func (api *ListMedia) Context(ctx string) *ListMedia {
	api.arguments["context"] = ctx
	return api
}

func (api *ListMedia) ContextView() *ListMedia {
	api.arguments["context"] = "view"
	return api
}

func (api *ListMedia) ContextEdit() *ListMedia {
	api.arguments["context"] = "edit"
	return api
}

func (api *ListMedia) ContextEmbed() *ListMedia {
	api.arguments["context"] = "embed"
	return api
}

func (api *ListMedia) Page(page int) *ListMedia {
	api.arguments["page"] = strconv.Itoa(page)
	return api
}

func (api *ListMedia) PerPage(perPage int) *ListMedia {
	api.arguments["per_page"] = strconv.Itoa(perPage)
	return api
}

func (api *ListMedia) Search(query string) *ListMedia {
	api.arguments["search"] = query
	return api
}

func (api *ListMedia) SearchColumns(columns ...string) *ListMedia {
	api.arguments["search_columns"] = strings.Join(columns, ",")
	return api
}

func (api *ListMedia) After(after time.Time) *ListMedia {
	api.arguments["after"] = after.Format(time.RFC3339)
	return api
}

func (api *ListMedia) Before(before time.Time) *ListMedia {
	api.arguments["before"] = before.Format(time.RFC3339)
	return api
}

func (api *ListMedia) ModifiedAfter(after time.Time) *ListMedia {
	api.arguments["modified_after"] = after.Format(time.RFC3339)
	return api
}

func (api *ListMedia) ModifiedBefore(before time.Time) *ListMedia {
	api.arguments["modified_before"] = before.Format(time.RFC3339)
	return api
}

func (api *ListMedia) Author(authorIDs ...int) *ListMedia {
	authors := make([]string, 0, len(authorIDs))
	for _, id := range authorIDs {
		authors = append(authors, strconv.Itoa(id))
	}
	api.arguments["author"] = strings.Join(authors, ",")
	return api
}

func (api *ListMedia) AuthorExclude(authorIDs ...int) *ListMedia {
	authors := make([]string, 0, len(authorIDs))
	for _, id := range authorIDs {
		authors = append(authors, strconv.Itoa(id))
	}
	api.arguments["author_exclude"] = strings.Join(authors, ",")
	return api
}

func (api *ListMedia) Exclude(excludeIDs ...int) *ListMedia {
	excludes := make([]string, 0, len(excludeIDs))
	for _, id := range excludeIDs {
		excludes = append(excludes, strconv.Itoa(id))
	}
	api.arguments["exclude"] = strings.Join(excludes, ",")
	return api
}

func (api *ListMedia) Include(includeIDs ...int) *ListMedia {
	includes := make([]string, 0, len(includeIDs))
	for _, id := range includeIDs {
		includes = append(includes, strconv.Itoa(id))
	}
	api.arguments["include"] = strings.Join(includes, ",")
	return api
}

func (api *ListMedia) Offset(offset int) *ListMedia {
	api.arguments["offset"] = strconv.Itoa(offset)
	return api
}

func (api *ListMedia) Order(order string) *ListMedia {
	api.arguments["order"] = order
	return api
}

func (api *ListMedia) OrderAsc() *ListMedia {
	api.arguments["order"] = "asc"
	return api
}

func (api *ListMedia) OrderDesc() *ListMedia {
	api.arguments["order"] = "desc"
	return api
}

func (api *ListMedia) OrderBy(orderBy string) *ListMedia {
	api.arguments["orderby"] = orderBy
	return api
}

func (api *ListMedia) OrderByDate() *ListMedia {
	api.arguments["orderby"] = "date"
	return api
}

func (api *ListMedia) OrderByID() *ListMedia {
	api.arguments["orderby"] = "id"
	return api
}

func (api *ListMedia) OrderByAuthor() *ListMedia {
	api.arguments["orderby"] = "author"
	return api
}

func (api *ListMedia) OrderByModified() *ListMedia {
	api.arguments["orderby"] = "modified"
	return api
}

func (api *ListMedia) OrderByParent() *ListMedia {
	api.arguments["orderby"] = "parent"
	return api
}

func (api *ListMedia) OrderByTitle() *ListMedia {
	api.arguments["orderby"] = "title"
	return api
}

func (api *ListMedia) OrderBySlug() *ListMedia {
	api.arguments["orderby"] = "slug"
	return api
}

func (api *ListMedia) OrderByRelevance() *ListMedia {
	api.arguments["orderby"] = "relevance"
	return api
}

func (api *ListMedia) OrderByInclude() *ListMedia {
	api.arguments["orderby"] = "include"
	return api
}

func (api *ListMedia) Parent(parentIDs ...int) *ListMedia {
	parents := make([]string, 0, len(parentIDs))
	for _, id := range parentIDs {
		parents = append(parents, strconv.Itoa(id))
	}
	api.arguments["parent"] = strings.Join(parents, ",")
	return api
}

func (api *ListMedia) ParentExclude(parentIDs ...int) *ListMedia {
	parents := make([]string, 0, len(parentIDs))
	for _, id := range parentIDs {
		parents = append(parents, strconv.Itoa(id))
	}
	api.arguments["parent_exclude"] = strings.Join(parents, ",")
	return api
}

func (api *ListMedia) Slug(slugs ...string) *ListMedia {
	api.arguments["slug"] = strings.Join(slugs, ",")
	return api
}

func (api *ListMedia) Status(statuses ...string) *ListMedia {
	api.arguments["status"] = strings.Join(statuses, ",")
	return api
}

func (api *ListMedia) MediaType(mediaType string) *ListMedia {
	api.arguments["media_type"] = mediaType
	return api
}

func (api *ListMedia) MediaTypeImage() *ListMedia {
	api.arguments["media_type"] = "image"
	return api
}

func (api *ListMedia) MediaTypeVideo() *ListMedia {
	api.arguments["media_type"] = "video"
	return api
}

func (api *ListMedia) MediaTypeAudio() *ListMedia {
	api.arguments["media_type"] = "audio"
	return api
}

func (api *ListMedia) MediaTypeText() *ListMedia {
	api.arguments["media_type"] = "text"
	return api
}

func (api *ListMedia) MediaTypeApplication() *ListMedia {
	api.arguments["media_type"] = "application"
	return api
}

func (api *ListMedia) MimeType(mimeType string) *ListMedia {
	api.arguments["mime_type"] = mimeType
	return api
}

func (api *ListMedia) Embed() *ListMedia {
	api.arguments["_embed"] = "true"
	return api
}

func (api *ListMedia) Fields(fields ...string) *ListMedia {
	api.arguments["_fields"] = strings.Join(fields, ",")
	return api
}

func (api *ListMedia) Do() (mediaList []Media, err error) {
	restyClient := api.client.httpClient.R()
	if api.client.auth.Username != "" && api.client.auth.Password != "" {
		restyClient.SetBasicAuth(api.client.auth.Username, api.client.auth.Password)
	}

	resp, err := restyClient.
		SetHeader("Accept", "application/json").
		SetResult(&mediaList).
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
		return mediaList, &wpError
	}

	return
}

// RetrieveMedia handles retrieving a single media item.
type RetrieveMedia struct {
	endpoint  string
	client    *RestClient
	arguments map[string]string
}

func (api *RetrieveMedia) Context(ctx string) *RetrieveMedia {
	api.arguments["context"] = ctx
	return api
}

func (api *RetrieveMedia) ContextView() *RetrieveMedia {
	api.arguments["context"] = "view"
	return api
}

func (api *RetrieveMedia) ContextEdit() *RetrieveMedia {
	api.arguments["context"] = "edit"
	return api
}

func (api *RetrieveMedia) ContextEmbed() *RetrieveMedia {
	api.arguments["context"] = "embed"
	return api
}

func (api *RetrieveMedia) Embed() *RetrieveMedia {
	api.arguments["_embed"] = "true"
	return api
}

func (api *RetrieveMedia) Fields(fields ...string) *RetrieveMedia {
	api.arguments["_fields"] = strings.Join(fields, ",")
	return api
}

func (api *RetrieveMedia) Do() (media *Media, err error) {
	endpoint := api.client.endpoint + api.endpoint

	restyClient := api.client.httpClient.R()
	if api.client.auth.Username != "" && api.client.auth.Password != "" {
		restyClient.SetBasicAuth(api.client.auth.Username, api.client.auth.Password)
	}

	resp, err := restyClient.
		SetHeader("Accept", "application/json").
		SetResult(&media).
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
		return media, &wpError
	}

	return
}

// CreateMedia handles uploading and creating a media attachment.
type CreateMedia struct {
	endpoint    string
	client      *RestClient
	media       MediaData
	filePath    string
	fileName    string
	fileData    []byte
	fileReader  io.Reader
	mimeType    string
	queryParams map[string]string
}

// File sets the source file path on the local filesystem.
func (api *CreateMedia) File(filePath string) *CreateMedia {
	api.filePath = filePath
	if api.fileName == "" {
		api.fileName = filepath.Base(filePath)
	}
	return api
}

// FileBytes sets the file upload from in-memory byte slice with a filename.
func (api *CreateMedia) FileBytes(fileName string, data []byte, mimeType ...string) *CreateMedia {
	api.fileName = fileName
	api.fileData = data
	if len(mimeType) > 0 {
		api.mimeType = mimeType[0]
	}
	return api
}

// FileReader sets the file upload from an io.Reader.
func (api *CreateMedia) FileReader(fileName string, r io.Reader, mimeType ...string) *CreateMedia {
	api.fileName = fileName
	api.fileReader = r
	if len(mimeType) > 0 {
		api.mimeType = mimeType[0]
	}
	return api
}

// FileName sets the filename for the attachment.
func (api *CreateMedia) FileName(fileName string) *CreateMedia {
	api.fileName = fileName
	return api
}

// MimeType specifies the MIME type of the uploaded file.
func (api *CreateMedia) MimeType(mimeType string) *CreateMedia {
	api.mimeType = mimeType
	return api
}

func (api *CreateMedia) Title(title string) *CreateMedia {
	api.media.Title = title
	return api
}

func (api *CreateMedia) AltText(altText string) *CreateMedia {
	api.media.AltText = altText
	return api
}

func (api *CreateMedia) Caption(caption string) *CreateMedia {
	api.media.Caption = caption
	return api
}

func (api *CreateMedia) Description(description string) *CreateMedia {
	api.media.Description = description
	return api
}

func (api *CreateMedia) Post(postID int) *CreateMedia {
	api.media.Post = postID
	return api
}

func (api *CreateMedia) Slug(slug string) *CreateMedia {
	api.media.Slug = slug
	return api
}

func (api *CreateMedia) Status(status PostStatus) *CreateMedia {
	api.media.Status = status
	return api
}

func (api *CreateMedia) Author(authorID int) *CreateMedia {
	api.media.Author = authorID
	return api
}

func (api *CreateMedia) CommentStatus(status OpenClosedStatus) *CreateMedia {
	api.media.CommentStatus = status
	return api
}

func (api *CreateMedia) PingStatus(status OpenClosedStatus) *CreateMedia {
	api.media.PingStatus = status
	return api
}

func (api *CreateMedia) Template(template string) *CreateMedia {
	api.media.Template = template
	return api
}

func (api *CreateMedia) Meta(meta map[string]any) *CreateMedia {
	api.media.Meta = meta
	return api
}

func (api *CreateMedia) Date(t time.Time) *CreateMedia {
	api.media.Date = &Date{Time: t}
	return api
}

func (api *CreateMedia) DateGMT(t time.Time) *CreateMedia {
	api.media.DateGMT = &Date{Time: t}
	return api
}

func (api *CreateMedia) Do() (media Media, err error) {
	restyClient := api.client.httpClient.R()
	if api.client.auth.Username != "" && api.client.auth.Password != "" {
		restyClient.SetBasicAuth(api.client.auth.Username, api.client.auth.Password)
	}

	var data []byte
	if api.fileData != nil {
		data = api.fileData
	} else if api.filePath != "" {
		data, err = os.ReadFile(api.filePath)
		if err != nil {
			return
		}
		if api.fileName == "" {
			api.fileName = filepath.Base(api.filePath)
		}
	} else if api.fileReader != nil {
		data, err = io.ReadAll(api.fileReader)
		if err != nil {
			return
		}
	}

	if api.fileName == "" {
		api.fileName = "upload.bin"
	}

	detectedMime := api.mimeType
	if detectedMime == "" {
		detectedMime = mime.TypeByExtension(filepath.Ext(api.fileName))
		if detectedMime == "" && len(data) > 0 {
			detectedMime = http.DetectContentType(data)
		}
	}
	if detectedMime == "" {
		detectedMime = "application/octet-stream"
	}

	// Build query params for metadata
	queryParams := make(map[string]string)
	for k, v := range api.queryParams {
		queryParams[k] = v
	}
	if api.media.Title != "" {
		queryParams["title"] = api.media.Title
	}
	if api.media.AltText != "" {
		queryParams["alt_text"] = api.media.AltText
	}
	if api.media.Caption != "" {
		queryParams["caption"] = api.media.Caption
	}
	if api.media.Description != "" {
		queryParams["description"] = api.media.Description
	}
	if api.media.Post != 0 {
		queryParams["post"] = strconv.Itoa(api.media.Post)
	}
	if api.media.Slug != "" {
		queryParams["slug"] = api.media.Slug
	}
	if api.media.Status != "" {
		queryParams["status"] = string(api.media.Status)
	}
	if api.media.Author != 0 {
		queryParams["author"] = strconv.Itoa(api.media.Author)
	}
	if api.media.CommentStatus != "" {
		queryParams["comment_status"] = string(api.media.CommentStatus)
	}
	if api.media.PingStatus != "" {
		queryParams["ping_status"] = string(api.media.PingStatus)
	}
	if api.media.Date != nil {
		queryParams["date"] = api.media.Date.Time.Format(time.RFC3339)
	}
	if api.media.DateGMT != nil {
		queryParams["date_gmt"] = api.media.DateGMT.Time.Format(time.RFC3339)
	}

	resp, err := restyClient.
		SetHeader("Content-Disposition", `attachment; filename="`+api.fileName+`"`).
		SetHeader("Content-Type", detectedMime).
		SetBody(data).
		SetQueryParams(queryParams).
		SetResult(&media).
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
		return media, &wpError
	}

	return
}

// UpdateMedia handles updating an existing media item's metadata.
type UpdateMedia struct {
	endpoint string
	client   *RestClient
	media    MediaData
}

func (api *UpdateMedia) ID(mediaID int) *UpdateMedia {
	api.media.ID = mediaID
	api.endpoint = "/wp/v2/media/" + strconv.Itoa(mediaID)
	return api
}

func (api *UpdateMedia) Title(title string) *UpdateMedia {
	api.media.Title = title
	return api
}

func (api *UpdateMedia) AltText(altText string) *UpdateMedia {
	api.media.AltText = altText
	return api
}

func (api *UpdateMedia) Caption(caption string) *UpdateMedia {
	api.media.Caption = caption
	return api
}

func (api *UpdateMedia) Description(description string) *UpdateMedia {
	api.media.Description = description
	return api
}

func (api *UpdateMedia) Post(postID int) *UpdateMedia {
	api.media.Post = postID
	return api
}

func (api *UpdateMedia) Slug(slug string) *UpdateMedia {
	api.media.Slug = slug
	return api
}

func (api *UpdateMedia) Status(status PostStatus) *UpdateMedia {
	api.media.Status = status
	return api
}

func (api *UpdateMedia) Author(authorID int) *UpdateMedia {
	api.media.Author = authorID
	return api
}

func (api *UpdateMedia) CommentStatus(status OpenClosedStatus) *UpdateMedia {
	api.media.CommentStatus = status
	return api
}

func (api *UpdateMedia) PingStatus(status OpenClosedStatus) *UpdateMedia {
	api.media.PingStatus = status
	return api
}

func (api *UpdateMedia) Template(template string) *UpdateMedia {
	api.media.Template = template
	return api
}

func (api *UpdateMedia) Meta(meta map[string]any) *UpdateMedia {
	api.media.Meta = meta
	return api
}

func (api *UpdateMedia) Date(t time.Time) *UpdateMedia {
	api.media.Date = &Date{Time: t}
	return api
}

func (api *UpdateMedia) DateGMT(t time.Time) *UpdateMedia {
	api.media.DateGMT = &Date{Time: t}
	return api
}

func (api *UpdateMedia) Do() (media Media, err error) {
	if api.media.ID != 0 && (api.endpoint == "/wp/v2/media" || !strings.Contains(api.endpoint, strconv.Itoa(api.media.ID))) {
		api.endpoint = "/wp/v2/media/" + strconv.Itoa(api.media.ID)
	}

	restyClient := api.client.httpClient.R()
	if api.client.auth.Username != "" && api.client.auth.Password != "" {
		restyClient.SetBasicAuth(api.client.auth.Username, api.client.auth.Password)
	}

	resp, err := restyClient.
		SetHeader("Content-Type", "application/json").
		SetResult(&media).
		SetBody(api.media).
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
		return media, &wpError
	}

	return
}

// DeleteMedia handles deleting a media attachment.
type DeleteMedia struct {
	endpoint string
	client   *RestClient
	mediaID  int
	force    bool
}

func (api *DeleteMedia) Force(force ...bool) *DeleteMedia {
	if len(force) > 0 {
		api.force = force[0]
	} else {
		api.force = true
	}
	return api
}

func (api *DeleteMedia) Do() (deletedMedia DeletedMedia, err error) {
	endpoint := api.client.endpoint + api.endpoint
	restyClient := api.client.httpClient.R()
	if api.client.auth.Username != "" && api.client.auth.Password != "" {
		restyClient.SetBasicAuth(api.client.auth.Username, api.client.auth.Password)
	}

	resp, err := restyClient.
		SetHeader("Content-Type", "application/json").
		SetResult(&deletedMedia).
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
		return deletedMedia, &wpError
	}

	if deletedMedia.Previous != nil && deletedMedia.Previous.ID != 0 {
		deletedMedia.Media = *deletedMedia.Previous
	}

	return
}
