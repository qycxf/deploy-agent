# deploy-agent

基础项目脚手架，使用 OpenAPI 作为接口契约。

## 目录说明

- `open-api/`：OpenAPI 源文件、组件和路径拆分
- `internal/api/`：OpenAPI 生成的 HTTP 路由与类型
- `internal/model/`：独立生成的模型定义
- `internal/server/`：服务启动与 OpenAPI 适配层

## 常用命令

### 生成 OpenAPI bundle
```bash
make bundle
```

### 生成代码
```bash
make generate
```

### 启动服务
```bash
go run . server
```

### 构建二进制
```bash
make build
```

### 运行测试
```bash
make test
```

## 说明

当前项目采用的是 OpenAPI-first 的 Go 服务结构：
- `open-api/` 维护契约
- `oapi-codegen` 生成接口和数据模型
- `internal/server/` 实现业务处理并挂载生成的路由
- `cmd/` 提供 CLI 入口
