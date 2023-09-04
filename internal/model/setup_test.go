package model

import "regexp"

var validUUID = regexp.MustCompile(`\S{8}-\S{4}-\S{4}-\S{4}-\S{12}`)
