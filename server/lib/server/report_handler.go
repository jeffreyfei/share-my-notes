package server

import (
	"net/http"
	"strconv"
)

// Handles healthcheck requests from load balancer
// Returns the number of unprocessed jobs to the load balancer
func (s *Server) reportStatusHandler(w http.ResponseWriter, r *http.Request) {
	w.Write([]byte(strconv.Itoa(s.buffer.JobCount())))
}
