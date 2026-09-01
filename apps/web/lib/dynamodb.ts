import { DynamoDBClient } from '@aws-sdk/client-dynamodb';
import { DynamoDBDocumentClient } from '@aws-sdk/lib-dynamodb';
import { awsCredentialsProvider } from '@vercel/oidc-aws-credentials-provider';

/**
 * Shared DynamoDB document client.
 *
 * On Vercel, `AWS_ROLE_ARN` is set and we use Vercel's OIDC credentials
 * provider, which exchanges the Vercel OIDC token for short-lived AWS
 * credentials via AssumeRoleWithWebIdentity. No long-lived keys.
 *
 * Locally (no `AWS_ROLE_ARN`), we fall back to the default AWS credential
 * chain (e.g. an `AWS_PROFILE`), so `npm run dev` works against the dev tables.
 */
const region = process.env.AWS_REGION || 'us-west-2';
const roleArn = process.env.AWS_ROLE_ARN;

const client = new DynamoDBClient({
  region,
  ...(roleArn
    ? { credentials: awsCredentialsProvider({ roleArn }) }
    : {}),
});

export const docClient = DynamoDBDocumentClient.from(client);
