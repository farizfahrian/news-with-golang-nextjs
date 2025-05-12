package response

type ErrorResponseDefault struct {
	Meta Meta `json:"meta"`
}

type Meta struct {
	Status  bool   `json:"status"`
	Message string `json:"message"`
}
