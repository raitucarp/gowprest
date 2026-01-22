package gowprest

import (
	"encoding/json"
	"strconv"
)

// Revision represents a WordPress post revision.
type Revision struct {
	Author      int     `json:"author,omitempty"`
	Date        *Date   `json:"date,omitempty"`
	DateGMT     *Date   `json:"date_gmt,omitempty"`
	GUID        *Object `json:"guid,omitempty"`
	ID          int     `json:"id,omitempty"`
	Modified    *Date   `json:"modified,omitempty"`
	ModifiedGMT *Date   `json:"modified_gmt,omitempty"`
	Parent      int     `json:"parent,omitempty"`
	Slug        string  `json:"slug,omitempty"`
	Title       *Object `json:"title,omitempty"`
	Content     *Object `json:"content,omitempty"`
	Excerpt     *Object `json:"excerpt,omitempty"`
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
		endpoint:  "/wp/v2/posts/" + strconv.Itoa(api.parentID) + "/revisions/" + strconv.Itoa(revisionID),
		client:    api.client,
		arguments: make(map[string]string),
	}
}

// Delete returns a DeletePostRevision struct to delete a specific revision.
func (api *PostRevisions) Delete(revisionID int) *DeletePostRevision {
	return &DeletePostRevision{
		endpoint: "/wp/v2/posts/" + strconv.Itoa(api.parentID) + "/revisions/" + strconv.Itoa(revisionID),
		client:   api.client,
	}
}

// ListPostRevisions handles listing revisions.
type ListPostRevisions struct {
	endpoint  string
	client    *RestClient
	arguments map[string]string
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

func (api *ListPostRevisions) Offset(offset int) *ListPostRevisions {
	api.arguments["offset"] = strconv.Itoa(offset)
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

func (api *ListPostRevisions) Do() (revisions []Revision, err error) {
	_, err = api.client.httpClient.R().
		SetHeader("Accept", "application/json").
		SetResult(&revisions).
		SetQueryParams(api.arguments).
		Get(api.client.endpoint + api.endpoint)

	return
}

// RetrievePostRevision handles retrieving a specific revision.
type RetrievePostRevision struct {
	endpoint  string
	client    *RestClient
	arguments map[string]string
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

func (api *RetrievePostRevision) Do() (revision *Revision, err error) {
	endpoint := api.client.endpoint + api.endpoint

	restyClient := api.client.httpClient.R()
	if api.client.auth.Username != "" && api.client.auth.Password != "" && api.arguments["context"] == "edit" {
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

// DeletePostRevision handles deleting a specific revision.
type DeletePostRevision struct {
	endpoint string
	client   *RestClient
	force    bool
}

func (api *DeletePostRevision) Force() *DeletePostRevision {
	api.force = true
	return api
}

func (api *DeletePostRevision) Do() (revision Revision, err error) {
	resp, err := api.client.httpClient.R().
		SetHeader("Content-Type", "application/json").
		SetBasicAuth(api.client.auth.Username, api.client.auth.Password).
		SetResult(&revision).
		SetQueryParam("force", strconv.FormatBool(api.force)).
		Delete(api.client.endpoint + api.endpoint)

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
