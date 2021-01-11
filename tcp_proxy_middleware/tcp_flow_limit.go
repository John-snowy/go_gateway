package tcp_proxy_middleware

import (
	"fmt"
	"github.com/go_gateway/dao"
	"github.com/go_gateway/public"
	"strings"
)

func TCPFlowLimitMiddleware() func(c *TcpSliceRouterContext) {
	return func(c *TcpSliceRouterContext) {
		serviceInterface := c.Get("service")
		if serviceInterface == nil {
			c.conn.Write([]byte("get service empty"))
			c.Abort()
			return
		}
		serviceDetail := serviceInterface.(*dao.ServiceDetail)
		if serviceDetail.AccessControl.ServiceFlowLimit > 0 {
			serviceLimiter, err := public.FlowLimiterHandler.GetLimiter(public.FlowServicePrefix+serviceDetail.Info.ServiceName, float64(serviceDetail.AccessControl.ServiceFlowLimit))
			if err != nil {
				c.conn.Write([]byte(err.Error()))
				c.Abort()
				return
			}
			if !serviceLimiter.Allow() {
				c.conn.Write([]byte(fmt.Sprintf("service flow limit %v", float64(serviceDetail.AccessControl.ServiceFlowLimit))))
				c.Abort()
				return

			}
		}

		splits := strings.Split(c.conn.RemoteAddr().String(), ":") // ip:port
		clientIP := ""
		if len(splits) == 2 {
			clientIP = splits[0]
		}
		if serviceDetail.AccessControl.ClientIPFlowLimit > 0 {
			clientLimiter, err := public.FlowLimiterHandler.GetLimiter(public.FlowServicePrefix+serviceDetail.Info.ServiceName+"_"+clientIP, float64(serviceDetail.AccessControl.ClientIPFlowLimit))
			if err != nil {
				c.conn.Write([]byte(err.Error()))
				c.Abort()
				return
			}
			if !clientLimiter.Allow() {
				c.conn.Write([]byte(fmt.Sprintf("%v flow limit %v", clientIP, float64(serviceDetail.AccessControl.ServiceFlowLimit))))
				c.Abort()
				return

			}
		}
		c.Next()

	}
}
