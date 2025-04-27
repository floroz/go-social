export interface SignupRequestDTO {
  first_name: string;
  last_name: string;
  username: string;
  email: string;
  password: string;
}

export interface LoginRequestDTO {
  email: string;
  password: string;
}
