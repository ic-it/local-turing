package internal

import (
	"encoding/json"
	"io"
)

/*
tests.json format:
{
    "test_name": [
        {
			"inputs": "input",
            "outputs": "output"
        },
        ...
    ]
}
*/

type (
	Tests map[string][]struct {
		Inputs  string `json:"inputs"`
		Outputs string `json:"outputs"`
	}
)

func TestsUnmarshal(r io.Reader) (Tests, error) {
	var tests Tests
	if err := json.NewDecoder(r).Decode(&tests); err != nil {
		return nil, err
	}
	return tests, nil
}

func TestsMarshal(tests Tests) ([]byte, error) {
	return json.Marshal(tests)
}
