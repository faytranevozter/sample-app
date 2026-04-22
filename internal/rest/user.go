package rest

import (
	"app/domain/model/auth"
	request_model "app/domain/model/request"
	"app/domain/model/response"
	"app/internal/rest/middleware"
	"context"
	"net/http"
	"net/url"

	"github.com/gin-gonic/gin"
)

type UserService interface {
	Login(ctx context.Context, payload request_model.LoginRequest) (statusCode int, response response.Base)
	Register(ctx context.Context, payload request_model.RegisterRequest) (statusCode int, response response.Base)
	GetMe(ctx context.Context, claim auth.JWTClaimUser) (statusCode int, response response.Base)

	SampleUserList(ctx context.Context, claim auth.JWTClaimUser, query url.Values) (statusCode int, response response.Base)
	SampleUserDetail(ctx context.Context, claim auth.JWTClaimUser, id string) (statusCode int, response response.Base)
	SampleUserExport(ctx context.Context, claim auth.JWTClaimUser, query url.Values) (statusCode int, response response.Base)
}

type UserHandler struct {
	Service    UserService
	Route      *gin.RouterGroup
	Middleware middleware.Middleware
}

func NewUserHandler(route *gin.RouterGroup, svc UserService, mdl middleware.Middleware) {
	handler := &UserHandler{
		Service:    svc,
		Route:      route,
		Middleware: mdl,
	}

	authRoute := handler.Route.Group("/auth")
	authRoute.POST("/login", handler.Login)
	authRoute.POST("/register", handler.Register)
	authRoute.GET("/me", handler.Middleware.Auth(), handler.GetMe)

	// sample
	sampleRoute := handler.Route.Group("/sample/user")
	sampleRoute.GET("/list", handler.Middleware.Auth(), handler.UserList)
	sampleRoute.GET("/detail/:id", handler.Middleware.Auth(), handler.UserDetail)
	sampleRoute.GET("/export", handler.Middleware.Auth(), handler.UserExport)
}

func (h *UserHandler) Login(c *gin.Context) {
	ctx := c.Request.Context()

	payload := request_model.LoginRequest{}
	err := c.ShouldBindJSON(&payload)
	if err != nil {
		c.AbortWithStatusJSON(http.StatusBadRequest, response.Error(http.StatusBadRequest, "invalid json data"))
		return
	}

	statusCode, response := h.Service.Login(ctx, payload)
	c.JSON(statusCode, response)
}

func (h *UserHandler) Register(c *gin.Context) {
	ctx := c.Request.Context()

	payload := request_model.RegisterRequest{}
	err := c.ShouldBindJSON(&payload)
	if err != nil {
		c.AbortWithStatusJSON(http.StatusBadRequest, response.Error(http.StatusBadRequest, "invalid json data"))
		return
	}

	statusCode, response := h.Service.Register(ctx, payload)
	c.JSON(statusCode, response)
}

func (h *UserHandler) GetMe(c *gin.Context) {
	ctx := c.Request.Context()

	statusCode, response := h.Service.GetMe(ctx, c.MustGet("token_data").(auth.JWTClaimUser))
	c.JSON(statusCode, response)
}

func (h *UserHandler) UserList(c *gin.Context) {
	ctx := c.Request.Context()

	statusCode, response := h.Service.SampleUserList(ctx, c.MustGet("token_data").(auth.JWTClaimUser), c.Request.URL.Query())
	c.JSON(statusCode, response)
}

func (h *UserHandler) UserDetail(c *gin.Context) {
	ctx := c.Request.Context()

	statusCode, response := h.Service.SampleUserDetail(ctx, c.MustGet("token_data").(auth.JWTClaimUser), c.Param("id"))
	c.JSON(statusCode, response)
}

func (h *UserHandler) UserExport(c *gin.Context) {
	ctx := c.Request.Context()

	statusCode, response := h.Service.SampleUserExport(ctx, c.MustGet("token_data").(auth.JWTClaimUser), c.Request.URL.Query())
	c.JSON(statusCode, response)
}
