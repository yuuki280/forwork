// internal/reflection/reflection.go
package reflection

import (
	"go/ast"
	"reflect"
	"runtime"
)

// IsExportedOrBuiltinType 检查类型是否为导出类型或内置类型
// 这个函数被服务注册过程使用，确保只有导出类型才能注册为RPC服务
func IsExportedOrBuiltinType(t reflect.Type) bool {
	return ast.IsExported(t.Name()) || t.PkgPath() == ""
}

// GetFunctionName 获取函数的全限定名
func GetFunctionName(fn interface{}) string {
	return runtime.FuncForPC(reflect.ValueOf(fn).Pointer()).Name()
}

// CreateInstance 创建指定类型的新实例
func CreateInstance(t reflect.Type) reflect.Value {
	if t.Kind() == reflect.Ptr {
		return reflect.New(t.Elem())
	}
	return reflect.New(t).Elem()
}

// CopyValue 复制一个反射值到另一个
func CopyValue(dst, src reflect.Value) bool {
	if !dst.CanSet() {
		return false
	}

	if src.Type().AssignableTo(dst.Type()) {
		dst.Set(src)
		return true
	}

	if src.Type().ConvertibleTo(dst.Type()) {
		dst.Set(src.Convert(dst.Type()))
		return true
	}

	return false
}

// WalkStructFields 遍历结构体的所有字段
func WalkStructFields(v reflect.Value, fn func(field reflect.StructField, value reflect.Value) bool) {
	t := v.Type()
	for i := 0; i < t.NumField(); i++ {
		field := t.Field(i)
		if !fn(field, v.Field(i)) {
			break
		}
	}
}
