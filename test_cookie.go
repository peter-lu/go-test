package main

import (
	"fmt"
	"log"
	"net/http"
	"time"
)

func main() {
	http.HandleFunc("/register", register)
	http.HandleFunc("/getCookie", getCookie)
	log.Fatal(http.ListenAndServe("localhost:8888", nil))
}

func register(w http.ResponseWriter, r *http.Request) {
	cName := http.Cookie{
		Name:       "name",      //  cookie的名称
		Value:      "golang",    //  cookie名称对应的值
		Path:       "/",         //
		Domain:     "",          //  cookie的作用域
		Expires:    time.Time{}, //  cookie的过期时间
		RawExpires: "",
		MaxAge:     0,     //  设置过期时间，对应浏览器cookie的MaxAge属性
		Secure:     false, //  设置Secure属性(说明：Cookie的Secure属性，意味着保持Cookie通信只限于加密传输，指示浏览器仅仅在通过安全/加密连接才能使用该Cookie。如果一个Web服务器从一个非安全连接里设置了一个带有secure属性的Cookie，当Cookie被发送到客户端时，它仍然能通过中间人攻击来拦截)
		HttpOnly:   true,  //  设置httpOnly属性（说明：Cookie的HttpOnly属性，指示浏览器不要在除HTTP（和 HTTPS)请求之外暴露Cookie。一个有HttpOnly属性的Cookie，不能通过非HTTP方式来访问，例如通过调用JavaScript(例如，引用 document.cookie），因此，不可能通过跨域脚本（一种非常普通的攻击技术）来偷走这种Cookie。尤其是Facebook 和 Google 正在广泛地使用HttpOnly属性。）
		SameSite:   0,
		Raw:        "",
		Unparsed:   nil,
	}
	cId := http.Cookie{
		Name:       "id",
		Value:      "21",
		Path:       "/",
		Domain:     "",
		Expires:    time.Time{},
		RawExpires: "",
		MaxAge:     0,
		Secure:     false,
		HttpOnly:   true,
		SameSite:   0,
		Raw:        "",
		Unparsed:   nil,
	}
	// 设置cookie
	w.Header().Set("Set-Cookie", cId.String())
	w.Header().Add("Set-Cookie", cName.String())
}

func getCookie(w http.ResponseWriter, r *http.Request) {
	// 读取cookie
	cookie := r.Header["Cookie"]
	fmt.Println(cookie)
	// 通过key获取cookie
	id, err := r.Cookie("id")
	if err != nil {
		fmt.Println(err.Error())
	} else {
		fmt.Println(id)
	}
	// 获取全部cookie
	cookies := r.Cookies()
	fmt.Println(cookies)
	return
}
