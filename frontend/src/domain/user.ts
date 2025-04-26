// Corresponds to the Go backend's domain.User struct

export interface User {
  id: number; // Assuming int64 maps to number in TS/JSON
  first_name: string;
  last_name: string;
  username: string;
  // Password is intentionally omitted (json:"-")
  email: string;
  created_at: string; // Dates usually come as ISO strings in JSON
  updated_at: string; // Dates usually come as ISO strings in JSON
  last_login?: string | null; // Optional and nullable date string
}

// We can add other related DTO interfaces here later if needed
// e.g., CreateUserDTO, LoginUserDTO
