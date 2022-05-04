package conn

import (
	"io/ioutil"
	"log"
	"os"
	"time"

	"github.com/jmoiron/sqlx"
	"github.com/joho/godotenv"
	_ "github.com/lib/pq"
	"github.com/valyala/fasthttp"
)

const (
	createTable = "sql/createTable.sql"
	addNewclient = "sql/addNewClient.sql"
	deleteClient = "sql/deleteClient.sql"
	addNewNotification = "sql/addNewNotification.sql"
	deleteNotification = "sql/deleteNotification.sql"
	addNewMessage = "sql/addNewMessage.sql"
	changeMessage = "sql/changeMessage.sql"
	checkValidNotification = "sql/checkValidNotification.sql"
	getClient = "sql/getClient.sql"
	checkMessage = "sql/checkMessage.sql"
	getAllNotification = "sql/getAllNotification.sql"
	getStatusMessages = "sql/getStatusMessages.sql"
	getAllMessage = "sql/getAllMessage.sql"
)

var (
	CreateTable string
	AddNewClient string
	DeleteClient string
	AddNewNotification string
	DeleteNotification string
	AddNewMessage string
	ChangeMessage string
	CheckValidNotification string
	GetClient string
	CheckMessage string
	GetAllNotification string
	GetStatusMessages string
	GetAllMessage string
	ServAddr string
	UrlApi string
	JWT string
	Client *fasthttp.Client
	DB *sqlx.DB
)

func parce_files() {
	data, err := ioutil.ReadFile(createTable)
	if err != nil {
		log.Fatalln(err)
	}
	CreateTable = string(data);

	data, err = ioutil.ReadFile(addNewclient)
	if err != nil {
		log.Fatalln(err)
	}
	AddNewClient = string(data);

	data, err = ioutil.ReadFile(deleteClient)
	if err != nil {
		log.Fatalln(err)
	}
	DeleteClient = string(data);

	data, err = ioutil.ReadFile(addNewNotification)
	if err != nil {
		log.Fatalln(err)
	}
	AddNewNotification = string(data);

	data, err = ioutil.ReadFile(deleteNotification)
	if err != nil {
		log.Fatalln(err)
	}
	DeleteNotification = string(data);

	data, err = ioutil.ReadFile(addNewMessage)
	if err != nil {
		log.Fatalln(err)
	}
	AddNewMessage = string(data);

	data, err = ioutil.ReadFile(changeMessage)
	if err != nil {
		log.Fatalln(err)
	}
	ChangeMessage = string(data);

	data, err = ioutil.ReadFile(checkValidNotification)
	if err != nil {
		log.Fatalln(err)
	}
	CheckValidNotification = string(data);

	data, err = ioutil.ReadFile(getClient)
	if err != nil {
		log.Fatalln(err)
	}
	GetClient = string(data);

	data, err = ioutil.ReadFile(checkMessage)
	if err != nil {
		log.Fatalln(err)
	}
	CheckMessage = string(data);

	data, err = ioutil.ReadFile(getAllNotification)
	if err != nil {
		log.Fatalln(err)
	}
	GetAllNotification = string(data);

	data, err = ioutil.ReadFile(getStatusMessages)
	if err != nil {
		log.Fatalln(err)
	}
	GetStatusMessages = string(data);

	data, err = ioutil.ReadFile(getAllMessage)
	if err != nil {
		log.Fatalln(err)
	}
	GetAllMessage = string(data);
}

func init() {
	parce_files()
	err := godotenv.Load("arg.env")

	if err != nil {
		log.Println(err)
	}

	ServAddr = os.Getenv("SERV_ADDR")
	if ServAddr == "" {
		ServAddr = "127.0.0.1:5000"
	}

	host := os.Getenv("DB_HOST")
	if host == "" {
		host = "127.0.0.1"
	}

	port := os.Getenv("DB_PORT")
	if port == "" {
		port = "5432"
	}

	user := os.Getenv("DB_USERNAME")
	if user == "" {
		user = "url_shortened"
	}

	password := os.Getenv("DB_PASSWORD")
	if password == "" {
		password = "url_shortened"
	}

	dbName := os.Getenv("DB_NAME")
	if dbName == "" {
		dbName = "url_shortened"
	}

	JWT = os.Getenv("JWT")
	if JWT == "" {
		JWT = "Bearer eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJleHAiOjE2ODE4MDc2MjcsImlzcyI6ImZhYnJpcXVlIiwibmFtZSI6IkFsZWtzZXkifQ.Ojk6divMeApFB3Jzk1dSvlfsjnrOQ2omuzwd3z-Tvf0"
	}

	UrlApi = os.Getenv("URL_API")
	if UrlApi == "" {
		UrlApi = "https://probe.fbrq.cloud/v1/send/"
	}

	DB, err = sqlx.Open("postgres", "host="+host+" port="+port+" user="+user+" password="+password+" dbname="+dbName+" sslmode=disable");
	if err != nil {
		log.Panicln(err)
	}
	err = DB.Ping()
	if err != nil {
		log.Panicln(err)
	}
	DB.MustExec(CreateTable)

	readTimeout, _ := time.ParseDuration("500ms")
	writeTimeout, _ := time.ParseDuration("500ms")
	maxIdleConnDuration, _ := time.ParseDuration("1h")
	Client = &fasthttp.Client{
		ReadTimeout:                   readTimeout,
		WriteTimeout:                  writeTimeout,
		MaxIdleConnDuration:           maxIdleConnDuration,
		NoDefaultUserAgentHeader:      true, // Don't send: User-Agent: fasthttp
		DisableHeaderNamesNormalizing: true, // If you set the case on your headers correctly you can enable this
		DisablePathNormalizing:        true,
		// increase DNS cache time to an hour instead of default minute
		Dial: (&fasthttp.TCPDialer{
			Concurrency:      4096,
			DNSCacheDuration: time.Hour,
		}).Dial,
	}
}