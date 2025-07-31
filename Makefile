build: 
	go build -o wallpeek

mac:
	GOOS=darwin GOARCH=amd64 go build -o wallpeek-mac

linux:
	GOOS=linux GOARCH=amd64 go build -o wallpeek-linux

windows:
	GOOS=windows GOARCH=amd64 go build -o wallpeek-windows.exe

clean:
	rm -f wallpeek wallpeek-mac wallpeek-linux wallpeek-windows.exe
