package template

// ================ Location ================ //

type LocationObject struct {
	Latitude  float64 `json:"latitude"`
	Longitude float64 `json:"longitude"`
	Name      string  `json:"name"`
	Address   string  `json:"address"`
}

type NuclearLocation struct {
	Type     string         `json:"type"`
	Location LocationObject `json:"location"`
}

// ============ Image ============ //

type ImageObject struct {
	Link string `json:"link"`
}

type NuclearImage struct {
	Type  string      `json:"type"`
	Image ImageObject `json:"link"`
}

// ============ Audio ============ //

type AudioObject struct {
	Link string `json:"link"`
}

type NuclearAudio struct {
	Type  string      `json:"type"`
	Audio AudioObject `json:"audio"`
}

// ============ Video ============ //

type VideoObject struct {
	Link    string `json:"link"`
	Caption string `json:"caption"`
}

type NuclearVideo struct {
	Type  string      `json:"type"`
	Video VideoObject `json:"video"`
}

// ============ Document ============ //

type DocumentObject struct {
	Link    string `json:"link"`
	Caption string `json:"caption"`
}

type NuclearDocument struct {
	Type     string         `json:"type"`
	Document DocumentObject `json:"document"`
}

// ======== Text ======== //

type NuclearText struct {
	Type string `json:"type"`
	Text string `json:"text"`
}

// ======== OTP Text (button) ======== //

type NuclearOtpText struct {
	Type string `json:"type"`
	Text string `json:"text"`
}

// ============ Payload ( button ) ============ //

type NuclearPayload struct {
	Type    string `json:"type"`
	Payload string `json:"payload"`
}
