package runner

import (
	"fmt"
	"net"
	"net/http"
	"runtime"
	"strings"
	"time"
)

func (userOptions *UserOptions) DoHealthCheck() string {
	var test strings.Builder
	test.WriteString(fmt.Sprintf("Version: %s\n", Version))
	test.WriteString(fmt.Sprintf("Operative System: %s\n", runtime.GOOS))
	test.WriteString(fmt.Sprintf("Architecture: %s\n", runtime.GOARCH))
	test.WriteString(fmt.Sprintf("Go Version: %s\n", runtime.Version()))
	test.WriteString(fmt.Sprintf("Compiler: %s\n", runtime.Compiler))

	var testResult string
	if userOptions.Proxy != "" {
		//auth := &proxy.Auth{}
		testResult = "Ok"
		// 设置 代理Client
		proxyClient, _ := NewClient(userOptions)
		testResult = "Ok"
		proxyRequest, err := http.NewRequest("GET", MainUrl, nil)
		if err != nil {
			testResult = fmt.Sprintf("Ko (%s)", err)
		}
		test.WriteString(fmt.Sprintf("Proxy client to %s => %s\n", userOptions.Proxy, testResult))
		testResult = "Ok"
		_, err = proxyClient.Do(proxyRequest)
		if err != nil {
			testResult = fmt.Sprintf("Ko check proxy (%s) ", err)
		}
		test.WriteString(fmt.Sprintf("Proxy connectivity to %s => %s\n", MainUrl, testResult))
	}
	c4, err := net.DialTimeout("tcp4", "zh.superchat.live:443", time.Duration(userOptions.Timeout)*time.Second)
	if err == nil && c4 != nil {
		c4.Close()
	}
	testResult = "Ok"
	if err != nil {
		testResult = fmt.Sprintf("Ko (%s)", err)
	}
	test.WriteString(fmt.Sprintf("IPv4  connectivity to zh.superchat.live:443 => %s\n", testResult))
	c6, err := net.DialTimeout("tcp6", "zh.superchat.live:443", time.Duration(userOptions.Timeout)*time.Second)
	if err == nil && c6 != nil {
		c6.Close()
	}
	testResult = "Ok"
	if err != nil {
		testResult = fmt.Sprintf("Ko (%s)", err)
	}
	test.WriteString(fmt.Sprintf("IPv6  connectivity to zh.superchat.live:443 => %s\n", testResult))
	return test.String()

}
