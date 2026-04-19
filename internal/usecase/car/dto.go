package car

type CarInput struct {
	Name   string
	VIN    string
	UserID int64
}

func NewCarInput(
	name string,
	vin string,
	userID int64,
) CarInput {
	return CarInput{
		Name:   name,
		VIN:    vin,
		UserID: userID,
	}
}
