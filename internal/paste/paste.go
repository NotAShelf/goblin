package paste

import (
	"fmt"
	"html/template"
	"io"
	"math/rand"
	"net/http"
	"os"
	"path/filepath"
	"time"

	"goblin/internal/metrics"
	"goblin/internal/util"

	"github.com/gorilla/mux"
	log "github.com/sirupsen/logrus"
	"github.com/spf13/viper"
)

var logger = log.WithFields(log.Fields{
	"package": "paste",
})

const idLength = 8

var letters = []rune("abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ")

var (
	pasteMap    = make(map[string]*Paste)
	templateDir string
	pasteDir    string
	pasteTmpl   *template.Template // Declare pasteTmpl at the package level
)

func generateUniqueID() string {
	idRunes := make([]rune, idLength)
	for i := range idRunes {
		idRunes[i] = letters[rand.Intn(len(letters))]
	}

	uniqueID := fmt.Sprintf("%s-%d", string(idRunes), time.Now().UnixNano())

	return uniqueID
}

type Paste struct {
	ID        string
	Content   string
	CreatedAt time.Time
	ExpiresAt time.Time
}

func InitializeTemplates() error {
	templateDir = viper.GetString("TemplateDir")
	pasteTmpl = template.Must(template.ParseFiles(filepath.Join(templateDir, "paste.html")))
	return nil
}

func CreatePasteHandler(w http.ResponseWriter, r *http.Request) {
	// Parse the request body
	content, err := io.ReadAll(r.Body)
	if err != nil {
		http.Error(w, "Failed to read request body", http.StatusBadRequest)
		metrics.IncrementPasteCreatedCounter("failure")
		return
	}

	// Check if the private flag is set
	private := viper.GetBool("Private")

	// Log the event
	log.Infof("Received request to create a paste (Private mode: %v)", private)

	// Log the content if private mode is not enabled
	if !private {
		log.Infof("Received content: %s", content)
	}

	// Generate a unique ID for the paste
	pasteID := generateUniqueID()

	// Save the paste in the map
	paste := &Paste{
		ID:        pasteID,
		Content:   string(content),
		CreatedAt: time.Now(),
		ExpiresAt: time.Now().Add(24 * time.Hour),
	}

	pasteMap[pasteID] = paste

	// Get the directory path from configuration
	pasteDir := viper.GetString("PasteDir")

	// Create the paste directory if it doesn't exist
	if err := util.CreateDirectoryIfNotExists(pasteDir); err != nil {
		logger.Errorf("Failed to create paste directory: %v", err)
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		metrics.IncrementPasteCreatedCounter("failure")
		return
	}

	// Save the paste content to a file
	pasteFilePath := filepath.Join(pasteDir, pasteID+".txt")
	pasteFile, err := os.Create(pasteFilePath)
	if err != nil {
		logger.Errorf("Failed to create paste file: %v", err)
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		metrics.IncrementPasteCreatedCounter("failure")
		return
	}
	defer pasteFile.Close()

	// Write the paste content to the file
	_, err = pasteFile.WriteString(string(content))
	if err != nil {
		logger.Errorf("Failed to write paste content to file: %v", err)
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		metrics.IncrementPasteCreatedCounter("failure")
		return
	}

	// Construct the complete paste URL including the protocol
	var protocol string
	if r.TLS != nil {
		protocol = "https"
	} else {
		protocol = "http"
	}
	pasteURL := fmt.Sprintf("%s://%s/%s", protocol, r.Host, pasteID)

	// Write the paste URL in the response
	response := fmt.Sprintf("Paste available at %s", pasteURL)
	w.Write([]byte(response))

	// Record metrics
	metrics.IncrementPasteCreatedCounter("success")
	metrics.ObservePasteLength("success", float64(len(content)))
}

func GetPasteHandler(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	pasteID := vars["id"]

	// Exclude the "/metrics" path
	if pasteID == "metrics" {
		http.NotFound(w, r)
		return
	}

	p, exists := pasteMap[pasteID]
	if !exists {
		http.Error(w, "Paste not found", http.StatusNotFound)
		return
	}

	if time.Now().After(p.ExpiresAt) {
		delete(pasteMap, pasteID)
		http.Error(w, "Paste has expired", http.StatusGone)
		return
	}

	log.Infof("Template file path: %s", templateDir)

	err := pasteTmpl.Execute(w, p)
	if err != nil {
		log.Errorf("Error executing template: %v", err)
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
	}
}

func ExpirePastes() {
	// Read the expiration duration from configuration
	expirationDuration := viper.GetDuration("ExpirationDuration")

	for {
		time.Sleep(expirationDuration)

		for pasteID, p := range pasteMap {
			if time.Now().After(p.ExpiresAt) {
				delete(pasteMap, pasteID)
			}
		}
	}
}
