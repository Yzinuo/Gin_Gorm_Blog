package utils

import (
	"errors"
	"log/slog"
	"net"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/lionsoul2014/ip2region/binding/golang/xdb"
	"xojoc.pw/useragent"
)

type IPutils struct{}

// 方便其他package调用方法
var IP  = new(IPutils)

func (*IPutils) GetIpaddress(c *gin.Context) (IPstring string){
	// 项目明确了用ningx转发，所以真正的客户端的IP被封存在http请求`X-Real-Ip`中
	IPstring = c.Request.Header.Get("X-Real-Ip")
	//如果失败，尝试从`X-Forwarded-For`中获取
	if IPstring == "" || len(IPstring) == 0 || strings.EqualFold("unknown", IPstring){
		IPs := c.Request.Header.Get("X-Forwarded-For")
		// X-Forwarded-For 有三个IP地址， 第一个是客户端的真实IP，第二个是代理服务器的IP，第三个是负载均衡器的IP
		splitIPs := strings.Split(IPs,",")
		if len(splitIPs)> 0{
			IPstring = splitIPs[0]
		}
	}

	// 如果前面两种方法都失败，直接使用remoteAddr获取地址 但是因为通过转发 这种方法获取的是nginx的地址
	//至少保证有一个ip地址，不为空
	if IPstring == "" || len(IPstring) == 0 || strings.EqualFold("unknown", IPstring){
		IPstring = c.Request.RemoteAddr
	}

	//检查是否是在本机
	if strings.HasPrefix(IPstring,"127.0.0.1") || strings.HasPrefix(IPstring,"[::1]"){
		ip, err := externalIP()
		if err != nil {
			slog.Error("GetIpAddress, externalIP, err: ", err)
		}
		IPstring = ip.String()
	}

	// 检查是否包含多个IP，取第一个IP
	if IPstring != ""&&len(IPstring) > 15 && strings.Index(IPstring,",") > 0 {
		IPstring = strings.Split(IPstring,",")[0]
	}

	return IPstring
}

// 获取地域信息: 中国|0|江苏省|苏州市|电信

// vectorIndex 减少一次硬盘的IO
var Vindex []byte

func (*IPutils) GetIPsource(IPaddress string) string{
	 Dbpath := "../assets/ip2region.xdb"
	 // 查询IP是基于文件查询

	 if Vindex == nil{
			var err error
			// 注意阴影变量的错误
			Vindex,err = xdb.LoadVectorIndexFromFile(Dbpath)
			if err != nil {
				slog.Error("加载IP2region的Vindex失败",err)
				return ""
			}
	 }

	 Seacher,err := xdb.NewWithVectorIndex(Dbpath,Vindex)
	 if err != nil {
		slog.Error("加载IP2region的Seacher失败",err)
		return ""
	 }
	 defer Seacher.Close()

	 region,err := Seacher.SearchByStr(IPaddress)
	 if err!= nil {
		slog.Error("查询IP2region失败",err)
		return ""
	 }
	 return region
}

// 获取region的简单信息：湖南长沙 电信
func (i *IPutils) GetIPsourceSimpleInfo(IPaddress string) string{
	regionstring := i.GetIPsource(IPaddress)
// 国家|区域|省份|城市|ISP
// 只有中国的数据绝大部分精确到了城市, 其他国家部分数据只能定位到国家, 后面的选项全部是 0
	region := strings.Split(regionstring,"|")
	if region[0] != "中国" && region[0] != "0"{
		return region[0]
	}

	if region[2] == "0" {
		region[2] = ""
	}
	if region[3] == "0" {
		region[2] = ""
	}
	if region[4] == "0" {
		region[2] = ""
	}

	if region[2] == "" && region[3] == "" && region[4] == "" {
		return region[0]
	}
	return region[2] + region[3] +" "+ region[4]
}

func (*IPutils) GetUserAgent(c *gin.Context) *useragent.UserAgent {
	return useragent.Parse(c.Request.UserAgent())
}

// 获取非 127.0.0.1 的局域网 IP
func externalIP() (net.IP, error) {
	// 获取服务器的网络接口列表
	ifaces, err := net.Interfaces()
	if err != nil {
		return nil, err
	}
	for _, iface := range ifaces {
		// 不在活动状态
		if iface.Flags&net.FlagUp == 0 {
			continue
		}
		// 环回
		if iface.Flags&net.FlagLoopback != 0 {
			continue
		}
		// 单播接口地址列表
		addrs, err := iface.Addrs()
		if err != nil {
			return nil, err
		}
		for _, addr := range addrs {
			ip := getIpFromAddr(addr)
			if ip == nil {
				continue
			}
			return ip, nil
		}
	}
	return nil, errors.New("connected to the network")
}

func getIpFromAddr(addr net.Addr) net.IP {
	var ip net.IP
	switch v := addr.(type) {
	case *net.IPNet:
		ip = v.IP
	case *net.IPAddr:
		ip = v.IP
	}
	if ip == nil || ip.IsLoopback() {
		return nil
	}
	ip = ip.To4()
	if ip == nil {
		return nil
	}
	return ip
}
