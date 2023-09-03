package core

import (
	"identityreconciliation/model"
	"identityreconciliation/repository"
	"identityreconciliation/service"
	"log"
	"net"
	"net/http"
	"os"
	"strings"

	"github.com/gorilla/handlers"
	"github.com/gorilla/mux"
	httpSwagger "github.com/swaggo/http-swagger"
)

type Core struct {
	storage *repository.Repository
	service *service.Service
	log     *log.Logger
}

func NewCoreService(storage *repository.Repository, service *service.Service) *Core {
	return &Core{
		storage: storage,
		service: service,
		log:     log.New(os.Stdout, "", log.Ldate|log.Ltime|log.Lmicroseconds|log.Llongfile),
	}
}

func (c *Core) CallRun() {
	c.run()
	//api.NewService(c.DB)
}

func (c *Core) run() {

	c.log.Println("starting core")
	router := mux.NewRouter()

	// Configure database connection
	if err := c.storage.InitDB("identity.db"); err != nil {
		c.log.Fatalf("failed to initialize database: %v", err)
	}
	defer c.storage.CloseDB()

	// Create tables based on struct definitions
	if err := c.storage.CreateTable(model.Contact{}); err != nil {
		c.log.Fatalf("failed to create table: %v", err)
	}

	corsHandler := handlers.CORS(
		handlers.AllowedOrigins([]string{"http://localhost:8080"}),
		handlers.AllowedMethods([]string{"GET", "POST", "PUT", "DELETE", "OPTIONS"}),
		handlers.AllowedHeaders([]string{"Content-Type", "Authorization"}),
	)
	handler := c.corsMiddleware(router)

	// Register API handlers using gorilla/mux
	//router.HandleFunc("/users/create", c.API.CreateUser).Methods("POST")
	//router.HandleFunc("/users/update", c.API.UpdateUser).Methods("PUT")
	router.HandleFunc("/identify", c.service.Identify).Methods("POST")

	//enable swagger
	c.EnableSwagger(c.getURL(), router)
	// Use the gorilla/mux router
	http.Handle("/", corsHandler(router))

	log.Fatal(http.ListenAndServe(":8080", handler))
}

func (c *Core) getURL() string {
	// No IP address present
	url := "http://localhost:8080"
	if strings.Contains(url, "://:") {
		conn, err := net.Dial("udp", "8.8.8.8:80")
		if err != nil {
			return url
		}
		defer conn.Close()
		localAddr := conn.LocalAddr().(*net.UDPAddr)
		outIp := localAddr.IP.String()
		s := strings.Split(url, "://:")
		url = s[0] + "://" + outIp + ":" + s[1]
	}
	c.log.Println("Swagger URL : " + url + "/swagger/index.html")
	return url
}

func (c *Core) EnableSwagger(url string, router *mux.Router) {

	/* r.PathPrefix("/swagger/").Handler(httpSwagger.Handler(
	httpSwagger.URL(url+"/swagger/doc.json"),
	httpSwagger.DeepLinking(true),
	httpSwagger.DocExpansion("none"),
	httpSwagger.DomID("swagger-ui"))).Methods(http.MethodGet) */

	swaggerURL := "/docs/swagger.json"
	router.PathPrefix("/swagger/").Handler(httpSwagger.Handler(
		httpSwagger.URL(swaggerURL), // The url pointing to API definition
	))
	// router.PathPrefix("/swagger/").Handler(httpSwagger.Handler(
	// 	httpSwagger.URL(swaggerURL), // The url pointing to API definition
	// ))
	router.PathPrefix("/docs/").Handler(http.StripPrefix("/docs/", http.FileServer(http.Dir("./docs"))))
	// Log URLs
	c.log.Println("Swagger UI (API Documentation): http://localhost:8080/swagger/")
	c.log.Println("Swagger JSON Specification: http://localhost:8080/docs" + swaggerURL)
}

// corsMiddleware is a middleware function to set the CORS headers in the response.
func (c *Core) corsMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// Set the Access-Control-Allow-Origin header to allow requests from http://localhost:8080
		//	for key, values := range r.Header {
		//		log.Printf("%s: %v\n", key, values)
		//	}
		w.Header().Set("Access-Control-Allow-Origin", "*")
		// Optionally, you can set other CORS headers, such as Access-Control-Allow-Methods, etc.
		w.Header().Set("Access-Control-Allow-Methods", "GET, POST, PUT, DELETE, OPTIONS")
		w.Header().Set("Access-Control-Allow-Headers", "Content-Type, Authorization")
		// Allow preflight requests (OPTIONS method) by setting appropriate headers for preflight responses
		if r.Method == "OPTIONS" {
			w.WriteHeader(http.StatusOK)
			return
		}
		// Call the next handler in the chain
		next.ServeHTTP(w, r)
	})
}
