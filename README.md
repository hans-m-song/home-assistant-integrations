# home-assistant-integrations

## Integrations

- Huawei HG659 Router
- Zever Solar Evershine TLC5000 Solar Inverter

## Deployment

### Docker

```bash
docker build -t home-assistant-integrations .
docker run -d home-assistant-integrations
```

### Kubernetes

```bash
helm add repo home-assistant-integrations https://home-assistant-integrations.charts.axatol.xyz
helm repo update
helm upgrade \
  home-assistant-integrations \
  home-assistant-integrations/home-assistant-integrations
  --install \
  --create-namespace \
  --atomic \
  --namespace home-assistant \
  --set providers.huaweiHg659.enabled=true \
  --set providers.huaweiHg659.address=http://192.168.1.1 \
  --set providers.huaweiHg659.enabled=true \
  --set providers.huaweiHg659.address=http://192.168.1.44

## Configuration

### Server

- `LOG_LEVEL` Zerolog log level (default: `info`)
- `LOG_FORMAT` Zerolog log format (default: `json`)
- `LISTEN_PORT` Server listen port (default: `8080`)
- `MQTT_URI` Broker address (required)

### Huawei HG659

- `HUAWEI_HG659_ENABLED` (default: `false`)
- `HUAWEI_HG659_ADDRESS` (required if enabled)
- `HUAWEI_HG659_ENTITY_NAME` (default: `huawei_hg659`)
- `HUAWEI_HG659_POLL_RATE` (default: `60s`)

### ZeverSolar TLC5000

- `ZEVER_SOLAR_TLC5000_ENABLED` (default: `false`)
- `ZEVER_SOLAR_TLC5000_ADDRESS` (required if enabled)
- `ZEVER_SOLAR_TLC5000_ENTITY_NAME` (default: `zever_solar_tlc5000`)
- `ZEVER_SOLAR_TLC5000_POLL_RATE` (default: `30s`)
```
