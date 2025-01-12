package service

import (
	pb "bookservice/proto"
	"bookservice/repository"
	"context"
)

type BookService struct {
	BookRepository repository.BookRepository
	pb.UnimplementedBookServiceServer
}

func NewBookService(bookRepository repository.BookRepository) *BookService {
	return &BookService{BookRepository: bookRepository}
}

func (bs *BookService) GetAllBook(ctx context.Context, req *pb.Empty) (*pb.GetAllBookResponse, error) {
	books, err := bs.BookRepository.GetAllBook()
	if err != nil {
		return nil, err
	}

	// Convert the books to the gRPC response format
	var pbBooks []*pb.Book
	for _, book := range books {
		pbBooks = append(pbBooks, &pb.Book{
			Id:            book.ID.Hex(),
			Title:         book.Title,
			Author:        book.Author,
			PublishedDate: book.PublishedDate,
			Status:        book.Status,
		})
	}

	return &pb.GetAllBookResponse{Books: pbBooks}, nil
}
