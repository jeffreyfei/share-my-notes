package server

import (
	"net/http"
	"strconv"
)

func (s *Server) reportStatusHandler(w http.ResponseWriter, r *http.Request) {
	w.Write([]byte(strconv.Itoa(s.buffer.JobCount())))
}
