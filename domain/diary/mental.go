// Mental値オブジェクト: メンタルスコアを表現するモデル

package diary

import (
	"encoding/json"
	"errors"
)

type Mental int

func NewMental(value int) (Mental, error) {
	if value < 1 || value > 10 {
		return 0, errors.New("mental value must be between 1 and 10")
	}
	return Mental(value), nil
}

func (m Mental) Value() int {
	return int(m)
}

func (m Mental) MarshalJSON() ([]byte, error) {
	return json.Marshal(int(m))
}

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
