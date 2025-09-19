package books

import (
	"context"
	"errors"
	"strings"
)

type Book struct {
    ID              int64
    BookName        string
    BookCategory    string
    TransactionType string
    Price           int64
    Status          bool
    PopularityScore int64
}

type BookWithInventory struct {
    Book              *Book
    AvailableQuantity int64
}

type CreateBookInput struct { // Estructura para la creación de un libro
	BookName        string
	BookCategory    string
	TransactionType string
	Price           int64
	Status          bool
	PopularityScore int64
	Stock           int64
}

type UpdateBookInput struct { // Estructura para la actualización de un libro
	BookName        *string
	BookCategory    *string
	TransactionType  *string
	Price            *int64
	Status           *bool
	PopularityScore  *int64
	Stock            *int64
}

type Service interface { // Interfaz del servicio de libros
	ListBook(ctx context.Context, onlyAvailable *bool) ([]*BookWithInventory, error)                    // Lista todos los libros, opcionalmente filtrados por estado
	GetBookByID(ctx context.Context, id int64, onlyAvailable *bool) (*BookWithInventory, error)         // Obtiene un libro por su ID, opcionalmente filtrado por estado
	CreateBook(ctx context.Context, input CreateBookInput) (int64, error)                        // Crea un nuevo libro
	UpdateBook(ctx context.Context, id int64, input UpdateBookInput) (*BookWithInventory, error) // Actualiza un libro existente
}

type Repository interface {
    ListBook(ctx context.Context, onlyAvailable *bool) ([]BookWithInventory, error)
    GetBookByID(ctx context.Context, id int64, onlyAvailable *bool) (BookWithInventory, error)
    CreateBook(ctx context.Context, book *Book, initialStock int64) (int64, error)
    UpdateBook(ctx context.Context, book *Book, stock *int64) error
}

type service struct { // Implementación del servicio de libros
	repo Repository // Repositorio para la gestión de libros
}

func NewService(repo Repository) Service { // Constructor para crear un nuevo servicio de libros
	return &service{repo: repo} // Retorna una instancia del servicio con el repositorio inyectado
}

func (s *service) ListBook(ctx context.Context, onlyAvailable *bool) ([]*BookWithInventory, error) {
	books, err := s.repo.ListBook(ctx, onlyAvailable)
	if err != nil {
		return nil, err
	}
	result := make([]*BookWithInventory, len(books))
	for i := range books {
		result[i] = &books[i]
	}
	return result, nil
}

func (s *service) GetBookByID(ctx context.Context, id int64, onlyAvailable *bool) (*BookWithInventory, error) {
	bookWithInventory, err := s.repo.GetBookByID(ctx, id, onlyAvailable)
	if err != nil {
		return nil, err
	}
	return &bookWithInventory, nil
}

func (s *service) CreateBook(ctx context.Context, input CreateBookInput) (int64, error) {
	tt := strings.ToLower(input.TransactionType)
	if tt != "venta" && tt != "arriendo" {
		return 0, errors.New("el tipo de transacción debe ser 'venta' o 'arriendo'") // Valida que el tipo de transacción sea válido
	}

	if input.Price < 0 {
		return 0, errors.New("el precio debe ser un valor positivo") // Valida que el precio sea un valor positivo
	}

	if strings.TrimSpace(input.BookName) == "" {
		return 0, errors.New("se necesita el nombre del libro") // Valida que el nombre del libro no esté vacío
	}

	book := &Book{
		BookName:        input.BookName,
		BookCategory:    input.BookCategory,
		TransactionType: tt,
		Price:           input.Price,
		Status:          input.Status,
		PopularityScore: input.PopularityScore,
	}

	// Crear el libro en la base de datos
	id, err := s.repo.CreateBook(ctx, book, input.Stock)
	if err != nil {
		return 0, err
	}
	book.ID = id

	return id, nil
}

func (s *service) UpdateBook(ctx context.Context, id int64, input UpdateBookInput) (*BookWithInventory, error) {
	book, err := s.repo.GetBookByID(ctx, id, nil)
	if err != nil {
		return nil, err
	}

	if input.BookName != nil {
		book.Book.BookName = *input.BookName
	}
	if input.BookCategory != nil {
		book.Book.BookCategory = *input.BookCategory
	}
	if input.TransactionType != nil {
		book.Book.TransactionType = *input.TransactionType
	}
	if input.Price != nil {
		book.Book.Price = *input.Price
	}
	if input.Status != nil {
		book.Book.Status = *input.Status
	}
	if input.PopularityScore != nil {
		book.Book.PopularityScore = *input.PopularityScore
	}

	err = s.repo.UpdateBook(ctx, book.Book, input.Stock)
	if err != nil {
		return nil, err
	}

	return &book, nil
}
