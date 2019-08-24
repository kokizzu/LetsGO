package graphql

import (
	"encoding/json"
	"log"
	"net/http"
	"testing"

	"github.com/graphql-go/graphql"
)

func TestItunes(t *testing.T) {
	http.HandleFunc("/graphql", handleQuery())
	_ = http.ListenAndServe(":12345", nil)
}

/**
This
*/
func handleQuery() func(w http.ResponseWriter, r *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		result := graphql.Do(graphql.Params{
			Schema:        createSchema(),
			RequestString: r.URL.Query().Get("query"),
		})
		err := json.NewEncoder(w).Encode(result)
		if err != nil {
			log.Println("Error Serializing result")
			panic(err)
		}
	}
}

/**
THis function use factory [graphql.NewSchema] where we pass [SchemaConfig] which basically require
a Object Query type
*/
func createSchema() graphql.Schema {
	schema, err := graphql.NewSchema(graphql.SchemaConfig{
		Query: createQuery(),
	})
	if err != nil {
		log.Println("Error creating Schema")
		panic(err)
	}
	return schema
}

/**
Graphql query object it's created using [graphql.NewObject] where we pass a [ObjectConfig]
this ObjectConfig it's created filling with:

	Name: Name of the Query
	Fields: map [string]Field where each Field is a select field in your query
	Description: string description of what the query is about.

*/
func createQuery() *graphql.Object {
	return graphql.NewObject(graphql.ObjectConfig{
		Name: "Query",
		Fields: graphql.Fields{
			"actors": loadActorsField(),
			"movie":  loadMovieField(),
		},
		Description: "A movie query with information of most famous movies",
	})
}

//######################
//	 GRAPHQL FIELDS    #
//######################

/**
A Field in graphql is through the type [Field] which is expected to be filled with:
	Name:Name of the field
	Type: Object type that define the attributes that it will used in the GraphQL engine
	Args:The filter arg that we expect to receive in the query and we will use in [Resolve] function
	Resolve: Function to be invoked when the query use this field [actors] and is where use the args to filter from all the data
	Description: description of the field

*/
func loadActorsField() *graphql.Field {
	return &graphql.Field{
		Name: "Actor filed",
		Type: graphql.NewList(actorType),
		Args: graphql.FieldConfigArgument{
			"movie": &graphql.ArgumentConfig{
				Type: graphql.NewNonNull(graphql.String),
			},
		},
		Resolve:     actorResolver,
		Description: "Songs field to find a particular song using album",
	}
}

/**
Same functionality than previous LoadActorField function but this one configure for Movies
*/
func loadMovieField() *graphql.Field {
	return &graphql.Field{
		Type: movieType,
		Args: graphql.FieldConfigArgument{
			"id": &graphql.ArgumentConfig{
				Type: graphql.NewNonNull(graphql.String),
			},
		},
		Resolve: movieResolver,
	}
}

//######################
//	GRAPHQL RESOLVERS  #
//######################

/**
Filter function that it use the filter arguments that we receive form the query using [ResolveParams]
we're able to get the attribute [movie] passed in the query to be used as a filter
*/
var actorResolver = func(params graphql.ResolveParams) (interface{}, error) {
	movie := params.Args["movie"].(string)
	var filterActors []Actor
	for _, actor := range actors {
		if actor.Movie == movie {
			filterActors = append(filterActors, actor)
		}
	}
	return filterActors, nil
}

/**
Filter function that it will receive the filter argument from [ResolveParams] of the query and we will find
if that particular filter data is in our current "database" in memory
*/
var movieResolver = func(params graphql.ResolveParams) (interface{}, error) {
	id := params.Args["id"].(string)
	for _, movie := range movies {
		if movie.ID == id {
			return movie, nil
		}
	}
	return nil, nil
}

//#########################
//	GRAPHQL OBJECT TYPES  #
//#########################

/**
Data types from GraphQL that will be used internally for the engine for all our CRUD with the service.
We need to define as before with the [ObjectConfig],and it will be used as [Type] when we define our [Field]

[Name] in a Type is mandatory for GraphQL library
*/
var actorType = graphql.NewObject(graphql.ObjectConfig{
	Name: "Actor",
	Fields: graphql.Fields{
		"id": &graphql.Field{
			Type: graphql.String,
		},
		"movie": &graphql.Field{
			Type: graphql.String,
		},
		"name": &graphql.Field{
			Type: graphql.String,
		},
		"age": &graphql.Field{
			Type: graphql.String,
		},
	},
})

/**
Data types from GraphQL that will be used internally for the engine for all our CRUD with the service.
We need to define as before with the [ObjectConfig],and it will be used as [Type] when we define our [Field]

[Name] in a Type is mandatory for GraphQL library
*/
var movieType = graphql.NewObject(graphql.ObjectConfig{
	Name: "Movie",
	Fields: graphql.Fields{
		"id": &graphql.Field{
			Type: graphql.String,
		},
		"artist": &graphql.Field{
			Type: graphql.String,
		},
		"title": &graphql.Field{
			Type: graphql.String,
		},
		"year": &graphql.Field{
			Type: graphql.String,
		},
		"genre": &graphql.Field{
			Type: graphql.String,
		},
		"type": &graphql.Field{
			Type: graphql.String,
		},
	},
})

//#################
//	  GO TYPES    #
//#################

/**
Go data types that we use to fill the data from our mocks database memory
*/
type Movie struct {
	ID     string `json:"id,omitempty"`
	Artist string `json:"artist"`
	Title  string `json:"title"`
	Year   string `json:"year"`
	Genre  string `json:"genre"`
	Type   string `json:"type"`
}

type Actor struct {
	ID    string `json:"id,omitempty"`
	Movie string `json:"movie"`
	Name  string `json:"name"`
	Age   string `json:"age"`
	Type  string `json:"type"`
}

//##################
//	  MOCK DATA    #
//##################

var movies = []Movie{
	{
		ID:     "ts-fearless",
		Artist: "1",
		Title:  "Fearless",
		Year:   "2008",
		Type:   "album",
	},
}

var actors = []Actor{
	{
		ID:    "1",
		Movie: "Titanic",
		Name:  "Brad Pitt",
		Age:   "60",
		Type:  "song",
	},
	{
		ID:    "2",
		Movie: "ts-fearless",
		Name:  "Keanu Reeves",
		Age:   "54",
		Type:  "song",
	},
}

/*curl -g 'http://localhost:12345/graphql?query={actors(movie:"ts-fearless"){name,age}}'
 */
