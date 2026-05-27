# deploy-agent 项目结构说明（目录职责 + OpenAPI 配置作用）

> 说明：本说明基于当前仓库中 `internal/*`、`open-api/*` 与 `Makefile` 的实际内容编写。

---

## 1. Go 代码目录与职责

### `cmd/`
- CLI 入口命令（使用 `spf13/cobra`）。
- 目前只有 `cmd/server.go`：暴露 `deploy-agent server`，最终调用 `internal/server.Start(...)`。

### `config/`
- 运行时配置读取（当前主要是端口读取）。
- `config/config.go`：`GetPort()` 从环境变量 `PORT` 读取，默认 `8080`。

### `internal/`
内部实现分层，按“API 生成物 / handler / 服务器启动”拆分：

#### `internal/api/`
- **oapi-codegen 生成产物**（`internal/api/server.gen.go`）。
- 主要内容：
  - `ServerInterface`：OpenAPI 路由对应的 gin handler 接口（例如 `PostAuthLogin`、`GetDeployments` 等）。
  - `RegisterHandlers(...)`：把 `ServerInterface` 的各个方法注册到 gin 路由上。
  - 还包含与 spec 相关的 `types`（models/params）——在当前实现中，handler 使用这些类型。

> 生成配置：`internal/api/cfg.yaml`

#### `internal/model/`
- **oapi-codegen 生成产物**（`internal/model/models.gen.go`）。
- 主要内容：OpenAPI `components/schemas` 对应的 Go struct 类型（例如 `Deployment`、`AuthRequest` 等）。
- 当前代码里，`internal/server/openapi_handler.go` 的 handler 仍使用 `internal/api` 的类型；`internal/model` 主要用于“模型生成物”沉淀。

> 生成配置：`internal/model/cfg.yaml`

#### `internal/server/`
- 手写的服务启动与业务 handler 实现。
- `server/server.go`：
  - 创建 gin router：`gin.Default()`
  - 创建 `OpenAPIHandler` 实例
  - 调用 `api.RegisterHandlers(r, handler)` 把 OpenAPI 接口挂到路由
  - 同时提供一个示例路由 `/ping`
- `server/openapi_handler.go`：
  - `OpenAPIHandler` 实现了 `internal/api.ServerInterface`
  - 具体业务逻辑（当前是内存 map 存储 deployments，并返回 spec 对应的结构体）

#### `internal/service/`
- 当前目录存在但本轮读取信息中未展示具体文件（你实际项目后续可在这里放业务服务层逻辑）。

---

## 2. OpenAPI 相关目录与职责（`open-api/`）

`open-api/` 目录存放 OpenAPI 契约的源文件与中间构建文件。

### `open-api/openapi.yaml`
- OpenAPI **入口文件**（root spec）。
- 当前职责：
  - 定义 `openapi`, `info`, `servers`
  - 把 `paths` 外部引用到 `paths/paths.yaml`
  - `components.schemas` 也通过外部 `$ref` 引用到 `components/schemas.yaml`

作用：让你用“入口文件”统一管理 spec，然后由 bundle/generate 工具进一步解析并内联引用。

### `open-api/paths/paths.yaml`
- 定义所有路由与操作（`/health`、`/auth/login`、`/deployments` 等）。
- 每个 API 操作里用 `$ref` 引用：
  - 请求体 schema：`../components/schemas.yaml#/schemas/...`
  - 响应 schema：`../components/schemas.yaml#/schemas/...`
  - request 的 examples：`../components/examples.yaml#/examples/...`

作用：把“HTTP API 形态”与“数据模型”分离维护。

### `open-api/components/schemas.yaml`
- 定义所有 `components.schemas` 的结构（HealthResponse、Deployment、Error…）。
- 例如：
  - `DeploymentList` 中 `items` 的类型约束
  - `CreateDeploymentRequest` 中 `config` 的 `additionalProperties`

作用：数据模型的权威来源。

### `open-api/components/examples.yaml`
- 定义 examples（`AuthRequestExample`、`CreateDeploymentExample` 等）。
- 在 `paths.yaml` 的 requestBody examples 中被引用。

作用：让文档/样例 payload 可复用且集中维护。

### `open-api/bundled.yaml`
- `swagger-cli bundle` 的输出产物。
- 它把入口 spec 中的外部 `$ref`（paths/components/examples）**内联展开**成一个完整 spec。
- `make generate` 使用它作为 `oapi-codegen` 输入。

作用：避免代码生成工具在解析外部引用时出现兼容性问题，并让生成更稳定、可复现。

---

## 3. 生成与验证脚本（`Makefile`）

### `make bundle`
- 使用 `swagger-cli` 把 `open-api/openapi.yaml` bundling 成 `open-api/bundled.yaml`。

### `make generate`
- 使用 `oapi-codegen`：
  - 生成 `internal/model/models.gen.go`（基于 `internal/model/cfg.yaml`）
  - 生成 `internal/api/server.gen.go`（基于 `internal/api/cfg.yaml`）

### `make test` / `make build`
- Go 层面编译与测试验证。

---

## 4. 建议的工程化边界（便于长期维护）
- OpenAPI 契约维护：只改 `open-api/openapi.yaml + paths/components/...`
- 生成物：`make bundle` + `make generate` 自动生成 `internal/api` / `internal/model`
- 业务逻辑：只在 `internal/server`（或进一步拆成 `internal/service`）手写实现
- 文档样例：只改 `open-api/components/examples.yaml`
