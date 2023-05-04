package data_preview

import (
	"fmt"
	"time"
	"encoding/json"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/athena"
)

type DataPreview struct{
	Header []string
	Rows [][]string
}

func DataPreviewAthena(vendor string, feedName string, fileName string) ([]byte, error){
	database := fmt.Sprintf("data_manager_output_%s", vendor)
	glue_table_name := fmt.Sprintf("%s.%s", database, feedName)
	results, err := queryAthena(database, glue_table_name, fileName)
	if err != nil {
		return nil, err
	}
	formattedResults, err := formatResults(results)
	if err != nil {
		return nil, err
	}
	jsonBytes, err := json.Marshal(formattedResults)
    if err != nil {
        panic(err)
    }

	return jsonBytes, nil
}

func formatResults(results *athena.ResultSet)(DataPreview, error){
	rows := results.Rows
	columnInfo := results.ResultSetMetadata.ColumnInfo
	var header []string
	for _, info := range columnInfo {
		header = append(header, *info.Name)
	}

	var outputRows [][]string
	for _, r := range rows {
		var outputRow []string
		for _, d := range r.Data {
			outputRow = append(outputRow, *d.VarCharValue)
		}
		outputRows = append(outputRows, outputRow)
	}

	dataPreview := DataPreview{header, outputRows}
	
	return dataPreview, nil
}

func queryAthena(database string, glue_table_name string, fileName string) (*athena.ResultSet, error){
	awscfg := &aws.Config{}
	awscfg.WithRegion("us-west-2")
	sess := session.Must(session.NewSession(awscfg))

	svc := athena.New(sess, aws.NewConfig().WithRegion("us-west-2"))
	var s athena.StartQueryExecutionInput
	query := fmt.Sprintf("SELECT * FROM %s WHERE file_name = '%s' LIMIT 25", glue_table_name, fileName)
	s.SetQueryString(query)

	var q athena.QueryExecutionContext
	q.SetDatabase(database)
	s.SetQueryExecutionContext(&q)

	var r athena.ResultConfiguration
	r.SetOutputLocation("s3://TestBucket")
	s.SetResultConfiguration(&r)

	result, err := svc.StartQueryExecution(&s)
	if err != nil {
		fmt.Println(err)
		return nil, err
	}
	fmt.Println("StartQueryExecution result:")
	fmt.Println(result.GoString())

	var qri athena.GetQueryExecutionInput
	qri.SetQueryExecutionId(*result.QueryExecutionId)

	var qrop *athena.GetQueryExecutionOutput
	duration := time.Duration(2) * time.Second // Pause for 2 seconds

	for {
		qrop, err = svc.GetQueryExecution(&qri)
		if err != nil {
			fmt.Println(err)
			return nil, err
		}
		if *qrop.QueryExecution.Status.State != "RUNNING" {
			break
		}
		fmt.Println("waiting.")
		time.Sleep(duration)

	}
	if *qrop.QueryExecution.Status.State == "SUCCEEDED" {

		var ip athena.GetQueryResultsInput
		ip.SetQueryExecutionId(*result.QueryExecutionId)

		op, err := svc.GetQueryResults(&ip)
		if err != nil {
			fmt.Println(err)
			return nil, err
		}
		return op.ResultSet, nil
	} else {
		fmt.Println(*qrop.QueryExecution.Status.State)
		return nil, nil
	}
	
}