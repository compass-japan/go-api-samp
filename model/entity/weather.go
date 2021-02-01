package entity

type (
	Weather struct {
		Dat      string
		Weather  int
		Location *Location
		Comment  string
	}

	Location struct {
		Id   int
		City string
	}
)
