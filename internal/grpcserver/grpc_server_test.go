package grpcserver

import (
	"context"
	"net"
	"strings"
	"testing"

	"github.com/konkovaanna23/shortener/internal/model"
	"github.com/konkovaanna23/shortener/internal/service"
	ss "github.com/konkovaanna23/shortener/pkg/shortenerservice"
	"github.com/stretchr/testify/assert"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/credentials/insecure"
	"google.golang.org/grpc/metadata"
	"google.golang.org/grpc/status"
)

var host = "localhost:4500"
var shortHost = "http://localhost:8080"

func setupServer(t *testing.T) *grpc.Server {
	listener, err := net.Listen("tcp", host)
	if err != nil {
		return nil
	}

	ctx := context.Background()
	converter := service.NewConverter(ctx, "http://localhost:8080", "", nil, 100, 10, 1)

	server := grpc.NewServer(
		grpc.UnaryInterceptor(UnaryInterceptor),
	)

	grpcServer := NewGrpcServer(converter, "", "")
	ss.RegisterShortenerServiceServer(server, grpcServer)

	go func() {
		if err := server.Serve(listener); err != nil {
			t.Logf("сервер остановлен: %v", err)
		}
	}()

	return server
}

func TestShortenURL_Success(t *testing.T) {

	server := setupServer(t)
	defer server.Stop()

	conn, err := grpc.NewClient(host, grpc.WithTransportCredentials(insecure.NewCredentials()))
	assert.NoError(t, err)
	defer conn.Close()

	client := ss.NewShortenerServiceClient(conn)

	url := "https://example.com"
	user := "user123"

	ctx := metadata.NewOutgoingContext(context.Background(), metadata.Pairs("authorization", "Bearer "+user))

	request := &ss.URLShortenRequest_builder{Url: url}

	_, err = client.ShortenURL(ctx, request.Build())

	assert.NoError(t, err)

}

func TestShortenURL_Unauthorized(t *testing.T) {

	server := setupServer(t)
	defer server.Stop()

	conn, err := grpc.NewClient(host, grpc.WithTransportCredentials(insecure.NewCredentials()))
	assert.NoError(t, err)
	defer conn.Close()

	client := ss.NewShortenerServiceClient(conn)

	ctx := context.Background()
	url := "https://example.com"
	request := &ss.URLShortenRequest_builder{Url: url}

	_, err = client.ShortenURL(ctx, request.Build())

	// Assert
	st, ok := status.FromError(err)
	assert.True(t, ok)
	assert.Equal(t, codes.Unauthenticated, st.Code())
}

func TestShortenURL_Conflict(t *testing.T) {

	server := setupServer(t)
	defer server.Stop()

	conn, err := grpc.NewClient(host, grpc.WithTransportCredentials(insecure.NewCredentials()))
	assert.NoError(t, err)
	defer conn.Close()

	client := ss.NewShortenerServiceClient(conn)

	ctx := context.Background()
	url := "https://example.com"
	user := "user123"

	ctx = metadata.NewOutgoingContext(ctx, metadata.Pairs("authorization", "Bearer "+user))

	request := &ss.URLShortenRequest_builder{Url: url}

	_, err = client.ShortenURL(ctx, request.Build())
	assert.NoError(t, err)

	_, err = client.ShortenURL(ctx, request.Build())
	assert.NoError(t, err)

	st, ok := status.FromError(err)
	assert.True(t, ok)
	assert.Equal(t, codes.AlreadyExists, st.Code())
}

func TestExpandURL_Success(t *testing.T) {

	server := setupServer(t)
	defer server.Stop()

	conn, err := grpc.NewClient(host, grpc.WithTransportCredentials(insecure.NewCredentials()))
	assert.NoError(t, err)
	defer conn.Close()

	client := ss.NewShortenerServiceClient(conn)

	url := "https://example.com"
	user := "user123"

	ctx := context.Background()
	ctx = metadata.NewOutgoingContext(ctx, metadata.Pairs("authorization", "Bearer "+user))

	request := &ss.URLShortenRequest_builder{Url: url}

	resp, err := client.ShortenURL(ctx, request.Build())
	assert.NoError(t, err)

	requestExpan := &ss.URLExpandRequest_builder{Id: strings.ReplaceAll(resp.GetResult(), shortHost+"/", "")}

	respExpand, err := client.ExpandURL(ctx, requestExpan.Build())

	assert.NoError(t, err)
	assert.Equal(t, url, respExpand.GetResult())
}

func TestListUserURLs_Success(t *testing.T) {

	server := setupServer(t)
	defer server.Stop()

	conn, err := grpc.NewClient(host, grpc.WithTransportCredentials(insecure.NewCredentials()))
	assert.NoError(t, err)
	defer conn.Close()

	client := ss.NewShortenerServiceClient(conn)

	ctx := context.Background()
	url := "https://example.com"
	user := "user123"

	ctx = metadata.NewOutgoingContext(ctx, metadata.Pairs("authorization", "Bearer "+user))

	request := &ss.URLShortenRequest_builder{Url: url}

	resp, err := client.ShortenURL(ctx, request.Build())
	assert.NoError(t, err)

	result := &model.DescriptionURL{
		Short: strings.ReplaceAll(resp.GetResult(), shortHost+"/", ""), Original: url,
	}

	ctx = metadata.NewOutgoingContext(ctx, metadata.Pairs("authorization", "Bearer "+user))

	respList, err := client.ListUserURLs(ctx, &ss.Empty{})

	assert.NoError(t, err)
	assert.Len(t, respList.GetUrl(), 1)
	assert.Equal(t, result.Short, strings.ReplaceAll(respList.GetUrl()[0].GetShortUrl(), shortHost+"/", ""))
	assert.Equal(t, result.Original, respList.GetUrl()[0].GetOriginalUrl())
}
