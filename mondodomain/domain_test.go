package mondodomain

import (
	"encoding/json"
	"testing"
	"time"
)

func TestMerchant_UnmarshalJSON_String(t *testing.T) {
	m := new(Merchant)
	err := json.Unmarshal([]byte(`"merch_0123456789"`), m)
	if err != nil {
		t.Fatalf("Failed to decode JSON string: %s", err)
	}

	expected := Merchant{ID: "merch_0123456789"}
	if *m != expected {
		t.Fatalf("Expected decoded %#v but got %#v", expected, m)
	}
}

func TestMerchant_UnmarshalJSON_Object(t *testing.T) {
	jsonMerch := `{
		"id": "merch_id",
		"name": "TestMerch",
		"created": "2016-02-09T12:30:35.000Z",
		"group_id": "grp_abcde",
		"logo": "http://company.com/logo.png",
		"emoji": "🔧",
		"category": "shopping"
	}`

	m := new(Merchant)
	err := json.Unmarshal([]byte(jsonMerch), m)
	if err != nil {
		t.Fatalf("Failed to decode JSON string: %s", err)
	}

	expected := Merchant{
		ID:       "merch_id",
		Name:     "TestMerch",
		Created:  time.Date(2016, 2, 9, 12, 30, 35, 0, time.UTC),
		GroupID:  "grp_abcde",
		Logo:     "http://company.com/logo.png",
		Emoji:    "🔧",
		Category: "shopping",
	}
	if *m != expected {
		t.Fatalf("Expected decoded %#v but got %#v", expected, m)
	}
}
