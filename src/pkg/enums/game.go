package enums

const fms2023 = "fms2023"

type Game string

func (e Game) String() string {
	return string(e)
}

type games struct{}

func (games) ForzaMotorsport2023() Game { return fms2023 }

var Games games
