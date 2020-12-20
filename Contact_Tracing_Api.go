package main

import (
	"fmt"
	
	"context"
	"encoding/json"
	"net/http"
	"time"
	"github.com/gorilla/mux"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
    "go.mongodb.org/mongo-driver/mongo/options"
   
)


var client *mongo.Client



type Marshaler interface {
    MarshalJSON() ([]byte, error)
}

type JSONTime time.Time

func (t JSONTime)MarshalJSON() ([]byte, error) {
    stamp := fmt.Sprintf("\"%s\"", time.Time(t).Format("Mon Jan _2"))
    return []byte(stamp), nil
}

type User struct {
	ID     primitive.ObjectID `json:"_id,omitempty" bson:"_id,omitempty"`
	Name   string             `json:"name,omitempty" bson:"name,omitempty"`
	DOB    string             `json:"DOB" bson:"dob,omitempty"`
	Mobile string             `json:"mobile" bson:"mobile,omitempty"`
	Email  string 			  `json:"email" bson: "email,omitempty"`
	Timestamp time.Time `json:"time_stamp" bson: "time_stamp"`
}

type Contact struct {
	ID        primitive.ObjectID `json:"_id,omitempty" bson:"_id,omitempty"`
	UserIdOne string `json:"idone,omitempty" bson:"idone,omitempty"`
	UserIdTwo string `json:"idtwo,omitempty" bson:"idtwo,omitempty"`
	Timestamp time.Time `json:"time_stamp" bson: "time_stamp"`
}

func Create_User(response http.ResponseWriter, request *http.Request) {
	response.Header().Set("content-type", "application/json")
	var user User
	_ = json.NewDecoder(request.Body).Decode(&user)

	user.Timestamp = time.Now()

	collection := client.Database("Contact_Tracing_Api").Collection("users")
	ctx, _ := context.WithTimeout(context.Background(), 5*time.Second)

	result, _ := collection.InsertOne(ctx, user)
	json.NewEncoder(response).Encode(result)
}

func Create_Contact(response http.ResponseWriter, request *http.Request) {
	response.Header().Set("content-type", "application/json")
	var contact Contact
	_ = json.NewDecoder(request.Body).Decode(&contact)

	contact.Timestamp = time.Now()

	collection := client.Database("Contact_Tracing_Api").Collection("contacts")
	ctx, _ := context.WithTimeout(context.Background(), 5*time.Second)

	result, _ := collection.InsertOne(ctx, contact)
	json.NewEncoder(response).Encode(result)
}

func Get_User(response http.ResponseWriter, request *http.Request) {
	response.Header().Set("content-type", "application/json")
	params := mux.Vars(request)

	id, _ := primitive.ObjectIDFromHex(params["id"])
	var user User
	collection := client.Database("Contact_Tracing_Api").Collection("users")
	ctx, _ := context.WithTimeout(context.Background(), 30*time.Second)

	err := collection.FindOne(ctx, User{ID: id}).Decode(&user)
	if err != nil {
		response.WriteHeader(http.StatusInternalServerError)
		response.Write([]byte(`{ "message": "` + err.Error() + `" }`))
		return
	}
	json.NewEncoder(response).Encode(user)
}

func Get_All_Users(response http.ResponseWriter, request *http.Request) {
	response.Header().Set("content-type", "application/json")
	var users []User
	collection := client.Database("Contact_Tracing_Api").Collection("users")
	ctx, _ := context.WithTimeout(context.Background(), 30*time.Second)
	cursor, err := collection.Find(ctx, bson.M{})
	if err != nil {
		response.WriteHeader(http.StatusInternalServerError)
		response.Write([]byte(`{ "message": "` + err.Error() + `" }`))
		return
	}
	defer cursor.Close(ctx)
	for cursor.Next(ctx) {
		var user User
		cursor.Decode(&user)
		users = append(users, user)
	}
	if err := cursor.Err(); err != nil {
		response.WriteHeader(http.StatusInternalServerError)
		response.Write([]byte(`{ "message": "` + err.Error() + `" }`))
		return
	}
	json.NewEncoder(response).Encode(users)
}








func Get_Contact(response http.ResponseWriter, request *http.Request) {
	response.Header().Set("content-type", "application/json")
	params := mux.Vars(request)

	id, _ := primitive.ObjectIDFromHex(params["id"])
	var contact Contact
	collection := client.Database("Contact_Tracing_Api").Collection("contacts")
	ctx, _ := context.WithTimeout(context.Background(), 30*time.Second)
	
	err := collection.FindOne(ctx, Contact{ID: id}).Decode(&contact)
	if err != nil {
		response.WriteHeader(http.StatusInternalServerError)
		response.Write([]byte(`{ "message": "` + err.Error() + `" }`))
		return
	}
	json.NewEncoder(response).Encode(contact)
}

func Get_All_Contacts(response http.ResponseWriter, request *http.Request) {
	response.Header().Set("content-type", "application/json")
	var contacts []Contact
	collection := client.Database("Contact_Tracing_Api").Collection("contacts")
	ctx, _ := context.WithTimeout(context.Background(), 30*time.Second)
	cursor, err := collection.Find(ctx, bson.M{})
	if err != nil {
		response.WriteHeader(http.StatusInternalServerError)
		response.Write([]byte(`{ "message": "` + err.Error() + `" }`))
		return
	}
	defer cursor.Close(ctx)
	for cursor.Next(ctx) {
		var contact Contact
		cursor.Decode(&contact)
		contacts = append(contacts, contact)
	}
	if err := cursor.Err(); err != nil {
		response.WriteHeader(http.StatusInternalServerError)
		response.Write([]byte(`{ "message": "` + err.Error() + `" }`))
		return
	}
	json.NewEncoder(response).Encode(contacts)
}

func main() {
	
	
	ctx, _ := context.WithTimeout(context.Background(), 10*time.Second)
	clientOptions := options.Client().ApplyURI("mongodb://localhost:27017")
	client, _ = mongo.Connect(ctx, clientOptions)

	router := mux.NewRouter()

	router.HandleFunc("/users", Get_All_Users).Methods("GET")
	router.HandleFunc("/users", Create_User).Methods("POST")

	router.HandleFunc("/users/{id}", Get_User).Methods("GET")

	router.HandleFunc("/contacts", Get_All_Contacts).Methods("GET")
	router.HandleFunc("/contacts", Create_Contact).Methods("POST")

	router.HandleFunc("/contacts?user={id}&infection_timestamp={ts}", Get_Contact).Methods("GET")

	http.ListenAndServe(":9090", router)
}

