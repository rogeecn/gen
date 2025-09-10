# gen 使用手册

## 项目简介

`gen` 是一个基于 GORM 的 Go 代码生成器，专注于 PostgreSQL，支持数据库表结构到 Go 结构体的自动映射、类型安全的 SQL 构建、丰富的 JSON/数组/时间类型扩展，并可自定义生成查询和接口方法。

## 主要特性

- **数据库表到 Go 结构体自动生成**
- **类型安全的 SQL 构建器**
- **PostgreSQL 专用 JSON/数组/时间类型支持**
- **可扩展的自定义方法和接口生成**
- **一键生成/更新代码，支持单元测试生成**

## 快速开始

### 1. 安装依赖

```bash
go get go.ipao.vip/gen
```

### 2. 连接数据库

仅支持 PostgreSQL：

```go
import (
  "gorm.io/driver/postgres"
  "gorm.io/gorm"
  "go.ipao.vip/gen"
)

db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{})
if err != nil {
  panic(err)
}
```

### 3. 创建 Generator 并生成代码（同包同目录生成，默认推荐）

```go
g := gen.NewGenerator(gen.Config{
  // 将 model 与 query 代码生成到同一目录（同一包），便于直接调用
  OutPath:      "./database",
  ModelPkgPath: "./database", // 与 OutPath 一致即可（默认也是同包同目录）
  // 其它配置项见下文
})
g.UseDB(db)
g.GenerateAllTable() // 生成所有表
g.Execute()          // 输出代码到 OutPath
```

### 4. 进阶用法

- 只生成指定表：

  ```go
  g.GenerateModel("user")
  ```

- 自定义方法/接口：

  ```go
  g.ApplyInterface(func(MyMethodInterface){}, MyModel{})
  ```

- 生成单元测试：

  ```go
  g := gen.NewGenerator(gen.Config{WithUnitTest: true, ...})
  ```

## 类型扩展与表达式

详见 `types/README.md`，支持：

- `types.JSON`、`types.JSONSet`、`types.JSONQuery`：PostgreSQL JSON 字段查询与更新
- `types.Date`、`types.Time`：日期/时间类型
- `types.JSONType[T]`：泛型 JSON 映射
- 数组表达式：包含、交集、被包含等

## 代码生成配置项

- `OutPath`：生成代码输出目录（推荐与 `ModelPkgPath` 相同，实现“同包同目录生成”）
- `ModelPkgPath`：model 代码输出目录（同包时与 `OutPath` 一致）
- `WithUnitTest`：是否生成单元测试
- `FieldNullable`：数据库可空字段是否生成指针
- `FieldCoverable`：有默认值字段是否生成指针
- `FieldSignable`：无符号整型映射
- `FieldWithIndexTag`/`FieldWithTypeTag`：是否生成 gorm tag
- `Mode`：生成模式（如 WithDefaultQuery、WithoutContext 等）

## 常用命令

- 构建：`go build ./...`
- 测试：`go test ./...`
- 静态检查：`go vet ./...`
- Lint：`golangci-lint run`

## 目录结构说明

- `generator.go`：主生成器逻辑
- `config.go`：生成器配置
- `field/`、`types/`：SQL 字段类型与表达式扩展
- `internal/`：代码生成与解析核心
- `helper/`：辅助工具

## 贡献与开发

- 遵循 Go 标准格式和命名
- 测试用例与实现代码同目录
- PR/Commit 规范见 `AGENTS.md`

---

如需更详细的类型/表达式用法，请参考 `types/README.md`。如需补充具体示例或接口说明，请告知！

## 渐进式教学示例

以下示例带你从零到一，完成连接数据库、生成代码，并在业务代码中以类型安全的方式进行查询。

1. 连接 PostgreSQL 并创建 Generator（同包同目录）

```go
import (
    "gorm.io/driver/postgres"
    "gorm.io/gorm"
    gen "go.ipao.vip/gen"
)

dsn := "host=127.0.0.1 user=postgres password=xxx dbname=app port=5432 sslmode=disable"
db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{})
if err != nil { panic(err) }

g := gen.NewGenerator(gen.Config{
    OutPath:      "./database",   // 查询与模型生成到同一目录（同一包）
    ModelPkgPath: "./database",
})
g.UseDB(db)
```

2. 生成模型与查询代码

```go
// 生成单表（如 users）
g.GenerateModel("users")
// 或生成全部表
// g.GenerateAllTable()
g.Execute()
```

3. 在业务代码中使用生成的查询（同包同目录）

生成包为 `your/module/database`：

```go
import (
  dbpkg "your/module/database"
)

// 全局初始化一次（设置默认 Query 入口与包级快捷变量）
dbpkg.SetDefault(db)

// 方式A：使用包级变量（更短的调用链）
tbl, q := dbpkg.StudentQuery.QueryContext(ctx)
stu, err := q.Where(tbl.ID.Eq(1)).First()

// 方式B：使用全局 Q（字段名无 Query 后缀）
tbl2, q2 := dbpkg.Q.Student.QueryContext(ctx)
list, err := q2.Where(tbl2.Age.Gt(18)).Order(tbl2.ID.Desc()).Limit(20).Find()

// 统计与分页
count, err := q2.Where(tbl2.Famous.Is(true)).Count()
page, total, err := q2.Where(tbl2.Score.Gte(90)).FindByPage(0, 10)

// JOIN 示例（需要已生成 Student/Teacher 两个表）
student := dbpkg.Q.Student
teacher := dbpkg.Q.Teacher
rows, err := student.LeftJoin(teacher, student.Instructor.EqCol(teacher.ID)).
    Where(teacher.ID.Gt(0)).
    Select(student.Name, teacher.Name.As("teacher_name")).
    Find()
```

4. 模型实例快捷操作（同包调用，无需引入 query 包）

```go
stu.Name = "Alice"
// Update/Save/Create/Delete 直接调用同包下的查询入口
_, err := stu.Update(ctx)
err := stu.Save(ctx)
err := stu.Create(ctx)
_, err := stu.Delete(ctx)
```

4. JSON/时间等 PostgreSQL 类型的使用

`go.ipao.vip/gen/types` 提供 PostgreSQL 类型与 JSON 表达式支持（详见 `types/README.md`）：

- JSON 查询：`types.JSONQuery("attributes").HasKey("role")`、`Equals/Likes(..., "name")`
- JSON 更新：`types.JSONSet("attributes").Set("{age}", 20)`（使用 JSONB_SET）
- 时间与日期：`types.Time`、`types.Date`
- 其他类型：`Inet/CIDR/MACAddr`、`Range`、`TSVector/TSQuery` 等

5. 使用字段表达式（可选）

如需更强表达能力，可用 `go.ipao.vip/gen/field`：

```go
import (
  f "go.ipao.vip/gen/field"
  t "go.ipao.vip/gen/types"
)

// JSONB：键路径比较与模糊匹配
data := f.NewJSONB("users", "attributes")
_ = db.Where(data.KeyILike("meta.title", "%hello%"))

// 网络：包含/被包含
inet := f.NewInet("hosts", "addr")
_ = db.Where(inet.ContainedByEq(t.MustInet("10.0.0.0/8")))
```

6. 常见配置进阶（可按需使用）

- `WithUnitTest: true`：生成基础单元测试模板
- `FieldNullable/FieldCoverable`：控制 NULL/默认值字段是否使用指针
- `Mode`：选择是否生成默认全局查询入口等

如需更深入的字段表达式与类型示例，请查看 `types/README.md` 与 `field/` 相关说明。

## 配置文件

```yaml
ignores:
  - migrations
imports:
  - go.ipao.vip/gen
  - gen-test/dto
field_type:
  comprehensive_types_table:
    json_val: types.JSONType[dto.Test]
field_relate:
  students:
    Class:
      # belong_to, has_one, has_many, many_to_many
      relation: belongs_to
      table: classes
      references: id # 关联表的主键/被引用键（通常是 id）
      foreign_key: class_id # 当前表上的外键列（如 students.class_id）
      Json: class
    Teachers:
      # belong_to, has_one, has_many, many_to_many
      relation: many_to_many
      table: teachers
      pivot: class_teacher
      foreign_key: class_id # 当前表（students）用于关联的键（转为结构体字段名 ClassID）
      join_foreign_key: class_id # 中间表中指向当前表的列（class_teacher.class_id）
      references: id # 关联表（teachers）被引用的列（转为结构体字段名 ID）
      join_references: teacher_id # 中间表中指向关联表的列（class_teacher.teacher_id）
      Json: teachers
  teachers:
    Classes:
      relation: many_to_many
      table: classes
      pivot: class_teacher
  classes:
    Teachers:
      relation: many_to_many
      table: teachers
      pivot: class_teacher
```

### 关联关系字段说明（对齐 GORM）

- relation

  - 取值：`belongs_to`、`has_one`、`has_many`、`many_to_many`。
  - 对应 GORM 的四种关系：Belongs To、Has One、Has Many、Many2Many。

- table

  - 关联的目标表名（即另一侧模型对应的表）。

- pivot（仅 many_to_many）

  - 多对多中间表名称，对应 GORM 标签 `many2many:<pivot>`。

- foreign_key（按关系含义不同）

  - 对应 GORM 标签 `foreignKey:<Field>`。
  - belongs_to：当前表上的外键列（例如 `students.class_id`），会映射为当前模型上的字段（如 `ClassID`）。
  - has_one / has_many：外键在对端表上（例如 `credit_cards.user_id`）。配置时仍在当前表的配置块里填“外键列名”，生成时会正确落到 GORM 标签中。

- references（按关系含义不同）

  - 对应 GORM 标签 `references:<Field>`。
  - belongs_to：对端表被引用的列（一般是 `id`），映射为对端模型字段名（如 `ID`）。
  - has_one / has_many：被对端外键引用的当前模型列（一般是当前模型的 `ID` 字段）。

- join_foreign_key（仅 many_to_many）

  - 对应 GORM 标签 `joinForeignKey:<Field>`，指中间表里“指向当前模型”的列（如 `class_teacher.class_id`）。

- join_references（仅 many_to_many）
  - 对应 GORM 标签 `joinReferences:<Field>`，指中间表里“指向关联模型”的列（如 `class_teacher.teacher_id`）。

说明：生成器会结合数据库的 NamingStrategy 将列名（如 `class_id`、`teacher_id`）转换为结构体字段名（如 `ClassID`、`TeacherID`），并据此写入正确的 GORM 标签。

### 与 GORM 标签的对应关系

- belongs_to 示例（students → classes）

  - YAML：
    - `foreign_key: class_id`
    - `references: id`
  - 生成的模型字段（示意）：
    - `Class Class  gorm:"foreignKey:ClassID;references:ID"`

- has_many 示例（users → credit_cards）

  - YAML（在 `users` 下配置 `CreditCards` 关系）：
    - `relation: has_many`
    - `table: credit_cards`
    - `foreign_key: user_id` （对端表上的外键列）
    - `references: id` （当前模型被引用的列）
  - 生成的模型字段（示意）：
    - `CreditCards []CreditCard gorm:"foreignKey:UserID;references:ID"`

- many_to_many 示例（students ⇄ teachers，经由 class_teacher）
  - YAML（在 `students` 下配置 `Teachers` 关系）：
    - `relation: many_to_many`
    - `table: teachers`
    - `pivot: class_teacher`
    - `foreign_key: class_id`
    - `join_foreign_key: class_id`
    - `references: id`
    - `join_references: teacher_id`
  - 生成的模型字段（示意）：
    - `Teachers []Teacher gorm:"many2many:class_teacher;foreignKey:ClassID;references:ID;joinForeignKey:ClassID;joinReferences:TeacherID"`

提示：GORM 在 many2many 下允许省略部分键，生成器也支持“只给必要字段”。若不确定，建议显式全部写出，避免命名不一致导致推断失败。

## 快速生成

```go
  var dsn = "host=host user=postgres password=password dbname=test port=5432 sslmode=disable TimeZone=Asia/Shanghai"
  db, err := gorm.Open(postgres.New(postgres.Config{DSN: dsn}))
  if err != nil { log.Fatal(err) }

  // 默认采用“同包同目录生成”：OutPath 与 ModelPkgPath 一致，文件统一生成到 ./database
  gen.GenerateWithDefault(db, ".transform.yaml")

  // 初始化默认入口与包级变量（StudentQuery/TeacherQuery 等）
  yourpkg/database.SetDefault(db)

## 最小完整示例（目录结构 + 代码）

以下示例演示一个最小可运行流程：连接数据库 → 生成代码（同包同目录）→ 在业务代码中直接查询。

1) 生成代码（同包同目录）

main.go（仅用于生成）：

```go
package main

import (
    "log"
    "go.ipao.vip/gen"
    "gorm.io/driver/postgres"
    "gorm.io/gorm"
)

func main() {
    dsn := "host=127.0.0.1 user=postgres password=postgres dbname=demo port=5432 sslmode=disable TimeZone=Asia/Shanghai"
    db, err := gorm.Open(postgres.New(postgres.Config{DSN: dsn}))
    if err != nil { log.Fatal(err) }

    // 默认同包同目录生成到 ./database
    gen.GenerateWithDefault(db, ".transform.yaml")
}
```

执行：

```bash
go run main.go
```

生成后的目录（示意）：

```
database/
  classes.gen.go
  classes.query.gen.go
  class_teacher.gen.go
  class_teacher.query.gen.go
  students.gen.go
  students.query.gen.go
  teachers.gen.go
  teachers.query.gen.go
  query.gen.go            # 入口（提供 Q 与包级 StudentQuery/TeacherQuery 等）
```

2) 业务代码中使用

```go
package svc

import (
    dbpkg "your/module/database"
    "gorm.io/gorm"
)

func Init(db *gorm.DB) {
    dbpkg.SetDefault(db)
}

func FindStudent(id int32) (*dbpkg.Student, error) {
    tbl, q := dbpkg.StudentQuery.QueryContext(nil)
    return q.Where(tbl.ID.Eq(id)).First()
}

func UpdateStudentName(id int32, name string) error {
    // 通过 Q.Student 也可以调用（Query 字段无后缀）
    tbl, q := dbpkg.Q.Student.QueryContext(nil)
    _, err := q.Where(tbl.ID.Eq(id)).Update(tbl.Name, name)
    return err
}

func SaveStudent(m *dbpkg.Student) error { return m.Save(nil) }
```

说明：

- 包级变量（如 `StudentQuery`）指向 `Q.Student`，提供更短的调用链。
- 模型实例方法（Update/Save/Create/Delete）无需导入查询包，直接在同包内调用。
```
