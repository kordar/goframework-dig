package goframework_dig

import (
	"log/slog"
	"os"

	"github.com/kordar/godb"
	"go.uber.org/dig"
)

var (
	digpool = godb.NewDbPool()
)

func GetDigInstance(db string) *dig.Container {
	return digpool.Handle(db).(*dig.Container)
}

func GetDig(namespace string) *dig.Container {
	if !HasDigInstance(namespace) {
		slog.Error("dig instance not exist", "namespace", namespace)
		os.Exit(1)
	}
	return GetDigInstance(namespace)
}

// AddDigInstance 添加dig句柄
func AddDigInstance(db string, options ...dig.Option) error {
	ins := NewDigIns(db, options...)
	return digpool.Add(ins)
}

// RemoveDigInstance 移除dig句柄
func RemoveDigInstance(db string) {
	digpool.Remove(db)
}

// HasDigInstance dig句柄是否存在
func HasDigInstance(db string) bool {
	return digpool != nil && digpool.Has(db)
}

func ProvideByNamespace(namespace string, constructor interface{}, opts ...dig.ProvideOption) error {
	instance := GetDig(namespace)
	return instance.Provide(constructor, opts...)
}

func ProvideEByNamespace(namespace string, constructor interface{}, opts ...dig.ProvideOption) {
	err := ProvideByNamespace(namespace, constructor, opts...)
	if err != nil {
		slog.Error("provide failed", "namespace", namespace, "err", err)
		os.Exit(1)
	}
}

func InvokeByNamespace(namespace string, function interface{}, opts ...dig.InvokeOption) error {
	instance := GetDig(namespace)
	return instance.Invoke(function, opts...)
}

func InvokeEByNamespace(namespace string, function interface{}, opts ...dig.InvokeOption) {
	err := InvokeByNamespace(namespace, function, opts...)
	if err != nil {
		slog.Error("invoke failed", "namespace", namespace, "err", err)
		os.Exit(1)
	}
}

func DecorateByNamespace(namespace string, decorator interface{}, opts ...dig.DecorateOption) error {
	instance := GetDig(namespace)
	return instance.Decorate(decorator, opts...)
}

func DecorateEByNamespace(namespace string, decorator interface{}, opts ...dig.DecorateOption) {
	err := DecorateByNamespace(namespace, decorator, opts...)
	if err != nil {
		slog.Error("decorate failed", "namespace", namespace, "err", err)
		os.Exit(1)
	}
}

func ScopeByNamespace(namespace string, name string, opts ...dig.ScopeOption) *dig.Scope {
	instance := GetDig(namespace)
	return instance.Scope(name, opts...)
}
