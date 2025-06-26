package main

import (
	"bytes"
	"encoding/json/jsontext"
	"encoding/json/v2"
	"fmt"
	"io"
)

type APIRequest struct {
	Type string         `json:"type"`
	Data map[string]any `json:",unknown"`
}

type Person struct {
	FirstName string `json:"first_name"`
	LastName  string `json:"last_name"`
}

type Animal struct {
	Name string `json:"name"`
	Age  int    `json:"age"`
}

func main() {
	personJSON := bytes.NewBufferString(`{"type": "person", "person": {"first_name": "Matthew", "last_name": "Sanabria"}}`)
	animalJSON := bytes.NewBufferString(`{"type": "animal", "animal": {"name": "Ava", "age": 7}}`)
	handler(personJSON)
	handler(animalJSON)
}

func handler(r io.Reader) {
	var apiRequest APIRequest

	dec := jsontext.NewDecoder(r)
	if err := json.UnmarshalDecode(dec, &apiRequest); err != nil {
		panic(err)
	}

	switch apiRequest.Type {
	case "person":
		b, err := json.Marshal(apiRequest.Data["person"])
		if err != nil {
			panic(err)
		}

		var p Person
		if err := json.Unmarshal(b, &p); err != nil {
			panic(err)
		}

		fmt.Println(p)
	case "animal":
		b, err := json.Marshal(apiRequest.Data["animal"])
		if err != nil {
			panic(err)
		}

		var a Animal
		if err := json.Unmarshal(b, &a); err != nil {
			panic(err)
		}
		fmt.Println(a)
	}
}
