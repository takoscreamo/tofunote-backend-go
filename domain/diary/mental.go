// Mental値オブジェクト: メンタルスコアを表現するモデル

package diary

type Mental struct {
	Value int
}

func NewMental(value int) Mental {
	if value < 1 || value > 10 {
		panic("Mental value must be between 1 and 10")
	}
	return Mental{Value: value}
}
