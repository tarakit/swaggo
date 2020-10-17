package controllers

import (
	"bapi/models"
	"encoding/json"
	"fmt"

	"github.com/astaxie/beego"
)

// Operations about Books
type BookController struct {
	beego.Controller
}

// @Title CreateBook
// @Description create books
// @Param	body		body 	models.Book	true		"body for book content"
// @Success 200 {int} models.Book.ID
// @Failure 403 body is empty
// @router / [post]
func (u *BookController) Post() {
	var book models.Book
	json.Unmarshal(u.Ctx.Input.RequestBody, &book)
	uid := models.AddBook(book)
	u.Data["json"] = map[string]string{"_id": uid}
	// u.Data["json"] = book
	u.ServeJSON()
}

// @Title GetAll
// @Description get all Books
// @Success 200 {object} models.Book
// @router / [get]
func (u *BookController) GetAll() {
	books := models.GetAllBooks()

	fmt.Println(books)
	u.Data["json"] = books
	u.ServeJSON()
}

// @Title Get
// @Description get book by uid
// @Param	uid		path 	string	true		"The key for staticblock"
// @Success 200 {object} models.Book
// @Failure 403 :uid is empty
// @router /:uid [get]
func (u *BookController) Get() {
	uid := u.GetString(":uid")
	if uid != "" {
		book, err := models.GetBook(uid)
		fmt.Print(uid)
		if err != nil {
			u.Data["json"] = book
		} else {
			u.Data["json"] = err.Error()
		}
	}
	u.ServeJSON()
}

// @Title Update
// @Description update the book
// @Param	uid		path 	string	true		"The uid you want to update"
// @Param	body		body 	models.Book	true		"body for user content"
// @Success 200 {object} models.Book
// @Failure 403 :uid is not int
// @router /:uid [put]
func (u *BookController) Put() {
	uid := u.GetString(":uid")
	if uid != "" {
		var book models.Book
		json.Unmarshal(u.Ctx.Input.RequestBody, &book)
		uu, err := models.UpdateBook(uid, &book)
		if err != nil {
			u.Data["json"] = uu
		} else {
			u.Data["json"] = err.Error()
		}
	}
	u.ServeJSON()
}

// @Title Delete
// @Description delete the book
// @Param	uid		path 	string	true		"The uid you want to delete"
// @Success 200 {string} delete success!
// @Failure 403 uid is empty
// @router /:uid [delete]
func (u *BookController) Delete() {
	uid := u.GetString(":uid")
	models.DeleteBook(uid)
	u.Data["json"] = "delete success!"
	u.ServeJSON()
}

// @Title Uplaod_Image
// @Description uplaod image to the system
// @Accept mpfd
// @Param	file	file 	file	true		"The file for uplading image"
// @Success 200 {string} upload success
// @Failure 403 book not exist
// @router /file [post]
func (u *BookController) UploadImage() {

	file, head, err := u.GetFile("file")

	if err != nil {
		u.Ctx.WriteString("Get file failed")
		return
	}
	defer file.Close()

	filename := head.Filename
	err = u.SaveToFile("file", "./photos/"+filename)
	if err != nil {
		u.Ctx.WriteString("Upload failed 1")
	} else {
		u.Ctx.WriteString("Upload succeeded")
	}

	// u.SaveToFile("file", filePath)
	// u.SaveToFile(handler.Filename, "./photos/")

	// models.UploadBookCover(file, handler, err)

}

// // @Title logout
// // @Description Logs out current logged in user session
// // @Success 200 {string} logout success
// // @router /logout [get]
// func (u *UserController) Logout() {
// 	u.Data["json"] = "logout success"
// 	u.ServeJSON()
// }
