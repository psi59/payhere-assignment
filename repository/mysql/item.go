package mysql

import (
	"context"
	"fmt"
	"strconv"
	"time"

	"github.com/daangn/gorean"
	"github.com/pkg/errors"
	"github.com/psi59/payhere-assignment/domain"
	"github.com/psi59/payhere-assignment/internal/db"
	"github.com/psi59/payhere-assignment/internal/valid"
	"github.com/psi59/payhere-assignment/repository"
	"gorm.io/gorm"
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

	itemNameChosung, err := GetChosung(item.Name)
	if err != nil {
		return errors.WithStack(err)
	}

	// 2. 아이템 생성
	record := &Item{
		ItemID:          item.ID,
		UserID:          item.UserID,
		Category:        item.Category,
		ItemName:        item.Name,
		ItemNameChosung: itemNameChosung,
		Price:           item.Price,
		Cost:            item.Cost,
		Description:     item.Description,
		Barcode:         item.Barcode,
		ItemSize:        item.Size,
		ExpiryAt:        item.ExpiryAt,
		CreatedAt:       item.CreatedAt,
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

func (r *ItemRepository) Delete(c context.Context, userID, itemID int) error {
	switch {
	case valid.IsNil(c):
		return domain.ErrNilContext
	case userID < 1:
		return fmt.Errorf("invalid userID: %d", userID)
	case itemID < 1:
		return fmt.Errorf("invalid itemID: %d", itemID)
	}
	conn, err := db.ConnFromContext(c)
	if err != nil {
		return errors.WithStack(err)
	}

	var record Item
	if err := conn.Where("user_id=?", userID).Where("item_id=?", itemID).Delete(&record).Error; err != nil {
		return errors.WithStack(err)
	}

	return nil
}

func (r *ItemRepository) Update(c context.Context, userID, itemID int, input *repository.UpdateItemInput) error {
	switch {
	case valid.IsNil(c):
		return domain.ErrNilContext
	case userID < 1:
		return fmt.Errorf("invalid userID: %d", userID)
	case itemID < 1:
		return fmt.Errorf("invalid itemID: %d", itemID)
	case valid.IsNil(input):
		return domain.ErrNilInput
	}
	if err := input.Validate(); err != nil {
		return errors.WithStack(err)
	}

	conn, err := db.ConnFromContext(c)
	if err != nil {
		return errors.WithStack(err)
	}

	updateItem, err := createItemByUpdateItemInput(input)
	if err != nil {
		return errors.WithStack(err)
	}
	if err := conn.Model(&Item{}).Where("user_id = ?", userID).Where("item_id = ?", itemID).Updates(updateItem).Error; err != nil {
		if IsDuplicateEntry(err) {
			return errors.Wrap(domain.ErrItemAlreadyExists, err.Error())
		}
		return errors.WithStack(err)
	}

	return nil
}

func (r *ItemRepository) Find(c context.Context, input *repository.FindItemInput) (*repository.FindItemOutput, error) {
	switch {
	case valid.IsNil(c):
		return nil, domain.ErrNilContext
	case valid.IsNil(input):
		return nil, domain.ErrNilInput
	}
	if err := valid.ValidateStruct(input); err != nil {
		return nil, errors.WithStack(err)
	}

	conn, err := db.ConnFromContext(c)
	if err != nil {
		return nil, errors.WithStack(err)
	}

	totalCountInput := *input
	totalCountInput.SearchAfter = 0
	totalCount, err := r.getCount(conn, &totalCountInput)
	if err != nil {
		return nil, errors.WithStack(err)
	}

	queryBuilder := r.createFindQuery(conn, input)
	rows := make([]Item, 0)
	if err := queryBuilder.Find(&rows).Error; err != nil {
		return nil, errors.WithStack(err)
	}

	items := make([]domain.Item, len(rows))
	var searchAfter int
	for i := 0; i < len(rows); i++ {
		items[i] = *rows[i].Domain()
	}
	if len(items) > 0 {
		searchAfter = rows[len(items)-1].ItemID
	}

	hasNextInput := *input
	hasNextInput.SearchAfter = searchAfter
	nextItemCount, err := r.getCount(conn, &hasNextInput)
	if err != nil {
		return nil, errors.WithStack(err)
	}

	return &repository.FindItemOutput{
		TotalCount:  totalCount,
		Items:       items,
		HasNext:     nextItemCount > 0,
		SearchAfter: searchAfter,
	}, nil
}

func (r *ItemRepository) getCount(conn *gorm.DB, input *repository.FindItemInput) (int, error) {
	queryBuilder := r.createFindQuery(conn, input)
	var cnt int64
	if err := queryBuilder.Count(&cnt).Error; err != nil {
		return 0, errors.WithStack(err)
	}

	return int(cnt), nil
}

func (r *ItemRepository) createFindQuery(conn *gorm.DB, input *repository.FindItemInput) *gorm.DB {
	queryBuilder := conn.Model(&Item{}).Limit(10).Where("user_id=?", input.UserID).Order("item_id ASC")
	if input.SearchAfter > 0 {
		queryBuilder = queryBuilder.Where("item_id > ?", input.SearchAfter)
	}
	if k := input.Keyword; len(k) > 0 {
		queryBuilder = queryBuilder.Where("MATCH(item_name, item_name_chosung) AGAINST(? IN BOOLEAN MODE)", strconv.Quote(k))
	}

	return queryBuilder
}

type Item struct {
	ItemID          int             `gorm:"item_id;primaryKey"`
	UserID          int             `gorm:"user_id"`
	Category        string          `gorm:"category"`
	ItemName        string          `gorm:"item_name"`
	ItemNameChosung string          `gorm:"item_name_chosung"`
	Price           int             `gorm:"price"`
	Cost            int             `gorm:"cost"`
	Description     string          `gorm:"description"`
	Barcode         string          `gorm:"barcode"`
	ItemSize        domain.ItemSize `gorm:"item_size"`
	CreatedAt       time.Time       `gorm:"created_at"`
	ExpiryAt        time.Time       `gorm:"expiry_at"`
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

func createItemByUpdateItemInput(input *repository.UpdateItemInput) (*Item, error) {
	var item Item
	if !valid.IsNil(input.Name) {
		item.ItemName = *input.Name
		c, err := GetChosung(item.ItemName)
		if err != nil {
			return nil, errors.WithStack(err)
		}
		item.ItemNameChosung = c
	}
	if !valid.IsNil(input.Description) {
		item.Description = *input.Description
	}
	if !valid.IsNil(input.Price) {
		item.Price = *input.Price
	}
	if !valid.IsNil(input.Cost) {
		item.Cost = *input.Cost
	}
	if !valid.IsNil(input.Category) {
		item.Category = *input.Category
	}
	if !valid.IsNil(input.Barcode) {
		item.Barcode = *input.Barcode
	}
	if !valid.IsNil(input.Size) {
		item.ItemSize = *input.Size
	}
	if !valid.IsNil(input.ExpiryAt) {
		item.ExpiryAt = *input.ExpiryAt
	}

	return &item, nil
}

func GetChosung(s string) (string, error) {
	cc, err := gorean.Chosung(s)
	if err != nil {
		return "", errors.WithStack(err)
	}

	var result string
	for _, c := range cc {
		if c == "" {
			result += " "
		} else {
			result += c
		}
	}

	return result, nil
}
