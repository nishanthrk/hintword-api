package tab_controller

type PayloadCollection struct {
	CollectionID string `json:"collection_id"`
	Name         string `json:"name" validate:"required"`
	Status       string `json:"status" validate:"oneof=ACTIVE INACTIVE"`
}

type PayloadTab struct {
	TabID        string `json:"tab_id"`
	Title        string `json:"title" validate:"required"`
	URL          string `json:"url" validate:"required"`
	FaviconURL   string `json:"favicon_url" validate:"required"`
	CollectionID string `json:"collection_id" validate:"required"`
	Sequence     int64  `json:"sequence" validate:"required"`
	Status       string `json:"status" validate:"oneof=ACTIVE INACTIVE"`
}
