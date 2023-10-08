package coding

import (
	"context"
	"fmt"
	"github.com/meow-pad/persian/errdef"
	"reflect"
)

// SafeRunAnyParams
//
//	@Description: 安全执行函数
//	@param function func(params ...any) error
//	@param params ...any 函数参数
//	@return resultErr error 执行结果 或 运行时异常
func SafeRunAnyParams(function func(params ...any) error, params ...any) (resultErr error) {
	if function == nil {
		resultErr = errdef.ErrInvalidParams
		return
	}
	defer func() {
		if err := recover(); err != nil {
			eErr := Cast[error](err)
			if eErr == nil {
				resultErr = fmt.Errorf("%v", err)
			} else {
				resultErr = eErr
			}
		}
	}()
	resultErr = function(params...)
	return
}

// SafeRunSimple
//
//	@Description: 安全执行函数
//	@param function func() 函数
//	@return resultErr 执行结果
func SafeRunSimple(function func()) (resultErr error) {
	if function == nil {
		return
	}
	defer func() {
		if err := recover(); err != nil {
			eErr := Cast[error](err)
			if eErr == nil {
				resultErr = fmt.Errorf("%v", err)
			} else {
				resultErr = eErr
			}
		}
	}()
	function()
	return
}

// SafeRunWithContext
//
//	@Description: 安全执行函数
//	@param function 带参数 `context.Context` 的函数
//	@param ctx	`context.Context` 对象
//	@return resultErr 执行结果 或 运行时异常
func SafeRunWithContext(function func(ctx context.Context) error, ctx context.Context) (resultErr error) {
	if function == nil {
		resultErr = errdef.ErrInvalidParams
		return
	}
	defer func() {
		if err := recover(); err != nil {
			eErr := Cast[error](err)
			if eErr == nil {
				resultErr = fmt.Errorf("%v", err)
			} else {
				resultErr = eErr
			}
		}
	}()
	resultErr = function(ctx)
	return
}

// SafeRunReflect
//
//	@Description: 以反射的方式执行函数
//	@param fn 函数
//	@param args 参数集合
//	@return resultErr 执行结果（结果集中最后一个error） 或 运行时异常
func SafeRunReflect(fn any, args ...any) (resultErr error) {
	// 函数
	fnVal := reflect.ValueOf(fn)
	if fnVal.IsNil() || !fnVal.IsValid() {
		return errdef.ErrInvalidParams
	}
	if !IsFuncType(fnVal) {
		return fmt.Errorf("SafeRunReflect: fn is not a fucntion")
	}
	// 返回值
	//t := fnVal.Type()
	//numOut := t.NumOut()
	// 限定没有返回值，或只返回一个error
	//if numOut < 0 || numOut > 1 || (numOut > 0 && !IsErrorType(t.Out(0))) {
	//	return fmt.Errorf("SafeRunReflect: fn should be func(...) of func(...)error")
	//}
	// 参数
	argValues := make([]reflect.Value, len(args))
	for i, arg := range args {
		argValues[i] = reflect.ValueOf(arg)
	}
	defer func() {
		if err := recover(); err != nil {
			eErr := Cast[error](err)
			if eErr == nil {
				resultErr = fmt.Errorf("%v", err)
			} else {
				resultErr = eErr
			}
		}
	}()
	// 执行
	result := fnVal.Call(argValues)
	rLen := len(result)
	if rLen > 0 {
		// 返回最后一个error
		resultErr = Cast[error](result[rLen-1].Interface())
	}
	return
}
