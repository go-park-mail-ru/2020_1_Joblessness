package xss

import (
	"bytes"
	"encoding/json"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/suite"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"
)

type testStruct struct {
	ID uint64 `json:"id,omitempty"`
	Login string `json:"login,omitempty"`
	Password string `json:"password,omitempty"`
	Tag string `json:"tag,omitempty"`
	Email string `json:"email,omitempty"`
	Phone string `json:"phone,omitempty"`
	Registered string `json:"registered,omitempty"`
	Avatar string `json:"avatar,omitempty"`
	Name string `json:"name,omitempty"`
	Site string `json:"site,omitempty"`
	Experiences []nestedArray `json:"experiences,omitempty"`
}

type nestedArray struct {
	ID uint64 `json:"id,omitempty"`
	Author nestedStruct `json:"author,omitempty"`
	Keywords string `json:"keywords,omitempty"`
}

type nestedStruct struct {
	ID uint64 `json:"id,omitempty"`
	Tag string `json:"tag,omitempty"`
}

type resStruct struct {
	ID uint64 `json:"id,omitempty"`
	Login string `json:"login,omitempty"`
	Password string `json:"password,omitempty"`
	Tag string `json:"tag,omitempty"`
	Email string `json:"email,omitempty"`
	Phone string `json:"phone,omitempty"`
	Registered time.Time `json:"registered,omitempty"`
	Avatar string `json:"avatar,omitempty"`
	Name string `json:"name,omitempty"`
	Site string `json:"site,omitempty"`
	Experiences []nestedArray `json:"experiences,omitempty"`
}


type userSuite struct {
	suite.Suite
	x *XssHandler
	organization testStruct
	orgExp resStruct
	organizationByte *bytes.Buffer
	wrongOrganizationByte *bytes.Buffer
}

func (suite *userSuite) SetupTest() {
	suite.x = NewXssHandler()

	organization := testStruct{
		ID: 12,
		Login:       "new username",
		Password:    "NewPassword123",
		Name:   "new name",
		Site:    "new site",
		Email:       "new email",
		Phone: "new phone number",
		Registered: "2006-01-02T15:04:05.999999999Z",
		Experiences:  []nestedArray{
			{
				ID:       1,
				Author:   nestedStruct{
					ID:  14,
					Tag: "czx",
				},
				Keywords: "awd",
			},
		},
	}
	organizationJSON, err := json.Marshal(organization)
	organizationByte := bytes.NewBuffer(organizationJSON)
	assert.NoError(suite.T(), err)

	err = json.NewDecoder(organizationByte).Decode(&suite.orgExp)
	suite.organizationByte = bytes.NewBuffer(organizationJSON)

	organization.Site = `<a href="javascript:alert('XSS1')" onmouseover="alert('XSS2')">new site<a>`
	organization.Experiences[0].Author.Tag = `<a href="javascript:alert('XSS1')" onmouseover="alert('XSS2')">czx<a>`
	organizationJSON, err = json.Marshal(organization)
	suite.wrongOrganizationByte = bytes.NewBuffer(organizationJSON)
}

func TestSuite(t *testing.T) {
	suite.Run(t, new(userSuite))
}

func (suite *userSuite) TestXssSuccess() {
	r, _ := http.NewRequest("POST", "/", suite.organizationByte)
	w := httptest.NewRecorder()

	suite.x.SanitizeMiddleware(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		var orgRes resStruct
		err := json.NewDecoder(r.Body).Decode(&orgRes)

		assert.NoError(suite.T(), err)
		assert.Equal(suite.T(), suite.orgExp, orgRes)
	})).ServeHTTP(w, r)
}

func (suite *userSuite) TestXssWrongTag() {
	r, _ := http.NewRequest("POST", "/", suite.wrongOrganizationByte)
	w := httptest.NewRecorder()

	suite.x.SanitizeMiddleware(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		var orgRes resStruct
		err := json.NewDecoder(r.Body).Decode(&orgRes)

		assert.NoError(suite.T(), err)
		assert.Equal(suite.T(), suite.orgExp, orgRes)
	})).ServeHTTP(w, r)
}

func (suite *userSuite) TestXssEmptyBody() {
	r, _ := http.NewRequest("POST", "/?a=1", bytes.NewBuffer([]byte{}))
	w := httptest.NewRecorder()

	suite.x.SanitizeMiddleware(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		a := r.FormValue("a")

		assert.Equal(suite.T(), "1", a)
	})).ServeHTTP(w, r)
}
