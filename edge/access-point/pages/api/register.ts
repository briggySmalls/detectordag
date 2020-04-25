import { NextApiRequest, NextApiResponse } from 'next'

// Handle form submission
export default (req: NextApiRequest, res: NextApiResponse) => {
  // Print the content
  console.log(req.body);
  // Assemble the response
  res.statusCode = 200
  res.setHeader('Content-Type', 'application/json')
  res.end(JSON.stringify({ name: 'John Doe' }))
}
