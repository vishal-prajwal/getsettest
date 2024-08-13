package client

type Sample interface {
}

// func NewSampleClient() (Sample, error) {// 	// The service name needs to match what we set in the nomad provisioning "svc-sample-grpc"
// 	conn, err := jungleegames.GrpcDial("svc-sample-grpc")
// 	if err != nil {
// 		return nil, errors.Wrap(err, "failed to connect to the sample grpc service")
// 	}

// 	return &grpcClient{
// 		conn: api.NewSampleServiceClient(conn),
// 	}, nil
// }

// type grpcClient struct {
// 	conn api.SampleServiceClient
// }
