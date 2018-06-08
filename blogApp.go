package main

// Import all necessary packages
import (
	"fmt"
	"html/template"
	"io"
	"log"
	"net/http"
	"path/filepath"
	"time"

	"github.com/Masterminds/sprig"
	"github.com/tanmaybaranwal/BlogAppUsingMGO/datastructure"
	"github.com/tanmaybaranwal/BlogAppUsingMGO/dbutilities"
	"gopkg.in/mgo.v2"
	"gopkg.in/mgo.v2/bson"
)

// STATIC_URL : "Path" to the static files (e.g. CSS files, JS files, Images)
const STATIC_URL string = "/static/"

// STATIC_ROOT : "Root" to the static files (e.g. CSS files, JS files, Images)
const STATIC_ROOT string = "static/"

var c *mgo.Collection = nil

// Context struct based on the "post" struct from the datastructure package
type Context struct {
	Posts    []datastructure.Post
	Comments []datastructure.Comment
	Static   string
}

/*
 *  Function for the Home page
 *  The home page should show a list with all tasks in the database
 */
func Home(w http.ResponseWriter, req *http.Request) {
	// Load all tasks
	posts := dbutilities.FindAll(c, []string{"-created_datetime"})

	// Creates an "object" of the Context structure
	var context Context

	// For each post append its values in the "Blogs" vector from the Context structure
	for _, x := range posts {
		content := datastructure.Post{
			ID:              x.ID,
			Author:          x.Author,
			Title:           x.Title,
			Content:         x.Content[0:100] + "...", // Shows a summary of the post (only the first 10 characters)
			IsVerified:      x.IsVerified,
			LastUpdated:     x.LastUpdated,
			CreatedDateTime: x.CreatedDateTime}
		// Append the current post
		context.Posts = append(context.Posts, content)
	}

	// Calls the function to render the home page (index)
	render(w, "index", context)
}

/*
 *  Function for the View page
 *  The View page should show a single post (title, content and author)
 */
func View(w http.ResponseWriter, req *http.Request) {
	// Gets the post title passed by the GET method
	id := req.URL.Path[len("/view/"):]
	fmt.Println(id)
	// If the request method is GET and the title is not empty
	if req.Method == "GET" && id != "" {
		// Search in the database based on the title of the post
		post := dbutilities.Find(c, bson.M{"_id": id})

		// Creates an "object" of the Context structure
		var context Context

		// Get the values from the post variable
		content := datastructure.Post{
			ID:              post.ID,
			Author:          post.Author,
			Title:           post.Title,
			Content:         post.Content,
			IsVerified:      post.IsVerified,
			LastUpdated:     post.LastUpdated,
			CreatedDateTime: post.CreatedDateTime}

		// Append the post content
		context.Posts = append(context.Posts, content)

		// Calls the function to render the View page
		render(w, "view", context)

		return
	}
	// Else, redirect to the home page (index)
	http.Redirect(w, req, "/", http.StatusFound)
}

/*
 *  Function for the Add post
 *  The add page should provide a way to the user enter with the post data (title, content and author)
 */
func Add(w http.ResponseWriter, req *http.Request) {
	// If the method is GET
	if req.Method == "GET" {
		// Creates an "object" of the Context structure
		var context Context
		// Creates an empty "object"
		content := datastructure.Post{Title: ""}
		context.Posts = append(context.Posts, content)
		// Calls the function to render the Edit page
		render(w, "add", context)
	} else {
		// If the method is not GET
		// Creates an object with the values passed by the user over the http.Request
		post := datastructure.Post{
			Author:          req.FormValue("author"),
			Title:           req.FormValue("title"),
			Content:         req.FormValue("content"),
			IsVerified:      false,
			LastUpdated:     time.Now().UTC(),
			CreatedDateTime: time.Now().UTC()}

		// If the title, content and author are filled, insert the post in the database
		if post.Author != "" && post.Title != "" && post.Content != "" {
			// Insert the post "object" in the collection
			dbutilities.Insert(c, post)
			// Redirect to the View page, passing by the GET method the title of the inserted post
			http.Redirect(w, req, "/view/"+req.FormValue("title"), http.StatusFound)
		}
	}
}

func insertDummy() {
	post := datastructure.Post{
		Author:          "tanmayb",
		Title:           "So a new post with ID",
		Content:         "This post contains ID which can we query on. Once it's done, I will try to print it's ID in console using the stub I got from stack overflow. Hail SO!",
		IsVerified:      false,
		LastUpdated:     time.Now().UTC(),
		CreatedDateTime: time.Now().UTC()}

	var id, success = dbutilities.Insert(c, post)
	if !success {
		log.Println("Couldn't save the record!")
	} else {
		fmt.Println("Post Saved with ID: " + id)
	}

	// posts := dbutilities.FindQuery(c, bson.M{"author": "tanmayb"}, []string{"-created_datetime"})

}

/*
 *  Function used to render the pages
 */
func render(w http.ResponseWriter, tmpl string, context Context) {
	// Fill the static variable in the context struct (passed by parameter)
	context.Static = STATIC_URL
	// Creates a template list, based on the base template and the template passed by parameter
	tmplList := []string{"templates/base.html",
		fmt.Sprintf("templates/%s.html", tmpl)}
	name := filepath.Base("templates/base.html")
	// t, err := template.ParseFiles(tmpl_list...)
	// If any error occurs, show it
	// if err != nil {
	// 	log.Println("Template parsing error: ", err)
	// }
	// Applies a parsed template to the specified data object

	t := template.Must(template.New(name).Funcs(sprig.FuncMap()).ParseFiles(tmplList...))
	var err = t.Execute(w, context)
	// If any error occurs, show it
	if err != nil {
		log.Println("Template executing error: ", err)
	}
}

/*
 *  Function used to deal with the static files
 */
func StaticHandler(w http.ResponseWriter, req *http.Request) {
	static_file := req.URL.Path[len(STATIC_URL):]
	// If the path is not empty
	if len(static_file) != 0 {
		f, err := http.Dir(STATIC_ROOT).Open(static_file)
		if err == nil {
			content := io.ReadSeeker(f)
			http.ServeContent(w, req, static_file, time.Now(), content)
			return
		}
	}
	http.NotFound(w, req)
}

/*
 *  Main function: connects to the database and opens a session to work in the "tasks" collection.
 */
func main() {
	// Establish a connection, obtain a session
	session := dbutilities.Connect("localhost")
	// Ensure that the session will be closed
	defer session.Close()

	c = session.DB("blog").C("posts")
	// insertDummy()
	// Assigns each page to each function
	http.HandleFunc("/", Home)
	http.HandleFunc("/add/", Add)
	http.HandleFunc("/view/", View)
	http.HandleFunc(STATIC_URL, StaticHandler)

	err := http.ListenAndServe(":8080", nil)
	if err != nil {
		log.Fatal("ListenAndServe: ", err)
	}
}
