package main

import (
	"database/sql"
	"github.com/gin-gonic/gin"
	_ "github.com/mattn/go-sqlite3"
	"github.com/speps/go-hashids"
	"time"
)

var db *sql.DB

func main() {

	//Set up database
	MakeDatabase()

	//Set up routers
	router := MakeRouter()
	router.Run(":8080")
}

func MakeDatabase() {
	db, _ = sql.Open("sqlite3", "./urls.db")
	defer db.Close()
	create := `create table if not exists urls (longUrl text, short text, id text primary key);`
	db.Exec(create)
}

//returns a router
func MakeRouter() *gin.Engine {
	r := gin.Default()
	r.GET("/", GetIndex)
	r.LoadHTMLFiles("public/index.html", "public/show.html")
	r.GET("/:id", ExpandUrl)
	r.POST("/create", CreateShortUrl)
	return r
}

//Handles GET request for the homepage. Returns index.html.
func GetIndex(c *gin.Context) {
	c.HTML(200, "index.html", gin.H{})
}

func getOriginalUrl(id string) string {
	db, _ := sql.Open("sqlite3", "./urls.db")
	defer db.Close()

	var longUrl string

	query, _ := db.Prepare("select longUrl from urls where id = ?")
	query.QueryRow(id).Scan(&longUrl)

	return longUrl
}

//Handles GET requests for shortened urls. Resolves the shortened url
//to its long url counterpart and redirects to that url.
func ExpandUrl(c *gin.Context) {
	id := c.Params.ByName("id")
	var longUrl string
	longUrl = getOriginalUrl(id)

	if longUrl == "" {
		c.AbortWithStatus(404)
	} else {
		c.Redirect(301, longUrl)
	}
}

//Handles POST request. for creating a new short url.
func CreateShortUrl(c *gin.Context) {
	db, _ := sql.Open("sqlite3", "./urls.db")
	defer db.Close()

	var longUrl string
	longUrl = c.PostForm("longUrl")
	var id, short string
	id = MakeUniqueSlug()
	short = "http:localhost:8080/" + id

	db.Exec("INSERT into urls(longUrl, short, id) values(?,?,?)", longUrl, short, id)

	c.HTML(200, "show.html", gin.H{"short": short})
}

//Generates a unique id (the slug component of the url) for the
//short urls
func MakeUniqueSlug() string {
	var id string
	hd := hashids.NewData()
	h, _ := hashids.NewWithData(hd)
	now := time.Now()
	id, _ = h.Encode([]int{int(now.Unix())})
	return id
}
