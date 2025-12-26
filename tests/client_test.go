package tests

import (
	"os"
	"testing"

	_ "github.com/joho/godotenv/autoload"
	"github.com/raitucarp/gowprest"
)

func TestNewClient(t *testing.T) {
	blogUrl := os.Getenv("BLOG_URL")
	client := gowprest.NewClient(blogUrl)
	defer client.Close()

	blogInfo, err := client.Discover()

	assert(t, err != nil, "Something error %v", err)
	assert(t, blogInfo.Home != blogUrl, "Blog url not equal, expected %s, got %s", blogUrl, blogInfo.Home)
}
