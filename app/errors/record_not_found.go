package errors

// RecordNotFound is an error that occurs when the record is not found
var ErrRecordNotFound = NewNotFound(NotFoundParam{})
