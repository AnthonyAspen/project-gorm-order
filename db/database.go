package main

import (
  "gorm.io/gorm"
  "gorm.io/driver/mysql"
  "fmt"
  "encoding/json"
  "time"
)

  type Model struct {
    ID        uint       `gorm:"primary_key auto_increment:true;column:id" json:"id"`
 CreatedAt time.Time  `gorm:"column:created_at" json:"created_at"`
 UpdatedAt time.Time  `gorm:"column:updated_at" json:"updated_at"`
 DeletedAt *time.Time `gorm:"column:deleted_at" json:"deleted_at"`
}
type Customer struct {
  Model
  FirstName string `json:"first_name"`
  SecondName string `json:"second_name"`
}
type Order struct {
  Model
  CustomerID uint     `gorm:"references:Customer"`   // Order has a customer Id
}
type OrderProduct struct {
  OrderID uint   `json:"order_id"`        // orderProduct has Order Id and Product Id
  ProductID uint `json:"product_id"`
}
type Product struct {
  Model
  Code  string `json:"code"`
  Price uint `json:"price"`
}


func connectToDataBase()(db *gorm.DB,err error){
  dsn := "root:root@tcp(127.0.0.1:3306)/Orders?charset=utf8mb4&parseTime=True&loc=Local"
  db, err = gorm.Open(mysql.Open(dsn), &gorm.Config{})
  if err != nil {
    return nil,err
  }
  return db,nil
}
// Customer funcs
func createCustomer(firstName string,secondName string)(err error){
  db, err := connectToDataBase()
  if err != nil {
    return err
  }
  db.AutoMigrate(&Customer{})
  db.Create(&Customer{FirstName: firstName, SecondName: secondName})


  return nil
}
// Order funcs gorm.Model,CustomerID int
func createOrder(customerID uint)(err error){
  db, err := connectToDataBase()
  if err != nil {
    return err
  }
  db.AutoMigrate(&Order{})
  db.Create(&Order{CustomerID: customerID})

  return nil
}
//OrderProduct funcs OrderId int, ProductID int
func createOrderProduct(orderId uint,productId uint)(err error){
  db, err := connectToDataBase()
  if err != nil {
    return err
  }
  db.AutoMigrate(&OrderProduct{})
  db.Create(&OrderProduct{OrderID: orderId,ProductID: productId})
  return nil
}

  type Result struct {
    OrderID uint
    Code string
    Price uint
  }
func showOrderProduct(orderId uint)(results [] Result,err error){
  db, err := connectToDataBase()
  if err != nil {
    return nil,err
  }

 db.Model(&OrderProduct{}).Select("order_products.order_id,products.code,products.price").Where("Order_id = ?",orderId).Joins("join products on products.id = order_products.product_id").Scan(&results)
  if err != nil{
    return nil,err
  }
  return results,nil
}

func getOrdersById(orderId uint)(result_json []byte,err error){
  result,err :=showOrderProduct(orderId)
  if err != nil {
    return nil,err
  }
  result_json,err = json.Marshal(result)
  if err != nil {
    return nil,err
  }
  return result_json,nil
}

// there if func to show every order if you don't know your order ID
func getCustomerOrders(first_name string,second_name string)(user_orders []byte,err error){
  db, err := connectToDataBase()
  if err != nil {
  }
  var customer Customer // this is a var to save a customer and to get his id
  db.Where("first_name = ? AND second_name = ?",first_name,second_name).First(&customer)
  

  var orders [] Order //this is a slice for every customer's order
  result_user := make ([][]Result,10)

  db.Where("customer_id = ?",customer.ID).Find(&orders)

  for _, ord := range orders {
    result,err := showOrderProduct(ord.ID)
    if err != nil{
      fmt.Println("there is an error in 111") //TODO error
    }
    result_user[ord.ID-1] = result
  }
  user_orders,err = json.Marshal(result_user)
  if err != nil {
    return nil,err
  } 
  return user_orders,nil
}
//Product func  gorm.Model,Code string, Price uint
func createProduct(code string,price uint)(err error){
  db, err := connectToDataBase()
  if err != nil {
    return err
  }
  
  db.AutoMigrate(&Product{})
  db.Create(&Product{Code: code,Price: price})

  return nil
}
