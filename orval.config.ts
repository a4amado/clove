import { defineConfig } from 'orval';

export default defineConfig({
    myApi: {
        input: {
            target: 'http://localhost/api/v1/docs/openapi.json',
        },
        output: {
            target: 'frontend/src/api/generated.ts',
            client: 'react-query',
            mode: 'split',
        },
    },
});