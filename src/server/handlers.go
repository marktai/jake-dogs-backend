package server

import (
	"corgi"
	"net/http"
)

func getImage(w http.ResponseWriter, r *http.Request) {
	image_url := corgi.GetImage()
	http.Redirect(w, r, image_url, 302)
}

func getGif(w http.ResponseWriter, r *http.Request) {
	gif_url := corgi.GetGif()
	http.Redirect(w, r, gif_url, 302)
}
