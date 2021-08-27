package main

import (
	"os"
	"strconv"
	"strings"
)

type query struct {
	target    int
	request   string
	path      []step
	attribute string
}

type step struct {
	name  string
	count int
}

const (
	tags = iota
	tagValue
	attr
	attrValue
	empty
)

func parseQuery() query {
	q := getQuery()

	q.path = q.getPath()

	if len(q.path) == 0 {
		return q
	}

	sa := q.getAttribute()
	if len(sa) == 1 {
		return q
	}

	q.path[len(q.path)-1].name = sa[0]

	q.attribute = sa[1]

	return q
}

func toTag(s string) int {
	switch s {
	case "tag":
		return tags
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

func getQuery() query {
	var q query
	args := os.Args[1:]
	switch len(args) {
	case 1:
		q = query{
			target:  tags,
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

func (q *query) getPath() []step {
	if q.request == "." {
		return []step{}
	}

	path := strings.Split(q.request, ".")

	if len(path) > 0 && path[0] == "" {
		path = path[1:]
	}

	steps := []step{}
	for i := range path {
		steps = append(steps, getStep(path[i]))
	}

	return steps
}

func getStep(s string) step {
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
		name = s[:]
	}

	return step{
		name:  name,
		count: count,
	}
}

func (q *query) getAttribute() []string {
	return strings.Split(q.path[len(q.path)-1].name, "#")
}
