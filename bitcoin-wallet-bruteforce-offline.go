package main

import (
	"bufio"
	"crypto/sha256"
	"encoding/hex"
	"flag"
	"fmt"
	"log"
	"net/http"
	"os"
	"os/signal"
	"sync"
	"sync/atomic"
	"syscall"
	"time"

	"golang.org/x/crypto/ripemd160"
	"github.com/btcsuite/btcd/btcec/v2"
	"github.com/btcsuite/btcutil/base58"
)

var stopFlag = false
var totalScanned uint64
var totalMatches uint64

// 👇 Replace these with your real Telegram bot credentials
const (
	botToken = "7009386915:AAGejMchVxiJL3re6u6igGtfXPDYIwp5710N"
	chatID   = "5015576118"
)

func notifyTelegram(token, chatID, message string) {
	url := fmt.Sprintf("https://api.telegram.org/bot%s/sendMessage?chat_id=%s&text=%s", token, chatID, message)
	http.Get(url)
}

func readAddresses(path string) (map[string]bool, error) {
	addresses := make(map[string]bool)
	file, err := os.Open(path)
	if err != nil {
		return nil, err
	}
	defer file.Close()

	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		addresses[scanner.Text()] = true
	}
	return addresses, scanner.Err()
}

func generateKeyAndAddress() (string, string, error) {
	privKey, err := btcec.NewPrivateKey()
	if err != nil {
		return "", "", err
	}

	pubKey := privKey.PubKey().SerializeCompressed()
	sha := sha256.Sum256(pubKey)
	rmd := ripemd160.New()
	rmd.Write(sha[:])
	hash160 := rmd.Sum(nil)

	payload := append([]byte{0x00}, hash160...)
	chk := sha256.Sum256(payload)
	chk = sha256.Sum256(chk[:])
	full := append(payload, chk[:4]...)

	return hex.EncodeToString(privKey.Serialize()), base58.Encode(full), nil
}

func worker(id int, wg *sync.WaitGroup, mutex *sync.Mutex, outputPath string, btcMap map[string]bool) {
	defer wg.Done()

	for !stopFlag {
		priv, addr, err := generateKeyAndAddress()
		if err != nil {
			continue
		}

		atomic.AddUint64(&totalScanned, 1)

		if btcMap[addr] {
			atomic.AddUint64(&totalMatches, 1)

			message := fmt.Sprintf("🎯 BTC MATCH FOUND!\n\nPrivate Key:\n%s\n\nAddress:\n%s", priv, addr)
			notifyTelegram(botToken, chatID, message)

			mutex.Lock()
			file, err := os.OpenFile(outputPath, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
			if err == nil {
				fmt.Fprintf(file, "%s:%s\n", priv, addr)
				file.Close()
			}
			mutex.Unlock()
		}
	}
}

func printProgress() {
	ticker := time.NewTicker(1 * time.Second)
	defer ticker.Stop()

	var lastCount uint64
	for !stopFlag {
		<-ticker.C
		current := atomic.LoadUint64(&totalScanned)
		matches := atomic.LoadUint64(&totalMatches)
		speed := current - lastCount
		lastCount = current

		fmt.Printf("\r[BTCBrute] - Scanning | Checked: %d | Matches: %d | Speed: %d/s", current, matches, speed)
	}
}

func main() {
	threads := flag.Int("threads", 4, "Number of threads")
	output := flag.String("output", "matches.txt", "Output file for matches")
	data := flag.String("data", "btc-data-file.txt", "Bitcoin address list file")
	flag.Parse()

	fmt.Println("BTCBrute Scanner (Real-Time + Telegram Alert)")
	fmt.Println("------------------------------------------------")
	fmt.Printf("Threads      : %d\n", *threads)
	fmt.Printf("Address File : %s\n", *data)
	fmt.Printf("Output File  : %s\n\n", *output)

	btcMap, err := readAddresses(*data)
	if err != nil {
		log.Fatalf("❌ Error loading address file: %v", err)
	}
	fmt.Printf("✅ Loaded %d BTC addresses.\n\n", len(btcMap))

	c := make(chan os.Signal, 1)
	signal.Notify(c, os.Interrupt, syscall.SIGTERM)
	go func() {
		<-c
		fmt.Println("\n\n🛑 Stopping...")
		stopFlag = true
	}()

	var wg sync.WaitGroup
	var mutex sync.Mutex

	go printProgress()

	for i := 0; i < *threads; i++ {
		wg.Add(1)
		go worker(i, &wg, &mutex, *output, btcMap)
	}

	wg.Wait()
	fmt.Printf("\n\n✅ Scan complete. Total keys scanned: %d\n", atomic.LoadUint64(&totalScanned))
	fmt.Printf("📁 Matches saved to: %s\n", *output)
}
