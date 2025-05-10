import { test, expect } from "@playwright/experimental-ct-react";
import LoginPage from "./LoginPage";
import { User } from "@/types/api";
import PlaywrightTestWrapper from "../../tests/PlaywrightTestWrapper";

test.use({ viewport: { width: 500, height: 900 } });

test("should allow successful login", async ({ mount, page }) => {
  // 1. Mock the API endpoint for successful login
  const mockUser: User = {
    id: 123,
    first_name: "Test",
    last_name: "User",
    username: "testuser",
    email: "test@example.com",
    created_at: new Date().toISOString(),
    updated_at: new Date().toISOString(),
  };

  await page.route("**/api/v1/auth/login", (route) => {
    route.fulfill({
      status: 200,
      contentType: "application/json",
      body: JSON.stringify({ data: mockUser }),
    });
  });

  // 2. Mount the component using the wrapper
  const component = await mount(
    <PlaywrightTestWrapper>
      <LoginPage />
    </PlaywrightTestWrapper>
  );

  // 3. Interact with the form
  await component.getByLabel("Email").fill("test@example.com");
  await component.getByLabel("Password").fill("password123");
  await component.getByRole("button", { name: /login/i }).click();

  // 4. Assert the outcome
  // Verify the button state changes during loading
  await expect(component.getByRole("button", { name: /logging in/i })).toBeVisible();
  
  // After successful login, the component should not show any error messages
  await expect(component.locator("text=Invalid email or password")).not.toBeVisible();
});

test("should display validation errors for invalid inputs", async ({ mount }) => {
  // Mount the component
  const component = await mount(
    <PlaywrightTestWrapper>
      <LoginPage />
    </PlaywrightTestWrapper>
  );

  // Try to submit with empty fields
  await component.getByRole("button", { name: /login/i }).click();

  // Check for validation errors
  await expect(component.getByText("Invalid email address.")).toBeVisible();
  await expect(component.getByText("Password is required.")).toBeVisible();
});

test("should handle login failure", async ({ mount, page }) => {
  // Mock the API endpoint for failed login
  await page.route("**/api/v1/auth/login", (route) => {
    route.fulfill({
      status: 401,
      contentType: "application/json",
      body: JSON.stringify({
        error: {
          message: "Invalid email or password",
          code: "INVALID_CREDENTIALS"
        }
      }),
    });
  });

  // Mount the component
  const component = await mount(
    <PlaywrightTestWrapper>
      <LoginPage />
    </PlaywrightTestWrapper>
  );

  // Fill in the form
  await component.getByLabel("Email").fill("test@example.com");
  await component.getByLabel("Password").fill("wrongpassword");
  await component.getByRole("button", { name: /login/i }).click();

  // Check for error message
  await expect(component.getByText("Invalid email or password")).toBeVisible();
});

test("should handle network error", async ({ mount, page }) => {
  // Mock a network error
  await page.route("**/api/v1/auth/login", (route) => {
    route.abort("failed");
  });

  // Mount the component
  const component = await mount(
    <PlaywrightTestWrapper>
      <LoginPage />
    </PlaywrightTestWrapper>
  );

  // Fill in the form
  await component.getByLabel("Email").fill("test@example.com");
  await component.getByLabel("Password").fill("password123");
  await component.getByRole("button", { name: /login/i }).click();

  // Check for network error message
  await expect(component.getByText(/network error|failed to fetch/i)).toBeVisible();
});
