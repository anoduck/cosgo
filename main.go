package main

import (
	"flag"
	"github.com/gorilla/mux"
	"github.com/gorilla/csrf"
	"log"
	"net/http"
	"os"
	"math/rand"
	"time"
)

var (
	mandrillApiUrl string
	mandrillKey    string
	casgoDestination string
	casgoAPIKey string
)

func main() {

// We can set the CASGO_API_KEY environment variable, or it defaults to a new random one!

	if os.Getenv("CASGO_API_KEY") == "" {
		log.Println("Generating Random API Key...")

// The length of the API key can be modified here.
		casgoAPIKey = GenerateAPIKey(20)

// Print new GenerateAPIKey
		log.Println("CASGO_API_KEY:",getKey())
	}else{
		casgoAPIKey = os.Getenv("CASGO_API_KEY")
// Print selected CASGO_API_KEY
		log.Println("CASGO_API_KEY:",getKey())
	}
//
	port := flag.String("port", "8080", "HTTP Port to listen on")
	flag.Parse()

	mandrillApiUrl = "https://mandrillapp.com/api/1.0/"
	mandrillKey = os.Getenv("MANDRILL_KEY")
	if mandrillKey == "" {
		log.Fatal("MANDRILL_KEY is Crucial. Type: export MANDRILL_KEY=123456789")
		os.Exit(1)
	}


	casgoDestination = os.Getenv("CASGO_DESTINATION")
	if casgoDestination == "" {
		log.Fatal("CASGO_DESTINATION is Crucial. Type: export CASGO_DESTINATION=\"your@email.com\"")
		os.Exit(1)
	}

	log.Printf("Starting Server on http://127.0.0.1:%s", *port)
	r := mux.NewRouter()




// Custom 404 redirect to /
	r.NotFoundHandler = http.HandlerFunc(RedirectHomeHandler)


// Should be called BlankPageHandler
	r.HandleFunc("/", HomeHandler)

// This is for behind a reverse proxy.
	r.HandleFunc("/" + casgoAPIKey + "/form", ContactHandler)
	r.HandleFunc("/" + casgoAPIKey + "/form/", ContactHandler)
//	r.HandleFunc("/contact/", ContactHandler)

// Fun
	r.HandleFunc("/{whatever}", LoveHandler)

// Magic URL Generator
	r.HandleFunc("/" + casgoAPIKey + "/send", EmailHandler)

	http.Handle("/", r)

// Switch to file log so we can ctrl+c and launch another instance :)
	log.Println("Switching Logs to debug.log")
	OpenLog()
	log.Println("info: Listening on", *port)
	log.Fatal(http.ListenAndServe(":"+*port, csrf.Protect([]byte("LI80PNK1xcT01jmQBsEyxyrNCrbyyFPjPU8CKnxwmCruxNijgnyb3hXXD3p1RBc0+LIRQUUbTtis6hc6LD4I/A=="), csrf.HttpOnly(true), csrf.Secure(false))(r)))

}


// Key Generator
func init() {
		rand.Seed(time.Now().UnixNano())
}
var runes = []rune("____ABCDEFGHIJKLMNOPQRSTUVWXYZ1234567890123456789012345678901234567890")
func GenerateAPIKey(n int) string {
		b := make([]rune, n)
		for i := range b {
				b[i] = runes[rand.Intn(len(runes))]
		}
return string(b)
}

// Which Key are we using again?
func getKey() string {
return casgoAPIKey
}

// This function opens a log file. "debug.log"
func OpenLog(){
f, err := os.OpenFile("./debug.log", os.O_RDWR | os.O_CREATE | os.O_APPEND, 0660)
if err != nil {
    log.Fatal("error opening file: %v", err)
		os.Exit(1)
}
log.SetOutput(f)
}

// This is the home page it is blank. "This server is broken"