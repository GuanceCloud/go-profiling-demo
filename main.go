package main

import (
	"bytes"
	"compress/gzip"
	"context"
	_ "embed"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"net/http/pprof"
	"os"
	"path/filepath"
	"runtime"
	"sort"
	"strconv"
	"strings"
	"sync"
	"sync/atomic"
	"time"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	"gopkg.in/DataDog/dd-trace-go.v1/ddtrace/tracer"
	"gopkg.in/DataDog/dd-trace-go.v1/profiler"
)

const BaseServiceName = "go-profiling-demo"

type ctxKeyStruct struct {
}

var serviceNameKey = ctxKeyStruct{}

//go:embed movies5000.json.gz
var moviesJSON []byte

var serviceId = func() *atomic.Int64 {
	return &atomic.Int64{}
}()

func resetServiceID() {
	serviceId.Store(0)
}

func getCurServID() string {
	return strconv.FormatInt(serviceId.Load(), 10)
}

func getNextServID() string {
	serviceId.Add(1)
	return getCurServID()
}

func getCurServName() string {
	return fmt.Sprintf("%s-%s", BaseServiceName, getCurServID())
}

func getNextServName() string {
	return fmt.Sprintf("%s-%s", BaseServiceName, getNextServID())
}

var movies = func() []Movie {
	movies, err := readMovies()
	if err != nil {
		panic(err)
	}
	return movies
}()

type Movie struct {
	Title       string  `json:"title"`
	VoteAverage float64 `json:"vote_average"`
	ReleaseDate string  `json:"release_date"`
}

func GetCallerFuncName() string {
	pcs := make([]uintptr, 1)
	if runtime.Callers(2, pcs) < 1 {
		return ""
	}
	frame, _ := runtime.CallersFrames(pcs).Next()

	base := filepath.Base(frame.Function)

	if strings.ContainsRune(base, '.') {
		return filepath.Ext(base)[1:]
	}
	return base
}

func readMovies() ([]Movie, error) {
	r, err := gzip.NewReader(bytes.NewReader(moviesJSON))
	if err != nil {
		return nil, fmt.Errorf("gzip new reader from *FILE fail: %w", err)
	}
	defer r.Close()

	var mov []Movie

	if err := json.NewDecoder(r).Decode(&mov); err != nil {
		return nil, fmt.Errorf("json unmarshal fail: %w", err)
	}

	return mov, nil
}

func isENVTrue(key string) bool {
	val := strings.ToLower(strings.TrimSpace(os.Getenv(key)))
	switch val {
	case "", "0", "false":
		return false
	}
	return true
}

func sendHtmlRequest(ctx context.Context, bodyText string, servName string) {
	newSpan, _ := tracer.StartSpanFromContext(ctx, GetCallerFuncName(),
		tracer.ServiceName(servName))
	defer newSpan.Finish()

	req, err := http.NewRequest(http.MethodGet, "https://tv189.com/", strings.NewReader(strings.Repeat(bodyText, 1000)))

	if err != nil {
		log.Println(err)
		return
	}

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		log.Println(err)
		return
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		log.Println(err)
		return
	}

	log.Println(string(body))
}

func fibonacci(ctx context.Context, n int, servName string) int {
	if n <= 2 {
		return 1
	}
	if n%31 == 0 {
		return fibonacciWithTrace(ctx, n-1, servName) + fibonacciWithTrace(ctx, n-2, servName)
	} else if n%37 == 0 {
		return fibonacciWithTrace(ctx, n-1, servName) + fibonacciWithTrace(ctx, n-2, servName)
	}
	return fibonacci(ctx, n-1, servName) + fibonacci(ctx, n-2, servName)
}

func fibonacciWithTrace(ctx context.Context, n int, servName string) int {
	span, newCtx := tracer.StartSpanFromContext(ctx, GetCallerFuncName(),
		tracer.ServiceName(servName))
	defer span.Finish()
	return fibonacci(newCtx, n-1, servName) + fibonacci(newCtx, n-2, servName)
}

func httpReqWithTrace(ctx context.Context) {
	span, newCtx := tracer.StartSpanFromContext(ctx, GetCallerFuncName(),
		tracer.ServiceName(getNextServName()))
	defer span.Finish()

	bodyText := `
黄河远上白云间，一片孤城万仞山。
羌笛何须怨杨柳，春风不度玉门关。
少小离家老大回，乡音无改鬓毛衰。
儿童相见不相识，笑问客从何处来。
`

	for i := 0; i < 10; i++ {
		sendHtmlRequest(newCtx, bodyText, getCurServName())
	}
}

func main() {

	if isENVTrue("DD_TRACE_ENABLED") {
		tracer.Start(
			tracer.WithUniversalVersion("v0.8.888"),
		)
		defer tracer.Stop()
	}

	if isENVTrue("DD_PROFILING_ENABLED") {
		err := profiler.Start(
			profiler.WithProfileTypes(
				profiler.CPUProfile,
				profiler.HeapProfile,

				// The profiles below are disabled by default to keep overhead
				// low, but can be enabled as needed.
				profiler.BlockProfile,
				profiler.MutexProfile,
				profiler.GoroutineProfile,
				profiler.MetricsProfile,
			),
		)

		if err != nil {
			log.Fatal(err)
		}

		defer profiler.Stop()
	}

	router := gin.New()
	//router.Use(gintrace.Middleware("go-profiling-demo"))

	// Access-Control-*
	router.Use(cors.New(cors.Config{
		AllowAllOrigins:  true,
		AllowMethods:     []string{"GET", "POST", "PUT", "PATCH", "DELETE", "HEAD", "OPTIONS"},
		AllowCredentials: true,
		AllowHeaders:     []string{"*"},
		MaxAge:           time.Hour * 24,
	}))

	router.GET("/movies", func(ctx *gin.Context) {
		resetServiceID()

		spanCtx, err := tracer.Extract(tracer.HTTPHeadersCarrier(ctx.Request.Header))
		if err != nil {
			log.Printf("unable to extract span context from request header: %s", err)
		}

		if spanCtx != nil {
			spanCtx.ForeachBaggageItem(func(k, v string) bool {
				log.Printf("span context extracted key value %s: %s\n", k, v)
				return true
			})
		}

		span := tracer.StartSpan("get_movies", tracer.ChildOf(spanCtx),
			tracer.ServiceName(getNextServName()))
		newCtx := tracer.ContextWithSpan(ctx.Request.Context(), span)
		defer span.Finish()

		var wg sync.WaitGroup
		wg.Add(2)

		go func(ctx context.Context) {

			defer wg.Done()
			param := 42
			log.Printf("fibonacci(%d) = %d\n", param, fibonacci(ctx, param, getNextServName()))
		}(newCtx)

		go func(ctx context.Context) {
			defer wg.Done()
			httpReqWithTrace(ctx)
		}(newCtx)

		q := ctx.Request.FormValue("q")

		moviesCopy := make([]Movie, len(movies))
		copy(moviesCopy, movies)

		//func() {
		//	request, err := http.NewRequestWithContext(tracer.ContextWithSpan(ctx.Request.Context(), span),
		//		http.MethodPost, "http://127.0.0.1:5888/foobar", nil)
		//	if err != nil {
		//		log.Println("unable to new request: ", err)
		//		return
		//	}
		//	err = tracer.Inject(span.Context(), tracer.HTTPHeadersCarrier(request.Header))
		//	if err != nil {
		//		log.Println("unable to inject span to request: ", err)
		//		return
		//	}
		//	resp, err := http.DefaultClient.Do(request)
		//	if err != nil {
		//		log.Println("unable to request go-http-client")
		//		return
		//	}
		//	defer resp.Body.Close()
		//
		//	body, err := io.ReadAll(resp.Body)
		//	if err != nil {
		//		log.Println("unable to read request body: ", err)
		//	}
		//
		//	fmt.Println("response: ", string(body))
		//}()

		sort.Slice(moviesCopy, func(i, j int) bool {
			time.Sleep(time.Microsecond * 10)
			t1, err := time.Parse("2006-01-02", moviesCopy[i].ReleaseDate)
			if err != nil {
				return false
			}
			t2, err := time.Parse("2006-01-02", moviesCopy[j].ReleaseDate)
			if err != nil {
				return true
			}
			return t1.After(t2)
		})

		if q != "" {
			q = strings.ToUpper(q)
			matchCount := 0
			for idx, m := range moviesCopy {
				if strings.Contains(strings.ToUpper(m.Title), q) && idx != matchCount {
					moviesCopy[matchCount] = moviesCopy[idx]
					matchCount++
				}
			}
			moviesCopy = moviesCopy[:matchCount]
		}

		encoder := json.NewEncoder(ctx.Writer)
		if err := encoder.Encode(moviesCopy); err != nil {
			log.Printf("encode into json fail: %s", err)
			ctx.Writer.WriteHeader(http.StatusInternalServerError)
		}
		wg.Wait()
	})

	pprofIndex := func(ctx *gin.Context) {
		pprof.Index(ctx.Writer, ctx.Request)
	}

	router.GET("/debug/pprof", func(ctx *gin.Context) {
		ctx.Redirect(http.StatusMovedPermanently, "/debug/pprof/")
	})

	pg := router.Group("/debug/pprof")
	pg.GET("/", pprofIndex)
	pg.GET("/:name", pprofIndex)
	pg.GET("/cmdline", func(ctx *gin.Context) {
		pprof.Cmdline(ctx.Writer, ctx.Request)
	})
	pg.GET("/profile", func(ctx *gin.Context) {
		pprof.Profile(ctx.Writer, ctx.Request)
	})
	pg.GET("/symbol", func(ctx *gin.Context) {
		pprof.Symbol(ctx.Writer, ctx.Request)
	})
	pg.GET("/trace", func(ctx *gin.Context) {
		pprof.Trace(ctx.Writer, ctx.Request)
	})

	if err := http.ListenAndServe(":8080", router); err != nil {
		log.Fatal(err)
	}
}
