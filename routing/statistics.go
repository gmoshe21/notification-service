package routing

import (
	"encoding/json"
	"log"
	"notification/conn"
	"notification/model"

	"github.com/valyala/fasthttp"
)

func getStatisticsAllNotification(ctx *fasthttp.RequestCtx) {
	/*
	рассылки {
		исполнились: 3
		исполняются: 2
		будут исполняться: 7
	}
	сообщения {
		доставились: 3
		не доставленны: 7
	}
	*/
	var stat model.Stat
	var result []byte

	err := conn.DB.QueryRow(conn.GetAllNotification).Scan(&stat.Notifications.Finished,
															&stat.Notifications.Sent,
															&stat.Notifications.Will_be_sent)
	if err != nil {
		log.Println("GET /statistic: Error:", err)
		ctx.Error("Server Error", fasthttp.StatusInternalServerError)
		return
	}

	err = conn.DB.QueryRow(conn.GetStatusMessages).Scan(&stat.Messages.Delivered,
															&stat.Messages.Failure)
	if err != nil {
		log.Println("GET /statistic: Error:", err)
		ctx.Error("Server Error", fasthttp.StatusInternalServerError)
		return
	}

	result, err = json.Marshal(stat)
	if err != nil {
		log.Println("GET /statistic: Error:", err)
		ctx.Error("Server Error", fasthttp.StatusInternalServerError)
		return
	}
	ctx.SetBody(result)
	ctx.SetStatusCode(200)
}

func getStatisticsNotification(ctx *fasthttp.RequestCtx) {
	var result []byte
	var allmessages []model.Message
	id := string(ctx.QueryArgs().Peek("id"))
	if id == "" {
		log.Println("GET /statistic/notification: Error: Bad request")
		ctx.Error("Bad request", fasthttp.StatusBadRequest)
		return
	}

	rows, err := conn.DB.Query(conn.GetAllMessage, id)
	if err != nil {
		log.Println("GET /statistic/notification: Error:", err)
		ctx.Error("Server Error", fasthttp.StatusInternalServerError)
		return
	}

	for rows.Next() {
		var message model.Message
		err := rows.Scan(message.Id,
							message.Data,
							message.Status,
							message.Id_notification,
							message.Id_client)
		if err != nil {
			log.Println("GET /statistic/notification: Error:", err)
		ctx.Error("Server Error", fasthttp.StatusInternalServerError)
			return
		}
		allmessages = append(allmessages, message)
	}

	result, err = json.Marshal(allmessages)
	if err != nil {
		log.Println("GET /statistic/notification: Error:", err)
		ctx.Error("Server Error", fasthttp.StatusInternalServerError)
		return
	}

	ctx.SetBody(result)
	ctx.SetStatusCode(200)
}