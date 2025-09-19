package books

import (
    "context"
    "database/sql"
    "errors"
    "strings"
)

type sqliteRepository struct{ dbconn *sql.DB }

func NewSQLiteRepository(db *sql.DB) Repository { return &sqliteRepository{dbconn: db} }

func (r *sqliteRepository) ListBook(ctx context.Context, onlyAvailable *bool) ([]BookWithInventory, error) {
    q := `
SELECT  b.id, b.book_name, b.book_category, b.transaction_type,
        b.price, b.status, b.popularity_score,
        COALESCE(i.available_quantity, 0) AS qty
FROM    Libro b
LEFT JOIN Inventario i ON i.book_id = b.id`
    args := []any{}
    if onlyAvailable != nil && *onlyAvailable {
        q += " WHERE COALESCE(i.available_quantity,0) > 0"
    }
    q += " ORDER BY b.id ASC"

    rows, err := r.dbconn.QueryContext(ctx, q, args...)
    if err != nil { return nil, err }
    defer rows.Close()

    var out []BookWithInventory
    for rows.Next() {
        var b Book
        var qty int64
        if err := rows.Scan(&b.ID, &b.BookName, &b.BookCategory, &b.TransactionType,
            &b.Price, &b.Status, &b.PopularityScore, &qty); err != nil {
            return nil, err
        }
        out = append(out, BookWithInventory{Book: &b, AvailableQuantity: qty})
    }
    return out, rows.Err()
}

func (r *sqliteRepository) GetBookByID(ctx context.Context, id int64, onlyAvailable *bool) (BookWithInventory, error) {
    q := `
SELECT  b.id, b.book_name, b.book_category, b.transaction_type,
        b.price, b.status, b.popularity_score,
        COALESCE(i.available_quantity, 0) AS qty
FROM    Libro b
LEFT JOIN Inventario i ON i.book_id = b.id
WHERE   b.id = ?`
    args := []any{id}
    if onlyAvailable != nil && *onlyAvailable {
        q += " AND COALESCE(i.available_quantity,0) > 0"
    }

    row := r.dbconn.QueryRowContext(ctx, q, args...)
    var b Book
    var qty int64
    if err := row.Scan(&b.ID, &b.BookName, &b.BookCategory, &b.TransactionType,
        &b.Price, &b.Status, &b.PopularityScore, &qty); err != nil {
        if errors.Is(err, sql.ErrNoRows) {
            return BookWithInventory{}, nil
        }
        return BookWithInventory{}, err
    }
    return BookWithInventory{Book: &b, AvailableQuantity: qty}, nil
}

func (r *sqliteRepository) CreateBook(ctx context.Context, b *Book, initialStock int64) (int64, error) {
    tx, err := r.dbconn.BeginTx(ctx, nil)
    if err != nil { return 0, err }
    defer func() { if err != nil { _ = tx.Rollback() } }()

    res, err := tx.ExecContext(ctx, `
INSERT INTO Libro (book_name, book_category, transaction_type, price, status, popularity_score)
VALUES (?, ?, ?, ?, ?, ?)`,
        strings.TrimSpace(b.BookName),
        strings.TrimSpace(b.BookCategory),
        strings.ToLower(strings.TrimSpace(b.TransactionType)),
        b.Price,
        b.Status,           
        b.PopularityScore,
    )
    if err != nil { return 0, err }

    sqlRes, ok := res.(sql.Result)
    if !ok {
        return 0, errors.New("failed to assert result to sql.Result")
    }
    bookID, err := sqlRes.LastInsertId()
    if err != nil { return 0, err }

    _, err = tx.ExecContext(ctx, `
INSERT INTO Inventario (book_id, available_quantity)
VALUES (?, ?)`,
        bookID, initialStock,
    )
    if err != nil { return 0, err }

    if err = tx.Commit(); err != nil { return 0, err }
    return bookID, nil
}

func (r *sqliteRepository) UpdateBook(ctx context.Context, b *Book, stock *int64) (err error) {
    tx, err := r.dbconn.BeginTx(ctx, nil)
    if err != nil { return err }
    defer func() { if err != nil { _ = tx.Rollback() } }()

    _, err = tx.ExecContext(ctx, `
UPDATE Libro
SET     book_name = ?,
        book_category = ?,
        transaction_type = ?,
        price = ?,
        status = ?,
        popularity_score = ?
WHERE   id = ?`,
        strings.TrimSpace(b.BookName),
        strings.TrimSpace(b.BookCategory),
        strings.ToLower(strings.TrimSpace(b.TransactionType)),
        b.Price,
        b.Status,
        b.PopularityScore,
        b.ID,
    )
    if err != nil { return err }

    if stock != nil {
        // Update inventory quantity for the book
        _, err = tx.ExecContext(ctx, `
UPDATE Inventario
SET available_quantity = ?
WHERE book_id = ?`,
            *stock, b.ID,
        )
        if err != nil { return err }

        // (Opcional) sincroniza status si quieres mantenerlo coherente
        // e.g., status = (available_quantity > 0)
        _, err = tx.ExecContext(ctx, `
UPDATE Libro
SET status = CASE WHEN (SELECT available_quantity FROM Inventario WHERE book_id = ?) > 0 THEN 1 ELSE 0 END
WHERE id = ?`, b.ID, b.ID)
        if err != nil { return err }
    }

    return tx.Commit()
}
