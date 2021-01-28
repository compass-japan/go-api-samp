package entity

type (
	Weather struct {
		Dat        string
		Weather    int
		LocationId int
		Comment    string
	}

	Location struct {
		Id   int
		City string
	}
)
