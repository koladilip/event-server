package api

import (
	"context"
	"net"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/koladilip/event-server/config"
	"github.com/koladilip/event-server/event"
	"github.com/koladilip/event-server/store"
	"github.com/segmentio/kafka-go"
	"go.uber.org/fx"
	"go.uber.org/zap"
)

func NewRouter() *gin.Engine {
	router := gin.New()
	router.Use(
		gin.Recovery(),
	)
	return router
}

func RegisterPublishAPI(logger *zap.Logger,
	kafkaWriter *kafka.Writer, router *gin.Engine) {
	router.POST("/publish", func(c *gin.Context) {
		jsonData, err := c.GetRawData()
		var sourceEvent event.SourceEvent
		if err == nil {
			sourceEvent, err = event.NewSourceEvent(jsonData)
		}
		if err != nil {
			c.JSON(400, gin.H{
				"message": "Error in paring the message",
			})
			return
		}
		err = kafkaWriter.WriteMessages(c.Request.Context(), kafka.Message{
			Topic: store.SourceEventTopic,
			Key:   []byte(sourceEvent.UserID),
			Value: jsonData,
		})
		if err != nil {
			logger.Error("Error in storing the message")
			c.JSON(503, gin.H{
				"message": "Error in storing the message",
			})
			return
		}
		c.JSON(200, gin.H{
			"message": "received",
		})
	})
}

func NewServer(ctx *config.BaseContext, config *config.Config, router *gin.Engine) *http.Server {
	return &http.Server{
		Addr:    config.Server.Port,
		Handler: router,
		BaseContext: func(_ net.Listener) context.Context {
			return ctx.Context
		},
	}
}

func StartServer(lc fx.Lifecycle, baseCtx *config.BaseContext,
	logger *zap.Logger, config *config.Config, server *http.Server) {
	lc.Append(fx.Hook{
		// To mitigate the impact of deadlocks in application startup and
		// shutdown, Fx imposes a time limit on OnStart and OnStop hooks. By
		// default, hooks have a total of 15 seconds to complete. Timeouts are
		// passed via Go's usual context.Context.
		OnStart: func(context.Context) error {
			logger.Info("Starting HTTP server.")
			// In production, we'd want to separate the Listen and Serve phases for
			// better error-handling.
			go server.ListenAndServe()
			return nil
		},
		OnStop: func(ctx context.Context) error {
			logger.Info("Stopping HTTP server.")
			baseCtx.Cancel()
			return server.Shutdown(ctx)
		},
	})
}

var Fx = fx.Options(
	fx.Provide(NewRouter),
	fx.Invoke(RegisterPublishAPI),
	fx.Provide(NewServer),
	fx.Invoke(StartServer),
)
