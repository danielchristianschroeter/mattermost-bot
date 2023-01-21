package utils

import (
	"fmt"
	"net/url"

	"github.com/gearnode/privatebin"
)

func PrivateBinPaste(msg string) (string, error) {
	if GetConfigBoolValue("PRIVATEBIN_ENABLE") {
		uri, err := url.Parse(GetConfigValue("PRIVATEBIN_HOST"))
		if err != nil {
			fmt.Printf("cannot parse host in PRIVATEBIN_HOST")
		}
		client := privatebin.NewClient(uri, "", "")

		resp, err := client.CreatePaste(
			msg,
			GetConfigValue("PRIVATEBIN_EXPIRE"),
			GetConfigValue("PRIVATEBIN_FORMATTER"),
			GetConfigBoolValue("PRIVATEBIN_OPENDISCUSSION"),
			GetConfigBoolValue("PRIVATEBIN_BURNAFTERREADING"),
			GetConfigValue("PRIVATEBIN_PASSWORD"))
		if err != nil {
			fmt.Printf("cannot create the paste: %v", err)
			return "", err
		}
		fmt.Printf("%s\n", resp.URL)
		return resp.URL, nil
	}
	return "", nil
}
