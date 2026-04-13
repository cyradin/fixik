export const extractUserMessage = async (e: any): Promise<string> => {
  const text = await new Response(e.response.body).text()
  const data = JSON.parse(text)

  return data.userMessage || data.error
}
