package entity

type (
	Weather struct {
		Dat        string
		Weather    int
		LocationId int //todo obj
		Comment    string
	}

	Location struct {
		Id   int
		City string
	}
)
