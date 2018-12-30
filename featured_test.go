package main

import "testing"

func TestCreateFeatured(t *testing.T) {
	f := NewFeaturedProperty("file3", "200 gram", "bham2")
	err := f.store()
	if err != nil {
		t.Logf("an error occurred creating a new prop %v\n", err)
		t.Fail()
	}
}
