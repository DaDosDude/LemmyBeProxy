package lemmy017

// GetCaptchaResponse matches real 0.17.2's shape. No omitempty: real
// 0.17.x's deserializer requires the key present (as null) even when
// the value is absent, same lesson as everywhere else in this package.
type GetCaptchaResponse struct {
	Ok *CaptchaResponse `json:"ok"`
}

type CaptchaResponse struct {
	Png  string `json:"png"`
	Wav  string `json:"wav"`
	Uuid string `json:"uuid"`
}
