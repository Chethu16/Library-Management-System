package main

import (
	"database/sql"
	"fmt"
	"net/http"

	_ "github.com/lib/pq"
)
func main() {
	db, err := sql.Open("postgres", "postgresql://library_xdr6_user:7S89MJlSGdp9FumFq6pMr0wicKUID7nt@dpg-d1lq3dje5dus73818cog-a.oregon-postgres.render.com/library_xdr6")
	if err != nil {
		fmt.Println("error connecting to database")
		return
	}
	defer db.Close()

	management := LibraryManagement{db}

	mux := http.NewServeMux()
	mux.HandleFunc("/register", management.RegisterUser)
	mux.HandleFunc("/login", management.LoginUser)
	mux.HandleFunc("/recharge", management.Recharge)
	mux.HandleFunc("/addbook", management.AddBook)
	mux.HandleFunc("/deletebook", management.DeleteBook)
	mux.HandleFunc("/borrowbook", management.BorrowBook)
	mux.HandleFunc("/returnbook",management.ReturnBook)
	mux.HandleFunc("/purchasebook",management.PurchaseBook)
	// Apply CORS
	handlerWithCORS := enableCORS(mux)

	fmt.Println("Server running at http://localhost:8000")
	http.ListenAndServe(":8000", handlerWithCORS)
}
