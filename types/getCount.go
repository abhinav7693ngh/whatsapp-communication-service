package types

type GetCountPayload struct {
	Clients   []string `json:"clients" validate:"required" message:"clients is required"`
	Status    []string `json:"status" validate:"required" message:"status is required"`
	StartTime int64    `json:"startTime" validate:"required" message:"startTime is required"`
	EndTime   int64    `json:"endTime" validate:"required" message:"endTime is required"`
}

type GetCountResponse struct {
	Count  int64  `json:"count"`
	Status string `json:"status"`
}
