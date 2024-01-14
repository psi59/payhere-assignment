package mysql

import (
	"context"
	"fmt"
	"time"

	"gorm.io/gorm"

	"github.com/psi59/payhere-assignment/internal/db"

	"github.com/pkg/errors"

	"github.com/psi59/payhere-assignment/internal/valid"

	"github.com/psi59/payhere-assignment/domain"
)

type ItemRepository struct{}

func NewItemRepository() *ItemRepository {
	return &ItemRepository{}
}

func (r *ItemRepository) Create(c context.Context, item *domain.Item) error {
	// 1. 파라메터 체크
	switch {
	case valid.IsNil(c):
		return domain.ErrNilContext
	case valid.IsNil(item):
		return domain.ErrNilItem
	}
	if err := item.Validate(); err != nil {
		return errors.WithStack(err)
	}
	conn, err := db.ConnFromContext(c)
	if err != nil {
		return errors.WithStack(err)
	}

	// 2. 아이템 생성
	record := &Item{
		ItemID:      item.ID,
		UserID:      item.UserID,
		Category:    item.Category,
		ItemName:    item.Name,
		Price:       item.Price,
		Cost:        item.Cost,
		Description: item.Description,
		Barcode:     item.Barcode,
		ItemSize:    item.Size,
		ExpiryAt:    item.ExpiryAt,
		CreatedAt:   item.CreatedAt,
	}
	if err := conn.Create(record).Error; err != nil {
		if IsDuplicateEntry(err) {
			return errors.Wrap(domain.ErrItemAlreadyExists, err.Error())
		}

		return errors.WithStack(err)
	}
	item.ID = record.ItemID

	// 3. 결과 반환
	return nil
}

func (r *ItemRepository) Get(c context.Context, userID, itemID int) (*domain.Item, error) {
	switch {
	case valid.IsNil(c):
		return nil, domain.ErrNilContext
	case userID < 1:
		return nil, fmt.Errorf("invalid userID: %d", userID)
	case itemID < 1:
		return nil, fmt.Errorf("invalid itemID: %d", itemID)
	}
	conn, err := db.ConnFromContext(c)
	if err != nil {
		return nil, errors.WithStack(err)
	}

	var record Item
	if err := conn.Where("user_id=?", userID).Where("item_id=?", itemID).Take(&record).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, errors.Wrap(domain.ErrItemNotFound, err.Error())
		}

		return nil, errors.WithStack(err)
	}

	return record.Domain(), nil
}

type Item struct {
	ItemID      int             `gorm:"item_id;primaryKey"`
	UserID      int             `gorm:"user_id"`
	Category    string          `gorm:"category"`
	ItemName    string          `gorm:"item_name"`
	Price       int             `gorm:"price"`
	Cost        int             `gorm:"cost"`
	Description string          `gorm:"description"`
	Barcode     string          `gorm:"barcode"`
	ItemSize    domain.ItemSize `gorm:"item_size"`
	CreatedAt   time.Time       `gorm:"created_at"`
	ExpiryAt    time.Time       `gorm:"expiry_at"`
}

func (i *Item) TableName() string {
	return "items"
}

func (i *Item) Domain() *domain.Item {
	return &domain.Item{
		ID:          i.ItemID,
		UserID:      i.UserID,
		Name:        i.ItemName,
		Description: i.Description,
		Price:       i.Price,
		Cost:        i.Cost,
		Category:    i.Category,
		Barcode:     i.Barcode,
		ExpiryAt:    i.ExpiryAt,
		Size:        i.ItemSize,
		CreatedAt:   i.CreatedAt,
	}
}
