package valid

import (
	"testing"

	"github.com/brianvoe/gofakeit/v6"
	"github.com/stretchr/testify/require"
)

func TestValidatePhoneNumber(t *testing.T) {
	t.Run("01012341234", func(t *testing.T) {
		phoneNumber := gofakeit.Regex(`^01\d{8,9}$`)
		t.Logf("phoneNumber: %s", phoneNumber)
		require.NoError(t, ValidatePhoneNumber(phoneNumber))
	})

	t.Run("invalid phoneNumber: uuid", func(t *testing.T) {
		require.Error(t, ValidatePhoneNumber(gofakeit.UUID()))
	})

	t.Run("invalid phoneNumber: length", func(t *testing.T) {
		phoneNumber := gofakeit.Regex(`01\d{0,7}`)
		t.Logf("phoneNumber: %s", phoneNumber)
		require.Error(t, ValidatePhoneNumber(phoneNumber))

		phoneNumber = gofakeit.Regex(`01\d{10,15}`)
		t.Logf("phoneNumber: %s", phoneNumber)
		require.Error(t, ValidatePhoneNumber(phoneNumber))
	})
}

func TestValidatePassword(t *testing.T) {
	t.Run("OK", func(t *testing.T) {
		require.NoError(t, ValidatePassword("Test1!"))
	})

	t.Run("empty password", func(t *testing.T) {
		require.Error(t, ValidatePassword(""))
	})

	t.Run("Lower case not contains", func(t *testing.T) {
		pwd := gofakeit.Password(false, true, true, true, true, 10)
		err := ValidatePassword(pwd)
		require.Error(t, err)
	})

	t.Run("Upper case not contains", func(t *testing.T) {
		pwd := gofakeit.Password(true, false, true, true, true, 10)
		err := ValidatePassword(pwd)
		require.Error(t, err)
	})

	t.Run("numeric not contains", func(t *testing.T) {
		pwd := gofakeit.Password(true, true, false, true, true, 10)
		err := ValidatePassword(pwd)
		require.Error(t, err)
	})

	t.Run("Symbol not contains", func(t *testing.T) {
		pwd := gofakeit.Password(true, true, true, false, false, 10)
		err := ValidatePassword(pwd)
		require.Error(t, err)
	})
}
