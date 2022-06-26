# line-webhook-receiver

## Setup

1. database
```bash
cp deployments/.env.example deployments/.env
cd deployments
docker-compose up -d
```

2. Setup configs
```bash
cp configs/config.yaml.example configs/config.yaml
```

3. Start service
```bash
go run ./
```

## Usage
1. Send message to line
```bash
curl -X POST {service_url}/messages/send --data "message=test message"
```

2. Query message list of the user
```bash
{service_url}/messages/{userLineId}
```