package webrestful

import (
	"net/http"
	"github.com/valyala/fasthttp"
	"reflect"
	"encoding/json"
)

func FastHttpHandler(ctx *fasthttp.RequestCtx) {
	defer func() {
		if p := recover(); p != nil {
			ctx.Error("System abnormality",fasthttp.StatusInternalServerError)
		}
	}()
	reqestInfo, ok, pathParams := getMapping(string(ctx.Request.RequestURI()), string(ctx.Method()))
	if ok == 404 {
		ctx.Error("System abnormality",fasthttp.StatusNotFound)
		return
	}
	if ok == 405 {
		ctx.Error("System abnormality",fasthttp.StatusMethodNotAllowed)
		return
	}
	method:=string(ctx.Method())
	if http.MethodPost ==method {
		fasthttpHandlePost(ctx, reqestInfo, pathParams)
	}
	if http.MethodGet == method {
		fasthttpHandleGet(ctx, reqestInfo, pathParams)
	}

}

func fasthttpHandlePost(ctx *fasthttp.RequestCtx, reqestInfo *RequestInfo, pathParams []string) {
	fv := reflect.ValueOf(reqestInfo.HandleFunc)
	body :=ctx.Request.Body()
	argsLen := fv.Type().NumIn()
	params := make([]reflect.Value, argsLen)
	var err error = nil
	for i := 0; i < argsLen; i++ {
		arg := fv.Type().In(i)
		if arg.Kind() != reflect.Ptr { //如果不是指针类型
			if arg.Kind() == reflect.Struct {
				v := reflect.New(arg)
				err = json.Unmarshal(body, v.Interface())
				params[i] = v.Elem()
				if err != nil {
					ctx.Error("Bad request!",fasthttp.StatusBadRequest)
					return
				}

			} else {
				params[i], err = getParam(arg.Kind(), pathParams[i])
				if err != nil {
					ctx.Error("Bad request!",fasthttp.StatusBadRequest)
					return
				}

			}

		} else { //如果是指针类型
			if arg.Elem().Kind() == reflect.Struct {
				v := reflect.New(arg.Elem()).Interface()
				err = json.Unmarshal(body, v)
				params[i] = reflect.ValueOf(v)
				if err != nil {
					ctx.Error("Bad request!",fasthttp.StatusBadRequest)
					return
				}
			} else {
				param := reflect.New(arg.Elem())
				getPtrParam(param, pathParams[i])
				params[i] = param
				if err != nil {
					ctx.Error("Bad request!",fasthttp.StatusBadRequest)
					return
				}
			}
		}
	}
	responses := fv.Call(params)
	fastHttpRes(responses, ctx)
	return
}

func fasthttpHandleGet(ctx *fasthttp.RequestCtx, reqestInfo *RequestInfo, pathParams []string) {
	fv := reflect.ValueOf(reqestInfo.HandleFunc)
	argsLen := fv.Type().NumIn()
	params := make([]reflect.Value, argsLen)
	var err error = nil
	for i := 0; i < argsLen; i++ {
		arg := fv.Type().In(i)
		if arg.Kind() != reflect.Ptr { //如果不是指针类型
			if arg.Kind() == reflect.Struct {

			} else {
				params[i], err = getParam(arg.Kind(), pathParams[i])
				if err != nil {
					ctx.Error("Bad request!",fasthttp.StatusBadRequest)
					return
				}

			}

		} else { //如果是指针类型
			if arg.Elem().Kind() == reflect.Struct {

			} else {
				param := reflect.New(arg.Elem())
				getPtrParam(param, pathParams[i])
				params[i] = param
				if err != nil {
					ctx.Error("Bad request!",fasthttp.StatusBadRequest)
					return
				}
			}
		}
	}
	responses := fv.Call(params)
	fastHttpRes(responses, ctx)
	return
}

func fastHttpRes(responses []reflect.Value, ctx *fasthttp.RequestCtx) {
	for _, res := range responses {
		if res.Kind() != reflect.Ptr {
			b, err := json.Marshal(res.Interface())
			if err != nil {
				ctx.Error(err.Error(),fasthttp.StatusInternalServerError)
				return
			}
			ctx.Write(b)
			ctx.SetStatusCode(fasthttp.StatusOK)
		} else {
			b, err := json.Marshal(res.Elem().Interface())
			if err != nil {
				ctx.Error(err.Error(),fasthttp.StatusInternalServerError)
				return
			}
			ctx.Write(b)
			ctx.SetStatusCode(fasthttp.StatusOK)
		}

	}
}