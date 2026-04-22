package user

import (
	"app/domain"
	"app/domain/model"
	"app/domain/model/auth"
	gorm_model "app/domain/model/gorm"
	request_model "app/domain/model/request"
	"app/domain/model/response"
	"app/helpers"
	jwt_helper "app/helpers/jsonwebtoken"
	"context"
	"database/sql"
	"errors"
	"net/http"
	"net/url"
	"time"

	exporter "github.com/faytranevozter/simple-exporter"
	"github.com/faytranevozter/simple-exporter/config"
	"github.com/google/uuid"
	"github.com/sirupsen/logrus"
	"golang.org/x/crypto/bcrypt"
	"gorm.io/gorm"
)

type UserRepository interface {
	StructScan(rows *sql.Rows, dest any) error
	FetchUser(ctx context.Context, options domain.UserFilter) (*sql.Rows, error)
	FetchOneUser(ctx context.Context, options domain.UserFilter) (*domain.User, error)
	CountUser(ctx context.Context, options domain.UserFilter) int64
	CreateUser(ctx context.Context, model *domain.User) (err error)
}

type Service struct {
	contextTimeout time.Duration
	userRepo       UserRepository
	// cacheRepo       CacheRepo
	// storageRepo     StorageRepository
}

func NewService(user UserRepository) *Service {
	return &Service{
		contextTimeout: time.Second * 10,
		userRepo:       user,
		// cacheRepo:       cache,
		// storageRepo:     storage,
	}
}

func (s *Service) Login(ctx context.Context, payload request_model.LoginRequest) (statusCode int, resp response.Base) {
	ctx, cancel := context.WithTimeout(ctx, s.contextTimeout)
	defer cancel()

	errValidation := make(map[string]string)
	// validating request
	if payload.Email == "" {
		errValidation["email"] = "email field is required"
	} else if !helpers.IsValidEmail(payload.Email) {
		errValidation["email"] = "email field is not valid"
	}

	if payload.Password == "" {
		errValidation["password"] = "password field is required"
	}

	if len(errValidation) > 0 {
		return http.StatusBadRequest, response.ErrorValidation(errValidation, "error validation")
	}

	// check the db
	user, err := s.userRepo.FetchOneUser(ctx, domain.UserFilter{
		Email: &payload.Email,
	})
	if err != nil && !errors.Is(err, gorm.ErrRecordNotFound) {
		return http.StatusInternalServerError, response.Error(500, err.Error())
	}

	if user == nil {
		return http.StatusBadRequest, response.Error(domain.ErrBadRequestCode, "user not found")
	}

	// check password
	if err = bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(payload.Password)); err != nil {
		return http.StatusBadRequest, response.Error(domain.ErrBadRequestCode, "Wrong password")
	}

	// generate token
	tokenString, err := jwt_helper.GenerateJWTToken(
		jwt_helper.GetJwtCredential().Member,
		auth.JWTClaimUser{
			UserID: user.ID.String(),
		},
	)
	if err != nil {
		return http.StatusBadRequest, response.Error(400, err.Error())
	}

	return http.StatusOK, response.Success(map[string]interface{}{
		"user":  user,
		"token": tokenString,
	})
}

func (s *Service) Register(ctx context.Context, payload request_model.RegisterRequest) (statusCode int, resp response.Base) {
	ctx, cancel := context.WithTimeout(ctx, s.contextTimeout)
	defer cancel()

	errValidation := make(map[string]string)
	// validating request
	if payload.Name == "" {
		errValidation["name"] = "name field is required"
	}

	if payload.Email == "" {
		errValidation["email"] = "email field is required"
	} else if !helpers.IsValidEmail(payload.Email) {
		errValidation["email"] = "email field is not valid"
	}

	if payload.Password == "" {
		errValidation["password"] = "password field is required"
	}

	if len(errValidation) > 0 {
		return http.StatusBadRequest, response.ErrorValidation(errValidation, "error validation")
	}

	// check the db
	user, err := s.userRepo.FetchOneUser(ctx, domain.UserFilter{
		Email: &payload.Email,
	})
	if err != nil && !errors.Is(err, gorm.ErrRecordNotFound) {
		return http.StatusInternalServerError, response.Error(500, err.Error())
	}

	helpers.Dump(user)

	if user != nil {
		return http.StatusBadRequest, response.Error(domain.ErrBadRequestCode, "username already taken")
	}

	hashedPassword, _ := bcrypt.GenerateFromPassword([]byte(payload.Password), bcrypt.DefaultCost)

	newUser := domain.User{
		ID:        uuid.New(),
		Name:      payload.Name,
		Email:     payload.Email,
		Password:  string(hashedPassword),
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
	}

	err = s.userRepo.CreateUser(ctx, &newUser)
	if err != nil {
		return http.StatusInternalServerError, response.Error(domain.ErrInternalServerCode, err.Error())
	}

	return http.StatusOK, response.Success(newUser)
}

func (s *Service) GetMe(ctx context.Context, claim auth.JWTClaimUser) (statusCode int, resp response.Base) {
	ctx, cancel := context.WithTimeout(ctx, s.contextTimeout)
	defer cancel()

	// parse uuid
	userID, _ := uuid.Parse(claim.UserID)

	// check the db
	user, err := s.userRepo.FetchOneUser(ctx, domain.UserFilter{
		ID: &userID,
	})
	if err != nil {
		return http.StatusInternalServerError, response.Error(domain.ErrInternalServerCode, err.Error())
	}

	if user == nil {
		return http.StatusBadRequest, response.Error(domain.ErrBadRequestCode, "user not found")
	}

	return http.StatusOK, response.Success(user)
}

func (s *Service) SampleUserList(ctx context.Context, claim auth.JWTClaimUser, urlQuery url.Values) (statusCode int, resp response.Base) {
	ctx, cancel := context.WithTimeout(ctx, s.contextTimeout)
	defer cancel()

	page, limit, offset := helpers.GetLimitOffset(urlQuery)

	options := domain.UserFilter{
		DefaultFilter: gorm_model.DefaultFilter{
			Limit:  &limit,
			Offset: &offset,
		},
	}

	if urlQuery.Get("email") != "" {
		email := urlQuery.Get("email")
		options.Email = &email
	}

	// ---------- createdAt filter ----------
	var start, end time.Time
	if t, e := time.Parse("2006-01-02 15:04:05", urlQuery.Get("created_at_start")); e == nil {
		start = t
	}

	if t, e := time.Parse("2006-01-02 15:04:05", urlQuery.Get("created_at_end")); e == nil {
		end = t
	}

	if !start.IsZero() && !end.IsZero() {
		options.CreatedAtRange = &model.DatetimeRange{
			Start: start,
			End:   end,
		}
	} else if !start.IsZero() {
		options.CreatedAtGte = &start
	} else if !end.IsZero() {
		options.CreatedAtLte = &end
	}
	// ---------- end createdAt filter ----------

	// count first
	totalDocuments := s.userRepo.CountUser(ctx, options)
	if totalDocuments == 0 {
		return http.StatusOK, response.Success(response.List{
			List:  []interface{}{},
			Page:  page,
			Limit: limit,
			Total: totalDocuments,
		})
	}

	// sorting here
	// options.Sorts = helpers.GetSorts(urlQuery, domain.UserAllowedSort)

	// check the db
	cur, err := s.userRepo.FetchUser(ctx, options)
	if err != nil {
		return http.StatusOK, response.Success(response.List{
			List:  []interface{}{},
			Page:  page,
			Limit: limit,
			Total: totalDocuments,
		})
	}
	defer cur.Close()

	list := make([]interface{}, 0)
	for cur.Next() {
		row := domain.User{}
		err := s.userRepo.StructScan(cur, &row)
		if err != nil {
			logrus.Error("User Decode ", err)
			return http.StatusOK, response.Success(response.List{
				List:  []interface{}{},
				Page:  page,
				Limit: limit,
				Total: totalDocuments,
			})
		}

		list = append(list, row)
	}

	return http.StatusOK, response.Success(response.List{
		List:  list,
		Page:  page,
		Limit: limit,
		Total: totalDocuments,
	})
}

func (s *Service) SampleUserDetail(ctx context.Context, claim auth.JWTClaimUser, id string) (statusCode int, resp response.Base) {
	ctx, cancel := context.WithTimeout(ctx, s.contextTimeout)
	defer cancel()

	userID, err := uuid.Parse(id)
	if err != nil {
		return http.StatusBadRequest, response.Error(domain.ErrBadRequestCode, "id is invalid")
	}

	// check the db
	user, err := s.userRepo.FetchOneUser(ctx, domain.UserFilter{
		ID: &userID,
	})
	if err != nil || user == nil {
		return http.StatusNotFound, response.Error(domain.ErrNotFoundCode, "User not found")
	}

	return http.StatusOK, response.Success(user)
}

func (s *Service) SampleUserExport(ctx context.Context, claim auth.JWTClaimUser, urlQuery url.Values) (statusCode int, resp response.Base) {
	ctx, cancel := context.WithTimeout(ctx, s.contextTimeout)
	defer cancel()

	options := domain.UserFilter{}

	// ---------- createdAt filter ----------
	var start, end time.Time
	if t, e := time.Parse("2006-01-02 15:04:05", urlQuery.Get("created_at_start")); e == nil {
		start = t
	}

	if t, e := time.Parse("2006-01-02 15:04:05", urlQuery.Get("created_at_end")); e == nil {
		end = t
	}

	if !start.IsZero() && !end.IsZero() {
		options.CreatedAtRange = &model.DatetimeRange{
			Start: start,
			End:   end,
		}
	} else if !start.IsZero() {
		options.CreatedAtGte = &start
	} else if !end.IsZero() {
		options.CreatedAtLte = &end
	}
	// ---------- end createdAt filter ----------

	// sorting here
	// options.Sorts = helpers.GetSorts(urlQuery, domain.UserAllowedSort)

	indo, _ := time.LoadLocation("Asia/Jakarta")

	// init the exporter
	exp := exporter.NewExporter(
		config.WithSheetHeader([]config.FieldConfig{
			{Key: "id", Label: "ID"},
			{Key: "name", Label: "Name"},
			{Key: "email", Label: "Email"},
			{Key: "created_at", Label: "Register At", As: "date", DateFormatLocation: indo},
		}),
		config.WithSheetFilter(true),
		config.WithSheetStyle(true),
	)

	// check the db
	cur, err := s.userRepo.FetchUser(ctx, options)
	if err != nil {
		return http.StatusInternalServerError, response.Error(domain.ErrInternalServerCode, err.Error())
	}
	defer cur.Close()

	for cur.Next() {
		row := domain.User{}
		err := s.userRepo.StructScan(cur, &row)
		if err != nil {
			logrus.Error("User Decode ", err)
			return http.StatusInternalServerError, response.Error(domain.ErrInternalServerCode, err.Error())
		}

		exp.AddRow(map[string]any{
			"id":         row.ID.String(),
			"name":       row.Name,
			"email":      row.Email,
			"created_at": row.CreatedAt,
			"updated_at": row.UpdatedAt,
		})
	}

	base64, _ := exp.ToBase64()

	return http.StatusOK, response.Success(map[string]any{
		"base64": base64,
	})
}
