package main

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"net/http"

	"github.com/gorilla/mux"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

func main() {
	//config := config.Config()

	route := mux.NewRouter()
	s := route.PathPrefix("/api").Subrouter() //Base Path
	//Routes
	//s := route
	//user api
	s.HandleFunc("/createProfile", createProfile).Methods("POST")
	s.HandleFunc("/getAllUsers", getAllUsers).Methods("GET")
	s.HandleFunc("/getUserProfile", getUserProfile).Methods("POST")
	s.HandleFunc("/updateProfile", updateProfile).Methods("PUT")
	s.HandleFunc("/deleteProfile/{id}", deleteProfile).Methods("DELETE")
	//state api
	s.HandleFunc("/stateApi", stateApi).Methods("GET")
	s.HandleFunc("/statesaveApi", statesaveApi).Methods("POST")
	s.HandleFunc("/stategetsingleApi", stategetsingleApi).Methods("POST")
	s.HandleFunc("/stateupdateApi", stateupdateApi).Methods("PUT")
	s.HandleFunc("/deletestateApi", deletestateApi).Methods("DELETE")
	s.HandleFunc("/enablestateApi", enablestateApi).Methods("PUT")
	s.HandleFunc("/disablestateApi", disablestateApi).Methods("PUT")
	s.HandleFunc("/filterstateApi", filterstateApi).Methods("POST")
	//district api
	s.HandleFunc("/districtApi/{id}", districtApi).Methods("GET")
	s.HandleFunc("/districtsaveApi", districtsaveApi).Methods("POST")
	s.HandleFunc("/districtgetsingleApi", districtgetsingleApi).Methods("POST")
	s.HandleFunc("/districtupdateApi", districtupdateApi).Methods("PUT")
	s.HandleFunc("/deletedistrictApi", deletedistrictApi).Methods("DELETE")
	s.HandleFunc("/enabledistrictApi", enabledistrictApi).Methods("PUT")
	s.HandleFunc("/disabledistrictApi", disabledistrictApi).Methods("PUT")
	s.HandleFunc("/filterdistrictApi", filterdistrictApi).Methods("POST")
	//block api
	s.HandleFunc("/blockApi/{id}", blockApi).Methods("GET")
	//grampanchayat api
	s.HandleFunc("/grampanchayatApi/{id}", grampanchayatApi).Methods("GET")
	//village api
	s.HandleFunc("/villageApi/{id}", villageApi).Methods("GET")
	log.Fatal(http.ListenAndServe(":8000", s)) // Run Server
}

//db connection
func db() *mongo.Client {
	clientOptions := options.Client().ApplyURI("mongodb://localhost:27018,localhost:27019,localhost:27020/?replicaSet=rsSample") // Connect to //MongoDB
	client, err := mongo.Connect(context.TODO(), clientOptions)
	if err != nil {
		log.Fatal(err)
	}
	// Check the connection
	//ping is library to connect to ip
	err = client.Ping(context.TODO(), nil)
	if err != nil {
		log.Fatal(err)
	}
	fmt.Println("Connected to MongoDB!")
	return client
}

// struct for storing data
type user struct {
	Name string `json:name`
	Age  int    `json:age`
	City string `json:city`
}

type State struct {
	ID           primitive.ObjectID   `json:"id" form:"id," bson:"_id,omitempty"`
	Name         string               `json:"name" bson:"name,omitempty"`
	Status       string               `json:"status" bson:"status,omitempty"`
	ActiveStatus bool                 `json:"activeStatus" bson:"activeStatus,omitempty"`
	Languages    []primitive.ObjectID `json:"languages"  bson:"languages,omitempty"`
	Version      int64                `json:"version"  bson:"version,omitempty"`
}

type RefState struct {
	State `bson:",inline"`
	Ref   struct {
	} `json:"ref"  bson:"ref,omitempty"`
}
type StateFilter struct {
	ID           []primitive.ObjectID `json:"id" form:"id," bson:"_id,omitempty"`
	Name         []string             `json:"name" bson:"name,omitempty"`
	Status       []string             `json:"status" bson:"status,omitempty"`
	ActiveStatus []bool               `json:"activeStatus" bson:"activeStatus,omitempty"`
	Languages    []primitive.ObjectID `json:"languages"  bson:"languages,omitempty"`
	Version      int64                `json:"version"  bson:"version,omitempty"`
	OmitIdState  struct {
		ID []primitive.ObjectID `json:"id" form:"id," bson:"_id,omitempty"`
	} `json:"omitIdState"`
	Regex struct {
		Name string `json:"name" bson:"name"`
	} `json:"regex" bson:"regex"`
}

type District struct {
	ID                  primitive.ObjectID   `json:"id" form:"id," bson:"_id,omitempty"`
	State               primitive.ObjectID   `json:"state"  bson:"state,omitempty"`
	Name                string               `json:"name" bson:"name,omitempty"`
	Status              string               `json:"status" bson:"status,omitempty"`
	ActiveStatus        bool                 `json:"activeStatus" bson:"activeStatus,omitempty"`
	AgroEcologicalZones []primitive.ObjectID `json:"agroEcologicalZones"  bson:"agroEcologicalZones,omitempty"`
	SoilTypes           []primitive.ObjectID `json:"soilTypes"  bson:"soilTypes,omitempty"`
	Version             int64                `json:"version"  bson:"version,omitempty"`
}

type DistrictFilter struct {
	ID             []primitive.ObjectID `json:"id" form:"id," bson:"_id,omitempty"`
	State          []primitive.ObjectID `json:"state"  bson:"state,omitempty"`
	ActiveStatus   []bool               `json:"activeStatus" bson:"activeStatus,omitempty"`
	Status         []string             `json:"status" bson:"status,omitempty"`
	OmitIddistrict struct {
		ID []primitive.ObjectID `json:"id" form:"id," bson:"_id,omitempty"`
	}
	Regex struct {
		Name string `json:"name" bson:"name"`
	} `json:"regex" bson:"regex"`
}

type RefDistrict struct {
	District `bson:",inline"`
	Ref      struct {
		State State `json:"state,omitempty" bson:"state,omitempty"`
	} `json:"ref,omitempty" bson:"ref,omitempty"`
}

type Block struct {
	ID           primitive.ObjectID `json:"id" form:"id," bson:"_id,omitempty"`
	District     primitive.ObjectID `json:"district"  bson:"district,omitempty"`
	Name         string             `json:"name" bson:"name,omitempty"`
	ActiveStatus bool               `json:"activeStatus" bson:"activeStatus,omitempty"`
	Version      int64              `json:"version"  bson:"version,omitempty"`
}

type GramPanchayat struct {
	ID           primitive.ObjectID `json:"id" form:"id," bson:"_id,omitempty"`
	Block        primitive.ObjectID `json:"block"  bson:"block,omitempty"`
	Name         string             `json:"name" bson:"name,omitempty"`
	ActiveStatus bool               `json:"activeStatus" bson:"activeStatus,omitempty"`
	Version      int64              `json:"version"  bson:"version,omitempty"`
}

type Village struct {
	ID            primitive.ObjectID `json:"id" form:"id," bson:"_id,omitempty"`
	GramPanchayat primitive.ObjectID `json:"gramPanchayat"  bson:"gramPanchayat,omitempty"`
	Name          string             `json:"name" bson:"name,omitempty"`
	ActiveStatus  bool               `json:"activeStatus" bson:"activeStatus,omitempty"`
	Location      struct {
		Longitude float64 `json:"longitude" bson:"longitude,omitempty"`
		Latitude  float64 `json:"latitude" bson:"latitude,omitempty"`
	} `json:"location" bson:"location,omitempty"`
	Population float64 `json:"population"  bson:"population,omitempty"`
	FieldAgent string  `json:"fieldAgent"  bson:"fieldAgent,omitempty"`
	Version    int64   `json:"version"  bson:"version,omitempty"`
}

var userCollection = db().Database("gopractice").Collection("users") // get collection "users" from db() which returns *mongo.Client

var blockCollection = db().Database("gopractice").Collection("block")

var districtCollection = db().Database("gopractice").Collection("district")

var gramPanchayatCollection = db().Database("gopractice").Collection("gramPanchayat")

var stateCollection = db().Database("gopractice").Collection("state")

var villageCollection = db().Database("gopractice").Collection("village")

//State Collection Constants
const (
	STATESTATUSACTIVE   = "Active"
	STATESTATUSDISABLED = "Disabled"
	STATESTATUSDELETED  = "Deleted"
)

//District Collection Constants
const (
	DISTRICTSTATUSACTIVE   = "Active"
	DISTRICTSTATUSDISABLED = "Disabled"
	DISTRICTSTATUSDELETED  = "Deleted"
)

// Create Profile or add new value

func createProfile(w http.ResponseWriter, r *http.Request) {

	w.Header().Set("Content-Type", "application/json") // for adding Content-type

	var person user
	err := json.NewDecoder(r.Body).Decode(&person) // storing in person variable of type user
	if err != nil {
		fmt.Print(err)
	}
	insertResult, err := userCollection.InsertOne(context.TODO(), person)
	if err != nil {
		log.Fatal(err)
	}

	fmt.Println("Inserted a single document: ", insertResult)
	//insertResult.InsertedID
	//map[type]type
	var x interface{} = insertResult.InsertedID
	str := fmt.Sprintf("%v", x)
	m := make(map[string]interface{})
	m[str] = person
	fmt.Println(m)
	dataB, err := json.Marshal(m)
	if err != nil {
		log.Fatal(err)
	}
	//json.NewEncoder(w).Encode(m)

	w.Write(dataB)

	// return the mongodb ID of generated document

}

func getAllUsers(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	var results []primitive.M                                   //slice for multiple documents
	cur, err := userCollection.Find(context.TODO(), bson.D{{}}) //returns a *mongo.Cursor
	if err != nil {

		fmt.Println(err)

	}
	for cur.Next(context.TODO()) { //Next() gets the next document for corresponding cursor

		var elem primitive.M
		err := cur.Decode(&elem)
		if err != nil {
			log.Fatal(err)
		}

		results = append(results, elem) // appending document pointed by Next()
	}
	cur.Close(context.TODO()) // close the cursor once stream of documents has exhausted
	json.NewEncoder(w).Encode(results)
}

// Get Profile of a particular User by Name

func getUserProfile(w http.ResponseWriter, r *http.Request) {

	w.Header().Set("Content-Type", "application/json")

	var body user
	e := json.NewDecoder(r.Body).Decode(&body)
	if e != nil {

		fmt.Print(e)
	}
	var result primitive.M //  an unordered representation of a BSON document which is a Map
	err := userCollection.FindOne(context.TODO(), bson.D{{"name", body.Name}}).Decode(&result)
	if err != nil {

		fmt.Println(err)

	}
	json.NewEncoder(w).Encode(result) // returns a Map containing document

}

//Update Profile of User

func updateProfile(w http.ResponseWriter, r *http.Request) {

	w.Header().Set("Content-Type", "application/json")

	type updateBody struct {
		Name string `json:"name"` //value that has to be matched
		City string `json:"city"` // value that has to be modified
	}
	var body updateBody
	e := json.NewDecoder(r.Body).Decode(&body)
	if e != nil {

		fmt.Print(e)
	}
	filter := bson.D{{"name", body.Name}} // converting value to BSON type
	after := options.After                // for returning updated document
	returnOpt := options.FindOneAndUpdateOptions{

		ReturnDocument: &after,
	}
	update := bson.D{{"$set", bson.D{{"city", body.City}}}}
	updateResult := userCollection.FindOneAndUpdate(context.TODO(), filter, update, &returnOpt)

	var result primitive.M
	_ = updateResult.Decode(&result)

	json.NewEncoder(w).Encode(result)
}

//Delete Profile of User

func deleteProfile(w http.ResponseWriter, r *http.Request) {

	w.Header().Set("Content-Type", "application/json")
	params := mux.Vars(r)["id"] //get Parameter value as string

	_id, err := primitive.ObjectIDFromHex(params) // convert params to mongodb Hex ID
	if err != nil {
		fmt.Printf(err.Error())
	}
	opts := options.Delete().SetCollation(&options.Collation{}) // to specify language-specific rules for string comparison, such as rules for lettercase
	res, err := userCollection.DeleteOne(context.TODO(), bson.D{{"_id", _id}}, opts)
	if err != nil {
		log.Fatal(err)
	}
	fmt.Printf("deleted %v documents\n", res.DeletedCount)
	json.NewEncoder(w).Encode(_id) // return number of documents deleted

}

//state api
func stateApi(w http.ResponseWriter, r *http.Request) {

	w.Header().Set("Content-Type", "application/json")
	var results []primitive.M                                    //slice for multiple documents
	cur, err := stateCollection.Find(context.TODO(), bson.D{{}}) //returns a *mongo.Cursor
	if err != nil {

		fmt.Println(err)

	}
	for cur.Next(context.TODO()) { //Next() gets the next document for corresponding cursor

		var elem primitive.M
		err := cur.Decode(&elem)
		if err != nil {
			log.Fatal(err)
		}

		results = append(results, elem) // appending document pointed by Next()
	}
	cur.Close(context.TODO()) // close the cursor once stream of documents has exhausted
	json.NewEncoder(w).Encode(results)

}

//state save api
func statesaveApi(w http.ResponseWriter, r *http.Request) {

	w.Header().Set("Content-Type", "application/json")

	var statevalue State
	fmt.Println(statevalue)
	err := json.NewDecoder(r.Body).Decode(&statevalue) // storing in person variable of type statevalue
	if err != nil {
		fmt.Print(err)
	}
	insertResult, err := stateCollection.InsertOne(context.TODO(), statevalue)
	if err != nil {
		log.Fatal(err)
	}

	json.NewEncoder(w).Encode(insertResult.InsertedID)
}

//state get single api
func stategetsingleApi(w http.ResponseWriter, r *http.Request) {

	w.Header().Set("Content-Type", "application/json")

	var body State
	//body.ID
	e := json.NewDecoder(r.Body).Decode(&body)
	if e != nil {

		fmt.Print(e)
	}
	var result []State //  an unordered representation of a BSON document which is a Map
	err := stateCollection.FindOne(context.TODO(), bson.D{{"_id", body.ID}}).Decode(&result)
	// fmt.Println(body.ID)
	// fmt.Println(result)
	if err != nil {

		fmt.Println(err)

	}
	json.NewEncoder(w).Encode(result)
}

//Update state api
func stateupdateApi(w http.ResponseWriter, r *http.Request) {

	w.Header().Set("Content-Type", "application/json")

	var body State
	e := json.NewDecoder(r.Body).Decode(&body)
	if e != nil {

		fmt.Print(e)
	}
	filter := bson.D{{"_id", body.ID}} // converting value to BSON type
	after := options.After             // for returning updated document
	returnOpt := options.FindOneAndUpdateOptions{

		ReturnDocument: &after,
	}

	//update := bson.M{{"$set": bson.M{body}}}
	update := bson.D{{"$set", bson.D{{"name", body.Name}}}}
	updateResult := stateCollection.FindOneAndUpdate(context.TODO(), filter, update, &returnOpt)

	var result State
	_ = updateResult.Decode(&result)

	json.NewEncoder(w).Encode(result)
}

//delete state api
func deletestateApi(w http.ResponseWriter, r *http.Request) {

	w.Header().Set("Content-Type", "application/json")
	var body State
	e := json.NewDecoder(r.Body).Decode(&body)
	if e != nil {

		fmt.Print(e)
	}
	filter := bson.D{{"_id", body.ID}} // converting value to BSON type
	after := options.After             // for returning updated document
	returnOpt := options.FindOneAndUpdateOptions{

		ReturnDocument: &after,
	}

	//update := bson.M{{"$set": bson.M{body}}}
	update := bson.M{"$set": bson.M{"status": STATESTATUSDELETED, "activeStatus": false}}
	updateResult := stateCollection.FindOneAndUpdate(context.TODO(), filter, update, &returnOpt)

	var result State
	_ = updateResult.Decode(&result)

	json.NewEncoder(w).Encode(result)

}

//Enable state api
func enablestateApi(w http.ResponseWriter, r *http.Request) {

	w.Header().Set("Content-Type", "application/json")
	var body State
	e := json.NewDecoder(r.Body).Decode(&body)
	if e != nil {

		fmt.Print(e)
	}
	filter := bson.D{{"_id", body.ID}} // converting value to BSON type
	after := options.After             // for returning updated document
	returnOpt := options.FindOneAndUpdateOptions{

		ReturnDocument: &after,
	}

	//update := bson.M{{"$set": bson.M{body}}}
	update := bson.M{"$set": bson.M{"status": STATESTATUSACTIVE, "activeStatus": true}}
	updateResult := stateCollection.FindOneAndUpdate(context.TODO(), filter, update, &returnOpt)

	var result State
	_ = updateResult.Decode(&result)

	json.NewEncoder(w).Encode(result)

}

//Disable state api
func disablestateApi(w http.ResponseWriter, r *http.Request) {

	w.Header().Set("Content-Type", "application/json")
	var body State
	e := json.NewDecoder(r.Body).Decode(&body)
	if e != nil {

		fmt.Print(e)
	}
	filter := bson.D{{"_id", body.ID}} // converting value to BSON type
	after := options.After             // for returning updated document
	returnOpt := options.FindOneAndUpdateOptions{

		ReturnDocument: &after,
	}

	//update := bson.M{{"$set": bson.M{body}}}
	update := bson.M{"$set": bson.M{"status": STATESTATUSDISABLED, "activeStatus": false}}
	updateResult := stateCollection.FindOneAndUpdate(context.TODO(), filter, update, &returnOpt)

	var result State
	_ = updateResult.Decode(&result)

	json.NewEncoder(w).Encode(result)

}

//filter state api
func filterstateApi(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	var body StateFilter
	e := json.NewDecoder(r.Body).Decode(&body)
	if e != nil {

		fmt.Print(e)
	}

	var result []RefState //  an unordered representation of a BSON document which is a Map
	mainPipeline := []bson.M{}
	query := []bson.M{}

	if len(body.ActiveStatus) > 0 {
		query = append(query, bson.M{"activeStatus": bson.M{"$in": body.ActiveStatus}})
	}
	if len(body.Status) > 0 {
		query = append(query, bson.M{"status": bson.M{"$in": body.Status}})
	}
	//Regex
	if body.Regex.Name != "" {
		query = append(query, bson.M{"name": primitive.Regex{Pattern: body.Regex.Name, Options: "xi"}})
	}
	if len(body.OmitIdState.ID) > 0 {
		query = append(query, bson.M{"_id": bson.M{"$in": body.OmitIdState.ID}})
	}

	if len(query) > 0 {
		mainPipeline = append(mainPipeline, bson.M{"$match": bson.M{"$and": query}})
	}

	//Aggregation
	//d.Shared.BsonToJSONPrintTag("state query =>", mainPipeline)
	cursor, err := stateCollection.Aggregate(context.TODO(), mainPipeline, nil)
	if err != nil {
		// return nil, err
		w.WriteHeader(500)
		return
	}
	//var states []models.RefState
	if err = cursor.All(context.TODO(), &result); err != nil {
		//return nil, err
		w.WriteHeader(500)
		return
	}
	//return result, nil
	err = json.NewEncoder(w).Encode(result)
	if err != nil {
		w.WriteHeader(500)
		return
	}
	w.WriteHeader(200)

}

//district api
func districtApi(w http.ResponseWriter, r *http.Request) {

	w.Header().Set("Content-Type", "application/json")
	params := mux.Vars(r)["id"] //get Parameter value as string

	id, err := primitive.ObjectIDFromHex(params) // convert params to mongodb Hex ID
	if err != nil {
		fmt.Printf(err.Error())
	}
	var result []District
	cur, err1 := districtCollection.Find(context.TODO(), bson.D{{"state", id}})
	if err1 != nil {

		fmt.Println(err1)

	}
	cur.All(context.TODO(), &result)
	json.NewEncoder(w).Encode(result)
}

//district save api
func districtsaveApi(w http.ResponseWriter, r *http.Request) {

	w.Header().Set("Content-Type", "application/json")

	var statevalue District
	fmt.Println(statevalue)
	err := json.NewDecoder(r.Body).Decode(&statevalue) // storing in person variable of type statevalue
	if err != nil {
		fmt.Print(err)
	}
	insertResult, err := districtCollection.InsertOne(context.TODO(), statevalue)
	if err != nil {
		log.Fatal(err)
	}

	json.NewEncoder(w).Encode(insertResult.InsertedID)
}

//district get single api
func districtgetsingleApi(w http.ResponseWriter, r *http.Request) {

	w.Header().Set("Content-Type", "application/json")

	var body District
	e := json.NewDecoder(r.Body).Decode(&body)
	if e != nil {

		fmt.Print(e)
	}
	var result District //  an unordered representation of a BSON document which is a Map
	err := districtCollection.FindOne(context.TODO(), bson.D{{"_id", body.ID}}).Decode(&result)
	if err != nil {

		fmt.Println(err)

	}
	json.NewEncoder(w).Encode(result)
}

//update district api
func districtupdateApi(w http.ResponseWriter, r *http.Request) {

	w.Header().Set("Content-Type", "application/json")

	var body District
	e := json.NewDecoder(r.Body).Decode(&body)
	if e != nil {

		fmt.Print(e)
	}
	filter := bson.D{{"_id", body.ID}} // converting value to BSON type
	after := options.After             // for returning updated document
	returnOpt := options.FindOneAndUpdateOptions{

		ReturnDocument: &after,
	}

	//update := bson.M{{"$set": bson.M{body}}}
	update := bson.D{{"$set", bson.D{{"name", body.Name}}}}
	updateResult := districtCollection.FindOneAndUpdate(context.TODO(), filter, update, &returnOpt)

	var result State
	_ = updateResult.Decode(&result)

	json.NewEncoder(w).Encode(result)
}

//delete district api
func deletedistrictApi(w http.ResponseWriter, r *http.Request) {

	w.Header().Set("Content-Type", "application/json")
	var body District
	e := json.NewDecoder(r.Body).Decode(&body)
	if e != nil {

		fmt.Print(e)
	}
	filter := bson.D{{"_id", body.ID}} // converting value to BSON type
	after := options.After             // for returning updated document
	returnOpt := options.FindOneAndUpdateOptions{

		ReturnDocument: &after,
	}

	update := bson.M{"$set": bson.M{"status": DISTRICTSTATUSDELETED, "activeStatus": false}}
	updateResult := districtCollection.FindOneAndUpdate(context.TODO(), filter, update, &returnOpt)

	var result District
	_ = updateResult.Decode(&result)

	json.NewEncoder(w).Encode(result)

}

//Enable district api
func enabledistrictApi(w http.ResponseWriter, r *http.Request) {

	w.Header().Set("Content-Type", "application/json")
	var body District
	e := json.NewDecoder(r.Body).Decode(&body)
	if e != nil {
		fmt.Print("e")
	}

	filter := bson.D{{"_id", body.ID}} // converting value to BSON type
	after := options.After             // for returning updated document
	returnOpt := options.FindOneAndUpdateOptions{

		ReturnDocument: &after,
	}

	update := bson.M{"$set": bson.M{"status": DISTRICTSTATUSACTIVE, "activeStatus": true}}
	updateResult := districtCollection.FindOneAndUpdate(context.TODO(), filter, update, &returnOpt)

	var result District
	_ = updateResult.Decode(&result)

	json.NewEncoder(w).Encode(result)

}

//disable district api
func disabledistrictApi(w http.ResponseWriter, r *http.Request) {

	w.Header().Set("Content-Type", "application/json")

	var body District
	e := json.NewDecoder(r.Body).Decode(&body)
	if e != nil {
		fmt.Print("e")
	}

	filter := bson.D{{"_id", body.ID}} // converting value to BSON type
	after := options.After             // for returning updated document
	returnOpt := options.FindOneAndUpdateOptions{

		ReturnDocument: &after,
	}

	update := bson.M{"$set": bson.M{"status": DISTRICTSTATUSDISABLED, "activeStatus": false}}
	updateResult := districtCollection.FindOneAndUpdate(context.TODO(), filter, update, &returnOpt)

	var result District
	_ = updateResult.Decode(&result)

	json.NewEncoder(w).Encode(result)
}

//filter district api
func filterdistrictApi(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	var body DistrictFilter
	e := json.NewDecoder(r.Body).Decode(&body)
	if e != nil {

		fmt.Print(e)
	}

	var result []RefDistrict //  an unordered representation of a BSON document which is a Map
	mainPipeline := []bson.M{}
	query := []bson.M{}

	if len(body.ActiveStatus) > 0 {
		query = append(query, bson.M{"activeStatus": bson.M{"$in": body.ActiveStatus}})
	}
	if len(body.Status) > 0 {
		query = append(query, bson.M{"status": bson.M{"$in": body.Status}})
	}
	//Regex
	if body.Regex.Name != "" {
		query = append(query, bson.M{"name": primitive.Regex{Pattern: body.Regex.Name, Options: "xi"}})
	}
	if len(body.OmitIddistrict.ID) > 0 {
		query = append(query, bson.M{"_id": bson.M{"$in": body.OmitIddistrict.ID}})
	}

	if len(query) > 0 {
		mainPipeline = append(mainPipeline, bson.M{"$match": bson.M{"$and": query}})
	}
	// type RefDistrict struct {
	// 	District `bson:",inline"`
	// 	Ref      struct {
	// 		State State `json:"state,omitempty" bson:"state,omitempty"`
	// 	} `json:"ref,omitempty" bson:"ref,omitempty"`
	// }

	//ref
	mainPipeline = append(mainPipeline,
		bson.M{"$lookup": bson.M{"from": "state", "as": "ref.state", "localField": "state", "foreignField": "_id"}},
		bson.M{"$addFields": bson.M{"ref.state": bson.M{"$arrayElemAt": []interface{}{"$ref.state", 0}}}})

	//Aggregation
	//d.Shared.BsonToJSONPrintTag("state query =>", mainPipeline)
	cursor, err := districtCollection.Aggregate(context.TODO(), mainPipeline, nil)
	if err != nil {
		// return nil, err
		w.WriteHeader(500)
		return
	}
	//var states []models.RefState
	if err = cursor.All(context.TODO(), &result); err != nil {
		//return nil, err
		w.WriteHeader(500)
		return
	}
	//return result, nil
	err = json.NewEncoder(w).Encode(result)
	if err != nil {
		w.WriteHeader(500)
		return
	}
	w.WriteHeader(200)

}

//block api
func blockApi(w http.ResponseWriter, r *http.Request) {

	w.Header().Set("Content-Type", "application/json")
	params := mux.Vars(r)["id"] //get Parameter value as string

	id, err := primitive.ObjectIDFromHex(params) // convert params to mongodb Hex ID
	if err != nil {
		fmt.Printf(err.Error())
	}
	var result []Block
	cur, err1 := blockCollection.Find(context.TODO(), bson.D{{"district", id}})
	if err1 != nil {

		fmt.Println(err1)

	}
	cur.All(context.TODO(), &result)
	json.NewEncoder(w).Encode(result)

}

//gramPanchayat api
func grampanchayatApi(w http.ResponseWriter, r *http.Request) {

	w.Header().Set("Content-Type", "application/json")
	params := mux.Vars(r)["id"] //get Parameter value as string

	id, err := primitive.ObjectIDFromHex(params) // convert params to mongodb Hex ID
	if err != nil {
		fmt.Printf(err.Error())
	}
	var result []GramPanchayat
	cur, err1 := gramPanchayatCollection.Find(context.TODO(), bson.D{{"block", id}})
	if err1 != nil {

		fmt.Println(err1)

	}
	cur.All(context.TODO(), &result)
	json.NewEncoder(w).Encode(result)

}

//village api
func villageApi(w http.ResponseWriter, r *http.Request) {

	w.Header().Set("Content-Type", "application/json")
	fmt.Println("start")
	params := mux.Vars(r)["id"] //get Parameter value as string
	fmt.Println(params)
	id, err := primitive.ObjectIDFromHex(params) // convert params to mongodb Hex ID
	if err != nil {
		fmt.Printf(err.Error())
	}
	var result []Village
	cur, err1 := villageCollection.Find(context.TODO(), bson.D{{"gramPanchayat", id}})
	if err1 != nil {

		fmt.Println(err1)

	}
	//	fmt.Println(cur)
	if err := cur.All(context.TODO(), &result); err != nil {
		fmt.Println("Err in reading cursor", err.Error())
		return
	}
	json.NewEncoder(w).Encode(result)

}
