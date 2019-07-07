package database

import (
	"fmt"
	"github.com/acidobinario/Angelhack-2019-Backend/models"
	"github.com/acidobinario/Angelhack-2019-Backend/random"
	"github.com/jinzhu/gorm"
	_ "github.com/jinzhu/gorm/dialects/mysql"
	"golang.org/x/crypto/bcrypt"
	"log"
	"os"
)

var connection *gorm.DB

//variables de entorno
const engineSql string = "mysql"

var username = os.Getenv("DBUSER")
var password = os.Getenv("DBPASS")
var database = os.Getenv("DBNAME")
var host = "tcp(" + os.Getenv("DBIP") + ":3306)"

func init() {
	InitializeDatabase()
}

func InitializeDatabase() {
	connection = ConnectORM(CreateString())
	if connection == nil {
		log.Fatal("OHSHI")
	}
	log.Println("[i] connected to DB.")

	log.Println("Migrating Users")
	connection.AutoMigrate(&models.User{})
	log.Println("Migrating Employees")
	connection.AutoMigrate(&models.Employee{})
	log.Println("Migrating Transactions")
	connection.AutoMigrate(&models.Transaction{})
	log.Println(" OK!")
	log.Println("[i] DB Automigrated.")
	connection.DB().SetMaxOpenConns(5)
}

func CloseConnection() {
	connection.Close()
	log.Println("closed connection with DB.")
}

func ConnectORM(stringConnection string) *gorm.DB {
	connection, err := gorm.Open(engineSql, stringConnection)
	if err != nil {
		log.Println(err)
		return nil
	}
	return connection
}

func CreateString() string {
	fmt.Println(username + ":" + password + "@" + host + "/" + database + "?charset=utf8&parseTime=True&loc=Local")
	return username + ":" + password + "@" + host + "/" + database + "?charset=utf8&parseTime=True&loc=Local"
}

///////////////////////////////

func UsernameExists(usename string) bool {
	var user models.User
	if err := connection.Table("users").Where("username = ?", username).First(&user).Error; gorm.IsRecordNotFoundError(err) {
		return false
	}
	return true
}

func CheckUser(username, password string) bool {
	var user models.User
	if err := connection.Table("users").Where("username = ?", username).First(&user).Error; gorm.IsRecordNotFoundError(err) {
		return false
	}
	//hashedPass, _ := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	//log.Println(string(hashedPass))
	if err := bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(password)); err != nil {
		//log.Println(err)
		return false
	} else {
		return true
	}
}

func GetUser(username string) models.User {
	var user models.User
	if err := connection.Table("users").Where("username = ?", username).First(&user).Error; gorm.IsRecordNotFoundError(err) {
		return models.User{}
	}
	return user
}

func addUser(user models.User) {
	connection.Debug().Table("users").Save(&user)
}

func CreateUser(user models.User) {
	addUser(user)
}

func CreateEmployee(user models.User, salary float64, company models.User) {
	addUser(user)
	createdUser := GetUser(user.Username)
	employee := models.Employee{
		CompanyId: company.ID,
		UserId:    createdUser.ID,
		Balance:   0,
		Code:      "",
		IsUsed:    false,
		Salary:    salary,
	}
	connection.Save(&employee)
}

func GetEmployees(companyUser string) (err error, employees []models.EmployeeList) {
	//get the company Id
	var Company models.User
	if err := connection.Debug().Table("users").Where("username = ?", companyUser).First(&Company).Error; gorm.IsRecordNotFoundError(err){
		return err, nil
	}
	// search for the employeeeeee that share the same Id
	var Employeesone []models.Employee
	if err := connection.Debug().Table("employees").Where("company_id = ?", Company.ID).Scan(&Employeesone).Error; gorm.IsRecordNotFoundError(err){
		return err, nil
	}
	var output []models.EmployeeList
	for _, v := range Employeesone {
		var tmpEmployee models.User
		connection.Debug().Table("users").Where("id = ?", v.UserId).First(&tmpEmployee)
		tmpEmployee.Password=""
		output = append(output, models.EmployeeList{
			User:     tmpEmployee,
			Employee: v,
		})
		//employees = append(employees, tmpEmployee)
	}
	return nil, output
}

func CreateTransaction(from uint, to uint, amount float64){
	newTransaction := models.Transaction{
		From:   from,
		To:     to,
		Amount: amount,
	}
	connection.Save(&newTransaction)
}

func GetBalance(username string)  {
	var user models.User
	connection.Table("user").Where("username = ?", username).First(&user)
}

func GenerateCode(id uint){
 	var thisEmployee models.Employee
	connection.Table("employees").Where("user_id = ?", id).First(&thisEmployee)
 	thisEmployee.Code = random.String(10)
 	thisEmployee.IsUsed = false
 	connection.Save(&thisEmployee)
}

func VerifyQr(code string) (stuff models.Employee, err error){
	// if the qr code is found on the !is_used
	var employeeeee models.Employee
	if err := connection.Debug().Table("employees").Where("is_used = false AND code = ?", code).First(&employeeeee).Error; gorm.IsRecordNotFoundError(err) {
		return models.Employee{}, err
	}
	return employeeeee, nil
}

func SetCodeUsed(aaaaaaaa models.Employee){
		// if the qr code is found on the !is_used

	connection.Debug().Table("employees").Where("is_used = false AND code = ?", aaaaaaaa.Code).Update("is_used", "1")
	connection.Debug().Table("employees").Where(aaaaaaaa).Update("balance", "0")
}