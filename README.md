# golang_helloworld
Hello World in Golang

### Environmental Variables
```
  HELLOWORLD_DISPLAYEXTERNALIP : "False"
  HELLOWORLD_DISPLAYGEOLOCATION : "False"
  HELLOWORLD_CRASHAPP : "False"
  HELLOWORLD_CRASHAPPCOUNT : "5"
  HELLOWORLD_DEBUG : "False"
  HELLOWORLD_SIMULATEREADY : "False"
  HELLOWORLD_WAITBEFOREREADY : "30"
  HELLOWORLD_WAITBEFOREREADY : "80"
```

## Cloud Foundry 

```
appname=golang-hello-world
cf push $appname
```

## Docker

`sh
docker build .
`

## Build 
`sh
go build -o hello .
`
