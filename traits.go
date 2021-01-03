package smartdevicemanagementexporter

type ThermostatTrait struct {
	SdmDevicesTraitsInfo struct {
		CustomName string `json:"customName"`
	} `json:"sdm.devices.traits.Info"`
	SdmDevicesTraitsHumidity struct {
		AmbientHumidityPercent float64 `json:"ambientHumidityPercent"`
	} `json:"sdm.devices.traits.Humidity"`
	SdmDevicesTraitsConnectivity struct {
		Status string `json:"status"`
	} `json:"sdm.devices.traits.Connectivity"`
	SdmDevicesTraitsFan struct {
		TimerMode string `json:"timerMode"`
	} `json:"sdm.devices.traits.Fan"`
	SdmDevicesTraitsThermostatMode struct {
		Mode           string   `json:"mode"`
		AvailableModes []string `json:"availableModes"`
	} `json:"sdm.devices.traits.ThermostatMode"`
	SdmDevicesTraitsThermostatEco struct {
		AvailableModes []string `json:"availableModes"`
		Mode           string   `json:"mode"`
		HeatCelsius    float64  `json:"heatCelsius"`
		CoolCelsius    float64  `json:"coolCelsius"`
	} `json:"sdm.devices.traits.ThermostatEco"`
	SdmDevicesTraitsThermostatHvac struct {
		Status string `json:"status"`
	} `json:"sdm.devices.traits.ThermostatHvac"`
	SdmDevicesTraitsSettings struct {
		TemperatureScale string `json:"temperatureScale"`
	} `json:"sdm.devices.traits.Settings"`
	SdmDevicesTraitsThermostatTemperatureSetpoint struct {
		HeatCelsius float64 `json:"heatCelsius"`
		CoolCelsius float64 `json:"coolCelsius"`
	} `json:"sdm.devices.traits.ThermostatTemperatureSetpoint"`
	SdmDevicesTraitsTemperature struct {
		AmbientTemperatureCelsius float64 `json:"ambientTemperatureCelsius"`
	} `json:"sdm.devices.traits.Temperature"`
}
