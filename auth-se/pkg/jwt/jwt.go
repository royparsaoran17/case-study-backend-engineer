package jwt

import (
	"crypto/md5"
	"encoding/hex"
	"github.com/dgrijalva/jwt-go"
	"github.com/pkg/errors"
	"time"
)

type RequestJwt struct {
	ID        string    `json:"id"`
	JWTKey    string    `json:"jwt_key"`
	CreatedAt time.Time `json:"created_at"`
}

func CreateToken(data interface{}, req RequestJwt) (*string, error) {

	hash := md5.New()
	hash.Write([]byte(req.ID + "-" + req.CreatedAt.String()))
	jti := hex.EncodeToString(hash.Sum(nil))

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"exp":  time.Now().Add(time.Hour * 1).Unix(),
		"data": data,
		"sub":  req.ID,
		"jti":  jti,
	})

	tokenString, err := token.SignedString([]byte(req.JWTKey))
	if err != nil {
		return nil, errors.Wrap(err, "generate signed url")
	}

	return &tokenString, nil
}

func VerifyToken(tokenString string, jwtKey string) (*jwt.MapClaims, error) {
	token, err := jwt.Parse(tokenString, func(t *jwt.Token) (interface{}, error) {
		return []byte(jwtKey), nil
	})

	if err != nil {

		if ve, ok := err.(*jwt.ValidationError); ok {

			if ve.Errors&jwt.ValidationErrorMalformed != 0 {
				return nil, ErrTokenMalformed
			}

			if ve.Errors&(jwt.ValidationErrorExpired|jwt.ValidationErrorNotValidYet) != 0 {
				return nil, ErrTokenExpired
			}
		}

	}

	if token.Valid {
		if claims, ok := token.Claims.(jwt.MapClaims); ok && claims["data"] != nil {
			return &claims, nil
		}
		return nil, ErrTokenMalformed

	} else if ve, ok := err.(*jwt.ValidationError); ok {

		if ve.Errors&jwt.ValidationErrorMalformed != 0 {
			return nil, ErrTokenMalformed
		}

		if ve.Errors&(jwt.ValidationErrorExpired|jwt.ValidationErrorNotValidYet) != 0 {
			return nil, ErrTokenExpired
		}
	}

	return nil, ErrTokenMalformed

}
