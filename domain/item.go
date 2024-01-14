package domain

import (
	"fmt"
	"time"

	"github.com/psi59/payhere-assignment/internal/valid"

	"github.com/pkg/errors"
)

type Item struct {
	ID          int
	UserID      int       `validate:"gt=0"`
	Name        string    `validate:"required,gte=1,lte=100"`
	Description string    `validate:"required"`
	Price       int       `validate:"required"`
	Cost        int       `validate:"required"`
	Category    string    `validate:"required,gte=1,lte=100"`
	Barcode     string    `validate:"required,gte=1,lte=100"`
	ExpiryAt    time.Time `validate:"required"`
	Size        ItemSize  `validate:"required"`
	CreatedAt   time.Time `validate:"required"`
}

const ErrNilItem ConstantError = "nil Item"

func NewItem(
	userID int,
	name string,
	description string,
	price int,
	cost int,
	category string,
	barcode string,
	expiryAt time.Time,
	size ItemSize,
) (*Item, error) {
	switch {
	case userID < 1:
		return nil, fmt.Errorf("invalid userID: %d", userID)
	case len(name) == 0:
		return nil, fmt.Errorf("empty name")
	case len(description) == 0:
		return nil, fmt.Errorf("empty description")
	case price < 1:
		return nil, fmt.Errorf("invalid price: %d", price)
	case cost < 1:
		return nil, fmt.Errorf("invalid cost: %d", cost)
	case len(barcode) == 0:
		return nil, fmt.Errorf("empty barcode")
	case expiryAt.IsZero():
		return nil, fmt.Errorf("zero expiryAt")
	}
	if err := size.Validate(); err != nil {
		return nil, errors.WithStack(err)
	}

	item := &Item{
		UserID:      userID,
		Name:        name,
		Description: description,
		Price:       price,
		Cost:        cost,
		Category:    category,
		Barcode:     barcode,
		ExpiryAt:    expiryAt,
		Size:        size,
		CreatedAt:   time.Now(),
	}

	return item, nil
}

func (i *Item) Validate() error {
	if err := valid.ValidateStruct(i); err != nil {
		return errors.WithStack(err)
	}
	if err := i.Size.Validate(); err != nil {
		return errors.WithStack(err)
	}

	return nil
}

type ItemSize string

const (
	ItemSizeSmall ItemSize = "small"
	ItemSizeLarge ItemSize = "large"
)

func (s ItemSize) Validate() error {
	switch s {
	case ItemSizeSmall, ItemSizeLarge:
		return nil
	default:
		return fmt.Errorf("undefined ItemSize")
	}
}
