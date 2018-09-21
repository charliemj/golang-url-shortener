package main

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/jinzhu/gorm"
	_ "github.com/jinzhu/gorm/dialects/sqlite"
	"github.com/speps/go-hashids"
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

	//Routes
	r := gin.Default()
	r.GET("/", GetIndex)
	r.LoadHTMLFiles("index.html", "show.html")
	r.GET("/:id", ExpandUrl)
	r.POST("/create", CreateShortUrl)
	r.Run(":8080")
}

func GetIndex(c *gin.Context) {
	var urls []MyUrl
	if err := db.Find(&urls).Error; err != nil {
		c.AbortWithStatus(404)
		fmt.Println(err)
	} else {
		c.HTML(200, "index.html") //, gin.H{"URLS": urls}) //for debugging
	}
}

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

func CreateShortUrl(c *gin.Context) {
	var longUrl string
	longUrl = c.PostForm("longUrl")

	fmt.Println(longUrl)
	hd := hashids.NewData()

	h, _ := hashids.NewWithData(hd)
	now := time.Now()
	var id string
	var short string
	id, _ = h.Encode([]int{int(now.Unix())})
	short = "http:localhost:8080/" + id
	fmt.Println(short)
	var newUrl MyUrl
	newUrl = MyUrl{LongUrl: longUrl, ShortUrl: short, ID: id}
	fmt.Println("lolhi")
	db.Create(&newUrl)
	c.HTML(200, "show.html", gin.H{"short": short})
}
