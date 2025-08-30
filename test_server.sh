#!/bin/bash

echo "=== Room Booker Testing Script ==="
echo

# Function to check if port is free
check_port() {
    if lsof -i:8080 > /dev/null 2>&1; then
        echo "Port 8080 is busy, killing processes..."
        lsof -ti:8080 | xargs kill -9 2>/dev/null || true
        sleep 2
    fi
}

# Function to start server
start_server() {
    echo "Starting server..."
    cd /Users/takeso/projects/goproj/roombooker
    go run cmd/server/main.go &
    SERVER_PID=$!
    sleep 3
    echo "Server started with PID: $SERVER_PID"
}

# Function to test endpoints
test_endpoints() {
    echo "Testing endpoints..."

    # Health check
    echo "1. Health check:"
    curl -s http://localhost:8080/health
    echo -e "\n"

    # Login page
    echo "2. Login page (first 3 lines):"
    curl -s http://localhost:8080/login | head -3
    echo -e "\n"

    # Main page
    echo "3. Main page (first 3 lines):"
    curl -s http://localhost:8080/ | head -3
    echo -e "\n"

    # CSS file
    echo "4. CSS file (first 2 lines):"
    curl -s http://localhost:8080/static/css/style.css | head -2
    echo -e "\n"

    # JS file
    echo "5. JS file (first 2 lines):"
    curl -s http://localhost:8080/static/js/app.js | head -2
    echo -e "\n"

    # API without auth (should fail)
    echo "6. API without auth (should return 401):"
    curl -s http://localhost:8080/api/offices
    echo -e "\n"
}

# Function to stop server
stop_server() {
    echo "Stopping server..."
    pkill -f "go run cmd/server/main.go"
    sleep 1
}

# Main execution
check_port
start_server
test_endpoints
stop_server

echo "=== Testing completed ==="
