package car

import "errors"

var(
	ErrNameCanNotBeNull = errors.New("название авто не может быть пустым")
	ErrVinCanNotBeNull = errors.New("вин код не может быть пустым")
)