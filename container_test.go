package goinject

import (
	"testing"
)

type (
	TestService struct {
		Name string
	}

	AnotherService struct {
		ID int
	}
)

func TestContainer_Register(t *testing.T) {
	tests := []struct {
		name    string
		service interface{}
		wantErr bool
	}{
		{
			name:    "successful registration of struct pointer",
			service: &TestService{Name: "test"},
			wantErr: false,
		},
		{
			name:    "error when registering value (not pointer)",
			service: TestService{Name: "test"},
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			c := New()
			err := c.Register(tt.service)
			if (err != nil) != tt.wantErr {
				t.Errorf("Register() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestContainer_RegisterFactory(t *testing.T) {
	tests := []struct {
		name    string
		factory interface{}
		wantErr bool
	}{
		{
			name: "successful factory registration",
			factory: func() *TestService {
				return &TestService{Name: "test"}
			},
			wantErr: false,
		},
		{
			name:    "error: factory is not a function",
			factory: "not a function",
			wantErr: true,
		},
		{
			name: "error: factory takes arguments",
			factory: func(name string) *TestService {
				return &TestService{Name: name}
			},
			wantErr: true,
		},
		{
			name: "error: factory returns non-pointer",
			factory: func() TestService {
				return TestService{Name: "test"}
			},
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			c := New()
			err := c.RegisterFactory(tt.factory)
			if (err != nil) != tt.wantErr {
				t.Errorf("RegisterFactory() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestContainer_Get(t *testing.T) {
	c := New()
	service := &TestService{Name: "test"}

	// Register service
	if err := c.Register(service); err != nil {
		t.Fatalf("failed to register service: %v", err)
	}

	// Test service retrieval
	var out TestService
	result, err := c.Get(&out)
	if err != nil {
		t.Errorf("Get() unexpected error = %v", err)
	}

	retrieved, ok := result.(*TestService)
	if !ok {
		t.Error("Get() did not return correct type")
	}

	if retrieved.Name != service.Name {
		t.Errorf("Get() got = %v, want %v", retrieved.Name, service.Name)
	}
}

func TestGenericGet(t *testing.T) {
	c := New()
	service := &TestService{Name: "test"}

	// Register service
	if err := c.Register(service); err != nil {
		t.Fatalf("failed to register service: %v", err)
	}

	// Test generic Get function
	result, err := Get[TestService](c)
	if err != nil {
		t.Errorf("Get[T]() unexpected error = %v", err)
	}

	if result.Name != service.Name {
		t.Errorf("Get[T]() got = %v, want %v", result.Name, service.Name)
	}
}

func TestGenericGetValue(t *testing.T) {
	c := New()
	service := &TestService{Name: "test"}

	// Register service
	if err := c.Register(service); err != nil {
		t.Fatalf("failed to register service: %v", err)
	}

	// Test generic GetValue function
	var result TestService
	err := GetValue(c, &result)
	if err != nil {
		t.Errorf("GetValue[T]() unexpected error = %v", err)
	}

	if result.Name != service.Name {
		t.Errorf("GetValue[T]() got = %v, want %v", result.Name, service.Name)
	}
}

func TestMustGet(t *testing.T) {
	c := New()
	service := &TestService{Name: "test"}

	// Register service
	if err := c.Register(service); err != nil {
		t.Fatalf("failed to register service: %v", err)
	}

	// Test successful case
	result := MustGet[TestService](c)
	if result.Name != service.Name {
		t.Errorf("MustGet[T]() got = %v, want %v", result.Name, service.Name)
	}

	// Test panic on missing service
	defer func() {
		if r := recover(); r == nil {
			t.Error("MustGet[T]() should panic when service not found")
		}
	}()

	_ = MustGet[AnotherService](c)
}
