// Mental値オブジェクト: メンタルスコアを表現するモデル

package diary

import (
	"encoding/json"
	"errors"
)

type Mental struct {
	Value int `json:"-"`
}

func NewMental(value int) (Mental, error) {
	if value < 1 || value > 10 {
		return Mental{}, errors.New("mental value must be between 1 and 10")
	}
	return Mental{Value: value}, nil
}

// Value returns the underlying int value
func (m Mental) GetValue() int {
	return m.Value
}

// MarshalJSON implements json.Marshaler interface
func (m Mental) MarshalJSON() ([]byte, error) {
	return json.Marshal(m.Value)
}

// UnmarshalJSON implements json.Unmarshaler interface
func (m *Mental) UnmarshalJSON(data []byte) error {
	var value int
	if err := json.Unmarshal(data, &value); err != nil {
		return err
	}

	mental, err := NewMental(value)
	if err != nil {
		return err
	}

	*m = mental
	return nil
}
