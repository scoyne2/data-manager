package feed

import (
	"database/sql"
	"fmt"
	"strconv"
	"os"
	"sync"
	_ "github.com/lib/pq"
	"github.com/jmoiron/sqlx/types"
)

type PostgressRepository struct {
	db *sql.DB
	sync.Mutex
}

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

func NewPostgressRepository() *PostgressRepository {
	postgresPort := getEnv("POSTGRES_PORT", "5432")
	port, err := strconv.Atoi(postgresPort)
    if err != nil {
        panic(err)
    }
	dbinfo := fmt.Sprintf("host=%s port=%d user=%s password=%s dbname=%s sslmode=disable",
    POSTGRES_HOST, port, POSTGRES_USER, POSTGRES_PASSWORD, POSTGRES_DB_NAME)

	db, err := sql.Open("postgres", dbinfo)
	if err != nil {
		panic(err)
	}

	return &PostgressRepository{
		db: db,
	}
}

func (pr *PostgressRepository) AddFeed(feed Feed) (string, error) {
	sqlStatement := `
	INSERT INTO feeds (vendor, feed_name, feed_method)
	SELECT $1, $2, $3
	WHERE
	NOT EXISTS (
	SELECT * FROM feeds WHERE vendor = $1 AND feed_name = $2 AND feed_method = $3
	);`

	_, err := pr.db.Exec(sqlStatement, feed.Vendor, feed.FeedName, feed.FeedMethod)
	if err != nil {
		return "", err
	}
	return "Feed Added", nil
}

func (pr *PostgressRepository) UpdateFeedStatus(fs FeedStatusUpdate) (string, error) {
	// check if feed exists first
	sqlStatementFd := `
	SELECT id, vendor, feed_name, feed_method
	FROM feeds
	WHERE vendor=$1 AND feed_name=$2;`
	feedRows := pr.db.QueryRow(sqlStatementFd, fs.Vendor, fs.FeedName)
	var fd Feed
	feedError := feedRows.Scan(&fd.ID, &fd.Vendor, &fd.FeedName, &fd.FeedMethod)
	if feedError == sql.ErrNoRows {
		return "", fmt.Errorf("could not complete update, feed %s does not exist for vendor %s", fs.FeedName, fs.Vendor)
	}

	// check if feed status exists
	sqlStatementFdS := `
	SELECT 
		COALESCE(fs.id, 0),
		COALESCE(fs.process_date, '1970-01-01'),
		COALESCE(fs.record_count, 0), 
		COALESCE(fs.error_count, 0),
		COALESCE(fs.feed_status, 'N/A'),
		COALESCE(fs.file_name, 'N/A'),
		f.id AS feed_id
	FROM feeds f
	LEFT JOIN feed_status fs
	ON fs.feed_id = f.id
    WHERE f.vendor=$1 AND f.feed_name=$2
	ORDER BY 2 DESC
	LIMIT 1
	;`
	feedStatusRows := pr.db.QueryRow(sqlStatementFdS, fs.Vendor, fs.FeedName)
	var f FeedStatus
	feedStatusErrors := feedStatusRows.Scan(&f.ID, &f.ProcessDate, &f.RecordCount, &f.ErrorCount, &f.Status, &f.FileName, &f.FeedID)
	if feedStatusErrors == sql.ErrNoRows {
		return "", fmt.Errorf("could not complete update, feed status %s does not exist for vendor %s", fs.FeedName, fs.Vendor)
	}
	if feedStatusErrors != nil && feedStatusErrors != sql.ErrNoRows {
		return "", fmt.Errorf("vendor: %s, feed_name: %s feed status query failed: %s", fs.Vendor, fs.FeedName, feedStatusErrors)
	}

	sqlUpdateAllStatement := `
	UPDATE feed_status
	SET is_current = False
	WHERE feed_id = $1;`
	_, updateAllErr := pr.db.Exec(sqlUpdateAllStatement, f.FeedID)
	if updateAllErr != nil {
		return "", updateAllErr
	}

	if f.FileName == fs.FileName {
		sqlUpdateStatement := `
		UPDATE feed_status
		SET process_date = $3, record_count = $4, error_count = $5, feed_status = $6, emr_application_id = NULLIF($7,''), emr_step_id = NULLIF($8,''), is_current = True
		WHERE feed_id = $1 AND file_name = $2;`
		_, updateErr := pr.db.Exec(sqlUpdateStatement, f.FeedID, fs.FileName, fs.ProcessDate, fs.RecordCount, fs.ErrorCount, fs.Status, fs.EMRApplicationID, fs.EMRStepID)
		if updateErr != nil {
			return "", fmt.Errorf("update to feed_status failed. %s", updateErr)
		}
	} else {
		sqlInsertStatement := `
		INSERT INTO feed_status (process_date, record_count, error_count, feed_status, file_name, feed_id, emr_application_id, emr_step_id, is_current)
		VALUES ($1, $2, $3, $4, $5, $6, NULLIF($7,''), NULLIF($8,''), True);`
		_, insertErr := pr.db.Exec(sqlInsertStatement, fs.ProcessDate, fs.RecordCount, fs.ErrorCount, fs.Status, fs.FileName, fd.ID, fs.EMRApplicationID, fs.EMRStepID)
		if insertErr != nil {
			return "", fmt.Errorf("update to feed_status failed. %s", insertErr)
		}
	}	

	return fmt.Sprintf("Feed status updated for id: %d", f.ID), nil
}

type feedStatusResultsDetailedDTO struct {
	ID         	  int            `db:"id"`
	ProcessDate   string         `db:"process_date"`
	RecordCount   int            `db:"record_count"`
	ErrorCount	  int            `db:"error_count"`
	Status	      string         `db:"feed_status"`
	Vendor     	  string         `db:"vendor"`
	FeedName   	  string         `db:"feed_name"`
	FeedMethod	  string         `db:"feed_method"`
	FileName   	  string         `db:"file_name"`
	PreviousFeeds types.JSONText `db:"previous_feeds"`
}

func (d feedStatusResultsDetailedDTO) ToFeedStatusResultsDetailed() (*FeedStatusResultsDetailed, error) {
	result := new(FeedStatusResultsDetailed)
	result.ID = d.ID
	result.ProcessDate = d.ProcessDate
	result.RecordCount = d.RecordCount
	result.ErrorCount = d.ErrorCount
	result.Status = d.Status
	result.Vendor = d.Vendor
	result.FeedName = d.FeedName
	result.FeedMethod = d.FeedMethod
	result.FileName = d.FileName
	if err := d.PreviousFeeds.Unmarshal(&result.PreviousFeeds); err != nil {
		return nil, fmt.Errorf("could not unmarshal previous_feeds, %v: %w", d.PreviousFeeds, err)
	}
	
	return result, nil
}

func (pr *PostgressRepository) GetFeedStatusDetails() ([]FeedStatusResultsDetailed, error) {
	sqlStatement := `
	WITH previous AS (
		SELECT
		    f.vendor, f.feed_name, f.feed_method, 
			json_agg(json_build_object('id', fs.id, 'process_date', fs.process_date,
			  'record_count', fs.record_count, 'error_count', fs.error_count, 'feed_status', fs.feed_status,
			  'vendor', f.vendor, 'feed_name', f.feed_name, 'feed_method', f.feed_method, 'file_name', fs.file_name,
			  'emr_logs', 'https://p-' || emr_step_id || '-' || emr_application_id || '.emrappui-prod.us-west-2.amazonaws.com/shs/history/' || emr_step_id || '/jobs/',
			  'data_quality_url', LOWER('http://' || $1 || '.s3-website-us-west-2.amazonaws.com/great_expectations/docs/expectations/' || replace(f.vendor,' ','_') || '_' ||  replace(f.feed_name,' ','_') || '_' || replace(replace(fs.file_name, '.txt', ''), '.csv', '') || '.html')
			)
			  ORDER BY fs.process_date DESC) previous_feeds
		FROM feed_status fs
		INNER JOIN feeds f
		ON fs.feed_id = f.id
		GROUP BY 1,2,3
	)
		
	SELECT 
	    fs.id, fs.process_date, fs.record_count, fs.error_count, fs.feed_status,
		f.vendor, f.feed_name, f.feed_method, fs.file_name, COALESCE(p.previous_feeds, '[{}]'::json) AS previous_feeds
	FROM feed_status fs
	INNER JOIN feeds f
	ON fs.feed_id = f.id
	LEFT JOIN previous p
	ON f.vendor = p.vendor AND f.feed_name = p.feed_name AND f.feed_method = p.feed_method
	WHERE fs.is_current = True
	AND f.vendor != 'test'
	;`
	rows, err := pr.db.Query(sqlStatement, S3_RESOURCES_BUCKET)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var feedStatuses []FeedStatusResultsDetailed
	for rows.Next() {
		var result feedStatusResultsDetailedDTO
		err = rows.Scan(&result.ID, &result.ProcessDate, &result.RecordCount, &result.ErrorCount, &result.Status, &result.Vendor, &result.FeedName, &result.FeedMethod, &result.FileName, &result.PreviousFeeds)
		if err != nil {
			return nil, err
		}
		feedStatus, err := result.ToFeedStatusResultsDetailed()
		if err != nil {
			return nil, err
		}
		feedStatuses = append(feedStatuses, *feedStatus)
	}
	if err := rows.Err(); err != nil {
		return nil, err
	}
	return feedStatuses, nil
}

func (pr *PostgressRepository) GetFeedStatusesAggregate(startDate string, endDate string) (FeedStatusAggregate, error) {
	sqlStatement := `
	SELECT 
		COUNT(DISTINCT fs.id) AS files,
		SUM(COALESCE(fs.record_count, 0)) AS rows,
		SUM(COALESCE(fs.error_count, 0)) AS errors
	FROM feed_status fs
	WHERE process_date::DATE BETWEEN TO_DATE($1,'YYYY-MM-DD') AND TO_DATE($2,'YYYY-MM-DD')`
	rows := pr.db.QueryRow(sqlStatement, startDate, endDate)
	var fsAgg FeedStatusAggregate
	err := rows.Scan(&fsAgg.Files, &fsAgg.Rows, &fsAgg.Errors)
	if err != nil {
		return FeedStatusAggregate{}, err
	}
	return fsAgg, nil
}