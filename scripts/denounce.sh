#!/bin/bash
set -eo pipefail

# brew install hivemq/mqtt-cli/mqtt-cli

topics=(
  homeassistant/binary_sensor/zever_solar_tlc5000/solar_status/config
  homeassistant/binary_sensor/huawei_hg659/router_internet_connected/config
  homeassistant/sensor/zever_solar_tlc5000/solar_energy_today_kwh/config
  homeassistant/sensor/zever_solar_tlc5000/solar_power_ac_w/config
  homeassistant/sensor/zever_solar_tlc5000/solar_last_updated/config
  homeassistant/sensor/huawei_hg659/router_device_uptime/config
  homeassistant/sensor/huawei_hg659/router_internet_uptime/config
  homeassistant/sensor/huawei_hg659/router_internet_err_reason/config
  homeassistant/sensor/huawei_hg659/router_internet_connection_status/config
  homeassistant/sensor/huawei_hg659/router_internet_self_test_message/config
)

for topic in ${topics[@]}; do
  mqtt publish --host=localhost --topic=$topic --message='' --retain
done
