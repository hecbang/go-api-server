package common

import (
	"github.com/jsonq"
)

type JsonQuery struct {
	Jq *jsonq.JsonQuery
}

func NewJsonQuery(filename string) *JsonQuery {
	var data map[string]interface{}
	LoadJson(filename, &data)
	return &JsonQuery{Jq: jsonq.NewQuery(data)}
}

func (this *JsonQuery) String(s ...string) string {
	retval, err := this.Jq.String(s...)
	if err != nil {
		panic(err.Error())
	}
	return retval
}

func (this *JsonQuery) Int(s ...string) int {
	retval, err := this.Jq.Int(s...)
	if err != nil {
		panic(err.Error())
	}
	return retval
}

func (this *JsonQuery) Bool(s ...string) bool {
	retval, err := this.Jq.Bool(s...)
	if err != nil {
		panic(err.Error())
	}
	return retval
}

func (this *JsonQuery) Float(s ...string) float64 {
	retval, err := this.Jq.Float(s...)
	if err != nil {
		panic(err.Error())
	}
	return retval
}
