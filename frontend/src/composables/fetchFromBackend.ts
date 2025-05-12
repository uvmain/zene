export async function backendFetchRequest(path: string, options = {}): Promise<Response> {
  const url = `/api/${path}`
  const response = await fetch(url, options)
  return response
}
