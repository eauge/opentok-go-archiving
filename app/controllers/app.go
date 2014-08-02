package controllers

import (
	"encoding/json"
	"fmt"
	"github.com/eauge/opentok"
	"github.com/revel/revel"
	"net/http"
	"os"
	"strconv"
)

var (
	apiKey     = readIntVariable("API_KEY", true)
	apiSecret  = readStringVariable("API_SECRET", true)
	ot         = opentok.OpenTok{ApiKey: apiKey, ApiSecret: apiSecret}
	options    = opentok.SessionOptions{}
	properties = opentok.TokenProperties{}
	session    *opentok.Session
	archiveId  string
)

type Message struct {
	ApiKey    int
	SessionId string
	Token     string
}

type Archives struct {
	Archives string
	Previous int
	Next     int
}

type App struct {
	*revel.Controller
}

type Html string

func (r Html) Apply(req *revel.Request, resp *revel.Response) {
	resp.WriteHeader(http.StatusOK, "text/html")
	resp.Out.Write([]byte(r))
}

func (c App) Index() revel.Result {
	return c.Render()
}

func (c App) HostView() revel.Result {
	session = getSession()

	token, err := session.GenerateToken(properties)
	if err != nil {
		panic("Token could not be created")
	}

	m := Message{
		ApiKey:    session.ApiKey,
		SessionId: session.Id,
		Token:     string(token),
	}
	return c.Render(m)
}

func (c App) ParticipantView() revel.Result {
	session = getSession()

	token, err := session.GenerateToken(properties)
	if err != nil {
		panic("Token could not be created")
	}

	m := Message{
		ApiKey:    session.ApiKey,
		SessionId: session.Id,
		Token:     string(token),
	}
	return c.Render(m)
}

func (c App) ArchiveStart() revel.Result {
	session = getSession()

	archive, err := session.StartArchive("Go-SDK")

	if err != nil {
		panic("Session could not be recorded" + err.Error())
	}
	archiveId = archive.Id
	return Html("")
}

func (c App) ArchiveStop() revel.Result {
	session = getSession()
	_, err := session.StopArchive(archiveId)

	if err != nil {
		panic("Session could not be stopped " + err.Error())
	}
	return Html("")
}

func (c App) ArchiveList(page int) revel.Result {
	var (
		count  = 10
		offset = page * count
	)
	archives, err := opentok.ListArchives(ot, count, offset)

	if err != nil {
		panic("List of archives could not be retrieved " + err.Error())
	}

	var archiveString []byte
	archiveString, err = json.Marshal(archives)

	if err != nil {
		panic("List of archives could not be decoded " + err.Error())
	}

	var next, previous int
	if len(archives) < count {
		next = page
	} else {
		next = page + 1
	}
	previous = page - 1

	if previous < 0 {
		previous = 0
	}

	a := Archives{
		Archives: string(archiveString),
		Previous: previous,
		Next:     next,
	}
	return c.Render(a)
}

func (c App) ArchiveDelete(id string) revel.Result {
	err := opentok.DeleteArchive(ot, id)

	if err != nil {
		panic("Archive could not be removed " + err.Error())
	}
	return c.Redirect("/ArchiveList/0")
}

func readIntVariable(variable string, mandatory bool) int {
	value := os.Getenv(variable)
	if len(value) == 0 && mandatory {
		panic(fmt.Sprintf("Environment variable : %s expected", variable))
	}

	intValue, err := strconv.Atoi(value)
	if err != nil {
		panic(fmt.Sprintf("Environment variable : %s was expected to be an integer : %s", variable, err))
	}
	return intValue
}

func getSession() *opentok.Session {
	if session == nil {
		var err error
		session, err = opentok.NewSession(ot, options)
		if err != nil {
			panic("Session could not be created")
		}

		if err = session.Create(); err != nil {
			panic("Session could not be created")
		}
	}
	return session
}

func readStringVariable(variable string, mandatory bool) string {
	value := os.Getenv(variable)
	if len(value) == 0 && mandatory {
		panic(fmt.Sprintf("Environment variable : %s expected", variable))
	}

	return value
}
