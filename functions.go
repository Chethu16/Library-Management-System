package main

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"net/http"
	"strconv"
)

type LibraryManagement struct {
	db *sql.DB
}

func enableCORS(h http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Access-Control-Allow-Origin", "*")
		w.Header().Set("Access-Control-Allow-Methods", "POST, GET, OPTIONS, PUT, DELETE")
		w.Header().Set("Access-Control-Allow-Headers", "Content-Type, Authorization")
		if r.Method == "OPTIONS" {
			w.WriteHeader(http.StatusNoContent)
			return
		}
		h.ServeHTTP(w, r)
	})
}

func (lm *LibraryManagement) RegisterUser(w http.ResponseWriter, r *http.Request) {
	var register User
	err := json.NewDecoder(r.Body).Decode(&register)
	if err != nil {
		fmt.Println("error in register")
		return
	}
	query := `INSERT INTO users(user_id,user_name,user_email,user_password,user_balance) VALUES ($1,$2,$3,$4,$5)`
	_, err = lm.db.Exec(query, register.UserID, register.UserName, register.UserEmail, register.UserPassword, register.UserBalance)
	if err != nil {
		fmt.Println("query execution error")
		return
	}
	err = json.NewEncoder(w).Encode(map[string]string{"message": "registerd succesfull"})
	if err != nil {
		fmt.Println("encode error")
		return
	}
}

func (lm *LibraryManagement) LoginUser(w http.ResponseWriter, r *http.Request) {
	var login User
	err := json.NewDecoder(r.Body).Decode(&login)
	if err != nil {
		fmt.Println("login json error")
		return
	}
	query := `SELECT user_password FROM users WHERE user_email=$1`
	var password string
	err = lm.db.QueryRow(query, login.UserEmail).Scan(&password)
	if err != nil {
		fmt.Println("login execution error", err)
		return
	}
	if login.UserPassword != password {
		fmt.Println("Incorrect password")
		return

	}
	err = json.NewEncoder(w).Encode(map[string]string{"message": "login succesfull"})
	if err != nil {
		fmt.Println("login encode error")
		return
	}

}

func (lm *LibraryManagement) Recharge(w http.ResponseWriter, r *http.Request) {
	var RechargeDetail User
	err := json.NewDecoder(r.Body).Decode(&RechargeDetail)
	if err != nil {
		fmt.Println("recharge json error")
		return
	}
	var user_id, balance string
	query := `SELECT user_id,user_balance FROM users WHERE user_email=$1`
	err = lm.db.QueryRow(query, RechargeDetail.UserEmail).Scan(&user_id, &balance)
	if err != nil {
		fmt.Println("Recharge querry execution error", err)
		return
	}
	if RechargeDetail.UserID != user_id {
		fmt.Println("incorrect user_id")
		return
	}
	var amount, databaseBalance int
	amount, err = strconv.Atoi(RechargeDetail.UserBalance)
	if err != nil {
		fmt.Println("error in converting balance")
		return
	}

	databaseBalance, err = strconv.Atoi(balance)
	if err != nil {
		fmt.Println("database balance converting error")
		return
	}
	totalbalance := strconv.Itoa(databaseBalance + amount)

	query1 := `UPDATE users SET user_balance=$1 WHERE user_email=$2`
	_, err = lm.db.Exec(query1, totalbalance, RechargeDetail.UserEmail)
	if err != nil {
		fmt.Println("query update error")
		return
	}
	err = json.NewEncoder(w).Encode(map[string]string{"message": "Recharge Succesfully"})
	if err != nil {
		fmt.Println("error in updateBalance error")
		return
	}
}

func (lm *LibraryManagement) AddBook(w http.ResponseWriter, r *http.Request) {
	var AddBookDetails Book
	err := json.NewDecoder(r.Body).Decode(&AddBookDetails)
	if err != nil {
		fmt.Println("Book json error")
		return
	}
	query := `SELECT EXISTS (SELECT 1 FROM books WHERE book_id=$1)`
	var isBookExists bool
	lm.db.QueryRow(query, AddBookDetails.BookID).Scan(&isBookExists)
	if !isBookExists {
		query1 := `INSERT INTO books(book_id,book_name,book_author_name,book_price,no_of_copies)VALUES($1,$2,$3,$4,$5)`
		_, err = lm.db.Exec(query1, AddBookDetails.BookID, AddBookDetails.BookName, AddBookDetails.BookAuthorName, AddBookDetails.BookPrice, AddBookDetails.NumberOfCopies)
		if err != nil {
			fmt.Println("query execution error", err)
			return
		}
		err = json.NewEncoder(w).Encode(map[string]string{"message": "Book Added Succesfully"})
		if err != nil {
			fmt.Println("encode error")
			return
		}
		return
	}
	query2 := `UPDATE books SET no_of_copies=$1 WHERE book_id=$2`
	query3 := `SELECT no_of_copies FROM books WHERE book_id=$1`
	var dbNoOfCopies string
	err = lm.db.QueryRow(query3, AddBookDetails.BookID).Scan(&dbNoOfCopies)
	if err != nil {
		fmt.Println("unable to fetch no of copies")
		return
	}
	var dbCopies, userCopies int
	dbCopies, err = strconv.Atoi(dbNoOfCopies)
	if err != nil {
		fmt.Println("error converting data")
		return
	}
	userCopies, err = strconv.Atoi(AddBookDetails.NumberOfCopies)
	if err != nil {
		fmt.Println("error converting data")
		return
	}
	total := strconv.Itoa(dbCopies + userCopies)

	_, err = lm.db.Exec(query2, total, AddBookDetails.BookID)
	if err != nil {
		fmt.Println("error updating table")
		return
	}
	err = json.NewEncoder(w).Encode(map[string]string{"message": "Book Added Succesfully"})
	if err != nil {
		fmt.Println("error encoding json")
		return
	}
}

func (lm *LibraryManagement) DeleteBook(w http.ResponseWriter, r *http.Request) {
	var DeleteBookDetails Book
	err := json.NewDecoder(r.Body).Decode(&DeleteBookDetails)
	if err != nil {
		fmt.Println("error in json")
		return
	}
	query1 := `SELECT EXISTS(SELECT 1 FROM books WHERE book_id=$1)`
	var isBookNotExist bool
	err = lm.db.QueryRow(query1, DeleteBookDetails.BookID).Scan(&isBookNotExist)
	if err != nil {
		fmt.Println("query execution error")
		return
	}
	if isBookNotExist {

		query := `DELETE FROM books WHERE book_id=$1`

		_, err = lm.db.Exec(query, DeleteBookDetails.BookID)
		if err != nil {
			fmt.Println("error in query execution")
			return
		}
		err = json.NewEncoder(w).Encode(map[string]string{"message": "Book deleted Succesfully"})
		if err != nil {
			fmt.Println("json encode error")
			return
		}

		return
	}
	err = json.NewEncoder(w).Encode(map[string]string{"message": "book is not found"})
	if err != nil {
		fmt.Println("encode error")
		return
	}
}

func (lm *LibraryManagement) BorrowBook(w http.ResponseWriter, r *http.Request) {
	var BorrowBookDetails Book
	err := json.NewDecoder(r.Body).Decode(&BorrowBookDetails)
	if err != nil {
		fmt.Println("error in json")
		return
	}
	query := `SELECT no_of_copies FROM books WHERE book_id=$1`
	query1 := `UPDATE books SET no_of_copies=$1 WHERE book_id=$2`
	var dbBorrowBook string
	err = lm.db.QueryRow(query, BorrowBookDetails.BookID).Scan(&dbBorrowBook)
	if err != nil {
		fmt.Println("query execution error", err)
		return
	}
	fmt.Println(dbBorrowBook)
	var dbBorrowBookCopies, userBorrowBookCopies int
	dbBorrowBookCopies, err = strconv.Atoi(dbBorrowBook)
	if err != nil {
		fmt.Println("converting error 1 ", err)
		return
	}
	userBorrowBookCopies, err = strconv.Atoi(BorrowBookDetails.NumberOfCopies)
	if err != nil {
		fmt.Println("converting error 2 ", err)
		return
	}
	total := strconv.Itoa(dbBorrowBookCopies - userBorrowBookCopies)
	_, err = lm.db.Exec(query1, total, BorrowBookDetails.BookID)
	if err != nil {
		fmt.Println("error updating table")
		return
	}
	err = json.NewEncoder(w).Encode(map[string]string{"message": "book borrowed succesfully "})
	if err != nil {
		fmt.Println("encode error")
		return
	}
}

func ReturnBook(w http.ResponseWriter, r *http.Request) {

}

func PurchaseBook(w http.ResponseWriter, r *http.Request) {

}
