package domain

// SearchType represents the type of search.
type SearchType int

const (
	// TagList represents a tag list inside target search.
	//
	// Example:
	//
	// 	<tagname attr1="value1" attr2="value2">
	// 		<inside attr="value">42</inside>
	//		<inside2 attr="value">some data</inside2>
	// 	</tagname>
	//
	// SearchType = `TagList` target tag = `tagname`
	//
	// result:
	//		inside
	//		inside2
	//
	TagList SearchType = iota
	// TagValue represents tags value search.
	//
	// Example:
	// 	<tagname attr1="value1" attr2="value2">
	// 		<inside attr="value">42</inside>
	//		<inside2 attr="value">some data</inside2>
	// 	</tagname>
	//
	// SearchType = `TagValue` target tag = `tagname.inside`
	// result:
	//		42
	//
	// SearchType = `TagValue` target tag = `tagname.inside2`
	//
	// result:
	//		some data
	TagValue
	// AttrList represents search all attributes names inside a tag.
	//
	// Example:
	// 	<tagname attr1="value1" attr2="value2">
	// 		<inside attr="value">42</inside>
	//		<inside2 attr="value">some data</inside2>
	// 	</tagname>
	//
	// SearchType = `AttrList` target tag = `tagname`
	//
	// result:
	//		attr1
	//		attr2
	AttrList
	// AttrValue represents search attributes value for a tag.
	//
	// Example:
	// 	<tagname attr1="value1" attr2="value2">
	// 		<inside attr="value">42</inside>
	//		<inside2 attr="value">some data</inside2>
	// 	</tagname>
	//
	// SearchType = `AttrValue` target tag = `tagname` attribute = `attr2`
	//
	// result:
	//		value2
	AttrValue
)
