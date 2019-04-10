package main

import "fmt"
import "log"
import "time"
import "context"
import "go.mongodb.org/mongo-driver/mongo"
import "go.mongodb.org/mongo-driver/mongo/options"
import "go.mongodb.org/mongo-driver/bson"
import "go.mongodb.org/mongo-driver/bson/primitive"

//	Documentacion: https://godoc.org/go.mongodb.org/mongo-driver/mongo
//	Referencia:    https://vkt.sh/go-mongodb-driver-cookbook/

func main()  {
	//	Creamos cliente
	client, err := mongo.NewClient(options.Client().ApplyURI("mongodb://mongo:27017"))
	if err != nil {
		log.Fatal("[ERROR]: No se pudo crear cliente\n", err)
	}
	//	Contexto de conexion (timeout de 10 segundos)
	ctx, _ := context.WithTimeout(context.Background(), 10*time.Second)
	
	//	Realizamos conexion del cliente con la base de datos
	//	Nota: la asignacion no lleva ":=" debido 
	//	que la variable err ya existia con anterioridad
	err = client.Connect(ctx)
	if err != nil {
		log.Fatal("[ERROR]: No se pudo conectar el cliente\n", err)
	}

	//	Accedemos a coleccion dentro de la base de datos
	collection := client.Database("testing").Collection("numbers")

	//	Eliminamos coleccion en base de datos
	err = collection.Drop(ctx)

	//	Insertamos datos en la coleccion
	res, err := collection.InsertOne(ctx, bson.M{"name": "uno", "value": 1.0})
	if err != nil {
		log.Fatal("[ERROR]: Insertar en coleccion\n", err)
	}

	//	Insertamos multiples datos en la coleccion
	multires, err := collection.InsertMany(ctx, []interface{}{
		bson.M{"name": "dos",    "value": 2.0},
		bson.M{"name": "tres",   "value": 3.0},
		bson.M{"name": "cuatro", "value": 4.0},
		bson.M{"name": "cinco",  "value": 5.0},
	})
	if err != nil {
		log.Fatal("[ERROR]: Insertar multiple en coleccion\n", err)
	}
	fmt.Println("-- MULTIRES --")
	fmt.Println(multires)

	//	Obtenemos ID de Dato insertado
	id := res.InsertedID
	fmt.Println(id)

	//	Multiples queries pueden ser realizados con un cursor
	//	para obtener un cursor se debe hacer...
	cur, err := collection.Find(ctx, bson.D{})
	if err != nil { 
		log.Fatal("[ERROR]: Obtener cursor",err) 
	}
	//	Podemos delegar el cierre del cursor al finalizar la ejecucion
	//	haciendo uso del modificador "defer"
	defer cur.Close(ctx)

	fmt.Println("-- FIND --")
	//	Accedemos a datos de cursor
	for cur.Next(ctx) {
		var result bson.M
		err := cur.Decode(&result)
		if err != nil {
			log.Fatal("[ERROR]: Al decodificar resultado de cursor",err) 
		}
		//	El resultado se decodifica en un MAP
		//  fmt.Println(result)
		//	Para acceder a un elemento del MAP se usan [llave]
		fmt.Printf( "%s -> %.4f \n", result["name"], result["value"] )
	}
	if err := cur.Err(); err != nil {
		log.Fatal("[ERROR]: Cursor")
	}

	//	Se define una estructura de resultado
	var algunDato struct {
		ID 		primitive.ObjectID 	`bson:"_id"`
		Name 	string 				`bson:"name"`
		Value 	float64 			`bson:"value"`
	}

	//	Se define filtro de busqueda
	filter := bson.M{"name" : "uno"}
	err = collection.FindOne(ctx, filter).Decode(&algunDato)
	if err != nil {
		log.Fatal("[ERROR]: Al buscar un elemento")
	}
	fmt.Println("-- FIND ONE --")
	fmt.Printf( "%s -> %.4f \n", algunDato.Name, algunDato.Value)

	//	Actualizamos algun elemento
	filter  = bson.M{"_id": algunDato.ID}
	update := bson.M{"$set": bson.M{"name" : "_uno"}}

	updateres, err := collection.UpdateOne(ctx, filter, update)
	if err != nil {
		log.Fatal("[ERROR]: Al Actualizar\n", err)
	}
	fmt.Println("-- UPDATE --")
	//	Obtenemos el total de elementos modificados
	//	updateres de tipo *UpdateResult
	//	https://godoc.org/go.mongodb.org/mongo-driver/mongo#UpdateResult
	fmt.Println(updateres.ModifiedCount)

	//	Eliminar documento de coleccion
	// delete document
	delres, err := collection.DeleteOne(ctx, bson.M{"name": "dos"})
	if err != nil {
		log.Fatal("[ERROR]: Al eliminar\n", err)
	}
	fmt.Printf("Total Eliminados: %d\n", delres.DeletedCount)
}