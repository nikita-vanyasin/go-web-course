package video

const (
	StatusCreated    int8 = 1
	StatusProcessing int8 = 2
	StatusReady      int8 = 3
	StatusDeleted    int8 = 4
	StatusError      int8 = 5
)

type Item struct {
	ID        string `json:"id"`
	Name      string `json:"name"`
	Duration  int64  `json:"duration"`
	Thumbnail string `json:"thumbnail"`
	URL       string `json:"url"`
	Status    int8   `json:"status"`
}
