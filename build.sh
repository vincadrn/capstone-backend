#!/bin/bash

go mod tidy
go build -o bin/application application.go