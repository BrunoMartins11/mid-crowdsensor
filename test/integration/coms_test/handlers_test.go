package coms_test

import (
	"github.com/BrunoMartins11/mid-crowdsensor/internal/coms"
	"github.com/BrunoMartins11/mid-crowdsensor/test"
	"github.com/thanhpk/randstr"
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestAddDeviceHandler(t *testing.T) {
	_ = test.SetupToken(t)

	req, err := http.NewRequest("POST", "/addDevice?device=" + randstr.String(16), nil)
	if err != nil {
		t.Fatal(err)
	}
	rr := httptest.NewRecorder()
	handler := http.HandlerFunc(coms.AddDeviceHandler)
	handler.ServeHTTP(rr, req)
	if status := rr.Code; status != http.StatusOK {
		t.Errorf("handler returned wrong status code: got %v want %v",
			status, http.StatusOK)
	}

}
