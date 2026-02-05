import type { NextApiRequest, NextApiResponse } from 'next';

interface HealthResponse {
  status: 'ok';
  timestamp: string;
  version: string;
}

export default function handler(
  req: NextApiRequest,
  res: NextApiResponse<HealthResponse>
) {
  if (req.method !== 'GET') {
    res.setHeader('Allow', ['GET']);
    return res.status(405).end(`Method ${req.method} Not Allowed`);
  }

  const response: HealthResponse = {
    status: 'ok',
    timestamp: new Date().toISOString(),
    version: '1.0.0'
  };

  res.status(200).json(response);
}