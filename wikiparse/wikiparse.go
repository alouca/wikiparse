package wikiparse

import (
	"compress/bzip2"
	"encoding/xml"
	"errors"
	"io"
	"os"
	"regexp"
	"strings"
)

var (
	categoriesRE = regexp.MustCompile(`\[\[Category:(.*?)\]\]`)
)

type Wikiparse struct {
	file      *os.File
	bz2Reader io.Reader
	decoder   *xml.Decoder
}

type Wikipage struct {
	Title    string `xml:"title"`
	Redirect string `xml:"redirect"`
	Id       int    `xml:"id"`
	Revision struct {
		Id       int    `xml:"id"`
		ParentId int    `xml:"parentid"`
		Text     string `xml:"text"`
		Minor    bool   `xml:"minor"`
		Sha1     string `xml:"sha1"`
	} `xml:"revision"`

	Contributor struct {
		Username string `xml:"username"`
		Id       int    `xml:"id"`
	} `xml:"contributor"`

	Timestamp string `xml:"timestamp"`

	Categories []string
}

func NewWikiParser(filename string) (*Wikiparse, error) {
	// check if the file exists
	if _, err := os.Stat(filename); err != nil {
		return nil, err
	}

	// make sure it ends with bz2
	if filename[len(filename)-3:] != "bz2" {
		return nil, errors.New("file must be bz2")
	}

	file, err := os.Open(filename)
	if err != nil {
		return nil, err
	}

	bz2Reader := bzip2.NewReader(file)
	decoder := xml.NewDecoder(bz2Reader)
	return &Wikiparse{
		file:      file,
		bz2Reader: bz2Reader,
		decoder:   decoder,
	}, nil
}

func (w *Wikiparse) Close() {
	w.file.Close()
}

func (w *Wikiparse) Next() (*Wikipage, error) {
	token, err := w.decoder.Token()
	if err != nil {
		if err == io.EOF {
			return nil, nil
		}
		return nil, err
	}

	if start, ok := token.(xml.StartElement); ok {
		if start.Name.Local == "page" {
			var page Wikipage
			err = w.decoder.DecodeElement(&page, &start)
			if err != nil {
				return nil, err
			}
			page.Categories = w.GetCategories(&page)
			return &page, nil
		}
	}

	return nil, nil
}

func (w *Wikiparse) GetCategories(page *Wikipage) []string {
	matches := categoriesRE.FindAllStringSubmatch(page.Revision.Text, -1)
	categories := make([]string, len(matches))
	for i, match := range matches {
		categories[i] = cleanCategoryText(match[1])
	}
	return categories
}

func cleanCategoryText(category string) string {
	category = strings.TrimSpace(strings.ToLower(category))
	category = strings.Replace(category, "|", "", -1)

	return category
}
