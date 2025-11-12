package controller

import (
	models "jwt-authentication/Models"
	"jwt-authentication/datatabase"
	"jwt-authentication/utils"
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"
)

type UserData struct{
	Username string
	Password string
}


func RegisterHandler(c *gin.Context) {
var user UserData



	if err:=c.ShouldBindJSON(&user); err!=nil{
		c.JSON(404, gin.H{"message":"cannot add data"})
		return
	}


	if strings.TrimSpace(user.Username)=="" ||  strings.TrimSpace(user.Password)==""{
		c.JSON(404,gin.H{"message":"Dont use empty string"})
		return
	}

	res:=datatabase.DB.Create(&models.User{Username: user.Username,Password: user.Password})

	if res.Error!=nil{
		c.JSON(401, gin.H{"message":"cannot add data to database"})
		return
	}

	c.JSON(200, gin.H{"message":"User registered Successfully"})



}

func LoginHandler(c *gin.Context) {
	var user UserData

	var dbData models.User

	if err:=c.ShouldBindJSON(&user); err!=nil{
		c.JSON(404, gin.H{"message":"cannot add data"})
		return
	}


	if strings.TrimSpace(user.Username)=="" ||  strings.TrimSpace(user.Password)==""{
		c.JSON(404,gin.H{"message":"Dont use empty string"})
		return
	}

	userInDb:=datatabase.DB.First(&dbData,"username = ?",user.Username)

	if userInDb.Error!=nil{
		c.JSON(404, gin.H{"message":"cannot find data"})
		return
	}

	if dbData.Password!=user.Password{
		c.JSON(404,gin.H{"message":"Invalid login credentials"})
		return
	}


	
	jwtToken,jwtErr:=utils.GenerateJwt(int(dbData.ID))
	
	if jwtErr!=nil{
		c.JSON(404, gin.H{"message":"No token found"})
		return
	}
	
	
	c.SetCookie("session",jwtToken, 5000, "/","localhost", false,true)
	c.JSON(200,gin.H{"message":"Login Successful"})


}

func Logout(c *gin.Context){
c.SetCookie("session", "", -1, "/", "localhost", false, true)
	c.JSON(http.StatusOK, gin.H{"message": "Logged out successfully"})
}

func HomeHandler(c *gin.Context){

	c.JSON(200,gin.H{"message":"Welcome to home"})

}