package controller

import (
	"Slot_booking/cache"
	"Slot_booking/entity"
	"Slot_booking/service"
	"Slot_booking/utils"
	"errors"
	"fmt"
	"github.com/gin-gonic/gin"
	"net/http"
)

type UserController interface {
	GetUser(ctx *gin.Context) (entity.User, error, int)
	AddUser(ctx *gin.Context) (error, int)
	GetAllUsers() ([]entity.User, error)
	ChangeRoleToUser(ctx *gin.Context) error
	ChangeRoleToAdmin(ctx *gin.Context) error
	BlockUser(ctx *gin.Context) error
}

type userController struct {
	service   service.UserService
	userCache cache.UserCache
}

func NewUserController(service service.UserService, cache cache.UserCache) UserController {
	return &userController{
		service:   service,
		userCache: cache,
	}
}

func (c *userController) GetUser(ctx *gin.Context) (entity.User, error, int) {
	userReq := ctx.Value("user_info")
	jwtData := userReq.(*utils.JWTClaim)

	key := fmt.Sprintf("user_data_%v", jwtData.User.ID)
	user, err := c.userCache.GetUser(ctx, key)
	if err == nil {
		return user, err, http.StatusOK
	}
	user, err = c.service.GetUser(jwtData.User.ID)
	if err != nil {
		err = errors.New("invalid request")
		return entity.User{}, err, http.StatusBadRequest
	}
	if err != nil {
		return entity.User{}, err, http.StatusInternalServerError
	}

	c.userCache.SetUser(ctx, key, user)
	return user, err, http.StatusOK
}

func (c *userController) AddUser(ctx *gin.Context) (error, int) {
	var user entity.User

	if err := ctx.ShouldBindJSON(&user); err != nil {
		return err, http.StatusBadRequest
	}
	if err := user.HashPassword(user.Password); err != nil {
		err = errors.New("password could not be created")
		return err, http.StatusInternalServerError
	}

	user, err := c.service.AddUser(user)
	if err != nil {
		err = errors.New("user already exists")
		return err, http.StatusInternalServerError
	}
	key := fmt.Sprintf("user_data_%v", user.ID)
	c.userCache.SetUser(ctx, key, user)
	return nil, http.StatusOK
}

func (c *userController) GetAllUsers() ([]entity.User, error) {
	return c.service.GetAllUsers()
}

func (c *userController) ChangeRoleToUser(ctx *gin.Context) error {
	userReq := ctx.Value("user_info")
	jwtData := userReq.(*utils.JWTClaim)

	user, err := c.service.ChangeRoleToUser(jwtData.User)
	fmt.Println(user)
	if err != nil {
		//err = errors.New("the user has already been assigned the 'user' role")
		return err
	}
	return nil
}

func (c *userController) ChangeRoleToAdmin(ctx *gin.Context) error {
	userReq := ctx.Value("user_info")
	jwtData := userReq.(*utils.JWTClaim)

	user, err := c.service.ChangeRoleToAdmin(jwtData.User)
	fmt.Println(user)
	if err != nil {
		//err = errors.New("the user has already been assigned the 'admin' role")
		return err
	}
	return nil
}

func (c *userController) BlockUser(ctx *gin.Context) error {
	userReq := ctx.Value("user_info")
	jwtData := userReq.(*utils.JWTClaim)

	user, err := c.service.BlockUser(jwtData.User)
	fmt.Println(user)
	if err != nil {
		err = errors.New("the user has already been blocked")
		return err
	}
	return nil
}
