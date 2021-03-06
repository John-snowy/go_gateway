package grpc_proxy_middleware

import (
	"errors"
	"fmt"
	"github.com/go_gateway/dao"
	"github.com/go_gateway/public"
	"google.golang.org/grpc"
	"google.golang.org/grpc/peer"
	"log"
	"strings"
)

func GrpcWhiteListMiddleware(serviceDetail *dao.ServiceDetail) func(srv interface{}, ss grpc.ServerStream, info *grpc.StreamServerInfo, handler grpc.StreamHandler) error {
	return func(srv interface{}, ss grpc.ServerStream, info *grpc.StreamServerInfo, handler grpc.StreamHandler) error {

		ipList := []string{}
		if serviceDetail.AccessControl.WhiteList != "" {
			ipList = strings.Split(serviceDetail.AccessControl.WhiteList, ",")
		}

		peerCtx, ok := peer.FromContext(ss.Context())
		if !ok {
			return errors.New("peer not found with context")
		}
		peerAddr := peerCtx.Addr.String()
		addrPos := strings.LastIndex(peerAddr, ":")
		clientIP := peerAddr[0:addrPos]

		if serviceDetail.AccessControl.OpenAuth == 1 && len(ipList) > 0 {
			if !public.InStringSlice(ipList, clientIP) {
				return errors.New(fmt.Sprintf("%s not in white ip list", clientIP))
			}
		}

		err := handler(srv, ss)
		if err != nil {
			log.Printf("RPC failed with error %v\n", err)
			return err
		}
		return nil
	}
}
