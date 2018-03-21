package rest

import (
	"strings"
)

func getMapping(requestUrl string, requestMethod string) (*RequestInfo, int, []string) {
	pathParams := []string{}
	reqestInfo, ok := UrlMapping[requestUrl]
	if ok {
		for index, _ := range reqestInfo {
			if reqestInfo[index].Method == requestMethod {
				return reqestInfo[index], 200, pathParams
			}
		}
		return nil, 200, nil
	} else {
		urls := strings.Split(requestUrl, "/")
		for _, v := range PathUrlMapping {
			for index, _ := range v {
				if len(urls) != len(v[index].UrlSplit) { //长度不一样就继续查找
					continue
				}
				for l, s := range v[index].UrlSplit {
					if strings.Index(s, "{") != -1 {
						pathParams = append(pathParams, urls[l])
						continue
					}
					if s != urls[l] {
						return nil, 404, nil
					}
				}
				if v[index].Method == requestMethod {
					return v[index], 200, pathParams
				}
				if index == len(v)-1 {
					return nil, 405, nil
				}
			}
		}
	}
	return nil, 404, nil
}
