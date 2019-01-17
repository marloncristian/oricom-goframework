package database

// EntityNotFound generic entity not found error
type EntityNotFound struct {
}

// Error generic description error method
func (err EntityNotFound) Error() string {
	return "Entity not found"
}

// InvalidArgument generic invalid argument error
type InvalidArgument struct {
	Description string
}

// Error generic description error method
func (err InvalidArgument) Error() string {
	return err.Description
}
