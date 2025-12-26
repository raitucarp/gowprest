package gowprest

import (
	"strings"

	"resty.dev/v3"
)

type BlogInfo struct {
	Name           string   `json:"name"`
	Description    string   `json:"description"`
	URL            string   `json:"url"`
	Home           string   `json:"home"`
	GmtOffset      string   `json:"gmt_offset"`
	TimezoneString string   `json:"timezone_string"`
	PageForPosts   int      `json:"page_for_posts"`
	PageOnFront    int      `json:"page_on_front"`
	ShowOnFront    string   `json:"show_on_front"`
	Namespaces     []string `json:"namespaces"`
	Authentication any      `json:"authentication"`
	SiteLogo       int      `json:"site_logo"`
	SiteIcon       int      `json:"site_icon"`
	SiteIconURL    string   `json:"site_icon_url"`
}

type Authentication struct {
	Username string
	Password string
}

type RestClient struct {
	baseURL  string
	endpoint string
	auth     Authentication

	httpClient *resty.Client
}

func (api *RestClient) Close() {
	api.httpClient.Close()
}

func (api *RestClient) WithBasicAuth(username, password string) *RestClient {
	api.auth.Username = username
	api.auth.Password = password
	return api
}

func (api *RestClient) Discover() (info BlogInfo, err error) {
	_, err = api.httpClient.R().
		SetHeader("Accept", "application/json").
		SetResult(&info).
		Get(api.endpoint)

	if err != nil {
		return
	}

	return
}

func NewClient(baseURL string) *RestClient {
	client := resty.New()
	return &RestClient{
		baseURL:    baseURL,
		endpoint:   strings.Trim(baseURL, "/") + "/wp-json",
		httpClient: client}
}
