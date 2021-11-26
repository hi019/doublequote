import { useAppDispatch } from "../../store/hooks";
import { useAuthCheckQuery } from "../../api";
import { setIsSignedIn } from "../slices/user";

export const AuthCheck = () => {
  const dispatch = useAppDispatch();
  const { error, isLoading } = useAuthCheckQuery();

  if (!isLoading) {
    dispatch(setIsSignedIn(error !== "Unauthorized."));
  }

  return null;
};
