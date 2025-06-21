# Manual Docker Testing Commands

## Step 1: Start Docker

First, make sure Docker Desktop is running on your Mac.

## Step 2: Basic Build Test

```bash
cd /Volumes/SanDisk/developer/IVY/IVY-backend

# Test basic build
docker build -t ivy-test .
```

## Step 3: Check Build Result

```bash
# Check if image was created
docker images ivy-test

# Inspect the image
docker inspect ivy-test
```

## Step 4: Test Container Run (Simple)

```bash
# Try to run the container
docker run --name ivy-test-container ivy-test

# Or run in background
docker run -d --name ivy-test-container -p 8080:8080 ivy-test

# Check if it's running
docker ps

# Check logs
docker logs ivy-test-container
```

## Step 5: Test with Environment Variables

```bash
# Create test environment file
cat > .env.test << EOF
APP_PORT=8080
JWT_SECRET=test-secret
DB_HOST=localhost
DB_PORT=5432
DB_USER=testuser
DB_PASSWORD=testpass
DB_NAME=testdb
EOF

# Run with environment
docker run -d \
  --name ivy-env-test \
  -p 8081:8080 \
  --env-file .env.test \
  ivy-test

# Check logs
docker logs ivy-env-test
```

## Step 6: Test Application Response

```bash
# Test if application responds
curl http://localhost:8080
curl http://localhost:8080/health  # if you have health endpoint

# Or using wget
wget -qO- http://localhost:8080
```

## Step 7: Debug Inside Container

```bash
# Get shell access (if available)
docker exec -it ivy-test-container /bin/sh

# Or check what's running
docker exec ivy-test-container ps aux
```

## Step 8: Cleanup

```bash
# Stop containers
docker stop ivy-test-container ivy-env-test
docker rm ivy-test-container ivy-env-test

# Remove test image
docker rmi ivy-test
```

## Common Issues & Solutions

### Build Fails

- **Issue**: Go modules not downloading
- **Solution**: Add `RUN go mod download` step

### Container Exits Immediately

- **Issue**: Missing command arguments
- **Solution**: Your app needs `CMD ["./main", "prod"]`

### Permission Denied

- **Issue**: Binary not executable
- **Solution**: Add `RUN chmod +x main` after copy

### Can't Connect to App

- **Issue**: App not listening on correct interface
- **Solution**: Ensure app listens on `0.0.0.0:8080` not `localhost:8080`
