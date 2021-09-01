package main

import "errors"

type tagParser struct {
	index int
	path  []string
}

func (p *parser) processTags(chunk []byte) {

}

func (p *parser) getTagsList(chunk []byte) ([]string, error) {
	for i := range chunk {
		if p.InsideTag {
			p.CurrentTag.Bytes = append(p.CurrentTag.Bytes, chunk[i])

			if chunk[i] == closeBracket {
				p.InsideTag = false

				tagName, err := getTagName(p.CurrentTag.Bytes)
				if err != nil {
					return nil, err
				}

				if tagName == p.searchQuery.query.path[p.searchQuery.count].name {
					p.searchQuery.count++

					if len(p.searchQuery.query.path) == p.searchQuery.count {

					}
				}
			}

			continue
		}

		if chunk[i] == openBracket {
			p.InsideTag = true
			p.CurrentTag = tag{
				Bytes: []byte{chunk[i]},
			}

			continue
		}
	}

	return nil, nil
}

func getTagName(t []byte) (string, error) {
	if len(t) < 3 {
		return "", errors.New("tag can't be less then 3 bytes")
	}

	if t[0] != openBracket {
		return "", errors.New("tag must start from open bracket symbol")
	}

	startName := 1   // name starts after open bracket
	if t[1] == '/' { // closed tag
		startName = 2
	}

	endName := startName

	for ; endName < len(t)-1; endName++ {
		if t[endName] == ' ' {
			break
		}
	}

	return string(t[startName:endName]), nil
}
