package internal

import (
	"encoding/json"
	"io"
)

type (
	Tests struct {
		Schema string            `json:"$schema"`
		Tests  map[string][]Test `json:"tests"`
	}
	Test struct {
		Inputs  string `json:"inputs"`
		Outputs string `json:"outputs"`
	}
)

func TestsUnmarshal(r io.Reader) (Tests, error) {
	var tests Tests
	if err := json.NewDecoder(r).Decode(&tests); err != nil {
		return tests, err
	}
	return tests, nil
}

func TestsMarshal(tests Tests) ([]byte, error) {
	return json.Marshal(tests)
}
