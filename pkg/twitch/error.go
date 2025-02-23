package twitch

type ErrResponse struct {
	Message string `json:"message"`
	Status  int    `json:"status"`
}
