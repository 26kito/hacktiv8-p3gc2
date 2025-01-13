package service

import (
	"bookservice/entity"
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

func (bs *BookService) InsertBook(ctx context.Context, req *pb.InsertBookRequest) (*pb.InsertBookResponse, error) {
	payload := entity.InsertBookRequest{
		Title:         req.Title,
		Author:        req.Author,
		PublishedDate: req.PublishedDate,
		Status:        req.Status,
	}

	res, err := bs.BookRepository.InsertBook(payload)
	if err != nil {
		return nil, err
	}

	return &pb.InsertBookResponse{Id: res.ID.Hex()}, nil
}

func (bs *BookService) GetBookById(ctx context.Context, req *pb.GetBookByIdRequest) (*pb.GetBookResponse, error) {
	book, err := bs.BookRepository.GetBookById(req.Id)
	if err != nil {
		return nil, err
	}

	return &pb.GetBookResponse{Book: &pb.Book{
		Id:            book.ID.Hex(),
		Title:         book.Title,
		Author:        book.Author,
		PublishedDate: book.PublishedDate,
		Status:        book.Status,
	}}, nil
}

func (bs *BookService) UpdateBook(ctx context.Context, req *pb.UpdateBookRequest) (*pb.UpdateBookResponse, error) {
	payload := entity.UpdateBookRequest{
		ID:            req.Id,
		Title:         req.Title,
		Author:        req.Author,
		PublishedDate: req.PublishedDate,
		Status:        req.Status,
	}

	res, err := bs.BookRepository.UpdateBook(payload)
	if err != nil {
		return nil, err
	}

	return &pb.UpdateBookResponse{Id: res.ID.Hex()}, nil
}
