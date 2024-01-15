package item

import (
	"testing"

	"github.com/psi59/payhere-assignment/domain"

	"github.com/stretchr/testify/assert"

	"github.com/brianvoe/gofakeit/v6"
)

func TestUpdateInput_Validate(t *testing.T) {
	t.Run("OK", func(t *testing.T) {
		name := gofakeit.Name()
		var input UpdateInput
		input.ItemID = 1
		input.User = &domain.User{}
		input.Name = &name
		err := input.Validate()
		assert.NoError(t, err)
	})

	t.Run("empty name", func(t *testing.T) {
		var name string
		var input UpdateInput
		input.Name = &name
		err := input.Validate()
		assert.Error(t, err)
	})

	t.Run("invalid input", func(t *testing.T) {
		input := UpdateInput{
			User:   &domain.User{},
			ItemID: 1,
		}
		err := input.Validate()
		assert.Error(t, err)
	})
}
