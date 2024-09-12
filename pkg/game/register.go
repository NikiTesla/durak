package durak

import (
	"encoding/json"
	"errors"
	"net/http"

	"golang.org/x/crypto/bcrypt"
)

// username to hashes password
// TODO move to storage layer
var usersCreds = map[string]string{}

func (g *Game) registerHandler(w http.ResponseWriter, r *http.Request) {
	var creds struct {
		Username string `json:"username"`
		Password string `json:"password"`
	}
	json.NewDecoder(r.Body).Decode(&creds)

	if err := g.validateRegisterCreds(creds.Username, creds.Password); err != nil {
		if errors.Is(err, ErrUsernameIsEmpty) || errors.Is(err, ErrPasswordIsEmpty) {
			http.Error(w, "please provide both username and password", http.StatusBadRequest)
			return
		}

		g.logger.WithError(err).Error("failed to validate user's registration credentials")
		http.Error(w, "credentials validation failed", http.StatusInternalServerError)
		return
	}

	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(creds.Password), bcrypt.DefaultCost)
	if err != nil {
		http.Error(w, "failed to hash password", http.StatusInternalServerError)
		return
	}

	usersCreds[creds.Username] = string(hashedPassword)

	if getRole(creds.Username) == playerRole {
		if err = g.createPlayer(creds.Username); err != nil {
			g.logger.WithError(err).Errorf("failed to create player with username %s", creds.Username)
			http.Error(w, "failed to create player", http.StatusInternalServerError)
			return
		}
		g.logger.Debugf("Player with %s was successfully created", creds.Username)
	}

	token, err := generateJWT(creds.Username)
	if err != nil {
		http.Error(w, "Failed to generate token", http.StatusInternalServerError)
		return
	}

	json.NewEncoder(w).Encode(map[string]string{"token": token})
}

func (g *Game) validateRegisterCreds(username, password string) error {
	if username == "" {
		return ErrPasswordIsEmpty
	}

	if password == "" {
		return ErrUsernameIsEmpty
	}

	return nil
}

func (g *Game) loginHandler(w http.ResponseWriter, r *http.Request) {
	var creds struct {
		Username string `json:"username"`
		Password string `json:"password"`
	}
	json.NewDecoder(r.Body).Decode(&creds)

	storedPassword, ok := usersCreds[creds.Username]
	if !ok || bcrypt.CompareHashAndPassword([]byte(storedPassword), []byte(creds.Password)) != nil {
		http.Error(w, "Invalid credentials", http.StatusUnauthorized)
		return
	}

	token, err := generateJWT(creds.Username)
	if err != nil {
		http.Error(w, "Failed to generate token", http.StatusInternalServerError)
		return
	}

	json.NewEncoder(w).Encode(map[string]string{"token": token})
}

func getRole(username string) string {
	if username == "admin" {
		return adminRole
	}
	return playerRole
}
