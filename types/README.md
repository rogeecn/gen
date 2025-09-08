# PostgreSQL 类型与表达式使用指南（types）

本文仅针对 PostgreSQL，介绍 `go.ipao.vip/gen/types` 中最常用的类型与表达式，并按“从易到难”的顺序给出示例，帮助你循序渐进上手。

适用于 GORM v2；示例中的 `DB` 为已经初始化好的 `*gorm.DB`。

```go
import (
    types  "go.ipao.vip/gen/types"
    // 可选：与字段表达式配合使用
    // field "go.ipao.vip/gen/field"
)
```

## 1. JSON 基础（查询）

定义 JSON/JSONB 字段并进行查询。

```go
type UserWithJSON struct {
    gorm.Model
    Name       string
    Attributes types.JSON // 建议在 PostgreSQL 中使用 JSONB 列类型
}

// 写入
DB.Create(&UserWithJSON{
    Name:       "json-1",
    Attributes: types.JSON([]byte(`{"name":"jinzhu","age":18,"tags":["t1","t2"],"orgs":{"orga":"orga"}}`)),
})

// 1) 判断键是否存在
DB.Where(types.JSONQuery("attributes").HasKey("role")).Find(&[]UserWithJSON{})
DB.Where(types.JSONQuery("attributes").HasKey("orgs", "orga")).Find(&[]UserWithJSON{})

// 2) 提取键值并比较（=）
DB.Where(types.JSONQuery("attributes").Equals("jinzhu", "name")).First(&UserWithJSON{})
DB.Where(types.JSONQuery("attributes").Equals("orgb", "orgs", "orgb")).First(&UserWithJSON{})

// 3) 提取键值并模糊匹配（LIKE）
DB.Where(types.JSONQuery("attributes").Likes("%jin%", "name")).Find(&[]UserWithJSON{})
```

要点：

- `HasKey(keys...)` 用于键是否存在；嵌套路径按多参数依次传入。
- `Equals/Likes(value, keys...)` 会对提取到的文本进行比较。

## 2. JSON 更新（JSONB_SET）

在 PostgreSQL 中更新 JSONB 子路径需使用花括号路径：`{age}`、`{orgs,orga}`、`{tags,0}`。

```go
// 增量更新多个子字段
DB.Model(&UserWithJSON{}).
   Where("name = ?", "json-1").
   UpdateColumn("attributes",
       types.JSONSet("attributes").
           Set("{age}", 20).
           Set("{tags,0}", "t3").
           Set("{orgs,orga}", "orgb"),
   )

// 写入数组/对象（会自动转 jsonb；也可手动传入 gorm.Expr("?::jsonb", raw)）
DB.Model(&UserWithJSON{}).
   Where("name = ?", "json-1").
   UpdateColumn("attributes", types.JSONSet("attributes").Set("{phones}", []string{"10085", "10086"}))
```

要点：

- 路径必须是 PostgreSQL 形式：`{key,sub,0}`。
- 传入的 Go 值会被序列化为 JSON 再写入。

## 3. 泛型 JSON（JSONType[T] / JSONSlice[T]）

以强类型的方式读写 JSON 数据。注意：这些泛型类型用于读写，不提供 JSON 路径查询能力。

```go
type Attribute struct {
    Age  int
    Orgs map[string]string
    Tags []string
}

type UserWithTypedJSON struct {
    gorm.Model
    Name       string
    Attributes types.JSONType[Attribute]
    Tags       types.JSONSlice[string]
}

u := UserWithTypedJSON{
    Name:       "u1",
    Attributes: types.NewJSONType(Attribute{Age: 18, Orgs: map[string]string{"orga": "orga"}, Tags: []string{"t1","t2"}}),
    Tags:       types.NewJSONSlice([]string{"t1", "t2"}),
}
DB.Create(&u)

// 更新
DB.Model(&u).Updates(UserWithTypedJSON{Tags: types.NewJSONSlice([]string{"t3", "t4"})})
```

## 4. 日期与时间（Date / Time）

```go
type UserWithDate struct {
    gorm.Model
    Name string
    Date types.Date
}
DB.Create(&UserWithDate{Name: "d1", Date: types.Date(time.Now())})
DB.Where("name = ? AND date = ?", "d1", types.Date(time.Now())).First(&UserWithDate{})

type UserWithTime struct {
    gorm.Model
    Name string
    Time types.Time // 支持纳秒
}
DB.Create(&UserWithTime{Name: "t1", Time: types.NewTime(1, 2, 3, 0)})
DB.Where("name = ? AND time = ?", "t1", types.NewTime(1, 2, 3, 0)).First(&UserWithTime{})
```

## 5. 常用 PostgreSQL 专用类型（速览）

以下类型均实现了 `sql.Scanner` / `driver.Valuer`，可直接作为模型字段使用。

网络（INET/CIDR/MACADDR）

```go
type Host struct {
    ID   uint
    Addr types.Inet
    Net  types.CIDR
    Mac  types.MACAddr
}
_ = DB.Create(&Host{
    Addr: types.MustInet("192.168.1.10"),
    Net:  types.MustCIDR("192.168.1.0/24"),
    Mac:  types.MustMACAddr("08:00:2b:01:02:03"),
}).Error
```

位串（BitString）

```go
type Bits struct { ID uint; B types.BitString }
_ = DB.Create(&Bits{B: types.NewBitString("10101010")}).Error
```

范围（Range）

```go
// Int4Range / Int8Range / NumRange / TsRange / TstzRange / DateRange
nr := types.NewNumRange(new(big.Rat).SetFloat64(1.0), new(big.Rat).SetFloat64(2.0), true, false)
type R struct { ID uint; R types.NumRange }
_ = DB.Create(&R{R: nr}).Error
```

几何（Point / Polygon / Box / Circle / Path）

```go
type Geo struct { ID uint; Pt types.Point; Poly types.Polygon }
_ = DB.Create(&Geo{Pt: types.NewPoint(1, 2), Poly: types.NewPolygon([]types.Point{{1,2},{3,4},{5,6}})}).Error
```

全文检索（TSQuery / TSVector）

```go
type Doc struct { ID uint; Vec types.TSVector }
_ = DB.Create(&Doc{Vec: types.NewTSVector("'quick':1 'brown':2 'fox':3")}).Error
// 结合 field 帮助器可进行匹配（可选）：
// q := types.NewTSQuery("fox & jump")
// DB.Where(field.NewTSVector("docs", "vec").Matches(q)).Find(&[]Doc{})
```

Money / XML / URL / BYTEA Hex

```go
type Pay struct { ID uint; Price types.Money; Meta types.XML; Site types.URL; Raw types.HexBytes }
_ = DB.Create(&Pay{Price: types.NewMoney("$123.45"), Meta: types.NewXML("<root/>")}).Error
```

UUID / BinUUID

```go
type U struct { ID uint; U1 types.UUID; U2 types.BinUUID }
u := U{U1: types.NewUUIDv4()}
_ = DB.Create(&u).Error
```

——

如需更多示例或想补充某个 PostgreSQL 类型的用法，请告知我们想覆盖的场景，我们会完善文档与示例代码。

## 6. Field 表达式进阶（PostgreSQL）

当你需要更强的可读性与可组合的表达式时，可使用 `go.ipao.vip/gen/field` 提供的类型安全字段助手。

```go
import (
    field "go.ipao.vip/gen/field"
    types "go.ipao.vip/gen/types"
)
```

- JSONB：键路径比较/匹配与键存在性

```go
data := field.NewJSONB("users", "attributes")
// 键存在（单键、任意键、全部键）
DB.Where(data.HasKey("role")).Find(&[]any{})
DB.Where(data.HasAnyKeys("role", "email_verified")).Find(&[]any{})
DB.Where(data.HasAllKeys("name", "age")).Find(&[]any{})

// 点分路径取值比较与匹配
DB.Where(data.KeyEq("profile.age", 18)).Find(&[]any{})
DB.Where(data.KeyILike("meta.title", "%hello%"))
DB.Where(data.KeyRegexp("tags.0", "^(foo|bar)$"))

// 直接提取文本作为子表达式使用
title := data.ExtractText("meta", "title")
DB.Where(title.ILike("%news%"))
```

- Range：重叠、包含、相邻

```go
nr := field.NewNumRange("events", "amount_range")
q := types.NewNumRange(new(big.Rat).SetFloat64(1.0), new(big.Rat).SetFloat64(2.0), true, false)
DB.Where(nr.Overlaps(q)).Or(nr.Contains(q)).Or(nr.Adjacent(q)).Find(&[]any{})
```

- INET/CIDR：网络包含关系

```go
inet := field.NewInet("hosts", "addr")
DB.Where(inet.Contains(types.MustInet("10.0.0.0/8")))     // >>
DB.Where(inet.ContainedByEq(types.MustInet("192.168.0.0/16"))) // <<=
```

- 几何：包含/距离/相交

```go
pt := field.NewPoint("geo", "pt")
poly := field.NewPolygon("geo", "poly")
DB.Where(poly.ContainsPoint(types.NewPoint(1, 2)))
DB.Order(pt.DistanceTo(types.NewPoint(10, 10)).Asc())
DB.Where(poly.Overlaps(types.NewPolygon([]types.Point{{0,0},{3,0},{3,3},{0,3}})))
```

- 全文检索：向量匹配

```go
vec := field.NewTSVector("docs", "vec")
q := types.NewTSQuery("fox & jump")
DB.Where(vec.Matches(q)).Find(&[]any{})
```

- 数组：包含/被包含/重叠（调用方提供数组值）

```go
arr := field.NewArray("t", "vals")
// 传入驱动实现或原始字面量（例如 `pq.Array([]int{1,2})` 或 `"{1,2}"`）
DB.Where(arr.Contains("{1,2}")).Or(arr.Overlaps("{2,3}"))
```
