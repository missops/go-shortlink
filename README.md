# go-shortlink

git clone https://github.com/hiningmeng/go-shortlink.git  shortlink


// 创建短链接口
1.  http://127.0.0.1:8000/api/shorten 

{
	"url": "http://www.hiningmeng.cn",
	"expiration_in_minutes": 30
}

//短链详细信息
2. http://127.0.0.1:8000/api/info?shortlink=1 


//短链跳转真实长链地址
3. http://127.0.0.1:8000/1
