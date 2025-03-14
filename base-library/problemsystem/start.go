package problemsystem

import (
	"encoding/json"
	"errors"
	"fmt"

	"github.com/go-resty/resty/v2"
	log "github.com/sirupsen/logrus"
)

// Start 开始判题
func (pd *Driver) Start(waitingID string, pid string) error {
	if len(pid) < 1 {
		err := errors.New("len(pid) < 1 error")
		log.Error(err)
		return err
	}

	data := &StartRequest{
		WaitingID: waitingID,
		PID:       pid,
		Email:     false,
	}

	binary, err := json.Marshal(data)
	if err != nil {
		err := fmt.Errorf("json.Marshal error {%v}", err)
		log.Error(err)
		return err
	}

	client := resty.New()
	resp, err := client.R().
		SetBody(binary).
		SetPathParams(map[string]string{"operate": "start"}).
		Post(pd.info.URL)
	if err != nil {
		err := fmt.Errorf("http post error {%v}", err)
		log.Error(err)
		return err
	}

	result := PResult{}
	if err := json.Unmarshal(resp.Body(), &result); err != nil {
		err := fmt.Errorf("json.Unmarshal {%s} error {%v}", string(resp.Body()), err)
		log.Error(err)
		return err
	}

	if result.Result == "false" {
		err := fmt.Errorf("start error {%v}", result.Msg)
		log.Error(err)
		return err
	}

	return nil
}
