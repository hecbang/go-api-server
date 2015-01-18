package controllers

import (
	"models/testing"
)

type Testing struct {
}

//tp 支持两种取值 app, db
func (this *Testing) Concurrence(tp string) {
	if tp == "db" {
		testing.DatabaseConcurrence()
	} else {
		testing.AppConcurrence()
	}
}
