# Clean
rm -f ./agent

# Build
go build ./cmd/agent

# Run 
export SERVER_URL="http://192.168.1.7:9090"
export TOKEN="1ZRXtM8ZCMVjP2jIXV99GYLeWfXgN4j0okzK0NaRpjosuIYU5wVoCEA=="
export LOG_LEVEL="DEBUG"
./agent
