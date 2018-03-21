package webrestful

import (
	"net/http"
	"reflect"
	"strings"
)

/**
add routes
**/
func Route(requestUrl string, reqMethod string, contentType string, handleFunc interface{}) {
	if !(reqMethod == http.MethodGet || reqMethod == http.MethodPost) {
		panic("util now only post and get method is supported")
	}
	if http.MethodPost == reqMethod && contentType != "application/json" {
		panic("post method must have application/json contentType")
	}
	_, ok, _ := getMapping(requestUrl, reqMethod)
	if ok == 200 {
		panic("url conflict:" + requestUrl)
	}
	t := reflect.TypeOf(handleFunc)
	argslen := t.NumIn()
	for i := 0; i < argslen; i++ {
		checkParamType(t.In(i))
		if t.In(i).Kind() == reflect.Struct || ((t.In(i).Kind() == reflect.Ptr) && (t.In(i).Elem().Kind() == reflect.Struct)) {
			if i != argslen-1 {
				panic("the struct param must be the last and only one is allowed," + requestUrl)
			}
		}

	}
	if strings.Index(requestUrl, "{") != -1 {
		s := strings.Split(requestUrl, "/")
		v, ok := PathUrlMapping[requestUrl]
		if ok {
			PathUrlMapping[requestUrl] = append(v, &RequestInfo{Method: reqMethod, ContentType: contentType, HandleFunc: handleFunc, UrlSplit: s})
		} else {
			PathUrlMapping[requestUrl] = []*RequestInfo{&RequestInfo{Method: reqMethod, ContentType: contentType, HandleFunc: handleFunc, UrlSplit: s}}
		}
	} else {
		v, ok := UrlMapping[requestUrl]
		if ok {
			UrlMapping[requestUrl] = append(v, &RequestInfo{Method: reqMethod, ContentType: contentType, HandleFunc: handleFunc})
		} else {
			UrlMapping[requestUrl] = []*RequestInfo{&RequestInfo{Method: reqMethod, ContentType: contentType, HandleFunc: handleFunc}}
		}
	}

}
func checkParamType(rt reflect.Type) {
	var paramType reflect.Kind
	if rt.Kind() != reflect.Ptr {
		paramType = rt.Kind()
	} else {
		paramType = rt.Elem().Kind()
	}
	switch paramType {
	case reflect.Int:
	case reflect.Int8:
	case reflect.Int32:
	case reflect.Int64:
	case reflect.String:
	case reflect.Struct:
	default:
		panic("param type not support ,int int8,int32,int64,bool,string,struct now support")
	}

}
