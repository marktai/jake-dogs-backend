package server

import (
	"jake_dogs"
	"net/http"
)

func getImage(w http.ResponseWriter, r *http.Request) {
	image_url := jake_dogs.GetImage()
	http.Redirect(w, r, image_url, 302)
}
