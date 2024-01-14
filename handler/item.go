package handler

import (
	"net/http"
	"strconv"
	"time"

	"github.com/psi59/payhere-assignment/usecase/item"

	"github.com/psi59/payhere-assignment/internal/valid"

	"github.com/pkg/errors"

	"github.com/psi59/payhere-assignment/internal/ginhelper"
	"github.com/psi59/payhere-assignment/internal/i18n"

	"github.com/psi59/payhere-assignment/domain"

	"github.com/gin-gonic/gin"
)

type ItemHandler struct {
	itemUsecase item.Usecase
}

func NewItemHandler(itemUsecase item.Usecase) (*ItemHandler, error) {
	if valid.IsNil(itemUsecase) {
		return nil, item.ErrNilUsecase
	}

	return &ItemHandler{itemUsecase: itemUsecase}, nil
}

func (h *ItemHandler) Create(ginCtx *gin.Context) {
	ctx := ginhelper.GetContext(ginCtx)

	// 1. 인증된 유저 확인
	user, ok := ctx.Value(domain.CtxKeyUser).(*domain.User)
	if !ok {
		ginhelper.Error(ginCtx, errors.New("unauthenticated request"))
		return
	}

	// 2. 요청 검증
	var req CreateItemRequest
	if err := ginCtx.BindJSON(&req); err != nil {
		ginhelper.Error(ginCtx, ginhelper.NewHTTPError(http.StatusBadRequest, i18n.InvalidRequest, errors.WithStack(err)))
		return
	}
	if err := valid.ValidateStruct(req); err != nil {
		ginhelper.Error(ginCtx, ginhelper.NewHTTPError(http.StatusBadRequest, i18n.InvalidRequest, errors.WithStack(err)))
		return
	}

	// 3. 아이템 생성
	createItemOutput, err := h.itemUsecase.Create(ctx, &item.CreateInput{
		User:        user,
		Name:        req.Name,
		Description: req.Description,
		Price:       req.Price,
		Cost:        req.Cost,
		Category:    req.Category,
		Barcode:     req.Barcode,
		Size:        req.Size,
		ExpiryAt:    req.ExpiryAt,
	})
	if err != nil {
		//// 3.1 에러 처리
		if errors.Is(err, domain.ErrItemAlreadyExists) {
			ginhelper.Error(ginCtx, ginhelper.NewHTTPError(http.StatusConflict, i18n.ItemAlreadyExists, errors.WithStack(err)))
			return
		}

		ginhelper.Error(ginCtx, errors.WithStack(err))
		return
	}
	itemDomain := createItemOutput.Item

	// 4. 응답 반환
	ginhelper.Success(ginCtx, CreateItemResponse{
		ID:          itemDomain.ID,
		Name:        itemDomain.Name,
		Description: itemDomain.Description,
		Price:       itemDomain.Price,
		Cost:        itemDomain.Cost,
		Category:    itemDomain.Category,
		Barcode:     itemDomain.Barcode,
		Size:        itemDomain.Size,
		ExpiryAt:    itemDomain.ExpiryAt,
		CreatedAt:   itemDomain.CreatedAt,
	})
}

func (h *ItemHandler) Get(ginCtx *gin.Context) {
	ctx := ginhelper.GetContext(ginCtx)

	// 1. 인증된 유저 확인
	user, ok := ctx.Value(domain.CtxKeyUser).(*domain.User)
	if !ok {
		ginhelper.Error(ginCtx, errors.New("unauthenticated request"))
		return
	}

	itemIDParam := ginCtx.Param("itemId")
	itemID, err := strconv.Atoi(itemIDParam)
	if err != nil {
		ginhelper.Error(ginCtx, ginhelper.NewHTTPError(http.StatusNotFound, i18n.ItemNotFound, errors.WithStack(err)))
		return
	}

	getItemOutput, err := h.itemUsecase.Get(ctx, &item.GetInput{
		User:   user,
		ItemID: itemID,
	})
	if err != nil {
		if errors.Is(err, domain.ErrItemNotFound) {
			ginhelper.Error(ginCtx, ginhelper.NewHTTPError(http.StatusNotFound, i18n.ItemNotFound, errors.WithStack(err)))
			return
		}

		ginhelper.Error(ginCtx, errors.WithStack(err))
		return
	}
	itemDomain := getItemOutput.Item

	ginhelper.Success(ginCtx, GetItemResponse{
		ID:          itemDomain.ID,
		Name:        itemDomain.Name,
		Description: itemDomain.Description,
		Price:       itemDomain.Price,
		Cost:        itemDomain.Cost,
		Category:    itemDomain.Category,
		Barcode:     itemDomain.Barcode,
		Size:        itemDomain.Size,
		ExpiryAt:    itemDomain.ExpiryAt,
		CreatedAt:   itemDomain.CreatedAt,
	})
}

type CreateItemRequest struct {
	Name        string          `json:"name" validate:"required,gte=1,lte=100"`
	Description string          `json:"description" validate:"required"`
	Price       int             `json:"price" validate:"gt=0"`
	Cost        int             `json:"cost" validate:"gt=0"`
	Category    string          `json:"category" validate:"required,gte=1,lte=100"`
	Barcode     string          `json:"barcode" validate:"required,gte=1,lte=100"`
	Size        domain.ItemSize `json:"size" validate:"required,oneof=small large"`
	ExpiryAt    time.Time       `json:"expiryAt" validate:"required"`
}

type CreateItemResponse struct {
	ID          int             `json:"id"`
	Name        string          `json:"name"`
	Description string          `json:"description"`
	Price       int             `json:"price"`
	Cost        int             `json:"cost"`
	Category    string          `json:"category"`
	Barcode     string          `json:"barcode"`
	Size        domain.ItemSize `json:"size"`
	ExpiryAt    time.Time       `json:"expiryAt"`
	CreatedAt   time.Time       `json:"createdAt"`
}

type GetItemResponse struct {
	ID          int             `json:"id"`
	Name        string          `json:"name"`
	Description string          `json:"description"`
	Price       int             `json:"price"`
	Cost        int             `json:"cost"`
	Category    string          `json:"category"`
	Barcode     string          `json:"barcode"`
	Size        domain.ItemSize `json:"size"`
	ExpiryAt    time.Time       `json:"expiryAt"`
	CreatedAt   time.Time       `json:"createdAt"`
}
