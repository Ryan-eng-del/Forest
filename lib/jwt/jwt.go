package jwt

import (
	"errors"
	lib "go-gateway/lib/conf"
	"time"

	"github.com/golang-jwt/jwt/v5"
)


var (
	ErrTokenExpired    = errors.New("token 已过期")
	ErrTokenNotValidYet = errors.New("token 不在有效")
	ErrTokenMalformed   = errors.New("token 非法")
	ErrTokenInvalid     = errors.New("token 不可用")
)

type JWT struct {
	SignKey []byte
	ExpirePeriod time.Duration
}


var JWTInstance *JWT

type CustomClaims struct {
	UserId   uint
	jwt.RegisteredClaims
}

type AppClaims struct {
	AppId string
	jwt.RegisteredClaims
}

func NewJWT() (*JWT) {
	if JWTInstance == nil {
		JWTInstance = &JWT{
			SignKey: []byte(lib.SecretKey),
			ExpirePeriod: lib.TokenExpirePeriod,
		}
	}
	return JWTInstance
}


func (j *JWT) GenerateTokenWithAppID(appId string) (string, error) {
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, AppClaims{
		AppId: appId,
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(j.ExpirePeriod)),
			IssuedAt:  jwt.NewNumericDate(time.Now()),
			NotBefore: jwt.NewNumericDate(time.Now()),
			Issuer:    "go-gateway/2024/3/16",
			Subject:   "token generated",
		},
	})
	return token.SignedString(j.SignKey)
}


func (j *JWT) GenerateTokenWithUserID(useId uint) (string, error) {
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, CustomClaims{
		UserId: useId,
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(j.ExpirePeriod)),
			IssuedAt:  jwt.NewNumericDate(time.Now()),
			NotBefore: jwt.NewNumericDate(time.Now()),
			Issuer:    "go-gateway/2024/3/16",
			Subject:   "token generated",
		},
	})

	return token.SignedString(j.SignKey)
}

func (j *JWT) ParseJWT(tokenStr string) (*CustomClaims, error) {
	token, err := jwt.ParseWithClaims(tokenStr, &CustomClaims{}, func(token *jwt.Token) (interface{}, error) {
		return j.SignKey, nil
	})

	if err != nil {
		return nil, err
	}

	if token != nil {
		if claims, ok := token.Claims.(*CustomClaims); ok && token.Valid {
			return claims, nil
		}

	}
	return nil, ErrTokenInvalid
}



func (j *JWT) ParseAppJWT(tokenStr string) (*AppClaims, error) {
	token, err := jwt.ParseWithClaims(tokenStr, &AppClaims{}, func(token *jwt.Token) (interface{}, error) {
		return j.SignKey, nil
	})

	if err != nil {
		return nil, err
	}

	if token != nil {
		if claims, ok := token.Claims.(*AppClaims); ok && token.Valid {
			return claims, nil
		}

	}
	return nil, ErrTokenInvalid
}

func (j *JWT) RefreshJWT(tokenStr string) (string, error) {
	jwt.WithTimeFunc(func() time.Time {
		return time.Unix(0, 0)
	})

	token, err := jwt.ParseWithClaims(tokenStr, &CustomClaims{}, func(token *jwt.Token) (interface{}, error) {
		return j.SignKey, nil
	})

	if err != nil {
		panic(err)
	}

	if claims, ok := token.Claims.(*CustomClaims); ok && token.Valid {
		jwt.WithTimeFunc(func() time.Time {
			return time.Now()
		})
		claims.RegisteredClaims.ExpiresAt = jwt.NewNumericDate(time.Now().Add(7 * 24 * time.Hour))
		return j.GenerateTokenWithUserID(claims.UserId)
	}

	return "", ErrTokenInvalid

}
