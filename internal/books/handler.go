package books

import (
	"net/http"
	"strconv"
	"log"
	"io"
	"bytes"
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

func toBookResponse(bwi *BookWithInventory) *bookResponse {
    b := bwi.Book

    // Si prefieres status booleano desde DB, usa: status := b.Status
    resp := &bookResponse{
        Id:              b.ID,
        BookName:        b.BookName,
        BookCategory:    b.BookCategory,
        TransactionType: b.TransactionType,
        Price:           b.Price,
        Status:          b.Status,
        PopularityScore: b.PopularityScore,
    }
    resp.Inventory.AvailableQuantity = bwi.AvailableQuantity
    return resp
}

type Handler struct {
	service Service
}

func NewHandler(service Service) *Handler {
	return &Handler{service: service}
}

func (h *Handler) RegisterRoutes(r *gin.RouterGroup) {
	r.GET("/books", h.ListBooks)
	r.GET("/books/:id", h.getBookByID)
	r.POST("/books", h.createBook)
	r.PATCH("/books/:id", h.updateBook)
}

func (h *Handler) createBook(c *gin.Context) {
    // Diagnóstico: asegúrate del content-type
    ct := c.GetHeader("Content-Type")
    log.Println("Content-Type:", ct)

    // Si quieres loguear el body sin romper el bind, hay que reponerlo:
    b, _ := io.ReadAll(c.Request.Body)
    log.Println("RAW BODY:", string(b))
    c.Request.Body = io.NopCloser(bytes.NewBuffer(b))

    var req CreateBookInput
    if err := c.ShouldBindJSON(&req); err != nil {
        c.JSON(http.StatusBadRequest, gin.H{"error": "JSON inválido: " + err.Error()})
        return
    }
    log.Printf("REQ (after bind): %#v", req)

	id, err := h.service.CreateBook(c.Request.Context(), req)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	bwi, _ := h.service.GetBookByID(c.Request.Context(), id, nil)
	c.JSON(http.StatusCreated, toBookResponse(bwi))
}


func (h *Handler) ListBooks(c *gin.Context) {
	statusParam := c.Query("status")
	var onlyAvailable *bool
	if statusParam != "" {
		parsed, err := strconv.ParseBool(statusParam)
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid status parameter"})
			return
		}
		onlyAvailable = &parsed
	}

	books, err := h.service.ListBook(c.Request.Context(), onlyAvailable)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	out := make([]*bookResponse, 0, len(books))
	for _, bwi := range books {
		out = append(out, toBookResponse(bwi))
	}
	c.JSON(http.StatusOK, gin.H{"books": out})
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
