package lemmy

// Jwt is nil when registration requires admin approval or email
// verification before login is possible — real Lemmy's own jwt field
// is Option<String> for exactly this reason (Register reuses this same
// response type, and retrolemmy.com's own "requireapplication" mode is
// a live example: jwt is genuinely absent until an admin approves the
// account). No omitempty: matches every other 0.17.x-shaped field in
// this codebase — the key must stay present, as null, not be dropped.
type LoginResponse struct {
	Jwt                 *string `json:"jwt"`
	RegistrationCreated bool    `json:"registration_created"`
	VerifyEmailSent     bool    `json:"verify_email_sent"`
}
