package entities

type Event struct {
	ID          string `json:"uuid"`
	Title       string `json:"title"`
	Description string `json:"description"`
	StartedAt   string `json:"start_at"`
	FinishedAt  string `json:"finished_at"`
	NotifyAt    string `json:"notify_at"`
	UserID      string `json:"user_id"`
}
