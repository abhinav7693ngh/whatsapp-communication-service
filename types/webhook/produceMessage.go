package webhook

type ProduceWaMessage struct {
	WaId   string `json:"id"`
	Status string `json:"status"`
}
