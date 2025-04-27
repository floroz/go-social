import { http, HttpResponse } from "msw";
import config from "@/config";
import { mockUser } from "@/mocks/data/user";

const API_BASE_URL = config.apiBaseUrl;

export const handlers = [
  // Mock the POST /v1/auth/signup endpoint
  http.post(`${API_BASE_URL}/v1/auth/signup`, async (/*{ request }*/) => {
    // Removed unused request param
    // You can optionally inspect the request body if needed
    // const body = await request.json() as { data: SignupRequestDTO };
    // console.log('MSW intercepted signup:', body.data);

    // Simulate a successful signup
    // Backend returns { data: User }
    return HttpResponse.json({ data: mockUser }, { status: 201 });

    // --- Example Error Response ---
    // return HttpResponse.json(
    //   { errors: [{ message: 'Username already exists' }] },
    //   { status: 400 }
    // );
  }),

  // Mock the POST /v1/auth/login endpoint
  http.post(`${API_BASE_URL}/v1/auth/login`, async (/*{ request }*/) => {
    // Simulate a successful login
    const mockLoginResponse = { token: "mock-jwt-token" };
    // Return the token wrapped in the standard { data: ... } structure
    return HttpResponse.json({ data: mockLoginResponse }, { status: 200 });

    // --- Example Error Response (401 Unauthorized) ---
    // return HttpResponse.json(
    //   {
    //     errors: [
    //       { code: "UNAUTHORIZED", message: "Invalid email or password." },
    //     ],
    //   },
    //   { status: 401 }
    // );

    // --- Example Error Response (400 Bad Request) ---
    // return HttpResponse.json(
    //   {
    //     errors: [
    //       {
    //         code: "VALIDATION_ERROR",
    //         message: "Email format is invalid.",
    //         field: "email",
    //       },
    //     ],
    //   },
    //   { status: 400 }
    // );
  }),

  // TODO: Add handlers for posts, comments etc. later
];
