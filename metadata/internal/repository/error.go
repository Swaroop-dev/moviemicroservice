package repository

import "errors"

//ErrornotFound is returned when requested record is not found
var ErrornotFound = errors.New("not found")
