package car

import "context"

type Repository interface {
    Create(ctx context.Context, car *Car) error
    Update(ctx context.Context, car *Car) error
    GetByID(ctx context.Context, id int64) (*Car, error)
    ListByUser(ctx context.Context, userID int64) ([]Car, error)
    Delete(ctx context.Context, id int64) error
}