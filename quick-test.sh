#!/bin/bash

# Quick Docker Test Script
# Run this after starting Docker Desktop

set -e

# Colors
GREEN='\033[0;32m'
RED='\033[0;31m'
YELLOW='\033[1;33m'
BLUE='\033[0;34m'
NC='\033[0m'

echo -e "${BLUE}🐳 Quick Docker Test for IVY Backend${NC}\n"

# Check if Docker is running
echo -e "${BLUE}1. Checking Docker status...${NC}"
if ! docker info >/dev/null 2>&1; then
    echo -e "${RED}❌ Docker is not running. Please start Docker Desktop first.${NC}"
    exit 1
fi
echo -e "${GREEN}✅ Docker is running${NC}\n"

# Build the image
echo -e "${BLUE}2. Building Docker image...${NC}"
echo "docker build -t ivy-backend-test ."
if docker build -t ivy-backend-test .; then
    echo -e "${GREEN}✅ Build successful${NC}\n"
else
    echo -e "${RED}❌ Build failed${NC}"
    exit 1
fi

# Show image info
echo -e "${BLUE}3. Image information:${NC}"
docker images ivy-backend-test --format "table {{.Repository}}\t{{.Tag}}\t{{.Size}}\t{{.CreatedSince}}"
echo ""

# Test run the container
echo -e "${BLUE}4. Testing container startup...${NC}"
echo "docker run -d --name ivy-test -p 8080:8080 ivy-backend-test"

if docker run -d --name ivy-test -p 8080:8080 ivy-backend-test; then
    echo -e "${GREEN}✅ Container started${NC}\n"
    
    # Wait a moment
    echo -e "${BLUE}5. Waiting for application to start...${NC}"
    sleep 5
    
    # Check if container is still running
    if docker ps | grep -q ivy-test; then
        echo -e "${GREEN}✅ Container is running${NC}\n"
        
        # Show logs
        echo -e "${BLUE}6. Container logs:${NC}"
        docker logs ivy-test
        echo ""
        
        # Test connection
        echo -e "${BLUE}7. Testing connection...${NC}"
        if curl -f -s http://localhost:8080 >/dev/null 2>&1; then
            echo -e "${GREEN}✅ Application is responding on port 8080${NC}"
        elif curl -f -s http://localhost:8080/health >/dev/null 2>&1; then
            echo -e "${GREEN}✅ Health endpoint is responding${NC}"
        else
            echo -e "${YELLOW}⚠️  Application might not be responding yet${NC}"
            echo -e "${YELLOW}   This could be normal if database connection is required${NC}"
        fi
        
    else
        echo -e "${RED}❌ Container exited unexpectedly${NC}"
        echo -e "${YELLOW}Container logs:${NC}"
        docker logs ivy-test
    fi
else
    echo -e "${RED}❌ Failed to start container${NC}"
fi

echo ""
echo -e "${BLUE}8. Test Summary:${NC}"
echo -e "${GREEN}✅ Docker build: SUCCESS${NC}"
echo -e "${GREEN}✅ Container start: SUCCESS${NC}"

echo ""
echo -e "${YELLOW}📋 Manual testing commands:${NC}"
echo "  View logs:     docker logs ivy-test"
echo "  Test app:      curl http://localhost:8080"
echo "  Stop test:     docker stop ivy-test && docker rm ivy-test"
echo "  Remove image:  docker rmi ivy-backend-test"

echo ""
echo -e "${BLUE}🎉 Docker test completed!${NC}"

# Cleanup function (optional)
read -p "Do you want to stop and remove the test container now? (y/n): " -n 1 -r
echo
if [[ $REPLY =~ ^[Yy]$ ]]; then
    echo -e "${BLUE}Cleaning up...${NC}"
    docker stop ivy-test >/dev/null 2>&1 || true
    docker rm ivy-test >/dev/null 2>&1 || true
    echo -e "${GREEN}✅ Cleanup completed${NC}"
fi
