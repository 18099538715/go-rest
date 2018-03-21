package webrestful

import (
	"encoding/json"
	"errors"
	"io/ioutil"
	"net/http"
	"reflect"
	"strconv"
)

type Handler struct {
}

func (c Handler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	reqestInfo, ok, pathParams := getMapping(r.RequestURI, r.Method)
	if ok == 404 {
		notFound(w)
		return
	}
	if ok == 405 {
		methodNotSupport(w)
		return
	}
	if http.MethodPost == r.Method {
		handlePost(w, r, reqestInfo, pathParams)
	}
	if http.MethodGet == r.Method {
		handleGet(w, r, reqestInfo, pathParams)
	}

}
func handlePost(w http.ResponseWriter, r *http.Request, reqestInfo *RequestInfo, pathParams []string) {
	defer func() {
		r.Body.Close()
	}()
	fv := reflect.ValueOf(reqestInfo.HandleFunc)
	body, _ := ioutil.ReadAll(r.Body)
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
					badRequest(w)
					return
				}

			} else {
				params[i], err = getParam(arg.Kind(), pathParams[i])
				if err != nil {
					badRequest(w)
					return
				}

			}

		} else { //如果是指针类型
			if arg.Elem().Kind() == reflect.Struct {
				v := reflect.New(arg.Elem()).Interface()
				err = json.Unmarshal(body, v)
				params[i] = reflect.ValueOf(v)
				if err != nil {
					badRequest(w)
					return
				}
			} else {
				param := reflect.New(arg.Elem())
				getPtrParam(param, pathParams[i])
				params[i] = param
				if err != nil {
					badRequest(w)
					return
				}
			}
		}
	}
	responses := fv.Call(params)
	res(responses, w)
	return
}
func res(responses []reflect.Value, w http.ResponseWriter) {
	for _, res := range responses {
		if res.Kind() != reflect.Ptr {
			b, err := json.Marshal(res.Interface())
			if err != nil {
				w.WriteHeader(500)
				w.Write([]byte(err.Error()))
				return
			}
			w.WriteHeader(200)
			w.Write(b)
		} else {
			b, err := json.Marshal(res.Elem().Interface())
			if err != nil {
				w.WriteHeader(500)
				w.Write([]byte(err.Error()))
				return
			}
			w.WriteHeader(200)
			w.Write(b)
		}

	}
}
func handleGet(w http.ResponseWriter, r *http.Request, reqestInfo *RequestInfo, pathParams []string) {
	fv := reflect.ValueOf(reqestInfo.HandleFunc)
	body, _ := ioutil.ReadAll(r.Body)
	argsLen := fv.Type().NumIn()
	params := make([]reflect.Value, argsLen)
	for i := 0; i < argsLen; i++ {
		arg := fv.Type().In(i)
		if arg.Kind() != reflect.Ptr { //如果不是指针类型
			v := reflect.New(arg)
			if arg.Kind() == reflect.Struct {
				json.Unmarshal(body, v.Interface())
				params[i] = v.Elem()
			} else {
				params[i] = reflect.ValueOf(pathParams[i])
			}

		} else { //如果是指针类型
			v := reflect.New(arg.Elem()).Interface()
			json.Unmarshal(body, v)
			params[i] = reflect.ValueOf(v)
		}
	}
	fv.Call(params)
	r.Body.Close()
}
func getParam(kind reflect.Kind, param string) (reflect.Value, error) {
	switch kind {
	case reflect.Int:
		{
			a, err := strconv.Atoi(param)
			if err != nil {
				return reflect.Value{}, err
			} else {
				return reflect.ValueOf(a), nil
			}
		}
	case reflect.Int8:
		{
			a, err := strconv.Atoi(param)
			if err != nil {
				return reflect.Value{}, err
			} else {
				return reflect.ValueOf(a), nil
			}
		}
	case reflect.Int32:
		{
			a, err := strconv.Atoi(param)
			if err != nil {
				return reflect.Value{}, err
			} else {
				return reflect.ValueOf(a), nil
			}
		}
	case reflect.Int64:
		{
			a, err := strconv.Atoi(param)
			if err != nil {
				return reflect.Value{}, err
			} else {
				return reflect.ValueOf(a), nil
			}
		}

	case reflect.String:
		return reflect.ValueOf(param), nil
	default:
		return reflect.Value{}, errors.New("参数获取失败")
	}
}
func getPtrParam(arg reflect.Value, param string) error {
	switch arg.Elem().Kind() {
	case reflect.Int:
		{
			a, err := strconv.ParseInt(param, 10, 64)
			if err != nil {
				return err
			} else {
				arg.Elem().SetInt(a)
				return nil
			}
		}
	case reflect.Int8:
		{
			a, err := strconv.ParseInt(param, 10, 64)
			if err != nil {
				return err
			} else {
				arg.Elem().SetInt(a)
				return nil
			}
		}
	case reflect.Int32:
		{
			a, err := strconv.ParseInt(param, 10, 64)
			if err != nil {
				return err
			} else {
				arg.Elem().SetInt(a)
				return nil
			}
		}
	case reflect.Int64:
		{
			a, err := strconv.ParseInt(param, 10, 64)
			if err != nil {
				return err
			} else {
				arg.Elem().SetInt(a)
				return nil
			}
		}

	case reflect.String:
		arg.Elem().SetString(param)
		return nil
	default:
		return errors.New("参数获取失败")
	}
}
