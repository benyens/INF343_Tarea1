package main // Define el paquete principal del programa

import (
	"database/sql" // Importa el paquete para trabajar con bases de datos SQL
	"uzm-server/internal/users" // Importa el paquete local 'users' que contiene la lógica relacionada con usuarios
	"uzm-server/internal/db" // Importa el paquete local 'db' que contiene la lógica para migrar la base de datos
	"github.com/gin-gonic/gin" // Importa el framework web Gin para crear servidores HTTP
	_ "modernc.org/sqlite"     // Importa el driver SQLite3 (el guion bajo indica importación solo para efectos secundarios)

	"log"  // Importa el paquete para registro de logs
)

func main() { // Función principal, punto de entrada del programa
	// Conexión a la base de datos SQLite
	dbconn, err := sql.Open("sqlite", "./uzm.db") // Abre una conexión
	if err != nil {                           // Verifica si hubo un error al abrir la conexión
		log.Fatal(err) // Registra el error y termina el programa
	}
	defer dbconn.Close() // Asegura que la conexión se cierre al finalizar la función

	// Crear tablas
	db.MakeMigrate(dbconn, "schema.sql") // Ejecuta una sentencia SQL para crear la tabla 'users' si no existe

	// DI
	userRepo := users.NewSQLiteRepository(dbconn) // Crea un repositorio de usuarios basado en SQLite
	userService := users.NewService(userRepo)  // Crea un servicio de usuarios utilizando el repositorio
	userHandler := users.NewHandler(userService) // Crea un manejador de usuarios utilizando el servicio

	// Inicializa el router Gin
	router := gin.Default() // Crea un router con las configuraciones por defecto
	router.Group("/api/v1")
	userHandler.RegisterRoutes(router) // Registra las rutas del manejador de usuarios bajo el grupo /api/v1

	router.GET("/health", func(c *gin.Context) { c.JSON(200, gin.H{"status":"ok"}) })

	router.Run(":8080") // Inicia el servidor en el puerto 8080
}
