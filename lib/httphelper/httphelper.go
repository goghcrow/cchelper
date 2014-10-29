package httphelper

import (
	"io/ioutil"
	"lib/dbgutil"
	"net/http"
	"net/http/cookiejar"
	"regexp"
)

func Get() {
	cookieJar, _ := cookiejar.New(nil)

	httpClient := &http.Client{Jar: cookieJar}
	httpReq, _ := http.NewRequest("GET", "https://passport.jd.com/uc/login", nil)
	httpRes, _ := httpClient.Do(httpReq)
	defer httpRes.Body.Close()
	body, _ := ioutil.ReadAll(httpRes.Body)
	uuidR, _ := regexp.Compile(`<input type="hidden" id="uuid" name="uuid" value="(.*?)"/>`)
	uuid := uuidR.FindStringSubmatch(string(body))
	dbgutil.FormatDisplay("uuid", uuid[1])

	dbgutil.FormatDisplay("url", httpReq.URL)
	cookies := cookieJar.Cookies(httpReq.URL)
	dbgutil.FormatDisplay("uuid", cookies)

}
