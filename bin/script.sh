#!/bin/bash

# Cameron Tillett
# University of Belize
# 22/ 09/ 2022
# System Administration Quiz 1

# A. Pass two Arguments
echo "Hello! Please type two arguments:"
read topleveldirectory identifier

# B. C. D. Prompt Continue. Response and triggers 
echo "I am about to create a directory structure named $topleveldirectory"
read -p "Do you want to continue? (Yes/No)" CONT

# If loop that checks the response
if [[ "$CONT" = "y" || "$CONT" = "Y" || "$CONT" = "yes" || "$CONT" = "Yes" ]]; then
	echo "Creating directory structure..."
	mkdir "$topleveldirectory"
	mkdir -p "$topleveldirectory/cmd/api"
	touch "$topleveldirectory/cmd/api/main.go"
	mkdir -p "$topleveldirectory/internals"
	mkdir -p "$topleveldirectory/migrations"
	mkdir -p "$topleveldirectory/remote"
	touch  "$topleveldirectory/go.mod"
	touch "$topleveldirectory/Makefile"

	go mod init "$topleveldirectory.$identifier"
	
	echo "//File: cmd/api/main.go

	package main

	import "fmt"

	func main() {
	fmt.Println("Hello World!")}" > "$topleveldirectory/cmd/api/main.go"

	echo "I have created a *main.go* file for you to test the directory structure."
	echo "Type *go run ./cmd/api* at the root directory of your project to test your project."
	echo "Thank you. :^)"

elif [[ "$CONT" = "n" || "$CONT" = "N" || "$CONT" = "No" || "$CONT" = "no" ]]; then
	echo "Its cool, I understand."
	echo "Abort."
fi
