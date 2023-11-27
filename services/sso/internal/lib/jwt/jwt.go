package jwt

import (
	"github.com/AlexisPell/m11s_grpc/services/sso/internal/domain/models"
	"github.com/golang-jwt/jwt/v5"
	"time"
)

func NewToken(user models.User, app models.App, duration time.Duration) (string, error) {
	token := jwt.New(jwt.SigningMethodHS256)

	claims := token.Claims.(jwt.MapClaims)
	claims["uid"] = user.ID
	claims["email"] = user.Email
	claims["exp"] = time.Now().Add(duration).Unix()
	claims["app_id"] = app.Id

	// sign token
	tokenStr, err := token.SignedString([]byte(app.Secret))
	if err != nil {
		return "", err
	}

	return tokenStr, nil
}
