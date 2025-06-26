package main

import (
	"bytes"
	"encoding/json/jsontext"
	"encoding/json/v2"
	"fmt"
)

type APIRequest struct {
	Type   string `json:"type"`
	Person Person `json:"person"`
}

type Person struct {
	FirstName string `json:"first_name"`
	LastName  string `json:"last_name"`
}

func main() {
	apiRequest := APIRequest{
		Type: "person",
		Person: Person{
			FirstName: "Matthew",
			LastName:  "Sanabria",
		},
	}
	apiRequestJSON := []byte(`{"type": "person", "person": {"first_name": "Matthew", "last_name": "Sanabria"}}`)

	// Marshal.
	b, err := json.Marshal(apiRequest)
	if err != nil {
		panic(err)
	}
	fmt.Printf("json.Marshal:\n%s\n\n", string(b))

	// MarshalEncode.
	var buf bytes.Buffer
	enc := jsontext.NewEncoder(&buf)
	if err := json.MarshalEncode(enc, apiRequest); err != nil {
		panic(err)
	}
	fmt.Printf("json.MarshalEncode:\n%s\n\n", buf.String())

	// Unmarshal.
	if err := json.Unmarshal(apiRequestJSON, &apiRequest); err != nil {
		panic(err)
	}
	fmt.Printf("json.Unmarshal:\n%#v\n\n", apiRequest)

	// UnmarshalDecode.
	apiRequestBytes := bytes.NewBuffer(apiRequestJSON)
	dec := jsontext.NewDecoder(apiRequestBytes)
	if err := json.UnmarshalDecode(dec, &apiRequest); err != nil {
		panic(err)
	}
	fmt.Printf("json.UnmarshalDecode:\n%#v\n\n", apiRequest)
}
