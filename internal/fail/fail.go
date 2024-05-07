// ---------------------------------------------------------------------------------------------- //
// -- Copyright (c) 2024 Braden Hitchcock - MIT License (https://opensource.org/licenses/MIT)  -- //
// ---------------------------------------------------------------------------------------------- //
package fail

import "fmt"

func Zero[V any](details string, err error) (V, error) {
	var v V
	return v, Format(details, err)
}

func Format(details string, err error) error {
	return fmt.Errorf("%s: %v", details, err)
}
