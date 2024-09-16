package goframework_dig

import (
	"github.com/kordar/godb"
	"go.uber.org/dig"
)

var (
	digpool = godb.NewDbPool()
)

func GetDigInstance(db string) *dig.Container {
	return digpool.Handle(db).(*dig.Container)
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
