package main

import (
	"bufio"
	"crypto/sha256"
	"encoding/hex"
	"errors"
	"fmt"
	"log"
	"math/big"
	"math/rand"
	"net"
	"strconv"
	"strings"
	"time"
)

func main() {
	// connect to server
	conn, err := net.Dial("tcp", "localhost:8080")
	if err != nil {
		log.Fatalln(err)
		return
	}
	defer conn.Close()

	reader := bufio.NewReader(conn)
	// read challenge
	challenge, err := reader.ReadString('\n')
	if err != nil {
		log.Fatalln(err)
		return
	}
	log.Printf("Challenge: %s", challenge)

	// read difficulty
	difficulty, err := reader.ReadString('\n')
	if err != nil {
		log.Fatalln(err)
		return
	}
	log.Printf("Difficulty: %s", difficulty)

	diff, err := strconv.Atoi(strings.TrimSuffix(difficulty, "\n"))
	if err != nil {
		log.Fatalln(err)
		return
	}
	// solve challenge using proof of work
	solution, err := solveProofOfWork(strings.TrimSuffix(challenge, "\n"), diff)
	if err != nil {
		log.Fatalln(err)
		return
	}

	// send solution to server
	_, err = conn.Write([]byte(solution + "\n"))
	if err != nil {
		log.Fatalln(err)
		return
	}

	// read quote
	res, err := bufio.NewReader(conn).ReadString('\n')
	if err != nil {
		log.Fatalln(err)
		return
	}
	log.Printf("Response: %s", res)
}

func solveProofOfWork(challenge string, difficulty int) (string, error) {
	// generate random nonce
	rand.Seed(time.Now().UnixNano())
	nonce := rand.Uint32()

	// start timer
	startTime := time.Now()

	// calculate target
	target := big.NewInt(1)
	target.Lsh(target, uint(256-difficulty))

	// loop until solution found
	for {
		// calculate hash
		hash := sha256.Sum256([]byte(fmt.Sprintf("%s%s", challenge, strconv.FormatUint(uint64(nonce), 10))))

		hashInt := new(big.Int).SetBytes(hash[:])
		if hashInt.Cmp(target) == -1 {
			// solution found
			elapsed := time.Since(startTime)
			log.Printf("Solution found in %v seconds\n", elapsed.Seconds())
			log.Printf("Nonce: %v, Hash: %v\n", nonce, hex.EncodeToString(hash[:]))
			return strconv.FormatUint(uint64(nonce), 10), nil
		}

		// generate nonce
		rand.Seed(time.Now().UTC().UnixNano())

		nonce = rand.Uint32()

		// check if time limit exceeded
		elapsed := time.Since(startTime)
		if elapsed.Seconds() > 60 {
			return "", errors.New("proof of work timed out")
		}
	}
}
