# xq

xq is a cli xml processor

inspired by [jq](https://github.com/stedolan/jq)

## usage

example structure

    <objects>
        <object>
            <title lang="EN">Name</title>
            <description>Name</description>
            <key attr="first">1</key>
        </object>
        <object>
            <title lang="RU">Имя</title>
            <description>Описание</description>
            <key attr="second">2</key>
        </object>
    </objects>

### tags

first level

    ~$ xq tags .

    objects

second level

    ~$ xq tags .objects

    object

third level

    ~$ xq tags .objects.object

    title
    description
    key

### get a tag value in a list

    ~$ xq .objects.object

    <object>
        <title lang="EN">
            Name
        </title>
        <description>
            Name
        </description>
        <key attr="first">
            1
        </key>
    </object>
    <object>
        <title lang="RU">
            Имя
        </title>
        <description>
            Описание
        </description>
        <key attr="second">
            2
        </key>
    </object>


another example

    ~$ xq .objects.object.title

    <title lang="EN">
        Name
    </title>
    <title lang="RU">
        Имя
    </title>

### get a value of a concrete tag 


    ~$ xq .objects.object[0]

    <object>
        <title lang="EN">
            Name
        </title>
        <description>
            Name
        </description>
        <key attr="first">
            1
        </key>
    </object>

### get an attribute value

for tags list

    ~$ xq .objects.object.title#lang

    EN
    RU

for concrete tag

    ~$ xq .objects.object[0].title#lang

    EN

## API Status

- [x] Add indentation for output
- [x] Colorize tags
- [x] Colorize attributes
- [x] Get tags list
- [x] Get attributes list
- [x] Get tag's data
- [x] Get attributes data