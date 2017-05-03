package orientation

// DecodeError is returned by Apply and Decode when image.Decode returns an error.
type DecodeError struct {
	Raw error
}

func (e *DecodeError) Error() string {
	return e.Raw.Error()
}

// Cause returns the underlying cause of DecodeError.
func (e *DecodeError) Cause() error {
	return e.Raw
}

// FormatError is returned by Apply and Decode when the format of decoded image is not jpeg.
type FormatError struct {
	Raw error
}

func (e *FormatError) Error() string {
	return e.Raw.Error()
}

// Cause returns the underlying cause of FormatError.
func (e *FormatError) Cause() error {
	return e.Raw
}

// TagError is returned by Apply and Tag when the EXIF can not be decoded, the orientation tag does not exits,
// or the orientation tag is not int.
type TagError struct {
	Raw error
}

func (e *TagError) Error() string {
	return e.Raw.Error()
}

// Cause returns the underlying cause of TagError.
func (e *TagError) Cause() error {
	return e.Raw
}

// OrientError is returned by Apply and Tag when the specified orientation tag is unknown.
type OrientError struct {
	Raw error
}

func (e *OrientError) Error() string {
	return e.Raw.Error()
}

// Cause returns the underlying cause of OrientError.
func (e *OrientError) Cause() error {
	return e.Raw
}
