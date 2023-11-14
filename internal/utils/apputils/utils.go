package apputils

import (
	"bytes"
	"encoding/json"
	"os"
	"time"

	"github.com/golang-jwt/jwt"
	"golang.org/x/crypto/bcrypt"

	"github.com/MehmetTalhaSeker/mts-blog-api/internal/model"
	"github.com/MehmetTalhaSeker/mts-blog-api/internal/utils/errorutils"
)

func InterfaceToStruct(in, out interface{}) error {
	buf := new(bytes.Buffer)

	err := json.NewEncoder(buf).Encode(in)
	if err != nil {
		return errorutils.New(errorutils.ErrJSONEncode, err)
	}

	err = json.NewDecoder(buf).Decode(out)
	if err != nil {
		return errorutils.New(errorutils.ErrJSONDecode, err)
	}

	return nil
}

func InterfaceUnmarshal(in, out interface{}) error {
	dbByte, err := json.Marshal(in)
	if err != nil {
		return errorutils.New(errorutils.ErrJSONMarshal, err)
	}

	err = json.Unmarshal(dbByte, out)
	if err != nil {
		return errorutils.New(errorutils.ErrJSONUnmarshal, err)
	}

	return nil
}

func CreateJWT(u *model.User) (string, error) {
	claims := &jwt.MapClaims{
		"expiresAt": time.Now().Add(48 * time.Hour),
		"email":     u.Email,
		"username":  u.Username,
		"role":      u.Role,
		"uid":       u.ID,
		"status":    u.Status,
	}

	secret := os.Getenv("JWT_SECRET")
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)

	return token.SignedString([]byte(secret))
}

func EncryptPassword(password string) (string, error) {
	cp, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		return "", errorutils.New(nil, err)
	}

	return string(cp), nil
}
