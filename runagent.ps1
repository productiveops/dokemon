# Clean
if (Test-Path .\agent.exe) {
    Remove-Item -Force .\agent.exe
}

# Build
go build .\cmd\agent\

# Run
$env:SERVER_URL="http://localhost:5173"
$env:TOKEN="uXQ/JvpmhtjiIm7MUdcUKC47/bIv7LyYDd3CLDvR0ixkUOYy4bIwog=="
$env:LOG_LEVEL="DEBUG"
.\agent.exe
