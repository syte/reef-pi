package doser

import (
	"bytes"
	"encoding/json"
	"github.com/reef-pi/reef-pi/controller/utils"
	"testing"
)

func TestATOAPI(t *testing.T) {
	telemetry := utils.TestTelemetry()
	store, err := utils.TestDB()
	if err != nil {
		t.Fatal("Failed to create test database. Error:", err)
	}
	c, err := New(true, store, telemetry)
	if err != nil {
		t.Fatal(err)
	}
	tr := utils.NewTestRouter()
	if err := c.Setup(); err != nil {
		t.Error(err)
	}
	c.Start()
	c.LoadAPI(tr.Router)
	body := new(bytes.Buffer)
	json.NewEncoder(body).Encode(&Pump{Name: "Foo", Pin: 1})
	if err := tr.Do("PUT", "/api/doser/pumps", body, nil); err != nil {
		t.Fatal("Failed to create dosing pump using api. Error:", err)
	}
	if err := tr.Do("GET", "/api/doser/pumps/1", new(bytes.Buffer), nil); err != nil {
		t.Fatal("Failed to delete get pump using api. Error:", err)
	}
	body.Reset()
	json.NewEncoder(body).Encode(&DosingRegiment{
		Schedule: Schedule{
			Hour:   "*",
			Minute: "*",
			Day:    "*",
			Second: "0",
		},
		Enable: true,
	})
	if err := tr.Do("POST", "/api/doser/pumps/1/schedule", body, nil); err != nil {
		t.Fatal("Failed to schedule dosing pump using api. Error:", err)
	}
	body.Reset()
	json.NewEncoder(body).Encode(&CalibrationDetails{})
	if err := tr.Do("POST", "/api/doser/pumps/1/calibrate", body, nil); err != nil {
		t.Fatal("Failed to calibrate dosing pump using api. Error:", err)
	}
	if err := tr.Do("GET", "/api/doser/pumps", new(bytes.Buffer), nil); err != nil {
		t.Fatal("Failed to list dosing pumps using api. Error:", err)
	}
	if err := tr.Do("DELETE", "/api/doser/pumps/1", new(bytes.Buffer), nil); err != nil {
		t.Fatal("Failed to delete dosing pump using api. Error:", err)
	}
	defer c.Stop()
}
