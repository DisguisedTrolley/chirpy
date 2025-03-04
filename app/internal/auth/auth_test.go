package auth

import (
	"testing"
	"time"

	"github.com/charmbracelet/log"
	"github.com/google/uuid"
)

func TestMakeJWT(t *testing.T) {
	userId := uuid.New()
	tokenSecret := "SuperSecret"
	expiresIn := 1 * time.Second

	t.Run("Test case 1", func(t *testing.T) {
		jwt, err := MakeJWT(userId, tokenSecret, expiresIn)
		if err != nil {
			log.Error(err)
			t.Errorf("Error creating jwt")
			return
		}

		uid, err := ValidateJWT(jwt, tokenSecret)
		if err != nil {
			t.Errorf("Expected the jwt to be valid")
			return
		}

		if uid != userId {
			t.Errorf("Decoded user id and i/p user id don't match")
			return
		}

		time.Sleep(1 * time.Second)

		_, err = ValidateJWT(jwt, tokenSecret)
		if err == nil {
			log.Error(err)
			t.Errorf("Expected jwt to have expired")
		}
	})
}
