package alarm

import (
	"bytes"
	"encoding/json"
	"net/http"
	"time"

	"github.com/pkg/errors"
	"github.com/sirupsen/logrus"
)

type Alarm interface {
	Notify(msg chan string)
}

var _ Alarm = &Weixin{}

type Weixin struct {
	hubAPI string
	towx   string

	c *http.Client
}

func NewWeixin(hubAPI, towx string) (Alarm, error) {
	if hubAPI == "" || towx == "" {
		return nil, errors.New("hub_api or towx is empty")
	}

	wx := &Weixin{
		hubAPI,
		towx,
		&http.Client{Timeout: 5 * time.Second},
	}

	// resp, err := wx.c.Head(hubAPI)
	// if err != nil {
	// 	return nil, err
	// }
	// defer resp.Body.Close()

	// if code := resp.StatusCode; code != http.StatusOK {
	// 	return nil, errors.Errorf("weixin_api status_code: %v", code)
	// }

	return wx, nil
}

// Notify to start a goroutine to run handle msg
func (wx *Weixin) Notify(msgs chan string) {
	go func() {

		for msg := range msgs {

			if err := wx.sendMsg(msg); err != nil {
				logrus.Errorf("send msg error: %v", err)
			}
		}
	}()
}

func (wx *Weixin) sendMsg(msg string) error {
	v := map[string]interface{}{
		"towx": wx.towx,
		"text": msg,
	}
	data, _ := json.Marshal(v)

	resp, err := wx.c.Post(wx.hubAPI, "application/json", bytes.NewReader(data))
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	if code := resp.StatusCode; code != http.StatusOK {
		return errors.Errorf("send msg status_code=%d", code)
	}

	return nil
}
