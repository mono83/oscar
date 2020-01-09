package lua

import (
	"io/ioutil"
	"net/http"
	"strconv"
	"time"

	"github.com/mono83/oscar/events"
	"github.com/mono83/oscar/util"
	"github.com/yuin/gopher-lua"
)

func lHTTPGet(L *lua.LState) int {
	tc := lContext(L)
	url := tc.Interpolate(L.CheckString(2))
	lTable := L.ToTable(3)

	headers := http.Header(map[string][]string{})
	if lTable != nil {
		lTable.ForEach(func(key lua.LValue, value lua.LValue) {
			if skey, ok := key.(lua.LString); ok {
				if svalue, ok := value.(lua.LString); ok {
					headers.Set(tc.Interpolate(string(skey)), tc.Interpolate(string(svalue)))
				}
			}
		})
	}

	// Building HTTP request
	tc.Tracef("Preparing HTTP GET request to %s", url)
	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		throwLua(L, tc, "HTTP Request build error: %s", err.Error())
		return 0
	}
	req.Header = headers

	// Filling request data into vars
	tc.Set("http.request.url", url)
	for name, hh := range req.Header {
		for _, h := range hh {
			if len(h) > 0 {
				tc.Set("http.request.header."+name, h)
			}
		}
	}

	// Sending HTTP request
	before := time.Now()
	resp, err := httpClient.Do(req)
	if err != nil {
		tc.Emit(events.RemoteRequest{Type: "http+get", URI: url, Elapsed: time.Now().Sub(before)})
		throwLua(L, tc, "HTTP Request failed: %s", err.Error())
		return 0
	}
	defer resp.Body.Close()

	bts, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		tc.Emit(events.RemoteRequest{Type: "http+get", URI: url, Elapsed: time.Now().Sub(before)})
		throwLua(L, tc, "Error reading HTTP response: %s", err.Error())
		return 0
	}

	delta := time.Now().Sub(before)

	// Filling response data into vars
	ray := util.RayExtractOrEmpty(resp.Header)
	tc.Emit(events.RemoteRequest{Type: "http+get", URI: url, Elapsed: time.Now().Sub(before), Ray: ray, Success: true})
	tc.Set("http.elapsed", strconv.Itoa(int(1000*delta.Seconds())))
	tc.Set("http.response.length", strconv.Itoa(len(bts)))
	tc.Set("http.response.code", strconv.Itoa(resp.StatusCode))
	tc.Set("http.response.body", string(bts))
	tc.Set("http.response.ray", ray)
	for name, hh := range resp.Header {
		for _, h := range hh {
			if len(h) > 0 {
				tc.Set("http.response.header."+name, h)
			}
		}
	}

	tc.Tracef(
		"HTTP request done in %s, received %d bytes with code %d",
		delta,
		len(bts),
		resp.StatusCode,
	)

	return 0
}
