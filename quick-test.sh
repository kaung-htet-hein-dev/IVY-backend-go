#!/bin/bash

# Quick Docker Test Script for Both Dev and Production
# Run this after starting Docker Desktop

set -e

# Colors
GREEN='\033[0;32m'
RED='\033[0;31m'
YELLOW='\033[1;33m'
BLUE='\033[0;34m'
PURPLE='\033[0;35m'
CYAN='\033[0;36m'
NC='\033[0m'

echo -e "${BLUE}🐳 Quick Docker Test for IVY Backend (Dev & Prod)${NC}\n"

# Function to show usage
show_usage() {
    echo -e "${CYAN}Usage: $0 [option]${NC}"
    echo -e "${CYAN}Options:${NC}"
    echo -e "  dev     - Test development Dockerfile only"
    echo -e "  prod    - Test production Dockerfile only"
    echo -e "  both    - Test both Dockerfiles (default)"
    echo -e "  clean   - Clean up test containers and images"
    echo ""
}

# Function to test a specific Dockerfile
test_dockerfile() {
    local dockerfile=$1
    local tag_name=$2
    local env_name=$3
    local port=$4
    local container_name=$5
    
    echo -e "${PURPLE}===========================================${NC}"
    echo -e "${PURPLE}Testing $env_name Environment${NC}"
    echo -e "${PURPLE}===========================================${NC}\n"
    
    # Build the image
    echo -e "${BLUE}Building $env_name image...${NC}"
    echo "docker build -f $dockerfile -t $tag_name ."
    
    if docker build -f "$dockerfile" -t "$tag_name" .; then
        echo -e "${GREEN}✅ $env_name build successful${NC}\n"
    else
        echo -e "${RED}❌ $env_name build failed${NC}"
        return 1
    fi
    
    # Show image info
    echo -e "${BLUE}Image information:${NC}"
    docker images "$tag_name" --format "table {{.Repository}}\t{{.Tag}}\t{{.Size}}\t{{.CreatedSince}}"
    echo ""
    
    # Test run the container
    echo -e "${BLUE}Testing $env_name container startup...${NC}"
    
    if docker run -d --name "$container_name" -p "$port:8080" "$tag_name"; then
        echo -e "${GREEN}✅ $env_name container started${NC}\n"
        
        # Wait a moment
        echo -e "${BLUE}Waiting for application to start...${NC}"
        sleep 5
        
        # Check if container is still running
        if docker ps | grep -q "$container_name"; then
            echo -e "${GREEN}✅ $env_name container is running${NC}\n"
            
            # Show logs
            echo -e "${BLUE}Container logs (last 10 lines):${NC}"
            docker logs "$container_name" | tail -10
            echo ""
            
            # Test connection
            echo -e "${BLUE}Testing connection on port $port...${NC}"
            if curl -f -s "http://localhost:$port" >/dev/null 2>&1; then
                echo -e "${GREEN}✅ $env_name application is responding${NC}"
            elif curl -f -s "http://localhost:$port/health" >/dev/null 2>&1; then
                echo -e "${GREEN}✅ $env_name health endpoint is responding${NC}"
            else
                echo -e "${YELLOW}⚠️  $env_name application might not be responding yet${NC}"
                echo -e "${YELLOW}   This could be normal if database connection is required${NC}"
            fi
            
        else
            echo -e "${RED}❌ $env_name container exited unexpectedly${NC}"
            echo -e "${YELLOW}Container logs:${NC}"
            docker logs "$container_name"
            return 1
        fi
        
        # Ask if user wants to keep container running
        echo ""
        read -p "Keep $env_name container running for manual testing? (y/n): " -n 1 -r
        echo
        if [[ ! $REPLY =~ ^[Yy]$ ]]; then
            echo -e "${BLUE}Stopping $env_name container...${NC}"
            docker stop "$container_name" >/dev/null 2>&1 || true
            docker rm "$container_name" >/dev/null 2>&1 || true
            echo -e "${GREEN}✅ $env_name cleanup completed${NC}"
        else
            echo -e "${CYAN}$env_name container left running at http://localhost:$port${NC}"
            echo -e "${CYAN}To stop: docker stop $container_name && docker rm $container_name${NC}"
        fi
        echo ""
        
    else
        echo -e "${RED}❌ Failed to start $env_name container${NC}"
        return 1
    fi
}

# Function to clean up test resources
cleanup_tests() {
    echo -e "${BLUE}Cleaning up test containers and images...${NC}"
    
    # Stop and remove containers
    docker stop ivy-dev-test ivy-prod-test ivy-test 2>/dev/null || true
    docker rm ivy-dev-test ivy-prod-test ivy-test 2>/dev/null || true
    
    # Remove test images
    docker rmi ivy-backend:dev-test ivy-backend:prod-test ivy-backend-test 2>/dev/null || true
    
    echo -e "${GREEN}✅ Cleanup completed${NC}"
}

# Parse command line arguments
case "${1:-both}" in
    "dev")
        TEST_MODE="dev"
        ;;
    "prod")
        TEST_MODE="prod"
        ;;
    "both")
        TEST_MODE="both"
        ;;
    "clean")
        cleanup_tests
        exit 0
        ;;
    "-h"|"--help"|"help")
        show_usage
        exit 0
        ;;
    *)
        echo -e "${RED}❌ Unknown option: $1${NC}"
        show_usage
        exit 1
        ;;
esac

# Check if Docker is running
echo -e "${BLUE}1. Checking Docker status...${NC}"
if ! docker info >/dev/null 2>&1; then
    echo -e "${RED}❌ Docker is not running. Please start Docker Desktop first.${NC}"
    exit 1
fi
echo -e "${GREEN}✅ Docker is running${NC}\n"

# Initialize success tracking
DEV_SUCCESS=false
PROD_SUCCESS=false

# Test based on mode
if [ "$TEST_MODE" = "dev" ] || [ "$TEST_MODE" = "both" ]; then
    echo -e "${CYAN}Testing Development Dockerfile...${NC}\n"
    if [ -f "Dockerfile.dev" ]; then
        if test_dockerfile "Dockerfile.dev" "ivy-backend:dev-test" "Development" "8081" "ivy-dev-test"; then
            DEV_SUCCESS=true
        fi
    else
        echo -e "${RED}❌ Dockerfile.dev not found${NC}\n"
    fi
fi

if [ "$TEST_MODE" = "prod" ] || [ "$TEST_MODE" = "both" ]; then
    echo -e "${CYAN}Testing Production Dockerfile...${NC}\n"
    if [ -f "Dockerfile.prod" ]; then
        if test_dockerfile "Dockerfile.prod" "ivy-backend:prod-test" "Production" "8082" "ivy-prod-test"; then
            PROD_SUCCESS=true
        fi
    elif [ -f "Dockerfile.production" ]; then
        if test_dockerfile "Dockerfile.production" "ivy-backend:prod-test" "Production" "8082" "ivy-prod-test"; then
            PROD_SUCCESS=true
        fi
    else
        echo -e "${RED}❌ Dockerfile.prod or Dockerfile.production not found${NC}\n"
    fi
fi

# Final Summary
if [ "$TEST_MODE" = "both" ]; then
    echo -e "${PURPLE}===========================================${NC}"
    echo -e "${PURPLE}FINAL TEST SUMMARY${NC}"
    echo -e "${PURPLE}===========================================${NC}\n"

    if [ "$DEV_SUCCESS" = true ]; then
        echo -e "${GREEN}✅ Development Dockerfile: SUCCESS${NC}"
    else
        echo -e "${RED}❌ Development Dockerfile: FAILED${NC}"
    fi

    if [ "$PROD_SUCCESS" = true ]; then
        echo -e "${GREEN}✅ Production Dockerfile: SUCCESS${NC}"
    else
        echo -e "${RED}❌ Production Dockerfile: FAILED${NC}"
    fi
    echo ""
fi

echo ""
echo -e "${YELLOW}📋 Manual testing commands:${NC}"
echo "Development:"
echo "  Build:    docker build -f Dockerfile.dev -t ivy-backend:dev ."
echo "  Run:      docker run -d --name ivy-dev -p 8081:8080 ivy-backend:dev"
echo "  Test:     curl http://localhost:8081"
echo ""
echo "Production:"
echo "  Build:    docker build -f Dockerfile.prod -t ivy-backend:prod ."
echo "  Run:      docker run -d --name ivy-prod -p 8082:8080 ivy-backend:prod"
echo "  Test:     curl http://localhost:8082"
echo ""
echo "Cleanup:"
echo "  Stop all:     docker stop \$(docker ps -q --filter 'name=ivy-*')"
echo "  Remove all:   docker rm \$(docker ps -aq --filter 'name=ivy-*')"
echo "  Clean script: ./quick-test.sh clean"

echo ""
echo -e "${BLUE}🎉 Docker test completed!${NC}"
echo ""
echo -e "${CYAN}Usage examples:${NC}"
echo "  ./quick-test.sh         # Test both environments"
echo "  ./quick-test.sh dev     # Test development only"
echo "  ./quick-test.sh prod    # Test production only"
echo "  ./quick-test.sh clean   # Clean up test resources"
echo "  ./quick-test.sh help    # Show help"
