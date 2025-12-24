package tests

import (
	"fmt"
	"os"
	"testing"

	_ "github.com/joho/godotenv/autoload"
	"github.com/raitucarp/gowprest"
)

func TestNewClient(t *testing.T) {
	blogUrl := os.Getenv("BLOG_URL")
	client := gowprest.NewClient(blogUrl)
	blogInfo, err := client.Discover()

	if err != nil {
		t.Errorf("Something error %v", err)
	}

	fmt.Println(blogInfo)
}
