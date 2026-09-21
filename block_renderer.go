package gowprest

import (
	"encoding/json"
	"strconv"
	"strings"
)

// RenderedBlock represents the server-rendered HTML output for a block.
type RenderedBlock struct {
	Rendered string `json:"rendered"`
}

// BlockRenderer handles requests to the WordPress block renderer API (/wp/v2/block-renderer/<name>).
type BlockRenderer struct {
	client *RestClient
}

// BlockRenderer returns a BlockRenderer service instance.
func (c *RestClient) BlockRenderer() *BlockRenderer {
	return &BlockRenderer{client: c}
}

// RenderedBlocks is a convenience alias for BlockRenderer.
func (c *RestClient) RenderedBlocks() *BlockRenderer {
	return c.BlockRenderer()
}

// Render returns a RenderBlock builder for the specified block name.
// Accepts either combined identifier ("core/calendar") or separate namespace and name ("core", "calendar").
func (api *BlockRenderer) Render(name string, blockName ...string) *RenderBlock {
	path := "/wp/v2/block-renderer"
	if len(blockName) > 0 {
		path += "/" + strings.Trim(name, "/") + "/" + strings.Trim(blockName[0], "/")
	} else {
		path += "/" + strings.Trim(name, "/")
	}
	builder := &RenderBlock{
		endpoint:   path,
		client:     api.client,
		attributes: make(map[string]any),
		context:    "edit", // WP block renderer requires context=edit
	}
	return builder
}

// Get returns a RenderBlock builder for GET requests.
func (api *BlockRenderer) Get(name string, blockName ...string) *RenderBlock {
	return api.Render(name, blockName...)
}

// Create returns a RenderBlock builder for POST requests (matching WP handbook "Create a Rendered Block").
func (api *BlockRenderer) Create(name string, blockName ...string) *RenderBlock {
	return api.Render(name, blockName...)
}

// RenderBlock handles rendering dynamic blocks.
type RenderBlock struct {
	endpoint   string
	client     *RestClient
	attributes map[string]any
	postID     int
	context    string
}

// Attributes sets the entire attributes map for the block.
func (api *RenderBlock) Attributes(attrs map[string]any) *RenderBlock {
	api.attributes = attrs
	return api
}

// SetAttribute sets a single attribute for the block.
func (api *RenderBlock) SetAttribute(key string, val any) *RenderBlock {
	if api.attributes == nil {
		api.attributes = make(map[string]any)
	}
	api.attributes[key] = val
	return api
}

// PostID sets the ID of the post context.
func (api *RenderBlock) PostID(postID int) *RenderBlock {
	api.postID = postID
	return api
}

// Context sets the context parameter (default is "edit").
func (api *RenderBlock) Context(ctx string) *RenderBlock {
	api.context = ctx
	return api
}

// ContextEdit sets the context parameter to "edit".
func (api *RenderBlock) ContextEdit() *RenderBlock {
	return api.Context("edit")
}

// Do executes a POST request to render the block with the configured attributes and post ID.
func (api *RenderBlock) Do() (rendered *RenderedBlock, err error) {
	return api.Post()
}

// Post executes a POST request to render the block.
func (api *RenderBlock) Post() (rendered *RenderedBlock, err error) {
	endpoint := api.client.endpoint + api.endpoint

	restyClient := api.client.httpClient.R()
	if api.client.auth.Username != "" && api.client.auth.Password != "" {
		restyClient = restyClient.SetBasicAuth(api.client.auth.Username, api.client.auth.Password)
	}

	body := make(map[string]any)
	if len(api.attributes) > 0 {
		body["attributes"] = api.attributes
	}
	if api.postID != 0 {
		body["post_id"] = api.postID
	}

	if api.context != "" {
		restyClient.SetQueryParam("context", api.context)
	}

	resp, err := restyClient.
		SetHeader("Content-Type", "application/json").
		SetHeader("Accept", "application/json").
		SetBody(body).
		SetResult(&rendered).
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

	return rendered, nil
}

// Get executes a GET request to render the block.
func (api *RenderBlock) Get() (rendered *RenderedBlock, err error) {
	endpoint := api.client.endpoint + api.endpoint

	restyClient := api.client.httpClient.R()
	if api.client.auth.Username != "" && api.client.auth.Password != "" {
		restyClient = restyClient.SetBasicAuth(api.client.auth.Username, api.client.auth.Password)
	}

	if api.context != "" {
		restyClient.SetQueryParam("context", api.context)
	}
	if api.postID != 0 {
		restyClient.SetQueryParam("post_id", strconv.Itoa(api.postID))
	}

	resp, err := restyClient.
		SetHeader("Accept", "application/json").
		SetResult(&rendered).
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

	return rendered, nil
}
