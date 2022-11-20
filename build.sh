#!/usr/bin/env bash

GOOS=linux go build -o ExchangeRate main.go

zip -r ExchangeRate.zip .

# Handle value shoudl be ExchangeRate