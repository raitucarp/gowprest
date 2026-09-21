package gowprest

import (
	"encoding/json"
	"strings"
)

// Block category constants.
const (
	BlockCategoryText    = "text"
	BlockCategoryMedia   = "media"
	BlockCategoryDesign  = "design"
	BlockCategoryWidgets = "widgets"
	BlockCategoryTheme   = "theme"
	BlockCategoryEmbed   = "embed"
)

// BlockStyle describes a style variant for a block type.
type BlockStyle struct {
	Name      string `json:"name,omitempty"`
	Label     string `json:"label,omitempty"`
	IsDefault bool   `json:"is_default,omitempty"`
}

// BlockType represents a registered Gutenberg block type.
type BlockType struct {
	Name                string         `json:"name,omitempty"`
	Title               string         `json:"title,omitempty"`
	Description         string         `json:"description,omitempty"`
	Icon                any            `json:"icon,omitempty"`
	Category            string         `json:"category,omitempty"`
	Keywords            []string       `json:"keywords,omitempty"`
	Parent              any            `json:"parent,omitempty"`
	Ancestor            any            `json:"ancestor,omitempty"`
	AllowedBlocks       any            `json:"allowed_blocks,omitempty"`
	Attributes          map[string]any `json:"attributes,omitempty"`
	ProvidesContext     any            `json:"provides_context,omitempty"`
	UsesContext         []string       `json:"uses_context,omitempty"`
	Supports            map[string]any `json:"supports,omitempty"`
	Styles              []BlockStyle   `json:"styles,omitempty"`
	Textdomain          *string        `json:"textdomain,omitempty"`
	Example             any            `json:"example,omitempty"`
	Selectors           any            `json:"selectors,omitempty"`
	IsDynamic           bool           `json:"is_dynamic,omitempty"`
	APIVersion          int            `json:"api_version,omitempty"`
	EditorScriptHandles []string       `json:"editor_script_handles,omitempty"`
	ScriptHandles       []string       `json:"script_handles,omitempty"`
	ViewScriptHandles   []string       `json:"view_script_handles,omitempty"`
	ViewScriptModuleIDs []string       `json:"view_script_module_ids,omitempty"`
	EditorStyleHandles  []string       `json:"editor_style_handles,omitempty"`
	StyleHandles        []string       `json:"style_handles,omitempty"`
	ViewStyleHandles    []string       `json:"view_style_handles,omitempty"`
	Variations          []any          `json:"variations,omitempty"`
	BlockHooks          any            `json:"block_hooks,omitempty"`
	EditorScript        *string        `json:"editor_script,omitempty"`
	Script              *string        `json:"script,omitempty"`
	ViewScript          *string        `json:"view_script,omitempty"`
	EditorStyle         *string        `json:"editor_style,omitempty"`
	Style               any            `json:"style,omitempty"`
	Links               map[string]any `json:"_links,omitempty"`
}

// BlockTypes handles requests to the WordPress block types API.
type BlockTypes struct {
	client *RestClient
}

// BlockTypes returns a BlockTypes service instance.
func (c *RestClient) BlockTypes() *BlockTypes {
	return &BlockTypes{client: c}
}

// List returns a ListBlockTypes builder to query block types.
func (api *BlockTypes) List() *ListBlockTypes {
	return &ListBlockTypes{
		endpoint:  "/wp/v2/block-types",
		client:    api.client,
		arguments: make(map[string]string),
	}
}

// ByNamespace returns a ListBlockTypes builder filtered by namespace.
func (api *BlockTypes) ByNamespace(namespace string) *ListBlockTypes {
	return api.List().Namespace(namespace)
}

// Retrieve returns a RetrieveBlockType builder to get a specific block type by name.
// Accepts either a combined identifier ("core/paragraph") or separate namespace and name ("core", "paragraph").
func (api *BlockTypes) Retrieve(identifier string, name ...string) *RetrieveBlockType {
	path := "/wp/v2/block-types"
	if len(name) > 0 {
		path += "/" + strings.Trim(identifier, "/") + "/" + strings.Trim(name[0], "/")
	} else {
		path += "/" + strings.Trim(identifier, "/")
	}
	return &RetrieveBlockType{
		endpoint:  path,
		client:    api.client,
		arguments: make(map[string]string),
	}
}

// ListBlockTypes handles querying block type collections.
type ListBlockTypes struct {
	endpoint  string
	client    *RestClient
	arguments map[string]string
}

func (api *ListBlockTypes) Namespace(namespace string) *ListBlockTypes {
	api.arguments["namespace"] = namespace
	return api
}

func (api *ListBlockTypes) Context(ctx string) *ListBlockTypes {
	api.arguments["context"] = ctx
	return api
}

func (api *ListBlockTypes) ContextView() *ListBlockTypes {
	return api.Context("view")
}

func (api *ListBlockTypes) ContextEdit() *ListBlockTypes {
	return api.Context("edit")
}

func (api *ListBlockTypes) ContextEmbed() *ListBlockTypes {
	return api.Context("embed")
}

func (api *ListBlockTypes) Fields(fields ...string) *ListBlockTypes {
	api.arguments["_fields"] = strings.Join(fields, ",")
	return api
}

func (api *ListBlockTypes) Do() (blockTypes []BlockType, err error) {
	endpoint := api.client.endpoint + api.endpoint

	restyClient := api.client.httpClient.R()
	if api.client.auth.Username != "" && api.client.auth.Password != "" {
		restyClient.SetBasicAuth(api.client.auth.Username, api.client.auth.Password)
	}

	resp, err := restyClient.
		SetHeader("Accept", "application/json").
		SetResult(&blockTypes).
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

	return blockTypes, nil
}

// RetrieveBlockType handles querying a single block type.
type RetrieveBlockType struct {
	endpoint  string
	client    *RestClient
	arguments map[string]string
}

func (api *RetrieveBlockType) Context(ctx string) *RetrieveBlockType {
	api.arguments["context"] = ctx
	return api
}

func (api *RetrieveBlockType) ContextView() *RetrieveBlockType {
	return api.Context("view")
}

func (api *RetrieveBlockType) ContextEdit() *RetrieveBlockType {
	return api.Context("edit")
}

func (api *RetrieveBlockType) ContextEmbed() *RetrieveBlockType {
	return api.Context("embed")
}

func (api *RetrieveBlockType) Fields(fields ...string) *RetrieveBlockType {
	api.arguments["_fields"] = strings.Join(fields, ",")
	return api
}

func (api *RetrieveBlockType) Do() (blockType *BlockType, err error) {
	endpoint := api.client.endpoint + api.endpoint

	restyClient := api.client.httpClient.R()
	if api.client.auth.Username != "" && api.client.auth.Password != "" {
		restyClient.SetBasicAuth(api.client.auth.Username, api.client.auth.Password)
	}

	resp, err := restyClient.
		SetHeader("Accept", "application/json").
		SetResult(&blockType).
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

	return blockType, nil
}
