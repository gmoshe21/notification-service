package routing

import (
	"encoding/json"
	"log"
	"notification/conn"
	"notification/model"

	"github.com/valyala/fasthttp"
)

func addNewClient(ctx *fasthttp.RequestCtx) {
	var client model.Client
	err := json.Unmarshal(ctx.PostBody(), &client)

	if err != nil {
		log.Println("POST /client: Error:", err)
		ctx.Error("Server Error", fasthttp.StatusInternalServerError)
		return
	}

	if client.Id == "" || client.Number == 0 || client.Num_kod == 0 || client.Teg == "" || client.Time_zone == "" {
		log.Println("POST /client: Error: Bad request")
		ctx.Error("Bad request", fasthttp.StatusBadRequest)
		return
	}
	log.Println("POST /client: Adding data to the database")
	err = conn.DB.QueryRow(conn.AddNewClient, 
						client.Id,
						client.Number,
						client.Num_kod,
						client.Teg,
						client.Time_zone).Scan()

	if err.Error() != "sql: no rows in result set" {
		log.Println("POST /client: Error:", err)
		ctx.Error("Server Error", fasthttp.StatusInternalServerError)
		return
	}

	ctx.SetStatusCode(200)
}

func changeClient(ctx *fasthttp.RequestCtx) {
	var client model.Client
	err := json.Unmarshal(ctx.PostBody(), &client)

	if err != nil {
		log.Println("PUT /client: Error:", err)
		ctx.Error("Server Error", fasthttp.StatusInternalServerError)
		return
	}

	if client.Id == ""|| client.Number == 0 || client.Num_kod == 0 || client.Teg == "" || client.Time_zone == "" {
		log.Println("PUT /client: Error: Bad request")
		ctx.Error("Bad request", fasthttp.StatusBadRequest)
		return
	}

	err = conn.DB.QueryRow(conn.DeleteClient, client.Id).Scan()

	if err.Error() != "sql: no rows in result set" {
		log.Println("PUT /client: Error:", err)
		ctx.Error("Server Error", fasthttp.StatusInternalServerError)
		return
	}

	err = conn.DB.QueryRow(conn.AddNewClient, 
						client.Id,
						client.Number,
						client.Num_kod,
						client.Teg,
						client.Time_zone).Scan()

	if err.Error() != "sql: no rows in result set" {
		log.Println("PUT /client: Error:", err)
		ctx.Error("Server Error", fasthttp.StatusInternalServerError)
		return
	}

	ctx.SetStatusCode(200)
}

func deleteClient(ctx *fasthttp.RequestCtx) {
	id := string(ctx.QueryArgs().Peek("id"))

	if id == "" {
		log.Println("DELETE /client: Error: Bad request")
		ctx.Error("Bad request", fasthttp.StatusBadRequest)
		return
	}

	err := conn.DB.QueryRow(conn.DeleteClient, id).Scan()

	if err.Error() != "sql: no rows in result set" {
		log.Println("DELETE /client: Error:", err)
		ctx.Error("Server Error", fasthttp.StatusInternalServerError)
		return
	}

	ctx.SetStatusCode(200)
}