export const customInstance = async <T>(
  { url, method, params, data, headers }: {
    url: string;
    method: 'GET' | 'POST' | 'PUT' | 'DELETE' | 'PATCH';
    params?: any;
    data?: any;
    headers?: any;
  },
): Promise<T> => {
  const baseUrl = 'http://localhost:8080/api/v1';
  
  const query = params ? `?${new URLSearchParams(params).toString()}` : '';

  const response = await fetch(`${baseUrl}${url}${query}`, {
    method,
    headers: {
      'Content-Type': 'application/json',
      ...headers,
    },
    body: data ? JSON.stringify(data) : undefined,
    credentials: 'include', 
  });

  if (!response.ok) {
    throw new Error(`API error: ${response.statusText}`);
  }

  return response.json();
};