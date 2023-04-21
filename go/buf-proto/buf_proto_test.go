package bufproto_test

import (
	"bytes"
	"context"
	"fmt"
	"net/http/httptest"
	"testing"

	connect_go "github.com/bufbuild/connect-go"
	"github.com/ever0de/playground/buf-proto/proto"
	"github.com/ever0de/playground/buf-proto/proto/protoconnect"
	"github.com/stretchr/testify/assert"
)

var (
	serverDB [][]byte
)

func init() {
	serverDB = make([][]byte, 0)
	for i := 0; i < 100; i++ {
		serverDB = append(serverDB, []byte(fmt.Sprintf("record %d", i)))
	}
}

type handler struct{}

var _ protoconnect.ServiceHandler = (*handler)(nil)

func (h *handler) Subscribe(
	ctx context.Context,
	req *connect_go.Request[proto.SubscriptionRequest],
	stream *connect_go.ServerStream[proto.SubscriptionResponse],
) error {
	for _, record := range serverDB {
		for _, key := range req.Msg.Key {
			if bytes.Equal(record, key) {
				if err := stream.Send(&proto.SubscriptionResponse{
					Record: []*proto.Record{
						{
							Key:   record,
							Value: record,
						},
					},
				}); err != nil {
					return err
				}
			}
		}
	}

	return nil
}

func TestConnectGoService(t *testing.T) {
	//server
	_, handler := protoconnect.NewServiceHandler(&handler{})
	server := httptest.NewUnstartedServer(handler)
	server.EnableHTTP2 = true
	server.StartTLS()
	t.Cleanup(server.Close)

	//client
	sCilent := server.Client()
	client := protoconnect.NewServiceClient(sCilent, server.URL)

	stream, err := client.Subscribe(
		context.Background(),
		connect_go.NewRequest(&proto.SubscriptionRequest{
			Id: []byte("client"),
			Key: [][]byte{
				[]byte("record 1"),
				[]byte("record 2"),
			},
		}),
	)
	assert.NoError(t, err)

	for i := 0; i < 2; i++ {
		assert.True(t, stream.Receive())
		records := stream.Msg().GetRecord()
		assert.Equal(t, 1, len(records))

		record := records[0]

		assert.Equal(t, []byte(fmt.Sprintf("record %d", (i+1))), record.GetKey())
	}
	assert.False(t, stream.Receive())
}
