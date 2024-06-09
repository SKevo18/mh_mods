package api_test

import (
	"os/exec"
	"time"

	"mhmods_gui/src/api"
)

func init() {
	cmd := exec.Command("bash", "webserver/start.sh --local")
	cmd.Start()

	time.Sleep(1 * time.Second)
	api.ApiServer = "http://localhost:8080"
}
