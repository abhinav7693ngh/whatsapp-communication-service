package template

type CompositeBody struct {
	Type string `json:"type"`
	// parameters can be NuclearText
	Parameters []interface{} `json:"parameters"`
}

type CompositeButton struct {
	Type string `json:"type"`
	// sub_type can be quick_reply and url
	Sub_Type string `json:"sub_type"`
	Index    string `json:"index"`
	// parameters can be NuclearPayload in case of quick_reply and NuclearOtpText in case of url
	Parameters []interface{} `json:"parameters"`
}

type CompositeHeader struct {
	Type string `json:"type"`
	// parameters can be NuclearLocation, NuclearImage
	Parameters []interface{} `json:"parameters"`
}
