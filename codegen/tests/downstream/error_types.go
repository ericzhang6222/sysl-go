// Code generated by sysl DO NOT EDIT.
package downstream

import (
	"fmt"
)

// Error fulfills the error type interface for Status
func (s Status) Error() string {
	type plain Status

	return fmt.Sprintf("%+v", plain(s))
}
