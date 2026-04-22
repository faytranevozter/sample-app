package user

import (
	"app/domain"
	"app/domain/model/auth"
	request_model "app/domain/model/request"
	"context"
	"database/sql"
	"errors"
	"net/http"
	"net/url"
	"os"
	"testing"

	"github.com/google/uuid"
	"golang.org/x/crypto/bcrypt"
	"gorm.io/gorm"
)

// mockUserRepo is a manual mock for UserRepository interface
type mockUserRepo struct {
	structScanFn   func(rows *sql.Rows, dest any) error
	fetchUserFn    func(ctx context.Context, options domain.UserFilter) (*sql.Rows, error)
	fetchOneUserFn func(ctx context.Context, options domain.UserFilter) (*domain.User, error)
	countUserFn    func(ctx context.Context, options domain.UserFilter) int64
	createUserFn   func(ctx context.Context, model *domain.User) error
}

func (m *mockUserRepo) StructScan(rows *sql.Rows, dest any) error {
	if m.structScanFn != nil {
		return m.structScanFn(rows, dest)
	}
	return nil
}

func (m *mockUserRepo) FetchUser(ctx context.Context, options domain.UserFilter) (*sql.Rows, error) {
	if m.fetchUserFn != nil {
		return m.fetchUserFn(ctx, options)
	}
	return nil, nil
}

func (m *mockUserRepo) FetchOneUser(ctx context.Context, options domain.UserFilter) (*domain.User, error) {
	if m.fetchOneUserFn != nil {
		return m.fetchOneUserFn(ctx, options)
	}
	return nil, nil
}

func (m *mockUserRepo) CountUser(ctx context.Context, options domain.UserFilter) int64 {
	if m.countUserFn != nil {
		return m.countUserFn(ctx, options)
	}
	return 0
}

func (m *mockUserRepo) CreateUser(ctx context.Context, model *domain.User) error {
	if m.createUserFn != nil {
		return m.createUserFn(ctx, model)
	}
	return nil
}

func TestNewService(t *testing.T) {
	repo := &mockUserRepo{}
	svc := NewService(repo)

	if svc == nil {
		t.Fatal("NewService() returned nil")
	}
	if svc.userRepo == nil {
		t.Error("userRepo should not be nil")
	}
}

// --- Login Tests ---

func TestLogin_EmptyEmail(t *testing.T) {
	svc := NewService(&mockUserRepo{})
	statusCode, resp := svc.Login(context.Background(), request_model.LoginRequest{
		Email:    "",
		Password: "password",
	})

	if statusCode != http.StatusBadRequest {
		t.Errorf("statusCode = %d, want %d", statusCode, http.StatusBadRequest)
	}
	if resp.Validation["email"] != "email field is required" {
		t.Errorf("validation[email] = %q", resp.Validation["email"])
	}
}

func TestLogin_InvalidEmail(t *testing.T) {
	svc := NewService(&mockUserRepo{})
	statusCode, resp := svc.Login(context.Background(), request_model.LoginRequest{
		Email:    "invalid-email",
		Password: "password",
	})

	if statusCode != http.StatusBadRequest {
		t.Errorf("statusCode = %d, want %d", statusCode, http.StatusBadRequest)
	}
	if resp.Validation["email"] != "email field is not valid" {
		t.Errorf("validation[email] = %q", resp.Validation["email"])
	}
}

func TestLogin_EmptyPassword(t *testing.T) {
	svc := NewService(&mockUserRepo{})
	statusCode, resp := svc.Login(context.Background(), request_model.LoginRequest{
		Email:    "test@example.com",
		Password: "",
	})

	if statusCode != http.StatusBadRequest {
		t.Errorf("statusCode = %d, want %d", statusCode, http.StatusBadRequest)
	}
	if resp.Validation["password"] != "password field is required" {
		t.Errorf("validation[password] = %q", resp.Validation["password"])
	}
}

func TestLogin_AllFieldsEmpty(t *testing.T) {
	svc := NewService(&mockUserRepo{})
	statusCode, resp := svc.Login(context.Background(), request_model.LoginRequest{})

	if statusCode != http.StatusBadRequest {
		t.Errorf("statusCode = %d, want %d", statusCode, http.StatusBadRequest)
	}
	if resp.Validation["email"] == "" {
		t.Error("email validation should be set")
	}
	if resp.Validation["password"] == "" {
		t.Error("password validation should be set")
	}
}

func TestLogin_UserNotFound(t *testing.T) {
	repo := &mockUserRepo{
		fetchOneUserFn: func(ctx context.Context, options domain.UserFilter) (*domain.User, error) {
			return nil, gorm.ErrRecordNotFound
		},
	}
	svc := NewService(repo)
	statusCode, resp := svc.Login(context.Background(), request_model.LoginRequest{
		Email:    "test@example.com",
		Password: "password",
	})

	if statusCode != http.StatusBadRequest {
		t.Errorf("statusCode = %d, want %d", statusCode, http.StatusBadRequest)
	}
	if resp.Message != "user not found" {
		t.Errorf("Message = %q, want %q", resp.Message, "user not found")
	}
}

func TestLogin_DatabaseError(t *testing.T) {
	repo := &mockUserRepo{
		fetchOneUserFn: func(ctx context.Context, options domain.UserFilter) (*domain.User, error) {
			return nil, errors.New("database connection lost")
		},
	}
	svc := NewService(repo)
	statusCode, _ := svc.Login(context.Background(), request_model.LoginRequest{
		Email:    "test@example.com",
		Password: "password",
	})

	if statusCode != http.StatusInternalServerError {
		t.Errorf("statusCode = %d, want %d", statusCode, http.StatusInternalServerError)
	}
}

func TestLogin_WrongPassword(t *testing.T) {
	hashedPassword, _ := bcrypt.GenerateFromPassword([]byte("correct-password"), bcrypt.DefaultCost)
	user := &domain.User{
		ID:       uuid.New(),
		Email:    "test@example.com",
		Password: string(hashedPassword),
	}

	repo := &mockUserRepo{
		fetchOneUserFn: func(ctx context.Context, options domain.UserFilter) (*domain.User, error) {
			return user, nil
		},
	}
	svc := NewService(repo)
	statusCode, resp := svc.Login(context.Background(), request_model.LoginRequest{
		Email:    "test@example.com",
		Password: "wrong-password",
	})

	if statusCode != http.StatusBadRequest {
		t.Errorf("statusCode = %d, want %d", statusCode, http.StatusBadRequest)
	}
	if resp.Message != "Wrong password" {
		t.Errorf("Message = %q, want %q", resp.Message, "Wrong password")
	}
}

func TestLogin_Success(t *testing.T) {
	os.Setenv("JWT_MEMBER_SECRET_KEY", "test-secret-key-for-testing")
	os.Setenv("JWT_MEMBER_TTL", "60")
	defer func() {
		os.Unsetenv("JWT_MEMBER_SECRET_KEY")
		os.Unsetenv("JWT_MEMBER_TTL")
	}()

	hashedPassword, _ := bcrypt.GenerateFromPassword([]byte("correct-password"), bcrypt.DefaultCost)
	user := &domain.User{
		ID:       uuid.New(),
		Name:     "Test User",
		Email:    "test@example.com",
		Password: string(hashedPassword),
	}

	repo := &mockUserRepo{
		fetchOneUserFn: func(ctx context.Context, options domain.UserFilter) (*domain.User, error) {
			return user, nil
		},
	}
	svc := NewService(repo)
	statusCode, resp := svc.Login(context.Background(), request_model.LoginRequest{
		Email:    "test@example.com",
		Password: "correct-password",
	})

	if statusCode != http.StatusOK {
		t.Errorf("statusCode = %d, want %d", statusCode, http.StatusOK)
	}
	if resp.ErrorCode != 0 {
		t.Errorf("ErrorCode = %d, want 0", resp.ErrorCode)
	}

	data, ok := resp.Data.(map[string]interface{})
	if !ok {
		t.Fatalf("Data type = %T, want map[string]interface{}", resp.Data)
	}
	if data["token"] == nil || data["token"] == "" {
		t.Error("token should not be empty")
	}
	if data["user"] == nil {
		t.Error("user should not be nil")
	}
}

// --- Register Tests ---

func TestRegister_EmptyName(t *testing.T) {
	svc := NewService(&mockUserRepo{})
	statusCode, resp := svc.Register(context.Background(), request_model.RegisterRequest{
		Name:     "",
		Email:    "test@example.com",
		Password: "password",
	})

	if statusCode != http.StatusBadRequest {
		t.Errorf("statusCode = %d, want %d", statusCode, http.StatusBadRequest)
	}
	if resp.Validation["name"] != "name field is required" {
		t.Errorf("validation[name] = %q", resp.Validation["name"])
	}
}

func TestRegister_EmptyEmail(t *testing.T) {
	svc := NewService(&mockUserRepo{})
	statusCode, resp := svc.Register(context.Background(), request_model.RegisterRequest{
		Name:     "Test",
		Email:    "",
		Password: "password",
	})

	if statusCode != http.StatusBadRequest {
		t.Errorf("statusCode = %d, want %d", statusCode, http.StatusBadRequest)
	}
	if resp.Validation["email"] != "email field is required" {
		t.Errorf("validation[email] = %q", resp.Validation["email"])
	}
}

func TestRegister_InvalidEmail(t *testing.T) {
	svc := NewService(&mockUserRepo{})
	statusCode, resp := svc.Register(context.Background(), request_model.RegisterRequest{
		Name:     "Test",
		Email:    "not-an-email",
		Password: "password",
	})

	if statusCode != http.StatusBadRequest {
		t.Errorf("statusCode = %d, want %d", statusCode, http.StatusBadRequest)
	}
	if resp.Validation["email"] != "email field is not valid" {
		t.Errorf("validation[email] = %q", resp.Validation["email"])
	}
}

func TestRegister_EmptyPassword(t *testing.T) {
	svc := NewService(&mockUserRepo{})
	statusCode, resp := svc.Register(context.Background(), request_model.RegisterRequest{
		Name:     "Test",
		Email:    "test@example.com",
		Password: "",
	})

	if statusCode != http.StatusBadRequest {
		t.Errorf("statusCode = %d, want %d", statusCode, http.StatusBadRequest)
	}
	if resp.Validation["password"] != "password field is required" {
		t.Errorf("validation[password] = %q", resp.Validation["password"])
	}
}

func TestRegister_AllFieldsEmpty(t *testing.T) {
	svc := NewService(&mockUserRepo{})
	statusCode, resp := svc.Register(context.Background(), request_model.RegisterRequest{})

	if statusCode != http.StatusBadRequest {
		t.Errorf("statusCode = %d, want %d", statusCode, http.StatusBadRequest)
	}
	if resp.Validation["name"] == "" {
		t.Error("name validation should be set")
	}
	if resp.Validation["email"] == "" {
		t.Error("email validation should be set")
	}
	if resp.Validation["password"] == "" {
		t.Error("password validation should be set")
	}
}

func TestRegister_UserAlreadyExists(t *testing.T) {
	existingUser := &domain.User{
		ID:    uuid.New(),
		Email: "test@example.com",
	}
	repo := &mockUserRepo{
		fetchOneUserFn: func(ctx context.Context, options domain.UserFilter) (*domain.User, error) {
			return existingUser, nil
		},
	}
	svc := NewService(repo)
	statusCode, resp := svc.Register(context.Background(), request_model.RegisterRequest{
		Name:     "Test",
		Email:    "test@example.com",
		Password: "password",
	})

	if statusCode != http.StatusBadRequest {
		t.Errorf("statusCode = %d, want %d", statusCode, http.StatusBadRequest)
	}
	if resp.Message != "username already taken" {
		t.Errorf("Message = %q, want %q", resp.Message, "username already taken")
	}
}

func TestRegister_DatabaseErrorOnFetch(t *testing.T) {
	repo := &mockUserRepo{
		fetchOneUserFn: func(ctx context.Context, options domain.UserFilter) (*domain.User, error) {
			return nil, errors.New("db error")
		},
	}
	svc := NewService(repo)
	statusCode, _ := svc.Register(context.Background(), request_model.RegisterRequest{
		Name:     "Test",
		Email:    "test@example.com",
		Password: "password",
	})

	if statusCode != http.StatusInternalServerError {
		t.Errorf("statusCode = %d, want %d", statusCode, http.StatusInternalServerError)
	}
}

func TestRegister_CreateError(t *testing.T) {
	repo := &mockUserRepo{
		fetchOneUserFn: func(ctx context.Context, options domain.UserFilter) (*domain.User, error) {
			return nil, gorm.ErrRecordNotFound
		},
		createUserFn: func(ctx context.Context, model *domain.User) error {
			return errors.New("create failed")
		},
	}
	svc := NewService(repo)
	statusCode, resp := svc.Register(context.Background(), request_model.RegisterRequest{
		Name:     "Test User",
		Email:    "newuser@example.com",
		Password: "password123",
	})

	if statusCode != http.StatusInternalServerError {
		t.Errorf("statusCode = %d, want %d", statusCode, http.StatusInternalServerError)
	}
	if resp.ErrorCode != domain.ErrInternalServerCode {
		t.Errorf("ErrorCode = %d, want %d", resp.ErrorCode, domain.ErrInternalServerCode)
	}
}

func TestRegister_Success(t *testing.T) {
	repo := &mockUserRepo{
		fetchOneUserFn: func(ctx context.Context, options domain.UserFilter) (*domain.User, error) {
			return nil, gorm.ErrRecordNotFound
		},
		createUserFn: func(ctx context.Context, model *domain.User) error {
			return nil
		},
	}
	svc := NewService(repo)
	statusCode, resp := svc.Register(context.Background(), request_model.RegisterRequest{
		Name:     "New User",
		Email:    "newuser@example.com",
		Password: "password123",
	})

	if statusCode != http.StatusOK {
		t.Errorf("statusCode = %d, want %d", statusCode, http.StatusOK)
	}
	if resp.ErrorCode != 0 {
		t.Errorf("ErrorCode = %d, want 0", resp.ErrorCode)
	}
	if resp.Data == nil {
		t.Error("Data should not be nil")
	}
}

// --- GetMe Tests ---

func TestGetMe_Success(t *testing.T) {
	userID := uuid.New()
	expectedUser := &domain.User{
		ID:    userID,
		Name:  "Test User",
		Email: "test@example.com",
	}

	repo := &mockUserRepo{
		fetchOneUserFn: func(ctx context.Context, options domain.UserFilter) (*domain.User, error) {
			return expectedUser, nil
		},
	}
	svc := NewService(repo)
	statusCode, resp := svc.GetMe(context.Background(), auth.JWTClaimUser{
		UserID: userID.String(),
	})

	if statusCode != http.StatusOK {
		t.Errorf("statusCode = %d, want %d", statusCode, http.StatusOK)
	}
	if resp.ErrorCode != 0 {
		t.Errorf("ErrorCode = %d, want 0", resp.ErrorCode)
	}
}

func TestGetMe_UserNotFound(t *testing.T) {
	repo := &mockUserRepo{
		fetchOneUserFn: func(ctx context.Context, options domain.UserFilter) (*domain.User, error) {
			return nil, nil
		},
	}
	svc := NewService(repo)
	statusCode, resp := svc.GetMe(context.Background(), auth.JWTClaimUser{
		UserID: uuid.New().String(),
	})

	if statusCode != http.StatusBadRequest {
		t.Errorf("statusCode = %d, want %d", statusCode, http.StatusBadRequest)
	}
	if resp.Message != "user not found" {
		t.Errorf("Message = %q, want %q", resp.Message, "user not found")
	}
}

func TestGetMe_DatabaseError(t *testing.T) {
	repo := &mockUserRepo{
		fetchOneUserFn: func(ctx context.Context, options domain.UserFilter) (*domain.User, error) {
			return nil, errors.New("db error")
		},
	}
	svc := NewService(repo)
	statusCode, _ := svc.GetMe(context.Background(), auth.JWTClaimUser{
		UserID: uuid.New().String(),
	})

	if statusCode != http.StatusInternalServerError {
		t.Errorf("statusCode = %d, want %d", statusCode, http.StatusInternalServerError)
	}
}

// --- SampleUserDetail Tests ---

func TestSampleUserDetail_InvalidID(t *testing.T) {
	svc := NewService(&mockUserRepo{})
	statusCode, resp := svc.SampleUserDetail(context.Background(), auth.JWTClaimUser{}, "invalid-uuid")

	if statusCode != http.StatusBadRequest {
		t.Errorf("statusCode = %d, want %d", statusCode, http.StatusBadRequest)
	}
	if resp.Message != "id is invalid" {
		t.Errorf("Message = %q, want %q", resp.Message, "id is invalid")
	}
}

func TestSampleUserDetail_UserNotFound(t *testing.T) {
	repo := &mockUserRepo{
		fetchOneUserFn: func(ctx context.Context, options domain.UserFilter) (*domain.User, error) {
			return nil, gorm.ErrRecordNotFound
		},
	}
	svc := NewService(repo)
	statusCode, resp := svc.SampleUserDetail(context.Background(), auth.JWTClaimUser{}, uuid.New().String())

	if statusCode != http.StatusNotFound {
		t.Errorf("statusCode = %d, want %d", statusCode, http.StatusNotFound)
	}
	if resp.Message != "User not found" {
		t.Errorf("Message = %q, want %q", resp.Message, "User not found")
	}
}

func TestSampleUserDetail_Success(t *testing.T) {
	userID := uuid.New()
	expectedUser := &domain.User{
		ID:    userID,
		Name:  "Test User",
		Email: "test@example.com",
	}

	repo := &mockUserRepo{
		fetchOneUserFn: func(ctx context.Context, options domain.UserFilter) (*domain.User, error) {
			return expectedUser, nil
		},
	}
	svc := NewService(repo)
	statusCode, resp := svc.SampleUserDetail(context.Background(), auth.JWTClaimUser{}, userID.String())

	if statusCode != http.StatusOK {
		t.Errorf("statusCode = %d, want %d", statusCode, http.StatusOK)
	}
	if resp.ErrorCode != 0 {
		t.Errorf("ErrorCode = %d, want 0", resp.ErrorCode)
	}
}

// --- SampleUserList Tests ---

func TestSampleUserList_EmptyResult(t *testing.T) {
	repo := &mockUserRepo{
		countUserFn: func(ctx context.Context, options domain.UserFilter) int64 {
			return 0
		},
	}
	svc := NewService(repo)
	statusCode, resp := svc.SampleUserList(context.Background(), auth.JWTClaimUser{}, url.Values{})

	if statusCode != http.StatusOK {
		t.Errorf("statusCode = %d, want %d", statusCode, http.StatusOK)
	}
	if resp.ErrorCode != 0 {
		t.Errorf("ErrorCode = %d, want 0", resp.ErrorCode)
	}
}

func TestSampleUserList_FetchError(t *testing.T) {
	repo := &mockUserRepo{
		countUserFn: func(ctx context.Context, options domain.UserFilter) int64 {
			return 5
		},
		fetchUserFn: func(ctx context.Context, options domain.UserFilter) (*sql.Rows, error) {
			return nil, errors.New("fetch error")
		},
	}
	svc := NewService(repo)
	statusCode, resp := svc.SampleUserList(context.Background(), auth.JWTClaimUser{}, url.Values{})

	if statusCode != http.StatusOK {
		t.Errorf("statusCode = %d, want %d", statusCode, http.StatusOK)
	}
	if resp.ErrorCode != 0 {
		t.Errorf("ErrorCode = %d, want 0", resp.ErrorCode)
	}
}

func TestSampleUserList_WithEmailFilter(t *testing.T) {
	var capturedOptions domain.UserFilter
	repo := &mockUserRepo{
		countUserFn: func(ctx context.Context, options domain.UserFilter) int64 {
			capturedOptions = options
			return 0
		},
	}
	svc := NewService(repo)
	query := url.Values{"email": []string{"test@example.com"}}
	svc.SampleUserList(context.Background(), auth.JWTClaimUser{}, query)

	if capturedOptions.Email == nil || *capturedOptions.Email != "test@example.com" {
		t.Error("email filter should be set")
	}
}

func TestSampleUserList_WithDateRangeFilter(t *testing.T) {
	var capturedOptions domain.UserFilter
	repo := &mockUserRepo{
		countUserFn: func(ctx context.Context, options domain.UserFilter) int64 {
			capturedOptions = options
			return 0
		},
	}
	svc := NewService(repo)
	query := url.Values{
		"created_at_start": []string{"2024-01-01 00:00:00"},
		"created_at_end":   []string{"2024-12-31 23:59:59"},
	}
	svc.SampleUserList(context.Background(), auth.JWTClaimUser{}, query)

	if capturedOptions.CreatedAtRange == nil {
		t.Error("CreatedAtRange should be set when both start and end are provided")
	}
}

func TestSampleUserList_WithOnlyStartDate(t *testing.T) {
	var capturedOptions domain.UserFilter
	repo := &mockUserRepo{
		countUserFn: func(ctx context.Context, options domain.UserFilter) int64 {
			capturedOptions = options
			return 0
		},
	}
	svc := NewService(repo)
	query := url.Values{
		"created_at_start": []string{"2024-01-01 00:00:00"},
	}
	svc.SampleUserList(context.Background(), auth.JWTClaimUser{}, query)

	if capturedOptions.CreatedAtGte == nil {
		t.Error("CreatedAtGte should be set when only start is provided")
	}
}

func TestSampleUserList_WithOnlyEndDate(t *testing.T) {
	var capturedOptions domain.UserFilter
	repo := &mockUserRepo{
		countUserFn: func(ctx context.Context, options domain.UserFilter) int64 {
			capturedOptions = options
			return 0
		},
	}
	svc := NewService(repo)
	query := url.Values{
		"created_at_end": []string{"2024-12-31 23:59:59"},
	}
	svc.SampleUserList(context.Background(), auth.JWTClaimUser{}, query)

	if capturedOptions.CreatedAtLte == nil {
		t.Error("CreatedAtLte should be set when only end is provided")
	}
}

// --- SampleUserExport Tests ---

func TestSampleUserExport_FetchError(t *testing.T) {
	repo := &mockUserRepo{
		fetchUserFn: func(ctx context.Context, options domain.UserFilter) (*sql.Rows, error) {
			return nil, errors.New("fetch error")
		},
	}
	svc := NewService(repo)
	statusCode, _ := svc.SampleUserExport(context.Background(), auth.JWTClaimUser{}, url.Values{})

	if statusCode != http.StatusInternalServerError {
		t.Errorf("statusCode = %d, want %d", statusCode, http.StatusInternalServerError)
	}
}

func TestSampleUserExport_WithDateRangeFilter(t *testing.T) {
	var capturedOptions domain.UserFilter
	repo := &mockUserRepo{
		fetchUserFn: func(ctx context.Context, options domain.UserFilter) (*sql.Rows, error) {
			capturedOptions = options
			return nil, errors.New("stop here")
		},
	}
	svc := NewService(repo)
	query := url.Values{
		"created_at_start": []string{"2024-01-01 00:00:00"},
		"created_at_end":   []string{"2024-12-31 23:59:59"},
	}
	svc.SampleUserExport(context.Background(), auth.JWTClaimUser{}, query)

	if capturedOptions.CreatedAtRange == nil {
		t.Error("CreatedAtRange should be set when both start and end are provided")
	}
}

func TestSampleUserExport_WithOnlyStartDate(t *testing.T) {
	var capturedOptions domain.UserFilter
	repo := &mockUserRepo{
		fetchUserFn: func(ctx context.Context, options domain.UserFilter) (*sql.Rows, error) {
			capturedOptions = options
			return nil, errors.New("stop here")
		},
	}
	svc := NewService(repo)
	query := url.Values{
		"created_at_start": []string{"2024-01-01 00:00:00"},
	}
	svc.SampleUserExport(context.Background(), auth.JWTClaimUser{}, query)

	if capturedOptions.CreatedAtGte == nil {
		t.Error("CreatedAtGte should be set when only start is provided")
	}
}
