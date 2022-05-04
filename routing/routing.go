package routing

import (
	"encoding/json"
	"log"
	"notification/conn"
	"notification/model"
	"strconv"
	"time"

	"github.com/valyala/fasthttp"
)

type request struct {
	Id uint64
	Number int
	Text string
}

type response struct {
	Code int
	Message string
}

var headerContentTypeJson = []byte("application/json")

func CheckNotification() {
	for {
		var listNotifications []model.Notification

		rows, err := conn.DB.Query(conn.CheckValidNotification)
		if err != nil {
			log.Println(err)
			return
		}
		
		for rows.Next() {
			var notification model.Notification
			var filter string
			err = rows.Scan(&notification.Id,
								&notification.Start_data,
								&notification.Text,
								&filter,
								&notification.Finish_data)
			if err != nil {
				log.Println(err)
				return
			}

			err = json.Unmarshal([]byte(filter), &notification.Filter)
			if err != nil {
				log.Println(err)
				return
			}
			
			listNotifications = append(listNotifications, notification)
		}

		for _, notification := range listNotifications {
			rows, err = conn.DB.Query(conn.GetClient, notification.Filter.Num_kod, notification.Filter.Teg)
			if err != nil {
				log.Println(err)
				return
			}
			
			for rows.Next() {
				var client model.Client
				err = rows.Scan(&client.Id,
									&client.Number,
									&client.Num_kod,
									&client.Teg,
									&client.Time_zone)
				if err != nil {
					log.Println(err)
					return
				}

				err = conn.DB.QueryRow(conn.CheckMessage, notification.Id, notification.Start_data, client.Id, "sent").Scan()
				if err.Error() == "sql: expected 1 destination arguments in Scan, not 0" {
					continue
				} else if err != nil && err.Error() != "sql: no rows in result set" {
					log.Println(err)
					return
				}

				pushMessage(notification.Text, notification.Id, client.Id, client.Number)
			}
		}
	}
}

func pushMessage(text string, id_notification string, id_client string, number int) {
	var message model.Message
	var code int

	message.Data = time.Now().Format(time.RFC3339)
	message.Id_notification = id_notification
	message.Id_client = id_client

	code, err := sendPostRequest(message.Id, text, number)
	if err != nil {
		log.Println(err)
		return
	}

	if code == 200 {
		err = conn.DB.QueryRow(conn.AddNewMessage,
			message.Data,
			"sent",
			message.Id_notification,
			message.Id_client).Scan(&message.Id)

		if err != nil {
			log.Println("pushMessage: Error:", err)
			return
		}
	} else {
		err = conn.DB.QueryRow(conn.AddNewMessage,
			message.Data,
			"fail",
			message.Id_notification,
			message.Id_client).Scan(&message.Id)

		if err != nil {
			log.Println("pushMessage: Error:", err)
			return
		}
	}
}

func sendPostRequest(id uint64, text string, number int) (int, error) {
	var statusCode int
	reqStruct := request{id,number,text}
	reqTimeout := time.Duration(1000000) * time.Millisecond

	reqBodyBytes, err := json.Marshal(reqStruct)
	if err != nil {
		log.Println(err)
		return 0, err
	}

	req := fasthttp.AcquireRequest()
	req.SetRequestURI(conn.UrlApi + strconv.FormatUint(id, 10))
	req.Header.SetMethod(fasthttp.MethodPost)
	req.Header.Add("Authorization", conn.JWT)
	req.Header.SetContentTypeBytes(headerContentTypeJson)
	req.SetBodyRaw([]byte(reqBodyBytes))
	resp := fasthttp.AcquireResponse()
	err = conn.Client.DoTimeout(req, resp, reqTimeout)
	fasthttp.ReleaseRequest(req)
	if err == nil {
		statusCode = resp.StatusCode()
		respBody := resp.Body()

		// log.Printf("DEBUG Response: %s\n", respBody)
		if statusCode == 200 {
			var respStruct response
			err = json.Unmarshal(respBody, &respStruct)
			if err != nil {
				log.Println("Error response:", err)
			}
			if respStruct.Code != 0 {
				log.Println("Error response:", respStruct.Message)
			}
		} else {
			log.Println("Error: invalid HTTP response code:", statusCode)
		}
	} else {
		log.Println("Error:", err)
		return 0, err
	}
	fasthttp.ReleaseResponse(resp)
	return statusCode, nil
}

func StartHttp() {
	handler := func (ctx *fasthttp.RequestCtx) {
		switch string(ctx.Path()) {

		case "/client/":
			switch string(ctx.Method()) {
			case "PUT":
				log.Println("PUT /client: request processing")
				changeClient(ctx)
			case "POST":
				log.Println("POST /client: request processing")
				addNewClient(ctx)
			case "DELETE":
				log.Println("DELETE /client: request processing")
				deleteClient(ctx)
			default:
				log.Println("Error: method does not exist")
			}
		

		case "/notification/":
			switch string(ctx.Method()) {
			case "PUT":
				log.Println("PUT /notification: request processing")
				changeNotification(ctx)
			case "POST":
				log.Println("POST /notification: request processing")
				addNewNotification(ctx)
			case "DELETE":
				log.Println("DELETE /notification: request processing")
				deleteNotification(ctx)
			default:
				log.Println("Error: method does not exist")
			}

		case "/statistic/":
			switch string(ctx.Method()) {
			case "GET":
				log.Println("GET /statistic: request processing")
				getStatisticsAllNotification(ctx)
			default:
				log.Println("Error: method does not exist")
			}

		default:

			switch string(ctx.Path()) {
			case "/statistic/notification/":
				switch string(ctx.Method()) {
				case "GET":
					log.Println("GET /statistic/notification: request processing")
					getStatisticsNotification(ctx)
				default:
					log.Println("Error: method does not exist")
				}

			default:
				log.Println("Error: path does not exist")
				ctx.Error("Path does not exist", fasthttp.StatusBadRequest)
			}
		}
	}
	log.Println("Server start")
	err := fasthttp.ListenAndServe(conn.ServAddr, handler)
	if err != nil {
		log.Println(err)
	}
}