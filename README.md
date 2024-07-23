# handheld-database-gui
A graphical interface made with SDL for low-end devices

# How to install GO?

1. sudo apt-get update
2. wget https://go.dev/dl/go1.21.0.linux-arm64.tar.gz
3. sudo tar -xvf go1.21.0.linux-arm64.tar.gz
4. sudo mv go /usr/local
5. export GOROOT=/usr/local/go
6. export GOPATH=$HOME/go
7. export PATH=$GOPATH/bin:$GOROOT/bin:$PATH
8. source ~/.profile

# How to run?

```
go run main.go
```