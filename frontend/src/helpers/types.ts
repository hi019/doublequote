import { DataResponse, ErrorResponse, InvalidParam } from "../api/types";

export const isErrorType = (error: any): error is ErrorResponse =>
  error.type && error.title;

export const isInvalidParamType = (error: any): error is InvalidParam =>
  error.name && error.reason;

export const isInvalidParamError = (error: any): error is Array<InvalidParam> =>
  Array.isArray(error) && error.every((e) => isInvalidParamType(e));

export const isDataType = (data: any): data is DataResponse =>
  typeof data == "object" && data.hasOwnProperty("data");
