package http

import (
	"bytes"
	"crypto/tls"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/http/httputil"
	"net/url"
	"strings"
)

type Request struct {
	Method          string
	RequestUri      string
	Path            string
	Query           string
	ProtocolVersion string
	Headers         map[string]string
	Cookies         map[string]string
	Body            []byte
}

type Response struct {
	Code   int
	Length int64
	Raw    []byte
}

func SetupTransport(proxyUrl string) {
	tr := &http.Transport{
		TLSClientConfig: &tls.Config{InsecureSkipVerify: true},
	}
	if proxyUrl != "" {
		purl, _ := url.Parse(proxyUrl)
		tr.Proxy = http.ProxyURL(purl)
	}
	http.DefaultTransport = tr
}

func Parse(bs []byte) Request {
	requestLine := bytes.Split(bs, []byte("\r\n"))[0]
	method, requestUri, protocolVersion := parseRequestLine(requestLine)
	path, query := parseRequestUri(requestUri)

	headers := parseHeaders(bs)

	cookies := map[string]string{}
	if rawCookies, ok := headers["Cookie"]; ok {
		delete(headers, "Cookie")
		parseRawCookies(cookies, rawCookies)
	}

	body := extractBody(bs)
	return Request{Method: method, RequestUri: requestUri, Path: path, Query: query,
		ProtocolVersion: protocolVersion, Headers: headers, Cookies: cookies, Body: body}
}

func parseRequestLine(requestLine []byte) (method, requestUri, protocolVersion string) {
	spaceSplitted := bytes.Split(requestLine, []byte(" "))
	method = string(spaceSplitted[0])
	requestUri = string(spaceSplitted[1])
	protocolVersion = string(spaceSplitted[2])
	return
}

func parseRequestUri(requestUri string) (path, query string) {
	if i := strings.Index(requestUri, "?"); i > 0 {
		path = requestUri[:i]
		query = requestUri[i+1:]
	} else {
		path = requestUri
	}
	return
}

func parseHeaders(rawReq []byte) (headers map[string]string) {
	headers = make(map[string]string)
	for _, rawHeader := range bytes.Split(rawReq, []byte("\r\n"))[1:] {
		if len(rawHeader) == 0 {
			break
		}
		name, val := parseHeader(rawHeader)
		headers[name] = val
	}
	return
}

func parseHeader(rawHeader []byte) (name, val string) {
	colonSplitted := bytes.SplitN(rawHeader, []byte(":"), 2)
	name = string(colonSplitted[0])
	val = string(colonSplitted[1])
	val = strings.TrimSpace(val)
	return
}

func extractBody(raw []byte) []byte {
	twoRns := []byte("\r\n\r\n")
	bodyIndex := bytes.Index(raw, twoRns) + len(twoRns)
	return raw[bodyIndex:]
}

func parseRawCookies(cookies map[string]string, raw string) {
	for _, c := range strings.Split(raw, "; ") {
		key := strings.Split(c, "=")[0]
		val := strings.Split(c, "=")[1]
		cookies[key] = strings.Replace(val, "\"", "%22", -1)
	}
}

func (r Request) asHttpReq(host string) *http.Request {
	url := host + r.RequestUri
	var body io.Reader
	if len(r.Body) > 0 {
		body = bytes.NewBuffer(r.Body)
	} else {
		body = nil
	}

	req, err := http.NewRequest(r.Method, url, body)
	if err != nil {
		panic(err)
	}

	for key, val := range r.Headers {
		req.Header.Set(key, val)
	}

	for key, val := range r.Cookies {
		c := &http.Cookie{Name: key, Value: val}
		req.AddCookie(c)
	}
	return req
}

func (r Request) Send(host string) (Response, error) {
	req := r.asHttpReq(host)

	client := &http.Client{}
	res, err := client.Do(req)
	if err != nil {
		return Response{}, err
	}
	raw, err := httputil.DumpResponse(res, true)

	contentLen := res.ContentLength
	if contentLen == -1 {
		contentLen = int64(len(extractBody(raw)))
	}

	return Response{res.StatusCode, contentLen, raw}, nil
}

func (r Request) Raw(host string) []byte {
	bs, _ := httputil.DumpRequestOut(r.asHttpReq(host), true)
	return bs
}

func (r Request) WithPath(path string) Request {
	result := r.Clone()
	result.RequestUri = strings.Replace(r.RequestUri, r.Path, path, 1)
	result.Path = path
	return result
}

func (r Request) WithQuery(query string) Request {
	result := r.Clone()
	result.RequestUri = strings.Replace(r.RequestUri, r.Query, query, 1)
	result.Query = query
	return result
}

func (r Request) WithBody(body []byte) Request {
	result := r.Clone()
	result.Body = body
	return result
}

func (r Request) WithHeader(key, val string) Request {
	result := r.Clone()
	result.Headers[key] = val
	return result
}

func (r Request) WithCookie(key, val string) Request {
	result := r.Clone()
	result.Cookies[key] = val
	return result
}

func (r Request) WithCookieString(val string) Request {
	result := r.Clone()
	result.Cookies = make(map[string]string)
	parseRawCookies(result.Cookies, val)
	return result
}

func (r Request) WithHeaderString(header string) Request {
	key, val := parseHeader([]byte(header))
	result := r.Clone()
	result.Headers[key] = val
	return result
}

func (r Request) Clone() Request {
	return Request{Method: r.Method, RequestUri: r.RequestUri, Path: r.Path, Query: r.Query,
		ProtocolVersion: r.ProtocolVersion, Headers: copyMap(r.Headers), Cookies: copyMap(r.Cookies), Body: r.Body}
}

func copyMap(hs map[string]string) map[string]string {
	res := make(map[string]string)
	for k, v := range hs {
		res[k] = v
	}
	return res
}

func (r Request) HasJsonBody() bool {
	ct, ok := r.Headers["Content-Type"]
	return ok && ct == "application/json"
}

func (r Request) HasJsonCookie(key string) bool {
	cookie, ok := r.Cookies[key]
	cookie = strings.Replace(cookie, "%22", "\"", -1)
	var data interface{}
	err := json.Unmarshal([]byte(cookie), &data)
	return ok && err == nil
}

func (r Request) HasFormUrlEncodedBody() bool {
	ct, ok := r.Headers["Content-Type"]
	return ok && ct == "application/x-www-form-urlencoded"
}

func (r Request) HasMultipartFormBody() bool {
	ct, ok := r.Headers["Content-Type"]
	return ok && strings.HasPrefix(ct, "multipart/form-data")
}

func (res Response) String() string {
	return fmt.Sprintf("[Code: %v, Len: %v]", res.Code, res.Length)
}
