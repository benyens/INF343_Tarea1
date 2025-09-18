CREATE TABLE IF NOT EXISTS Usuario (
	id INTEGER PRIMARY KEY AUTOINCREMENT,
	first_name TEXT NOT NULL,
	last_name TEXT NOT NULL UNIQUE,
	email TEXT NOT NULL UNIQUE,
	password TEXT NOT NULL,
	usm_pesos REAL DEFAULT 0
);

CREATE TABLE IF NOT EXISTS Libro (
    id INTEGER PRIMARY KEY AUTOINCREMENT,
    book_name TEXT NOT NULL,
    book_category TEXT NOT NULL,
    transaction_type TEXT NOT NULL,
    price REAL NOT NULL,
    status BOOLEAN DEFAULT 0,
    popularity_score INTEGER DEFAULT 0
);

CREATE TABLE IF NOT EXISTS Inventario (
    id INTEGER PRIMARY KEY AUTOINCREMENT,
    available_quantity INTEGER NOT NULL
); 

CREATE TABLE IF NOT EXISTS Prestamo (
    id INTEGER PRIMARY KEY AUTOINCREMENT,
    user_id INTEGER,
    book_id INTEGER,
    start_date DATE NOT NULL,
    return_date DATE NOT NULL,
    status TEXT NOT NULL,
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