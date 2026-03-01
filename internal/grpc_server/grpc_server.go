package grpc_server

import (
	"context"
	"errors"
	"strings"
	"time"

	"github.com/konkovaanna23/shortener/internal/handler/audit"
	"github.com/konkovaanna23/shortener/internal/model"
	"github.com/konkovaanna23/shortener/internal/service"
	ss "github.com/konkovaanna23/shortener/pkg/shortenerservice"
	"github.com/sirupsen/logrus"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/metadata"
	"google.golang.org/grpc/status"
)

const userKey = "user_id"

type GrpcServer struct {
	ss.UnimplementedShortenerServiceServer

	converter *service.Converter
	auditor   *audit.Publisher
}

func NewGrpcServer(converter *service.Converter, auditFile, auditURL string) *GrpcServer {
	return &GrpcServer{
		converter: converter,
		auditor:   audit.NewAuditor(auditFile, auditURL),
	}
}

func UnaryInterceptor(ctx context.Context, req interface{}, info *grpc.UnaryServerInfo, handler grpc.UnaryHandler) (interface{}, error) {
	md, ok := metadata.FromIncomingContext(ctx)
	if !ok {
		return nil, status.Errorf(codes.Unauthenticated, "нет метаданных")
	}

	values := md["authorization"]
	if len(values) == 0 {
		return nil, status.Errorf(codes.Unauthenticated, "требуется заголовок authorization")
	}

	authHeader := values[0]
	parts := strings.Split(authHeader, " ")
	if len(parts) != 2 || strings.ToLower(parts[0]) != "bearer" {
		return nil, status.Errorf(codes.Unauthenticated, "некорректный формат Authorization: ожидается 'Bearer <token>'")
	}

	token := parts[1]

	ctx = context.WithValue(ctx, userKey, token)

	return handler(ctx, req)
}

func (g *GrpcServer) ShortenURL(ctx context.Context, req *ss.URLShortenRequest) (*ss.URLShortenResponse, error) {
	user, ok := ctx.Value(userKey).(string)
	if !ok {
		return nil, status.Errorf(codes.Unauthenticated, "пользователь не найден в контексте")
	}
	var urlRequest *service.URLRequest
	logrus.Info("ShortenURL Заданный URL:", req.GetUrl(), " заданный user:", user)
	urlRequest = &service.URLRequest{
		URL: req.GetUrl(),
	}
	result, err := g.converter.AddURLForRequest(urlRequest, user)
	if err != nil {
		if errors.Is(err, model.ErrorConflictURL) {
			return nil, status.Errorf(codes.AlreadyExists, "URL уже существует")
		} else {
			return nil, status.Errorf(codes.Internal, "ошибка добавления URL: %v", err)
		}
	}
	response := &ss.URLShortenResponse_builder{
		Result: result.URLShort,
	}
	g.auditor.Publish(audit.Event{Time: time.Now(),
		Action: "GRPC:ShortenURL",
		UserID: user,
		URL:    urlRequest.URL})
	return response.Build(), nil

}

func (g *GrpcServer) ExpandURL(ctx context.Context, req *ss.URLExpandRequest) (*ss.URLExpandResponse, error) {
	logrus.Info("ExpandURL Заданный URL:", req.GetId())
	sourceURL, err := g.converter.GetURL(req.GetId())
	if err != nil {
		return nil, status.Errorf(codes.Internal, "ошибка получения URL: %v", err)
	}
	response := &ss.URLExpandResponse_builder{
		Result: sourceURL,
	}
	g.auditor.Publish(audit.Event{Time: time.Now(),
		Action: "GRPC:ExpandURL",
		URL:    sourceURL,
	})
	return response.Build(), nil

}

func (g *GrpcServer) ListUserURLs(ctx context.Context, req *ss.Empty) (*ss.UserURLsResponse, error) {
	user, ok := ctx.Value(userKey).(string)
	if !ok {
		return nil, status.Errorf(codes.Unauthenticated, "пользователь не найден в контексте")
	}
	logrus.Info("ListUserURLs Заданный user:", user)
	result, err := g.converter.GetURLsForUser(user)
	if err != nil {
		return nil, status.Errorf(codes.Internal, "ошибка получения списка URL: %v", err)
	}

	resultURLs := make([]*ss.URLData, len(result))
	for i, r := range result {
		resultURLs[i] = &ss.URLData{}
		resultURLs[i].SetShortUrl(r.Short)
		resultURLs[i].SetOriginalUrl(r.Original)
	}

	response := &ss.UserURLsResponse{}

	response.SetUrl(resultURLs)

	return response, nil

}
