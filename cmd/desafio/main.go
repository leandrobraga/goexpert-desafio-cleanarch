package main

import (
	"database/sql"
	"fmt"
	"net"
	"net/http"

	graphql_handler "github.com/99designs/gqlgen/graphql/handler"
	"github.com/99designs/gqlgen/graphql/playground"
	"github.com/go-chi/chi/middleware"
	"github.com/go-chi/chi/v5"
	"github.com/leandrobraga/goexpert-desafio-cleanarch/configs"
	"github.com/leandrobraga/goexpert-desafio-cleanarch/internal/infra/database"
	"github.com/leandrobraga/goexpert-desafio-cleanarch/internal/infra/event"
	"github.com/leandrobraga/goexpert-desafio-cleanarch/internal/infra/event/handler"
	"github.com/leandrobraga/goexpert-desafio-cleanarch/internal/infra/graph"
	"github.com/leandrobraga/goexpert-desafio-cleanarch/internal/infra/grpc/pb"
	"github.com/leandrobraga/goexpert-desafio-cleanarch/internal/infra/grpc/service"
	"github.com/leandrobraga/goexpert-desafio-cleanarch/internal/infra/web"
	"github.com/leandrobraga/goexpert-desafio-cleanarch/internal/infra/web/webserver"
	"github.com/leandrobraga/goexpert-desafio-cleanarch/internal/usecase"
	"github.com/leandrobraga/goexpert-desafio-cleanarch/pkg/events"
	ampq "github.com/rabbitmq/amqp091-go"
	"google.golang.org/grpc"
	"google.golang.org/grpc/reflection"

	_ "github.com/go-sql-driver/mysql"
)

func main() {
	configs, err := configs.LoadConfig(".")
	if err != nil {
		panic(err)
	}

	db, err := sql.Open(configs.DBDriver, fmt.Sprintf("%s:%s@tcp(%s:%s)/%s", configs.DBUser, configs.DBPassword, configs.DBHost, configs.DBPort, configs.DBName))
	if err != nil {
		panic(err)
	}
	defer db.Close()

	rabbitMQChannel := getRabbitMQChannel()

	eventDispatcher := events.NewEventDispatcher()
	eventDispatcher.Register("OrderCreated", &handler.OrderCreatedHandler{
		RabbitMQChannel: rabbitMQChannel,
	})

	orderRepository := database.NewOrderRepository(db)
	orderCreated := event.NewOrderCreated("OrderCreated")

	createOrderUseCase := usecase.NewCreateOrderUseCase(orderRepository, orderCreated, eventDispatcher)
	listOrdersUsecase := usecase.NewListOrdersUseCase(orderRepository)

	webserver := webserver.NewWebServer(configs.WebServerPort)
	webOrderHandler := web.NewWebOrderHandler(orderRepository, orderCreated, eventDispatcher)
	// webserver.AddHandler("/order", webOrderHandler.Create)
	// webserver.AddHandler("/orderlist", webOrderHandler.List)
	webserver.Router.Route("/order", func(r chi.Router) {
		r.Use(middleware.Logger)
		r.Get("/", webOrderHandler.List)
		r.Post("/", webOrderHandler.Create)
	})
	fmt.Println("Starting web server on port", configs.WebServerPort)
	go webserver.Start()

	grpcServer := grpc.NewServer()
	orderService := service.NewOrderService(*createOrderUseCase, *listOrdersUsecase)
	pb.RegisterOrderServiceServer(grpcServer, orderService)
	//Para usar o cli grpc
	reflection.Register(grpcServer)
	fmt.Println("Starting gRPC server on port", configs.GRPCServerPort)
	lis, err := net.Listen("tcp", fmt.Sprintf(":%s", configs.GRPCServerPort))
	if err != nil {
		panic(err)
	}
	go grpcServer.Serve(lis)

	srv := graphql_handler.NewDefaultServer(graph.NewExecutableSchema(graph.Config{Resolvers: &graph.Resolver{
		CreateOrderUseCase: *createOrderUseCase,
		ListOrdersUseCase:  *listOrdersUsecase,
	}}))
	http.Handle("/", playground.Handler("GraphQL playground", "/query"))
	http.Handle("/query", srv)
	fmt.Println("Starting GraphQL server on port", configs.GraphQLServerPort)
	http.ListenAndServe(":"+configs.GraphQLServerPort, nil)

}

func getRabbitMQChannel() *ampq.Channel {
	// COLOCAR COMO VARIAVEL DE AMBIENTE
	conn, err := ampq.Dial("amqp://guest:guest@localhost:5672/")
	if err != nil {
		panic(err)
	}
	ch, err := conn.Channel()
	if err != nil {
		panic(err)
	}
	return ch
}
