package main

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/jinzhu/gorm"
	_ "github.com/jinzhu/gorm/dialects/sqlite"
	"github.com/speps/go-hashids"
	//"net/http"
	"time"
)

var db *gorm.DB
var err error

type MyUrl struct {
	ID       string `json:"id,omitempty"`
	LongUrl  string `json:"longUrl,omitempty"`
	ShortUrl string `json:"shortUrl,omitempty"`
}

func main() {

	//Set up database
	db, err = gorm.Open("sqlite3", "./gorm.db")
	if err != nil {
		fmt.Println(err)
	}
	defer db.Close()

	db.AutoMigrate(&MyUrl{})

	// p1 := MyUrl{LongUrl: "http://charlie.codes", ID: "1"}
	// p2 := MyUrl{LongUrl: "http://kar.moore.com", ID: "2"}
	// db.Create(&p1)
	// db.Create(&p2)

	//Routing
	router := MakeRouter()
	router.Run(":8080")
}

//returns a router
func MakeRouter() *gin.Engine {
	r := gin.Default()
	r.GET("/", GetIndex)
	r.LoadHTMLFiles("index.html", "show.html")
	r.GET("/:id", ExpandUrl)
	r.POST("/create", CreateShortUrl)

	return r
}

//Handles GET request for the homepage. Returns index.html
func GetIndex(c *gin.Context) {
	c.HTML(200, "index.html", gin.H{})
}

//Handles GET requests for shortened urls. Resolves the shortened url
//to its long url counterpart and redirects to that url.
func ExpandUrl(c *gin.Context) {
	id := c.Params.ByName("id")
	var url MyUrl
	if err := db.Where("id = ?", id).First(&url).Error; err != nil {
		c.AbortWithStatus(404)
		fmt.Println(err)
	} else {
		c.Redirect(301, url.LongUrl)
	}
}

//Handles POST request. for creating a new short url.
func CreateShortUrl(c *gin.Context) {
	var longUrl string
	longUrl = c.PostForm("longUrl")
	var id string
	var short string
	id = MakeUniqueSlug()
	short = "http:localhost:8080/" + id
	var newUrl MyUrl
	newUrl = MyUrl{LongUrl: longUrl, ShortUrl: short, ID: id}
	db.Create(&newUrl)
	c.HTML(200, "show.html", gin.H{"short": short})
}

//Generates a unique id (the slug component of the url) for the
//short urls
func MakeUniqueSlug() string {
	var id string
	var url MyUrl
	var isAlreadyMade bool
	isAlreadyMade = true

	hd := hashids.NewData()
	h, _ := hashids.NewWithData(hd)
	//check if id is already in db and keep trying to make on until it is not
	for isAlreadyMade {
		now := time.Now()
		id, _ = h.Encode([]int{int(now.Unix())})

		//get all rows where the id is = the id we just made
		rows, _ := db.Where(map[string]interface{}{"id": id}).Find(&url).Rows()

		//if there are any rows with that id, we remake the id
		count := 0
		for rows.Next() {
			count++
		}
		if count == 0 {

			isAlreadyMade = false
		} else {
			isAlreadyMade = true
		}
	}
	return id
}
