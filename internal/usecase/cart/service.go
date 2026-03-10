package cart

import (
	"context"
	"partsBot/internal/domain/cart"
	"partsBot/internal/domain/order"
	"partsBot/pkg/errors"
	"partsBot/pkg/money"
)

type Service struct {
	cartRepo  cart.Repository
	orderRepo order.Repository
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

	if cartAgg == nil {
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

	cart, err := s.cartRepo.GetByUserID(ctx, userID)
	if err != nil {
		return nil, err
	}

	if cart == nil {
		return nil, errors.ErrCartIsEmpty
	}

	order, err := cart.NewOrderFromCart(address)
	if err != nil {
		return nil, err
	}

	err = s.orderRepo.Save(ctx, order) //сделать атомарно
	if err != nil {
		return nil, err
	}

	cart.Clear()

	err = s.cartRepo.Save(ctx, cart)
	if err != nil {
		return nil, err
	}

	return order, nil
}
