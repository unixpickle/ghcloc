package ghcloc

import (
	"encoding/json"
	"errors"
	"io/ioutil"
	"net/http"
)

func (self *Repository) Request(url string) ([]byte, error) {
	self.waitReqSem()
	defer self.doneReqSem()
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
	data, err := ioutil.ReadAll(res.Body)
	if err != nil {
		return nil, err
	}

	// Check for raw error value
	var raw map[string]interface{}
	if err := json.Unmarshal(data, &raw); err == nil {
		if message, ok := raw["message"]; ok {
			if s, ok := message.(string); ok {
				return nil, errors.New(s)
			}
		}
	}
	return data, nil
}

func (self *Repository) waitReqSem() {
	self.ReqSem <- struct{}{}
}

func (self *Repository) doneReqSem() {
	<-self.ReqSem
}
