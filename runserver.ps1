# Clean
if (Test-Path .\server.exe) {
    Remove-Item -Force .\server.exe
}

Set-Location web
npm run build
Set-Location ..

# Build
go build .\cmd\server\

# Run
$env:DB_CONNECTION_STRING="c:\temp\dokemondata\db"
$env:DATA_PATH="c:\temp\dokemondata"
$env:LOG_LEVEL="DEBUG"
$env:SSL_ENABLED="0"
.\server.exe
