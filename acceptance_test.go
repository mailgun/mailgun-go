package mailgun

import (
	"fmt"
	"os"
	"testing"
)

// Return the variable missing which caused the test to be skipped
func SkipNetworkTest() string {
	for _, env := range []string{"MG_DOMAIN", "MG_API_KEY", "MG_EMAIL_TO", "MG_PUBLIC_API_KEY"} {
		if os.Getenv(env) == "" {
			return fmt.Sprintf("'%s' missing from environment skipping...", env)
		}
	}
	return ""
}

func spendMoney(t *testing.T, tFunc func()) {
	ok := os.Getenv("MG_SPEND_MONEY")
	if ok != "" {
		tFunc()
	} else {
		t.Log("Money spending not allowed, not running function.")
	}
}

