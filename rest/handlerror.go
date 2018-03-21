package rest

import (
	"net/http"
)

func notFound(w http.ResponseWriter) {
	w.WriteHeader(404)
	w.Write([]byte("not found"))
}

func methodNotSupport(w http.ResponseWriter) {
	w.WriteHeader(405)
	w.Write([]byte("method not support"))
}

func badRequest(w http.ResponseWriter) {
	w.WriteHeader(400)
	w.Write([]byte("400 Bad Request"))
}
