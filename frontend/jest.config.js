/** @type {import('jest').Config} */
export default {
  setupFiles: ['./jest.setup.ts'],
  transform: {
    '^.+\\.svelte$': 'svelte-jester',
    '^.+\\.ts$': ['ts-jest', {
      useESM: true,
    }]
  },
  moduleFileExtensions: ['js', 'ts', 'svelte'],
  testEnvironment: 'node',
  extensionsToTreatAsEsm: ['.ts', '.svelte']
}