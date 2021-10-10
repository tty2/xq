package main

import (
	"os"
	"strconv"
	"strings"

	"github.com/tty2/xq/internal/domain"
)

type query struct {
	request    string
	path       []domain.Step
	attribute  string
	searchType domain.SearchType
}

func getQuery() query {
	var q query
	args := os.Args[1:]
	switch len(args) {
	case 1:
		q = query{
			request: args[0],
		}
		q.searchType = domain.TagValue
	case 2:
		q = query{
			request: args[1],
		}
		switch args[0] {
		case "tags":
			q.searchType = domain.TagList
		case "attr":
			q.searchType = domain.AttrList
		}
	default:
		q = query{
			request:    ".",
			searchType: domain.TagValue,
		}
	}

	return q
}

func (q *query) parse() {
	q.path = q.getPath()
	if len(q.path) == 0 {
		return
	}

	q.attribute = q.getAttribute()
	if q.attribute != "" {
		q.searchType = domain.AttrValue
	}
}

func (q *query) getPath() []domain.Step {
	if q.request == "." {
		return []domain.Step{}
	}

	path := strings.Split(q.request, ".")

	if len(path) > 0 && path[0] == "" {
		path = path[1:]
	}

	steps := []domain.Step{}
	for i := range path {
		step := getStep(path[i])
		steps = append(steps, step)
	}

	return steps
}

func getStep(s string) domain.Step {
	var name string
	var inBrackets bool
	var num []byte
	var i int
	for ; i < len(s); i++ {
		if s[i] == '#' {
			break
		}
		if inBrackets {
			if s[i] == ']' {
				break
			}
			num = append(num, s[i])
		}
		if s[i] == '[' {
			name = s[:i]
			inBrackets = true
		}
	}

	count, err := strconv.Atoi(string(num))
	if err != nil {
		count = -1
		name = s[:i]
	}

	return domain.Step{
		Name:  name,
		Index: count,
	}
}

func (q *query) getAttribute() string {
	sa := strings.Split(q.request, "#")
	if len(sa) != 2 {
		return ""
	}

	return sa[1]
}
