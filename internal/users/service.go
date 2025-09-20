package users

import "context"

type Usuario struct {
	ID    int64
	FirstName string
	LastName  string
	Email string
	Password string
	USMPesos int64
}

type Service interface {
	RegisterUser(ctx context.Context, user *Usuario) (int64, error) // Registra un nuevo usuario y devuelve su ID
	LoginUser(ctx context.Context, email, password string) (*Usuario, error) // Autentica a un usuario y devuelve su información
	GetUserByID(ctx context.Context, id int64) (*Usuario, error) // Obtiene un usuario por su ID
	UpdateUserUSMPesos(ctx context.Context, userID int64, amount int64) error // Actualiza la cantidad de USM Pesos de un usuario
	ListUsers(ctx context.Context) ([]*Usuario, error) // Lista todos los usuarios
	GetUserByEmail(ctx context.Context, email string) (*Usuario, error) // Obtiene un usuario por su email
}

type service struct {
	repo Repository
}

func NewService(repo Repository) Service {
	return &service{repo: repo}
}

func (s *service) RegisterUser(ctx context.Context, user *Usuario) (int64, error) {
	return s.repo.CreateUser(ctx, user)
}

func (s *service) LoginUser(ctx context.Context, email, password string) (*Usuario, error) {
	return s.repo.LoginUser(ctx, email, password)
}

func (s *service) GetUserByID(ctx context.Context, id int64) (*Usuario, error) {
	return s.repo.GetUserByID(ctx, id)
}

func (s *service) UpdateUserUSMPesos(ctx context.Context, userID int64, amount int64) error {
	return s.repo.UpdateUserUSMPesos(ctx, userID, amount)
}

func (s *service) ListUsers(ctx context.Context) ([]*Usuario, error) {
	return s.repo.ListUsers(ctx)
}
func (s *service) GetUserByEmail(ctx context.Context, email string) (*Usuario, error) {
	return s.repo.GetUserByEmail(ctx, email)
}



