package main

import "testing"

func TestCreateGalleryImage(t *testing.T) {
	c := NewImageCard("test4", "fourthfile", "image4")
	err := c.store()
	if err != nil {
		t.Logf("an error occurred creating a new prop %v\n", err)
		t.Fail()
	}
}
