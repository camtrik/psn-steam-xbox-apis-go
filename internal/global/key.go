package global

import "time"

// Redis Settings
const (
	USER_TITLES_KEY = "user_titles:%s"

	DEFAULT_EXPIRATION = 2 * time.Hour
)

// PSN
const (
	AUTH_BASE_URL   = "https://ca.account.sony.com/api/authz/v3/oauth"
	TROPHY_BASE_URL = "https://m.np.playstation.com/api/trophy"
)
