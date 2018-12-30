package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"os"
)

/*
	ROUTES
	/admin/featured/delete/{{link}}
	/admin/featured
	/admin/gallery
*/

const featuredPath = "./featured/"

// FeaturedProperty holds the data of a featured property in a file
type FeaturedProperty struct {
	FileName string `json:"fileName"`
	Address  string `json:"address"`
	Link     string `json:"link"`
}

// Featured holds all of the featured properties
var Featured map[int]FeaturedProperty

// NewFeaturedProperty creates a new featured property object
func NewFeaturedProperty(fileName, address, link string) *FeaturedProperty {
	return &FeaturedProperty{
		FileName: fileName,
		Address:  address,
		Link:     link,
	}
}

func featured(w http.ResponseWriter, req *http.Request) {
	// TODO: check for login first
	if req.Method == http.MethodPost {
		//TODO: add a new featured property here
	}
	tpls.ExecuteTemplate(w, "featuredprop.html", struct {
		Properties       map[int]FeaturedProperty
		PropertiesAmount int
	}{
		Properties:       Featured,
		PropertiesAmount: len(Featured),
	})
}

func (f *FeaturedProperty) store() error {
	if _, err := os.Stat(featuredPath); os.IsNotExist(err) {
		os.MkdirAll(featuredPath, os.ModePerm)
	}

	file, err := os.Create(fmt.Sprintf("%s%s.txt", featuredPath, f.FileName))
	if err != nil {
		return err
	}
	defer file.Close()

	store, err := json.Marshal(f)
	if err != nil {
		return err
	}

	_, err = file.Write(store)
	if err != nil {
		return err
	}

	return nil
}

func (f *FeaturedProperty) remove() error {
	return os.Remove(fmt.Sprintf("%s%s.txt", featuredPath, f.FileName))
}

// GetFeaturedProperties gets all of the featured properties
func GetFeaturedProperties() []FeaturedProperty {
	var fps []FeaturedProperty
	filelist, err := ioutil.ReadDir(featuredPath)
	if err != nil {
		log.Fatalln(err)
	}

	for _, fileinfo := range filelist {
		if fileinfo.Mode().IsRegular() {
			contents, err := ioutil.ReadFile(featuredPath + fileinfo.Name())
			if err != nil {
				log.Fatalln("fatal err:", err)
			}
			// fmt.Println("Bytes read: ", len(bytes))
			// fmt.Println("String read: ", string(contents))
			var q FeaturedProperty
			err = json.Unmarshal(contents, &q)
			if err != nil {
				log.Fatalln("Error unmarshaling featured properties:", err)
			}
			fps = append(fps, q)
		}
	}
	return fps
}

// ToMap takes a list of questions and changes it to a map
func ToMap(fps []FeaturedProperty) map[int]FeaturedProperty {
	m := make(map[int]FeaturedProperty)
	for i, f := range fps {
		m[i+1] = f
	}
	return m
}
