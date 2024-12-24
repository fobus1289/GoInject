package goinject

import (
	"errors"
	"reflect"
	"sync"
)

type (
	typeof = reflect.Type
)

var (
	ErrServiceNotFound            = errors.New("service not found")
	ErrFactoryMustBeAFunction     = errors.New("factory must be a function")
	ErrFactoryMustReturnOneValue  = errors.New("factory must return one value")
	ErrFactoryMustTakeNoArguments = errors.New("factory must take no arguments")
	ErrOutputMustBeAPointer       = errors.New("output must be a pointer")
)

type Container struct {
	factories map[typeof]func() any
	providers map[typeof]any
	mu        sync.RWMutex
}

// New creates a new Container instance.
// It returns a pointer to the Container.
//
// Example:
//
//	container := goinject.New()
func New() *Container {
	return &Container{
		factories: make(map[typeof]func() any),
		providers: make(map[typeof]any),
	}
}

// RegisterFactory registers a factory function that returns a new instance of the given type.
// It returns an error if the factory is not a function or does not return a pointer.
//
// Example:
//
//	container.RegisterFactory(func() *User {
//	    return &User{ID: 1, Name: "John", Age: 25, Salary: 50000.0}
//	})
func (c *Container) RegisterFactory(factory any) error {

	c.mu.RLock()
	defer c.mu.RUnlock()

	factoryValue := reflect.ValueOf(factory)

	factoryType := factoryValue.Type()
	{
		if factoryType.Kind() != reflect.Func {
			return ErrFactoryMustBeAFunction
		}

		if factoryType.NumIn() != 0 {
			return ErrFactoryMustTakeNoArguments
		}

		if factoryType.NumOut() != 1 {
			return ErrFactoryMustReturnOneValue
		}
	}

	typeof := factoryType.Out(0)

	if typeof.Kind() != reflect.Ptr {
		return ErrOutputMustBeAPointer
	}

	c.factories[typeof] = func() any {
		return factoryValue.Call(nil)[0].Interface()
	}

	return nil
}

// Register registers a singleton instance of the given type.
// It returns an error if the input is not a pointer.
//
// Example:
//
//	container.Register(&User{ID: 1, Name: "John", Age: 25, Salary: 50000.0})
func (c *Container) Register(service any) error {
	typeof := reflect.TypeOf(service)
	{
		if typeof.Kind() != reflect.Ptr {
			return ErrOutputMustBeAPointer
		}
	}

	c.mu.RLock()
	defer c.mu.RUnlock()

	c.providers[typeof] = service

	return nil
}

// Get retrieves a dependency of the given type from the container.
// It returns an error if the dependency is not found.
//
// Example:
//
//	var user User
//	err := container.Get(&user)
//	if err != nil {
//	    log.Fatal(err)
//	}
//	fmt.Println(user.Name) // Prints: John
func (c *Container) Get(out any) (any, error) {

	typeof := reflect.TypeOf(out)
	{
		if typeof.Kind() != reflect.Ptr {
			return nil, ErrOutputMustBeAPointer
		}
	}

	c.mu.RLock()
	defer c.mu.RUnlock()

	service, ok := c.providers[typeof]

	if !ok {
		if factory := c.factories[typeof]; factory != nil {
			service = factory()
		} else {
			return nil, ErrServiceNotFound
		}
	}

	if !ok {
		c.providers[typeof] = service
	}

	return service, nil
}

// GetValue retrieves a dependency and copies its value into the provided pointer.
// It returns an error if the dependency is not found.
//
// Example:
//
//	var user User
//	err := container.GetValue(&user)
//	if err != nil {
//	    log.Fatal(err)
//	}
//	fmt.Println(user.Name) // Prints: John
func (c *Container) GetValue(out any) error {

	service, err := c.Get(out)
	{
		if err != nil {
			return err
		}
	}

	servicePtr := reflect.ValueOf(service).Elem()

	setOutValue := reflect.ValueOf(out).Elem()

	setOutValue.Set(servicePtr)

	return nil
}
