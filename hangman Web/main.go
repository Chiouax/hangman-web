package main

import (
	"encoding/json"
	"fmt"
	"html/template"
	"log"
	"math/rand"
	"net/http"
	"os"
	"strings"
	"sync"
	"time"
)

type GameState struct {
	Word           string
	HiddenWord     string
	TriedLetters   []string
	RemainingLives int
	IsGameOver     bool
	Message        string
	PlayerName     string
	Difficulty     string
	IsVictory      bool
	IsCorrect      bool
	mu             sync.Mutex
}

type Score struct {
	PlayerName string
	Word       string
	Difficulty string
	IsVictory  bool
	Date       string
}

var (
	gameStates = make(map[string]*GameState)
	templates  = template.Must(template.ParseGlob("templates/*.html"))
	scoresFile = "scores.json"
)

func main() {
	// Initialisation du générateur de nombres aléatoires
	rand.Seed(time.Now().UnixNano())

	// Routes statiques avec logging
	fs := http.FileServer(http.Dir("static"))
	http.Handle("/static/", http.StripPrefix("/static/", logStaticFiles(fs)))

	// Routes des pages
	http.HandleFunc("/", handleHome)
	http.HandleFunc("/game", handleGame)
	http.HandleFunc("/end", handleEnd)
	http.HandleFunc("/scores", handleScores)

	// Routes d'action
	http.HandleFunc("/start", handleStartGame)
	http.HandleFunc("/play", handlePlayTurn)

	log.Println("Serveur démarré sur http://localhost:8080")
	log.Fatal(http.ListenAndServe(":8080", nil))
}

func logStaticFiles(h http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		log.Printf("Requête fichier statique: %s", r.URL.Path)
		h.ServeHTTP(w, r)
	})
}

func handleHome(w http.ResponseWriter, r *http.Request) {
	session, _ := r.Cookie("session")
	if session != nil {
		if game, exists := gameStates[session.Value]; exists && !game.IsGameOver {
			http.Redirect(w, r, "/game", http.StatusSeeOther)
			return
		}
	}

	templates.ExecuteTemplate(w, "home.html", nil)
}

func handleStartGame(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "Méthode non autorisée", http.StatusMethodNotAllowed)
		return
	}

	playerName := r.FormValue("playerName")
	difficulty := r.FormValue("difficulty")

	if playerName == "" || difficulty == "" {
		http.Error(w, "Paramètres manquants", http.StatusBadRequest)
		return
	}

	word := selectWord(difficulty)
	hiddenWord := strings.Repeat("_", len(word))

	game := &GameState{
		PlayerName:     playerName,
		Difficulty:     difficulty,
		RemainingLives: 6, // Commence à 6 pour les images 0-6
		Word:           word,
		HiddenWord:     hiddenWord,
		TriedLetters:   make([]string, 0),
	}

	sessionID := generateSessionID()
	gameStates[sessionID] = game

	log.Printf("Nouvelle partie créée - Joueur: %s, Mot: %s, Vies: %d",
		playerName, word, game.RemainingLives)

	http.SetCookie(w, &http.Cookie{
		Name:     "session",
		Value:    sessionID,
		Path:     "/",
		HttpOnly: true,
	})

	http.Redirect(w, r, "/game", http.StatusSeeOther)
}

func handleGame(w http.ResponseWriter, r *http.Request) {
	session, err := r.Cookie("session")
	if err != nil {
		log.Printf("Erreur cookie session: %v", err)
		http.Redirect(w, r, "/", http.StatusSeeOther)
		return
	}

	game, exists := gameStates[session.Value]
	if !exists || game.IsGameOver {
		log.Printf("Partie non trouvée ou terminée. Exists: %v, IsGameOver: %v",
			exists, game != nil && game.IsGameOver)
		http.Redirect(w, r, "/", http.StatusSeeOther)
		return
	}

	log.Printf("Affichage jeu - Joueur: %s, Vies: %d, Mot caché: %s",
		game.PlayerName, game.RemainingLives, game.HiddenWord)
	log.Printf("Image chargée: hangman-%d.png", game.RemainingLives)

	templates.ExecuteTemplate(w, "game.html", game)
}

func handlePlayTurn(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "Méthode non autorisée", http.StatusMethodNotAllowed)
		return
	}

	session, err := r.Cookie("session")
	if err != nil {
		http.Error(w, "Session non trouvée", http.StatusBadRequest)
		return
	}

	game, exists := gameStates[session.Value]
	if !exists {
		http.Error(w, "Partie non trouvée", http.StatusBadRequest)
		return
	}

	guess := r.FormValue("guess")
	game.mu.Lock()
	defer game.mu.Unlock()

	guess = strings.ToUpper(guess)

	if len(guess) == 1 {
		letter := guess
		// Vérifie si la lettre a déjà été essayée
		for _, tried := range game.TriedLetters {
			if tried == letter {
				game.Message = "Vous avez déjà essayé cette lettre !"
				game.IsCorrect = false
				http.Redirect(w, r, "/game", http.StatusSeeOther)
				return
			}
		}

		game.TriedLetters = append(game.TriedLetters, letter)
		found := false
		wordRunes := []rune(game.Word)
		hiddenRunes := []rune(game.HiddenWord)

		for i, char := range wordRunes {
			if string(char) == letter {
				hiddenRunes[i] = char
				found = true
			}
		}

		game.HiddenWord = string(hiddenRunes)

		if found {
			game.Message = "Bonne lettre !"
			game.IsCorrect = true
		} else {
			game.RemainingLives--
			game.Message = "Lettre incorrecte !"
			game.IsCorrect = false
			log.Printf("Vie perdue ! Restant: %d", game.RemainingLives)
		}
	} else {
		if strings.EqualFold(guess, game.Word) {
			game.HiddenWord = game.Word
			game.IsGameOver = true
			game.IsVictory = true
			game.Message = "Félicitations ! Vous avez trouvé le mot !"
			game.IsCorrect = true
		} else {
			game.RemainingLives--
			game.Message = "Ce n'est pas le bon mot !"
			game.IsCorrect = false
			log.Printf("Vie perdue ! Restant: %d", game.RemainingLives)
		}
	}

	if !strings.Contains(game.HiddenWord, "_") {
		game.IsGameOver = true
		game.IsVictory = true
		game.Message = "Félicitations ! Vous avez trouvé le mot !"
		game.IsCorrect = true
		saveScore(game)
		http.Redirect(w, r, "/end", http.StatusSeeOther)
		return
	}

	if game.RemainingLives <= 0 {
		game.IsGameOver = true
		game.IsVictory = false
		game.Message = "Perdu ! Le mot était : " + game.Word
		game.IsCorrect = false
		saveScore(game)
		http.Redirect(w, r, "/end", http.StatusSeeOther)
		return
	}

	http.Redirect(w, r, "/game", http.StatusSeeOther)
}

func handleEnd(w http.ResponseWriter, r *http.Request) {
	session, err := r.Cookie("session")
	if err != nil {
		http.Redirect(w, r, "/", http.StatusSeeOther)
		return
	}

	game, exists := gameStates[session.Value]
	if !exists || !game.IsGameOver {
		http.Redirect(w, r, "/", http.StatusSeeOther)
		return
	}

	templates.ExecuteTemplate(w, "end.html", game)
}

func handleScores(w http.ResponseWriter, r *http.Request) {
	scores := []Score{}
	file, err := os.ReadFile(scoresFile)
	if err != nil && !os.IsNotExist(err) {
		log.Printf("Erreur lecture scores: %v", err)
		http.Error(w, "Erreur lors de la lecture des scores", http.StatusInternalServerError)
		return
	}

	if len(file) > 0 {
		if err := json.Unmarshal(file, &scores); err != nil {
			log.Printf("Erreur parsing scores: %v", err)
			http.Error(w, "Erreur lors de la lecture des scores", http.StatusInternalServerError)
			return
		}
	}

	templates.ExecuteTemplate(w, "scores.html", struct{ Scores []Score }{scores})
}

func saveScore(game *GameState) error {
	scores := []Score{}

	file, err := os.ReadFile(scoresFile)
	if err != nil && !os.IsNotExist(err) {
		return err
	}

	if len(file) > 0 {
		if err := json.Unmarshal(file, &scores); err != nil {
			return err
		}
	}

	scores = append(scores, Score{
		PlayerName: game.PlayerName,
		Word:       game.Word,
		Difficulty: game.Difficulty,
		IsVictory:  game.IsVictory,
		Date:       time.Now().Format("2006-01-02 15:04:05"),
	})

	data, err := json.Marshal(scores)
	if err != nil {
		return err
	}

	return os.WriteFile(scoresFile, data, 0644)
}

func selectWord(difficulty string) string {
	words := map[string][]string{
		"easy": {
			"SECOUPE", "SPATIAL", "LASER", "FUTURE", "TELEPHONE",
			"TABLE", "FLEUR", "PORTE", "ROUTE", "POMME",
		},
		"medium": {
			"VOITURE", "MONTAGNE", "POISSON", "ORDINATEUR", "TELEPHONE",
			"MUSIQUE", "FENETRE", "JOURNAL", "CINEMA", "RESTAURANT",
		},
		"hard": {
			"DEVELOPPEMENT", "ARCHITECTURE", "TECHNOLOGIE", "INFORMATIQUE", "PHILOSOPHIE",
			"EXTRAORDINAIRE", "EXPLORATION", "IMAGINATION", "INTELLIGENCE", "COMPLEXITE",
		},
	}

	difficultyWords, exists := words[difficulty]
	if !exists {
		difficultyWords = words["medium"]
	}

	return difficultyWords[rand.Intn(len(difficultyWords))]
}

func generateSessionID() string {
	b := make([]byte, 16)
	rand.Read(b)
	return fmt.Sprintf("%x", b)
}
