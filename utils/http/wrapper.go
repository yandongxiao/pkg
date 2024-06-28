package http

import (
	"fmt"
	"net/http"
	url2 "net/url"
	"os/exec"
	"strings"
)

type Wrapper struct {
	// TODO: If this a private link, e.g. a GitHub private repo, we should support to use token.
	url string
}

func NewWrapper(url string) *Wrapper {
	return &Wrapper{url: url}
}

func (w *Wrapper) Reachable() error {
	url, err := url2.Parse(w.url)
	if err != nil {
		return err
	}

	// we can not send a HEAD request, because some servers do not support HEAD requests
	resp, err := http.Get(url.String())
	if err != nil {
		return err
	}
	if resp.StatusCode != 200 {
		return fmt.Errorf("url %s return %d", w.url, resp.StatusCode)
	}
	return nil
}

// ReachableByCurl checks if the url is reachable by curl.
// Note: If Reachable() is failed, you do not need to call ReachableByCurl(). The unit test has shown that Reachable()
// and ReachableByCurl() will retuen the same response code
func (w *Wrapper) ReachableByCurl() error {
	_, err := url2.Parse(w.url)
	if err != nil {
		return err
	}

	// create a process to invoke curl to check the url
	// -v is used to print the response headers to stderr
	cmd := exec.Command("curl", "-v", "-L", "-o", "/dev/null", w.url)
	output, err := cmd.CombinedOutput()
	if err != nil {
		return err
	}

	// 检查命令执行状态码
	if cmd.ProcessState.ExitCode() != 0 {
		return fmt.Errorf("the status code of curl command is %d", cmd.ProcessState.ExitCode())
	}

	strs := strings.Split(string(output), "HTTP/1.1 ")
	if len(strs) <= 1 {
		strs = strings.Split(string(output), "HTTP/2 ")
	}

	// NOTE: you must use the last element of the slice, because the last element is the response header
	statusCode := strs[len(strs)-1][:3]
	if statusCode != "200" {
		return fmt.Errorf("the http response header code by curl %s return %s", w.url, statusCode)
	}
	return nil
}
