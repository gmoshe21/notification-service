package routing

import (
	"encoding/json"
	"log"
	"notification/conn"
	"notification/model"

	"github.com/valyala/fasthttp"
)

func addNewNotification(ctx *fasthttp.RequestCtx) {
	var notification model.Notification
	var filter []byte

	err := json.Unmarshal(ctx.PostBody(), &notification)

	if err != nil {
		log.Println("POST /notification: Error:", err)
		ctx.Error("Server Error", fasthttp.StatusInternalServerError)
		return
	}

	if notification.Id == "" || notification.Start_data == "" || notification.Text == "" || notification.Filter.Num_kod == "" || notification.Filter.Teg == "" || notification.Finish_data == "" {
		log.Println("POST /notification: Error: Bad request")
		ctx.Error("Bad request", fasthttp.StatusBadRequest)
		return
	}

	filter, err = json.Marshal(notification.Filter)

	if err != nil {
		log.Println("POST /notification: Error:", err)
		ctx.Error("Server Error", fasthttp.StatusInternalServerError)
		return
	}

	err = conn.DB.QueryRow(conn.AddNewNotification,
						notification.Id,
						notification.Start_data,
						notification.Text,
						filter,
						notification.Finish_data).Scan()
	
	if err.Error() != "sql: no rows in result set" {
		log.Println("POST /notification: Error:", err)
		ctx.Error("Server Error", fasthttp.StatusInternalServerError)
		return
	}

	ctx.SetStatusCode(200)
}

func changeNotification(ctx *fasthttp.RequestCtx) {
	var notification model.Notification
	var filter []byte

	err := json.Unmarshal(ctx.PostBody(), &notification)
	
	if err != nil {
		log.Println("PUT /notification: Error:", err)
		ctx.Error("Server Error", fasthttp.StatusInternalServerError)
		return
	}

	if notification.Id == "" || notification.Start_data == "" || notification.Text == "" ||  notification.Filter.Num_kod == "" || notification.Filter.Teg == "" || notification.Finish_data == "" {
		log.Println("PUT /notification: Error: Bad request")
		ctx.Error("Bad request", fasthttp.StatusBadRequest)
		return
	}

	filter, err = json.Marshal(notification.Filter)

	if err != nil {
		log.Println("PUT /notification: Error:", err)
		ctx.Error("Server Error", fasthttp.StatusInternalServerError)
		return
	}

	err = conn.DB.QueryRow(conn.DeleteNotification, notification.Id).Scan()
	if err.Error() != "sql: no rows in result set" {
		log.Println("PUT /notification: Error:", err)
		ctx.Error("Server Error", fasthttp.StatusInternalServerError)
		return
	}

	err = conn.DB.QueryRow(conn.AddNewNotification,
		notification.Id,
		notification.Start_data,
		notification.Text,
		filter,
		notification.Finish_data).Scan()

	if err.Error() != "sql: no rows in result set" {
	log.Println("PUT /notification: Error:", err)
	ctx.Error("Server Error", fasthttp.StatusInternalServerError)
	return
	}

	ctx.SetStatusCode(200)
}

func deleteNotification(ctx *fasthttp.RequestCtx) {
	id := string(ctx.QueryArgs().Peek("id"))

	if id == "" {
		log.Println("DELETE /notification: Error: Bad request")
		ctx.Error("Bad request", fasthttp.StatusBadRequest)
		return
	}

	err := conn.DB.QueryRow(conn.DeleteNotification, id).Scan()

	if err.Error() != "sql: no rows in result set" {
		log.Println("DELETE /notification: Error:", err)
		ctx.Error("Server Error", fasthttp.StatusInternalServerError)
		return
	}

	ctx.SetStatusCode(200)
}