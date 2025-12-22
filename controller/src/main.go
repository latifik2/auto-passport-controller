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

	response := &types.Response{}

	for {
		passports := MakePassport(ac)
		jsonb, _ := json.Marshal(passports)

		resp, err := httpClient.Post(apiEndpoint, "application/json", bytes.NewReader(jsonb))
		if err != nil {
			slog.Error(fmt.Sprintf("Failed to send passports to API gateway: %v", err))
		}

		if resp.StatusCode < 200 || resp.StatusCode >= 300 {
			slog.Error("Error while sending passports to API", "statusCode", resp.Status)
		}

		body, err := io.ReadAll(resp.Body)
		if err != nil {
			slog.Error(fmt.Sprintf("Failed to read API gateway response body: %v", err))
		}

		if err := json.Unmarshal(body, response); err != nil {
			slog.Error("Failed to unmarshal response", "err", err, "body", string(body))
		}

		if response.Status == "ok" {
			slog.Info("Successefully recieved response from API gateway", "msg", response.Message)
		}
		if response.Status == "fail" {
			slog.Error("Posting passports to API gateway failed", "msg", response.Message)
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
