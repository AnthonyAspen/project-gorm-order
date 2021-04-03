package main

import (
  "net/http"
  "log"
)

func addHandler(w http.ResponseWriter,r *http.Request){
  

}
func getCustomerOrdersHandler(w http.ResponseWriter,r *http.Request){
  first_name := r.URL.Query().Get("first")
  second_name := "Customer"


  user_orders,err := getCustomerOrders(first_name,second_name)
  if err != nil{
   log.Println(err)
    w.WriteHeader(http.StatusNotFound)
  }

  
  w.Header().Add("Access-Control-Allow-Origin", "*")
  w.Header().Add("Access-Control-Allow-Methods", "*")
  w.Header().Add("Access-Control-Allow-Headers", "*")
  w.Write(user_orders)

}
func main(){
  http.HandleFunc("/add",addHandler)
  http.HandleFunc("/get",getCustomerOrdersHandler)
  log.Fatal(http.ListenAndServe(":8080",nil))
}
