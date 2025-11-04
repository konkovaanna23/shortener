package handler

import (
	"bytes"
	"compress/gzip"
	"fmt"
	"io"
	"net/http"
	"strings"
	"time"

	"github.com/sirupsen/logrus"
)

type ResponseCompressLogger struct {
	http.ResponseWriter
	header http.Header
	status int
	size   int
	buf    *bytes.Buffer
}

func (l *ResponseCompressLogger) WriteHeader(code int) {
	if l.status != 0 {
		return
	}
	l.status = code
	l.ResponseWriter.WriteHeader(code)
}

func (l *ResponseCompressLogger) Header() http.Header {
	return l.header
}

func (l *ResponseCompressLogger) Write(b []byte) (int, error) {
	if l.status == 0 {
		l.status = http.StatusOK
	}
	size, err := l.buf.Write(b)
	l.size += size
	return size, err
}

func NewResponseCompressLogger(w http.ResponseWriter) *ResponseCompressLogger {
	return &ResponseCompressLogger{
		header: make(http.Header),
		status: http.StatusOK,
		buf:    new(bytes.Buffer),
	}
}

func (s *Server) LoggingMiddleware(next http.Handler) http.HandlerFunc {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		start := time.Now()
		uri := r.URL.RequestURI()
		method := r.Method

		newRequest, err := s.decompressRequest(r)
		if err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}

		writer := NewResponseCompressLogger(w)

		next.ServeHTTP(writer, newRequest)

		acceptEncoding := r.Header.Get("Accept-Encoding")
		contentType := writer.header.Get("Content-Type")
		var needCompress bool
		if strings.Contains(acceptEncoding, "gzip") &&
			(contentType == "application/json" || contentType == "text/html") {
			needCompress = true
		}

		for k, vv := range writer.header {
			for _, v := range vv {
				w.Header().Add(k, v)
			}
		}
		var body []byte
		if needCompress {
			w.Header().Set("Content-Encoding", "gzip")
			var err error
			body, err = s.compessResponse(writer)
			if err != nil {
				http.Error(w, err.Error(), http.StatusBadRequest)
				return
			}
		} else {
			body = writer.buf.Bytes()
		}

		w.WriteHeader(writer.status)

		if _, err := w.Write(body); err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}

		logrus.WithFields(logrus.Fields{
			"method":      method,
			"uri":         uri,
			"status":      writer.status,
			"duration":    time.Since(start).String(),
			"content_len": writer.size,
		}).Info("Запрос выполнен")
	})
}

func (s *Server) decompressRequest(r *http.Request) (*http.Request, error) {
	if r.Header.Get("Content-Encoding") == "gzip" &&
		(r.Header.Get("Content-Type") == "application/json" || r.Header.Get("Content-Type") == "text/html") {
		gz, err := gzip.NewReader(r.Body)
		if err != nil {
			return nil, fmt.Errorf("%s", "Ошибка декодирования gzip")
		}
		defer gz.Close()

		decompressedBody, err := io.ReadAll(gz)
		if err != nil {
			return nil, fmt.Errorf("%s", "Ошибка чтения разжатого тела")
		}
		newReq := r.Clone(r.Context())
		newReq.Body = io.NopCloser(bytes.NewBuffer(decompressedBody))
		newReq.ContentLength = int64(len(decompressedBody))
		newReq.Header.Del("Content-Encoding")
		return newReq, nil
	} else {
		return r, nil
	}
}

func (s *Server) compessResponse(resp *ResponseCompressLogger) ([]byte, error) {
	var gzipBuf bytes.Buffer
	gz := gzip.NewWriter(&gzipBuf)
	_, err := gz.Write(resp.buf.Bytes())
	if err != nil {
		return nil, err
	}
	err = gz.Close()
	if err != nil {
		return nil, err
	}
	return gzipBuf.Bytes(), nil
}
