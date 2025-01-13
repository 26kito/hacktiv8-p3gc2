package service

import (
	"bookservice/entity"
	pb "bookservice/proto"
	"bookservice/repository"
	"context"
	"errors"
	"fmt"
	"time"

	"github.com/golang-jwt/jwt/v5"
	"google.golang.org/grpc/status"
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

func (bs *BookService) DeleteBook(ctx context.Context, req *pb.GetBookByIdRequest) (*pb.Empty, error) {
	err := bs.BookRepository.DeleteBook(req.Id)
	if err != nil {
		return nil, err
	}

	return &pb.Empty{}, nil
}

func (bs *BookService) BorrowBook(ctx context.Context, req *pb.BorrowBookRequest) (*pb.BorrowBookResponse, error) {
	userId := ctx.Value("user").(jwt.MapClaims)["user_id"].(string)

	payload := entity.BorrowBookRequest{
		BookID:     req.BookId,
		BorrowDate: req.BorrowDate,
	}

	if err := validateBorrowBookPayload(req); err != nil {
		return nil, status.Error(400, err.Error())
	}

	res, err := bs.BookRepository.BorrowBook(userId, payload)
	if err != nil {
		return nil, status.Error(500, err.Error())
	}

	return &pb.BorrowBookResponse{Id: res.ID.Hex()}, nil
}

func (bs *BookService) ReturnBook(ctx context.Context, req *pb.ReturnBookRequest) (*pb.Empty, error) {
	userId := ctx.Value("user").(jwt.MapClaims)["user_id"].(string)

	payload := entity.ReturnBookRequest{
		BookID:     req.BookId,
		ReturnDate: req.ReturnDate,
	}

	err := bs.BookRepository.ReturnBook(userId, payload)
	if err != nil {
		return nil, status.Error(500, err.Error())
	}

	return &pb.Empty{}, nil
}

func (bs *BookService) UpdateBookStatus(ctx context.Context, req *pb.Empty) (*pb.Empty, error) {
	err := bs.BookRepository.UpdateBookStatus()
	if err != nil {
		return nil, status.Error(500, err.Error())
	}

	return &pb.Empty{}, nil
}

func validateBorrowBookPayload(payload *pb.BorrowBookRequest) error {
	if payload.BookId == "" {
		return fmt.Errorf("book_id is required")
	}

	if payload.BorrowDate == "" {
		return fmt.Errorf("borrow_date is required")
	}

	parsedDate, err := time.Parse("2006-01-02", payload.BorrowDate)
	if err != nil {
		return errors.New("invalid borrow_date format, expected YYYY-MM-DD")
	}

	// Get the current date (without time part)
	currentDate := time.Now().Truncate(24 * time.Hour)

	// Validate that borrow_date is not earlier than the current date
	if parsedDate.Before(currentDate) {
		return errors.New("borrow_date cannot be earlier than the current date")
	}

	return nil
}
