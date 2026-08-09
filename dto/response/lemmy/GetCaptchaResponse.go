package lemmy

// CaptchaResponse matches real Lemmy's shape exactly.
type CaptchaResponse struct {
	Png  string `json:"png"`
	Wav  string `json:"wav"`
	Uuid string `json:"uuid"`
}

// GetCaptchaResponse — Ok is nil when captchas are disabled, matching
// real Lemmy's own semantics exactly. LemmyBackend genuinely passes
// through real Lemmy's own captcha generation. PiefedBackend always
// returns nil — this proxy has no image/audio captcha generation
// capability of its own, and Piefed has no equivalent feature to
// forward, so nil is an honest "no captcha required" rather than a
// faked one. Correct either way for retrolemmy.com specifically, which
// has captcha_enabled: false.
type GetCaptchaResponse struct {
	Ok *CaptchaResponse `json:"ok,omitempty"`
}
