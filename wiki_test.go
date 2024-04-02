package main

import (
	"testing"
)

func TestPageSaveLoadDelete(t *testing.T) {
	p := Page{Title: "testWiki", Body: []byte("test wiki body")}
	p.save()

	p1, err := loadPage("testWiki")

	if err != nil {
		t.Errorf("%s", err.Error())
	}

	if string(p1.Body) != string(p.Body) {
		t.Errorf("Body is incorrect saved: %s loaded: %s", string(p.Body), string(p1.Body))
	}

	if p1.Title != p.Title {
		t.Errorf("Title is incorrect saved: %s loaded: %s", string(p.Title), string(p1.Title))
	}

	err = deletePage("testWiki")

	if err != nil {
		t.Errorf("Error when deleting Wikipage")
	}
}
