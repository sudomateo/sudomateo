package main

import (
	"bytes"
	"encoding/json/jsontext"
	"encoding/json/v2"
	"fmt"
)

var _ json.MarshalerTo = (*APIRequest)(nil)
var _ json.UnmarshalerFrom = (*APIRequest)(nil)

type APIRequest struct {
	Type   string `json:"type"`
	Person Person `json:"person"`
}

type Person struct {
	FirstName string `json:"first_name"`
	LastName  string `json:"last_name"`
}

func (a *APIRequest) UnmarshalJSONFrom(dec *jsontext.Decoder) error {
	if dec.PeekKind() != '{' {
		return fmt.Errorf("expected JSON object, got %v", dec.PeekKind())
	}

	if _, err := dec.ReadToken(); err != nil {
		return err
	}

	for dec.PeekKind() != '}' {
		tok, err := dec.ReadToken()
		if err != nil {
			return err
		}

		fieldName := tok.String()

		switch fieldName {
		case "type":
			tok, err := dec.ReadToken()
			if err != nil {
				return err
			}
			a.Type = tok.String()

		case "person":
			if err := json.UnmarshalDecode(dec, &a.Person); err != nil {
				return fmt.Errorf("failed to unmarshal person: %w", err)
			}

		default:
			if err := dec.SkipValue(); err != nil {
				return fmt.Errorf("failed to skip unknown field %q: %w", fieldName, err)
			}
		}
	}

	if _, err := dec.ReadToken(); err != nil {
		return err
	}

	return nil
}

func (a *APIRequest) MarshalJSONTo(enc *jsontext.Encoder) error {
	if err := enc.WriteToken(jsontext.BeginObject); err != nil {
		return err
	}

	if err := enc.WriteToken(jsontext.String("type")); err != nil {
		return err
	}
	if err := enc.WriteToken(jsontext.String(a.Type)); err != nil {
		return err
	}

	if err := enc.WriteToken(jsontext.String("person")); err != nil {
		return err
	}

	if err := json.MarshalEncode(enc, &a.Person); err != nil {
		return fmt.Errorf("failed to marshal person: %w", err)
	}

	if err := enc.WriteToken(jsontext.EndObject); err != nil {
		return err
	}

	return nil
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

	// MarshalDecode.
	var buf bytes.Buffer
	enc := jsontext.NewEncoder(&buf)
	if err := json.MarshalEncode(enc, apiRequest); err != nil {
		panic(err)
	}
	fmt.Printf("json.MarshalEncode:\n%s\n\n", buf.String())

	// UnmarshalDecode.
	apiRequestBytes := bytes.NewBuffer(apiRequestJSON)
	dec := jsontext.NewDecoder(apiRequestBytes)
	if err := json.UnmarshalDecode(dec, &apiRequest); err != nil {
		panic(err)
	}
	fmt.Printf("json.UnmarshalDecode:\n%#v\n\n", apiRequest)
}
