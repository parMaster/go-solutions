package sandbox

import (
	"encoding/json"
	"encoding/xml"
	"fmt"
	"testing"

	log "github.com/go-pkgz/lgr"
	"github.com/stretchr/testify/assert"
)

type User struct {
	Name string `json:"name"`
	Age  int    `json:"age,omitempty"`
	Pass string `json:"-"`
}

func Test_JSON(t *testing.T) {

	testUser := User{
		Name: "User1",
		Age:  123,
		Pass: "secret",
	}

	assert.NotNil(t, testUser)

	j, err := json.Marshal(testUser)
	assert.NoError(t, err)
	assert.NotEmpty(t, j)

	log.Printf("%s", j)

	var u User

	err = json.Unmarshal(j, &u)

	assert.Equal(t, testUser.Name, u.Name)
	assert.Equal(t, "sandbox.User", fmt.Sprintf("%T", u))

	j, err = json.MarshalIndent(testUser, "", "\t")

	log.Printf("%s", j)
}

type XMLUser struct {
	Name    string `xml:"name"`
	Age     int    `xml:"age,omitempty"`
	Pass    string `xml:"-"`
	Comment string `xml:",comment"` // this is crazy
}

func Test_XML(t *testing.T) {

	testUser := XMLUser{
		Name:    "User1",
		Age:     123,
		Pass:    "secret",
		Comment: "This is a comment",
	}

	assert.NotNil(t, testUser)

	x, err := xml.MarshalIndent([]XMLUser{testUser, testUser}, "", "\t")
	assert.NoError(t, err)

	log.Printf("%s %s", xml.Header, x)

	/*
		<?xml version="1.0" encoding="UTF-8"?>
		<XMLUser>
			<name>User1</name>
			<age>123</age>
			<!--This is a comment-->
		</XMLUser>
		<XMLUser>
			<name>User1</name>
			<age>123</age>
			<!--This is a comment-->
		</XMLUser>
	*/

	type XMLUsers struct {
		Users []XMLUser `xml:"user"`
	}

	x, err = xml.MarshalIndent(XMLUsers{Users: []XMLUser{testUser, testUser}}, "", "\t")
	assert.NoError(t, err)

	log.Printf("%s %s", xml.Header, x)
	/*
		<?xml version="1.0" encoding="UTF-8"?>
		<XMLUsers>
		<user>
			<name>User1</name>
			<age>123</age>
			<!--This is a comment-->
		</user>
		<user>
			<name>User1</name>
			<age>123</age>
			<!--This is a comment-->
		</user>
		</XMLUsers>
	*/
}
