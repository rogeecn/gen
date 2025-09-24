# Gen - GORM 代码生成器

## 项目概述

**Gen** 是一个基于 GORM 的 Go 代码生成器，专注于 PostgreSQL，支持数据库表结构到 Go 结构体的自动映射、类型安全的 SQL 构建、丰富的 JSON/数组/时间类型扩展，并可自定义生成查询和接口方法。

**核心特性：**

- 数据库表到 Go 结构体自动生成
- 类型安全的 SQL 构建器
- PostgreSQL 专用 JSON/数组/时间类型支持
- 可扩展的自定义方法和接口生成
- 一键生成/更新代码，支持单元测试生成

## 技术架构

### 核心模块结构

```
gen/
├── generator.go          # 主生成器逻辑
├── config.go             # 生成器配置
├── do.go                 # 数据操作接口
├── interface.go          # 接口定义
├── condition.go          # 查询条件
├── import.go             # 导入管理
├── field/                # SQL 字段类型与表达式扩展
│   ├── array.go         # 数组类型支持
│   ├── association.go   # 关联关系处理
│   ├── export.go        # 字段导出功能
│   ├── jsonb.go         # JSONB 类型支持
│   ├── time.go          # 时间类型支持
│   └── ...
├── types/                # PostgreSQL 特殊类型支持
│   ├── json.go          # JSON 类型
│   ├── json_expr.go     # JSON 表达式
│   ├── array.go         # 数组类型
│   ├── datetime_*.go    # 日期时间类型
│   ├── network_*.go     # 网络类型
│   ├── range_*.go       # 范围类型
│   └── ...
├── internal/             # 内部实现
│   ├── generate/        # 代码生成核心
│   ├── model/           # 模型定义
│   ├── parser/          # 解析器
│   ├── template/        # 模板
│   └── utils/           # 工具
└── helper/              # 辅助工具
```

### 关键组件

**Generator (`generator.go`):**

- 主要代码生成逻辑
- 支持"同包同目录"生成模式
- 处理数据库连接、表解析、代码生成

**Config (`config.go`):**

- 生成器配置管理
- 支持多种生成模式（WithDefaultQuery, WithoutContext, WithQueryInterface）
- 字段映射配置（nullable, coverable, signable 等）

**Field 模块 (`field/`):**

- SQL 字段表达式构建
- 类型安全的字段操作
- PostgreSQL 特殊字段类型支持

**Types 模块 (`types/`):**

- PostgreSQL 专用类型系统
- JSON/JSONB 类型及表达式
- 数组、网络、范围、几何等类型支持

## 使用模式

### 1. 同包同目录生成（推荐）

```go
g := gen.NewGenerator(gen.Config{
    OutPath:      "./database",
    ModelPkgPath: "./database", // 与 OutPath 一致
})
```

### 2. 类型安全查询

```go
// 使用生成的查询接口
tbl, q := dbpkg.StudentQuery.QueryContext(ctx)
stu, err := q.Where(tbl.ID.Eq(1)).First()

// 或使用全局 Q
tbl2, q2 := dbpkg.Q.Student.QueryContext(ctx)
list, err := q2.Where(tbl2.Age.Gt(18)).Find()
```

### 3. PostgreSQL 类型扩展

```go
// JSON 查询
types.JSONQuery("attributes").HasKey("role")

// JSON 更新
types.JSONSet("attributes").Set("{age}", 20)

// 时间类型
types.Date, types.Time
```

## 开发规范

### 构建与测试

- `go build ./...` - 构建所有包
- `go test ./...` - 运行所有测试
- `go vet ./...` - 静态检查
- `golangci-lint run` - 代码检查

### 代码风格

- 遵循 Go 标准格式
- 测试文件与实现文件同目录（`*_test.go`）
- 导出标识符使用 `PascalCase`
- 包名使用小写，无下划线

### 测试要求

- 使用标准 `testing` 包
- 表格驱动测试优先
- 保持测试确定性（避免网络/DB 依赖）
- 运行 `go test -cover ./...` 检查覆盖率

## 配置文件支持

项目支持 YAML 配置文件定义：

- 字段类型映射 (`field_type`)
- 关联关系定义 (`field_relate`)
- 导入包管理 (`imports`)
- 忽略表配置 (`ignores`)

## 安全注意事项

- 不提交敏感信息或凭证
- 使用环境变量处理本地实验配置
- 运行 `go vet` 和 `golangci-lint` 捕获常见问题
- 代码生成路径中验证输入

## 主要依赖

- `gorm.io/gorm` v1.25.12 - ORM 框架
- `golang.org/x/tools` - 代码解析工具
- `github.com/google/uuid` - UUID 支持
- `gopkg.in/yaml.v3` - YAML 配置支持

## 贡献指南

- 保持提交范围小且逻辑清晰
- 提交信息使用祈使语气
- PR 描述清晰，包含行为变更
- 遵循现有代码风格和架构模式
