package user

import (
	"encoding/json"
	"fmt"
	"linklab/device-control-v2/base-library/parameter/request"
	"linklab/device-control-v2/base-library/parameter/response"
	"net/http"

	"github.com/gorilla/websocket"
)

// GetEdgePodLog 获取边缘Pod的运行日志
func (ud *Driver) GetEdgePodLog(token string, clientid string, pod string, container string) (*websocket.Conn, error) {

	wsReq := &request.UserPodLog{
		ClientID:  clientid,
		Pod:       pod,
		Container: container,
	}
	wsUrl := fmt.Sprintf("%v?%v", ud.info.EdgeNode.PodLogURL, wsReq.QueryRaw())
	wsHandler, resp, err := websocket.DefaultDialer.Dial(wsUrl, http.Header{"Authorization": []string{token}})
	if err != nil {
		if resp != nil {
			respbody := make([]byte, resp.ContentLength)
			resp.Body.Read(respbody)
			respMsg := response.Response{}
			if err := json.Unmarshal(respbody, &respMsg); err != nil {
				err := fmt.Errorf("json.Unmarshal error {%v}", err)
				return nil, err
			}
			err = fmt.Errorf("%v", respMsg.Msg)
		}
		return nil, err
	}

	return wsHandler, nil
}
