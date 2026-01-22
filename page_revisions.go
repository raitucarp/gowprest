package gowprest

import (
	"encoding/json"
	"strconv"
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
		endpoint:  "/wp/v2/pages/" + strconv.Itoa(api.parentID) + "/revisions/" + strconv.Itoa(revisionID),
		client:    api.client,
		arguments: make(map[string]string),
	}
}

// Delete returns a DeletePageRevision struct to delete a specific revision.
func (api *PageRevisions) Delete(revisionID int) *DeletePageRevision {
	return &DeletePageRevision{
		endpoint: "/wp/v2/pages/" + strconv.Itoa(api.parentID) + "/revisions/" + strconv.Itoa(revisionID),
		client:   api.client,
	}
}

// ListPageRevisions handles listing revisions.
type ListPageRevisions struct {
	endpoint  string
	client    *RestClient
	arguments map[string]string
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

func (api *ListPageRevisions) Offset(offset int) *ListPageRevisions {
	api.arguments["offset"] = strconv.Itoa(offset)
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

func (api *ListPageRevisions) Do() (revisions []Revision, err error) {
	_, err = api.client.httpClient.R().
		SetHeader("Accept", "application/json").
		SetResult(&revisions).
		SetQueryParams(api.arguments).
		Get(api.client.endpoint + api.endpoint)

	return
}

// RetrievePageRevision handles retrieving a specific revision.
type RetrievePageRevision struct {
	endpoint  string
	client    *RestClient
	arguments map[string]string
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

func (api *RetrievePageRevision) Do() (revision *Revision, err error) {
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

// DeletePageRevision handles deleting a specific revision.
type DeletePageRevision struct {
	endpoint string
	client   *RestClient
	force    bool
}

func (api *DeletePageRevision) Force() *DeletePageRevision {
	api.force = true
	return api
}

func (api *DeletePageRevision) Do() (revision Revision, err error) {
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
