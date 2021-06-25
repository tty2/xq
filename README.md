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

if tag is a container for other tags

    ~$ xq .objects.object

    <title lang="EN">Name</title>
    <description>Name</description>
    <key attr="first">1</key>

    <title lang="RU">Имя</title>
    <description>Описание</description>
    <key attr="second">2</key>


if tag is a container for scalar data

    ~$ xq .objects.object.title

    Name
    Имя

### get a value of a concrete tag 


if tag is a container for other tags

    ~$ xq .objects.object[0]

    <title lang="EN">Name</title>
    <description>Name</description>
    <key attr="first">1</key>

if tag is a container for scalar data

    ~$ xq .objects.object[0].title

    Name

### get an attribute value

for tags list

    ~$ xq .objects.object.title:lang

    EN
    RU

for concrete tag

    ~$ xq .objects.object[0].title:lang

    EN

## API Status

- [x] Add indentation for output
- [x] Colorize tags
- [ ] Colorize attributes
- [ ] Get tags list
- [ ] Get attributes list
- [ ] Get tag's data
- [ ] Get attributes data