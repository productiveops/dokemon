# Clean
rm -f ./server
cd web
npm run build
cd ..

# Build
go build ./cmd/server

# Run 
export DB_CONNECTION_STRING="/tmp/db"
export DATA_PATH="/tmp"
export LOG_LEVEL="DEBUG"
./server
