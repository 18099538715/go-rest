package webrestful

import (
	"strings"
)

func getMapping(requestUrl string, requestMethod string) (*RequestInfo, int, []string) {
	reqestInfo, ok := UrlMapping[requestUrl]
	if ok {
		for index, _ := range reqestInfo {
			if reqestInfo[index].Method == requestMethod {
				return reqestInfo[index], 200, []string{}
			}
		}
		return nil, 405, nil
	} else {
		urls := strings.Split(requestUrl, "/")
		for _, v := range PathUrlMapping {
			pathParams := []string{}
			sourceUrls := v[0].UrlSplit
			if len(urls) != len(sourceUrls) { //长度不一样就可以继续下一个
				break
			}
			flag := false
			for l, s := range sourceUrls {
				if strings.Index(s, "{") != -1 {
					pathParams = append(pathParams, urls[l])
					continue
				}
				if s != urls[l] { //这个时候也可以下一个pathurlmapping元素了
					flag = true
					break
				}
			}
			if flag {
				break
			}
			for m, _ := range v {
				if v[m].Method == requestMethod {
					return v[m], 200, pathParams
				}
			}
			return nil, 405, nil
		}
	}
	return nil, 404, nil
}
