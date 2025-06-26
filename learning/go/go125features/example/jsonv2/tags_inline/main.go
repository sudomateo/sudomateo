package main

import (
	"encoding/json/jsontext"
	"encoding/json/v2"
	"os"
)

type APIRequest struct {
	Type   string `json:"type"`
	Person Person `json:",inline"`
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

	enc := jsontext.NewEncoder(os.Stdout, jsontext.WithIndent("  "))
	if err := json.MarshalEncode(enc, apiRequest); err != nil {
		panic(err)
	}
}
