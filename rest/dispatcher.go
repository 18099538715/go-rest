package rest

//不带路径参数的url映射
var UrlMapping map[string][]*RequestInfo = make(map[string][]*RequestInfo)

//带路径参数的url映射
var PathUrlMapping map[string][]*RequestInfo = make(map[string][]*RequestInfo)
