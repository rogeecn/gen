# GORM Data Types

## JSON

PostgreSQL supported

```go
import types "go.ipao.vip/gen/types"

type UserWithJSON struct {
	gorm.Model
	Name       string
	Attributes types.JSON
}

DB.Create(&UserWithJSON{
	Name:       "json-1",
	Attributes: types.JSON([]byte(`{"name": "jinzhu", "age": 18, "tags": ["tag1", "tag2"], "orgs": {"orga": "orga"}}`)),
}

// Check JSON has keys
types.JSONQuery("attributes").HasKey(value, keys...)

db.Find(&user, types.JSONQuery("attributes").HasKey("role"))
db.Find(&user, types.JSONQuery("attributes").HasKey("orgs", "orga"))
// MySQL
// SELECT * FROM `users` WHERE JSON_EXTRACT(`attributes`, '$.role') IS NOT NULL
// SELECT * FROM `users` WHERE JSON_EXTRACT(`attributes`, '$.orgs.orga') IS NOT NULL

// PostgreSQL
// SELECT * FROM "user" WHERE "attributes"::jsonb ? 'role'
// SELECT * FROM "user" WHERE "attributes"::jsonb -> 'orgs' ? 'orga'


// Check JSON extract value from keys equal to value
types.JSONQuery("attributes").Equals(value, keys...)

DB.First(&user, types.JSONQuery("attributes").Equals("jinzhu", "name"))
DB.First(&user, types.JSONQuery("attributes").Equals("orgb", "orgs", "orgb"))
// MySQL
// SELECT * FROM `user` WHERE JSON_EXTRACT(`attributes`, '$.name') = "jinzhu"
// SELECT * FROM `user` WHERE JSON_EXTRACT(`attributes`, '$.orgs.orgb') = "orgb"

// PostgreSQL
// SELECT * FROM "user" WHERE json_extract_path_text("attributes"::json,'name') = 'jinzhu'
// SELECT * FROM "user" WHERE json_extract_path_text("attributes"::json,'orgs','orgb') = 'orgb'
```

NOTE: This project targets PostgreSQL.

## Date

```go
import types "go.ipao.vip/gen/types"

type UserWithDate struct {
	gorm.Model
	Name string
	Date types.Date
}

user := UserWithDate{Name: "jinzhu", Date: types.Date(time.Now())}
DB.Create(&user)
// INSERT INTO `user_with_dates` (`name`,`date`) VALUES ("jinzhu","2020-07-17 00:00:00")

DB.First(&result, "name = ? AND date = ?", "jinzhu", types.Date(curTime))
// SELECT * FROM user_with_dates WHERE name = "jinzhu" AND date = "2020-07-17 00:00:00" ORDER BY `user_with_dates`.`id` LIMIT 1
```

## Time

PostgreSQL supported. Time with nanoseconds is supported.

```go
import types "go.ipao.vip/gen/types"

type UserWithTime struct {
    gorm.Model
    Name string
    Time types.Time
}

user := UserWithTime{Name: "jinzhu", Time: types.NewTime(1, 2, 3, 0)}
DB.Create(&user)
// INSERT INTO `user_with_times` (`name`,`time`) VALUES ("jinzhu","01:02:03")

DB.First(&result, "name = ? AND time = ?", "jinzhu", types.NewTime(1, 2, 3, 0))
// SELECT * FROM user_with_times WHERE name = "jinzhu" AND time = "01:02:03" ORDER BY `user_with_times`.`id` LIMIT 1
```

NOTE: This project targets PostgreSQL only.

## JSON_SET

PostgreSQL supported

```go
import (
	types "go.ipao.vip/gen/types"
	"gorm.io/gorm"
)

type UserWithJSON struct {
	gorm.Model
	Name       string
	Attributes types.JSON
}

DB.Create(&UserWithJSON{
	Name:       "json-1",
	Attributes: types.JSON([]byte(`{"name": "json-1", "age": 18, "tags": ["tag1", "tag2"], "orgs": {"orga": "orga"}}`)),
})

type User struct {
	Name string
	Age  int
}

friend := User{
	Name: "Bob",
	Age:  21,
}

// Set fields of JSON column
types.JSONSet("attributes").Set("age", 20).Set("tags[0]", "tag2").Set("orgs.orga", "orgb")

DB.Model(&UserWithJSON{}).Where("name = ?", "json-1").UpdateColumn("attributes", types.JSONSet("attributes").Set("age", 20).Set("tags[0]", "tag3").Set("orgs.orga", "orgb"))
DB.Model(&UserWithJSON{}).Where("name = ?", "json-1").UpdateColumn("attributes", types.JSONSet("attributes").Set("phones", []string{"10085", "10086"}))
DB.Model(&UserWithJSON{}).Where("name = ?", "json-1").UpdateColumn("attributes", types.JSONSet("attributes").Set("phones", gorm.Expr("CAST(? AS JSON)", `["10085", "10086"]`)))
DB.Model(&UserWithJSON{}).Where("name = ?", "json-1").UpdateColumn("attributes", types.JSONSet("attributes").Set("friend", friend))
// MySQL
// UPDATE `user_with_jsons` SET `attributes` = JSON_SET(`attributes`, '$.tags[0]', 'tag3', '$.orgs.orga', 'orgb', '$.age', 20) WHERE name = 'json-1'
// UPDATE `user_with_jsons` SET `attributes` = JSON_SET(`attributes`, '$.phones', CAST('["10085", "10086"]' AS JSON)) WHERE name = 'json-1'
// UPDATE `user_with_jsons` SET `attributes` = JSON_SET(`attributes`, '$.phones', CAST('["10085", "10086"]' AS JSON)) WHERE name = 'json-1'
// UPDATE `user_with_jsons` SET `attributes` = JSON_SET(`attributes`, '$.friend', CAST('{"Name": "Bob", "Age": 21}' AS JSON)) WHERE name = 'json-1'
```

NOTE: MariaDB does not support CAST(? AS JSON).

NOTE: Path in PostgreSQL is different.

```go
// Set fields of JSON column
types.JSONSet("attributes").Set("{age}", 20).Set("{tags, 0}", "tag2").Set("{orgs, orga}", "orgb")

DB.Model(&UserWithJSON{}).Where("name = ?", "json-1").UpdateColumn("attributes", types.JSONSet("attributes").Set("{age}", 20).Set("{tags, 0}", "tag2").Set("{orgs, orga}", "orgb"))
DB.Model(&UserWithJSON{}).Where("name = ?", "json-1").UpdateColumn("attributes", types.JSONSet("attributes").Set("{phones}", []string{"10085", "10086"}))
DB.Model(&UserWithJSON{}).Where("name = ?", "json-1").UpdateColumn("attributes", types.JSONSet("attributes").Set("{phones}", gorm.Expr("?::jsonb", `["10085", "10086"]`)))
DB.Model(&UserWithJSON{}).Where("name = ?", "json-1").UpdateColumn("attributes", types.JSONSet("attributes").Set("{friend}", friend))
// PostgreSQL
// UPDATE "user_with_jsons" SET "attributes" = JSONB_SET(JSONB_SET(JSONB_SET("attributes", '{age}', '20'), '{tags, 0}', '"tag2"'), '{orgs, orga}', '"orgb"') WHERE name = 'json-1'
// UPDATE "user_with_jsons" SET "attributes" = JSONB_SET("attributes", '{phones}', '["10085","10086"]') WHERE name = 'json-1'
// UPDATE "user_with_jsons" SET "attributes" = JSONB_SET("attributes", '{phones}', '["10085","10086"]'::jsonb) WHERE name = 'json-1'
// UPDATE "user_with_jsons" SET "attributes" = JSONB_SET("attributes", '{friend}', '{"Name": "Bob", "Age": 21}') WHERE name = 'json-1'
```

## JSONType[T]

PostgreSQL supported

```go
import types "go.ipao.vip/gen/types"

type Attribute struct {
	Sex   int
	Age   int
	Orgs  map[string]string
	Tags  []string
	Admin bool
	Role  string
}

type UserWithJSON struct {
	gorm.Model
	Name       string
	Attributes types.JSONType[Attribute]
}

var user = UserWithJSON{
	Name: "hello",
	Attributes: types.NewJSONType(Attribute{
        Age:  18,
        Sex:  1,
        Orgs: map[string]string{"orga": "orga"},
        Tags: []string{"tag1", "tag2", "tag3"},
    }),
}

// Create
DB.Create(&user)

// First
var result UserWithJSON
DB.First(&result, user.ID)

// Update
jsonMap = UserWithJSON{
	Attributes: types.NewJSONType(Attribute{
        Age:  18,
        Sex:  1,
        Orgs: map[string]string{"orga": "orga"},
        Tags: []string{"tag1", "tag2", "tag3"},
    }),
}

DB.Model(&user).Updates(jsonMap)
```

NOTE: it's not support json query

## JSONSlice[T]

PostgreSQL supported

```go
import types "go.ipao.vip/gen/types"

type Tag struct {
	Name  string
	Score float64
}

type UserWithJSON struct {
	gorm.Model
	Name       string
	Tags       types.JSONSlice[Tag]
}

var tags = []Tag{{Name: "tag1", Score: 0.1}, {Name: "tag2", Score: 0.2}}
var user = UserWithJSON{
	Name: "hello",
	Tags: types.NewJSONSlice(tags),
}

// Create
DB.Create(&user)

// First
var result UserWithJSON
DB.First(&result, user.ID)

// Update
var tags2 = []Tag{{Name: "tag3", Score: 10.1}, {Name: "tag4", Score: 10.2}}
jsonMap = UserWithJSON{
	Tags: types.NewJSONSlice(tags2),
}

DB.Model(&user).Updates(jsonMap)
```

NOTE: it's not support json query and `db.Pluck` method

## JSONArray

mysql supported

```go
import "gorm.io/datatypes"

type Param struct {
    ID          int
    Letters     string
    Config      types.JSON
}

//Create
DB.Create(&Param{
    Letters: "JSONArray-1",
    Config:      types.JSON("[\"a\", \"b\"]"),
})

DB.Create(&Param{
    Letters: "JSONArray-2",
    Config:      types.JSON("[\"a\", \"c\"]"),
})

//Query
var retMultiple []Param
DB.Where(types.JSONArrayQuery("config").Contains("c")).Find(&retMultiple)
}
```

## PostgreSQL-Specific Types (go.ipao.vip/gen/types)

The `go.ipao.vip/gen/types` package provides first-class Go types for many PostgreSQL data types. These types implement `driver.Valuer` and `sql.Scanner`, and expose small constructor helpers for convenient usage.

Import in your code:

```go
import (
    genTypes "go.ipao.vip/gen/types"
)
```

### Network

- Types: `genTypes.Inet`, `genTypes.CIDR`, `genTypes.MACAddr`
- Constructors: `NewInet(string) (Inet, error)`, `MustInet(string) Inet`, `NewCIDR`, `MustCIDR`, `NewMACAddr`, `MustMACAddr`

Example (model + create + query):

```go
type Host struct {
    ID   uint
    Addr genTypes.Inet
    Net  genTypes.CIDR
    Mac  genTypes.MACAddr
}

h := Host{
    Addr: genTypes.MustInet("192.168.1.10"),
    Net:  genTypes.MustCIDR("192.168.1.0/24"),
    Mac:  genTypes.MustMACAddr("08:00:2b:01:02:03"),
}
_ = DB.Create(&h).Error
```

### Bit String

- Type: `genTypes.BitString`
- Constructor: `NewBitString(string) BitString`

```go
type Bits struct { ID uint; B genTypes.BitString }
_ = DB.Create(&Bits{B: genTypes.NewBitString("10101010")}).Error
```

### Geometry

- Types: `genTypes.Point`, `Box`, `Path`, `Polygon`, `Circle`
- Constructors: `NewPoint`, `NewBox`, `NewPath`, `NewPolygon`, `NewCircle`

```go
type Geo struct {
    ID   uint
    Pt   genTypes.Point
    Area genTypes.Polygon
}

g := Geo{Pt: genTypes.NewPoint(1, 2), Area: genTypes.NewPolygon([]genTypes.Point{{1,1},{2,2},{3,1}})}
_ = DB.Create(&g).Error
```

### Ranges

- Types: `genTypes.Int4Range`, `Int8Range`, `NumRange`, `TsRange`, `TstzRange`, `DateRange`
- Constructors: `NewInt4Range`, `NewInt8Range`, `NewNumRange`, `NewTsRange`, `NewTstzRange`, `NewDateRange`

```go
type Period struct { ID uint; R genTypes.Int4Range }
p := Period{R: genTypes.NewInt4Range(1, 10, true, false)} // [1,10)
_ = DB.Create(&p).Error
```

With generated field helpers (if you use gorm/gen fields):

```go
import (
    genField "go.ipao.vip/gen/field"
)

// WHERE r && [5,7)
_ = DB.Where(genField.NewInt4Range("periods", "r").Overlaps(genTypes.NewInt4Range(5, 7, true, false))).Find(&[]Period{}).Error
```

### Full Text

- Types: `genTypes.TSVector`, `TSQuery`
- Constructors: `NewTSVector`, `NewTSQuery`

```go
type Doc struct { ID uint; Vec genTypes.TSVector }
_ = DB.Create(&Doc{Vec: genTypes.NewTSVector("'quick':1 'brown':2 'fox':3")}).Error
```

With field helper for matching (generated by gorm/gen fields):

```go
import genField "go.ipao.vip/gen/field"
q := genTypes.NewTSQuery("fox & jump")
_ = DB.Where(genField.NewTSVector("docs", "vec").Matches(q)).Find(&[]Doc{}).Error
```

### XML & Money

- Types: `genTypes.XML`, `genTypes.Money`
- Constructors: `NewXML`, `NewMoney`

```go
type Pay struct { ID uint; Price genTypes.Money; Meta genTypes.XML }
_ = DB.Create(&Pay{Price: genTypes.NewMoney("$123.45"), Meta: genTypes.NewXML("<root><item>ok</item></root>")}).Error
```

### Field Helpers (PostgreSQL)

Below are handy `go.ipao.vip/gen/field` helpers for common PostgreSQL types. All methods are strongly typed unless noted.

```go
import (
    genField "go.ipao.vip/gen/field"
    genTypes "go.ipao.vip/gen/types"
)

// JSON/JSONB: key-path comparisons and pattern matching
// Path uses dot notation: "a.b.c" and performs json_extract_path_text internally.
_ = DB.Where(genField.NewJSONB("t", "data").KeyEq("a.b.c", 10)).Find(&[]any{}).Error
_ = DB.Where(genField.NewJSONB("t", "data").KeyILike("meta.title", "%hello%"))).Error
_ = DB.Where(genField.NewJSON("t", "data").KeyRegexp("tags.0", "^(foo|bar)$")).Error

// Money: typed comparisons and ranges
price := genField.NewMoney("payments", "price")
_ = DB.Where(price.Gt(genTypes.NewMoney("$10.00")).And(price.Lte(genTypes.NewMoney("$99.99")))).Find(&[]any{}).Error
_ = DB.Where(price.Between(genTypes.NewMoney("$10"), genTypes.NewMoney("$20"))).Error

// String/Bytes: case-insensitive matching (ILIKE)
name := genField.NewString("users", "name")
_ = DB.Where(name.ILike("%smith%"))
payload := genField.NewBytes("docs", "payload")
_ = DB.Where(payload.ILike("%xml%"))

// INET/CIDR: network containment operators
inet := genField.NewInet("hosts", "addr")
_ = DB.Where(inet.Contains(genTypes.NewInet("10.0.0.0/8")))        // addr >> 10.0.0.0/8
_ = DB.Where(inet.ContainedByEq(genTypes.NewInet("10.0.0.0/8")))    // addr <<= 10.0.0.0/8

// MACADDR/BitString: In/NotIn
mac := genField.NewMACAddr("nics", "mac")
_ = DB.Where(mac.In(genTypes.NewMACAddr("aa:bb:cc:dd:ee:ff")))

// Ranges: overlaps/containment/adjacency
nr := genField.NewNumRange("tbl", "r")
_ = DB.Where(nr.Overlaps(genTypes.NewNumRange(1.0, 2.0, true, false)))
_ = DB.Where(nr.Contains(genTypes.NewNumRange(1.2, 1.8, true, true)))
_ = DB.Where(nr.Adjacent(genTypes.NewNumRange(2.0, 3.0, false, true)))

// Geometry: overlaps, contains, distance, within
poly := genField.NewPolygon("g", "poly")
pt := genField.NewPoint("g", "pt")
_ = DB.Where(poly.ContainsPoint(genTypes.NewPoint(1, 2)))
_ = DB.Order(pt.DistanceTo(genTypes.NewPoint(10, 10)).Asc())

// Array: contains/contained-by/overlaps (caller provides driver.Valuer or raw literal)
arr := genField.NewArray("t", "vals")
_ = DB.Where(arr.Contains(`{1,2}`))
_ = DB.Where(arr.Overlaps(`{2,3}`))
```

### Bytea Hex Helper

- Type: `genTypes.HexBytes` for handling `BYTEA` hex strings like "\\xDEADBEEF"

```go
hb, _ := genTypes.NewHexBytesFromHex("DEADBEEF")
type B struct { ID uint; Data genTypes.HexBytes }
_ = DB.Create(&B{Data: hb}).Error
```

### JSONType[T] / JSONSlice[T]

The generic JSON helpers are available under `go.ipao.vip/gen/types` too, mirroring the creation style:

```go
type Attr struct { Name string; Age int }
type U struct { ID uint; A types.JSONType[Attr]; L types.JSONSlice[string] }
u := U{A: types.NewJSONType(Attr{Name:"n", Age:1}), L: types.NewJSONSlice([]string{"a","b"})}
_ = DB.Create(&u).Error
```

## UUID

MySQL, PostgreSQL, SQLServer and SQLite are supported.

```go
import "gorm.io/datatypes"

type UserWithUUID struct {
    gorm.Model
    Name string
    UserUUID types.UUID
}

// Generate a new random UUID (version 4).
userUUID := types.NewUUIDv4()

user := UserWithUUID{Name: "jinzhu", UserUUID: userUUID}
DB.Create(&user)
// INSERT INTO `user_with_uuids` (`name`,`user_uuid`) VALUES ("jinzhu","ca95a578-816c-4812-babd-a7602b042460")

var result UserWithUUID
DB.First(&result, "name = ? AND user_uuid = ?", "jinzhu", userUUID)
// SELECT * FROM user_with_uuids WHERE name = "jinzhu" AND user_uuid = "ca95a578-816c-4812-babd-a7602b042460" ORDER BY `user_with_uuids`.`id` LIMIT 1

// Use the datatype's Equals() to compare the UUIDs.
if userCreate.UserUUID.Equals(userFound.UserUUID) {
	fmt.Println("User UUIDs match as expected.")
} else {
	fmt.Println("User UUIDs do not match. Something is wrong.")
}

// Use the datatype's String() function to get the UUID as a string type.
fmt.Printf("User UUID is %s", userFound.UserUUID.String())

// Check the UUID value with datatype's IsNil() and IsEmpty() functions.
if userFound.UserUUID.IsNil() {
	fmt.Println("User UUID is a nil UUID (i.e. all bits are zero)")
}
if userFound.UserUUID.IsEmpty() {
	fmt.Println(
		"User UUID is empty (i.e. either a nil UUID or a zero length string)",
	)
}
```
