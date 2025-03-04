package auth

import (
	"fmt"
	"net/http"
	"testing"
	"time"

	"github.com/google/uuid"
)

func TestGetBearerToken(t *testing.T) {
	userID := uuid.New()
	tokenSecret := "SuperSecret"
	expiresIn := time.Second
	jwt, err := MakeJWT(userID, tokenSecret, expiresIn)
	if err != nil {
		t.Errorf("Error creating jwt: %s", err)
		return
	}

	cases := []struct {
		name   string
		header http.Header
		expErr bool
	}{
		{
			name: "No jwt after bearer",
			header: http.Header{
				"Authorization": []string{"Bearer "},
			},
			expErr: true,
		},
		{
			name: "Valid jwt",
			header: http.Header{
				"Authorization": []string{fmt.Sprintf("Bearer %s", jwt)},
			},
			expErr: false,
		},
	}

	for _, tt := range cases {
		t.Run(tt.name, func(t *testing.T) {
			_, err := GetBearerToken(tt.header)
			if (err != nil) != tt.expErr {
				t.Errorf("Invalid jwt found")
			}
		})
	}
}
