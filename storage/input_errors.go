package storage

import "fmt"

type InvalidSchemeError struct {
	Scheme string
}

func NewInvalidSchemeError(scheme string) InvalidSchemeError {
	return InvalidSchemeError{scheme}
}

func (err InvalidSchemeError) Error() string {
	return fmt.Sprintf("scheme '%s' isn't allowed", err.Scheme)
}

type InvalidHostError struct {
	Host string
}

func NewInvalidHostError(host string) InvalidHostError {
	return InvalidHostError{host}
}

func (err InvalidHostError) Error() string {
	return fmt.Sprintf("host '%s' isn't allowed", err.Host)
}

type InvalidSizeError struct {
	Width  int
	Height int
}

func NewInvalidSizeError(width, height int) InvalidSizeError {
	return InvalidSizeError{width, height}
}

func (err InvalidSizeError) Error() string {
	return fmt.Sprintf("size %d * %d isn't allowed", err.Width, err.Height)
}

type InvalidMethodError struct {
	Method string
}

func NewInvalidMethodError(method string) InvalidMethodError {
	return InvalidMethodError{method}
}

func (err InvalidMethodError) Error() string {
	return fmt.Sprintf("method '%s' isn't allowed", err.Method)
}

type InvalidFormatError struct {
	Format string
}

func NewInvalidFormatError(format string) InvalidFormatError {
	return InvalidFormatError{format}
}

func (err InvalidFormatError) Error() string {
	return fmt.Sprintf("format '%s' isn't allowed", err.Format)
}

type InvalidQualityError struct {
	Quality int
}

func NewInvalidQualityError(quality int) InvalidQualityError {
	return InvalidQualityError{quality}
}

func (err InvalidQualityError) Error() string {
	return fmt.Sprintf("quality %d isn't allowed", err.Quality)
}
