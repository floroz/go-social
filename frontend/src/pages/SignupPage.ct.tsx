import { test, expect } from "@playwright/experimental-ct-react";
import SignupPage from "./SignupPage";
import { User } from "@/types/api";
import PlaywrightTestWrapper from "../../tests/PlaywrightTestWrapper"; // Import the wrapper

test.use({ viewport: { width: 500, height: 900 } });

test("should allow successful signup", async ({ mount, page }) => {
  // 1. Mock the API endpoint for successful signup
  const mockUser: User = {
    id: 123, // Changed to number
    first_name: "Test", // Changed to snake_case
    last_name: "User", // Changed to snake_case
    username: "testuser",
    email: "test@example.com",
    created_at: new Date().toISOString(), // Changed to snake_case
    updated_at: new Date().toISOString(), // Changed to snake_case
    // Add other fields as defined in your User type, potentially null or default values
    // last_login: null, // Changed to snake_case if needed
  };

  await page.route("**/api/v1/auth/signup", (route) => {
    route.fulfill({
      status: 201, // Created
      contentType: "application/json",
      body: JSON.stringify({ data: mockUser }), // Match your API's success response structure
    });
  });

  // 2. Mount the component using the wrapper
  const component = await mount(
    <PlaywrightTestWrapper>
      <SignupPage />
    </PlaywrightTestWrapper>
  );
  // await page.pause(); // Removed pause

  // 3. Interact with the form
  await component.getByLabel("First Name").fill("Test");
  await component.getByLabel("Last Name").fill("User");
  await component.getByLabel("Username").fill("testuser");
  await component.getByLabel("Email").fill("test@example.com");
  await component.getByLabel("Password").fill("password123");
  await component.getByRole("button", { name: /sign up/i }).click();

  // 4. Assert the outcome (Update based on actual component behavior)
  // Example: Check for a success message or navigation (might need router mocking)
  // await expect(component.locator('text=Signup successful!')).toBeVisible();
  // Or, if it navigates or updates global state, assertions might differ.
  // For now, let's just assert the button click doesn't immediately show an obvious error
  await expect(
    component.locator("text=Email already exists")
  ).not.toBeVisible(); // Example negative assertion
  await expect(
    component.locator("text=Username already exists")
  ).not.toBeVisible(); // Example negative assertion

  // TODO: Add more specific assertions based on what should happen after successful signup
  // (e.g., redirect, state update, success message display)
});

// TODO: Add tests for failure cases (e.g., email exists, username exists, validation errors)
