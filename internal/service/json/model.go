package json

type SettingsDTO struct {
	Id    int64  `json:"id"`
	Name  string `json:"name"`
	Color string `json:"color"`
	Lives int8   `json:"lives"`
}
