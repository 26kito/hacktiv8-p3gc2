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

// @Summary Get all books
// @Description Get all books
// @Tags books
// @ID get-all-books
// @Accept json
// @Produce json
// @Success 200 {object} entity.ResponseOK
// @Failure 400 {object} entity.ResponseError
// @Router /books [get]
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

// @Summary Insert a new book
// @Description Insert a new book
// @Tags books
// @ID insert-book
// @Accept json
// @Produce json
// @Param book body entity.InsertBookRequest true "Book data"
// @Success 200 {object} entity.ResponseOK
// @Failure 400 {object} entity.ResponseError
// @Router /books [post]
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

// @Summary Get book by ID
// @Description Get book by ID
// @Tags books
// @ID get-book-by-id
// @Accept json
// @Produce json
// @Param id path string true "Book ID"
// @Success 200 {object} entity.ResponseOK
// @Failure 400 {object} entity.ResponseError
// @Router /books/{id} [get]
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

// @Summary Update a book
// @Description Update a book
// @Tags books
// @ID update-book
// @Accept json
// @Produce json
// @Param id path string true "Book ID"
// @Param book body entity.UpdateBookRequest true "Book data"
// @Success 200 {object} entity.ResponseOK
// @Failure 400 {object} entity.ResponseError
// @Router /books/{id} [put]
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

// @Summary Delete a book
// @Description Delete a book
// @Tags books
// @ID delete-book
// @Accept json
// @Produce json
// @Param id path string true "Book ID"
// @Success 200 {object} entity.ResponseOK
// @Failure 400 {object} entity.ResponseError
// @Router /books/{id} [delete]
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

// @Summary Borrow a book
// @Description Borrow a book
// @Tags books
// @ID borrow-book
// @Accept json
// @Produce json
// @Security ApiKeyAuth
// @Param id path string true "Book ID"
// @Param borrow_date body entity.BorrowBookRequest true "Borrow data"
// @Success 200 {object} entity.ResponseOK
// @Failure 400 {object} entity.ResponseError
// @Router /books/{id}/borrow [post]
func (bs *BookService) BorrowBook(c echo.Context) error {
	var payload entity.BorrowBookRequest
	bookId := c.Param("id")
	token := c.Request().Header.Get("Authorization")

	md := metadata.Pairs("authorization", token)
	ctx := metadata.NewOutgoingContext(context.Background(), md)

	if err := c.Bind(&payload); err != nil {
		return c.JSON(400, entity.ResponseError{
			Message: err.Error(),
		})
	}

	req := &pb.BorrowBookRequest{
		BookId:     bookId,
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

// @Summary Return a book
// @Description Return a book
// @Tags books
// @ID return-book
// @Accept json
// @Produce json
// @Security ApiKeyAuth
// @Param id path string true "Book ID"
// @Param return_date body entity.ReturnBookRequest true "Return data"
// @Success 200 {object} entity.ResponseOK
// @Failure 400 {object} entity.ResponseError
// @Router /books/{id}/return [post]
func (bs *BookService) ReturnBook(c echo.Context) error {
	var payload entity.ReturnBookRequest
	bookId := c.Param("id")
	token := c.Request().Header.Get("Authorization")

	md := metadata.Pairs("authorization", token)
	ctx := metadata.NewOutgoingContext(context.Background(), md)

	if err := c.Bind(&payload); err != nil {
		return c.JSON(400, entity.ResponseError{
			Message: err.Error(),
		})
	}

	req := &pb.ReturnBookRequest{
		BookId:     bookId,
		ReturnDate: payload.ReturnDate,
	}

	_, err := bs.BookClient.ReturnBook(ctx, req)
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
