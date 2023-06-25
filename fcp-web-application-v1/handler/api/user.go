package api

import (
	"a21hc3NpZ25tZW50/model"
	"a21hc3NpZ25tZW50/service"
	"fmt"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt"
)

type UserAPI interface {
	Register(c *gin.Context)
	Login(c *gin.Context)
	GetUserTaskCategory(c *gin.Context)
}

type userAPI struct {
	userService service.UserService
}

func NewUserAPI(userService service.UserService) *userAPI {
	return &userAPI{userService}
}

func (u *userAPI) Register(c *gin.Context) {
	var user model.UserRegister

	if err := c.BindJSON(&user); err != nil {
		c.JSON(http.StatusBadRequest, model.NewErrorResponse("invalid decode json"))
		return
	}

	if user.Email == "" || user.Password == "" || user.Fullname == "" {
		c.JSON(http.StatusBadRequest, model.NewErrorResponse("register data is empty"))
		return
	}

	var recordUser = model.User{
		Fullname: user.Fullname,
		Email:    user.Email,
		Password: user.Password,
	}

	recordUser, err := u.userService.Register(&recordUser)
	if err != nil {
		c.JSON(http.StatusInternalServerError, model.NewErrorResponse("error internal server"))
		return
	}

	c.JSON(http.StatusCreated, model.NewSuccessResponse("register success"))
}

func (u *userAPI) Login(c *gin.Context) {
	var user model.UserLogin

	if err := c.BindJSON(&user); err != nil {
		c.JSON(http.StatusBadRequest, model.ErrorResponse{Error: "invalid decode json"})
		return
	}

	var recordUser = model.User{
		Email:    user.Email,
		Password: user.Password,
	}
	token, err := u.userService.Login(&recordUser)

	if err != nil {
		c.JSON(http.StatusBadRequest, model.ErrorResponse{Error: "error internal server"})
		return
	}

	expirationTime := time.Now().Add(5 * time.Minute)
	claim := model.Claims{
		Email: recordUser.Email,
		StandardClaims: jwt.StandardClaims{
			ExpiresAt: expirationTime.Unix(),
		},
	}

	tokenn := jwt.NewWithClaims(jwt.SigningMethodHS256, claim)

	tokenString, err := tokenn.SignedString(model.JwtKey)
	if err != nil {
		c.JSON(http.StatusBadRequest, model.ErrorResponse{Error: "error internal server"})
		return
	}

	// if err != nil {
	// 	cookie = "NotSet"
	c.SetCookie("session_token", tokenString, 3600, "/", "localhost", false, true)
	// }
	cookie, err := c.Cookie("session_token")

	fmt.Println(tokenString)
	fmt.Printf("Cookie value: %s \n", cookie)
	fmt.Println(token)

	c.JSON(http.StatusOK, gin.H{"user_id": recordUser.ID, "message": "login success"})
}

func (u *userAPI) GetUserTaskCategory(c *gin.Context) {
	cookie, _ := c.Cookie("session_token")

	if cookie == "" {
		c.JSON(http.StatusUnauthorized, model.ErrorResponse{Error: "unauthorized"})
		return
	}

	usertaskcat, err := u.userService.GetUserTaskCategory()

	if err != nil {
		c.JSON(http.StatusInternalServerError, model.ErrorResponse{Error: "error internal server"})
	}

	c.JSON(http.StatusOK, usertaskcat)
}
