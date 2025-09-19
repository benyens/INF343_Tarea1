package books

import (
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
)

type bookResponse struct {
	Id         int64   `json:"id"`
	BookName   string  `json:"book_name"`
	BookCategory string `json:"book_category"`
	TransactionType string `json:"transaction_type"`
	Price      int64 `json:"price"`
	Status     bool    `json:"status"`
	PopularityScore int64 `json:"popularity_score"`
	Inventory       struct {
        AvailableQuantity int64 `json:"available_quantity"`
    } `json:"inventory"`
}

func toBookResponse(book *BookWithInventory) *bookResponse {
	return &bookResponse{
		Id:         book.ID,
		BookName:   book.BookName,
		BookCategory: book.BookCategory,
		TransactionType: book.TransactionType,
		Price:      book.Price,
		Status:     book.Status,
		PopularityScore: book.PopularityScore,
		Inventory: struct {
			AvailableQuantity int64 `json:"available_quantity"`
		}{
			AvailableQuantity: book.AvailableQuantity,
		},
	}

}

type Handler struct {
	service Service
}

func NewHandler(service Service) *Handler {
	return &Handler{service: service}
}

func (h *Handler) RegisterRoutes(r *gin.Engine) {
	r.GET("/books", h.ListBooks)
	r.GET("/books/:id", h.getBookByID)
	r.POST("/books", h.createBook)
	r.PATCH("/books/:id", h.updateBook)
}

func (h *Handler) createBook(c *gin.Context) {
	var input CreateBookInput
	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request payload"})
		return
	}
	bookID, err := h.service.CreateBook(c.Request.Context(), input)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	book, err := h.service.GetBookByID(c.Request.Context(), bookID, nil)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusCreated, toBookResponse(book))
}

func (h *Handler) ListBooks(c *gin.Context) {
	statusParam := c.Query("status")
	var status *bool
	if statusParam != "" {
		parsedStatus, err := strconv.ParseBool(statusParam)
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid status parameter"})
			return
		}
		status = &parsedStatus
	}
	books, err := h.service.ListBook(c.Request.Context(), status)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, books)
}

func (h *Handler) getBookByID(c *gin.Context) {
	idParam := c.Param("id")
	id, err := strconv.ParseInt(idParam, 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid book ID"})
		return
	}
	book, err := h.service.GetBookByID(c.Request.Context(), id, nil)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, toBookResponse(book))
}
func (h *Handler) updateBook(c *gin.Context) {
	idParam := c.Param("id")
	id, err := strconv.ParseInt(idParam, 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid book ID"})
		return
	}
	var input UpdateBookInput
	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request payload"})
		return
	}
	book, err := h.service.UpdateBook(c.Request.Context(), id, input)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, toBookResponse(book))
}
