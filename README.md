## Structure

- main.go: entrypoint
- httpclient: main script
    - service: logic 
- cmd: cli 
- mock: auxiliary for mocking

## Assumptions
- We are going to send Get request 
- `parallel` default value is 10

## Install
In project directory
```shell
go install
```
---
#### Run
```shell
adjust https://google.com
```
OR
```shell
adjust -parallel=4 https://google.com
```
## Result
```shell
https://google.com 3399f485fbb02147b3f2fc8840f4e2fd 
```
## Test
```shell
go test ./...
```
