# goframework-dig

对 [`go.uber.org/dig`](https://github.com/uber-go/dig) 的轻量封装，配合 [`godb`](https://github.com/kordar/godb) 提供「多命名空间、多实例池化管理」能力。

在一个进程内可以维护多个 `dig.Container` 实例，每个实例通过字符串命名（namespace），统一注册、获取和调用。

> 本库也是 [`github.com/kordar/dig-starter`](https://github.com/kordar/dig-starter) 的底层依赖，如果你只需要简单的默认命名空间和配置化装配，可以直接使用 `dig-starter`。

## 安装

```bash
go get github.com/kordar/goframework-dig
```

## 核心概念

- 实例池（`digpool`）
  - 使用 `godb.NewDbPool()` 管理一组实现了 `godb.Client` 接口的对象；
  - 每个对象内部持有一个 `*dig.Container`。
- 命名空间（namespace）
  - 本质上就是实例池中的 key，通常是字符串，如 `"default"`、`"tenant-a"` 等；
  - 所有操作都以 namespace 作为入口。
- E 结尾的 API
  - 如 `ProvideEByNamespace`、`InvokeEByNamespace`；
  - 内部会记录错误日志并直接 `os.Exit(1)`，适合「启动即失败退出」的场景。

## 快速开始

```go
package main

import (
	"fmt"

	goframeworkdig "github.com/kordar/goframework-dig"
)

type Config struct {
	Name string
}

type Service struct {
	Cfg Config
}

func NewConfig() Config {
	return Config{Name: "demo"}
}

func NewService(cfg Config) *Service {
	return &Service{Cfg: cfg}
}

func main() {
	// 1. 创建 dig 容器实例（命名空间为 "default"）
	_ = goframeworkdig.AddDigInstance("default")

	// 2. 注册依赖
	_ = goframeworkdig.ProvideByNamespace("default", NewConfig)
	_ = goframeworkdig.ProvideByNamespace("default", NewService)

	// 3. 调用
	_ = goframeworkdig.InvokeByNamespace("default", func(s *Service) {
		fmt.Println("service name:", s.Cfg.Name)
	})
}
```

## API 说明

包名：`github.com/kordar/goframework-dig`

### 实例池管理

- `AddDigInstance(namespace string, options ...dig.Option) error`
  - 创建一个新的 `*dig.Container`，并以 `namespace` 为 key 加入实例池。
- `RemoveDigInstance(namespace string)`
  - 从实例池中移除指定命名空间的容器。
- `HasDigInstance(namespace string) bool`
  - 判断指定命名空间是否存在实例。
- `GetDigInstance(namespace string) *dig.Container`
  - 直接从实例池中取出容器，不做存在性检查。
- `GetDig(namespace string) *dig.Container`
  - 带存在性检查的版本：
    - 如果实例不存在，使用 `slog.Error` 记录 `"dig instance not exist"`，并 `os.Exit(1)`。

### Provide：注册依赖

- `ProvideByNamespace(namespace string, constructor interface{}, opts ...dig.ProvideOption) error`
  - 等价于在指定命名空间的 `*dig.Container` 上调用 `Provide`。
- `ProvideEByNamespace(namespace string, constructor interface{}, opts ...dig.ProvideOption)`
  - 在 `ProvideByNamespace` 出错时：
    - 记录 `slog.Error("provide failed", ...)`；
    - 调用 `os.Exit(1)`。

### Invoke：解析并调用

- `InvokeByNamespace(namespace string, function interface{}, opts ...dig.InvokeOption) error`
  - 等价于在指定命名空间的 `*dig.Container` 上调用 `Invoke`。
- `InvokeEByNamespace(namespace string, function interface{}, opts ...dig.InvokeOption)`
  - 在 `InvokeByNamespace` 出错时：
    - 记录 `slog.Error("invoke failed", ...)`；
    - 调用 `os.Exit(1)`。

### Decorate：装饰已注册依赖

- `DecorateByNamespace(namespace string, decorator interface{}, opts ...dig.DecorateOption) error`
  - 等价于在指定命名空间的 `*dig.Container` 上调用 `Decorate`。
- `DecorateEByNamespace(namespace string, decorator interface{}, opts ...dig.DecorateOption)`
  - 在 `DecorateByNamespace` 出错时：
    - 记录 `slog.Error("decorate failed", ...)`；
    - 调用 `os.Exit(1)`。

### Scope：子容器

- `ScopeByNamespace(namespace string, name string, opts ...dig.ScopeOption) *dig.Scope`
  - 获取指定命名空间容器的 `Scope`。

## 多命名空间示例

```go
// 创建两个独立的 dig 容器
_ = goframeworkdig.AddDigInstance("tenant-a")
_ = goframeworkdig.AddDigInstance("tenant-b")

// 在不同命名空间注册不同实现
_ = goframeworkdig.ProvideByNamespace("tenant-a", NewServiceForTenantA)
_ = goframeworkdig.ProvideByNamespace("tenant-b", NewServiceForTenantB)

// 按需调用
_ = goframeworkdig.InvokeByNamespace("tenant-a", func(s *Service) { /* ... */ })
_ = goframeworkdig.InvokeByNamespace("tenant-b", func(s *Service) { /* ... */ })
```

## 与 dig-starter 的关系

- `goframework-dig`：专注于
  - `dig.Container` 的多命名空间管理；
  - 统一的错误处理和日志输出。
- `dig-starter`：在此基础上进一步封装
  - 默认命名空间（`defaultNamespace`）；
  - 基于配置的模块化加载（`DigModule`）。

如果你只需要一个简单的全局容器，可以考虑直接用 `dig-starter`；如果希望把容器管理下沉到更底层或需要多容器场景，推荐直接使用 `goframework-dig`。

## License

MIT（如仓库中 LICENSE 所示）。
