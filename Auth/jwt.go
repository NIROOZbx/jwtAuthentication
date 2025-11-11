package main

import (
	"fmt"
	"log"
	"net/http"
	"strings"

	"time"

	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v5"
)

var acessToken = []byte("my_secret_key")
var refreshToken =[]byte("new_refresh_token")

type Claims struct {
	Username string `json:"username"`
	jwt.RegisteredClaims
}

func loginHandler(c *gin.Context) {

	var User struct {
		Name     string `json:"user_name"`
		Password string `json:"password"`
	}
	if err := c.ShouldBindJSON(&User); err != nil {
		c.JSON(404, gin.H{"message": "Could not add login data"})
		return
	}
	if User.Name != "nirooz" || User.Password != "1234" {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid credentials"})
		return
	}
	sessionDuration := 168 * time.Hour
	expirationTime := time.Now().Add(1*time.Minute)

	claims := &Claims{
		Username: User.Name,
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(expirationTime),
		},
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	

	tokenString, tokenErr := token.SignedString(acessToken)

	if tokenErr != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Could not create token"})
		return
	}

	claims1:=&Claims{
		Username: User.Name,
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(168*time.Hour)),
		},
	}

	refreshTK:=jwt.NewWithClaims(jwt.SigningMethodHS256,claims1)

	fullRefToken,_:=refreshTK.SignedString(refreshToken)



	c.SetCookie("refreshToken", fullRefToken, int(sessionDuration.Seconds()), "/", "localhost", false, true)

	c.JSON(http.StatusOK, gin.H{
		"Acesstoken": tokenString,
		"RefreshToken":fullRefToken,
	})
}

func AuthMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {

		authHeader:=c.GetHeader("Authorization")

		if authHeader==""{
			c.JSON(http.StatusUnauthorized, gin.H{"message": "Authorization header is required"})
            c.Abort()
            return
		}

		tokenString := strings.TrimPrefix(authHeader, "Bearer ")
	
		if tokenString==authHeader{
			c.JSON(http.StatusUnauthorized, gin.H{"message": "Invalid token format"})
            c.Abort()
			return
		}

		claims:=&Claims{}

		token,err:=jwt.ParseWithClaims(tokenString,claims,func(t *jwt.Token)(interface{},error){
			return acessToken,nil
		} )

		if err!=nil || !token.Valid{
			c.JSON(http.StatusUnauthorized, gin.H{"message": "Unauthorized: Invalid token"})
            c.Abort()
            return
		}


		c.Set("username",claims.Username)
		c.Next()
	

	}
}

func homeHandler(c *gin.Context) {
	username := c.MustGet("username").(string)
	c.JSON(http.StatusOK, gin.H{"message": "Welcome " + username})
}

func refreshHandler(c *gin.Context){

	tokenString,err:=c.Cookie("refreshToken")

	if err!=nil{
	c.JSON(http.StatusUnauthorized, gin.H{"message": "Refresh token not found"})
        return
	}

	claims:=&Claims{}

	token,jwtErr:=jwt.ParseWithClaims(tokenString, claims,func(t *jwt.Token) (any, error) {
		return refreshToken, nil
	})

	if jwtErr!=nil || !token.Valid{
		c.SetCookie("refreshToken", "",-1, "/", "localhost", false, true)
		c.JSON(http.StatusUnauthorized, gin.H{"message": "Unauthorized: Invalid refresh token"})
        return
	}

	if time.Until(claims.ExpiresAt.Time)<0{
		c.JSON(http.StatusUnauthorized, gin.H{"message": "Refresh token expired"})
        return
	}

	accessExpiration:=time.Now().Add(15*time.Minute)

	newClaims:=&Claims{
		Username: claims.Username,
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(accessExpiration),
		},
	}

	accessToken:=jwt.NewWithClaims(jwt.SigningMethodHS256,newClaims)

	newAccessToken,_:=accessToken.SignedString(acessToken)

	c.JSON(http.StatusOK, gin.H{
        "accessToken": newAccessToken,
    })


}


func main() {

	router := gin.Default()

	router.POST("/login", loginHandler)

	router.GET("/home", AuthMiddleware(), homeHandler)

	router.POST("/refresh", refreshHandler)

	log.Println("Starting server on :8080")

	router.Run()

}
