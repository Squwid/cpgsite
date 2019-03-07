package main

import (
	"encoding/json"
	"fmt"
	"io"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"strings"

	uuid "github.com/satori/go.uuid"
)

// The two different card types that go into the gallery
const (
	CardImage     = "card-image"
	CardImageText = "card-image-text"
	CardText      = "card-text"
)

const cardPath = "./cards/"
const imagePath = "./assets/gallery/"

// Card holds the structure for each card inside of the gallery
type Card struct {
	ID            string
	Type          string
	TypeFormatted string
	Title         string
	TextTitle     string
	Paragraph     string
	Link          string

	ShowText  bool
	ShowImage bool
	Class     string
}

// GalleryImages holds a map of numbers to images
var GalleryImages map[int]Card

// NewImageCard returns a pointer to an image card
func NewImageCard(uuid, title, link string) *Card {
	return &Card{
		Type:          CardImage,
		ID:            uuid,
		Title:         title,
		Link:          link,
		ShowImage:     true,
		ShowText:      false,
		Class:         "card-img",
		TypeFormatted: "Image",
	}
}

// NewTextCard creates a card that only shows text
func NewTextCard(uuid, title, textTitle, textParagraph string) *Card {
	return &Card{
		Type:          CardText,
		ID:            uuid,
		Title:         title,
		ShowImage:     false,
		ShowText:      true,
		TextTitle:     textTitle,
		Paragraph:     textParagraph,
		TypeFormatted: "Text",
	}
}

// NewImageTextCard creates a new card with text and an image
func NewImageTextCard(uuid, title, link, paragraph, textTitle string) *Card {
	return &Card{
		Type:          CardImageText,
		ID:            uuid,
		Title:         title,
		Link:          link,
		Paragraph:     paragraph,
		TextTitle:     textTitle,
		ShowImage:     true,
		ShowText:      true,
		Class:         "card-img-top",
		TypeFormatted: "Image w/Text",
	}
}

func gallery(w http.ResponseWriter, req *http.Request) {
	tpls.ExecuteTemplate(w, "gallery.html", struct {
		Images map[int]Card
	}{
		Images: GalleryImages,
	})
}

func adminGallery(w http.ResponseWriter, req *http.Request) {
	if !loggedIn(req) {
		http.Redirect(w, req, "/login", http.StatusSeeOther)
		return
	}
	tpls.ExecuteTemplate(w, "admingallery.html", struct {
		PropertiesAmount int
		Images           map[int]Card
		ImageAmount      int
	}{
		PropertiesAmount: len(Featured),
		Images:           GalleryImages,
		ImageAmount:      len(GalleryImages),
	})
}

func addImage(w http.ResponseWriter, req *http.Request) {
	if !loggedIn(req) {
		http.Redirect(w, req, "/login", http.StatusSeeOther)
		return
	}
	if req.Method == http.MethodPost {
		cardType := req.FormValue("imagetype")
		imageTitle := req.FormValue("imagetitle")
		textTitle := req.FormValue("texttitle")
		textParagraph := req.FormValue("textparagraph")
		// fmt.Printf("%s\n%s\n%s\n%s\n", cardType, imageTitle, textTitle, textParagraph)

		// Check and create card types
		if cardType == CardImage {
			// Now get the image file
			file, fileHead, err := req.FormFile("fileinput")
			if err != nil {
				logger.Printf("error opening image file %v\n", err)
				http.Redirect(w, req, "/admin/gallery", http.StatusSeeOther)
				return
			}
			defer file.Close()
			f, err := os.OpenFile(imagePath+fileHead.Filename, os.O_WRONLY|os.O_CREATE, 0666)
			if err != nil {
				logger.Printf("error storing file: %v\n", err)
				http.Redirect(w, req, "/admin/gallery", http.StatusSeeOther)
				return
			}
			defer f.Close()
			io.Copy(f, file)
			uuid, err := uuid.NewV4()
			if err != nil {
				panic(err) // todo: why would this error
			}
			c := NewImageCard(uuid.String(), imageTitle, imagePath+fileHead.Filename)
			err = c.store()
			if err != nil {
				logger.Printf("error storing card: %v\n", err)
				return
			}
			GalleryImages = CardsToMap(GetImages())
		} else if cardType == CardImageText {
			file, fileHead, err := req.FormFile("fileinput")
			if err != nil {
				logger.Printf("error opening image file %v\n", err)
				http.Redirect(w, req, "/admin/gallery", http.StatusSeeOther)
				return
			}
			defer file.Close()
			f, err := os.OpenFile(imagePath+fileHead.Filename, os.O_WRONLY|os.O_CREATE, 0666)
			if err != nil {
				logger.Printf("error storing file: %v\n", err)
				http.Redirect(w, req, "/admin/gallery", http.StatusSeeOther)
				return
			}
			defer f.Close()
			io.Copy(f, file)
			uuid, err := uuid.NewV4()
			if err != nil {
				panic(err)
			}
			c := NewImageTextCard(uuid.String(), imageTitle, imagePath+fileHead.Filename, textParagraph, textTitle)
			err = c.store()
			if err != nil {
				logger.Printf("error storing card: %v\n", err)
				return
			}
			GalleryImages = CardsToMap(GetImages())
		} else if cardType == CardText {
			uuid, err := uuid.NewV4()
			if err != nil {
				panic(err)
			}
			c := NewTextCard(uuid.String(), imageTitle, textTitle, textParagraph)
			err = c.store()
			if err != nil {
				logger.Printf("error storing card: %v\n", err)
				return
			}
			GalleryImages = CardsToMap(GetImages())
		}
	}
	http.Redirect(w, req, "/admin/gallery", http.StatusSeeOther)
}

func deleteImage(w http.ResponseWriter, req *http.Request) {
	if !loggedIn(req) {
		http.Redirect(w, req, "/login", http.StatusSeeOther)
		return
	}
	s := strings.Split(req.URL.Path, "/")
	// fmt.Println(req.URL.Path)
	path := s[4]
	removeImageByID(path)
	GalleryImages = CardsToMap(GetImages())
	http.Redirect(w, req, "/admin/gallery", http.StatusSeeOther)
	return
}

func removeImageByID(id string) {
	logger.Printf("trying to remove image %s\n", id)
	for _, c := range GalleryImages {
		if c.ID == id {
			err1, err2 := c.remove()
			if err1 != nil {
				panic(err1)
			}
			if err2 != nil {
				panic(err2)
			}
			logger.Printf("successfully removed image %s\n", id)
		}
	}
}

func (c *Card) store() error {
	if _, err := os.Stat(cardPath); os.IsNotExist(err) {
		os.MkdirAll(cardPath, os.ModePerm)
	}

	file, err := os.Create(fmt.Sprintf("%s%s.json", cardPath, c.ID))
	if err != nil {
		return err
	}
	defer file.Close()

	store, err := json.Marshal(c)
	if err != nil {
		return err
	}

	_, err = file.Write(store)
	if err != nil {
		return err
	}

	return nil
}

func (c *Card) remove() (error, error) {
	if c.Link == "" {
		return nil, os.Remove(fmt.Sprintf("%s%s.json", cardPath, c.ID))
	}
	return os.Remove(c.Link), os.Remove(fmt.Sprintf("%s%s.json", cardPath, c.ID))
}

// GetImages gets all of the images from the
func GetImages() []Card {
	var cards []Card
	filelist, err := ioutil.ReadDir(cardPath)
	if err != nil {
		log.Fatalln(err)
	}

	for _, fileinfo := range filelist {
		if fileinfo.Mode().IsRegular() {
			contents, err := ioutil.ReadFile(cardPath + fileinfo.Name())
			if err != nil {
				log.Fatalln("fatal err:", err)
			}
			// fmt.Println("Bytes read: ", len(bytes))
			// fmt.Println("String read: ", string(contents))
			var c Card
			err = json.Unmarshal(contents, &c)
			if err != nil {
				log.Fatalln("Error unmarshaling cards:", err)
			}
			cards = append(cards, c)
		}
	}
	return cards
}

// CardsToMap creates a map of ints to cards from a list of cards
func CardsToMap(cards []Card) map[int]Card {
	m := make(map[int]Card)
	for i, c := range cards {
		m[i+1] = c
	}
	return m
}
