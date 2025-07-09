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

func (lm *LibraryManagement)ReturnBook(w http.ResponseWriter, r *http.Request) {
	 var ReturnBookDetails Book
	 err:=json.NewDecoder(r.Body).Decode(&ReturnBookDetails)
	 if err !=nil{
		fmt.Println("json error")
		return
	 }
	 var dbreturnbookDetails string
	 query :=`SELECT no_of_copies FROM books WHERE book_id=$1`
	 err =lm.db.QueryRow(query,ReturnBookDetails.BookID).Scan(&dbreturnbookDetails)
	 if err !=nil{
		fmt.Println("Unable to fetch")
		return
	 }
	 query1:=`UPDATE books SET no_of_copies=$1 WHERE book_id=$2`
	 var dbreturnBookcopies, userreturnBookcopies int
	 dbreturnBookcopies,err = strconv.Atoi(dbreturnbookDetails)
	 if err!=nil{
		fmt.Println("Converting error")
		return
	 }
	 userreturnBookcopies,err = strconv.Atoi(ReturnBookDetails.NumberOfCopies)
	 if err!=nil{
		fmt.Println("Converting error 2")
		return
	 }
	 total:=strconv.Itoa(dbreturnBookcopies+userreturnBookcopies)
	 _,err=lm.db.Exec(query1,total,ReturnBookDetails.BookID)
	 if err !=nil{
		fmt.Println("query execution error",err)
		return
	 }
	err= json.NewEncoder(w).Encode(map[string]string{"message":"Book returnd Succesfully"})
	if err !=nil{
		fmt.Println("Encode error")
		return
	}

}

func (lm *LibraryManagement)PurchaseBook(w http.ResponseWriter, r *http.Request) {
	var PurchaseBookDetails Purchase
	err:=json.NewDecoder(r.Body).Decode(&PurchaseBookDetails)
	if err!=nil{
		fmt.Println("json error")
		return
	}
	query:=`SELECT user_balance FROM users WHERE user_id=$1`
	var dbBalance string
	err=lm.db.QueryRow(query,PurchaseBookDetails.UserId).Scan(&dbBalance)
	if err!=nil{
		fmt.Println("unable to fetch 1",err)
		return
	}
	query1:=`SELECT no_of_copies,book_price FROM books WHERE book_id=$1`
	var dbbookcopies, dbbookprice string
	err=lm.db.QueryRow(query1,PurchaseBookDetails.BookId).Scan(&dbbookcopies,&dbbookprice)
	if err!=nil{
		fmt.Println("unable to fetch 2")
		return
	}
	if dbbookcopies < PurchaseBookDetails.NumberOfCopies{
		err=json.NewEncoder(w).Encode(map[string]string{"message":"error"})
		if err!=nil{
			fmt.Println("erorr in encode 1")
			return
		}
		return
	}else
	{
	var Dbnumberofcopies,Usernumberofcopies,Dbbookprice ,DBuserbalance int
	Dbnumberofcopies,err= strconv.Atoi(dbbookcopies)
	if err!=nil{
		fmt.Println("error in converting 1")
		return
	}
	Usernumberofcopies,err=strconv.Atoi(PurchaseBookDetails.NumberOfCopies)
	if err!=nil{
		fmt.Println("converting error 2")
		return
	}
	Dbbookprice,err=strconv.Atoi(dbbookprice)
	if err!=nil{
		fmt.Println("converting error 3")
		return
	}
	DBuserbalance,err=strconv.Atoi(dbBalance)
	if err !=nil{
		fmt.Println("converting error 4")
		return
	}
	total:= Usernumberofcopies*Dbbookprice
	total3:= DBuserbalance-total

	total2 :=strconv.Itoa(Dbnumberofcopies-Usernumberofcopies)

	query2:=`UPDATE books SET no_of_copies=$1  WHERE book_id=$2`
	_,err=lm.db.Exec(query2,total2,PurchaseBookDetails.BookId)
		if err!=nil{
			fmt.Println("query execution error 1",err)
			return
		}
		query3:=`UPDATE users SET user_balance=$1  WHERE user_id=$2`
		_,err=lm.db.Exec(query3,total3,PurchaseBookDetails.UserId)
		if err !=nil{
			fmt.Println("query execution error 2",err)
			return
		}
		err=json.NewEncoder(w).Encode(map[string]string{"message":"Book purchased Succesfully"})
		if  err!=nil{
			fmt.Println("encode error 2")
			return
		}
	}
	}
	

	





	










