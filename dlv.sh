go build -o app.exe
dlv debug --headless --listen=:2345 --api-version=2 --accept-multiclient
dlv --listen=:2345 --headless=true --api-version=2 --accept-multiclient exec ./app.exe