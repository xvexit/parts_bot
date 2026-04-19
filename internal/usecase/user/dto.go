package user

type UserInput struct {
	Name     string
	Password string
	Phone    string
	Email    string
}

func NewUserInput(
	name string,
	phone string,
	password string,
	email string,
) UserInput {
	return UserInput{
		Name:     name,
		Password: password,
		Phone:    phone,
		Email:    email,
	}
}
