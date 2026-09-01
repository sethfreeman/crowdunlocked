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
const onVercel = !!process.env.VERCEL;

// On Vercel we MUST use OIDC. If the role ARN is missing there, fail loudly
// instead of silently falling through to the default chain (which has no
// credentials in the Vercel runtime and yields a confusing 500).
if (onVercel && !roleArn) {
  console.error(
    'AWS_ROLE_ARN is not set in this Vercel environment. DynamoDB calls will ' +
      'fail. Set AWS_ROLE_ARN (and AWS_REGION) for the Production/Preview scopes.'
  );
}

const client = new DynamoDBClient({
  region,
  ...(roleArn
    ? { credentials: awsCredentialsProvider({ roleArn }) }
    : {}),
});

export const docClient = DynamoDBDocumentClient.from(client);
