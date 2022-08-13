package main

import (
	"fmt"
	"log"
	"net/http"
	"github.com/ErfanMomeniii/AuthForm/configs"
	"gorm.io/driver/mysql"
	validation "github.com/go-ozzo/ozzo-validation"
	"github.com/go-ozzo/ozzo-validation/is"
	"gorm.io/gorm"
)
type User struct {
	gorm.Model
	Username string
	Password string
	Email string
}
func register(w http.ResponseWriter,r *http.Request, db *gorm.DB){
	
	if r.URL.Path != "/register"{
		http.Error(w,"404 not found",http.StatusNotFound);
		return 
	}
	
	if r.Method != "POST" {
		http.Error(w,"Method not supported",http.StatusNotFound);
		return
	}
	
	if err:= r.ParseForm(); err!=nil{
		http.Error(w,"form has a problem",http.StatusBadRequest);
		return 
	}

	if r.FormValue("password") != r.FormValue("confirm_password") {
		http.Error(w,"password and confirmation for password should be the same",http.StatusForbidden);
		return 
	}

	user := &User{
		Username : r.FormValue("username"),
		Password: r.FormValue("password"),
		Email : r.FormValue("email"),
	}

	if err := validation.ValidateStruct(
		validation.Field(&user.Username, validation.Required),
		validation.Field(&user.Password, validation.Required),
		validation.Field(&user.Email, validation.Required,is.Email),
	);err!=nil {
		http.Error(w,"form data has a problem",http.StatusBadRequest);
		return 
	}

	db.Create(user)

	fmt.Fprintf(w,"User was registered");
}
func main() {

	fileServer:=http.FileServer(http.Dir("./static"));
	
	dsn := "root:@tcp(127.0.0.1:3306)/authform?charset=utf8mb4&parseTime=True&loc=Local"
  	db, err := gorm.Open(mysql.Open(dsn), &gorm.Config{})

	if err != nil {
		log.Fatalf("we cant connect to database")
	}

	db.AutoMigrate(&User{})
	
	http.Handle("/",fileServer)
	http.HandleFunc("/register",func(w http.ResponseWriter,r *http.Request){
	 register(w,r,db);
	});
	
	if err := http.ListenAndServe(string(":"+configs.Port),nil);err!=nil {
		log.Fatalf("%v",err);
	}
}
