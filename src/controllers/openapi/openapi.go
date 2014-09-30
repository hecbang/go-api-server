package openapi

import (
	"fmt"
	"io/ioutil"
	"net/http"
)

type OpenApi struct {
	Data []byte
}

//设置请求数据到Data属性中，便于数据后续被具体的业务逻辑处理
func (this *OpenApi) SetRequestData(w http.ResponseWriter, r *http.Request) {
	if err := r.ParseForm(); err != nil {
		panic(err)
	}
	result, err := ioutil.ReadAll(r.Body)
	defer r.Body.Close()
	if err != nil {
		panic(err)
	}
	//make []byte to string type
	this.Data = result

}

//开放API请求入口点
func (this *OpenApi) Index(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintln(w, string(this.Data))
}
