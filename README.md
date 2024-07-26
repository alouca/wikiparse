# wikiparse
A simple golang library to parse Wikipedia XML dumps

## Usage Example

```
package main

import (
	"flag"
	"fmt"

	"github.com/alouca/wikiparse/wikiparse"
)

var (
	file = flag.String("file", "", "path to the file")
)

func init() {
	flag.Parse()
}

func main() {
	parser, err := wikiparse.NewWikiParser(*file)
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

```
