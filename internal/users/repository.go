package users

import ( 
	"context"
	"database/sql"
)

type Repository interface {
	CreateUser(ctx context.Context, user *Usuario) (int64, error) // Crea un nuevo usuario y devuelve su ID
	LoginUser(ctx context.Context, email, password string) (*Usuario, error ) // Obtiene un usuario por su correo electrónico y contraseña
	GetUserByID(ctx context.Context, id int64) (*Usuario, error) // Obtiene un usuario por su ID
	UpdateUserUSMPesos(ctx context.Context, userID int64, amount int64) error // Actualiza la cantidad de USM Pesos de un usuario
	ListUsers(ctx context.Context) ([]*Usuario, error) // Lista todos los usuarios
	GetUserByEmail(ctx context.Context, email string) (*Usuario, error) // Obtiene un usuario por su email
}

type sqliteRepository struct { // Implementación del repositorio utilizando SQLite
	db *sql.DB // Conexión a la base de datos
}

func NewSQLiteRepository(db *sql.DB) Repository { // Constructor para crear un nuevo repositorio SQLite
	return &sqliteRepository{db: db} // Retorna una instancia del repositorio SQLite
}

func (r *sqliteRepository) CreateUser(ctx context.Context, user *Usuario) (int64, error) {
	result, err := r.db.ExecContext(ctx, "INSERT INTO Usuario (first_name, last_name, email, password, usm_pesos) VALUES (?, ?, ?, ?, ?)",
		user.FirstName, user.LastName, user.Email, user.Password, user.USMPesos)
	if err != nil {
		return 0, err
	}
	return result.LastInsertId()
}

func (r *sqliteRepository) LoginUser(ctx context.Context, email, password string) (*Usuario, error) {
	row := r.db.QueryRowContext(ctx, "SELECT id, first_name, last_name, email, password, usm_pesos FROM Usuario WHERE email = ? AND password = ?", email, password)
	user := &Usuario{}
	err := row.Scan(&user.ID, &user.FirstName, &user.LastName, &user.Email, &user.Password, &user.USMPesos)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, nil // No se encontró el usuario
		}
		return nil, err
	}
	return user, nil
}

func (r *sqliteRepository) GetUserByEmail(ctx context.Context, email string) (*Usuario, error) {
	row := r.db.QueryRowContext(ctx, "SELECT id, first_name, last_name, email, password, usm_pesos FROM Usuario WHERE email = ?", email)
	user := &Usuario{}
	err := row.Scan(&user.ID, &user.FirstName, &user.LastName, &user.Email, &user.Password, &user.USMPesos)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, nil // No se encontró el usuario
		}
		return nil, err
	}
	return user, nil
}

func (r *sqliteRepository) GetUserByID(ctx context.Context, id int64) (*Usuario, error) {
	row := r.db.QueryRowContext(ctx, "SELECT id, first_name, last_name, email, password, usm_pesos FROM Usuario WHERE id = ?", id)
	user := &Usuario{}
	err := row.Scan(&user.ID, &user.FirstName, &user.LastName, &user.Email, &user.Password, &user.USMPesos)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, nil // No se encontró el usuario
		}
		return nil, err
	}
	return user, nil
}

func (r *sqliteRepository) UpdateUserUSMPesos(ctx context.Context, userID int64, amount int64) error {
	_, err := r.db.ExecContext(ctx, "UPDATE Usuario SET usm_pesos = usm_pesos + ? WHERE id = ?", amount, userID)
	return err
}

func (r *sqliteRepository) ListUsers(ctx context.Context) ([]*Usuario, error) {
	rows, err := r.db.QueryContext(ctx, "SELECT id, first_name, last_name, email, password, usm_pesos FROM Usuario")
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var users []*Usuario
	for rows.Next() {
		user := &Usuario{}
		err := rows.Scan(&user.ID, &user.FirstName, &user.LastName, &user.Email, &user.Password, &user.USMPesos)
		if err != nil {
			return nil, err
		}
		users = append(users, user)
	}
	if err = rows.Err(); err != nil {
		return nil, err
	}
	return users, nil
}

