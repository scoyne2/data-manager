#!/usr/bin/env node
import 'source-map-support/register';
import * as cdk from '@aws-cdk/core';
import { DataManagerLambdaStack } from '../lib/data-manager-lambda-stack';

const app = new cdk.App();
new DataManagerLambdaStack(app, 'DataManagerLambdaStack');