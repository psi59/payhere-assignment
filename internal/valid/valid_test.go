package valid

import (
	"crypto/sha256"
	"encoding/base64"
	"io"
	"os"
	"testing"

	"github.com/stretchr/testify/require"
)

func TestName(t *testing.T) {
	f, err := os.Open("/Users/psi59/Downloads/IMG_250sqa.gif")
	require.NoError(t, err)

	h := sha256.New()
	_, err = io.Copy(h, f)
	require.NoError(t, err)
	b := h.Sum(nil)
	t.Log(base64.StdEncoding.EncodeToString(b))
}

func TestValidatePhoneNumber(t *testing.T) {
	t.Run("01012341234", func(t *testing.T) {
		require.NoError(t, ValidatePhoneNumber("01012341234"))
	})

	t.Run("0161231234", func(t *testing.T) {
		require.NoError(t, ValidatePhoneNumber("01612341234"))
	})

	t.Run("0212341111", func(t *testing.T) {
		require.Error(t, ValidatePhoneNumber("0212341111"))
	})
}

func TestValidatePassword(t *testing.T) {
	t.Run("OK", func(t *testing.T) {
		require.NoError(t, ValidatePassword("Test1!"))
	})

	t.Run("empty password", func(t *testing.T) {
		require.Error(t, ValidatePassword(""))
	})

	t.Run("numeric not contains", func(t *testing.T) {
		require.Error(t, ValidatePassword("Test!"))
	})

	t.Run("Upper case not contains", func(t *testing.T) {
		require.Error(t, ValidatePassword("test1!"))
	})

	t.Run("Lower case not contains", func(t *testing.T) {
		require.Error(t, ValidatePassword("TEST1!"))
	})

	t.Run("Symbol not contains", func(t *testing.T) {
		require.Error(t, ValidatePassword("test1"))
	})
}
