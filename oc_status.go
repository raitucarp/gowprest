package gowprest

import (
	"encoding/json"
	"fmt"
)

type OpenClosedStatus string

const (
	StatusOpen   OpenClosedStatus = "open"
	StatusClosed OpenClosedStatus = "closed"
)

func (s OpenClosedStatus) MarshalJSON() ([]byte, error) {
	return json.Marshal(string(s))
}

func (s *OpenClosedStatus) UnmarshalJSON(data []byte) error {
	var ocStatus string
	if err := json.Unmarshal(data, &ocStatus); err != nil {
		return err
	}

	switch ocStatus {
	case string(StatusOpen), string(StatusClosed):
		*s = OpenClosedStatus(ocStatus)
		return nil
	default:
		return fmt.Errorf("invalid suit value: %s", ocStatus)
	}
}
