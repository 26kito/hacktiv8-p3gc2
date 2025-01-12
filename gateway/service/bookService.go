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
