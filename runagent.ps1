# Clean
if (Test-Path .\agent.exe) {
    Remove-Item -Force .\agent.exe
}

# Build
go build .\cmd\agent\

# Run
$env:SERVER_URL="http://localhost:5173"
$env:LOG_LEVEL="DEBUG"
.\agent.exe
