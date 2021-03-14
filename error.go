package main

// Error is selfmade..
type Error interface {
	error
	Status() int
}

// StatusError is ..
type StatusError struct {
	Code int
	Err  error
}

func (se StatusError) Error() string {
	return se.Err.Error()
}

// Status is ..
func (se StatusError) Status() int {
	return se.Code
}
