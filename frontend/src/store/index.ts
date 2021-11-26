import { configureStore } from "@reduxjs/toolkit";
import { api } from "../api";
import { setupListeners } from "@reduxjs/toolkit/query";
import { userSlice } from "../shared/slices/user";
import { loadState, saveState } from "./localStorage";

const persistedState = loadState();

export const store = configureStore({
  reducer: {
    [api.reducerPath]: api.reducer,
    user: userSlice.reducer,
  },
  preloadedState: persistedState,
  middleware: (getDefaultMiddleware) =>
    getDefaultMiddleware().concat(api.middleware),
});

store.subscribe(() => {
  saveState({
    user: store.getState().user,
  });
});

setupListeners(store.dispatch);

export type RootState = ReturnType<typeof store.getState>;
export type AppDispatch = typeof store.dispatch;
