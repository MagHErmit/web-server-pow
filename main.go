package main

import (
	"bufio"
	"crypto/sha256"
	"encoding/hex"
	"fmt"
	"github.com/joho/godotenv"
	"log"
	"math/big"
	"math/rand"
	"net"
	"os"
	"strconv"
	"strings"
	"time"
)

var difficulty int

func main() {
	// load env file
	err := godotenv.Load()

	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}

	diff := os.Getenv("DIFFICULTY")
	if diff == "" {
		difficulty = 10
	} else {
		difficulty, err = strconv.Atoi(diff)
		if err != nil {
			log.Printf("Error parsing difficulty to int: %v", err)
		}
	}

	quotesFile := os.Getenv("QUOTES_PATH")
	if quotesFile == "" {
		quotesFile = "quotes/quotes.txt"
	}

	log.Printf("Starting server on port %s with difficulty %d", port, difficulty)
	log.Printf("Quotes file: %s", quotesFile)

	quotes := loadQuotes(quotesFile)
	ln, err := net.Listen("tcp", ":"+port)
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
	defer ln.Close()
	log.Printf("Server started on :%s", port)
	for {
		conn, err := ln.Accept()
		if err != nil {
			fmt.Println(err)
			continue
		}
		go handleConnection(conn, quotes)
	}
}

func handleConnection(conn net.Conn, quotes []string) {
	defer conn.Close()
	challenge := generateChallenge(difficulty)
	_, err := fmt.Fprintf(conn, "%s\n", challenge)
	if err != nil {
		log.Printf("Error writing to connection: %v", err)
		return
	}
	scanner := bufio.NewScanner(conn)
	for scanner.Scan() {
		nonce := scanner.Text()
		if validateProofOfWork(challenge, nonce, difficulty) {
			quote := getRandomQuote(quotes)
			_, err := fmt.Fprintf(conn, "Quote: %s\n", quote)
			if err != nil {
				log.Printf("Error writing to connection: %v", err)
			}
			return
		}
		_, err := fmt.Fprintf(conn, "Invalid solution. Try again.\nChallenge: %s\n", challenge)
		if err != nil {
			log.Printf("Error writing to connection: %v", err)
			return
		}
	}
}

func generateChallenge(difficulty int) string {
	hash := sha256.Sum256([]byte(time.Now().String()))
	for i := 0; i < difficulty; i++ {
		hash = sha256.Sum256(hash[:])
	}
	return hex.EncodeToString(hash[:])
}

func validateProofOfWork(challenge, nonce string, difficulty int) bool {
	hash := sha256.Sum256([]byte(fmt.Sprintf("%s%s", challenge, nonce)))
	target := big.NewInt(1)
	target.Lsh(target, uint(256-difficulty))
	hashInt := new(big.Int).SetBytes(hash[:])
	return hashInt.Cmp(target) == -1
}

func getRandomQuote(quotes []string) string {
	source := rand.New(rand.NewSource(time.Now().UnixNano()))
	return quotes[source.Intn(len(quotes))]
}

func loadQuotes(file string) []string {
	f, err := os.Open(file)
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
	defer f.Close()
	scanner := bufio.NewScanner(f)
	var quotes []string
	for scanner.Scan() {
		quote := strings.TrimSpace(scanner.Text())
		if quote != "" {
			quotes = append(quotes, quote)
		}
	}
	return quotes
}
