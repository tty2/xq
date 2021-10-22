package tags

import (
	"testing"

	"github.com/stretchr/testify/require"
	"github.com/tty2/xq/internal/domain"
	"github.com/tty2/xq/internal/domain/symbol"
)

func TestUpdatePath(t *testing.T) {
	t.Parallel()

	t.Run("err: incorrect xml", func(t *testing.T) {
		t.Parallel()

		p := Processor{
			currentPath: []string{"1", "2", "3"},
			currentTag: tag{
				closed: true,
				name:   "4",
			},
		}

		rq := require.New(t)

		err := p.updatePath()
		rq.Error(err)
	})

	t.Run("ok: correct close tag", func(t *testing.T) {
		t.Parallel()

		p := Processor{
			currentPath: []string{"1", "2", "3"},
			currentTag: tag{
				closed: true,
				name:   "3",
			},
		}

		rq := require.New(t)

		err := p.updatePath()
		rq.NoError(err)
		rq.Len(p.currentPath, 2)
		rq.Equal("1", p.currentPath[0])
		rq.Equal("2", p.currentPath[1])
	})

	t.Run("ok: new tag", func(t *testing.T) {
		t.Parallel()

		p := Processor{
			currentPath: []string{"1", "2", "3"},
			currentTag: tag{
				name: "4",
			},
		}

		rq := require.New(t)

		err := p.updatePath()
		rq.NoError(err)
		rq.Len(p.currentPath, 4)
		rq.Equal("1", p.currentPath[0])
		rq.Equal("2", p.currentPath[1])
		rq.Equal("3", p.currentPath[2])
		rq.Equal("4", p.currentPath[3])
	})
}

func TestUpdateTagsList(t *testing.T) {
	t.Parallel()

	t.Run("skip: current path less than query", func(t *testing.T) {
		t.Parallel()

		p := Processor{
			query: query{
				path: []domain.Step{
					{
						Name:  "1",
						Index: -1,
					},
					{
						Name:  "2",
						Index: -1,
					},
					{
						Name:  "3",
						Index: -1,
					},
					{
						Name:  "4",
						Index: -1,
					},
				},
			},
			currentPath: []string{"1", "2", "3"},
			currentTag: tag{
				name: "7",
			},
			printList: []string{},
		}

		rq := require.New(t)

		p.updatePrintList()
		rq.Len(p.printList, 0)
	})

	t.Run("skip: current path greater than query", func(t *testing.T) {
		t.Parallel()

		p := Processor{
			query: query{
				path: []domain.Step{
					{
						Name:  "1",
						Index: -1,
					},
					{
						Name:  "2",
						Index: -1,
					},
					{
						Name:  "3",
						Index: -1,
					},
					{
						Name:  "4",
						Index: -1,
					},
				},
			},
			currentPath: []string{"1", "2", "3", "4", "5"},
			currentTag: tag{
				name: "10",
			},
			printList: []string{"6", "7"},
		}

		rq := require.New(t)

		p.updatePrintList()
		rq.Len(p.printList, 2)
	})

	t.Run("skip: current path contains current tag name", func(t *testing.T) {
		t.Parallel()

		p := Processor{
			query: query{
				path: []domain.Step{
					{
						Name:  "1",
						Index: -1,
					},
					{
						Name:  "2",
						Index: -1,
					},
					{
						Name:  "3",
						Index: -1,
					},
					{
						Name:  "4",
						Index: -1,
					},
				},
			},
			currentPath: []string{"1", "2", "3", "4"},
			currentTag: tag{
				name: "7",
			},
			printList: []string{"6", "7"},
		}

		rq := require.New(t)

		p.updatePrintList()
		rq.Len(p.printList, 2)
	})

	t.Run("skip: step back from closed tag: last query name is the same as current tag", func(t *testing.T) {
		t.Parallel()

		p := Processor{
			query: query{
				path: []domain.Step{
					{
						Name:  "1",
						Index: -1,
					},
					{
						Name:  "2",
						Index: -1,
					},
					{
						Name:  "3",
						Index: -1,
					},
					{
						Name:  "4",
						Index: -1,
					},
				},
			},
			currentPath: []string{"1", "2", "3", "4"},
			currentTag: tag{
				name: "4",
			},
			printList: []string{"6", "7"},
		}

		rq := require.New(t)

		p.updatePrintList()
		rq.Len(p.printList, 2)
	})

	t.Run("add tag", func(t *testing.T) {
		t.Parallel()

		p := Processor{
			query: query{
				path: []domain.Step{
					{
						Name:  "1",
						Index: -1,
					},
					{
						Name:  "2",
						Index: -1,
					},
					{
						Name:  "3",
						Index: -1,
					},
					{
						Name:  "4",
						Index: -1,
					},
				},
			},
			currentPath: []string{"1", "2", "3", "4"},
			currentTag: tag{
				name: "8",
			},
			printList: []string{"6", "7"},
		}

		rq := require.New(t)

		p.updatePrintList()
		rq.Len(p.printList, 3)
	})
}

func TestProcessCurrentTag(t *testing.T) {
	t.Parallel()

	t.Run("err: too short", func(t *testing.T) {
		t.Parallel()

		p := Processor{
			currentTag: tag{
				bytes: []byte("<>"),
			},
		}

		rq := require.New(t)

		err := p.processCurrentTag()
		rq.Error(err)
	})

	t.Run("err: don't start from open bracket", func(t *testing.T) {
		t.Parallel()

		p := Processor{
			currentTag: tag{
				bytes: []byte("tagname"),
			},
		}

		rq := require.New(t)

		err := p.processCurrentTag()
		rq.Error(err)
	})

	t.Run("ok: without attributes", func(t *testing.T) {
		t.Parallel()

		p := Processor{
			currentTag: tag{
				bytes: []byte("<tagname>"),
			},
		}

		rq := require.New(t)

		err := p.processCurrentTag()
		rq.NoError(err)
		rq.Equal("tagname", p.currentTag.name)
		rq.False(p.currentTag.closed)
	})

	t.Run("ok: without attributes: close tag", func(t *testing.T) {
		t.Parallel()

		p := Processor{
			currentTag: tag{
				bytes: []byte("</tagname>"),
			},
		}

		rq := require.New(t)

		err := p.processCurrentTag()
		rq.NoError(err)
		rq.Equal("tagname", p.currentTag.name)
		rq.True(p.currentTag.closed)
	})

	t.Run("ok: without attributes: single tag", func(t *testing.T) {
		t.Parallel()

		p := Processor{
			currentTag: tag{
				bytes: []byte("<tagname />"),
			},
		}

		rq := require.New(t)

		err := p.processCurrentTag()
		rq.NoError(err)
		rq.Equal("tagname", p.currentTag.name)
	})
}

func TestAddSymbolIntoTag(t *testing.T) {
	t.Parallel()

	t.Run("open bracket", func(t *testing.T) {
		t.Parallel()

		p := Processor{
			insideTag: true,
			currentTag: tag{ // initialization is the same as in p.proccess second condition with open bracket
				bytes:    []byte{symbol.OpenBracket},
				brackets: 1,
			},
		}

		rq := require.New(t)
		rq.Equal(1, p.currentTag.brackets)

		err := p.addSymbolIntoTag(symbol.OpenBracket)
		rq.NoError(err)
		rq.Equal(2, p.currentTag.brackets)
		rq.True(p.insideTag)
	})

	t.Run("alphabet symbol", func(t *testing.T) {
		t.Parallel()

		p := Processor{
			insideTag: true,
			currentTag: tag{
				bytes:    []byte("<ta"),
				brackets: 1,
			},
		}

		rq := require.New(t)

		err := p.addSymbolIntoTag('g')
		rq.NoError(err)
		rq.Equal("<tag", string(p.currentTag.bytes))
		rq.True(p.insideTag)
	})

	t.Run("comment", func(t *testing.T) {
		t.Parallel()

		p := Processor{
			insideTag: true,
			currentTag: tag{
				bytes:    []byte("<!-- some comment here <b> with tags inside </b"),
				brackets: 2, // <<><  =>  <<
			},
		}

		rq := require.New(t)

		err := p.addSymbolIntoTag(symbol.CloseBracket)
		rq.NoError(err)
		rq.Equal("<!-- some comment here <b> with tags inside </b>", string(p.currentTag.bytes))
		rq.True(p.insideTag)
	})

	t.Run("alphabet symbol", func(t *testing.T) {
		t.Parallel()

		p := Processor{
			insideTag:   true,
			currentPath: []string{"tags"},
			currentTag: tag{
				bytes:    []byte("<tag"),
				brackets: 1,
			},
			query: query{
				path: []domain.Step{ // query.path can't be empty: it is checked on init processor
					{
						Name:  "tags",
						Index: -1,
					},
				},
			},
		}

		rq := require.New(t)

		rq.Equal("", p.currentTag.name)
		rq.Len(p.printList, 0)
		rq.Len(p.currentPath, 1)

		err := p.addSymbolIntoTag(symbol.CloseBracket)
		rq.NoError(err)
		rq.Equal("<tag>", string(p.currentTag.bytes))
	})

	t.Run("service tag", func(t *testing.T) {
		t.Parallel()

		p := Processor{
			insideTag: true,
			currentTag: tag{
				bytes:    []byte("<!-- some comment here <b> with tags inside </b>"),
				brackets: 1, // <
			},
		}

		rq := require.New(t)

		err := p.addSymbolIntoTag(symbol.CloseBracket)
		rq.NoError(err)
		rq.Equal("<!-- some comment here <b> with tags inside </b>>", string(p.currentTag.bytes))
		rq.Equal(0, p.currentTag.brackets)
	})

	t.Run("single tag", func(t *testing.T) {
		t.Parallel()

		p := Processor{
			insideTag: true,
			currentTag: tag{
				bytes:    []byte(`<tagname attr="value" /`),
				brackets: 1, // <
			},
			query: query{
				path: []domain.Step{
					{
						Name:  "tagname",
						Index: -1,
					},
				},
			},
		}

		rq := require.New(t)

		err := p.addSymbolIntoTag(symbol.CloseBracket)
		rq.NoError(err)
		rq.Equal(`<tagname attr="value" />`, string(p.currentTag.bytes))
		rq.Equal(0, p.currentTag.brackets)
	})
}

func TestProcess(t *testing.T) {
	t.Parallel()

	t.Run("ok", func(t *testing.T) {
		t.Parallel()
		rq := require.New(t)

		p := Processor{
			query: query{
				path: []domain.Step{},
			},
		}

		err := p.process([]byte(`t</tagname attr="value"`))
		rq.NoError(err)
	})
}

func TestNewProcessor(t *testing.T) {
	t.Parallel()

	t.Run("error: empty path", func(t *testing.T) {
		t.Parallel()
		rq := require.New(t)

		p, err := NewProcessor([]domain.Step{}, "", domain.TagList)
		rq.Error(err)
		rq.Nil(p)
	})

	t.Run("ok", func(t *testing.T) {
		t.Parallel()
		rq := require.New(t)

		p, err := NewProcessor([]domain.Step{
			{
				Name:  "tagname",
				Index: -1,
			},
		}, "test", domain.TagList)

		rq.NoError(err)
		rq.Len(p.query.path, 1)
		rq.Equal("test", p.query.attribute)
		rq.Equal(domain.TagList, p.query.searchType)
	})
}

func TestPickAttributesNames(t *testing.T) {
	t.Parallel()

	t.Run("nil: empty tag", func(t *testing.T) {
		t.Parallel()

		tag := []byte{}
		res := pickAttributesNames(tag)

		rq := require.New(t)
		rq.Nil(res)
	})

	t.Run("nil: first byte is not an open bracket", func(t *testing.T) {
		t.Parallel()

		tag := []byte("tagname")
		res := pickAttributesNames(tag)

		rq := require.New(t)
		rq.Nil(res)
	})

	t.Run("ok", func(t *testing.T) {
		t.Parallel()

		tag := []byte("<tagname attr1='value1' attr2='value2' attr3='value3'")
		res := pickAttributesNames(tag)

		rq := require.New(t)
		rq.Len(res, 3)
		rq.Equal("attr1", res[0])
		rq.Equal("attr2", res[1])
		rq.Equal("attr3", res[2])
	})

	t.Run("ok: other quotes", func(t *testing.T) {
		t.Parallel()

		tag := []byte(`<tagname attr1="value1" attr2="value2" attr3="value3"`)
		res := pickAttributesNames(tag)

		rq := require.New(t)
		rq.Len(res, 3)
		rq.Equal("attr1", res[0])
		rq.Equal("attr2", res[1])
		rq.Equal("attr3", res[2])
	})
}
