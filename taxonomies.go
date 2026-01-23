package gowprest

import (
	"encoding/json"
)

type TaxonomyCapabilities struct {
	ManageTerms string `json:"manage_terms"`
	EditTerms   string `json:"edit_terms"`
	DeleteTerms string `json:"delete_terms"`
	AssignTerms string `json:"assign_terms"`
}

type TaxonomyLabels struct {
	Name                    string `json:"name"`
	SingularName            string `json:"singular_name"`
	SearchItems             string `json:"search_items"`
	PopularItems            string `json:"popular_items"`
	AllItems                string `json:"all_items"`
	ParentItem              string `json:"parent_item"`
	ParentItemColon         string `json:"parent_item_colon"`
	EditItem                string `json:"edit_item"`
	ViewItem                string `json:"view_item"`
	UpdateItem              string `json:"update_item"`
	AddItem                 string `json:"add_item"`
	NewItemName             string `json:"new_item_name"`
	SeparateItemsWithCommas string `json:"separate_items_with_commas"`
	AddOrRemoveItems        string `json:"add_or_remove_items"`
	ChooseFromMostUsed      string `json:"choose_from_most_used"`
	NotFound                string `json:"not_found"`
	NoTerms                 string `json:"no_terms"`
	FilterByItem            string `json:"filter_by_item"`
	ItemsListNavigation     string `json:"items_list_navigation"`
	ItemsList               string `json:"items_list"`
	MostUsed                string `json:"most_used"`
	BackToItems             string `json:"back_to_items"`
	ItemLink                string `json:"item_link"`
	ItemLinkDescription     string `json:"item_link_description"`
	MenuName                string `json:"menu_name"`
	NameAdminBar            string `json:"name_admin_bar"`
	Archives                string `json:"archives"`
}

type TaxonomyVisibility struct {
	Public            bool `json:"public"`
	PubliclyQueryable bool `json:"publicly_queryable"`
	ShowUI            bool `json:"show_ui"`
	ShowInMenu        bool `json:"show_in_menu"`
	ShowInNavMenus    bool `json:"show_in_nav_menus"`
	ShowInQuickEdit   bool `json:"show_in_quick_edit"`
	ShowAdminColumn   bool `json:"show_admin_column"`
}

type Taxonomy struct {
	Capabilities  TaxonomyCapabilities `json:"capabilities,omitempty"`
	Description   string               `json:"description,omitempty"`
	Hierarchical  bool                 `json:"hierarchical,omitempty"`
	Labels        TaxonomyLabels       `json:"labels,omitempty"`
	Name          string               `json:"name,omitempty"`
	Slug          string               `json:"slug,omitempty"`
	ShowCloud     bool                 `json:"show_cloud,omitempty"`
	Types         []string             `json:"types,omitempty"`
	RestBase      string               `json:"rest_base,omitempty"`
	RestNamespace string               `json:"rest_namespace,omitempty"`
	Visibility    TaxonomyVisibility   `json:"visibility,omitempty"`
}

type Taxonomies struct {
	client *RestClient
}

func (c *RestClient) Taxonomies() *Taxonomies {
	return &Taxonomies{client: c}
}

type ListTaxonomies struct {
	endpoint  string
	client    *RestClient
	arguments map[string]string
}

func (api *Taxonomies) List() *ListTaxonomies {
	return &ListTaxonomies{
		endpoint:  "/wp/v2/taxonomies",
		client:    api.client,
		arguments: make(map[string]string),
	}
}

func (api *ListTaxonomies) ContextView() *ListTaxonomies {
	api.arguments["context"] = "view"
	return api
}

func (api *ListTaxonomies) ContextEdit() *ListTaxonomies {
	api.arguments["context"] = "edit"
	return api
}

func (api *ListTaxonomies) ContextEmbed() *ListTaxonomies {
	api.arguments["context"] = "embed"
	return api
}

func (api *ListTaxonomies) Type(postType string) *ListTaxonomies {
	api.arguments["type"] = postType
	return api
}

func (api *ListTaxonomies) Do() (taxonomies map[string]Taxonomy, err error) {
	endpoint := api.client.endpoint + api.endpoint

	restyClient := api.client.httpClient.R()
	if api.client.auth.Username != "" && api.client.auth.Password != "" && api.arguments["context"] == "edit" {
		restyClient.SetBasicAuth(api.client.auth.Username, api.client.auth.Password)
	}

	resp, err := restyClient.
		SetHeader("Accept", "application/json").
		SetResult(&taxonomies).
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

		return taxonomies, &wpError
	}

	return
}

type RetrieveTaxonomy struct {
	endpoint  string
	client    *RestClient
	arguments map[string]string
}

func (api *Taxonomies) Retrieve(taxonomy string) *RetrieveTaxonomy {
	return &RetrieveTaxonomy{
		endpoint:  "/wp/v2/taxonomies/" + taxonomy,
		client:    api.client,
		arguments: make(map[string]string),
	}
}

func (api *RetrieveTaxonomy) ContextView() *RetrieveTaxonomy {
	api.arguments["context"] = "view"
	return api
}

func (api *RetrieveTaxonomy) ContextEdit() *RetrieveTaxonomy {
	api.arguments["context"] = "edit"
	return api
}

func (api *RetrieveTaxonomy) ContextEmbed() *RetrieveTaxonomy {
	api.arguments["context"] = "embed"
	return api
}

func (api *RetrieveTaxonomy) Do() (taxonomy *Taxonomy, err error) {
	endpoint := api.client.endpoint + api.endpoint

	restyClient := api.client.httpClient.R()
	if api.client.auth.Username != "" && api.client.auth.Password != "" && api.arguments["context"] == "edit" {
		restyClient.SetBasicAuth(api.client.auth.Username, api.client.auth.Password)
	}

	resp, err := restyClient.
		SetHeader("Accept", "application/json").
		SetResult(&taxonomy).
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

		return taxonomy, &wpError
	}

	return
}
