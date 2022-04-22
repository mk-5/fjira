mkdir -p build
go build -o build/fjira cmd/fjira-cli/main.go
cp build/fjira .
chmod +x fjira
