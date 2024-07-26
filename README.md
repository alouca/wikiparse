# wikiparse
A simple golang library to parse Wikipedia XML dumps

## Usage Example

`
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
`