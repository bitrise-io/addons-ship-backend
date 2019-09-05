package main

import (
	"fmt"

	"github.com/pkg/errors"
	"github.com/simonmarton/common-colors/processimage"
)

func main() {
	url := "https://concrete-userfiles-production.s3.us-west-2.amazonaws.com/repositories/7a0c4ec18ba27f55/avatar/avatar.jpg"
	_, err := processimage.FromURL(url)
	fmt.Println(errors.WithStack(err))
}
