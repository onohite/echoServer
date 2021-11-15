package discord

type DiscordResp struct {
	ID            string      `json:"id"`
	Username      string      `json:"username"`
	Avatar        string      `json:"avatar"`
	Discriminator string      `json:"discriminator"`
	PublicFlags   int         `json:"public_flags"`
	Flags         int         `json:"flags"`
	Banner        interface{} `json:"banner"`
	BannerColor   interface{} `json:"banner_color"`
	AccentColor   interface{} `json:"accent_color"`
	Locale        string      `json:"locale"`
	MfaEnabled    bool        `json:"mfa_enabled"`
	Email         string      `json:"email"`
	Verified      bool        `json:"verified"`
}
