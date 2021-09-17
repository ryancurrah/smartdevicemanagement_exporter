package smartdevicemanagementexporter

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"time"

	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promauto"
	"github.com/ryancurrah/smartdevicemanagement_exporter/partnerconnmanager"
	"golang.org/x/oauth2"
	"golang.org/x/oauth2/google"
	"google.golang.org/api/option"
	"google.golang.org/api/smartdevicemanagement/v1"
)

var (
	authorizationCodeCheckSleep                         = time.Second * 10
	namespace                                           = "sdm"
	thermostatTrait                                     = "sdm.devices.types.THERMOSTAT"
	thermostatTemperatureAmbientTemperatureCelsiusGauge = promauto.NewGaugeVec(prometheus.GaugeOpts{
		Namespace: namespace,
		Name:      "thermostat_temperature_ambientTemperatureCelsius",
		Help:      "Temperature in degrees Celsius, measured at the device.",
	}, []string{
		"CustomName",
		"Name",
		"Room",
		"Type",
	})
	thermostatThermostatTemperatureSetpointHeatCelsiusGauge = promauto.NewGaugeVec(prometheus.GaugeOpts{
		Namespace: namespace,
		Name:      "thermostat_thermostatTemperatureSetpoint_heatCelsius",
		Help:      "Target temperature in Celsius for thermostat HEAT and HEATCOOL modes.",
	}, []string{
		"CustomName",
		"Name",
		"Room",
		"Type",
	})
	thermostatThermostatTemperatureSetpointCoolCelsiusGauge = promauto.NewGaugeVec(prometheus.GaugeOpts{
		Namespace: namespace,
		Name:      "thermostat_thermostatTemperatureSetpoint_coolCelsius",
		Help:      "Target temperature in Celsius for thermostat COOL and HEATCOOL modes.",
	}, []string{
		"CustomName",
		"Name",
		"Room",
		"Type",
	})
	//https://developers.google.com/nest/device-access/traits/device/fan
	thermostatThermostatFanGauge = promauto.NewGaugeVec(prometheus.GaugeOpts{
		Namespace: namespace,
		Name:      "thermostat_thermostatFan",
		Help:      "Thermostat fan status 0 (OFF) - 1 (ON)",
	}, []string{
		"CustomName",
		"Name",
		"Room",
		"Type",
	})
	//https://developers.google.com/nest/device-access/traits/device/thermostat-hvac
	thermostatThermostatHVACGauge = promauto.NewGaugeVec(prometheus.GaugeOpts{
		Namespace: namespace,
		Name:      "thermostat_thermostatHVAC",
		Help:      "Thermostat fan status -1 (COOLING) 0 (OFF) - 1 (HEATING)",
	}, []string{
		"CustomName",
		"Name",
		"Room",
		"Type",
	})
	thermostatHumidityAmbientHumidityPercentGauge = promauto.NewGaugeVec(prometheus.GaugeOpts{
		Namespace: namespace,
		Name:      "thermostat_humidity_ambientHumidityPercent",
		Help:      "Percent humidity, measured at the device.",
	}, []string{
		"CustomName",
		"Name",
		"Room",
		"Type",
	})
)

type SmartDeviceManagementExporter struct {
	AuthorizationCodeChan chan partnerconnmanager.AuthorizationCode
	Config                *oauth2.Config
	Ctx                   context.Context
	ProjectID             string
	RefreshTokenFile      string
	RecordMetricsDelay    time.Duration
	client                *smartdevicemanagement.Service
}

func (s *SmartDeviceManagementExporter) IsClientRunning() bool {
	return s.client != nil
}

func (s *SmartDeviceManagementExporter) Start() error {
	s.recordMetrics()

	refreshTokenFile, err := ioutil.ReadFile(s.RefreshTokenFile)
	if err != nil {
		switch err.(type) {
		case *os.PathError:
			// Do nothing
		default:
			return err
		}
	}

	if refreshTokenFile != nil {
		err = s.authenticateWithRefreshToken(refreshTokenFile)
		if err != nil {
			return err
		}
	}

	for {
		select {
		case authorizationCode := <-s.AuthorizationCodeChan:
			err := s.authenticateWithCode(authorizationCode)
			if err != nil {
				return err
			}
		default:
			time.Sleep(authorizationCodeCheckSleep)
		}
	}
}

func (s *SmartDeviceManagementExporter) authenticateWithCode(authorizationCode partnerconnmanager.AuthorizationCode) error {
	s.Config.RedirectURL = authorizationCode.RedirectURI

	token, err := s.Config.Exchange(s.Ctx, authorizationCode.Code)
	if err != nil {
		return err
	}

	if token.RefreshToken == "" {
		return errors.New("no refresh token was provided by the api when authenticating")
	}

	tokenStr, err := json.Marshal(&token)
	if err != nil {
		return err
	}

	err = ioutil.WriteFile(s.RefreshTokenFile, []byte(tokenStr), 0644)
	if err != nil {
		return err
	}

	service, err := smartdevicemanagement.NewService(s.Ctx, option.WithTokenSource(s.Config.TokenSource(s.Ctx, token)))
	if err != nil {
		return err
	}

	s.client = service

	return nil
}

func (s *SmartDeviceManagementExporter) authenticateWithRefreshToken(refreshTokenFile []byte) error {
	token := oauth2.Token{}

	err := json.Unmarshal(refreshTokenFile, &token)
	if err != nil {
		return err
	}

	service, err := smartdevicemanagement.NewService(s.Ctx, option.WithTokenSource(s.Config.TokenSource(s.Ctx, &token)))
	if err != nil {
		return err
	}

	s.client = service

	return nil
}

func (s *SmartDeviceManagementExporter) recordMetrics() {
	go func() {
		for {
			time.Sleep(s.RecordMetricsDelay)

			if s.client == nil {
				continue
			}

			var devices []*smartdevicemanagement.GoogleHomeEnterpriseSdmV1Device

			err := s.client.Enterprises.Devices.List(projectID(s.ProjectID)).Pages(s.Ctx, func(d *smartdevicemanagement.GoogleHomeEnterpriseSdmV1ListDevicesResponse) error {
				devices = append(devices, d.Devices...)
				return nil
			})
			if err != nil {
				log.Printf("unable to get devices: %s", err.Error())
				continue
			}

			for _, device := range devices {
				switch device.Type {
				case thermostatTrait:
					thermostatTrait := ThermostatTrait{}

					err = json.Unmarshal(device.Traits, &thermostatTrait)
					if err != nil {
						log.Printf("unable to unmarshal device trait: %s", err.Error())
						continue
					}

					labels := prometheus.Labels{
						"CustomName": thermostatTrait.SdmDevicesTraitsInfo.CustomName,
						"Name":       device.Name,
						"Type":       device.Type,
					}

					if len(device.ParentRelations) > 0 {
						labels["Room"] = device.ParentRelations[0].DisplayName
					} else {
						labels["Room"] = ""
					}

					thermostatTemperatureAmbientTemperatureCelsiusGauge.With(labels).Set(thermostatTrait.SdmDevicesTraitsTemperature.AmbientTemperatureCelsius)
					thermostatThermostatTemperatureSetpointHeatCelsiusGauge.With(labels).Set(thermostatTrait.SdmDevicesTraitsThermostatTemperatureSetpoint.HeatCelsius)
					thermostatThermostatTemperatureSetpointCoolCelsiusGauge.With(labels).Set(thermostatTrait.SdmDevicesTraitsThermostatTemperatureSetpoint.CoolCelsius)
					thermostatHumidityAmbientHumidityPercentGauge.With(labels).Set(thermostatTrait.SdmDevicesTraitsHumidity.AmbientHumidityPercent)
					thermostatThermostatFanGauge.With(labels).Set(convertFanStatus(thermostatTrait.SdmDevicesTraitsFan.TimerMode))
					thermostatThermostatHVACGauge.With(labels).Set(convertFanStatus(thermostatTrait.SdmDevicesTraitsThermostatHvac.Status))
				default:
					log.Printf("unknown device type: %s", device.Type)
					continue
				}
			}
		}
	}()
}

//https://developers.google.com/nest/device-access/traits/device/fan
func convertFanStatus(status string) float64 {
	if status == "ON" {
		return 1
	}
	return 0
}
//https://developers.google.com/nest/device-access/traits/device/thermostat-hvac
func convertHVACStatus(status string) float64 {
	if status == "HEATING" {
		return 1
	}
	if status == "COOLING" {
		return -1
	}
	return 0
}

func projectID(pid string) string {
	return fmt.Sprintf("enterprises/%s", pid)
}

func LoadOauth2Config(filename string) (*oauth2.Config, error) {
	jsonKey, err := ioutil.ReadFile(filename)
	if err != nil {
		return nil, err
	}

	return google.ConfigFromJSON(jsonKey)
}
