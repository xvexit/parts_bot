package cart

import (
	"context"
	"partsBot/internal/domain/cart"
	"partsBot/internal/domain/order"
	"partsBot/internal/infrastructure/db"
	"partsBot/pkg/errors"
	"partsBot/pkg/money"
)

type Service struct {
	cartRepo  cart.Repository
	orderRepo order.Repository
	txManager *db.TxManager
}

func NewService(
	cartRepo cart.Repository,
	orderRepo order.Repository,
	txManager *db.TxManager,
) *Service {
	return &Service{
		cartRepo:  cartRepo,
		orderRepo: orderRepo,
		txManager: txManager,
	}
}

func (s *Service) AddItem(
	ctx context.Context,
	userID int64,
	dto CartItemDto,
) error {

	cartAgg, err := s.cartRepo.GetByUserID(ctx, userID)
	if err != nil {
		return err
	}

	if cartAgg.IsEmpty() { //удалить
		cartAgg, err = cart.NewCart(userID)
		if err != nil {
			return err
		}
	}

	price, err := money.New(dto.Price)
	if err != nil {
		return err
	}

	item, err := cart.NewCartItem(
		dto.PartID,
		dto.Name,
		dto.Brand,
		dto.ImageURL,
		dto.Quantity,
		dto.DeliveryDay,
		price,
	)
	if err != nil {
		return err
	}

	if err := cartAgg.AddItem(*item); err != nil {
		return err
	}

	return s.cartRepo.Save(ctx, cartAgg)
}

func (s *Service) RemoveItem(сtx context.Context, userID int64, partID string) error {
	cart, err := s.cartRepo.GetByUserID(сtx, userID)
	if err != nil {
		return err
	}

	if err := cart.DeleteItem(partID); err != nil {
		return err
	}

	if err := s.cartRepo.Save(сtx, cart); err != nil {
		return err
	}

	return nil
}

func (s *Service) GetCart(ctx context.Context, userID int64) ([]cart.CartItem, error) {
	cartAgg, err := s.cartRepo.GetByUserID(ctx, userID)
	if err != nil {
		return nil, err
	}

	if cartAgg == nil {
		cartAgg, _ = cart.NewCart(userID)
	}

	return cartAgg.Items(), nil
}

func (s *Service) ClearCart(сtx context.Context, userID int64) error {
	cart, err := s.cartRepo.GetByUserID(сtx, userID)
	if err != nil {
		return err
	}

	cart.Clear()

	return s.cartRepo.Save(сtx, cart)
}

func (s *Service) Checkout(
	ctx context.Context,
	userID int64,
	address string,
) (*order.Order, error) {

	var ord *order.Order

	err := s.txManager.WithinTx(ctx, func(txCtx context.Context) error {

		cart, err := s.cartRepo.GetByUserIDForUpdate(txCtx, userID)
		if err != nil {
			return err
		}

		if cart.IsEmpty() {
			return errors.ErrCartIsEmpty
		}

		ord, err = cart.NewOrderFromCart(address)
		if err != nil {
			return err
		}

		if err := s.orderRepo.Create(txCtx, ord); err != nil {
			return err
		}

		cart.Clear()

		return s.cartRepo.Save(txCtx, cart)
	})

	if err != nil {
		return nil, err
	}

	return ord, nil
}

// updateStatus