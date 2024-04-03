package helpers

import (
	"net/http"
)

func GenerateRouteLink(r *http.Request, path string) string {
	return "http://" + r.Host + path
}
