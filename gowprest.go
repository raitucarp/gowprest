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
	Authentication []string `json:"authentication"`
	SiteLogo       int      `json:"site_logo"`
	SiteIcon       int      `json:"site_icon"`
	SiteIconURL    string   `json:"site_icon_url"`
}

type Authentication struct {
	Username string
	Password string
}

type RestClient struct {
	baseURL string
	auth    Authentication

	httpClient *resty.Client
}

func (c *RestClient) endpoint() string {
	return strings.Trim(c.baseURL, "/") + "/wp-json/"
}

func (c *RestClient) close() {
	c.httpClient.Close()
}

func (c *RestClient) WithBasicAuth(username, password string) *RestClient {
	c.auth.Username = username
	c.auth.Password = password
	return c
}

func (c *RestClient) Discover() (info BlogInfo, err error) {
	_, err = c.httpClient.R().
		SetHeader("Accept", "application/json").
		SetResult(&info).
		Get(c.endpoint())

	defer c.httpClient.Close()

	if err != nil {
		return
	}

	return
}

func NewClient(baseURL string) *RestClient {
	client := resty.New()
	return &RestClient{baseURL: baseURL, httpClient: client}
}
