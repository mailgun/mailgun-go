package mtypes

const (
	APIKeysEndpoint           = "keys"
	APIKeysRegenerateEndpoint = APIKeysEndpoint + "/public"
	APIKeysVersion            = 1
)

type GetAPIKeyListResponse struct {
	Items []APIKey `json:"items"`
}

type CreateAPIKeyResponse struct {
	Key APIKey `json:"key"`
}

type DeleteAPIKeyResponse struct {
	Message string `json:"message"`
}

type APIKey struct {
	ID             string      `json:"id"`
	Description    string      `json:"description"`
	Kind           string      `json:"kind"`
	Role           string      `json:"role"`
	CreatedAt      RFC2822Time `json:"created_at"`
	UpdatedAt      RFC2822Time `json:"updated_at"`
	DomainName     string      `json:"domain_name"`
	Requestor      string      `json:"requestor"`
	UserName       string      `json:"user_name"`
	IsDisabled     bool        `json:"is_disabled"`
	ExpiresAt      RFC2822Time `json:"expires_at"`
	Secret         string      `json:"string"`
	DisabledReason string      `json:"disabled_reason"`
}
