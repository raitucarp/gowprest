package gowprest

import (
	"encoding/json"
	"strings"
)

// Standard WordPress core post type constants.
const (
	PostTypePost           = "post"
	PostTypePage           = "page"
	PostTypeAttachment     = "attachment"
	PostTypeNavMenuItem    = "nav_menu_item"
	PostTypeWPBlock        = "wp_block"
	PostTypeWPTemplate     = "wp_template"
	PostTypeWPTemplatePart = "wp_template_part"
	PostTypeWPNavigation   = "wp_navigation"
	PostTypeWPFontFamily   = "wp_font_family"
	PostTypeWPFontFace     = "wp_font_face"
)

// PostTypeCapabilities describes the permissions associated with a post type.
type PostTypeCapabilities struct {
	EditPost             string `json:"edit_post,omitempty"`
	ReadPost             string `json:"read_post,omitempty"`
	DeletePost           string `json:"delete_post,omitempty"`
	EditPosts            string `json:"edit_posts,omitempty"`
	EditOthersPosts      string `json:"edit_others_posts,omitempty"`
	DeletePosts          string `json:"delete_posts,omitempty"`
	PublishPosts         string `json:"publish_posts,omitempty"`
	ReadPrivatePosts     string `json:"read_private_posts,omitempty"`
	Read                 string `json:"read,omitempty"`
	DeletePrivatePosts   string `json:"delete_private_posts,omitempty"`
	DeletePublishedPosts string `json:"delete_published_posts,omitempty"`
	DeleteOthersPosts    string `json:"delete_others_posts,omitempty"`
	EditPrivatePosts     string `json:"edit_private_posts,omitempty"`
	EditPublishedPosts   string `json:"edit_published_posts,omitempty"`
	CreatePosts          string `json:"create_posts,omitempty"`
}

// PostTypeLabels describes human-readable labels for display in admin and UI.
type PostTypeLabels struct {
	Name                   string  `json:"name,omitempty"`
	SingularName           string  `json:"singular_name,omitempty"`
	AddNew                 string  `json:"add_new,omitempty"`
	AddNewItem             string  `json:"add_new_item,omitempty"`
	EditItem               string  `json:"edit_item,omitempty"`
	NewItem                string  `json:"new_item,omitempty"`
	ViewItem               string  `json:"view_item,omitempty"`
	ViewItems              string  `json:"view_items,omitempty"`
	SearchItems            string  `json:"search_items,omitempty"`
	NotFound               string  `json:"not_found,omitempty"`
	NotFoundInTrash        string  `json:"not_found_in_trash,omitempty"`
	ParentItemColon        *string `json:"parent_item_colon,omitempty"`
	AllItems               string  `json:"all_items,omitempty"`
	Archives               string  `json:"archives,omitempty"`
	Attributes             string  `json:"attributes,omitempty"`
	InsertIntoItem         string  `json:"insert_into_item,omitempty"`
	UploadedToThisItem     string  `json:"uploaded_to_this_item,omitempty"`
	FeaturedImage          string  `json:"featured_image,omitempty"`
	SetFeaturedImage       string  `json:"set_featured_image,omitempty"`
	RemoveFeaturedImage    string  `json:"remove_featured_image,omitempty"`
	UseFeaturedImage       string  `json:"use_featured_image,omitempty"`
	FilterItemsList        string  `json:"filter_items_list,omitempty"`
	FilterByDate           string  `json:"filter_by_date,omitempty"`
	ItemsListNavigation    string  `json:"items_list_navigation,omitempty"`
	ItemsList              string  `json:"items_list,omitempty"`
	ItemPublished          string  `json:"item_published,omitempty"`
	ItemPublishedPrivately string  `json:"item_published_privately,omitempty"`
	ItemRevertedToDraft    string  `json:"item_reverted_to_draft,omitempty"`
	ItemTrashed            string  `json:"item_trashed,omitempty"`
	ItemScheduled          string  `json:"item_scheduled,omitempty"`
	ItemUpdated            string  `json:"item_updated,omitempty"`
	ItemLink               string  `json:"item_link,omitempty"`
	ItemLinkDescription    string  `json:"item_link_description,omitempty"`
	MenuName               string  `json:"menu_name,omitempty"`
	NameAdminBar           string  `json:"name_admin_bar,omitempty"`
}

// PostTypeVisibility describes visibility settings for a post type.
type PostTypeVisibility struct {
	ShowUI         bool `json:"show_ui,omitempty"`
	ShowInNavMenus bool `json:"show_in_nav_menus,omitempty"`
}

// PostType represents a WordPress post type.
type PostType struct {
	Capabilities  PostTypeCapabilities `json:"capabilities,omitempty"`
	Description   string               `json:"description,omitempty"`
	Hierarchical  bool                 `json:"hierarchical,omitempty"`
	Viewable      bool                 `json:"viewable,omitempty"`
	Labels        PostTypeLabels       `json:"labels,omitempty"`
	Name          string               `json:"name,omitempty"`
	Slug          string               `json:"slug,omitempty"`
	Supports      map[string]any       `json:"supports,omitempty"`
	HasArchive    any                  `json:"has_archive,omitempty"`
	Taxonomies    []string             `json:"taxonomies,omitempty"`
	RestBase      string               `json:"rest_base,omitempty"`
	RestNamespace string               `json:"rest_namespace,omitempty"`
	Visibility    PostTypeVisibility   `json:"visibility,omitempty"`
	Icon          *string              `json:"icon,omitempty"`
	Template      []any                `json:"template,omitempty"`
	TemplateLock  any                  `json:"template_lock,omitempty"`
	Links         map[string]any       `json:"_links,omitempty"`
	Embedded      map[string]any       `json:"_embedded,omitempty"`
}

// PostTypes handles requests to the WordPress post types API.
type PostTypes struct {
	client *RestClient
}

// PostTypes returns a PostTypes service instance.
func (c *RestClient) PostTypes() *PostTypes {
	return &PostTypes{client: c}
}

// Types is an alias for PostTypes.
func (c *RestClient) Types() *PostTypes {
	return &PostTypes{client: c}
}

// List returns a ListPostTypes builder to query post types.
func (api *PostTypes) List() *ListPostTypes {
	return &ListPostTypes{
		endpoint:  "/wp/v2/types",
		client:    api.client,
		arguments: make(map[string]string),
	}
}

// Retrieve returns a RetrievePostType builder to get a specific post type by slug.
func (api *PostTypes) Retrieve(postType string) *RetrievePostType {
	return &RetrievePostType{
		endpoint:  "/wp/v2/types/" + postType,
		client:    api.client,
		arguments: make(map[string]string),
	}
}

// ListPostTypes handles querying post type collections.
type ListPostTypes struct {
	endpoint  string
	client    *RestClient
	arguments map[string]string
}

func (api *ListPostTypes) Context(ctx string) *ListPostTypes {
	api.arguments["context"] = ctx
	return api
}

func (api *ListPostTypes) ContextView() *ListPostTypes {
	api.arguments["context"] = "view"
	return api
}

func (api *ListPostTypes) ContextEdit() *ListPostTypes {
	api.arguments["context"] = "edit"
	return api
}

func (api *ListPostTypes) ContextEmbed() *ListPostTypes {
	api.arguments["context"] = "embed"
	return api
}

func (api *ListPostTypes) Embed() *ListPostTypes {
	api.arguments["_embed"] = "true"
	return api
}

func (api *ListPostTypes) Fields(fields ...string) *ListPostTypes {
	api.arguments["_fields"] = strings.Join(fields, ",")
	return api
}

func (api *ListPostTypes) Do() (types map[string]PostType, err error) {
	endpoint := api.client.endpoint + api.endpoint

	restyClient := api.client.httpClient.R()
	if api.client.auth.Username != "" && api.client.auth.Password != "" {
		restyClient.SetBasicAuth(api.client.auth.Username, api.client.auth.Password)
	}

	var raw json.RawMessage
	resp, err := restyClient.
		SetHeader("Accept", "application/json").
		SetResult(&raw).
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

	body := strings.TrimSpace(string(raw))
	if strings.HasPrefix(body, "[") {
		return make(map[string]PostType), nil
	}

	err = json.Unmarshal(raw, &types)
	if err != nil {
		return nil, err
	}
	return types, nil
}

// RetrievePostType handles querying a single post type.
type RetrievePostType struct {
	endpoint  string
	client    *RestClient
	arguments map[string]string
}

func (api *RetrievePostType) Context(ctx string) *RetrievePostType {
	api.arguments["context"] = ctx
	return api
}

func (api *RetrievePostType) ContextView() *RetrievePostType {
	api.arguments["context"] = "view"
	return api
}

func (api *RetrievePostType) ContextEdit() *RetrievePostType {
	api.arguments["context"] = "edit"
	return api
}

func (api *RetrievePostType) ContextEmbed() *RetrievePostType {
	api.arguments["context"] = "embed"
	return api
}

func (api *RetrievePostType) Embed() *RetrievePostType {
	api.arguments["_embed"] = "true"
	return api
}

func (api *RetrievePostType) Fields(fields ...string) *RetrievePostType {
	api.arguments["_fields"] = strings.Join(fields, ",")
	return api
}

func (api *RetrievePostType) Do() (postType *PostType, err error) {
	endpoint := api.client.endpoint + api.endpoint

	restyClient := api.client.httpClient.R()
	if api.client.auth.Username != "" && api.client.auth.Password != "" {
		restyClient.SetBasicAuth(api.client.auth.Username, api.client.auth.Password)
	}

	resp, err := restyClient.
		SetHeader("Accept", "application/json").
		SetResult(&postType).
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

	return postType, nil
}
