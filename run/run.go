package run

import (
	"fmt"
	"go.uber.org/zap"
	"google.golang.org/grpc"
	"log"
	"net"
	"os"
	"os/signal"
	"syscall"
	"usdt/config"
	"usdt/internal/db"
	"usdt/internal/infrastructure/requestAPI/garantex"
	"usdt/internal/modules/controller"
	"usdt/internal/modules/service"
	"usdt/internal/modules/storage"
	proto "usdt/internal/proto/usdt_proto"
)

func shutdown(adapter *db.DbAdapter, logger *zap.Logger, grpcServer *grpc.Server) {
	logger.Info("Получен сигнал завершения работы. Начинаем graceful shutdown...")
	logger.Info("Остановка gRPC сервера...")
	grpcServer.GracefulStop() // Используем ctx
	logger.Info("gRPC сервер остановлен.")
	logger.Info("Закрытие соединения с базой данных...")
	err := adapter.Close()
	if err != nil {
		logger.Error("Ошибка при закрытии соединения с базой данных:", zap.Error(err))
	}
	logger.Info("Соединение с базой данных закрыто.")
	logger.Info("Сервис завершил работу.")
}

func Run(adapter *db.DbAdapter, logger *zap.Logger, conf config.Config, grpcServer *grpc.Server) {
	storageusddt := storage.NewUsdtStorage(adapter)
	api := garantex.NewGrantexAPI("")
	serviceusdt := service.NewUsdtService(storageusddt, api)
	controllerusdt := controller.NewController(serviceusdt, logger)
	proto.RegisterAuthServiceServer(grpcServer, controllerusdt)
	lis, err := net.Listen("tcp", fmt.Sprintf(":%s", conf.Port))
	if err != nil {
		log.Fatalf("failed to listen: %v", err)
	}

	// Обработка сигналов
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)

	go func() {
		<-quit
		shutdown(adapter, logger, grpcServer)
	}()

	logger.Info(fmt.Sprintf("USDT service started on port: %s", conf.Port))
	if err = grpcServer.Serve(lis); err != nil {
		logger.Error(fmt.Sprintf("failed to serve: %v", err))
	}
}
