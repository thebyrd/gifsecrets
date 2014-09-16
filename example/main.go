package main

import (
	"net/http"
	"os"
	"strings"

	"github.com/codegangsta/negroni"
	"github.com/gorilla/mux"
	"github.com/thebyrd/gifsecrets"
	"github.com/unrolled/render"
)

var r = render.New(render.Options{
	IndentJSON:    true,
	IsDevelopment: true,
})

// encodeSecretHandler encodes a given secret in a gif
func encodeSecretHandler(rw http.ResponseWriter, req *http.Request) {
	response := map[string]string{
		"status": "success",
	}

	secret, err := gifsecrets.Decode("../coin.gif")
	if err != nil {
		response["status"] = "error"
		response["error"] = err.Error()
	}
	response["secret"] = secret

	r.JSON(rw, http.StatusOK, response)
}

//
// Boiler Plate that you should ignore...
//
type Route struct {
	Method  string
	Path    string
	Handler http.HandlerFunc
}

var Routes = []*Route{
	&Route{"GET", "/", http.FileServer(http.Dir("static")).ServeHTTP},
	&Route{"POST", "/secret", encodeSecretHandler},
}

type SlashHandler struct {
	Handler http.Handler
}

func (sh *SlashHandler) ServeHTTP(rw http.ResponseWriter, req *http.Request) {
	if req.URL.Path != "/" {
		req.URL.Path = strings.TrimRight(req.URL.Path, "/")
		req.RequestURI = req.URL.RequestURI()
	}

	sh.Handler.ServeHTTP(rw, req)
}

func main() {
	router := mux.NewRouter()
	for _, r := range Routes {
		route := router.NewRoute()
		route.Path(r.Path).Methods(r.Method)
		route.HandlerFunc(r.Handler)
	}

	app := negroni.Classic()
	app.UseHandler(&SlashHandler{router})

	port := os.Getenv("PORT")
	if port == "" {
		port = "3000"
	}

	app.Run(":" + port)
}
