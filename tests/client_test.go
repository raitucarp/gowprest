package tests

import (
	"fmt"
	"testing"

	"github.com/raitucarp/gowprest"
)

func TestNewClient(t *testing.T) {
	client := gowprest.NewClient("http://localhost:8888")
	blogInfo, err := client.Discover()

	if err != nil {
		t.Errorf("Something error %v", err)
	}

	fmt.Println(blogInfo)
}
