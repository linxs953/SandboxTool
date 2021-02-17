package utils

import (
	"encoding/json"
	"fmt"
	"strings"

	"github.com/dgrijalva/jwt-go"
	"github.com/rs/zerolog/log"
)

func GetTokenByUnionid(unionid string) string {
	if unionid == "" {
		log.Print("unionid is nil")
		return ""
	}
	var tokenStruct TokenByUnionid
	getTokenByUnionidAPI := GetStringEnv("TOKENBYUNIONIDAPI", "")
	secret := GetStringEnv("SECRETKEY", "secret")
	if getTokenByUnionidAPI == "" {
		log.Print("Get [TOKENBYUNIONIDAPI] env error")
		return ""
	}
	if secret == "" {
		log.Print("Get [SECRETKEY] env error")
		return ""
	}
	reqURL := fmt.Sprintf(getTokenByUnionidAPI, unionid)
	bodyByte := Get(reqURL, 200, false)
	if string(bodyByte) == "" {
		log.Print("occur error")
		return ""
	}
	err := json.Unmarshal(bodyByte, &tokenStruct)
	if err != nil {
		log.Error().Err(err).Msg("Parse json error")
		return ""
	}
	tokenString := tokenStruct.Data.JWT
	if tokenString == "" {
		log.Print("Get token from resp error")
		fmt.Println(reqURL)
		return ""
	}
	if strings.Contains(tokenString, "user not found") {
		log.Printf("Can not find user by unionid %s", unionid)
		return ""
	}
	token := UpdateJwt(tokenString, secret)
	if token == "" {
		log.Print("update token parm error")
		return ""
	}
	return token
}

func UpdateJwt(tokenString string, secret string) string {
	payload := DecodeToken(tokenString, secret)
	if payload == nil {
		log.Print("Decode token return nil")
		return ""
	}
	payload["iss"] = "forc"
	jwtToken := jwt.NewWithClaims(jwt.SigningMethodHS256, payload)
	token, err := jwtToken.SignedString([]byte(secret))
	if err != nil {
		log.Error().Err(err).Msg("Signed token string error")
		return ""
	}
	return token
}

func DecodeToken(tokenString string, secret string) jwt.MapClaims {
	hmacSampleSecret := []byte(secret)
	var token *jwt.Token
	var err error
	token, err = jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("Unexpected signing method: %v", token.Header["alg"])
		}
		return hmacSampleSecret, nil
	})
	if token == nil {
		log.Print("decode token get nil")
		return nil
	}
	if claims, ok := token.Claims.(jwt.MapClaims); ok && token.Valid {
		return claims
	} else {
		fmt.Printf("%v", err)
		return nil
	}
}

func GenerateToken(puid string, storageid string) string {
	secret := GetStringEnv("SECRETKEY", "")
	if secret == "" {
		log.Print("Get [SECRET] env error")
		return ""
	}
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"puid":       puid,
		"exp":        GetTimeStamp(),
		"app": "pandateacher.com",
		"iss": "passport.pandateacher.com",
		"path": "/",
		"scene": "pc_scan",
		"storage_id": storageid,
	})
	tokenString, err := token.SignedString([]byte(secret))
	if err != nil {
		log.Error().Err(err).Msg("jwt map to string error")
		return ""
	}
	return tokenString
}

func ParseToken(unionid string) interface{} {
	if unionid == "" {
		log.Print("unionid can not be empty")
		return nil
	}
	token := GetTokenByUnionid(unionid)
	secet := GetStringEnv("SECRETKEY", "")
	if secet == "" {
		log.Print("get  secret env error")
		return nil
	}
	if token == "" {
		log.Print("token is nil")
		return nil
	}
	jwtMap := DecodeToken(token, secet)
	if puid := jwtMap["puid"]; puid == "" {
		log.Print("Puid is nil")
		return nil
	} else {
		return puid
	}
}
