# UZM Server – Proyecto Sistema Distribuidos
Benjamín Ferrada Larach | 202273061-7
Erick Jakín Ávila |
Renato Martínez Pierola | 

Este servidor implementa la API para gestionar **Usuarios, Libros, Inventario, Ventas y Préstamos**.  
Está desarrollado en **Go + Gin + SQLite** siguiendo arquitectura limpia (service, repository, handler).

---

## Ejecutar el servidor

Compila y corre el API:

```bash
go build -o bin/api.exe ./cmd/api
./bin/api.exe
```
O correr directamente el API: 

```bash
go run ./cmd/api
```

Compilar y correr el cliente:
```bash
go build -o bin/client.exe ./client
./bin/client.exe
```

Deberías ver en consola:

```
Servidor escuchando en http://localhost:8080
```

---

## Endpoints principales

### Health
```bash
curl http://localhost:8080/health
```
Respuesta:
```json
{"status":"ok"}
```

### Catálogo de libros (por defecto solo disponibles)
```bash
curl http://localhost:8080/api/books
```

### Catálogo incluyendo agotados
```bash
curl "http://localhost:8080/api/books?available=0"
```

### Obtener libro por ID
```bash
curl http://localhost:8080/api/books/1
```

### Crear libro
```bash
curl -X POST http://localhost:8080/api/books \
 -H "Content-Type: application/json" \
 -d '{
   "book_name": "El principito",
   "book_category": "Infantil",
   "transaction_type": "venta",
   "price": 10000,
   "status": true,
   "stock": 6
 }'
```

Respuesta esperada:
```json
{
  "id": 1,
  "book_name": "El principito",
  "book_category": "Infantil",
  "transaction_type": "venta",
  "price": 10000,
  "status": "Disponible",
  "popularity_score": 0,
  "inventory": { "available_quantity": 6 }
}
```

### Actualizar libro (PATCH)
```bash
curl -X PATCH http://localhost:8080/api/books/1 \
 -H "Content-Type: application/json" \
 -d '{"price":12000,"stock":0}'
```

Después, listar catálogo:
```bash
curl http://localhost:8080/api/books
```
→ El libro con stock 0 ya no aparece.

---

## Validaciones

- Transaction type inválido:
```bash
curl -X POST http://localhost:8080/api/books \
 -H "Content-Type: application/json" \
 -d '{"book_name":"X","book_category":"Y","transaction_type":"foo","price":1,"stock":1}'
```
Respuesta: **400 Bad Request**.

- Precio negativo:
```bash
curl -X POST http://localhost:8080/api/books \
 -H "Content-Type: application/json" \
 -d '{"book_name":"X","book_category":"Y","transaction_type":"venta","price":-1,"stock":1}'
```
Respuesta: **400 Bad Request**.

---

## Pruebas automáticas

### Unit tests
Ejecutar:
```bash
go test ./...
```

Ejemplos incluidos:
- `handler_test.go`: prueba endpoints con `httptest` y fake service.
- `repo_test.go`: prueba `CreateBook` y `GetBookByID` contra SQLite en disco (`test.db`).

---

## Seed de datos (opcional)

Archivo `seed.sql`:

```sql
DELETE FROM Inventario;
DELETE FROM Libro;

INSERT INTO Libro (id, book_name, book_category, transaction_type, price, status, popularity_score)
VALUES
(1,'1984','Ficción','venta',9000,1,0),
(2,'Rayuela','Ficción','arriendo',3000,1,0),
(3,'Base de Datos','Académico','venta',15000,1,0);

INSERT INTO Inventario (book_id, available_quantity)
VALUES
(1,5),(2,0),(3,2);
```

Ejecutar:
```bash
sqlite3 uzm.db ".read seed.sql"
```

---

## Guión rápido de pruebas (manual)

1. Compilar y correr:
   ```bash
   go build -o bin/api ./cmd/api
   ./bin/api
   ```

2. Health:
   ```bash
   curl http://localhost:8080/health
   ```

3. Crear libro:
   ```bash
   curl -X POST http://localhost:8080/api/books -d '{"book_name":"Prueba","book_category":"Test","transaction_type":"venta","price":5000,"stock":2}' -H "Content-Type: application/json"
   ```

4. Listar catálogo:
   ```bash
   curl http://localhost:8080/api/books
   ```

5. Actualizar stock a 0:
   ```bash
   curl -X PATCH http://localhost:8080/api/books/1 -d '{"stock":0}' -H "Content-Type: application/json"
   curl http://localhost:8080/api/books   # el libro ya no aparece
   ```

---

## Notas

- Usa **Go 1.25.1** o superior.  
- La base de datos es **SQLite** (`uzm.db` por defecto).  
- El esquema se migra automáticamente al iniciar (`MustMigrate`).  
- Todos los endpoints están bajo el prefijo `/api`.
