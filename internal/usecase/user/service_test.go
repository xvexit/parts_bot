package user

// import (
// 	"context"
// 	"testing"
// 	"time"

// 	domain "partsBot/internal/domain/user"
// )

// type stubUserRepo struct {
// 	createFn  func(ctx context.Context, u *domain.User) error
// 	updateFn  func(ctx context.Context, u *domain.User) error
// 	deleteFn  func(ctx context.Context, id int64) error
// 	getByIDFn func(ctx context.Context, id int64) (*domain.User, error)
// 	getByTgFn func(ctx context.Context, tgID int64) (*domain.User, error)
// }

// func (s *stubUserRepo) Create(ctx context.Context, u *domain.User) error {
// 	return s.createFn(ctx, u)
// }

// func (s *stubUserRepo) Update(ctx context.Context, u *domain.User) error {
// 	return s.updateFn(ctx, u)
// }

// func (s *stubUserRepo) Delete(ctx context.Context, id int64) error {
// 	return s.deleteFn(ctx, id)
// }

// func (s *stubUserRepo) GetByID(ctx context.Context, id int64) (*domain.User, error) {
// 	return s.getByIDFn(ctx, id)
// }

// func (s *stubUserRepo) GetByTgID(ctx context.Context, tgID int64) (*domain.User, error) {
// 	return s.getByTgFn(ctx, tgID)
// }

// func TestRegister(t *testing.T) {

// 	repo := &stubUserRepo{
// 		createFn: func(ctx context.Context, u *domain.User) error {
// 			u.SetID(1)
// 			return nil
// 		},
// 	}

// 	service := Service{repo: repo}

// 	dto := UserInput{
// 		ID:    1,
// 		Name:  "Ivan",
// 		Phone: "7999",
// 	}

// 	u, err := service.Register(context.Background(), dto)

// 	if err != nil {
// 		t.Fatal(err)
// 	}

// 	if u.ID() != 1 {
// 		t.Fatal("id not set")
// 	}
// }

// func TestChangeName(t *testing.T) {

// 	u := domain.RestoreUser(1, 1, "Ivan", "7999", time.Now())

// 	repo := &stubUserRepo{
// 		getByIDFn: func(ctx context.Context, id int64) (*domain.User, error) {
// 			return u, nil
// 		},
// 		updateFn: func(ctx context.Context, u *domain.User) error {
// 			return nil
// 		},
// 	}

// 	service := Service{repo: repo}

// 	res, err := service.ChangeName(context.Background(), "Petr", 1)

// 	if err != nil {
// 		t.Fatal(err)
// 	}

// 	if res.Name() != "Petr" {
// 		t.Fatal("name not changed")
// 	}
// }

// func TestChangePhone(t *testing.T) {

// 	u := domain.RestoreUser(1, 1, "Ivan", "7999", time.Now())

// 	repo := &stubUserRepo{
// 		getByIDFn: func(ctx context.Context, id int64) (*domain.User, error) {
// 			return u, nil
// 		},
// 		updateFn: func(ctx context.Context, u *domain.User) error {
// 			return nil
// 		},
// 	}

// 	service := Service{repo: repo}

// 	res, err := service.ChangePhone(context.Background(), "8888", 1)

// 	if err != nil {
// 		t.Fatal(err)
// 	}

// 	if res.Phone() != "8888" {
// 		t.Fatal("phone not changed")
// 	}
// }

// func TestDeleteUser(t *testing.T) {

// 	repo := &stubUserRepo{
// 		deleteFn: func(ctx context.Context, id int64) error {
// 			return nil
// 		},
// 	}

// 	service := Service{repo: repo}

// 	err := service.Delete(context.Background(), 1)

// 	if err != nil {
// 		t.Fatal(err)
// 	}
// }

// func TestGetByID(t *testing.T) {

// 	u := domain.RestoreUser(1, 1, "Ivan", "7999", time.Now())

// 	repo := &stubUserRepo{
// 		getByIDFn: func(ctx context.Context, id int64) (*domain.User, error) {
// 			return u, nil
// 		},
// 	}

// 	service := Service{repo: repo}

// 	res, err := service.GetByID(context.Background(), 1)

// 	if err != nil {
// 		t.Fatal(err)
// 	}

// 	if res.ID() != 1 {
// 		t.Fatal("wrong id")
// 	}
// }

// func TestGetByTgID(t *testing.T) {

// 	u := domain.RestoreUser(1, 1, "Ivan", "7999", time.Now())

// 	repo := &stubUserRepo{
// 		getByTgFn: func(ctx context.Context, id int64) (*domain.User, error) {
// 			return u, nil
// 		},
// 	}

// 	service := Service{repo: repo}

// 	res, err := service.GetByTgID(context.Background(), 1)

// 	if err != nil {
// 		t.Fatal(err)
// 	}

// 	if res.TelegramID() != 1 {
// 		t.Fatal("wrong tg id")
// 	}
// }
