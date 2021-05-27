package errors

import (
	"fmt"
)

func ExampleWrap() {
	err := Wrap(ErrNotFound, "message")
	fmt.Println(err)
	// Output: message
}

func ExampleWrap_multiple() {
	err := Wrap(ErrNotFound, "original message")
	err = Wrap(err, "prefixed message")
	fmt.Println(err)
	// Output: prefixed message: original message
}
