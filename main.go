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
	Email string  `gorm:"primaryKey"`
}

func (user User) Validate() error {
    return validation.ValidateStruct(&user,
        validation.Field(&user.Username, validation.Required, validation.Length(5, 20)),
		validation.Field(&user.Password,validation.Required),
		validation.Field(&user.Email, validation.Required,is.Email),
    )
}

func register(w http.ResponseWriter, r *http.Request, db *gorm.DB){
	
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

	user := User{
		Username : r.FormValue("username"),
		Password: r.FormValue("password"),
		Email : r.FormValue("email"),
	}

	err := user.Validate();

	if err != nil{
		http.Error(w,"form data has a problem",http.StatusBadRequest)
		return 
	}

	db.Create(&user)

	fmt.Fprintf(w,"User was registered");
}

func login(w http.ResponseWriter, r *http.Request, db *gorm.DB){
	
	if r.URL.Path != "/login"{
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
	
	user := User {
		Password: r.FormValue("password"),
		Email: r.FormValue("email"),
	}
	
	if err := validation.ValidateStruct(&user,
		validation.Field(&user.Password, validation.Required),
		validation.Field(&user.Email, validation.Required,is.Email),
	);err!=nil {
		http.Error(w,"form data has a problem",http.StatusBadRequest)
		return 
	}

	result := db.Where("email = ?", r.FormValue("email")).Where("password = ?", r.FormValue("password")).Find(&user);
	
	if result.Error!=nil{
		http.Error(w,"we cant find user with this data",http.StatusNotFound)
		return
	}

	fmt.Fprintf(w,"%#v has been singing in",user.Username)
}

func main() {

	fileServer:=http.FileServer(http.Dir("./static"));
	
	dsn := "root:@tcp(127.0.0.1:"+configs.DatabasePort+")/authform?charset=utf8mb4&parseTime=True&loc=Local"
  	db, err := gorm.Open(mysql.Open(dsn), &gorm.Config{})
	
	if err != nil {
		log.Fatalf("we cant connect to database")
	}

	db.AutoMigrate(&User{})
	
	database,err:= db.DB()
	defer database.Close()

	http.Handle("/",fileServer)
	http.HandleFunc("/register",func(w http.ResponseWriter,r *http.Request){
	 register(w,r,db);
	});
	http.HandleFunc("/login",func(w http.ResponseWriter,r *http.Request){
		login(w,r,db);
	   });

	if err := http.ListenAndServe(string(":"+configs.Port),nil);err!=nil {
		log.Fatalf("%v",err);
	}
}
