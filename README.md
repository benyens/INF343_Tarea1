# 1 Correr el server 
```bash
go run ./cmd/api 
```
# 2 Probar "salud" del server
```bash
curl http://localhost:8080/health
```
# 3 Pruebas
## Crear Libro
curl -s -X POST http://localhost:8080/api/v1/books \
-H "Content-Type: application/json" \
-d '{"book_name": "Al sur de la frontera, al oeste del sol", "book_category": "Documental y road movie político",
"transaction_type": "venta", 
"price": 10,
"status": "Disponible",
"popularity_score": 4,
"inventory": {
    "available_quantity": 6
    }
}