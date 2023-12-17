# Clean
rm -f ./agent

# Build
go build ./cmd/agent

# Run 
export SERVER_URL="http://192.168.1.7:9090"
export TOKEN="kUmRk3SvRKSRuuov0L4T9h8dgOMBnmvXIoLAjDON6rVATLlE9dd9pQ=="
export LOG_LEVEL="DEBUG"
./agent
