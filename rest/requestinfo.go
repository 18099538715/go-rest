package rest

type RequestInfo struct {
	Method      string      //'请求的方法'
	ContentType string      //'请求内容类型'
	UrlSplit    []string    //'分割的url地址，将url以/进行分割'
	HandleFunc  interface{} //处理的方法
}
