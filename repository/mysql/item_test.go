package mysql

import (
	"context"
	"testing"
	"time"

	"github.com/psi59/payhere-assignment/repository"

	"github.com/jinzhu/copier"

	"github.com/brianvoe/gofakeit/v6"
	"github.com/psi59/payhere-assignment/domain"

	"github.com/psi59/payhere-assignment/internal/db"
	"github.com/stretchr/testify/assert"
)

func TestItemRepository_Create(t *testing.T) {
	ctx := db.ContextWithConn(context.TODO(), conn)
	userRepo := NewUserRepository()
	user := newTestUser(t)
	err := userRepo.Create(ctx, user)
	assert.NoError(t, err)
	itemRepo := NewItemRepository()

	t.Run("OK", func(t *testing.T) {
		item := newTestItem(t, user.ID)
		err := itemRepo.Create(ctx, item)
		assert.NoError(t, err)
		assert.True(t, item.ID > 0)
	})

	t.Run("nil context", func(t *testing.T) {
		item := newTestItem(t, user.ID)
		err := itemRepo.Create(nil, item)
		assert.Error(t, err)
	})

	t.Run("nil item", func(t *testing.T) {
		err := itemRepo.Create(ctx, nil)
		assert.Error(t, err)
	})

	t.Run("context without conn", func(t *testing.T) {
		item := newTestItem(t, user.ID)
		err := itemRepo.Create(context.TODO(), item)
		assert.Error(t, err)
	})

	t.Run("invalid item", func(t *testing.T) {
		item := newTestItem(t, user.ID)
		item.UserID = 0
		err := itemRepo.Create(ctx, item)
		assert.Error(t, err)
	})

	t.Run("중복 아이템", func(t *testing.T) {
		item := newTestItem(t, user.ID)
		var dupl domain.Item
		err := copier.Copy(&dupl, item)
		assert.NoError(t, err)

		err = itemRepo.Create(ctx, item)
		assert.NoError(t, err)
		assert.True(t, item.ID > 0)

		err = itemRepo.Create(ctx, &dupl)
		assert.ErrorIs(t, err, domain.ErrItemAlreadyExists)
	})
}

func TestItemRepository_Get(t *testing.T) {
	ctx := db.ContextWithConn(context.TODO(), conn)
	userRepo := NewUserRepository()
	user := newTestUser(t)
	err := userRepo.Create(ctx, user)
	assert.NoError(t, err)
	itemRepo := NewItemRepository()
	item := newTestItem(t, user.ID)
	err = itemRepo.Create(ctx, item)
	assert.NoError(t, err)

	t.Run("OK", func(t *testing.T) {
		got, err := itemRepo.Get(ctx, item.UserID, item.ID)
		assert.NoError(t, err)
		assert.Equal(t, item, got)
	})

	t.Run("nil context", func(t *testing.T) {
		got, err := itemRepo.Get(nil, item.UserID, item.ID)
		assert.Error(t, err)
		assert.Nil(t, got)
	})

	t.Run("invalid userID", func(t *testing.T) {
		got, err := itemRepo.Get(ctx, 0, item.ID)
		assert.Error(t, err)
		assert.Nil(t, got)
	})

	t.Run("invalid itemID", func(t *testing.T) {
		got, err := itemRepo.Get(ctx, item.UserID, 0)
		assert.Error(t, err)
		assert.Nil(t, got)
	})

	t.Run("context without conn", func(t *testing.T) {
		got, err := itemRepo.Get(context.TODO(), item.UserID, item.ID)
		assert.Error(t, err)
		assert.Nil(t, got)
	})

	t.Run("item not found", func(t *testing.T) {
		got, err := itemRepo.Get(ctx, gofakeit.Number(1000, 2000), item.ID)
		assert.Error(t, err, domain.ErrItemNotFound)
		assert.Nil(t, got)

		got, err = itemRepo.Get(ctx, item.UserID, gofakeit.Number(1000, 2000))
		assert.Error(t, err, domain.ErrItemNotFound)
		assert.Nil(t, got)
	})
}

func TestItemRepository_Delete(t *testing.T) {
	ctx := db.ContextWithConn(context.TODO(), conn)
	userRepo := NewUserRepository()
	user := newTestUser(t)
	err := userRepo.Create(ctx, user)
	assert.NoError(t, err)
	itemRepo := NewItemRepository()
	item := newTestItem(t, user.ID)
	err = itemRepo.Create(ctx, item)
	assert.NoError(t, err)

	t.Run("OK", func(t *testing.T) {
		err := itemRepo.Delete(ctx, item.UserID, item.ID)
		assert.NoError(t, err)
	})

	t.Run("nil context", func(t *testing.T) {
		err := itemRepo.Delete(nil, item.UserID, item.ID)
		assert.Error(t, err)

	})

	t.Run("invalid userID", func(t *testing.T) {
		err := itemRepo.Delete(ctx, 0, item.ID)
		assert.Error(t, err)

	})

	t.Run("invalid itemID", func(t *testing.T) {
		err := itemRepo.Delete(ctx, item.UserID, 0)
		assert.Error(t, err)

	})

	t.Run("context without conn", func(t *testing.T) {
		err := itemRepo.Delete(context.TODO(), item.UserID, item.ID)
		assert.Error(t, err)

	})
}

func TestItemRepository_Update(t *testing.T) {
	ctx := db.ContextWithConn(context.TODO(), conn)
	userRepo := NewUserRepository()
	user := newTestUser(t)
	err := userRepo.Create(ctx, user)
	assert.NoError(t, err)
	itemRepo := NewItemRepository()

	t.Run("OK", func(t *testing.T) {
		item := newTestItem(t, user.ID)
		err = itemRepo.Create(ctx, item)
		assert.NoError(t, err)

		name := gofakeit.Drink()
		description := gofakeit.SentenceSimple()
		price := gofakeit.Number(1000, 10000)
		cost := gofakeit.Number(1000, 10000)
		category := gofakeit.SentenceSimple()
		barcode := gofakeit.RandomString([]string{"coffee", "tea", "desert"})
		size := domain.ItemSizeLarge
		expiryAt := time.Unix(gofakeit.FutureDate().Unix(), 0).UTC()
		input := &repository.UpdateItemInput{
			Name:        &name,
			Description: &description,
			Price:       &price,
			Cost:        &cost,
			Category:    &category,
			Barcode:     &barcode,
			Size:        &size,
			ExpiryAt:    &expiryAt,
		}
		err := itemRepo.Update(ctx, item.UserID, item.ID, input)
		assert.NoError(t, err)

		var expected domain.Item
		err = copier.Copy(&expected, item)
		assert.NoError(t, err)
		got, err := itemRepo.Get(ctx, item.UserID, item.ID)
		assert.NoError(t, err)

		expected.Name = name
		expected.Description = description
		expected.Price = price
		expected.Cost = cost
		expected.Category = category
		expected.Barcode = barcode
		expected.Size = size
		expected.ExpiryAt = expiryAt

		assert.Equal(t, &expected, got)
	})

	t.Run("부분 업데이트", func(t *testing.T) {
		item := newTestItem(t, user.ID)
		err = itemRepo.Create(ctx, item)
		assert.NoError(t, err)

		name := gofakeit.Drink()
		expiryAt := time.Unix(gofakeit.FutureDate().Unix(), 0).UTC()
		input := &repository.UpdateItemInput{
			Name:     &name,
			ExpiryAt: &expiryAt,
		}
		err := itemRepo.Update(ctx, item.UserID, item.ID, input)
		assert.NoError(t, err)

		var expected domain.Item
		err = copier.Copy(&expected, item)
		assert.NoError(t, err)
		got, err := itemRepo.Get(ctx, item.UserID, item.ID)
		assert.NoError(t, err)

		expected.Name = name
		expected.ExpiryAt = expiryAt

		assert.Equal(t, &expected, got)
	})

	t.Run("nil context", func(t *testing.T) {
		item := newTestItem(t, user.ID)
		err = itemRepo.Create(ctx, item)
		assert.NoError(t, err)

		name := gofakeit.Drink()
		expiryAt := time.Unix(gofakeit.FutureDate().Unix(), 0).UTC()
		input := &repository.UpdateItemInput{
			Name:     &name,
			ExpiryAt: &expiryAt,
		}

		err = itemRepo.Update(nil, item.UserID, item.ID, input)
		assert.Error(t, err)

	})

	t.Run("invalid userID", func(t *testing.T) {
		item := newTestItem(t, user.ID)
		err = itemRepo.Create(ctx, item)
		assert.NoError(t, err)

		name := gofakeit.Drink()
		expiryAt := time.Unix(gofakeit.FutureDate().Unix(), 0).UTC()
		input := &repository.UpdateItemInput{
			Name:     &name,
			ExpiryAt: &expiryAt,
		}
		err := itemRepo.Update(ctx, 0, item.ID, input)
		assert.Error(t, err)
	})

	t.Run("invalid itemID", func(t *testing.T) {
		item := newTestItem(t, user.ID)
		err = itemRepo.Create(ctx, item)
		assert.NoError(t, err)

		name := gofakeit.Drink()
		expiryAt := time.Unix(gofakeit.FutureDate().Unix(), 0).UTC()
		input := &repository.UpdateItemInput{
			Name:     &name,
			ExpiryAt: &expiryAt,
		}
		err := itemRepo.Update(ctx, item.UserID, 0, input)
		assert.Error(t, err)
	})

	t.Run("nil input", func(t *testing.T) {
		item := newTestItem(t, user.ID)
		err = itemRepo.Create(ctx, item)
		assert.NoError(t, err)
		err := itemRepo.Update(ctx, item.UserID, item.ID, nil)
		assert.Error(t, err)
	})

	t.Run("invalid input", func(t *testing.T) {
		item := newTestItem(t, user.ID)
		err = itemRepo.Create(ctx, item)
		assert.NoError(t, err)
		err := itemRepo.Update(ctx, item.UserID, item.ID, &repository.UpdateItemInput{})
		assert.Error(t, err)
	})

	t.Run("context without conn", func(t *testing.T) {
		item := newTestItem(t, user.ID)
		err = itemRepo.Create(ctx, item)
		assert.NoError(t, err)

		name := gofakeit.Drink()
		expiryAt := time.Unix(gofakeit.FutureDate().Unix(), 0).UTC()
		input := &repository.UpdateItemInput{
			Name:     &name,
			ExpiryAt: &expiryAt,
		}
		err := itemRepo.Update(context.TODO(), item.UserID, item.ID, input)
		assert.Error(t, err)
	})

	t.Run("이름 중복", func(t *testing.T) {
		item := newTestItem(t, user.ID)
		err = itemRepo.Create(ctx, item)
		assert.NoError(t, err)

		item2 := newTestItem(t, user.ID)
		err = itemRepo.Create(ctx, item2)
		assert.NoError(t, err)

		expiryAt := time.Unix(gofakeit.FutureDate().Unix(), 0).UTC()
		input := &repository.UpdateItemInput{
			Name:     &item2.Name,
			ExpiryAt: &expiryAt,
		}
		err := itemRepo.Update(ctx, item.UserID, item.ID, input)
		assert.Error(t, err)
		assert.ErrorIs(t, err, domain.ErrItemAlreadyExists)
	})
}

func TestItemRepository_Find(t *testing.T) {
	ctx := db.ContextWithConn(context.TODO(), conn)
	userRepo := NewUserRepository()
	user := newTestUser(t)
	err := userRepo.Create(ctx, user)
	assert.NoError(t, err)
	itemRepo := NewItemRepository()

	for _, itemName := range []string{
		"슈크림 라떼",
		"카페 아메리카노",
		"카페 라떼",
		"사케라또 아포가토",
		"스파클링 시트러스 에스프레소",
		"클래식 아포가토",
		"사케라또 비안코 오버 아이스",
		"아이스 다크 초콜릿 모카",
		"아이스 바닐라 빈 라떼",
		"코르타도",
		"에스프레소",
		"아이스 리벤더 카페 브레베",
		"프렌치 애플 타르트 나이트로",
		"벨벳 다크 모카 나이트로",
		"리저브 나이트로",
		"콜드 브루 몰트",
		"콜드 브루 플로트",
		"리저브 콜드 브루",
		"아이스 에콰도르 로하",
		"아이스 선드라이드 브라질 아이피 에스테이트",
		"아이스 슬라웨시 토라자 사판 빌리지",
		"아이스 에이지드 수마트라 빈티지 2021",
		"아이스 코스타리카 나랑호",
		"콜드 브루 오트 라떼",
		"돌체 콜드 브루",
		"바닐라 크림 콜드 브루",
		"콜드 브루",
		"나이트로 바닐라 크림",
		"나이트로 콜드 브루",
		"제주 비자림 콜드 브루",
		"아이스 블론드 에스프레소 라떼",
		"아이스 블론드 바닐라 더블 샷 마키아또",
		"아이스 블론드 스타벅스 돌체 라떼",
		"아이스 블론드 카페 라떼",
		"아이스 블론드 카페 아메리카노",
		"아이스 별다방 바닐라 라떼",
		"바닐라 플랫 화이트",
		"아이스 스타벅스 돌체 라떼",
		"아이스 카페 모카",
		"아이스 카페 아메리카노",
		"아이스 카페 라떼",
		"아이스 카푸치노",
		"아이스 카라멜 마키아또",
		"아이스 화이트 초콜릿 모카",
		"커피 스타벅스 더블 샷",
		"바닐라 스타벅스 더블 샷",
		"헤이즐넛 스타벅스 더블샷",
		"에스프레스",
		"에스프레소 마키아또",
		"에스프레소 콘 파나",
		"제주 별다방 땅콩 라떼",
		"아이스 디카페인 스타벅스 돌체 라떼",
		"아이스 디카페인 카라멜 마키아또",
		"아이스 디카페인 카페 라떼",
		"아이스 디카페인 카페 아메리카노",
		"아이스 1/2 디카페인 스타벅스 돌체 라떼",
		"아이스 1/2디카페인 카라멜 마키아또",
		"아이스 1/2디카페인 카페 라떼",
		"아이스 1/2디카페인 카페 아메리카노",
		"돌체 카라멜 칩 커피 프라푸치노",
		"더블 에스프레소 칩 프라푸치노",
		"제주 유기농 말차로 만든 크림 프라푸치노",
		"자바 칩 프라푸치노",
		"화이트 딸기 크림 프라푸치노",
		"초콜릿 크림 칩 프라푸치노",
		"화이트 초콜릿 모카 프라푸치노",
		"모카 프라푸치노",
		"카라멜 프라푸치노",
		"에스프레소 프라푸치노",
		"바닐라 크림 프라푸치노",
		"제주 까망 크림 프라푸치노",
		"제주 쑥떡 크림 프라푸치노",
		"제주 별다아 땅콩 프라푸치노",
		"화이트 타이거 프라푸치노",
		"돌체 딸기 크림 프라푸치노",
		"트리플 초콜릿 칩 커피 프라푸치노",
		"트리플 초콜릿 칩 크림 프라푸치노",
		"딸기 레몬 블렌디드",
		"민트 초콜릿 칩 블렌디드",
		"딸기 딜라이트 요거트 블렌디드",
		"피치&레몬 블렌디드",
		"망고 바나나 블렌디드",
		"망고 패션 후르츠 블렌디드",
		"제주 천혜향 블랙 티 블렌디드",
		"쿨 라임 피지오",
		"블랙 티 레모네이드 피지오",
		"패션 탱고 티 레모네이드 피지오",
		"스타벅스 파인애플 선셋 아이스티",
		"아이스 패션 푸르트 티",
		"아이스 유자 민트 티",
		"아이스 돌체 블랙 밀크 티",
		"피치 젤리 아이스티",
		"아이스 제주 유기농 말차로 만든 라떼",
		"아이스 차이 티 라떼",
		"아이스 라임 패션 티",
		"아이스 자몽 허니 블랙티",
		"아이스 제주 유기 녹차",
		"아이스 잉글리쉬 브렉퍼스트 티",
		"아이스 얼 그레이 티",
		"아이스 유스베리 티",
		"아이스 히비스커스 블렌드 티",
		"아이스 민트 블렌드 티",
		"아이스 캐모마일 블렌드 티",
		"아이스 별궁 오미자 유스베리 티",
		"아이스 콩고물 블랙 밀크 티",
		"아이스 푸를 청귤 민트 티",
		"아이스 허니 얼 그레이 밀크 티",
		"아이스 피치 히비스커스 티",
		"오늘의 커피",
		"아이스 커피",
		"아이스 시그니처 초콜릿",
		"스팀 우유",
		"우유",
		"제주 쑥쑥 라떼",
		"아이스 제주 까망 라떼",
		"제주 청귤 레모네이드",
		"플러피 판다 아이스 초콜릿",
		"스타벅스 슬래머",
		"파이팅 청귤",
		"도와줘 흑흑",
		"퍼플베리 굿",
		"기운내라임",
		"한방에 쭉 감당",
		"햇사과 주스",
		"수박주스",
		"딸리주스",
		"망고주스",
		"케일&사과주스",
		"한라봉 주스",
		"토마토주스",
		"블루베리 요거트",
		"치아씨드 요거트",
	} {
		item := newTestItem(t, user.ID)
		item.Name = itemName
		err = itemRepo.Create(ctx, item)
		assert.NoError(t, err)
	}

	t.Run("OK", func(t *testing.T) {
		page1, err := itemRepo.Find(ctx, &repository.FindItemInput{
			UserID:      user.ID,
			Keyword:     "라떼",
			SearchAfter: 0,
		})
		assert.NoError(t, err)
		assert.Equal(t, 19, page1.TotalCount)
		assert.Equal(t, 10, len(page1.Items))

		page2, err := itemRepo.Find(ctx, &repository.FindItemInput{
			UserID:      user.ID,
			Keyword:     "라떼",
			SearchAfter: page1.SearchAfter,
		})
		assert.NoError(t, err)
		assert.Equal(t, 19, page2.TotalCount)
		assert.Equal(t, 9, len(page2.Items))
	})

	t.Run("nil context", func(t *testing.T) {
		got, err := itemRepo.Find(nil, &repository.FindItemInput{
			UserID:      user.ID,
			Keyword:     "라떼",
			SearchAfter: 0,
		})
		assert.Error(t, err)
		assert.Nil(t, got)
	})

	t.Run("nil input", func(t *testing.T) {
		got, err := itemRepo.Find(ctx, nil)
		assert.Error(t, err)
		assert.Nil(t, got)
	})

	t.Run("context without conn", func(t *testing.T) {
		got, err := itemRepo.Find(context.TODO(), &repository.FindItemInput{
			UserID:      user.ID,
			Keyword:     "라떼",
			SearchAfter: 0,
		})
		assert.Error(t, err)
		assert.Nil(t, got)
	})

	t.Run("invalid input", func(t *testing.T) {
		got, err := itemRepo.Find(ctx, &repository.FindItemInput{
			UserID:      0,
			Keyword:     "라떼",
			SearchAfter: 0,
		})
		assert.Error(t, err)
		assert.Nil(t, got)
	})
}

func newTestItem(t *testing.T, userID int) *domain.Item {
	item, err := domain.NewItem(
		userID,
		gofakeit.UUID(),
		gofakeit.SentenceSimple(),
		gofakeit.Number(5000, 10000),
		gofakeit.Number(3000, 5000),
		gofakeit.RandomString([]string{"coffee", "tea", "desert"}),
		gofakeit.Numerify("##################"),
		time.Unix(gofakeit.FutureDate().Unix(), 0).UTC(),
		domain.ItemSize(gofakeit.RandomString([]string{string(domain.ItemSizeSmall), string(domain.ItemSizeLarge)})),
	)
	assert.NoError(t, err)
	item.CreatedAt = time.Unix(time.Now().Unix(), 0).UTC()

	return item
}
