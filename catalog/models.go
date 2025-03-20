package catalog

type Product struct {
	ID          string  `json:"id"`
	Name        string  `json:"name"`
	Description string  `json:"description"`
	Price       float64 `json:"price"`
}

type productDocument struct {
	Name        string  `json:"name"`
	Description string  `json:"description"`
	Price       float64 `json:"price"`
}

type searchResponse struct {
	Hits struct {
		Hits []struct {
			Source Product `json:"_source"`
		} `json:"hits"`
	} `json:"hits"`
}
