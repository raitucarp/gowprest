package gowprest

import (
	"encoding/json"
	"fmt"
)

type Format string

const (
	FormatStandard Format = "standard"
	FormatAside    Format = "aside"
	FormatAudio    Format = "audio"
	FormatChat     Format = "chat"
	FormatGallery  Format = "gallery"
	FormatImage    Format = "image"
	FormatLink     Format = "link"
	FormatQuote    Format = "quote"
	FormatStatus   Format = "status"
	FormatVideo    Format = "video"
)

func (s Format) MarshalJSON() ([]byte, error) {
	return json.Marshal(string(s))
}
func (s *Format) UnmarshalJSON(data []byte) error {
	var format string
	if err := json.Unmarshal(data, &format); err != nil {
		return err
	}
	switch format {
	case string(FormatStandard):
		*s = FormatStandard
	case string(FormatAside):
		*s = FormatAside
	case string(FormatAudio):
		*s = FormatAudio
	case string(FormatChat):
		*s = FormatChat
	case string(FormatGallery):
		*s = FormatGallery
	case string(FormatImage):
		*s = FormatImage
	case string(FormatLink):
		*s = FormatLink
	case string(FormatQuote):
		*s = FormatQuote
	case string(FormatStatus):
		*s = FormatStatus
	case string(FormatVideo):
		*s = FormatVideo
	default:
		return fmt.Errorf("invalid format: %s", format)
	}
	return nil
}
