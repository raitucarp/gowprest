package tests

import (
	"os"
	"testing"

	_ "github.com/joho/godotenv/autoload"
	"github.com/raitucarp/gowprest"
	"github.com/stretchr/testify/assert"
)

func TestNewClient(t *testing.T) {
	blogUrl := os.Getenv("BLOG_URL")
	client := gowprest.NewClient(blogUrl)
	defer client.Close()

	blogInfo, err := client.Discover()

	assert.Equal(t, err, nil, "Something error %v", err)
	assert.Equal(t, blogInfo.Home, blogUrl, "Blog url not equal, expected %s, got %s", blogUrl, blogInfo.Home)
}
