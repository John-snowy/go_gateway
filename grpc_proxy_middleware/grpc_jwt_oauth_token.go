package grpc_proxy_middleware

import (
	"errors"
	"github.com/go_gateway/dao"
	"github.com/go_gateway/public"
	"google.golang.org/grpc"
	"google.golang.org/grpc/metadata"
	"log"
	"strings"
)

func GrpcJwtOAuthTokenMiddleware(serviceDetail *dao.ServiceDetail) func(srv interface{}, ss grpc.ServerStream, info *grpc.StreamServerInfo, handler grpc.StreamHandler) error {
	return func(srv interface{}, ss grpc.ServerStream, info *grpc.StreamServerInfo, handler grpc.StreamHandler) error {
		md, ok := metadata.FromIncomingContext(ss.Context())
		if !ok {
			return errors.New("miss metadata from context")
		}
		authToken := ""
		auth := md.Get("Authorization")
		if len(auth) > 0 {
			authToken = auth[0]
		}
		// decode jwt token
		// app_id  app_list  => appInfo
		// appInfo  =>  gin.context
		token := strings.ReplaceAll(authToken, "Bearer ", "")
		appMatched := false
		if token != "" {
			claims, err := public.JwtDecode(token)
			if err != nil {
				return err
			}
			appList := dao.AppManagerHandler.GetAppList()
			for _, appInfo := range appList {
				if appInfo.AppID == claims.Issuer {
					md.Set("app", public.Obj2Json(appInfo))
					//c.Set("appDetail", appInfo)
					appMatched = true
					break
				}
			}
		}
		if serviceDetail.AccessControl.OpenAuth == 1 && !appMatched {
			return errors.New("not match valid app")
		}

		if err := handler(srv, ss); err != nil {
			log.Printf("RPC failed with error %v\n", err)
			return err
		}
		return nil

	}
}
