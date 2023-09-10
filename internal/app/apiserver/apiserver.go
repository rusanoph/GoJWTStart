package apiserver

import (
	"fmt"
	"net/http"
	"time"

	"github.com/dgrijalva/jwt-go"
)

type ApiServer struct {
	config *Config
	secret []byte
}

func NewApiServer(config *Config) *ApiServer {
	return &ApiServer{
		config: config,
		secret: []byte(config.SecretKey),
	}
}

func (as *ApiServer) Start() error {

	http.Handle("/home", as.ValidateJWT(as.handleHome))
	http.HandleFunc("/jwt", as.handleGetJwt)

	return http.ListenAndServe(as.config.BindAddr, nil)
}

func (as *ApiServer) CreateJWT() (string, error) {
	token := jwt.New(jwt.SigningMethodHS256)
	claims := token.Claims.(jwt.MapClaims)

	claims["exp"] = time.Now().Add(time.Hour).Unix()

	tokenStr, err := token.SignedString(as.secret)
	if err != nil {
		fmt.Println(err.Error())
		return "", err
	}

	return tokenStr, nil
}

func (as *ApiServer) ValidateJWT(next func(w http.ResponseWriter, r *http.Request)) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {

		if r.Header["Token"] != nil {
			token, err := jwt.Parse(r.Header["Token"][0], func(t *jwt.Token) (interface{}, error) {
				if _, ok := t.Method.(*jwt.SigningMethodHMAC); !ok {
					w.WriteHeader(http.StatusUnauthorized)
					w.Write([]byte("not authorized wrong token format"))
				}
				return as.secret, nil
			})

			if err != nil {
				w.WriteHeader(http.StatusUnauthorized)
				w.Write([]byte("not authorized " + err.Error()))
			}

			if token.Valid {
				next(w, r)
			}
		} else {
			w.WriteHeader(http.StatusUnauthorized)
			w.Write([]byte("not authorized token is nil"))
		}
	})
}

func (as *ApiServer) handleGetJwt(w http.ResponseWriter, r *http.Request) {
	if r.Header["Access"] != nil {
		if r.Header["Access"][0] == as.config.ApiKey {
			token, err := as.CreateJWT()
			if err != nil {
				return
			}
			fmt.Fprint(w, token)
		}
	} else {
		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte("no Access"))
	}
}

func (as *ApiServer) handleHome(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintf(w, "super secret area")
}
