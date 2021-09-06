package main

import (
	"os"
	"strconv"
	"strings"

	"github.com/tty2/xq/internal/domain"
)

type query struct {
	request   string
	path      []domain.Step
	attribute string
}

func getQuery() query {
	var q query
	args := os.Args[1:]
	switch len(args) {
	case 1:
		q = query{
			request: args[0],
		}
	case 2:
		q = query{
			request: args[1],
		}
	default:
		q = query{
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
