Type: inuse_space
Time: 2026-01-18 17:05:56 MSK
Duration: 60.01s, Total samples = 2568.07kB
Showing nodes accounting for -534.98kB, 20.83% of 2568.07kB total
      flat  flat%   sum%        cum   cum%
 1536.12kB 59.82% 59.82%  1536.12kB 59.82%  reflect.New
-1032.02kB 40.19% 19.63% -1032.02kB 40.19%  github.com/jackc/chunkreader/v2.(*ChunkReader).newBuf (inline)
-1024.05kB 39.88% 20.25% -1024.05kB 39.88%  encoding/json.(*decodeState).literalStore
 -528.17kB 20.57% 40.81%  -528.17kB 20.57%  unicode/utf16.Encode
  513.13kB 19.98% 20.83%   513.13kB 19.98%  reflect.growslice
  512.02kB 19.94%  0.89%   512.02kB 19.94%  net/textproto.(*Reader).ReadLine (inline)
 -512.01kB 19.94% 20.83%  -512.01kB 19.94%  github.com/konkovaanna23/shortener/internal/model.(*Storage).randomString
         0     0% 20.83% -1032.02kB 40.19%  database/sql.(*DB).retry
         0     0% 20.83% -1032.02kB 40.19%  database/sql.(*Stmt).QueryContext
         0     0% 20.83% -1032.02kB 40.19%  database/sql.(*Stmt).QueryContext.func1
         0     0% 20.83% -1032.02kB 40.19%  database/sql.(*Stmt).QueryRow (inline)
         0     0% 20.83% -1032.02kB 40.19%  database/sql.(*Stmt).QueryRowContext
         0     0% 20.83% -1032.02kB 40.19%  database/sql.ctxDriverStmtQuery
         0     0% 20.83% -1032.02kB 40.19%  database/sql.rowsiFromStatement
         0     0% 20.83%  1025.20kB 39.92%  encoding/json.(*decodeState).array
         0     0% 20.83%   512.07kB 19.94%  encoding/json.(*decodeState).object
         0     0% 20.83%  1025.20kB 39.92%  encoding/json.(*decodeState).unmarshal
         0     0% 20.83%  1025.20kB 39.92%  encoding/json.(*decodeState).value
         0     0% 20.83%  1025.20kB 39.92%  encoding/json.Unmarshal
         0     0% 20.83%  1536.12kB 59.82%  encoding/json.indirect
         0     0% 20.83%    -1047kB 40.77%  github.com/go-chi/chi/v5.(*Mux).ServeHTTP
         0     0% 20.83%    -1047kB 40.77%  github.com/go-chi/chi/v5.(*Mux).routeHTTP
         0     0% 20.83% -1032.02kB 40.19%  github.com/jackc/chunkreader/v2.(*ChunkReader).Next      
         0     0% 20.83% -1032.02kB 40.19%  github.com/jackc/pgconn.(*PgConn).ExecPrepared
         0     0% 20.83% -1032.02kB 40.19%  github.com/jackc/pgconn.(*PgConn).execExtendedSuffix     
         0     0% 20.83% -1032.02kB 40.19%  github.com/jackc/pgconn.(*PgConn).peekMessage
         0     0% 20.83% -1032.02kB 40.19%  github.com/jackc/pgconn.(*ResultReader).readUntilRowDescription
         0     0% 20.83% -1032.02kB 40.19%  github.com/jackc/pgproto3/v2.(*Frontend).Receive
         0     0% 20.83% -1032.02kB 40.19%  github.com/jackc/pgx/v4.(*Conn).Query
         0     0% 20.83% -1032.02kB 40.19%  github.com/jackc/pgx/v4/stdlib.(*Conn).QueryContext      
         0     0% 20.83% -1032.02kB 40.19%  github.com/jackc/pgx/v4/stdlib.(*Stmt).QueryContext      
         0     0% 20.83%    -1047kB 40.77%  github.com/konkovaanna23/shortener/internal/handler.(*Server).newJSONBatchURL
         0     0% 20.83%    -1047kB 40.77%  github.com/konkovaanna23/shortener/internal/handler.NewServer.AuthMiddleware.func1.1
         0     0% 20.83%    -1047kB 40.77%  github.com/konkovaanna23/shortener/internal/handler/middleware.CompressMiddleware.func1
         0     0% 20.83%    -1047kB 40.77%  github.com/konkovaanna23/shortener/internal/handler/middleware.LoggingMiddleware.func1
         0     0% 20.83%  -512.01kB 19.94%  github.com/konkovaanna23/shortener/internal/model.(*Storage).GenerateShortURL (inline)
         0     0% 20.83% -1544.03kB 60.12%  github.com/konkovaanna23/shortener/internal/service.(*Converter).AddURLForBatch
         0     0% 20.83% -1032.02kB 40.19%  github.com/konkovaanna23/shortener/internal/service.(*Converter).TranStoreURLInDB
         0     0% 20.83%  -528.17kB 20.57%  github.com/sirupsen/logrus.(*Entry).Log
         0     0% 20.83%  -528.17kB 20.57%  github.com/sirupsen/logrus.(*Entry).log
         0     0% 20.83%  -528.17kB 20.57%  github.com/sirupsen/logrus.(*Entry).write
         0     0% 20.83%  -528.17kB 20.57%  github.com/sirupsen/logrus.(*Logger).Info (inline)       
         0     0% 20.83%  -528.17kB 20.57%  github.com/sirupsen/logrus.(*Logger).Log
         0     0% 20.83%  -528.17kB 20.57%  github.com/sirupsen/logrus.Info (inline)
         0     0% 20.83%  -528.17kB 20.57%  internal/poll.(*FD).Write
         0     0% 20.83%  -528.17kB 20.57%  internal/poll.(*FD).writeConsole
         0     0% 20.83%   512.02kB 19.94%  net/http.(*conn).readRequest
         0     0% 20.83%  -534.98kB 20.83%  net/http.(*conn).serve
         0     0% 20.83%    -1047kB 40.77%  net/http.HandlerFunc.ServeHTTP
         0     0% 20.83%   512.02kB 19.94%  net/http.readRequest
         0     0% 20.83%    -1047kB 40.77%  net/http.serverHandler.ServeHTTP
         0     0% 20.83%  -528.17kB 20.57%  os.(*File).Write
         0     0% 20.83%  -528.17kB 20.57%  os.(*File).write (inline)
         0     0% 20.83%   513.13kB 19.98%  reflect.Value.Grow
         0     0% 20.83%   513.13kB 19.98%  reflect.Value.grow