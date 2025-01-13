package service

import (
	"context"
	"gateway/entity"
	pb "gateway/proto/book"

	"github.com/labstack/echo/v4"
	"google.golang.org/grpc/metadata"
	"google.golang.org/grpc/status"
)

type BookService struct {
	BookClient pb.BookServiceClient
}

func (bs *BookService) GetAllBook(c echo.Context) error {
	resp, err := bs.BookClient.GetAllBook(context.Background(), &pb.Empty{})
	if err != nil {
		return c.JSON(400, entity.ResponseError{
			Message: err.Error(),
		})
	}

	return c.JSON(200, entity.ResponseOK{
		Message: "Success",
		Data:    resp.Books,
	})
}

func (bs *BookService) InsertBook(c echo.Context) error {
	var payload entity.InsertBookRequest

	if err := c.Bind(&payload); err != nil {
		return c.JSON(400, entity.ResponseError{
			Message: err.Error(),
		})
	}

	req := &pb.InsertBookRequest{
		Title:         payload.Title,
		Author:        payload.Author,
		PublishedDate: payload.PublishedDate,
		Status:        payload.Status,
	}

	resp, err := bs.BookClient.InsertBook(context.Background(), req)
	if err != nil {
		return c.JSON(400, entity.ResponseError{
			Message: err.Error(),
		})
	}

	return c.JSON(200, entity.ResponseOK{
		Message: "Book created successfully",
		Data:    resp,
	})
}

func (bs *BookService) GetBookById(c echo.Context) error {
	id := c.Param("id")
	req := &pb.GetBookByIdRequest{
		Id: id,
	}

	resp, err := bs.BookClient.GetBookById(context.Background(), req)
	if err != nil {
		return c.JSON(400, entity.ResponseError{
			Message: err.Error(),
		})
	}

	return c.JSON(200, entity.ResponseOK{
		Message: "Success",
		Data:    resp,
	})
}

func (bs *BookService) UpdateBook(c echo.Context) error {
	id := c.Param("id")
	var payload entity.UpdateBookRequest

	if err := c.Bind(&payload); err != nil {
		return c.JSON(400, entity.ResponseError{
			Message: err.Error(),
		})
	}

	req := &pb.UpdateBookRequest{
		Id:            id,
		Title:         payload.Title,
		Author:        payload.Author,
		PublishedDate: payload.PublishedDate,
		Status:        payload.Status,
	}

	_, err := bs.BookClient.UpdateBook(context.Background(), req)

	if err != nil {
		return c.JSON(400, entity.ResponseError{
			Message: err.Error(),
		})
	}

	return c.JSON(200, entity.ResponseOK{
		Message: "Book updated successfully",
		Data:    nil,
	})
}

func (bs *BookService) DeleteBook(c echo.Context) error {
	id := c.Param("id")
	req := &pb.GetBookByIdRequest{
		Id: id,
	}

	_, err := bs.BookClient.DeleteBook(context.Background(), req)
	if err != nil {
		return c.JSON(400, entity.ResponseError{
			Message: err.Error(),
		})
	}

	return c.JSON(200, entity.ResponseOK{
		Message: "Book deleted successfully",
		Data:    nil,
	})
}

func (bs *BookService) BorrowBook(c echo.Context) error {
	var payload entity.BorrowBookRequest
	id := c.Param("id")
	token := c.Request().Header.Get("Authorization")

	md := metadata.Pairs("authorization", token)
	ctx := metadata.NewOutgoingContext(context.Background(), md)

	if err := c.Bind(&payload); err != nil {
		return c.JSON(400, entity.ResponseError{
			Message: err.Error(),
		})
	}

	req := &pb.BorrowBookRequest{
		BookId:     id,
		UserId:     payload.UserID,
		BorrowDate: payload.BorrowDate,
	}

	res, err := bs.BookClient.BorrowBook(ctx, req)
	if err != nil {
		st, _ := status.FromError(err)
		errMessage := st.Message()

		return c.JSON(400, entity.ResponseError{
			Message: errMessage,
		})
	}

	return c.JSON(200, entity.ResponseOK{
		Message: "Book borrowed successfully",
		Data:    res.Id,
	})
}

func (bs *BookService) ReturnBook(c echo.Context) error {
	id := c.Param("id")
	var payload entity.ReturnBookRequest

	if err := c.Bind(&payload); err != nil {
		return c.JSON(400, entity.ResponseError{
			Message: err.Error(),
		})
	}

	req := &pb.ReturnBookRequest{
		Id:         id,
		ReturnDate: payload.ReturnDate,
	}

	_, err := bs.BookClient.ReturnBook(context.Background(), req)
	if err != nil {
		st, _ := status.FromError(err)
		errMessage := st.Message()

		return c.JSON(400, entity.ResponseError{
			Message: errMessage,
		})
	}

	return c.JSON(200, entity.ResponseOK{
		Message: "Book returned successfully",
		Data:    nil,
	})
}
