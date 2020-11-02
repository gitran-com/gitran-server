package util

import (
	"fmt"
	"net/url"

	giturls "github.com/whilp/git-urls"
)

//ParseGitURL parse a git url
func ParseGitURL(url string) (*url.URL, error) {
	return giturls.Parse(url)
}

//GenSSHGitURL gen a git url in ssh protocol using host and path
func GenSSHGitURL(host string, path string) string {
	return fmt.Sprintf("git@%s:%s", host, path)
}

//GenHTTPSGitURL gen a git url in https protocol using host and path
func GenHTTPSGitURL(host string, path string) string {
	return fmt.Sprintf("https://%s%s", host, path)
}
