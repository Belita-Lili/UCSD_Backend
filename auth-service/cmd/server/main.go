package main

import (
	"database/sql"
	"log"
	"os"

	"github.com/LiliBeta/auth-service/internal/auth"
	"github.com/LiliBeta/auth-service/internal/server"

	_ "github.com/go-sql-driver/mysql"
)

func main() {
	// Configuración de la base de datos
	db, err := sql.Open("mysql", "user:password@tcp(localhost:3306)/auth_db")
	if err != nil {
		log.Fatal(err)
	}
	defer db.Close()

	// Configuración de Keycloak
	keycloakAuth := auth.NewKeycloakAuthRepository(
		os.Getenv("KEYCLOAK_URL"),
		os.Getenv("KEYCLOAK_REALM"),
		os.Getenv("KEYCLOAK_CLIENT_ID"),
		os.Getenv("KEYCLOAK_CLIENT_SECRET"),
	)

	// Configurar router
	router := server.SetupRouter(db, keycloakAuth)

	// Iniciar servidor
	if err := router.Run(":8080"); err != nil {
		log.Fatal(err)
	}
}
