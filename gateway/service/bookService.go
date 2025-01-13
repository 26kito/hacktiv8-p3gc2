package service

import (
	"context"
	"gateway/entity"
	pb "gateway/proto/book"

	"github.com/labstack/echo/v4"
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
