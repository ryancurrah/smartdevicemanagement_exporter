# smartdevicemanagement_exporter

Google Nest [Smart Device Management](https://developers.google.com/nest/device-access) exporter for Prometheus.

## Prerequisites

Guides you through the prerequisite items required to start using this exporter.

### Nest account migration

Ensure you migrate your Nest account to Google if you have not already done so https://www.blog.google/products/google-nest/its-time-nest-users-can-now-switch-google-accounts.

### Sign up for Google Device Access Console

Follow the getting started guide https://developers.google.com/nest/device-access/get-started.

When you are done you should have the following items completed.

1. Register for the Device Access program.
2. Activate a supported Nest device with a Google account.
3. Create a Google Cloud Platform (GCP) project to enable the SDM API and get an OAuth 2.0 client ID.
4. Create a Device Access project to receive a Project ID.

## Installation

You can download the binary for your `OS` and `Architecture` on the releases page. A Docker image is also provided which is also found on the releases page.

https://github.com/ryancurrah/smartdevicemanagement_exporter/releases

## Configuration

All the configuration parameters for this exporter.

- `-listen-address`. Envvar: `LISTEN_ADDRESS`. Default: `:8080`. Address to listen on for HTTP requests.                       
- `-project-id`. Envvar: `PROJECT_ID`. Default: . ID of the Smart Device Management project (Required).                    
- `-credentials`. Envvar: `CREDENTIALS`. Default: `client_secret.json`. Location on disk to the Oauth2 credentials JSON file.         
- `-refresh-token`. Envvar: `REFRESH_TOKEN`. Default: `refresh_token.json`. Location on disk to store the Oauth2 refresh token JSON file. 
- `-record-metrics-delay`. Envvar: `RECORD_METRICS_DELAY`. Default: `1m`. Delay between queries to the Smart Device Management API for recording metrics.

## Endpoints

All the HTTP endpoints for this exporter.

- `/metrics`. Prometheus metrics are exposed here for scraping.
- `/authstatus`. Provides the current authentication status of the exporter with Google APIs.
- `/authorize`. Starts or resets the authentication of the exporter with the Partner Connection Manager.

## Usage

Assuming you have stored the `client_secret.json` file in the current working directory start the exporter.

```shell
./smartdevicemanagement_exporter -project-id e14044ff-a995-4733-9672-ddad34c970d5
2021/01/03 14:19:48 running
```

Now authorize the exporter to use the Google APIs by visiting the authorize endpoint in your browser. You only have to do this once or when the `refresh_token.json` file is no longer valid.

```shell
open http://127.0.0.1:8080/authorize
```

Once your done filling out the prompts that authorize the exporter you will be redirected back to the exporter on a page that says: `authorization code received from partner connection manager`.

You can check the `stdout` and `/authstatus` page to see if the exporter has been successfully authorized.

```shell
open http://127.0.0.1:8080/authstatus
```

Now the exporter will start recording metrics.

```shell
open http://127.0.0.1:8080/metrics
```

## Metrics

Currently only Thermostat devices are supported.

```
# HELP sdm_thermostat_humidity_ambientHumidityPercent Percent humidity, measured at the device.
# TYPE sdm_thermostat_humidity_ambientHumidityPercent gauge
sdm_thermostat_humidity_ambientHumidityPercent{CustomName="",Name="DU1AeyLeA",Room="Hallway",Type="sdm.devices.types.THERMOSTAT"} 34

# HELP sdm_thermostat_temperature_ambientTemperatureCelsius Temperature in degrees Celsius, measured at the device.
# TYPE sdm_thermostat_temperature_ambientTemperatureCelsius gauge
sdm_thermostat_temperature_ambientTemperatureCelsius{CustomName="",Name="DU1AeyLeA",Room="Hallway",Type="sdm.devices.types.THERMOSTAT"} 22.189987

# HELP sdm_thermostat_thermostatTemperatureSetpoint_coolCelsius Target temperature in Celsius for thermostat COOL and HEATCOOL modes.
# TYPE sdm_thermostat_thermostatTemperatureSetpoint_coolCelsius gauge
sdm_thermostat_thermostatTemperatureSetpoint_coolCelsius{CustomName="",Name="DU1AeyLeA",Room="Hallway",Type="sdm.devices.types.THERMOSTAT"} 0

# HELP sdm_thermostat_thermostatTemperatureSetpoint_heatCelsius Target temperature in Celsius for thermostat HEAT and HEATCOOL modes.
# TYPE sdm_thermostat_thermostatTemperatureSetpoint_heatCelsius gauge
sdm_thermostat_thermostatTemperatureSetpoint_heatCelsius{CustomName="",Name="DU1AeyLeA",Room="Hallway",Type="sdm.devices.types.THERMOSTAT"} 22
```
