package httpCodes

const StatusError = "Error"
const StatusOK = "OK"

type Response struct {
	Status string `json:"status"`
	Error  string `json:"error,omitempty"`
}
