package main

import (
	"database/sql"
	"fmt"
	"strconv"
	"log"
	"os"
	"time"

	"github.com/lib/pq"
)

var POSTGRES_HOST string = os.Getenv("POSTGRES_HOST")
var POSTGRES_USER string = os.Getenv("POSTGRES_USER")
var POSTGRES_PASSWORD string= os.Getenv("POSTGRES_PASSWORD")
var POSTGRES_DB_NAME string = os.Getenv("POSTGRES_DB_NAME")
var S3_RESOURCES_BUCKET string = os.Getenv("RESOURCES_BUCKET_NAME")

func getEnv(key, fallback string) string {
    if value, ok := os.LookupEnv(key); ok {
        return value
    }
    return fallback
}

func main() {
	postgresPort := getEnv("POSTGRES_PORT", "5432")
	port, err := strconv.Atoi(postgresPort)
    if err != nil {
        panic(err)
    }
	dbinfo := fmt.Sprintf("host=%s port=%d user=%s password=%s dbname=%s sslmode=disable",
    POSTGRES_HOST, port, POSTGRES_USER, POSTGRES_PASSWORD, POSTGRES_DB_NAME)

	db, err := sql.Open("postgres", dbinfo)
	if err != nil {
		log.Fatal(err)
	}
	defer db.Close()

	// 1. Check for feeds that violated the SLA
	violatedQuery := `
	SELECT feed_id 
	FROM feed_sla 
	WHERE (schedule='daily' AND last_load_date < $1) OR 
		  (schedule='weekly' AND last_load_date < $2) OR 
		  (schedule='monthly' AND last_load_date < $3)
	`
	now := time.Now()
	yesterday := now.AddDate(0, 0, -1)
	lastWeek := now.AddDate(0, 0, -7)
	lastMonth := now.AddDate(0, -1, 0)

	rows, err := db.Query(violatedQuery, yesterday, lastWeek, lastMonth)
	if err != nil {
		log.Fatal(err)
	}
	defer rows.Close()

	var violatedFeedIDs []int
	for rows.Next() {
		var feedID int
		err := rows.Scan(&feedID)
		if err != nil {
			log.Fatal(err)
		}
		violatedFeedIDs = append(violatedFeedIDs, feedID)
	}

	if err = rows.Err(); err != nil {
		log.Fatal(err)
	}

	if len(violatedFeedIDs) > 0 {
		updateQuery := `UPDATE feed_sla SET sla_missed=true WHERE feed_id = ANY($1)`
		_, err := db.Exec(updateQuery, pq.Array(violatedFeedIDs))
		if err != nil {
			log.Fatal(err)
		}
		log.Printf("Updated sla_missed to true for feed IDs: %v", violatedFeedIDs)
	}

	// 2. Check for feeds that adhered to their SLA
	adheredQuery := `
	SELECT feed_id 
	FROM feed_sla 
	WHERE (schedule='daily' AND last_load_date >= $1) OR 
		  (schedule='weekly' AND last_load_date >= $2) OR 
		  (schedule='monthly' AND last_load_date >= $3)
	`
	rows, err = db.Query(adheredQuery, yesterday, lastWeek, lastMonth)
	if err != nil {
		log.Fatal(err)
	}
	defer rows.Close()

	var adheredFeedIDs []int
	for rows.Next() {
		var feedID int
		err := rows.Scan(&feedID)
		if err != nil {
			log.Fatal(err)
		}
		adheredFeedIDs = append(adheredFeedIDs, feedID)
	}

	if len(adheredFeedIDs) > 0 {
		updateQuery := `UPDATE feed_sla SET sla_missed=false WHERE feed_id = ANY($1)`
		_, err := db.Exec(updateQuery, pq.Array(adheredFeedIDs))
		if err != nil {
			log.Fatal(err)
		}
		log.Printf("Updated sla_missed to false for feed IDs: %v", adheredFeedIDs)
	}
}