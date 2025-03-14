package device

import (
	"context"
	"errors"
	"fmt"
	"linklab/device-control-v2/base-library/database"
	"linklab/device-control-v2/base-library/database/table"
	"linklab/device-control-v2/base-library/parameter/response"
	"net/http"

	"github.com/gin-gonic/gin"
	log "github.com/sirupsen/logrus"
	"go.mongodb.org/mongo-driver/bson"
)

func boardlist(c *gin.Context) {

	userid := c.GetString("id")
	if len(userid) < 1 {
		err := errors.New("userid not exist")
		log.Error(err)
		c.SecureJSON(http.StatusOK, response.Response{Code: -1, Msg: err.Error()})
		return
	}

	cursor, err := database.Mdriver.FindElem("boards", bson.D{})
	if err != nil {
		err := fmt.Errorf("database.Mdriver.FindElem error {%v}", err)
		log.Error(err)
		c.SecureJSON(http.StatusOK, response.Response{Code: -1, Msg: err.Error()})
		return
	}

	var boards []table.Board
	if err := cursor.All(context.TODO(), &boards); err != nil {
		err := fmt.Errorf("database.Mdriver.FindElem error {%v}", err)
		log.Error(err)
		c.SecureJSON(http.StatusOK, response.Response{Code: -1, Msg: err.Error()})
		return
	}

	if len(boards) < 1 {
		err := errors.New("support board length 0 error")
		log.Error(err)
		c.SecureJSON(http.StatusOK, response.Response{Code: -1, Msg: err.Error()})
		return
	}

	result := response.BoardList{
		Boards: []response.BoardStatus{},
	}
	for _, board := range boards {
		result.Boards = append(result.Boards, response.BoardStatus{
			BoardName: board.BoardName,
			BoardType: board.BoardType,
		})
	}
	c.SecureJSON(http.StatusOK, response.Response{Code: 0, Msg: "success", Data: result})
}
