package main

import (
	"fmt"
	"math/rand"
	"net/http"
	"os"
	"strings"
	"text/template"
	"time"
)

type Hangman struct {
	Mot_cache     string
	Mot_a_trouver string
	Inputletter   string
	Life          int
	Win           bool
	Erreur        string
	Difficulte    string
	StockLetter   []string
}

func Home(rw http.ResponseWriter, r *http.Request, Pts *Hangman) {
	tmp, _ := template.ParseFiles("./index.html", "./template/header.html", "./template/footer.html")
	tmp.Execute(rw, Pts)
}

func Info(w http.ResponseWriter, r *http.Request) {
	tmp, _ := template.ParseFiles("./page/info.html", "./template/header.html", "./template/footer.html")
	tmp.Execute(w, Info)
}

func main() {
	HangPts := Hangman{Mot_cache: "", Mot_a_trouver: "", Inputletter: "", Life: 10}
	Pts := &HangPts
	Pts.Life = 10
	fmt.Println(Pts.Life)
	Pts.Erreur = ""
	Pts.Mot_a_trouver = ""
	Pts.Mot_cache = ""
	Pts.Difficulte = ""
	Pts.Win = false
	Pts.StockLetter = []string{}

	http.HandleFunc("/", func(rw http.ResponseWriter, r *http.Request) {
		Home(rw, r, Pts)
	})

	http.HandleFunc("/info", func(w http.ResponseWriter, r *http.Request) {
		Info(w, r)

	})

	// InitialiseStuct(Pts)

	http.HandleFunc("/hangman", func(w http.ResponseWriter, r *http.Request) {
		difficulty := r.FormValue("PLAY")
		if difficulty == "easy" || difficulty == "medium" || difficulty == "hard" {
			fmt.Println(difficulty)
			Pts.Difficulte = difficulty
			Pts.Mot_a_trouver = Pickword("photos_mots/" + difficulty + ".txt")
			Pts.Mot_cache = strings.Repeat("_", len(Pts.Mot_a_trouver))

		} else {

			Pts.Inputletter = r.FormValue("letter")
			Pts.StockLetter = append(Pts.StockLetter, Pts.Inputletter)
			fmt.Println(Pts.Inputletter)
			if len(Pts.Inputletter) > 0 {
				if !IsNotLetter(Pts) {
					Pts.Erreur = ""
					IfLetterInTheWord(Pts)
				}
				if Equal(Pts) || WordWin(Pts) {
					fmt.Println("yes")
					Pts.Win = true
				}
			}
		}
		// http.Redirect(w, r, "/", http.StatusFound)
		tmp, _ := template.ParseFiles("./page/hangman.html", "./template/header.html", "./template/footer.html")
		tmp.Execute(w, Pts)
	})

	fs := http.FileServer((http.Dir("static/")))
	http.Handle("/static/", http.StripPrefix("/static/", fs))
	fmt.Println("http://localhost:8080")
	http.ListenAndServe(":8080", nil)
}

func IfLetterInTheWord(Pts *Hangman) {
	count := 0
	tabr := []rune(Pts.Mot_cache)
	tabrletter := []rune(Pts.Inputletter)
	for j := 0; j < len(Pts.Mot_a_trouver); j++ {
		if Pts.Inputletter == (string(Pts.Mot_a_trouver[j])) { //vérifi si lettre est dans le mot
			tabr[j] = tabrletter[0] //vérifi si lettre est dans tab_rune
			count++
		}
	}
	if count == 0 {
		Pts.Life--
		fmt.Println(Pts.Life)
	}
	Pts.Mot_cache = string(tabr)
}

func IsNotLetter(Pts *Hangman) bool {

	if (Pts.Inputletter[0] < 97 || Pts.Inputletter[0] > 122) && (Pts.Inputletter[0] < 65 || Pts.Inputletter[0] > 90) {
		Pts.Erreur = "veuillez entrer une lettre"
		return true
	}
	return false
}

func WordWin(Pts *Hangman) bool {

	if Pts.Inputletter == Pts.Mot_a_trouver {
		return true
	}
	return false
}
func IsWord(word string, wordentry string, tab_found []rune, life int) (bool, []rune, int) { // verifie si le mot en entrée est bon
	if word == wordentry { //vérifi si mot entrée par utilisateur = mot à trouvé
		tab_found = []rune(word)
		return true, tab_found, life
	}
	return false, tab_found, life - 2
}

func Equal(Pts *Hangman) bool {

	if Pts.Mot_cache == Pts.Mot_a_trouver {
		return true
	}
	return false
}

func ReadFileName(name string) string {
	file, err := os.Open(name)
	if err != nil {
		fmt.Println("ERROR: open " + name + ": no such file or directory\n")
		os.Exit(1)
	}
	res := ""
	arr := make([]byte, 1000)
	n, _ := file.Read(arr)
	for i := 0; i < n; i++ {
		res += string(arr[i])
	}
	return res

}
func ReadWordsOnFiles(wordsn string) []string {
	tab := strings.Split(wordsn, "\n")
	return tab
}

func Pickword(difficulte string) string {
	rand.Seed(time.Now().UnixNano())
	fileName := ReadFileName(difficulte)
	tabword := ReadWordsOnFiles(fileName)
	aleaindexonfiles := tabword[rand.Intn(len(tabword))]
	fmt.Println(aleaindexonfiles)
	fmt.Println(len(tabword))
	return aleaindexonfiles
}
