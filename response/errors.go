package response

type ErrorResponse struct {
	Name       string
	StatusCode int
}

var errorsList = []*ErrorResponse{
	{"unknown_error", 500},
	{"internal_error", 500},
	{"invalid_payload", 422},
}

func ErrorByName(name string) ErrorResponse {
	for _, err := range errorsList {
		if err.Name == name {
			return *err
		}
	}
	return ErrorByName("unknown_error")
}
