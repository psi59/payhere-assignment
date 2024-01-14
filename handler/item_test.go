package handler

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/psi59/payhere-assignment/internal/i18n"
	"golang.org/x/text/language"

	"github.com/psi59/payhere-assignment/internal/mocks/ucmocks"
	"github.com/psi59/payhere-assignment/usecase/item"
	"go.uber.org/mock/gomock"

	"github.com/psi59/payhere-assignment/internal/ginhelper"

	"github.com/brianvoe/gofakeit/v6"
	"github.com/psi59/payhere-assignment/domain"

	"github.com/gin-gonic/gin"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestNewItemHandler(t *testing.T) {
	itemUsecase := &item.Service{}
	type args struct {
		itemUsecase item.Usecase
	}
	tests := []struct {
		name    string
		args    args
		wantErr bool
	}{
		{
			name: "OK",
			args: args{
				itemUsecase: itemUsecase,
			},
			wantErr: false,
		},
		{
			name: "nil itemUsecase",
			args: args{
				itemUsecase: nil,
			},
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := NewItemHandler(tt.args.itemUsecase)
			if tt.wantErr {
				assert.Error(t, err)
				assert.Nil(t, got)
			} else {
				assert.NoError(t, err)
				assert.NotNil(t, got)
			}
		})
	}
}

func TestItemHandler_Create(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	itemUsecase := ucmocks.NewMockItemTokenUsecase(ctrl)
	r := gin.New()
	handler, err := NewItemHandler(itemUsecase)
	assert.NoError(t, err)
	assert.NotNil(t, handler)

	userDomain := newTestUser(t, gofakeit.Password(true, true, true, true, true, 10))
	r.POST("/", ginhelper.ContextMiddleware(), func(ginCtx *gin.Context) {
		ctx := ginhelper.GetContext(ginCtx)
		ctx = context.WithValue(ctx, domain.CtxKeyUser, userDomain)
		ginhelper.SetContext(ginCtx, ctx)
		ginCtx.Next()
	}, handler.Create)

	r.PUT("/", handler.Create)

	t.Run("OK", func(t *testing.T) {
		createItemRequest := &CreateItemRequest{
			Name:        gofakeit.Drink(),
			Description: gofakeit.SentenceSimple(),
			Price:       gofakeit.Number(1, 10000),
			Cost:        gofakeit.Number(1, 10000),
			Category:    gofakeit.RandomString([]string{"coffee", "tea", "desert"}),
			Barcode:     gofakeit.Numerify("################"),
			Size:        domain.ItemSizeSmall,
			ExpiryAt:    time.Unix(gofakeit.FutureDate().Unix(), 0).UTC(),
		}
		itemDomain := &domain.Item{
			ID:          gofakeit.Number(1, 10),
			Name:        createItemRequest.Name,
			Description: createItemRequest.Description,
			Price:       createItemRequest.Price,
			Cost:        createItemRequest.Cost,
			Category:    gofakeit.RandomString([]string{"coffee", "tea", "desert"}),
			Barcode:     createItemRequest.Barcode,
			ExpiryAt:    createItemRequest.ExpiryAt,
			Size:        createItemRequest.Size,
			CreatedAt:   time.Unix(gofakeit.FutureDate().Unix(), 0).UTC(),
		}

		itemUsecase.EXPECT().Create(gomock.Any(), &item.CreateInput{
			User:        userDomain,
			Name:        createItemRequest.Name,
			Description: createItemRequest.Description,
			Price:       createItemRequest.Price,
			Cost:        createItemRequest.Cost,
			Category:    createItemRequest.Category,
			Barcode:     createItemRequest.Barcode,
			Size:        createItemRequest.Size,
			ExpiryAt:    createItemRequest.ExpiryAt,
		}).Return(&item.CreateOutput{Item: itemDomain}, nil)

		responseWriter := httptest.NewRecorder()
		buf := bytes.NewBuffer(nil)
		err := json.NewEncoder(buf).Encode(createItemRequest)
		require.NoError(t, err)
		httpRequest, err := http.NewRequest(http.MethodPost, "/", buf)
		require.NoError(t, err)
		r.ServeHTTP(responseWriter, httpRequest)

		responseData := &CreateItemResponse{}
		resp := ginhelper.Response{Data: responseData}
		err = json.NewDecoder(responseWriter.Body).Decode(&resp)
		require.NoError(t, err)
		assert.Equal(t, http.StatusOK, responseWriter.Code)
		assert.Equal(t, &CreateItemResponse{
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
		}, resp.Data)
	})

	t.Run("binding error", func(t *testing.T) {
		responseWriter := httptest.NewRecorder()
		buf := bytes.NewBuffer(nil)
		err := json.NewEncoder(buf).Encode(map[string]any{
			"name": 123,
		})
		require.NoError(t, err)
		httpRequest, err := http.NewRequest(http.MethodPost, "/", buf)
		require.NoError(t, err)
		r.ServeHTTP(responseWriter, httpRequest)

		resp := ginhelper.Response{}
		err = json.NewDecoder(responseWriter.Body).Decode(&resp)
		require.NoError(t, err)
		assert.Equal(t, http.StatusBadRequest, responseWriter.Code)
		assert.Equal(t, http.StatusBadRequest, resp.Meta.Code)
		assert.Equal(t, i18n.T(language.English, i18n.InvalidRequest, nil), resp.Meta.Message)
	})

	t.Run("invalid request", func(t *testing.T) {
		createItemRequest := &CreateItemRequest{
			Name:        "",
			Description: gofakeit.SentenceSimple(),
			Price:       gofakeit.Number(1, 10000),
			Cost:        gofakeit.Number(1, 10000),
			Category:    gofakeit.RandomString([]string{"coffee", "tea", "desert"}),
			Barcode:     gofakeit.Numerify("################"),
			Size:        domain.ItemSizeSmall,
			ExpiryAt:    time.Unix(gofakeit.FutureDate().Unix(), 0).UTC(),
		}

		responseWriter := httptest.NewRecorder()
		buf := bytes.NewBuffer(nil)
		err := json.NewEncoder(buf).Encode(createItemRequest)
		require.NoError(t, err)
		httpRequest, err := http.NewRequest(http.MethodPost, "/", buf)
		require.NoError(t, err)
		r.ServeHTTP(responseWriter, httpRequest)

		responseData := &CreateItemResponse{}
		resp := ginhelper.Response{Data: responseData}
		err = json.NewDecoder(responseWriter.Body).Decode(&resp)
		require.NoError(t, err)
		assert.Equal(t, http.StatusBadRequest, responseWriter.Code)
		assert.Equal(t, http.StatusBadRequest, resp.Meta.Code)
		assert.Equal(t, i18n.T(language.English, i18n.InvalidRequest, nil), resp.Meta.Message)
	})

	t.Run("인증되지 않은 요청일 경우", func(t *testing.T) {
		createItemRequest := &CreateItemRequest{
			Name:        gofakeit.Drink(),
			Description: gofakeit.SentenceSimple(),
			Price:       gofakeit.Number(1, 10000),
			Cost:        gofakeit.Number(1, 10000),
			Category:    gofakeit.RandomString([]string{"coffee", "tea", "desert"}),
			Barcode:     gofakeit.Numerify("################"),
			Size:        domain.ItemSizeSmall,
			ExpiryAt:    time.Unix(gofakeit.FutureDate().Unix(), 0).UTC(),
		}

		responseWriter := httptest.NewRecorder()
		buf := bytes.NewBuffer(nil)
		err := json.NewEncoder(buf).Encode(createItemRequest)
		require.NoError(t, err)
		httpRequest, err := http.NewRequest(http.MethodPut, "/", buf)
		require.NoError(t, err)
		r.ServeHTTP(responseWriter, httpRequest)

		responseData := &CreateItemResponse{}
		resp := ginhelper.Response{Data: responseData}
		err = json.NewDecoder(responseWriter.Body).Decode(&resp)
		require.NoError(t, err)
		assert.Equal(t, http.StatusInternalServerError, responseWriter.Code)
		assert.Equal(t, http.StatusInternalServerError, resp.Meta.Code)
		assert.Equal(t, i18n.T(language.English, i18n.InternalError, nil), resp.Meta.Message)
	})

	t.Run("중복된 아이템일 경우", func(t *testing.T) {
		createItemRequest := &CreateItemRequest{
			Name:        gofakeit.Drink(),
			Description: gofakeit.SentenceSimple(),
			Price:       gofakeit.Number(1, 10000),
			Cost:        gofakeit.Number(1, 10000),
			Category:    gofakeit.RandomString([]string{"coffee", "tea", "desert"}),
			Barcode:     gofakeit.Numerify("################"),
			Size:        domain.ItemSizeSmall,
			ExpiryAt:    time.Unix(gofakeit.FutureDate().Unix(), 0).UTC(),
		}

		itemUsecase.EXPECT().Create(gomock.Any(), &item.CreateInput{
			User:        userDomain,
			Name:        createItemRequest.Name,
			Description: createItemRequest.Description,
			Price:       createItemRequest.Price,
			Cost:        createItemRequest.Cost,
			Category:    createItemRequest.Category,
			Barcode:     createItemRequest.Barcode,
			Size:        createItemRequest.Size,
			ExpiryAt:    createItemRequest.ExpiryAt,
		}).Return(nil, domain.ErrItemAlreadyExists)

		responseWriter := httptest.NewRecorder()
		buf := bytes.NewBuffer(nil)
		err := json.NewEncoder(buf).Encode(createItemRequest)
		require.NoError(t, err)
		httpRequest, err := http.NewRequest(http.MethodPost, "/", buf)
		require.NoError(t, err)
		r.ServeHTTP(responseWriter, httpRequest)

		responseData := &CreateItemResponse{}
		resp := ginhelper.Response{Data: responseData}
		err = json.NewDecoder(responseWriter.Body).Decode(&resp)
		require.NoError(t, err)
		assert.Equal(t, http.StatusConflict, responseWriter.Code)
		assert.Equal(t, http.StatusConflict, resp.Meta.Code)
		assert.Equal(t, i18n.T(language.English, i18n.ItemAlreadyExists, nil), resp.Meta.Message)
	})

	t.Run("아이템 usecase에서 알 수 없는 에러를 반환할 경우", func(t *testing.T) {
		createItemRequest := &CreateItemRequest{
			Name:        gofakeit.Drink(),
			Description: gofakeit.SentenceSimple(),
			Price:       gofakeit.Number(1, 10000),
			Cost:        gofakeit.Number(1, 10000),
			Category:    gofakeit.RandomString([]string{"coffee", "tea", "desert"}),
			Barcode:     gofakeit.Numerify("################"),
			Size:        domain.ItemSizeSmall,
			ExpiryAt:    time.Unix(gofakeit.FutureDate().Unix(), 0).UTC(),
		}

		itemUsecase.EXPECT().Create(gomock.Any(), &item.CreateInput{
			User:        userDomain,
			Name:        createItemRequest.Name,
			Description: createItemRequest.Description,
			Price:       createItemRequest.Price,
			Cost:        createItemRequest.Cost,
			Category:    createItemRequest.Category,
			Barcode:     createItemRequest.Barcode,
			Size:        createItemRequest.Size,
			ExpiryAt:    createItemRequest.ExpiryAt,
		}).Return(nil, gofakeit.Error())

		responseWriter := httptest.NewRecorder()
		buf := bytes.NewBuffer(nil)
		err := json.NewEncoder(buf).Encode(createItemRequest)
		require.NoError(t, err)
		httpRequest, err := http.NewRequest(http.MethodPost, "/", buf)
		require.NoError(t, err)
		r.ServeHTTP(responseWriter, httpRequest)

		responseData := &CreateItemResponse{}
		resp := ginhelper.Response{Data: responseData}
		err = json.NewDecoder(responseWriter.Body).Decode(&resp)
		require.NoError(t, err)
		assert.Equal(t, http.StatusInternalServerError, responseWriter.Code)
		assert.Equal(t, http.StatusInternalServerError, resp.Meta.Code)
		assert.Equal(t, i18n.T(language.English, i18n.InternalError, nil), resp.Meta.Message)
	})
}

func TestItemHandler_Get(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	itemUsecase := ucmocks.NewMockItemTokenUsecase(ctrl)
	r := gin.New()
	handler, err := NewItemHandler(itemUsecase)
	assert.NoError(t, err)
	assert.NotNil(t, handler)

	userDomain := newTestUser(t, gofakeit.Password(true, true, true, true, true, 10))
	r.GET("/items/:itemId", ginhelper.ContextMiddleware(), func(ginCtx *gin.Context) {
		ctx := ginhelper.GetContext(ginCtx)
		ctx = context.WithValue(ctx, domain.CtxKeyUser, userDomain)
		ginhelper.SetContext(ginCtx, ctx)
		ginCtx.Next()
	}, handler.Get)
	r.GET("/unauthorized/:itemId", handler.Get)

	t.Run("OK", func(t *testing.T) {
		itemDomain := newTestItem(t, userDomain.ID)
		itemUsecase.EXPECT().Get(gomock.Any(), &item.GetInput{
			User:   userDomain,
			ItemID: itemDomain.ID,
		}).Return(&item.GetOutput{Item: itemDomain}, nil)
		responseWriter := httptest.NewRecorder()
		httpRequest, err := http.NewRequest(http.MethodGet, fmt.Sprintf("/items/%d", itemDomain.ID), nil)
		require.NoError(t, err)
		r.ServeHTTP(responseWriter, httpRequest)

		responseData := &GetItemResponse{}
		resp := ginhelper.Response{Data: responseData}
		err = json.NewDecoder(responseWriter.Body).Decode(&resp)
		require.NoError(t, err)
		assert.Equal(t, http.StatusOK, responseWriter.Code)
		assert.Equal(t, &GetItemResponse{
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
		}, responseData)
	})

	t.Run("invalid itemID", func(t *testing.T) {
		responseWriter := httptest.NewRecorder()
		httpRequest, err := http.NewRequest(http.MethodGet, fmt.Sprintf("/items/%s", gofakeit.UUID()), nil)
		require.NoError(t, err)
		r.ServeHTTP(responseWriter, httpRequest)

		responseData := &GetItemResponse{}
		resp := ginhelper.Response{Data: responseData}
		err = json.NewDecoder(responseWriter.Body).Decode(&resp)
		require.NoError(t, err)
		assert.Equal(t, http.StatusNotFound, responseWriter.Code)
		assert.Equal(t, http.StatusNotFound, resp.Meta.Code)
		assert.Equal(t, i18n.T(language.English, i18n.ItemNotFound, nil), resp.Meta.Message)
	})

	t.Run("unauthorized", func(t *testing.T) {
		itemDomain := newTestItem(t, userDomain.ID)
		responseWriter := httptest.NewRecorder()
		httpRequest, err := http.NewRequest(http.MethodGet, fmt.Sprintf("/unauthorized/%d", itemDomain.ID), nil)
		require.NoError(t, err)
		r.ServeHTTP(responseWriter, httpRequest)

		responseData := &GetItemResponse{}
		resp := ginhelper.Response{Data: responseData}
		err = json.NewDecoder(responseWriter.Body).Decode(&resp)
		require.NoError(t, err)
		assert.Equal(t, http.StatusInternalServerError, responseWriter.Code)
		assert.Equal(t, http.StatusInternalServerError, resp.Meta.Code)
		assert.Equal(t, i18n.T(language.English, i18n.InternalError, nil), resp.Meta.Message)
	})

	t.Run("item not found", func(t *testing.T) {
		itemDomain := newTestItem(t, userDomain.ID)
		itemUsecase.EXPECT().Get(gomock.Any(), &item.GetInput{
			User:   userDomain,
			ItemID: itemDomain.ID,
		}).Return(nil, domain.ErrItemNotFound)
		responseWriter := httptest.NewRecorder()
		httpRequest, err := http.NewRequest(http.MethodGet, fmt.Sprintf("/items/%d", itemDomain.ID), nil)
		require.NoError(t, err)
		r.ServeHTTP(responseWriter, httpRequest)

		responseData := &GetItemResponse{}
		resp := ginhelper.Response{Data: responseData}
		err = json.NewDecoder(responseWriter.Body).Decode(&resp)
		require.NoError(t, err)
		assert.Equal(t, http.StatusNotFound, responseWriter.Code)
		assert.Equal(t, http.StatusNotFound, resp.Meta.Code)
		assert.Equal(t, i18n.T(language.English, i18n.ItemNotFound, nil), resp.Meta.Message)
	})

	t.Run("unexpected error", func(t *testing.T) {
		itemDomain := newTestItem(t, userDomain.ID)
		itemUsecase.EXPECT().Get(gomock.Any(), &item.GetInput{
			User:   userDomain,
			ItemID: itemDomain.ID,
		}).Return(nil, gofakeit.Error())
		responseWriter := httptest.NewRecorder()
		httpRequest, err := http.NewRequest(http.MethodGet, fmt.Sprintf("/items/%d", itemDomain.ID), nil)
		require.NoError(t, err)
		r.ServeHTTP(responseWriter, httpRequest)

		responseData := &GetItemResponse{}
		resp := ginhelper.Response{Data: responseData}
		err = json.NewDecoder(responseWriter.Body).Decode(&resp)
		require.NoError(t, err)
		assert.Equal(t, http.StatusInternalServerError, responseWriter.Code)
		assert.Equal(t, http.StatusInternalServerError, resp.Meta.Code)
		assert.Equal(t, i18n.T(language.English, i18n.InternalError, nil), resp.Meta.Message)
	})
}

func newTestItem(t *testing.T, userID int) *domain.Item {
	itemDomain, err := domain.NewItem(
		userID,
		gofakeit.Drink(),
		gofakeit.SentenceSimple(),
		gofakeit.Number(5000, 10000),
		gofakeit.Number(3000, 5000),
		gofakeit.RandomString([]string{"coffee", "tea", "desert"}),
		gofakeit.Numerify("##################"),
		time.Unix(gofakeit.FutureDate().Unix(), 0).UTC(),
		domain.ItemSize(gofakeit.RandomString([]string{string(domain.ItemSizeSmall), string(domain.ItemSizeLarge)})),
	)
	assert.NoError(t, err)
	itemDomain.CreatedAt = time.Unix(time.Now().Unix(), 0).UTC()
	itemDomain.ID = gofakeit.Number(1, 10000)

	return itemDomain
}
