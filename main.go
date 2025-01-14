package main

import (
	"net/http"

	"github.com/opchaves/tudo/internal/config"
	"github.com/opchaves/tudo/internal/server"
)

func main() {
	s := server.CreateNewServer(nil)
	s.MountHandlers()
	http.ListenAndServe(":"+config.Port, s.Router)
}
