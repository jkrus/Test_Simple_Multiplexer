package models

type (
	Info struct {
		URL  string `json:"url"`
		Body []byte `json:"body"`
	}

	InfoDTO struct {
		Urls []string `json:"urls"`
	}
)
