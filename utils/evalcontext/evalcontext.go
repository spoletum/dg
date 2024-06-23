package evalcontext

import (
	"fmt"
	"log/slog"
	"strings"

	"github.com/zclconf/go-cty/cty"
)

// fromEnvironment returns a map of cty.Value from the environment. For testability purposes, the environment is passed as an array of strings to mimic the os.Environ() function.
func FromEnvironment(source []string) (map[string]cty.Value, error) {

	result := make(map[string]cty.Value)

	for _, e := range source {

		pair := strings.SplitN(e, "=", 2)

		if len(pair) == 2 {
			slog.Debug("Adding environment variable to the evaluation context", "name", pair[0])
			result[pair[0]] = cty.StringVal(pair[1])
		} else {
			return nil, fmt.Errorf("could not parse environment variable: %s", e)
		}
	}
	return result, nil
}
