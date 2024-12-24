package goinject

// Get retrieves a dependency of type T from the container.
// It returns a pointer to the dependency and an error if not found.
//
// Example:
//
//	type UserService struct {
//	    Name string
//	}
//
//	container := goinject.New()
//	container.Register(&UserService{Name: "John"})
//
//	userService, err := goinject.Get[UserService](container)
//	if err != nil {
//	    log.Fatal(err)
//	}
//	fmt.Println(userService.Name) // Prints: John
// GetValue retrieves a dependency and copies its value into the provided pointer.
// It returns an error if the dependency is not found.
//
// Example:
//
//	var user UserService
//	err := goinject.GetValue(container, &user)
//	if err != nil {
//	    log.Fatal(err)
//	}
//	fmt.Println(user.Name) // Prints: John

// MustGet retrieves a dependency of type T from the container.
// It panics if the dependency is not found.
//
// Example:
//
//	userService := goinject.MustGet[UserService](container)
//	fmt.Println(userService.Name) // Prints: John
func Get[T any](c *Container) (*T, error) {

	var out T

	v, err := c.Get(&out)
	{
		if err != nil {
			return nil, err
		}
	}

	o, ok := v.(*T)
	{
		if !ok {
			return nil, ErrOutputMustBeAPointer
		}
	}

	return o, nil
}

// GetValue retrieves a dependency and copies its value into the provided pointer.
// It returns an error if the dependency is not found.
//
// Example:
//
//	var user UserService
//	err := goinject.GetValue(container, &user)
//	if err != nil {
//	    log.Fatal(err)
//	}
//	fmt.Println(user.Name) // Prints: John
func GetValue[T any](c *Container, out *T) error {
	return c.GetValue(out)
}

// MustGet retrieves a dependency of type T from the container.
// It panics if the dependency is not found.
//
// Example:
//
//	userService := goinject.MustGet[UserService](container)
//	fmt.Println(userService.Name) // Prints: John
func MustGet[T any](c *Container) *T {

	v, err := Get[T](c)
	{
		if err != nil {
			panic(err)
		}
	}

	return v
}
