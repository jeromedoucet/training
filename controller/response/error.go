package response

// Error is used to return details over on issue when possible.
type Error struct {
	Message string `json:"message"`
}
