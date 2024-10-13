package middleware

//JWT和Session混用  JWT是暴露在浏览器的，虽然加密但是也是安全的，不宜存储敏感信息。
//先查找有无session 有的话从session读取信息，没有的话从authorization中读取token，创造新的session

