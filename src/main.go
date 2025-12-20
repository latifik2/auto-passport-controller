package main

import (
	"auto-passport/collector"
	"auto-passport/targets"
	"auto-passport/utils"
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"log/slog"
	"net/http"
	"os"
	"time"

	"github.com/latifik2/auto-passport-controller/types"
)

func main() {
	logger := slog.New(slog.NewTextHandler(os.Stdout, nil))
	slog.SetDefault(logger)

	config := utils.Config{}
	config.ReadConfig()

	fmt.Println(config)

	clientset := collector.GetK8sClientSet()

	ac := collector.AirflowCollector{
		K8sClientSet: clientset,
		Config:       &config,
	}

	apiEndpoint := config.ApiEndpoint

	httpClient := &http.Client{}

	for {
		passports := MakePassport(ac)
		jsonb, _ := json.Marshal(passports)

		resp, err := httpClient.Post(apiEndpoint, "application/json", bytes.NewReader(jsonb))
		if err != nil {
			slog.Error(fmt.Sprintf("Failed to send passports to API gateway: %v", err))
		}

		body, err := io.ReadAll(resp.Body)
		if err != nil {
			slog.Error(fmt.Sprintf("Failed to read API gateway response body: %v", err))
		}

		

		resp.Body.Close()

		time.Sleep(time.Second * 60)
	}
}

func MakePassport(ac collector.AbstractCollector) []types.CommonPassport {
	staticTargets := ac.GetStaticTargets()
	dynamicTargets := ac.GetDynamicTargets()

	var passports []types.CommonPassport

	passports = append(passports, ac.GetMetadata(targets.ToTargets(staticTargets))...)
	passports = append(passports, ac.GetMetadata(targets.ToTargets(dynamicTargets))...)
	return passports
}
