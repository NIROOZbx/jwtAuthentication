package utils

import (
	
	"time"

	"github.com/golang-jwt/jwt/v5"
)

 type Claims struct {
 	Userid int
 	jwt.RegisteredClaims
 }

 var SecretKey=[]byte("secret_token")

 func GenerateJwt(userid int) (string,error){


	expirationTime:=time.Now().Add(30*time.Minute)


	claims:=&Claims{
		Userid: userid,
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(expirationTime),
		},
	}


	token:=jwt.NewWithClaims(jwt.SigningMethodHS256,claims)

	fullTokenString,err:=token.SignedString(SecretKey)


return fullTokenString,err


 }



 