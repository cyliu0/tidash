package pd

import (
	"encoding/json"
	"io/ioutil"
	"net/http"
	"net/url"

	"github.com/sirupsen/logrus"
)

type PDClient struct {
	ApiAddr string
	*http.Client
}

var GPDClinet *PDClient

func InitPDClient(apiAddr string) {
	GPDClinet = &PDClient{
		ApiAddr: apiAddr,
		Client:  &http.Client{},
	}
}

func (pdc *PDClient) get(path string, query *url.Values) (body []byte, err error) {
	url := pdc.ApiAddr + path
	if query == nil {
		url = url + "?" + query.Encode()
	}
	logrus.Debugf("pdc.get URL: %v", url)
	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		logrus.Errorf("http.NewRequest failed, err: %v", err)
		return
	}
	resp, err := pdc.Do(req)
	if err != nil {
		logrus.Errorf("pdc.Do failed, err: %v", err)
		return
	}
	defer resp.Body.Close()
	body, err = ioutil.ReadAll(resp.Body)
	if err != nil {
		logrus.Errorf("ioutil.ReadAll failed, err: %v", err)
		return
	}
	logrus.Debugf("pdc.get body: %s", string(body))
	return
}

func (pdc *PDClient) GetTrend() (trend Trend, err error) {
	body, err := pdc.get("/trend", &url.Values{})
	if err != nil {
		logrus.Errorf("pdc.get failed, err: %v", err)
		return
	}
	err = json.Unmarshal(body, &trend)
	if err != nil {
		logrus.Errorf("json.Unmarshal failed, err: %v, body: %v", err, string(body))
	}
	return
}
