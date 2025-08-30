#!/bin/bash
echo "Starting server..."
go run ./cmd/server &
SERVER_PID=$!
sleep 3

echo "Testing health endpoint..."
curl -s http://localhost:8080/health

echo -e "\nTesting main page..."
curl -s http://localhost:8080/ | grep -o '<title>.*</title>'

echo -e "\nTesting test page..."
curl -s http://localhost:8080/test | grep -o '<title>.*</title>'

echo -e "\nTesting OIDC start (should redirect)..."
curl -s -I http://localhost:8080/auth/oidc/start | grep -E "(HTTP|Location)"

echo -e "\nKilling server..."
kill $SERVER_PID
wait $SERVER_PID 2>/dev/null
echo "Done."
