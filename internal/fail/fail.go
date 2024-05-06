// ---------------------------------------------------------------------------------------------- //
// -- Copyright (c) 2024 Braden Hitchcock - MIT License (https://opensource.org/licenses/MIT)  -- //
// ---------------------------------------------------------------------------------------------- //
package fail

import "fmt"

func Zero[V any](details string, err error) (V, error) {
	var v V
	return v, fmt.Errorf("%s: %v", details, err)
}
