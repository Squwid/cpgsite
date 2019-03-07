package main

import (
	"html/template"
	"log"
	"net/http"
	"os"
)

/*
	DEV NOTES
	only page needing to be done would be the email service for inquiries back to Ian
	maybe email templates? Definitely need to format the emails being sent at some point
	Definitely need to change email name senders to something on our end. Need to contact AWS
	Change emails from mine to ians
*/

var (
	logger *log.Logger
	tpls   *template.Template
)

func init() {
	tpls = template.Must(template.ParseGlob("views/*"))
	logger = log.New(os.Stdout, "[CPG] ", log.Ldate|log.Ltime)

	// init featured properties
	Featured = make(map[int]FeaturedProperty)
	Featured = ToMap(GetFeaturedProperties())

	// init gallery
	GalleryImages = make(map[int]Card)
	GalleryImages = CardsToMap(GetImages())
}

func main() {
	http.Handle("/assets/", http.StripPrefix("/assets", http.FileServer(http.Dir("./assets"))))

	http.HandleFunc("/", index)
	http.HandleFunc("/gallery", gallery)
	http.HandleFunc("/socialmedia", comingsoon)
	http.HandleFunc("/aboutus", aboutUs)
	http.HandleFunc("/contact", comingsoon)
	http.HandleFunc("/login", login)

	// Admin Functions
	http.HandleFunc("/admin/featured", featured)
	http.HandleFunc("/admin/gallery/", adminGallery)
	http.HandleFunc("/admin/gallery/addimage", addImage)
	http.HandleFunc("/admin/gallery/delete/", deleteImage)

	http.ListenAndServe(":80", nil)
}
