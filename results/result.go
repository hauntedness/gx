package results

import "fmt"

type Result[T any] struct {
	value T
	err   error
}

func From[T any](value T, err error) Result[T] {
	return Result[T]{
		value: value,
		err:   err,
	}
}

// Catch map the result of f, wrap panic as error
func Catch[T any](f func() T) (result Result[T]) {
	defer func() {
		if v := recover(); v != nil {
			if err, ok := v.(error); ok {
				result = Err[T](err)
			} else {
				result = Err[T](fmt.Errorf("%v", v))
			}
		}
	}()
	result = Ok(f())
	return result
}

// CatchE is similar to Catch, but also accept error
func CatchE[T any](f func() (T, error)) (result Result[T]) {
	defer func() {
		if v := recover(); v != nil {
			if err, ok := v.(error); ok {
				result = Err[T](err)
			} else {
				result = Err[T](fmt.Errorf("%v", v))
			}
		}
	}()
	result = From(f())
	return result
}

func Must[T any](value T, err error) T {
	if err != nil {
		panic(err)
	}
	return value
}

func Ok[T any](value T) Result[T] {
	return Result[T]{
		value: value,
		err:   nil,
	}
}

func Err[T any](err error) Result[T] {
	if err == nil {
		panic("error is nil")
	}
	var value T
	return Result[T]{
		value: value,
		err:   err,
	}
}

func (r Result[T]) Get() T {
	if r.err != nil {
		panic(r.err)
	}
	return r.value
}

func (r Result[T]) Unpack() (T, error) {
	return r.value, r.err
}

func (r Result[T]) Err() error {
	return r.err
}

func (r Result[T]) OrElse(other T) T {
	if r.err != nil {
		return other
	}
	return r.value
}

func (r Result[T]) OrEmpty() T {
	if r.err != nil {
		var other T
		return other
	}
	return r.value
}

func (r Result[T]) OrFrom(fn func() T) T {
	if r.err != nil {
		return fn()
	}
	return r.value
}

func (r Result[T]) IsOk() bool {
	return r.err == nil
}

func (r Result[T]) IsErr() bool {
	return r.err != nil
}
