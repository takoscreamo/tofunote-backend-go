package diary

import (
	"encoding/json"
	"testing"
)

func TestMental_MarshalJSON(t *testing.T) {
	mental, err := NewMental(8)
	if err != nil {
		t.Fatalf("NewMental failed: %v", err)
	}

	data, err := json.Marshal(mental)
	if err != nil {
		t.Fatalf("json.Marshal failed: %v", err)
	}

	expected := "8"
	if string(data) != expected {
		t.Errorf("Expected %s, got %s", expected, string(data))
	}
}

func TestMental_UnmarshalJSON(t *testing.T) {
	var mental Mental
	err := json.Unmarshal([]byte("8"), &mental)
	if err != nil {
		t.Fatalf("json.Unmarshal failed: %v", err)
	}

	if int(mental) != 8 {
		t.Errorf("Expected 8, got %d", int(mental))
	}
}

func TestDiary_MarshalJSON(t *testing.T) {
	mental, _ := NewMental(7)
	diary := Diary{
		ID:     1,
		UserID: 200,
		Date:   "2025-01-20",
		Mental: mental,
		Diary:  "テスト日記",
	}

	data, err := json.Marshal(diary)
	if err != nil {
		t.Fatalf("json.Marshal failed: %v", err)
	}

	jsonStr := string(data)
	if jsonStr == `{"ID":1,"UserID":200,"Date":"2025-01-20","Mental":{"Value":7},"Diary":"テスト日記"}` {
		t.Errorf("Mental should be serialized as int, not object. Got: %s", jsonStr)
	}
	if jsonStr == `{"ID":1,"UserID":200,"Date":"2025-01-20","Mental":7,"Diary":"テスト日記"}` {
		t.Logf("Mental is correctly serialized as int: %s", jsonStr)
	} else {
		t.Errorf("Unexpected JSON format: %s", jsonStr)
	}
}
