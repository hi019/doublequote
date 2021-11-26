export interface SignupRequest {
  email: string;
  password: string;
}

export interface LoginRequest {
  email: string;
  password: string;
}

export interface VerifyEmailRequest {
  token: string;
}

export interface InvalidParam {
  name: string;
  reason: string;
}

export interface ErrorResponse {
  type: string;
  title: string;
  invalidParams?: Array<InvalidParam>;
}

export interface DataResponse {
  data: unknown;
}
