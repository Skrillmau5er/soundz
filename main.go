package main

import (
    "fmt"
    "net/http"
    "os"
    "time"

    tea "github.com/charmbracelet/bubbletea"
)

const url = "https://charm.sh/"

type model struct {
	httpStatus int
	err error
}

func checkServer() tea.Msg {
	c := &http.Client{Timeout: 10 * time.Seconds}
	res, err := c.Get(url)

	if err != nil {
		return errMsg{res.StatusCode}
	}
	return statusMsg(res.StatusCode)
}

type statusMsg int

type errMsg struct{ err error }

func (e errMsg) Error() string { return e.err.Error()}