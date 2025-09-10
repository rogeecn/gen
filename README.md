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

### 3. 创建 Generator 并生成代码

```go
g := gen.NewGenerator(gen.Config{
  OutPath:      "./dao/query",
  ModelPkgPath: "./dao/model",
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

- `OutPath`：生成代码输出目录
- `ModelPkgPath`：生成 model 代码的包名
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

1. 连接 PostgreSQL 并创建 Generator

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
    OutPath:      "./dao/query", // 查询代码输出目录
    ModelPkgPath: "./dao/model", // model 结构体输出目录
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

3. 在业务代码中使用生成的查询

假设已生成包 `dao/query` 与 `dao/model`：

```go
import (
  q "your/module/dao/query"
)

qry := q.Use(db)      // 构造 Query 入口
u := qry.Users        // 获取 users 的查询句柄（字段名与生成的 Model 名匹配）

// 3.1 基础查询
user, err := u.Where(u.ID.Eq(1)).First()
list, err := u.Where(u.Age.Gt(18)).Order(u.ID.Desc()).Limit(20).Find()

// 3.2 统计与分页
count, err := u.Where(u.Famous.Is(true)).Count()
page, total, err := u.Where(u.Score.Gte(90)).FindByPage(0, 10)

// 3.3 JOIN（示例：student LEFT JOIN teacher）
// 需要在生成器中包含 student/teacher 两个表
student := qry.Student
teacher := qry.Teacher
rows, err := student.LeftJoin(teacher, student.Instructor.EqCol(teacher.ID)).
    Where(teacher.ID.Gt(0)).
    Select(student.Name, teacher.Name.As("teacher_name")).
    Find()
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
# `User` belongs to `Company`, `User.CompanyID` is the foreign key
# User has one CreditCard, CreditCard.UserID is the foreign key
# User has many CreditCards, CreditCard.UserID is the foreign key
field_relate:
  students:
    Class:
      # belong_to, has_one, has_many, many_to_many
      relation: belongs_to
      table: classes
      references: id # 关联表ID
      foreign_key: class_id # 当前表ID
      Json: class
    Teachers:
      # belong_to, has_one, has_many, many_to_many
      relation: many_to_many
      table: teachers
      pivot: class_teacher
      foreign_key: class_id # 当前表ID
      join_foreign_key: class_id # 关联中间表ID
      references: id # 关联表ID
      join_references: teacher_id # 关联跳转表ID
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

## 快速生成

```go
  var dsn = "host=host user=postgres password=password dbname=test port=5432 sslmode=disable TimeZone=Asia/Shanghai"
	db, err := gorm.Open(postgres.New(postgres.Config{DSN: dsn}))
	if err != nil {
		log.Fatal(err)
	}

	gen.GenerateWithDefault(db, "")
```
