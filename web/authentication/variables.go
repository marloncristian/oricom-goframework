package authentication

var (
	tokenSecret string
)

// Initialize initilizes the global variables from configuration
func Initialize(secret string) {
	tokenSecret = secret
}
