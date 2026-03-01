module.exports = {
  envoy: {
    input: {
      target: '../server/docs/swagger.json', 
    },
    output: {
      mode: 'tags-split',        
      target: './src/api/generated/envoy.ts',
      schemas: './src/api/generated/model',
      client: 'react-query',       
      httpClient: 'fetch',       
      override: {
        mutator: {
          path: './src/api/custom-fetch.ts', 
          name: 'customInstance',
        },
      },
    },
  },
};