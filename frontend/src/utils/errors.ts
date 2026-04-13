const defaultErrorMessage = 'Что-то пошло не так, попробуйте позже'

export const extractUserMessage = async (e: any): Promise<string> => {
  if (!e.response || !e.response.body) {
    return defaultErrorMessage
  }

  const text = await new Response(e.response.body).text()
  const data = JSON.parse(text)

  return data.userMessage || defaultErrorMessage
}
