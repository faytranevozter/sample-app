package response

import (
	"app/domain"
	"testing"
)

func TestError(t *testing.T) {
	resp := Error(500, "internal error")

	if resp.ErrorCode != 500 {
		t.Errorf("ErrorCode = %d, want 500", resp.ErrorCode)
	}
	if resp.Message != "internal error" {
		t.Errorf("Message = %q, want %q", resp.Message, "internal error")
	}
	if resp.Data != nil {
		t.Errorf("Data = %v, want nil", resp.Data)
	}
	if resp.Validation != nil {
		t.Errorf("Validation = %v, want nil", resp.Validation)
	}
}

func TestErrorValidation(t *testing.T) {
	validation := map[string]string{
		"email":    "email is required",
		"password": "password is required",
	}
	resp := ErrorValidation(validation, "validation error")

	if resp.ErrorCode != domain.ErrValidationCode {
		t.Errorf("ErrorCode = %d, want %d", resp.ErrorCode, domain.ErrValidationCode)
	}
	if resp.Message != "validation error" {
		t.Errorf("Message = %q, want %q", resp.Message, "validation error")
	}
	if resp.Data != nil {
		t.Errorf("Data = %v, want nil", resp.Data)
	}
	if len(resp.Validation) != 2 {
		t.Errorf("Validation length = %d, want 2", len(resp.Validation))
	}
	if resp.Validation["email"] != "email is required" {
		t.Errorf("Validation[email] = %q, want %q", resp.Validation["email"], "email is required")
	}
	if resp.Validation["password"] != "password is required" {
		t.Errorf("Validation[password] = %q, want %q", resp.Validation["password"], "password is required")
	}
}

func TestSuccess(t *testing.T) {
	data := map[string]string{"name": "test"}
	resp := Success(data)

	if resp.ErrorCode != 0 {
		t.Errorf("ErrorCode = %d, want 0", resp.ErrorCode)
	}
	if resp.Message != "OK" {
		t.Errorf("Message = %q, want %q", resp.Message, "OK")
	}
	if resp.Validation != nil {
		t.Errorf("Validation = %v, want nil", resp.Validation)
	}

	respData, ok := resp.Data.(map[string]string)
	if !ok {
		t.Fatalf("Data type = %T, want map[string]string", resp.Data)
	}
	if respData["name"] != "test" {
		t.Errorf("Data[name] = %q, want %q", respData["name"], "test")
	}
}

func TestSuccess_NilData(t *testing.T) {
	resp := Success(nil)

	if resp.ErrorCode != 0 {
		t.Errorf("ErrorCode = %d, want 0", resp.ErrorCode)
	}
	if resp.Message != "OK" {
		t.Errorf("Message = %q, want %q", resp.Message, "OK")
	}
	if resp.Data != nil {
		t.Errorf("Data = %v, want nil", resp.Data)
	}
}

func TestError_DifferentCodes(t *testing.T) {
	tests := []struct {
		name    string
		code    int
		message string
	}{
		{"bad request", 400, "bad request"},
		{"not found", 404, "not found"},
		{"conflict", 409, "conflict"},
		{"server error", 500, "server error"},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			resp := Error(tt.code, tt.message)
			if resp.ErrorCode != tt.code {
				t.Errorf("ErrorCode = %d, want %d", resp.ErrorCode, tt.code)
			}
			if resp.Message != tt.message {
				t.Errorf("Message = %q, want %q", resp.Message, tt.message)
			}
		})
	}
}
