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

    ~$ xq tags

    objects

second level

    ~$ xq tags objects

    object

third level

    ~$ xq tags objects.object

    title
    description
    key

### get a tag value in a list

if tag is a container for other tags

    ~$ xq objects.object

    <title lang="EN">Name</title>
    <description>Name</description>
    <key attr="first">1</key>

    <title lang="RU">Имя</title>
    <description>Описание</description>
    <key attr="second">2</key>

if tag is a container for other tags, for the first object only

    ~$ xq objects.object[0]

    <title lang="EN">Name</title>
    <description>Name</description>
    <key attr="first">1</key>

if tag is a container for scalar data

    ~$ xq objects.object.title

    Name
    Имя

### get a value of a concrete tag 

    ~$ xq objects.object[0].title

    Name

### get an attribute value

for tags list

    ~$ xq objects.object.title:lang

    EN
    RU

for concrete tag

    ~$ xq objects.object[0].title:lang

    EN