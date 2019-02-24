package attendance

import (
	"github.com/edwintcloud/classmate/api/services/server"
	"github.com/globalsign/mgo"
	"github.com/labstack/echo"
)

var (
	db = struct {
		persons *mgo.Collection
		classes *mgo.Collection
	}{}
	s *server.Server
)

// Register registers routes with echo
func Register(svr *server.Server) {

	// setup server var
	s = svr

	// setup db collections
	db.persons = s.Db.C("persons")
	db.classes = s.Db.C("classes")

	s.Echo.POST("/api/v1/persons", CreatePerson)

	s.Echo.GET("/", func(c echo.Context) error {
		return c.JSON(200, server.Success())
	})
}

// CreatePerson is the a new person route
func CreatePerson(c echo.Context) error {
	person := Person{}

	// bind req body to person
	err := c.Bind(&person)
	if err != nil {
		return c.JSON(400, server.Error(err, 400))
	}

	// save password so we can use to authenticate
	password := person.Password

	// create new person
	err = person.Create()
	if err != nil {
		return c.JSON(400, server.Error(err, 400))
	}

	// authenticate person
	err = person.Authenticate(password)
	if err != nil {
		return c.JSON(500, server.Error(err, 500))
	}

	// set Password to ""
	person.Password = []byte("")

	// return person
	return c.JSON(200, person)
}

// LoginPerson generates a jwt for subsequent interaction with the server
func LoginPerson(c echo.Context) error {
	person := Person{}

	// bind req body to person
	err := c.Bind(&person)
	if err != nil {
		return c.JSON(400, server.Error(err, 400))
	}

	// authenticate person
	err = person.Authenticate(person.Password)
	if err != nil {
		return c.JSON(421, server.Error(err, 421))
	}

	// set Password to ""
	person.Password = []byte("")

	// return person
	return c.JSON(200, person)
}