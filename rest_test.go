package main

import (
	"github.com/freshautomations/telegram-moderator-bot/context"
	"net/http"
	"net/http/httptest"
	"testing"
)

//Todo: Fix and enhance testing. Something always falls behind...
func Test_MainHandler(t *testing.T) {
	req, err := http.NewRequest("GET", "/", nil)
	if err != nil {
		t.Fatal(err)
	}

	ctx, err := Initialization(&context.InitialContext{
		WebserverIp:    "",
		WebserverPort:  0,
		ConfigFile:     "tmb.conf.template",
		LocalExecution: true,
	})
	if err != nil {
		t.Fatal(err)
	}

	rr := httptest.NewRecorder()
	handler := context.Handler{ctx, MainHandler}

	handler.ServeHTTP(rr, req)

	if status := rr.Code; status != http.StatusOK {
		t.Errorf("handler returned wrong status code: got %v want %v",
			status, http.StatusOK)
	}

	expected := "{\"message\":\"EOF\"}\n"
	if rr.Body.String() != expected {
		t.Errorf("handler returned unexpected body: got %v want %v",
			rr.Body.String(), expected)
	}
}
