package piefed

// RegisterRequest — Piefed's exact registration API shape is an
// unverified assumption, matching every other endpoint's established
// 1:1 naming mirror of Lemmy's own API. Worth confirming directly
// against a live Piefed instance's registration flow.
type RegisterRequest struct {
	Username       string `json:"username" validate:"required"`
	Password       string `json:"password" validate:"required"`
	PasswordVerify string `json:"password_verify" validate:"required"`
	ShowNsfw       bool   `json:"show_nsfw"`
	Email          string `json:"email,omitempty"`
	Answer         string `json:"answer,omitempty"`
}
