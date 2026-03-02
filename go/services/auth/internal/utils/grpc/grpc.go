package grpcutil

import (
	"crypto/tls"
	"strings"
	"time"

	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials"
	"google.golang.org/grpc/keepalive"

	insecurepkg "google.golang.org/grpc/credentials/insecure"
)

func Dial(addr string, insecure bool, opts ...grpc.DialOption) (*grpc.ClientConn, func(), error) {
	kp := keepalive.ClientParameters{
		Time:                2 * time.Minute,
		PermitWithoutStream: true,
	}

	dialOptions := []grpc.DialOption{
		grpc.WithKeepaliveParams(kp),
	}

	if insecure {
		dialOptions = append(dialOptions, grpc.WithTransportCredentials(insecurepkg.NewCredentials()))
	} else {
		config := &tls.Config{
			MinVersion:         tls.VersionTLS12,
			InsecureSkipVerify: false,
			ServerName:         strings.Split(addr, ":")[0],
		}
		dialOptions = append(dialOptions, grpc.WithTransportCredentials(credentials.NewTLS(config)))
	}

	// append param dial options
	dialOptions = append(dialOptions, opts...)

	conn, err := grpc.NewClient(addr, dialOptions...)
	var cleanUpFunc func()
	if err != nil {
		cleanUpFunc = func() {}
	} else {
		cleanUpFunc = func() {
			conn.Close()
		}
	}

	return conn, cleanUpFunc, err
}
