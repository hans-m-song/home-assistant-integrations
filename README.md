# zeversolar-monitor

Pulls data from a Zever Solar Evershine TLC5000 and submits to Home Assistant

## deployment

### Docker

```bash
docker build -t zeversolar-monitor .
docker run -d zeversolar-monitor
```

### Kubernetes

Refer to the manifests folder for a kubernetes deployment.

## configuration

Environment variables:

- `SOURCE_ENDPOINT`: inverter location to pull data from (default: `http://localhost/home.cgi`)
- `MQTT_TOPIC`: used to formulate topic to announce entities, i.e. (default: `zeversolar`)
- `MQTT_USER`: mqtt user
- `MQTT_PASS`: mqtt pass
- `SERVER_PORT`: server stub listening port (default: `8000`)
- `POLL_RATE`: delay between poll attempts in milliseconds (default: `5000`)
