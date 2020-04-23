package main

import (
	"github.com/kataras/golog"
	traceutils "github.com/opentracing-contrib/go-grpc"
	"github.com/opentracing/opentracing-go"
	"github.com/uber/jaeger-client-go"
	jaegercfg "github.com/uber/jaeger-client-go/config"
	jaegerlog "github.com/uber/jaeger-client-go/log"
	"github.com/uber/jaeger-lib/metrics"
	"google.golang.org/grpc"
	searchPostgres "joblessness/haha/search/repository/postgres"
	"joblessness/haha/utils/database"
	searchRpc "joblessness/searchService/rpc"
	searchServer "joblessness/searchService/server"
	"log"
	"net"
)

func main() {
	jaegerCfgInstance := jaegercfg.Configuration{
		ServiceName: "session",
		Sampler: &jaegercfg.SamplerConfig{
			Type:  jaeger.SamplerTypeConst,
			Param: 1,
		},
		Reporter: &jaegercfg.ReporterConfig{
			LogSpans:           true,
			LocalAgentHostPort: "localhost:6831",
		},
	}

	tracer, closer, err := jaegerCfgInstance.NewTracer(
		jaegercfg.Logger(jaegerlog.StdLogger),
		jaegercfg.Metrics(metrics.NullFactory),
	)

	if err != nil {
		log.Fatal("cannot create tracer", err)
	}

	opentracing.SetGlobalTracer(tracer)
	defer closer.Close()

	db, err := database.OpenDatabase()
	if err != nil {
		golog.Error(err.Error())
		return
	}
	defer db.Close()

	repo := searchPostgres.NewSearchRepository(db)
	list, err := net.Listen("tcp", "127.0.0.1:8002")
	if err != nil {
		golog.Error(err.Error())
		return
	}

	server := grpc.NewServer(grpc.UnaryInterceptor(traceutils.OpenTracingServerInterceptor(tracer)))
	searchRpc.RegisterSearchServer(server, searchServer.NewSearchServer(repo))
	server.Serve(list)
}
