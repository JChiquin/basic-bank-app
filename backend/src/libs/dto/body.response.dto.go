package dto

//BodyResponse with fields according confluence
type BodyResponse struct {
	Message string              `json:"message"`
	Errors  []map[string]string `json:"errors"`
	Data    interface{}         `json:"data"`
}

//NewBodyResponse is a constructor for BodeResponseDTO
func NewBodyResponse(message string, errors []map[string]string, data interface{}) *BodyResponse {
	return &BodyResponse{
		Message: message,
		Errors:  errors,
		Data:    data,
	}
}
