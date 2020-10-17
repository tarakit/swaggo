package models

import (
	"bapi/helper"
	"context"
	"errors"
	"fmt"
	"io"
	"log"
	"mime/multipart"
	"os"

	"github.com/google/uuid"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

var (
	BookList []*Book
)

func init() { // TODO Need to reconnect to Mongo again to get the updated the Documents
	// var BookList = []*Book/

	// create a context
	ctx := context.TODO()

	// var books []Book

	collection := helper.ConnectDB()

	cur, _ := collection.Find(ctx, bson.M{})

	// defer : is used to close the conneciton when it is done
	defer cur.Close(ctx)

	// travel through the collection of data
	for cur.Next(ctx) {

		var book Book

		err := cur.Decode(&book) // decode to deserialized the books from Collection of MongoDB
		if err != nil {
			log.Fatal(err)
		}

		BookList = append(BookList, &book)
	}

	if err := cur.Err(); err != nil {
		log.Fatal(err)
	}

}

// swagger:model BookResponse
type Book struct {
	ID     primitive.ObjectID `json:"_id,omitempty" bson:"_id,omitempty"`
	Title  string             `json:"title,omitempty" bson:"title,omitempty"`
	Author *Author            `json:"author,omitempty" bson:"author,omitempty"`
}

type Author struct {
	FirstName string `json:"firstName,omitempty" bson:"firstName,omitempty"`
	LastName  string `json:"lastName,omitempty" bson:"lastName,omitempty"`
}

func AddBook(u Book) string {

	// var book Book
	ctx := context.TODO()

	// book =
	// _ = json.NewDecoder(r.Body).Decode(&book)

	collection := helper.ConnectDB()

	result, err := collection.InsertOne(ctx, u)
	if err != nil {
		fmt.Print(err)
	}
	// json.NewEncoder(w).Encode(result)

	return fmt.Sprint(result)
}

func GetBook(uid string) (u *Book, err error) {

	var book *Book

	collection := helper.ConnectDB()

	id, _ := primitive.ObjectIDFromHex(uid)

	filter := bson.M{"_id": id}
	fmt.Println("Before Decoding")
	err = collection.FindOne(context.TODO(), filter).Decode(&book)

	fmt.Print(book.Author)
	fmt.Println(book.Title)

	if err != nil {
		return
	}
	u = book
	fmt.Println(u.Title)

	return u, errors.New("Book not exists")
}

func GetAllBooks() []*Book {
	return BookList
}

func UpdateBook(uid string, uu *Book) (a *Book, err error) {

	id, _ := primitive.ObjectIDFromHex(uid)

	var book Book

	collection := helper.ConnectDB()

	filter := bson.M{"_id": id}

	update := bson.D{
		{"$set", bson.D{
			{"title", uu.Title},
			{"author", uu.Author}}, // put the Author object rather than the document of Author
		},
	}

	err = collection.FindOneAndUpdate(context.TODO(), filter, update).Decode(&book)

	if err != nil {
		return
	}
	book.ID = id
	a = &book

	return a, errors.New("Book Not Exist")
}

func UploadBookCover(file multipart.File, handler *multipart.FileHeader, err error) string {

	// generate uuid concatenated with chosen file name
	id := uuid.New()
	FILENAME := fmt.Sprint(id, handler.Filename)

	// f, err := os.OpenFile("C:/Users/kitta/go/src/github.com/RestAPI/photos/"+handler.Filename, os.O_WRONLY|os.O_CREATE, 0666)
	f, err := os.OpenFile("./photos/"+FILENAME, os.O_WRONLY|os.O_CREATE, 0666)
	if err != nil {
		return err.Error()
	}
	defer f.Close()
	io.Copy(f, file)

	return "upload success"
}

func DeleteBook(uid string) {
	id, _ := primitive.ObjectIDFromHex(uid)

	// var book Book

	collection := helper.ConnectDB()

	filter := bson.M{"_id": id}

	err := collection.FindOneAndDelete(context.TODO(), filter)

	if err != nil {
		return
	}
}
