export GOPROXY=https://goproxy.io,direct
export GOARCH=amd64
export GOOS=linux

cd   ../src
go build -o ../bin/gameserver -ldflags "-s -w" main.go
read -n1 -p "Press any key to continue..."