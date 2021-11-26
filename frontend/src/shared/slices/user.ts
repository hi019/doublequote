import { createSlice, PayloadAction } from "@reduxjs/toolkit";

export interface UserState {
  readonly isSignedIn: boolean;
}

const initialState: UserState = {
  isSignedIn: false,
};

export const userSlice = createSlice({
  name: "user",
  initialState,
  reducers: {
    setIsSignedIn: (state, action: PayloadAction<boolean>) => ({
      ...state,
      isSignedIn: action.payload,
    }),
  },
});

export const { setIsSignedIn } = userSlice.actions;
