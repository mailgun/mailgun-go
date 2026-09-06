package mailgun

import (
	"fmt"
	"os"
)

// Return the variable missing which caused the test to be skipped
// TODO(vtopc): merge with `//go:build integration`
func SkipNetworkTest() string {
	for _, env := range []string{"MG_DOMAIN", "MG_API_KEY", "MG_EMAIL_TO"} {
		if os.Getenv(env) == "" {
			return fmt.Sprintf("'%s' missing from environment skipping...", env)
		}
	}
	return ""
}
