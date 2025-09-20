package main

import (
	"bufio"
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
	"os"
	"strconv"
	"strings"
	"time"
	"sort"
)

// ===== Config =====

func baseURL() string {
	if v := os.Getenv("UZM_BASE_URL"); v != "" {
		return v
	}
	return "http://localhost:8080"
}

var httpc = &http.Client{Timeout: 10 * time.Second}

// ===== DTOs =====

type Book struct {
	ID              int64  `json:"id"`
	BookName        string `json:"book_name"`
	BookCategory    string `json:"book_category"`
	TransactionType string `json:"transaction_type"`
	Price           int64  `json:"price"`
	Status          any    `json:"status"`
	PopularityScore int64  `json:"popularity_score"`
	Inventory       struct {
		AvailableQuantity int64 `json:"available_quantity"`
	} `json:"inventory"`
}

type BooksList struct {
	Books []Book `json:"books"`
}

type CreateBookReq struct {
	BookName        string `json:"book_name"`
	BookCategory    string `json:"book_category"`
	TransactionType string `json:"transaction_type"` // "venta" | "arriendo"
	Price           int64  `json:"price"`
	Status          bool   `json:"status"`
	Stock           int64  `json:"stock"`
}

type UpdateBookReq struct {
	BookName        *string `json:"book_name,omitempty"`
	BookCategory    *string `json:"book_category,omitempty"`
	TransactionType *string `json:"transaction_type,omitempty"`
	Price           *int64  `json:"price,omitempty"`
	Status          *bool   `json:"status,omitempty"`
	PopularityScore *int64  `json:"popularity_score,omitempty"`
	Stock           *int64  `json:"stock,omitempty"`
}

type User struct {
	ID        int64  `json:"id"`
	FirstName string `json:"first_name"`
	LastName  string `json:"last_name"`
	Email     string `json:"email"`
	USMPesos  int64  `json:"usm_pesos"`
	Password  string `json:"password,omitempty"`
}

type UsersList struct {
	Users []User `json:"users"`
}

type CreateUserReq struct {
	FirstName string `json:"first_name"`
	LastName  string `json:"last_name"`
	Email     string `json:"email"`
	Password  string `json:"password"`
	USMPesos  int64  `json:"usm_pesos"`
}

type UpdatePesosReq struct {
	Amount int64 `json:"amount"`
}

// ===== HTTP helpers =====

func getJSON(path string, v any) error {
	req, _ := http.NewRequest(http.MethodGet, baseURL()+path, nil)
	res, err := httpc.Do(req)
	if err != nil {
		return err
	}
	defer res.Body.Close()
	if res.StatusCode >= 400 {
		var m map[string]any
		_ = json.NewDecoder(res.Body).Decode(&m)
		return fmt.Errorf("GET %s -> %s: %v", path, res.Status, m)
	}
	return json.NewDecoder(res.Body).Decode(v)
}

func postJSON(path string, body any, v any) error {
	b, _ := json.Marshal(body)
	req, _ := http.NewRequest(http.MethodPost, baseURL()+path, bytes.NewReader(b))
	req.Header.Set("Content-Type", "application/json")
	res, err := httpc.Do(req)
	if err != nil {
		return err
	}
	defer res.Body.Close()
	if res.StatusCode >= 400 {
		var m map[string]any
		_ = json.NewDecoder(res.Body).Decode(&m)
		return fmt.Errorf("POST %s -> %s: %v", path, res.Status, m)
	}
	if v != nil {
		return json.NewDecoder(res.Body).Decode(v)
	}
	return nil
}

func patchJSON(path string, body any, v any) error {
	b, _ := json.Marshal(body)
	req, _ := http.NewRequest(http.MethodPatch, baseURL()+path, bytes.NewReader(b))
	req.Header.Set("Content-Type", "application/json")
	res, err := httpc.Do(req)
	if err != nil {
		return err
	}
	defer res.Body.Close()
	if res.StatusCode >= 400 {
		var m map[string]any
		_ = json.NewDecoder(res.Body).Decode(&m)
		return fmt.Errorf("PATCH %s -> %s: %v", path, res.Status, m)
	}
	if v != nil {
		return json.NewDecoder(res.Body).Decode(v)
	}
	return nil
}

// ===== CLI helpers =====

var in = bufio.NewReader(os.Stdin)

func prompt(s string) string {
	fmt.Print(s)
	t, _ := in.ReadString('\n')
	return strings.TrimSpace(t)
}

func mustAtoi64(s string) int64 {
	v, _ := strconv.ParseInt(strings.TrimSpace(s), 10, 64)
	return v
}

func pause() {
	fmt.Print("\n(Enter para continuar) ")
	_, _ = in.ReadString('\n')
}

// ===== Menú 1: Registrarse / Iniciar sesion =====

func menuInicio() (user *User, ok bool) {
	for {
		fmt.Println("\n=== Menú ===")
		fmt.Println("1) Registrarse")
		fmt.Println("2) Iniciar sesion")
		fmt.Println("3) Terminar ejecución")

		switch prompt("> ") {
		case "1":
			u := crearUsuario()
			if u != nil {
				fmt.Printf("Usuario creado, ahora puedes iniciar sesion.\n")
			}
		case "2":
			Correo := prompt("Correo: ")

			if strings.TrimSpace(Correo) == "" {
				fmt.Println("Correo inválido")
				continue
			}

			var me User
			if err := getJSON(fmt.Sprintf("/api/users/email/%s", Correo), &me); err != nil || me.ID == 0 {
				fmt.Println("Usuario no encontrado:", err)
				continue
			}

			passwordCorrect := me.Password
			passwordTry := prompt("Password: ")
			if passwordTry==passwordCorrect {
				fmt.Println("Contraseña incorrecta")
				continue
			}

			fmt.Printf("Bienvenido, %s %s!\n", me.FirstName, me.LastName)
			return &me, true
		case "3":
			return nil, false
			default:
				fmt.Println("Opción inválida")
			}
		}
	}


func crearUsuario() *User {
	fn := prompt("Nombre: ")
	ln := prompt("Apellido: ")
	em := prompt("Email: ")
	pw := prompt("Password: ")
	usm := mustAtoi64(prompt("USM Pesos iniciales (0 si no aplica): "))

	req := CreateUserReq{FirstName: fn, LastName: ln, Email: em, Password: pw, USMPesos: usm}
	var created User
	if err := postJSON("/api/users", req, &created); err != nil {
		// si tu handler devuelve { "user_id": n } en vez del user completo,
		// hacemos un segundo POST sin decodificar respuesta estricta (solo para crear):
		if e2 := postJSON("/api/users", req, nil); e2 != nil {
			fmt.Println("Error creando usuario:", err)
			pause()
			return nil
		}
		// no sabemos el id; pedirá entrar con ID manual luego
		return nil
	}
	return &created
}

// ===== Menú 2: Operaciones soportadas por la API actual =====

func menuPrincipal(user *User) {
	for {
		fmt.Println("\n=== Menú Principal ===")
		fmt.Println("1) Ver catálogo (solo disponibles)")
		fmt.Println("2) Ver catálogo (incluye agotados)")
		fmt.Println("3) Ver carro de compras")
		fmt.Println("4) Mis prestamos")
		fmt.Println("5) Ver mi cuenta")
		fmt.Println("6) Ver libros populares")
		fmt.Println("7) Salir")
		switch prompt("> ") {
		case "1":
			listarLibros(false)
		case "2":
			listarLibros(true)
		case "3":
			verCarro(user)
		case "4":
			misPrestamos(user)
		case "5":
			verMiCuenta(user)
		case "6":
			verPopulares()
		case "7":
			return
		default:
			fmt.Println("Opción inválida")
		}
	}
}


func listarLibros(includeAll bool) {
	url := "/api/books"
	if includeAll { url += "?status=false" }
	var out BooksList
	if err := getJSON(url, &out); err != nil {
		fmt.Println("Error:", err)
		pause()
		return
	}
	if len(out.Books) == 0 {
		fmt.Println("(sin libros)")
		pause()
		return
	}
	fmt.Println("-----------------------------------------------------------------")
	fmt.Printf("| %-7s | %-20s | %-10s | %-9s | %-5s | %-5s |\n", "ID", "Nombre", "Categoría", "Tipo", "Valor", "Stock")
	fmt.Println("-----------------------------------------------------------------")
	for _, b := range out.Books {
		fmt.Printf("| %-7d | %-20s | %-10s | %-9s | %-5d | %-5d |\n",
			b.ID, b.BookName, b.BookCategory, b.TransactionType, b.Price, b.Inventory.AvailableQuantity)
	}
	fmt.Println("-----------------------------------------------------------------")
	pause()
}

func verCarro(user *User) {
	fmt.Println("(no implementado)")
	pause()
}

func misPrestamos(user *User) {
	fmt.Println("(no implementado)")
	pause()
}

func verPopulares() {
    var out BooksList
    if err := getJSON("/api/books", &out); err != nil {
        fmt.Println("Error:", err)
        pause()
        return
    }

    // Ordenar por PopularityScore descendente
    sort.Slice(out.Books, func(i, j int) bool {
        return out.Books[i].PopularityScore > out.Books[j].PopularityScore
    })

    fmt.Println("-----------------------------------------------------------------")
    fmt.Printf("| %-7s | %-20s | %-10s | %-5s | %-5s |\n", "ID", "Nombre", "Categoría", "Valor", "Popularidad")
    fmt.Println("-----------------------------------------------------------------")
    for _, b := range out.Books {
        fmt.Printf("| %-7d | %-20s | %-10s | %-5d | %-5d |\n",
            b.ID, b.BookName, b.BookCategory, b.Price, b.PopularityScore)
    }
    fmt.Println("-----------------------------------------------------------------")
    pause()
}



func verMiCuenta(user *User) {
	var me User
	if err := getJSON(fmt.Sprintf("/api/users/%d", user.ID), &me); err != nil {
		fmt.Println("Error:", err)
		pause()
		return
	}
	for {
		fmt.Println("1. Consultar saldo")
		fmt.Printf("2. Abonar usm pesos")
		fmt.Printf("3. Ver historial de compras y arriendos")
		fmt.Printf("4. Salir")
		switch prompt("> ") {
		case "1":
			fmt.Printf("Saldo actual: %d USM Pesos\n", me.USMPesos)
			pause()
		case "2":
			abonarMiCuenta(user)
		case "3":
			fmt.Println("(no implementado)")
			pause()
		case "4":
			return
		default:
			fmt.Println("Opción inválida")

	}
}
}

func abonarMiCuenta(user *User) {
	amt := mustAtoi64(prompt("Ingrese la cantidad de usm pesos a abonar (+/-): "))
	if err := patchJSON(fmt.Sprintf("/api/users/%d/usm_pesos", user.ID), UpdatePesosReq{Amount: amt}, nil); err != nil {
		fmt.Println("Error:", err)
	} else {
		user.USMPesos += amt
		fmt.Printf("Nuevo saldo de %d usm pesos.\n", user.USMPesos)
	}
	pause()
}

// ===== main =====

func main() {
	for {
		if u, ok := menuInicio(); ok && u != nil {
			menuPrincipal(u)
			// Al salir del menú principal, vuelve al menú inicial
		} else {
			fmt.Println("Muchas gracias por visitarnos")
			return
		}
	}
}
