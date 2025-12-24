package gowprest

import (
	"encoding/json"
	"fmt"
)

type PostStatus string

const (
	StatusDraft     PostStatus = "draft"
	StatusPending   PostStatus = "pending"
	StatusPrivate   PostStatus = "private"
	StatusPublished PostStatus = "publish"
)

func (s PostStatus) MarshalJSON() ([]byte, error) {
	return json.Marshal(string(s))
}

func (s *PostStatus) UnmarshalJSON(data []byte) error {
	var postStatus string
	if err := json.Unmarshal(data, &postStatus); err != nil {
		return err
	}

	switch postStatus {
	case string(StatusDraft), string(StatusPending), string(StatusPrivate), string(StatusPublished):
		*s = PostStatus(postStatus)
		return nil
	default:
		return fmt.Errorf("invalid suit value: %s", postStatus)
	}
}
