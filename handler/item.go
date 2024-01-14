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

func (h *ItemHandler) Delete(ginCtx *gin.Context) {
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

	if err := h.itemUsecase.Delete(ctx, &item.DeleteInput{User: user, ItemID: itemID}); err != nil {
		if errors.Is(err, domain.ErrItemNotFound) {
			ginhelper.Error(ginCtx, ginhelper.NewHTTPError(http.StatusNotFound, i18n.ItemNotFound, errors.WithStack(err)))
			return
		}

		ginhelper.Error(ginCtx, errors.WithStack(err))
		return
	}

	ginCtx.Status(http.StatusNoContent)
}

func (h *ItemHandler) Update(ginCtx *gin.Context) {
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

	var req UpdateItemRequest
	if err := ginCtx.BindJSON(&req); err != nil {
		ginhelper.Error(ginCtx, ginhelper.NewHTTPError(http.StatusBadRequest, i18n.InvalidRequest, errors.WithStack(err)))
		return
	}
	if err := valid.ValidateStruct(req); err != nil {
		ginhelper.Error(ginCtx, ginhelper.NewHTTPError(http.StatusBadRequest, i18n.InvalidRequest, errors.WithStack(err)))
		return
	}
	if !req.ShouldUpdate() {
		ginhelper.Error(ginCtx, ginhelper.NewHTTPError(http.StatusBadRequest, i18n.InvalidRequest, errors.WithStack(err)))
		return
	}

	if err := h.itemUsecase.Update(ctx, &item.UpdateInput{
		User:        user,
		ItemID:      itemID,
		Name:        req.Name,
		Description: req.Description,
		Price:       req.Price,
		Cost:        req.Cost,
		Category:    req.Category,
		Barcode:     req.Barcode,
		Size:        req.Size,
		ExpiryAt:    req.ExpiryAt,
	}); err != nil {
		if errors.Is(err, domain.ErrItemNotFound) {
			ginhelper.Error(ginCtx, ginhelper.NewHTTPError(http.StatusNotFound, i18n.ItemNotFound, errors.WithStack(err)))
			return
		}
		if errors.Is(err, domain.ErrItemAlreadyExists) {
			ginhelper.Error(ginCtx, ginhelper.NewHTTPError(http.StatusConflict, i18n.ItemAlreadyExists, errors.WithStack(err)))
			return
		}

		ginhelper.Error(ginCtx, errors.WithStack(err))
		return
	}

	ginCtx.Status(http.StatusNoContent)
}

func (h *ItemHandler) Find(ginCtx *gin.Context) {
	ctx := ginhelper.GetContext(ginCtx)

	// 1. 인증된 유저 확인
	user, ok := ctx.Value(domain.CtxKeyUser).(*domain.User)
	if !ok {
		ginhelper.Error(ginCtx, errors.New("unauthenticated request"))
		return
	}

	var req FindItemRequest
	if err := ginCtx.BindQuery(&req); err != nil {
		ginhelper.Error(ginCtx, ginhelper.NewHTTPError(http.StatusBadRequest, i18n.InvalidRequest, errors.WithStack(err)))
		return
	}

	findOutput, err := h.itemUsecase.Find(ctx, &item.FindInput{
		User:        user,
		Keyword:     req.Keyword,
		SearchAfter: req.SearchAfter,
	})
	if err != nil {
		ginhelper.Error(ginCtx, errors.WithStack(err))
		return
	}

	items := make([]GetItemResponse, len(findOutput.Items))
	for i := 0; i < len(findOutput.Items); i++ {
		items[i] = GetItemResponse{
			ID:          findOutput.Items[i].ID,
			Name:        findOutput.Items[i].Name,
			Description: findOutput.Items[i].Description,
			Price:       findOutput.Items[i].Price,
			Cost:        findOutput.Items[i].Cost,
			Category:    findOutput.Items[i].Category,
			Barcode:     findOutput.Items[i].Barcode,
			Size:        findOutput.Items[i].Size,
			ExpiryAt:    findOutput.Items[i].ExpiryAt,
			CreatedAt:   findOutput.Items[i].CreatedAt,
		}
	}

	ginhelper.Success(ginCtx, FindItemResponse{
		TotalCount:  findOutput.TotalCount,
		Items:       items,
		HasNext:     findOutput.HasNext,
		SearchAfter: findOutput.SearchAfter,
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

type UpdateItemRequest struct {
	Name        *string          `json:"name" validate:"omitempty,gte=1,lte=100"`
	Description *string          `json:"description" validate:"omitempty,required"`
	Price       *int             `json:"price" validate:"omitempty,gt=0"`
	Cost        *int             `json:"cost" validate:"omitempty,gt=0"`
	Category    *string          `json:"category" validate:"omitempty,required,gte=1,lte=100"`
	Barcode     *string          `json:"barcode" validate:"omitempty,required,gte=1,lte=100"`
	Size        *domain.ItemSize `json:"size" validate:"omitempty,required,oneof=small large"`
	ExpiryAt    *time.Time       `json:"expiryAt" validate:"omitempty,required"`
}

func (r *UpdateItemRequest) ShouldUpdate() bool {
	return !valid.IsNil(r.Name) ||
		!valid.IsNil(r.Description) ||
		!valid.IsNil(r.Price) ||
		!valid.IsNil(r.Cost) ||
		!valid.IsNil(r.Category) ||
		!valid.IsNil(r.Barcode) ||
		!valid.IsNil(r.Size) ||
		!valid.IsNil(r.ExpiryAt)
}

type FindItemRequest struct {
	Keyword     string `form:"keyword"`
	SearchAfter int    `form:"searchAfter"`
}

type FindItemResponse struct {
	TotalCount  int               `json:"totalCount"`
	Items       []GetItemResponse `json:"items"`
	HasNext     bool              `json:"hasNext"`
	SearchAfter int               `json:"searchAfter"`
}
