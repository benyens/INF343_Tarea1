CREATE TABLE IF NOT EXISTS Usuario (
	id INTEGER PRIMARY KEY AUTOINCREMENT,
	first_name TEXT NOT NULL,
	last_name TEXT NOT NULL UNIQUE,
	email TEXT NOT NULL UNIQUE,
	password TEXT NOT NULL,
	usm_pesos INTEGER DEFAULT 0
);

CREATE TABLE IF NOT EXISTS Libro (
    id INTEGER PRIMARY KEY AUTOINCREMENT,
    book_name TEXT NOT NULL,
    book_category TEXT NOT NULL,
    transaction_type TEXT NOT NULL,
    price INTEGER NOT NULL,
    status BOOLEAN NOT NULL DEFAULT 1,
    popularity_score INTEGER DEFAULT 0
);

CREATE TABLE IF NOT EXISTS Inventario (
    book_id INTEGER PRIMARY KEY,
    available_quantity INTEGER NOT NULL,
    FOREIGN KEY (book_id) REFERENCES Libro(id)
); 

CREATE TABLE IF NOT EXISTS Prestamo (
    id INTEGER PRIMARY KEY AUTOINCREMENT,
    user_id INTEGER,
    book_id INTEGER,
    start_date TEXT NOT NULL,
    return_date TEXT NOT NULL,
    status TEXT NOT NULL CHECK (status IN ('pendiente','finalizado')),
    FOREIGN KEY (user_id) REFERENCES Usuario(id),
    FOREIGN KEY (book_id) REFERENCES Libro(id)
);

CREATE TABLE IF NOT EXISTS Venta (
    id INTEGER PRIMARY KEY AUTOINCREMENT,
    user_id INTEGER,
    book_id INTEGER,
    sale_date DATE NOT NULL,
    FOREIGN KEY (user_id) REFERENCES Usuario(id),
    FOREIGN KEY (book_id) REFERENCES Libro(id)
)