package context

import (
	"encoding/xml"
	"net/http"
)

var xmlContentType = []string{"application/xml; charset=utf-8"}
var plainContentType = []string{"text/plain; charset=utf-8"}

//Render render from bytes
func (ctx *Context) Render(bytes []byte) {
	//debug
	//fmt.Println("response msg = ", string(bytes))
	if ctx.HttType == "fasthttp" {
		//debug
		//fmt.Println("response msg = ", string(bytes))
		ctx.FastHttpCtx.Response.SetStatusCode(200)
		_, err := ctx.FastHttpCtx.Write(bytes)
		if err != nil {
			panic(err)
		}
	}else{
		//debug
		//fmt.Println("response msg = ", string(bytes))
		ctx.Writer.WriteHeader(200)
		_, err := ctx.Writer.Write(bytes)
		if err != nil {
			panic(err)
		}
	}
}

//String render from string
func (ctx *Context) String(str string) {
	if ctx.HttType == "fasthttp" {
		if val := ctx.FastHttpCtx.Response.Header.Peek("Content-Type"); len(val) == 0 {
			ctx.FastHttpCtx.Response.Header.Set("Content-Type",plainContentType[0])
		}
	}else{
		writeContextType(ctx.Writer, plainContentType)
	}
	ctx.Render([]byte(str))
}

//XML render to xml
func (ctx *Context) XML(obj interface{}) {
	if ctx.HttType == "fasthttp" {
		if val := ctx.FastHttpCtx.Response.Header.Peek("Content-Type"); len(val) == 0 {
			ctx.FastHttpCtx.Response.Header.Set("Content-Type",xmlContentType[0])
		}
	}else {
		writeContextType(ctx.Writer, xmlContentType)
	}
	bytes, err := xml.Marshal(obj)
	if err != nil {
		panic(err)
	}
	ctx.Render(bytes)
}

func writeContextType(w http.ResponseWriter, value []string) {
	header := w.Header()
	if val := header["Content-Type"]; len(val) == 0 {
		header["Content-Type"] = value
	}
}
