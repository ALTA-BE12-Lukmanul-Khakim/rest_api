package main

import (
	"net/http"
	"time"

	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
	"github.com/labstack/gommon/log"
	"golang.org/x/crypto/bcrypt"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
)

type User struct {
	Name     string `json:"name" form:"name"`
	Email    string `json:"email" form:"email"`
	Password string `json:"password" form:"password"`
	Hp       string `json:"hp" form:"hp"`
}

type Vendor struct {
	Name_co   string `json:"name_co" form:"name_co"`
	Expedisi  string `json:"expedisi" form:"expedisi"`
	Transpot  string `json:"tansport" form:"transport"`
	Time_go   time.Time
	Time_come time.Time
	Is_done   bool
}

func connectDBGorm() *gorm.DB {
	dsn := "root:@tcp(127.0.0.1:3306)/restapi_db?charset=utf8mb4&parseTime=True&loc=Local"
	db, _ := gorm.Open(mysql.Open(dsn), &gorm.Config{})
	return db
}

func Regist(db *gorm.DB) echo.HandlerFunc {
	return func(c echo.Context) error {
		var user User
		if err := c.Bind(&user); err != nil {
			return c.JSON(http.StatusBadRequest, map[string]interface{}{
				"message": "cannot read data",
			})
		}
		newuser := User{
			Name:     user.Name,
			Email:    user.Email,
			Password: user.Password,
			Hp:       user.Hp,
		}
		Pass, err := bcrypt.GenerateFromPassword([]byte(user.Password), bcrypt.DefaultCost) //dycrypt password
		if err != nil {
			return err
		}
		newuser.Password = string(Pass)
		if err := db.Create(&newuser); err != nil {
			return c.JSON(http.StatusInternalServerError, map[string]interface{}{
				"message": "cannot insert data",
			})
		}

		return c.JSON(http.StatusCreated, map[string]interface{}{
			"message": "success insert new user",
			"data":    newuser,
		})
	}
}

func Login(db *gorm.DB) echo.HandlerFunc {
	return func(c echo.Context) error {
		email := c.Param("email")
		password := c.Param("password")
		var resQry User

		//check password
		err := bcrypt.CompareHashAndPassword([]byte(resQry.Password), []byte(password))
		if err != nil {
			return err
		}
		//check email
		if err := db.First(&resQry, "email = ?", email).Error; err != nil {
			return c.JSON(http.StatusInternalServerError, map[string]interface{}{
				"message": "wrong email",
			})
		}

		return c.JSON(http.StatusOK, map[string]interface{}{
			"message": "success get specific data",
			"data":    resQry,
		})
	}
}

func GetAllvendor(db *gorm.DB) echo.HandlerFunc {
	return func(c echo.Context) error {
		var ven []Vendor
		if err := db.Find(&ven).Error; err != nil {
			log.Error(err.Error())
			return c.JSON(http.StatusInternalServerError, map[string]interface{}{
				"message": "error on database",
			})
		}

		return c.JSON(http.StatusOK, map[string]interface{}{
			"message": "success get all data",
			"data":    ven,
		})
	}
}

func GetDataVendor(db *gorm.DB) echo.HandlerFunc {
	return func(c echo.Context) error {
		exp := c.Param("expedisi")

		var ven Vendor
		if err := db.First(&ven, "expedisi = ?", exp).Error; err != nil {
			return c.JSON(http.StatusInternalServerError, map[string]interface{}{
				"message": "cannot select data",
			})
		}
		return c.JSON(http.StatusOK, map[string]interface{}{
			"message": "success get specific data",
			"data":    ven,
		})
	}
}

func AddVendor(db *gorm.DB) echo.HandlerFunc {
	return func(c echo.Context) error {
		var newven Vendor

		if err := c.Bind(newven).Error; err != nil {
			log.Error(err)
			return c.JSON(http.StatusBadRequest, map[string]interface{}{
				"message": "cannot read data",
			})
		}
		if err := db.Create(&newven).Error; err != nil {
			log.Error(err)
			return c.JSON(http.StatusInternalServerError, map[string]interface{}{
				"message": "cannot insert data",
			})
		}
		return c.JSON(http.StatusCreated, map[string]interface{}{
			"message": "success insert new user",
			"data":    newven,
		})
	}
}

func migrate(db *gorm.DB) {
	db.AutoMigrate(&User{})
	db.AutoMigrate(&Vendor{})

}

func main() {

	e := echo.New()
	db := connectDBGorm()
	migrate(db)
	e.Use(middleware.Logger())

	e.POST("/users", Regist(db))
	e.GET("/users", Login(db))
	e.GET("/vendors", GetAllvendor(db))
	e.GET("/vendors/:expedisi", GetDataVendor(db))
	e.POST("/vendors", AddVendor(db))
	e.Start(":8000")

}