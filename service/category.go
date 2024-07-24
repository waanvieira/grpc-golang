package service

import (
	"context"
	"io"

	"github.com/waanvieira/grpc-go/internal/database"
	"github.com/waanvieira/grpc-go/internal/pb"
)

type CategoryService struct {
	// Aqui é o serivço que temos que utilizar do protocol buffer
	pb.UnimplementedCategoryServiceServer
	// Nossa classe do DB que vamos acessar para ter acesso ao DB
	CategoryDB database.Category
}

// Basicamente nosso método construtor, que vai receber o categoryDB na consrtução para retornar a nossa struct (classe)
func NewCategoryService(categoryDB database.Category) *CategoryService {
	return &CategoryService{
		CategoryDB: categoryDB,
	}
}

// Função para gerar o nosso GRPC com a request e a response, essa função já é pré feita, implementamos a interface
// do nosso arquivo course_category_grpc.pb
// c *CategoryService - definimos que esse método pertence a essa struct
// ctx - recebemos um context
// input - recebemos a nossa entrada com o formato criado na nossa proto/course_category.proto
// Retornamos o nosso objeto - CategoryResponse também criado no nosso  proto/course_category.proto
// func (c *CategoryService) CreateCategory(ctx context.Context, input *pb.CreateCategoryRequest) (*pb.CategoryResponse, error) {
// 	// Criamos a nossa categoria chamando o nosso método create da nossa classe de database
// 	// Basicamente a mesma coisa que uma api só que com entrada de dados diferente e diferente modo de saida
// 	category, err := c.CategoryDB.Create(input.Name, input.Description)
// 	if err != nil {
// 		return nil, err
// 	}
// 	// Hidratamos a nossa struct de category do nosso pb com os dados retornados do banco
// 	// Para retornar no formato da nossa struc definida do nosso course_category
// 	categoryResponse := &pb.Category{
// 		Id:          category.ID,
// 		Name:        category.Name,
// 		Description: category.Description,
// 	}
// 	// Nesse caso estamos retornando o CategoryResponse que retorna uma categoria,
// 	// Então injetamos a categoria com o CategoryReponse
// 	// Nesse caso também podemos retornar com o formato de Category também
// 	// Abaixo deixar um exemplo de como poderia utilizar também retornando o category
// 	return &pb.CategoryResponse{
// 		Category: categoryResponse,
// 	}, nil

// deixando dessa forma o retorno seria esse
// {
// 	"category": {
// 	  "description": "descrição test",
// 	  "id": "8cf1850f-cd5f-4447-92e7-0139ac78fe87",
// 	  "name": "categoria teste"
// 	}
//   }
// }

// outro exemplo de uso
func (c *CategoryService) CreateCategory(ctx context.Context, input *pb.CreateCategoryRequest) (*pb.Category, error) {
	category, err := c.CategoryDB.Create(input.Name, input.Description)
	if err != nil {
		return nil, err
	}
	categoryResponse := &pb.Category{
		Id:          category.ID,
		Name:        category.Name,
		Description: category.Description,
	}
	return categoryResponse, nil
	// O retorno desse é sse aqui
	// {
	// 	"description": "description",
	// 	"id": "c86ceec3-142d-4ab0-b89e-1e357210b4ee",
	// 	"name": "Name"
	//   }
}

func (c *CategoryService) ListCategories(ctx context.Context, input *pb.Blank) (*pb.CategoryList, error) {
	categories, err := c.CategoryDB.GetAll()
	if err != nil {
		return nil, err
	}

	var categoriesReponse []*pb.Category

	for _, category := range categories {
		categoryResponse := &pb.Category{
			Id:          category.ID,
			Name:        category.Name,
			Description: category.Description,
		}

		categoriesReponse = append(categoriesReponse, categoryResponse)
	}

	return &pb.CategoryList{Categories: categoriesReponse}, nil
}

// O nome do método tem que ser igual o nome da message declarado no proto para poder fazer o bind
func (c *CategoryService) ListCategoryByID(ctx context.Context, input *pb.FindCategoryByID) (*pb.Category, error) {
	category, err := c.CategoryDB.GetCategoryByID(input.Id)
	if err != nil {
		return nil, err
	}
	return &pb.Category{
		Id:          category.ID,
		Name:        category.Name,
		Description: category.Description,
	}, nil
}

func (c *CategoryService) CreateCategoryStream(stream pb.CategoryService_CreateCategoryStreamServer) error {
	categories := &pb.CategoryList{}
	// Criando um loop infinito, para o loop ficar responsável em ficar enviando meus dados o tempo todo
	for {
		// Aqui começamos a receber os dados, a strem de dados
		category, err := stream.Recv()
		// Verifica se chegou ao final e não tem mais nada para enviar, aqui envia todas os meus dados no retorno
		if err == io.EOF {
			return stream.SendAndClose(categories)
		}
		if err != nil {
			return err
		}
		// Fazemos a inserção dos dados no banco
		categoryResult, err := c.CategoryDB.Create(category.Name, category.Description)

		if err != nil {
			return err
		}
		// Populamos a lista vazia que criamos, categories, ao final iremos hidratar essa variavel e retorna-la ao final da execução
		categories.Categories = append(categories.Categories, &pb.Category{
			Id:          categoryResult.ID,
			Name:        categoryResult.Name,
			Description: categoryResult.Description,
		})
	}
}

func (c *CategoryService) CreateCategoryStreamBidirectional(stream pb.CategoryService_CreateCategoryStreamBidirectionalServer) error {
	for {
		// Aqui recebemos os dados enviados pelo grpc
		dataReceivedViaStreamRecv, err := stream.Recv()
		if err == io.EOF {
			return nil
		}

		if err != nil {
			return err
		}

		// Criamos a categoria com os dados recebidos pelo grpc
		categoryCreatedInDB, err := c.CategoryDB.Create(dataReceivedViaStreamRecv.Name, dataReceivedViaStreamRecv.Description)
		if err != nil {
			return err
		}

		err = stream.Send(&pb.Category{
			Id:          categoryCreatedInDB.ID,
			Name:        categoryCreatedInDB.Name,
			Description: categoryCreatedInDB.Description,
		})

		if err != nil {
			return err
		}

	}
}
