package domain

import (
	"testing"
	"time"

	"golang.org/x/crypto/bcrypt"

	"github.com/stretchr/testify/require"

	"github.com/brianvoe/gofakeit/v6"
)

func TestNewUser(t *testing.T) {
	now := time.Unix(time.Now().Unix(), 0)
	type args struct {
		phoneNumber string
		name        string
		password    string
		createdAt   time.Time
	}
	tests := []struct {
		name    string
		args    args
		wantErr bool
	}{
		{
			name: "OK",
			args: args{
				phoneNumber: "01012341234",
				name:        gofakeit.Name(),
				password:    gofakeit.Password(true, true, true, true, true, 72),
				createdAt:   now,
			},
			wantErr: false,
		},
		{
			name: "empty name",
			args: args{
				phoneNumber: "01012341234",
				name:        "",
				password:    gofakeit.Password(true, true, true, true, true, 10),
				createdAt:   now,
			},
			wantErr: false,
		},
		{
			name: "invalid phoneNumber",
			args: args{
				phoneNumber: gofakeit.LetterN(10),
				name:        "",
				password:    gofakeit.Password(true, true, true, true, true, 10),
				createdAt:   now,
			},
			wantErr: true,
		},
		{
			name: "invalid password",
			args: args{
				phoneNumber: "01012341234",
				name:        gofakeit.Name(),
				password:    gofakeit.Password(true, false, false, false, false, 10),
				createdAt:   now,
			},
			wantErr: true,
		},
		{
			name: "invalid password",
			args: args{
				phoneNumber: "01012341234",
				name:        gofakeit.Name(),
				password:    gofakeit.Password(true, true, true, true, true, 73),
				createdAt:   now,
			},
			wantErr: true,
		},
		{
			name: "zero createdAt",
			args: args{
				phoneNumber: "01012341234",
				name:        "",
				password:    gofakeit.Password(true, true, true, true, true, 10),
				createdAt:   time.Time{},
			},
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := NewUser(tt.args.phoneNumber, tt.args.password, tt.args.createdAt)
			if tt.wantErr {
				require.Error(t, err)
				require.Nil(t, got)
				t.Logf("error: %v", err)
			} else {
				require.NoError(t, err)
				require.NotNil(t, got)
				t.Log(got)
			}
		})
	}
}

func TestUser_Validate(t *testing.T) {
	now := time.Unix(time.Now().Unix(), 0)
	password := gofakeit.Password(true, true, true, true, true, 10)
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	require.NoError(t, err)

	type fields struct {
		ID          int
		Name        string
		PhoneNumber string
		Password    string
		CreatedAt   time.Time
	}
	tests := []struct {
		name    string
		fields  fields
		wantErr bool
	}{
		{
			name: "OK",
			fields: fields{
				ID:          gofakeit.Number(1, 100),
				Name:        gofakeit.Name(),
				PhoneNumber: "01012341234",
				Password:    string(hashedPassword),
				CreatedAt:   now,
			},
			wantErr: false,
		},
		{
			name: "zero UserID",
			fields: fields{
				ID:          0,
				Name:        gofakeit.Name(),
				PhoneNumber: "01012341234",
				Password:    string(hashedPassword),
				CreatedAt:   now,
			},
			wantErr: true,
		},
		{
			name: "invalid PhoneNumber",
			fields: fields{
				ID:          gofakeit.Number(1, 100),
				Name:        gofakeit.Name(),
				PhoneNumber: gofakeit.Name(),
				Password:    string(hashedPassword),
				CreatedAt:   now,
			},
			wantErr: true,
		},
		{
			name: "plain password",
			fields: fields{
				ID:          gofakeit.Number(1, 100),
				Name:        gofakeit.Name(),
				PhoneNumber: "01012341234",
				Password:    password,
				CreatedAt:   now,
			},
			wantErr: true,
		},
		{
			name: "zero CreatedAt",
			fields: fields{
				ID:          gofakeit.Number(1, 100),
				Name:        gofakeit.Name(),
				PhoneNumber: "01012341234",
				Password:    string(hashedPassword),
				CreatedAt:   time.Time{},
			},
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			u := &User{
				ID:          tt.fields.ID,
				PhoneNumber: tt.fields.PhoneNumber,
				Password:    tt.fields.Password,
				CreatedAt:   tt.fields.CreatedAt,
			}
			err := u.Validate()
			if tt.wantErr {
				require.Error(t, err)
			} else {
				require.NoError(t, err)
			}
		})
	}
}
