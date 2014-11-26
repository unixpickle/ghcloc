package ghcloc

import (
	"io/ioutil"
	"net/http"
)

func (self *Repository) Request(url string) ([]byte, error) {
	client := &http.Client{}
	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return nil, err
	}
	if self.UserAuth {
		req.SetBasicAuth(self.AuthUser, self.AuthPass)
	}
	res, err := client.Do(req)
	if err != nil {
		return nil, err
	}
	return ioutil.ReadAll(res.Body)
}
