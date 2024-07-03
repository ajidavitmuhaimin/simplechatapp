package main

//import "fmt"
import "net/http"
import "github.com/gorilla/websocket"

type Message struct{
   Username string `json:"username"`
   Message string `json:"message"`
}

var upgrader=websocket.Upgrader{
   ReadBufferSize: 1024,
   WriteBufferSize: 1024,
   CheckOrigin: func(r *http.Request)bool{
      return true
   },
}

var clients=make(map[*websocket.Conn]bool)
var broadcast=make(chan Message)

func main(){

   //http.HandleFunc("/send",sendHandler)
   http.HandleFunc("/conn",connHandler)
   go sendHandler()
   http.ListenAndServe(":1395",nil)

}


func connHandler(w http.ResponseWriter, r *http.Request){
   ws,err:=upgrader.Upgrade(w,r,nil)
   if err!=nil{
      panic(err)
   }
   clients[ws]=true
   //loop forever
   for{
      var msg Message
      err:=ws.ReadJSON(&msg)
      if err!=nil{
         panic(err)
      }
      broadcast <-msg
   }

}


func sendHandler(){
   
   for{
      msg:= <-broadcast
      for client:=range clients{
         err:=client.WriteJSON(msg)
         if err!=nil{
            panic(err)
         }
      }
   }

}
