package errs

// NonErrorRecordSkipped is special error that tells omniparser '1.0' to skip a record.
type NonErrorRecordSkipped string

// Error() implements error interface.
func (e NonErrorRecordSkipped) Error() string { return string(e) }

// IsNonErrorRecordSkipped checks if an err is of NonErrorRecordSkipped type or not.
func IsNonErrorRecordSkipped(err error) bool {
	switch err.(type) {
	case NonErrorRecordSkipped:
		return true
	default:
		return false
	}
}
