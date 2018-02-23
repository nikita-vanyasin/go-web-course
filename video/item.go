package video

const (
	STATUS_CREATED    int8 = 1
	STATUS_PROCESSING int8 = 2
	STATUS_READY      int8 = 3
	STATUS_DELETED    int8 = 4
	STATUS_ERROR      int8 = 5
)

type Item struct {
	ID        string `json:"id"`
	Name      string `json:"name"`
	Duration  int64  `json:"duration"`
	Thumbnail string `json:"thumbnail"`
	Url       string `json:"url"`
	Status    int8   `json:"status"`
}
