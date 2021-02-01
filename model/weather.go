package model

type wtr uint16

const (
	Sunny wtr = iota + 1
	Cloudy
	Rainy
	Snowy
)

var (
	wtrMap = map[wtr]string{
		0: "",
		1: "Sunny",
		2: "Cloudy",
		3: "Rainy",
		4: "Snowy",
	}
)

func ToWeather(i int) string {
	if v, ok := wtrMap[wtr(i)]; ok {
		return v
	}
	return ""
}
