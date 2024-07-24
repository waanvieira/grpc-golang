package main

import (
	"database/sql"
	"net"

	_ "github.com/mattn/go-sqlite3"
	"github.com/waanvieira/grpc-go/internal/database"
	"github.com/waanvieira/grpc-go/internal/pb"
	"github.com/waanvieira/grpc-go/service"
	"google.golang.org/grpc"
	"google.golang.org/grpc/reflection"
)

func main() {
	db, err := sql.Open("sqlite3", "./db.sqlite")
	if err != nil {
		panic(err)
	}

	defer db.Close()
	// Iniciamos a nossa classe passando o db como construtor, com isso teremos acesso a todos os
	// Metodos do banco na nossa classe e conexão com o banco
	categoryDB := database.NewCategory(db)
	// Injetamos o categoryDB no nosso seriço para ter acesso a todos os métodos, nesse caso criar uma nova categoria
	categoryService := service.NewCategoryService(*categoryDB)

	// iniciamos o nosso servidor grpc
	grpcServer := grpc.NewServer()
	// Aqui estamos ataxando, relacionando o nosso serviço com o nosso service
	pb.RegisterCategoryServiceServer(grpcServer, categoryService)
	// Iniciando o reflection para porder usar o servidor GRPC Evans https://github.com/ktr0731/evans
	// go install github.com/ktr0731/evans@latest
	// Iniciando o reflection, reflection o que consegue ler e processar a sua própria informação
	reflection.Register(grpcServer)

	// Abrindo uma porta tcp para fazer a nossa comunicação
	lis, err := net.Listen("tcp", ":50051")
	if err != nil {
		panic(err)
	}

	// Iniciando o servidor grpc na nossa porta tcp 50051
	if err := grpcServer.Serve(lis); err != nil {
		panic(err)
	}
}
