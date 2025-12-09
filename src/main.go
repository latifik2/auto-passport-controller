package main

import (
	"auto-passport/collector"
	"auto-passport/targets"
	"auto-passport/utils"
	"encoding/json"
	"fmt"
	"log/slog"
	"os"
	"time"
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
		Config:       config,
	}

	for {
		passports := MakePassport(ac)

		// for _, passport := range passports {
		// 	// fmt.Printf("Service Type: %s, Host: %s, Version: %s, Severity: %s\n",
		// 	// 	passport.ServiceType,
		// 	// 	passport.Infrastructure.Host,
		// 	// 	passport.Version,
		// 	// 	passport.Severity)

		// 	b, _ := json.Marshal(passport)
		// 	slog.Info(string(b))
		// }

		b, _ := json.Marshal(passports)
		slog.Info(string(b))

		time.Sleep(time.Second * 60)
	}
}

func MakePassport(ac collector.AbstractCollector) []collector.CommonPassport {
	staticTargets := ac.GetStaticTargets()
	dynamicTargets := ac.GetDynamicTargets()

	var passports []collector.CommonPassport

	passports = append(passports, ac.GetMetadata(targets.ToTargets(staticTargets))...)
	passports = append(passports, ac.GetMetadata(targets.ToTargets(dynamicTargets))...)
	return passports
}
