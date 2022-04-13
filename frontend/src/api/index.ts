import {
  BaseQueryFn,
  createApi,
  fetchBaseQuery,
} from "@reduxjs/toolkit/query/react";
import type { User } from "../types";
import {
  ErrorResponse,
  LoginRequest,
  SignupRequest,
  VerifyEmailRequest,
} from "./types";
import { isDataType, isErrorType } from "../helpers/types";

const baseQuery = fetchBaseQuery({
  baseUrl: process.env.REACT_APP_API_URL,
  mode: process.env.REACT_APP_API_URL?.startsWith("http")
    ? "cors"
    : "same-origin",
  credentials: "include",
});

const customBaseQuery: BaseQueryFn<any, unknown, ErrorResponse> = async (
  args,
  a,
  extraOptions
) => {
  const { data, error, meta } = await baseQuery(args, a, extraOptions);

  if (meta?.response?.ok) {
    return { data: isDataType(data) ? data.data : {} };
  }

  if (isErrorType(error?.data) && error !== undefined) {
    return {
      error: {
        type: error.data.type,
        title: error.data.title,
        invalidParams: error.data.invalidParams,
      },
    };
  }

  return { data: {} };
};

export const api = createApi({
  reducerPath: "api",
  baseQuery: customBaseQuery,

  endpoints: (builder) => ({
    signup: builder.mutation<User, SignupRequest>({
      query: (req) => ({
        url: "pub/register",
        method: "POST",
        body: req,
      }),
    }),
    signin: builder.mutation<User, LoginRequest>({
      query: (req) => ({
        url: "pub/login",
        method: "POST",
        body: req,
      }),
    }),
    authCheck: builder.query<void, void>({
      query: () => ({
        url: "authcheck",
        method: "GET",
      }),
    }),
    verifyEmail: builder.mutation<void, VerifyEmailRequest>({
      query: (req) => ({
        url: "verify-email",
        method: "POST",
        body: req,
      }),
    }),
  }),
});

export const {
  useSignupMutation,
  useSigninMutation,
  useAuthCheckQuery,
  useVerifyEmailMutation,
} = api;
