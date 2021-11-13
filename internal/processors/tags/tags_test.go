package tags

import (
	"testing"

	"github.com/stretchr/testify/require"
	"github.com/tty2/xq/internal/domain"
	"github.com/tty2/xq/internal/domain/symbol"
)

func TestSetName(t *testing.T) {
	t.Parallel()

	rq := require.New(t)

	t.Run("err: short", func(t *testing.T) {
		t.Parallel()

		tg := tag{
			bytes: []byte("<>"),
		}

		err := tg.setName()
		rq.Error(err)
	})

	t.Run("err: not a tag: there are no open bracket", func(t *testing.T) {
		t.Parallel()

		tg := tag{
			bytes: []byte("tagname>"),
		}

		err := tg.setName()
		rq.Error(err)
	})

	t.Run("err: not a tag: there are no close bracket", func(t *testing.T) {
		t.Parallel()

		tg := tag{
			bytes: []byte("<tagname"),
		}

		err := tg.setName()
		rq.Error(err)
	})

	t.Run("ok: no attributes", func(t *testing.T) {
		t.Parallel()

		tg := tag{
			bytes: []byte("<tagname>"),
		}

		err := tg.setName()
		rq.NoError(err)

		rq.Equal("tagname", tg.name)
	})

	t.Run("ok: close tag", func(t *testing.T) {
		t.Parallel()

		tg := tag{
			bytes: []byte("</tagname>"),
		}

		err := tg.setName()
		rq.NoError(err)

		rq.Equal("tagname", tg.name)
	})

	t.Run("ok: with attributes", func(t *testing.T) {
		t.Parallel()

		tg := tag{
			bytes: []byte("<tagname attr='value'>"),
		}

		err := tg.setName()
		rq.NoError(err)

		rq.Equal("tagname", tg.name)
	})
}

func TestDecrementPath(t *testing.T) {
	t.Parallel()

	rq := require.New(t)

	t.Run("short", func(t *testing.T) {
		t.Parallel()

		p := Processor{
			currentPath: []string{},
			currentTag: tag{
				closed: true,
				name:   "4",
			},
		}

		err := p.decrementPath()
		rq.Nil(err)
	})

	t.Run("err: incorrect xml", func(t *testing.T) {
		t.Parallel()

		p := Processor{
			currentPath: []string{"1", "2", "3"},
			currentTag: tag{
				closed: true,
				name:   "4",
			},
		}

		err := p.decrementPath()
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

		err := p.decrementPath()
		rq.NoError(err)
		rq.Len(p.currentPath, 2)
		rq.Equal("1", p.currentPath[0])
		rq.Equal("2", p.currentPath[1])
	})
}

func TestTagInQueryPath(t *testing.T) {
	t.Parallel()
	rq := require.New(t)

	t.Run("false: short current path", func(t *testing.T) {
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
				},
			},
			currentPath: []string{"1", "2", "3"},
		}

		ok := p.tagInQueryPath()
		rq.False(ok)
	})

	t.Run("false: long current path", func(t *testing.T) {
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
				},
			},
			currentPath: []string{"1", "2", "3", "4", "5"},
		}

		ok := p.tagInQueryPath()
		rq.False(ok)
	})

	t.Run("false: different path", func(t *testing.T) {
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
				},
			},
			currentPath: []string{"1", "2", "5", "4"},
		}

		ok := p.tagInQueryPath()
		rq.False(ok)
	})

	t.Run("true", func(t *testing.T) {
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
				},
			},
			currentPath: []string{"1", "2", "3", "4"},
		}

		ok := p.tagInQueryPath()
		rq.True(ok)
	})
}

func TestUpdatePrintListForTags(t *testing.T) {
	t.Parallel()

	searchType := domain.TagList
	rq := require.New(t)

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
				searchType: searchType,
			},
			currentPath: []string{"1", "2", "3"},
			printList:   []string{},
		}

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
				searchType: searchType,
			},
			currentPath: []string{"1", "2", "3", "4", "5", "7"},
			printList:   []string{"6", "7"},
		}

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
				searchType: searchType,
			},
			currentPath: []string{"1", "2", "3", "4", "7"},
			currentTag: tag{
				name: "7",
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
				searchType: domain.TagList,
			},
			currentPath: []string{"1", "2", "3", "4", "8"},
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

func TestUpdatePrintListForAttrList(t *testing.T) {
	t.Parallel()

	searchType := domain.AttrList
	rq := require.New(t)

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
				searchType: searchType,
			},
			currentPath: []string{"1", "2", "3"},
			printList:   []string{},
			currentTag: tag{
				bytes: []byte("<tagname attr1='value1' attr2=>"),
			},
		}

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
				searchType: searchType,
			},
			currentPath: []string{"1", "2", "3", "4", "5"},
			printList:   []string{"6", "7"},
			currentTag: tag{
				bytes: []byte("<tagname attr1='value1' attr2=>"),
			},
		}

		p.updatePrintList()
		rq.Len(p.printList, 2)
	})

	t.Run("skip: different path", func(t *testing.T) {
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
				searchType: searchType,
			},
			currentPath: []string{"1", "2", "3", "5"},
			printList:   []string{"6", "7"},
			currentTag: tag{
				bytes: []byte("<tagname attr1='value1' attr2=>"),
			},
		}

		p.updatePrintList()
		rq.Len(p.printList, 2)
	})

	t.Run("ok", func(t *testing.T) {
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
				searchType: searchType,
			},
			currentPath: []string{"1", "2", "3", "4"},
			printList:   []string{"attr1", "attr3"},
			currentTag: tag{
				bytes: []byte("<tagname attr1='value1' attr2='value2'>"),
			},
		}

		p.updatePrintList()
		rq.Len(p.printList, 3)
	})
}

func TestUpdatePrintListForAttrValue(t *testing.T) {
	t.Parallel()

	searchType := domain.AttrValue
	rq := require.New(t)

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
				attribute:  "attr1",
				searchType: searchType,
			},
			currentPath: []string{"1", "2", "3"},
			printList:   []string{},
			currentTag: tag{
				bytes: []byte("<tagname attr1='value1' attr2='value2'>"),
			},
		}

		p.updatePrintList()
		rq.Len(p.printList, 0)
	})

	t.Run("skip: there are no attribute", func(t *testing.T) {
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
				attribute:  "attr3",
				searchType: searchType,
			},
			currentPath: []string{"1", "2", "3", "4"},
			printList:   []string{},
			currentTag: tag{
				bytes: []byte("<tagname attr1='value1' attr2='value2'>"),
			},
		}

		p.updatePrintList()
		rq.Len(p.printList, 0)
	})

	t.Run("skip: empty string", func(t *testing.T) {
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
				attribute:  "attr1",
				searchType: searchType,
			},
			currentPath: []string{"1", "2", "3", "4"},
			printList:   []string{},
			currentTag: tag{
				bytes: []byte("<tagname attr1='' attr2='value2'>"),
			},
		}

		p.updatePrintList()
		rq.Len(p.printList, 0)
	})

	t.Run("ok", func(t *testing.T) {
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
				attribute:  "attr1",
				searchType: searchType,
			},
			currentPath: []string{"1", "2", "3", "4"},
			printList:   []string{},
			currentTag: tag{
				bytes: []byte("<tagname attr1='value1' attr2='value2'>"),
			},
		}

		p.updatePrintList()
		rq.Len(p.printList, 1)
		rq.Equal("value1", p.printList[0])
	})
}

func TestUpdatePrintListForTagValue(t *testing.T) {
	t.Parallel()

	searchType := domain.TagValue
	rq := require.New(t)

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
				attribute:  "attr1",
				searchType: searchType,
			},
			currentPath: []string{"1", "2", "3"},
			printList:   []string{},
			currentTag: tag{
				bytes: []byte("<tagname attr1='value1' attr2='value2'>"),
			},
		}

		p.updatePrintList()
		rq.Len(p.printList, 0)
		rq.Equal(0, p.indentation)
	})

	t.Run("ok: the same length", func(t *testing.T) {
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
				attribute:  "attr1",
				searchType: searchType,
			},
			currentPath: []string{"1", "2", "3", "4"},
			printList:   []string{},
			currentTag: tag{
				bytes: []byte("<tagname attr1='value1' attr2='value2'>"),
			},
		}

		p.updatePrintList()
		rq.Len(p.printList, 1)
		rq.Equal("<tagname attr1='value1' attr2='value2'>", p.printList[0])
		rq.Equal(0, p.indentation)
	})

	t.Run("ok: 2 items deeper", func(t *testing.T) {
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
				attribute:  "attr1",
				searchType: searchType,
			},
			currentPath: []string{"1", "2", "3", "4", "5", "6"},
			printList:   []string{},
			currentTag: tag{
				bytes: []byte("<tagname attr1='value1' attr2='value2'>"),
			},
		}

		p.updatePrintList()
		rq.Len(p.printList, 1)
		rq.Equal("    <tagname attr1='value1' attr2='value2'>", p.printList[0])
		rq.Equal(2, p.indentation)
	})
}

func TestSkip(t *testing.T) {
	t.Parallel()
	rq := require.New(t)

	t.Run("true: xml tag", func(t *testing.T) {
		t.Parallel()

		p := Processor{
			currentTag: tag{
				bytes: []byte(`<?xml version="1.0" encoding="UTF-8"?>`),
			},
		}

		rq.True(p.skip())
	})

	t.Run("true: commend", func(t *testing.T) {
		t.Parallel()

		p := Processor{
			currentTag: tag{
				bytes: []byte(`<!--2021.06.14 03:07:43-->`),
			},
		}

		rq.True(p.skip())
	})

	t.Run("false", func(t *testing.T) {
		t.Parallel()

		p := Processor{
			currentTag: tag{
				bytes: []byte(`<tagname attr1='value1' attr2='value2'>`),
			},
		}

		rq.False(p.skip())
	})
}

func TestCurrentTagIsSingle(t *testing.T) {
	t.Parallel()
	rq := require.New(t)

	t.Run("false: too short", func(t *testing.T) {
		t.Parallel()

		p := Processor{
			currentTag: tag{
				bytes: []byte(`<b>`),
			},
		}

		rq.False(p.currentTagIsSingle())
	})

	t.Run("false", func(t *testing.T) {
		t.Parallel()

		p := Processor{
			currentTag: tag{
				bytes: []byte(`<tagname>`),
			},
		}

		rq.False(p.currentTagIsSingle())
	})

	t.Run("false: close", func(t *testing.T) {
		t.Parallel()

		p := Processor{
			currentTag: tag{
				bytes: []byte(`</b>`),
			},
		}

		rq.False(p.currentTagIsSingle())
	})

	t.Run("true", func(t *testing.T) {
		t.Parallel()

		p := Processor{
			currentTag: tag{
				bytes: []byte(`<tagname/>`),
			},
		}

		rq.True(p.currentTagIsSingle())
	})

	t.Run("true", func(t *testing.T) {
		t.Parallel()

		p := Processor{
			currentTag: tag{
				bytes: []byte(`<tagname />`),
			},
		}

		rq.True(p.currentTagIsSingle())
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

	t.Run("ok", func(t *testing.T) {
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
			query: query{
				searchType: domain.TagList,
				path: []domain.Step{
					{
						Name:  "1",
						Index: -1,
					},
				},
			},
			currentTag: tag{
				bytes: []byte("<tagname />"),
			},
			currentPath: []string{"1"},
			printList:   []string{},
		}

		rq := require.New(t)

		err := p.processCurrentTag()
		rq.NoError(err)
		rq.Equal("tagname", p.currentTag.name)
		rq.Len(p.printList, 1)
		rq.Equal("tagname", p.printList[0])
		rq.Len(p.currentPath, 1)
		rq.Equal("1", p.currentPath[0])
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

		p.addSymbolIntoTag(symbol.OpenBracket)
		rq.Equal(2, p.currentTag.brackets)
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

		p.addSymbolIntoTag('g')
		rq.Equal("<tag", string(p.currentTag.bytes))
		rq.Equal(1, p.currentTag.brackets)
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

		p.addSymbolIntoTag(symbol.CloseBracket)
		rq.Equal("<!-- some comment here <b> with tags inside </b>", string(p.currentTag.bytes))
		rq.Equal(1, p.currentTag.brackets)
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

		p.addSymbolIntoTag(symbol.CloseBracket)
		rq.Equal("<tag>", string(p.currentTag.bytes))
		rq.Equal(0, p.currentTag.brackets)
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

func TestProcess(t *testing.T) {
	t.Parallel()

	t.Run("tagname", func(t *testing.T) {
		t.Parallel()
		rq := require.New(t)

		p := Processor{
			query: query{
				path:       []domain.Step{},
				searchType: domain.TagList,
			},
		}

		err := p.process([]byte(`attr0="value0"><tagname attr="value"><!--comment--></tag attr="invalid tag name">`))
		rq.Error(err)
		rq.Len(p.printList, 2)
		rq.Equal("tagname", p.printList[0])
		rq.Equal("tag", p.printList[1])
	})

	t.Run("ok", func(t *testing.T) {
		t.Parallel()
		rq := require.New(t)

		p := Processor{
			query: query{
				path:       []domain.Step{},
				searchType: domain.TagValue,
			},
		}

		err := p.process([]byte("<tagname attr='value'>\ndata</tag attr='invalid tag name'>"))
		rq.Error(err)
		rq.Len(p.printList, 3)
		rq.Equal(`  <tagname attr='value'>`, p.printList[0])
		rq.Equal(`    data`, p.printList[1])
		rq.Equal(`  </tag attr='invalid tag name'>`, p.printList[2])
	})
}

func TestIntoQueryPath(t *testing.T) {
	t.Parallel()
	rq := require.New(t)

	t.Run("false: too short", func(t *testing.T) {
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
		}

		rq.False(p.intoQueryPath())
	})

	t.Run("false: different", func(t *testing.T) {
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
			currentPath: []string{"1", "2", "8", "4", "5"},
		}

		rq.False(p.intoQueryPath())
	})

	t.Run("true: the same", func(t *testing.T) {
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
		}

		rq.True(p.intoQueryPath())
	})

	t.Run("true: greater", func(t *testing.T) {
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
		}

		rq.True(p.intoQueryPath())
	})
}

// nolint lll: there are long strings on purpose
func TestProcessWithIndex(t *testing.T) {
	t.Parallel()

	t.Run("tag value: first", func(t *testing.T) {
		t.Parallel()
		rq := require.New(t)

		p := Processor{
			query: query{
				path: []domain.Step{
					{
						Name:  "objects",
						Index: -1,
					},
					{
						Name:  "object",
						Index: 0,
					},
				},
				searchType: domain.TagValue,
			},
			index: index{
				set: true,
			},
		}

		err := p.process([]byte(`<objects><object><tg></tg></object><object><tg1></tg1></object><object><tg2></tg2></object></objects>`))
		rq.NoError(err)
		rq.Len(p.printList, 4)
		rq.Equal("<object>", p.printList[0])
		rq.Equal("  <tg>", p.printList[1])
		rq.Equal("  </tg>", p.printList[2])
		rq.Equal("</object>", p.printList[3])
	})

	t.Run("tag value: second", func(t *testing.T) {
		t.Parallel()
		rq := require.New(t)

		p := Processor{
			query: query{
				path: []domain.Step{
					{
						Name:  "objects",
						Index: -1,
					},
					{
						Name:  "object",
						Index: 1,
					},
				},
				searchType: domain.TagValue,
			},
			index: index{
				set: true,
			},
		}

		err := p.process([]byte(`<objects><object><tg></tg></object><object><tg1></tg1></object><object><tg2></tg2></object></objects>`))
		rq.NoError(err)
		rq.Len(p.printList, 4)
		rq.Equal("<object>", p.printList[0])
		rq.Equal("  <tg1>", p.printList[1])
		rq.Equal("  </tg1>", p.printList[2])
		rq.Equal("</object>", p.printList[3])
	})

	t.Run("tag value: single", func(t *testing.T) {
		t.Parallel()
		rq := require.New(t)

		p := Processor{
			query: query{
				path: []domain.Step{
					{
						Name:  "objects",
						Index: -1,
					},
					{
						Name:  "single",
						Index: 0,
					},
				},
				searchType: domain.TagValue,
			},
			index: index{
				set: true,
			},
		}

		err := p.process([]byte(`<objects><object><tg></tg></object><single /><data><tg2></tg2></data></objects>`))
		rq.NoError(err)
		rq.Len(p.printList, 1)
		rq.Equal("<single />", p.printList[0])
	})

	t.Run("tag name: first", func(t *testing.T) {
		t.Parallel()
		rq := require.New(t)

		p := Processor{
			query: query{
				path: []domain.Step{
					{
						Name:  "objects",
						Index: -1,
					},
					{
						Name:  "object",
						Index: 0,
					},
				},
				searchType: domain.TagList,
			},
			index: index{
				set: true,
			},
		}

		err := p.process([]byte(`<objects><object><tg></tg><data></data></object><object><tg1></tg1></object><object><tg2></tg2></object></objects>`))
		rq.NoError(err)
		rq.Len(p.printList, 2)
		rq.Equal("tg", p.printList[0])
		rq.Equal("data", p.printList[1])
	})

	t.Run("tag name: second", func(t *testing.T) {
		t.Parallel()
		rq := require.New(t)

		p := Processor{
			query: query{
				path: []domain.Step{
					{
						Name:  "objects",
						Index: -1,
					},
					{
						Name:  "object",
						Index: 1,
					},
				},
				searchType: domain.TagList,
			},
			index: index{
				set: true,
			},
		}

		err := p.process([]byte(`<objects><object><tg></tg><data></data></object><object><tg1></tg1></object><object><tg2></tg2></object></objects>`))
		rq.NoError(err)
		rq.Len(p.printList, 1)
		rq.Equal("tg1", p.printList[0])
	})

	t.Run("tag name: single", func(t *testing.T) {
		t.Parallel()
		rq := require.New(t)

		p := Processor{
			query: query{
				path: []domain.Step{
					{
						Name:  "objects",
						Index: -1,
					},
					{
						Name:  "single",
						Index: 0,
					},
				},
				searchType: domain.TagList,
			},
			index: index{
				set: true,
			},
		}

		err := p.process([]byte(`<objects><object><tg></tg></object><single /><data><tg2></tg2></data></objects>`))
		rq.NoError(err)
		rq.Len(p.printList, 0)
	})

	t.Run("attr list: first", func(t *testing.T) {
		t.Parallel()
		rq := require.New(t)

		p := Processor{
			query: query{
				path: []domain.Step{
					{
						Name:  "objects",
						Index: -1,
					},
					{
						Name:  "object",
						Index: 0,
					},
				},
				searchType: domain.AttrList,
			},
			index: index{
				set: true,
			},
		}

		err := p.process([]byte(`<objects><object attr1="value1" attr2="value2"><tg></tg><data></data></object><object attr3="value3" atrr4="value4"><tg1></tg1></object><object attr5="value5" atrr6="value6"><tg2></tg2></object></objects>`))
		rq.NoError(err)
		rq.Len(p.printList, 2)
		rq.Equal("attr1", p.printList[0])
		rq.Equal("attr2", p.printList[1])
	})

	t.Run("attr list: second", func(t *testing.T) {
		t.Parallel()
		rq := require.New(t)

		p := Processor{
			query: query{
				path: []domain.Step{
					{
						Name:  "objects",
						Index: -1,
					},
					{
						Name:  "object",
						Index: 1,
					},
				},
				searchType: domain.AttrList,
			},
			index: index{
				set: true,
			},
		}

		err := p.process([]byte(`<objects><object attr1="value1" attr2="value2"><tg></tg><data></data></object><object attr3="value3" attr4="value4"><tg1></tg1></object><object attr5="value5" atrr6="value6"><tg2></tg2></object></objects>`))
		rq.NoError(err)
		rq.Len(p.printList, 2)
		rq.Equal("attr3", p.printList[0])
		rq.Equal("attr4", p.printList[1])
	})

	t.Run("attr list: single", func(t *testing.T) {
		t.Parallel()
		rq := require.New(t)

		p := Processor{
			query: query{
				path: []domain.Step{
					{
						Name:  "objects",
						Index: -1,
					},
					{
						Name:  "single",
						Index: 0,
					},
				},
				searchType: domain.AttrList,
			},
			index: index{
				set: true,
			},
		}

		err := p.process([]byte(`<objects><object><tg></tg></object><single attr3="value3" attr4="value4" /><data><tg2></tg2></data></objects>`))
		rq.NoError(err)
		rq.Len(p.printList, 2)
		rq.Equal("attr3", p.printList[0])
		rq.Equal("attr4", p.printList[1])
	})

	t.Run("attr value: first", func(t *testing.T) {
		t.Parallel()
		rq := require.New(t)

		p := Processor{
			query: query{
				path: []domain.Step{
					{
						Name:  "objects",
						Index: -1,
					},
					{
						Name:  "object",
						Index: 0,
					},
				},
				searchType: domain.AttrValue,
				attribute:  "attr1",
			},
			index: index{
				set: true,
			},
		}

		err := p.process([]byte(`<objects><object attr1="value1" attr2="value2"><tg></tg><data></data></object><object attr1="value3" attr2="value4"><tg1></tg1></object><object attr1="value5" attr2="value6"><tg2></tg2></object></objects>`))
		rq.NoError(err)
		rq.Len(p.printList, 1)
		rq.Equal("value1", p.printList[0])
	})

	t.Run("attr value: second", func(t *testing.T) {
		t.Parallel()
		rq := require.New(t)

		p := Processor{
			query: query{
				path: []domain.Step{
					{
						Name:  "objects",
						Index: -1,
					},
					{
						Name:  "object",
						Index: 1,
					},
				},
				searchType: domain.AttrValue,
				attribute:  "attr1",
			},
			index: index{
				set: true,
			},
		}

		err := p.process([]byte(`<objects><object attr1="value1" attr2="value2"><tg></tg><data></data></object><object attr1="value3" attr2="value4"><tg1></tg1></object><object attr1="value5" attr2="value6"><tg2></tg2></object></objects>`))
		rq.NoError(err)
		rq.Len(p.printList, 1)
		rq.Equal("value3", p.printList[0])
	})

	t.Run("attr value: single", func(t *testing.T) {
		t.Parallel()
		rq := require.New(t)

		p := Processor{
			query: query{
				path: []domain.Step{
					{
						Name:  "objects",
						Index: -1,
					},
					{
						Name:  "single",
						Index: 0,
					},
				},
				searchType: domain.AttrValue,
				attribute:  "attr3",
			},
			index: index{
				set: true,
			},
		}

		err := p.process([]byte(`<objects><object><tg></tg></object><single attr3="value3" attr4="value4" /><data><tg2></tg2></data></objects>`))
		rq.NoError(err)
		rq.Len(p.printList, 1)
		rq.Equal("value3", p.printList[0])
	})

	t.Run("tag value: first tg", func(t *testing.T) {
		t.Parallel()
		rq := require.New(t)

		p := Processor{
			query: query{
				path: []domain.Step{
					{
						Name:  "objects",
						Index: -1,
					},
					{
						Name:  "object",
						Index: 0,
					},
					{
						Name:  "tg",
						Index: -1,
					},
				},
				searchType: domain.TagValue,
			},
			index: index{
				set: true,
			},
		}

		err := p.process([]byte(`<objects><object><tg></tg></object><object><tg1></tg1></object><object><tg2></tg2></object></objects>`))
		rq.NoError(err)
		rq.Len(p.printList, 2)
		rq.Equal("<tg>", p.printList[0])
		rq.Equal("</tg>", p.printList[1])
	})

	t.Run("tag value: second tg", func(t *testing.T) {
		t.Parallel()
		rq := require.New(t)

		p := Processor{
			query: query{
				path: []domain.Step{
					{
						Name:  "objects",
						Index: -1,
					},
					{
						Name:  "object",
						Index: 1,
					},
					{
						Name:  "tg1",
						Index: -1,
					},
				},
				searchType: domain.TagValue,
			},
			index: index{
				set: true,
			},
		}

		err := p.process([]byte(`<objects><object><tg></tg></object><object><tg1></tg1></object><object><tg2></tg2></object></objects>`))
		rq.NoError(err)
		rq.Len(p.printList, 2)
		rq.Equal("<tg1>", p.printList[0])
		rq.Equal("</tg1>", p.printList[1])
	})
}
