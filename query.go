package main

import (
	"os"
	"strconv"
	"strings"

	"github.com/tty2/xq/internal/domain"
)

type query struct {
	target    int
	request   string
	path      []domain.Step
	attribute string
}

const (
	_tags = iota
	tagValue
	attr
	attrValue
	empty
)

func getQuery() query {
	var q query
	args := os.Args[1:]
	switch len(args) {
	case 1:
		q = query{
			target:  _tags,
			request: args[0],
		}
	case 2:
		q = query{
			target:  toTag(args[0]),
			request: args[1],
		}
	default:
		q = query{
			target:  empty,
			request: ".",
		}
	}

	return q
}

func (q *query) parse() {
	q.path = q.getPath()

	if len(q.path) == 0 {
		return
	}

	sa := q.separateAttribute()
	if len(sa) == 1 {
		return
	}

	q.path[len(q.path)-1].Name = sa[0]

	q.attribute = sa[1]
}

func toTag(s string) int {
	switch s {
	case "tag":
		return _tags
	case "value":
		return tagValue
	case "attribute":
		return attr
	case "aValue":
		return attrValue
	default:
		return empty
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
		steps = append(steps, getStep(path[i]))
	}

	return steps
}

func getStep(s string) domain.Step {
	var name string
	var inBrackets bool
	var num []byte
	for i := range s {
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
		name = s
	}

	return domain.Step{
		Name:  name,
		Index: count,
	}
}

func (q *query) separateAttribute() []string {
	return strings.Split(q.path[len(q.path)-1].Name, "#")
}
