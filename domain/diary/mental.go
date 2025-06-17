// Mental値オブジェクト: メンタルスコアを表現するモデル

package diary

import "errors"

type Mental struct {
	Value int
}

func NewMental(value int) (Mental, error) {
	if value < 1 || value > 10 {
		return Mental{}, errors.New("mental value must be between 1 and 10")
	}
	return Mental{Value: value}, nil
}
