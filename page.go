package main

import (
	"cmp"
	"io/fs"
	"os"
	"slices"
	"strings"
)

type Page struct {
	Title string
	Body  []byte
}

func (p *Page) save() error {
	filename := "./content/" + p.Title + ".txt"
	return os.WriteFile(filename, p.Body, 0600)
}

func loadPage(title string) (*Page, error) {
	filename := "./content/" + title + ".txt"
	body, err := os.ReadFile(filename)
	if err != nil {
		return nil, err
	}
	return &Page{Title: title, Body: body}, nil
}

func loadAllPages(search string) ([]Page, error) {
	var pageList []Page
	path := "./content/"
	fileSystem := os.DirFS(path)

	err := fs.WalkDir(fileSystem, ".", func(path string, d fs.DirEntry, err error) error {
		if !d.IsDir() {
			p := &Page{}
			p, err := loadPage(strings.TrimSuffix(d.Name(), ".txt"))
			if err != nil {
				return err
			}
			if strings.Contains(strings.ToLower(p.Title), strings.ToLower(search)) {
				pageList = append(pageList, *p)
			}

		}
		return nil
	})

	if err != nil {
		return nil, err
	}
	slices.SortFunc(pageList, func(a, b Page) int {
		return cmp.Compare(strings.ToLower(a.Title), strings.ToLower(b.Title))
	})

	return pageList, nil
}
