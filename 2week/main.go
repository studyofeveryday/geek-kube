package main

import (
	"flag"
	"fmt"
	"log"
	"net"
	"net/http"
	"runtime"
	"strings"

	"github.com/golang/glog"
)

func Test(w http.ResponseWriter,r *http.Request)  {
	test := r.Header.Get("test")
	w.Header().Set("test",test)
	version := runtime.Version()
	w.Header().Set("version",version)
	w.WriteHeader(http.StatusOK)
	fmt.Fprintf(w, "Hello World")
	ip := getIp(r)
	glog.Infof("ip:%s",ip)
	glog.Infof("httpCode:%d",200)
}

func getIp(r *http.Request) string {
	xForwardedFor := r.Header.Get("X-Forwarded-For")
	ip := strings.TrimSpace(strings.Split(xForwardedFor, ",")[0])
	if ip != "" {
		return ip
	}

	ip = strings.TrimSpace(r.Header.Get("X-Real-Ip"))
	if ip != "" {
		return ip
	}

	if ip, _, err := net.SplitHostPort(strings.TrimSpace(r.RemoteAddr)); err == nil {
		if ip == "::1" {
			return "127.0.0.1"
		}
		return ip
	}
	return ""
}

func health(w http.ResponseWriter,r *http.Request)  {
	w.Write([]byte("200"))
}

func main()  {
	flag.Parse()
	defer glog.Flush()
	log.Println("server init")
	http.HandleFunc("/test",Test)
	http.HandleFunc("/health",health)

	http.ListenAndServe(":8080",nil)
}

