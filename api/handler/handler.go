package handler

import (
	"encoding/json"
	"github.com/uCibar/bootcamp-radio/api/response"
	"log"
	"net/http"
)

type handler struct {
	logger *log.Logger
}

func newHandler(logger *log.Logger) *handler {
	return &handler{logger: logger}
}

func (h *handler) writeLog(msg string) {
	if h.logger == nil {
		return
	}

	h.logger.Println(msg)
}

func (h *handler) readBody(r *http.Request, b interface{}) error {
	d := json.NewDecoder(r.Body)
	return d.Decode(b)
}

func (h *handler) writeResponse(w http.ResponseWriter, res *response.Response) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(res.Code)

	if err := json.NewEncoder(w).Encode(res); err != nil {
		h.writeLog(err.Error())
	}
}
