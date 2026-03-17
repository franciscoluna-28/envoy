import createClient from 'openapi-fetch';
import type { paths } from './api';

const client = createClient<paths>({
  baseUrl: 'http://localhost:8080/api/v1/',
  credentials: 'include',
});

export default client;
