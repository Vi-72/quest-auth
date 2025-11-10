# Quest Auth - Deployment Guide

## 🐳 Docker Deployment

### Build Docker Image
```bash
docker build -t quest-auth:latest .
```

### Run with Docker
```bash
docker run -d \
  --name quest-auth \
  -p 8080:8080 \
  -p 9090:9090 \
  -e DB_HOST=postgres \
  -e DB_PASSWORD=secret \
  -e JWT_SECRET_KEY=your-secret \
  quest-auth:latest
```

---

## 🐙 Docker Compose

### Start Services
```bash
docker compose up -d
```

### Stop Services
```bash
docker compose down
```

### View Logs
```bash
docker compose logs -f quest-auth
```

---

## ☸️ Kubernetes (Basic)

### Deployment YAML
```yaml
apiVersion: apps/v1
kind: Deployment
metadata:
  name: quest-auth
spec:
  replicas: 3
  selector:
    matchLabels:
      app: quest-auth
  template:
    metadata:
      labels:
        app: quest-auth
    spec:
      containers:
      - name: quest-auth
        image: quest-auth:latest
        ports:
        - containerPort: 8080
        - containerPort: 9090
        env:
        - name: DB_HOST
          value: postgres-service
        - name: JWT_SECRET_KEY
          valueFrom:
            secretKeyRef:
              name: quest-auth-secrets
              key: jwt-secret
```

### Service YAML
```yaml
apiVersion: v1
kind: Service
metadata:
  name: quest-auth-service
spec:
  selector:
    app: quest-auth
  ports:
  - name: http
    port: 8080
    targetPort: 8080
  - name: grpc
    port: 9090
    targetPort: 9090
  type: LoadBalancer
```

---

## 🔍 Health Checks

### HTTP Health Check
```bash
curl http://localhost:8080/health
```

### Docker Health Check
```dockerfile
HEALTHCHECK --interval=30s --timeout=3s \
  CMD curl -f http://localhost:8080/health || exit 1
```

---

## 🗄️ Database Migrations

### Run Migrations
Migrations run automatically on application start.

### Manual Migration
```bash
go run ./cmd/app
```

---

## 📊 Monitoring

### Logs
```bash
# Docker
docker logs -f quest-auth

# Docker Compose
docker compose logs -f quest-auth

# File
tail -f app.log
```

---

## 🔐 Production Checklist

- [ ] Change JWT_SECRET_KEY to strong random value
- [ ] Enable DB SSL mode (DB_SSLMODE=require)
- [ ] Use HTTPS for HTTP endpoints
- [ ] Set up database backups
- [ ] Configure log aggregation
- [ ] Set up monitoring and alerts
- [ ] Review security settings

---

**Last Updated:** November 10, 2025
