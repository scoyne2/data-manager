package datapreview

import (
	"fmt"
	"time"


	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/athena"
)

type DataPreview struct{
	Header []string
	Rows [][]string
}

func DataPreviewAthena(vendor string, feedName string, fileName string, s3_bucket string) (DataPreview, error){
	database := fmt.Sprintf("data_manager_output_%s", vendor)
	glue_table_name := fmt.Sprintf("%s.%s", database, feedName)
	raw, err := queryAthena(database, glue_table_name, fileName, s3_bucket)
	if err != nil {
		return DataPreview{}, err
	}
	result, err2 := formatResults(raw)
	if err2 != nil {
		return DataPreview{}, err2
	}

	return result, nil
}

func formatResults(results *athena.ResultSet)(DataPreview, error){
	rows := results.Rows

	var outputRows [][]string
	for _, r := range rows {
		var outputRow []string
		for _, d := range r.Data {
			outputRow = append(outputRow, *d.VarCharValue)
		}
		outputRows = append(outputRows, outputRow)
	}

	dataPreview := DataPreview{outputRows[0], outputRows[1:]}
	
	return dataPreview, nil
}

func queryAthena(database string, glue_table_name string, fileName string, s3_bucket string) (*athena.ResultSet, error){
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
	outputLocation := fmt.Sprintf("s3://%s/data_preview/%s/%s/", s3_bucket, glue_table_name, fileName)
	r.SetOutputLocation(outputLocation)
	s.SetResultConfiguration(&r)

	result, err := svc.StartQueryExecution(&s)
	if err != nil {
		fmt.Println(err)
		return nil, fmt.Errorf("location: %s, Error running query: %s", outputLocation, err)
	}
	fmt.Println("StartQueryExecution result:")
	fmt.Println(result.GoString())

	var qri athena.GetQueryExecutionInput
	qri.SetQueryExecutionId(*result.QueryExecutionId)

	duration := time.Duration(2) * time.Second
	time.Sleep(duration)

	var ip athena.GetQueryResultsInput
	ip.SetQueryExecutionId(*result.QueryExecutionId)

	op, err := svc.GetQueryResults(&ip)
	if err != nil {
		fmt.Println(err)
		return nil, err
	}
	return op.ResultSet, nil	
}