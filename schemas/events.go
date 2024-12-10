package schemas

type YMDDate struct {
	Year  int `param:"year"`
	Month int `param:"month"`
	Day   int `param:"day"`
}

type YMDate struct {
	Year  int `param:"year"`
	Month int `param:"month"`
}

type YDate struct {
	Year int `param:"year"`
}
