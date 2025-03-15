package mailgun

import (
	"fmt"
	"os"
)

// Return the variable missing which caused the test to be skipped
func SkipNetworkTest() string {
	for _, env := range []string{"MG_DOMAIN", "MG_API_KEY", "MG_EMAIL_TO"} {
		if os.Getenv(env) == "" {
			return fmt.Sprintf("'%s' missing from environment skipping...", env)
		}
	}
	return ""
}
