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
