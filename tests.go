package main

import (
	"fmt"

	"github.com/RouteHub-Link/DomainUtils/validator"
)

func test() {
	_validator := validator.DefaultValidator()
	fmt.Println("Hello, World!")

	links := []string{
		"http://google.com",
		"https://www.google.com",
		"http://www.google.com",
		"http://www.google.com/",
		"www.google.com",
		"http://username:",
		"https://www.google.com/file.txt",
	}

	for _, link := range links {
		_, err := _validator.ValidateURL(link)

		if err != nil {
			fmt.Println(err)
		}

		_, err = _validator.ValidateSite(link)

		if err != nil {
			fmt.Println(err)
		}
	}
}

func dnsTest() {
	_validator := validator.DefaultValidator()
	_validator.ValidateOwnershipOverDNSTxtRecord("https://routehub.link", "routehub_domainkey", "e322c8a8ffef929ce17002ec521eeee2", "1.1.1.1:53")
}
