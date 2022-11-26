import * as cdk from '@aws-cdk/core';
import * as path from "path";
import * as lambda from '@aws-cdk/aws-lambda';
import * as s3 from '@aws-cdk/aws-s3';
import * as lambdaEventSources from '@aws-cdk/aws-lambda-event-sources';


export class DataManagerLambdaStack extends cdk.Stack {
  constructor(scope: cdk.Construct, id: string, props?: cdk.StackProps) {
    super(scope, id, props);
    // Build the lambda
    const lambdaFunction = this.buildLambda('data-manager-lambda', path.join(__dirname, '../../lambda'), 'main');

    // Create the S3 bucket for trigger
    const bucketName = process.env.TRIGGER_BUCKET_NAME
    const bucket = new s3.Bucket(this, 'data-manager-trigger-bucket', {
      autoDeleteObjects: false,
      removalPolicy: cdk.RemovalPolicy.RETAIN,
      bucketName: bucketName
    });

    // Create S3 trigger, when any file is dropped in the bucket it will send an event to the Lambda
    const s3PutEventSource = new lambdaEventSources.S3EventSource(bucket, {
      events: [
        s3.EventType.OBJECT_CREATED_PUT
      ]
    });

    lambdaFunction.addEventSource(s3PutEventSource);

  }

  /**
   * buildLambda build the code and create the lambda
   * @param id - CDK id for this lambda
   * @param lambdaPath - Location of the code
   * @param handler - name of the handler to call for this lambda
   */
  buildLambda(id: string, lambdaPath: string, handler: string): lambda.Function {
    const environment = {
      CGO_ENABLED: '0',
      GOOS: 'linux',
      GOARCH: 'amd64',
    };
    return new lambda.Function(this, id, {
      code: lambda.Code.fromAsset(lambdaPath, {
        bundling: {
          image: lambda.Runtime.GO_1_X.bundlingDockerImage,
          user: "root",
          environment,
          command: [
            'bash', '-c', [
              'make vendor',
              'make lambda-build',
            ].join(' && ')
          ]
        },
      }),
      handler,
      runtime: lambda.Runtime.GO_1_X,
    });
  }
}