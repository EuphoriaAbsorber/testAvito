package Config

var Port = ":8080"
var BasePath = "/api"
var PathDocs = BasePath + "/docs"

var PathFillDB = BasePath + "/filldb"
var PathGetUsers = BasePath + "/users"
var PathUserBanner = BasePath + "/user_banner"

var Headers = map[string]string{
	"Access-Control-Allow-Credentials": "true",
	"Access-Control-Allow-Headers":     "Origin, Content-Type, accept, csrf",
	"Access-Control-Allow-Methods":     "GET, POST, DELETE, OPTIONS",
	"Content-Type":                     "application/json",
}
