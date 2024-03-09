package lib

import (
	"fmt"
	funcLib "go-gateway/lib/func"
	dlog "go-gateway/lib/log"
	"io"
	"net/http"

	"net/url"

	"strings"
	"time"
)

func HttpGET(trace *dlog.TraceContext, urlString string, urlParams url.Values, msTimeout int, header http.Header) (*http.Response, []byte, error) {
	startTime := time.Now().UnixNano()
	client := http.Client{
		Timeout: time.Duration(msTimeout) * time.Millisecond,
	}
	urlString = AddGetDataToUrl(urlString, urlParams)
	req, err := http.NewRequest("GET", urlString, nil)
	if err != nil {
		dlog.Log.TagWarn(trace, dlog.DLTagHTTPFailed, map[string]interface{}{
			"url":       urlString,
			"proc_time": float32(time.Now().UnixNano()-startTime) / 1.0e9,
			"method":    "GET",
			"args":      urlParams,
			"err":       err.Error(),
		})
		return nil, nil, err
	}
	if len(header) > 0 {
		req.Header = header
	}
	req = addTrace2Header(req, trace)
	resp, err := client.Do(req)
	if err != nil {
		dlog.Log.TagWarn(trace, dlog.DLTagHTTPFailed, map[string]interface{}{
			"url":       urlString,
			"proc_time": float32(time.Now().UnixNano()-startTime) / 1.0e9,
			"method":    "GET",
			"args":      urlParams,
			"err":       err.Error(),
		})
		return nil, nil, err
	}
	defer resp.Body.Close()
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		dlog.Log.TagWarn(trace, dlog.DLTagHTTPFailed, map[string]interface{}{
			"url":       urlString,
			"proc_time": float32(time.Now().UnixNano()-startTime) / 1.0e9,
			"method":    "GET",
			"args":      urlParams,
			"result":    funcLib.Substr(string(body), 0, 1024),
			"err":       err.Error(),
		})
		return nil, nil, err
	}
	dlog.Log.TagInfo(trace, dlog.DLTagHTTPSuccess, map[string]interface{}{
		"url":       urlString,
		"proc_time": float32(time.Now().UnixNano()-startTime) / 1.0e9,
		"method":    "GET",
		"args":      urlParams,
		"result":    funcLib.Substr(string(body), 0, 1024),
	})
	return resp, body, nil
}

func HttpPOST(trace *dlog.TraceContext, urlString string, urlParams url.Values, msTimeout int, header http.Header, contextType string) (*http.Response, []byte, error) {
	startTime := time.Now().UnixNano()
	client := http.Client{
		Timeout: time.Duration(msTimeout) * time.Millisecond,
	}
	if contextType == "" {
		contextType = "application/x-www-form-urlencoded"
	}
	urlParamEncode := urlParams.Encode()
	req, err := http.NewRequest("POST", urlString, strings.NewReader(urlParamEncode))
	if err != nil {
		return nil, nil, err
	}
	if len(header) > 0 {
		req.Header = header
	}
	req = addTrace2Header(req, trace)
	req.Header.Set("Content-Type", contextType)
	resp, err := client.Do(req)
	if err != nil {
		dlog.Log.TagWarn(trace, dlog.DLTagHTTPFailed, map[string]interface{}{
			"url":       urlString,
			"proc_time": float32(time.Now().UnixNano()-startTime) / 1.0e9,
			"method":    "POST",
			"args":      funcLib.Substr(urlParamEncode, 0, 1024),
			"err":       err.Error(),
		})
		return nil, nil, err
	}
	defer resp.Body.Close()
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		dlog.Log.TagWarn(trace, dlog.DLTagHTTPFailed, map[string]interface{}{
			"url":       urlString,
			"proc_time": float32(time.Now().UnixNano()-startTime) / 1.0e9,
			"method":    "POST",
			"args":      funcLib.Substr(urlParamEncode, 0, 1024),
			"result":    funcLib.Substr(string(body), 0, 1024),
			"err":       err.Error(),
		})
		return nil, nil, err
	}
	dlog.Log.TagInfo(trace, dlog.DLTagHTTPSuccess, map[string]interface{}{
		"url":       urlString,
		"proc_time": float32(time.Now().UnixNano()-startTime) / 1.0e9,
		"method":    "POST",
		"args":      funcLib.Substr(urlParamEncode, 0, 1024),
		"result":    funcLib.Substr(string(body), 0, 1024),
	})
	return resp, body, nil
}

func HttpJSON(trace *dlog.TraceContext, urlString string, jsonContent string, msTimeout int, header http.Header) (*http.Response, []byte, error) {
	startTime := time.Now().UnixNano()
	client := http.Client{
		Timeout: time.Duration(msTimeout) * time.Millisecond,
	}
	req, err := http.NewRequest("POST", urlString, strings.NewReader(jsonContent))
	if err != nil {
		return nil, nil, err
	}
	if len(header) > 0 {
		req.Header = header
	}
	req = addTrace2Header(req, trace)
	req.Header.Set("Content-Type", "application/json")
	resp, err := client.Do(req)
	if err != nil {
		dlog.Log.TagWarn(trace, dlog.DLTagHTTPFailed, map[string]interface{}{
			"url":       urlString,
			"proc_time": float32(time.Now().UnixNano()-startTime) / 1.0e9,
			"method":    "POST",
			"args":      funcLib.Substr(jsonContent, 0, 1024),
			"err":       err.Error(),
		})
		return nil, nil, err
	}
	defer resp.Body.Close()
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		dlog.Log.TagWarn(trace, dlog.DLTagHTTPFailed, map[string]interface{}{
			"url":       urlString,
			"proc_time": float32(time.Now().UnixNano()-startTime) / 1.0e9,
			"method":    "POST",
			"args":      funcLib.Substr(jsonContent, 0, 1024),
			"result":    funcLib.Substr(string(body), 0, 1024),
			"err":       err.Error(),
		})
		return nil, nil, err
	}
	dlog.Log.TagInfo(trace, dlog.DLTagHTTPSuccess, map[string]interface{}{
		"url":       urlString,
		"proc_time": float32(time.Now().UnixNano()-startTime) / 1.0e9,
		"method":    "POST",
		"args":      funcLib.Substr(jsonContent, 0, 1024),
		"result":    funcLib.Substr(string(body), 0, 1024),
	})
	return resp, body, nil
}

func AddGetDataToUrl(urlString string, data url.Values) string {
	if strings.Contains(urlString, "?") {
		urlString = urlString + "&"
	} else {
		urlString = urlString + "?"
	}
	return fmt.Sprintf("%s%s", urlString, data.Encode())
}

func addTrace2Header(request *http.Request, trace *dlog.TraceContext) *http.Request {
	traceId := trace.TraceId
	cSpanId := dlog.NewSpanId()
	if traceId != "" {
		request.Header.Set("header-rid", traceId)
	}
	if cSpanId != "" {
		request.Header.Set("header-spanid", cSpanId)
	}
	trace.CSpanId = cSpanId
	return request
}
