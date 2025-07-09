package main

type User struct {
	UserID string `json:"user_id"`
	UserName string `json:"user_name"`
	UserEmail string `json:"user_email"`
	UserPassword string `json:"user_password"`
	UserBalance string `json:"user_balance"`
}

type Book struct {
	BookID string `json:"book_id"`
	BookName string `json:"book_name"`
	BookAuthorName string `json:"book_author_name"`
	BookPrice string `json:"book_price"`
	NumberOfCopies string `json:"no_of_copies"`
}

type Purchase struct{
	UserId string `json:"user_id"`
	BookId string `json:"book_id"`
	NumberOfCopies string `json:"no_of_copies"`

}



