package main

import (
	"fmt"

	"github.com/alouca/wikiparse/wikiparse"
)

func main() {
	file := "/Users/alouca/Downloads/enwiki-latest-pages-articles-multistream.xml.bz2"

	parser, err := wikiparse.NewWikiParser(file)
	if err != nil {
		panic(err)
	}

	for {
		page, err := parser.Next()
		if err != nil {
			panic(err)
		}
		if page == nil {
			continue
		}
		fmt.Printf("%v\n\t%+v\n", page.Title, page.Categories)
	}
}
