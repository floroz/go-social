import { beforeAll, afterEach, afterAll } from "vitest";
import { server } from "./mocks/server"; // Import the MSW server

// Establish API mocking before all tests.
beforeAll(() => server.listen({ onUnhandledRequest: "error" }));

// Reset any request handlers that we may add during the tests,
// so they don't affect other tests.
afterEach(() => server.resetHandlers());

// Clean up after the tests are finished.
afterAll(() => server.close());

// Optional: Add any other global setup, like cleaning up jsdom
// afterEach(() => {
//   cleanup(); // if using @testing-library/react
// });
