# home-assistant-integrations

## Integrations

- Zever Solar Evershine TLC5000 Solar Inverter
- Huawei HG659 Router

## Deployment

### Docker

```bash
docker build -t home-assistant-integrations .
docker run -d home-assistant-integrations
```

### Kubernetes

Refer to the manifests folder for a kubernetes deployment.

## Configuration

Environment variables:

- `ZEVERSOLAR_TLC5000_ENDPOINT`: inverter network address to pull data from (default: `http://192.168.1.44/home.cgi`), leave empty to disable
- `ZEVERSOLAR_POLL_RATE`: delay between poll attempts in milliseconds (default: `5000`)
- `HUAWEI_HG659_ENDPOINT`: router admin portal address to pull data from (most likely `http://192.168.1.1/home.cgi`), leave empty to disable
- `HUAWEI_HG659_POLL_RATE`: delay between poll attempts in milliseconds (default: `5000`)
- `MQTT_ENDPOINT`: mqtt address to submit data to (**required**)
- `MQTT_USER`: mqtt instance username
- `MQTT_PASS`: mqtt instance password
- `MQTT_ANNOUNCE_RATE`: interval to announce integrations to mqtt broker (default: `60000`)
- `HTTP_PORT`: server stub listening port (default: `8000`, max: `65536`)
