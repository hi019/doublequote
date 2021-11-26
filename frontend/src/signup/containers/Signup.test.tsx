import "@testing-library/react";
import { Signup } from "./Signup";
import { render } from "../../testUtils";
import React from "react";
import { screen, waitFor } from "@testing-library/react";
import userEvent from "@testing-library/user-event";

test("renders without crashing", async () => {
  render(<Signup />);

  expect(screen.getByText("Let's get you started")).toBeInTheDocument();
  expect(screen.getByPlaceholderText("Email")).toBeInTheDocument();
  expect(screen.getByPlaceholderText("Password")).toBeInTheDocument();
  expect(screen.getByText("Sign up")).toBeInTheDocument();
});

test("requires email and password", async () => {
  render(<Signup />);
  userEvent.click(screen.getByRole("button"));

  await waitFor(() => {
    const errors = screen.queryAllByRole("alert");

    expect(errors).toHaveLength(2);

    expect(errors[0]).toHaveTextContent("An email is required.");
    expect(errors[1]).toHaveTextContent("A password is required.");
  });
});

test("submits signup info", async () => {
  render(<Signup />);
  userEvent.type(screen.getByPlaceholderText("Email"), "test@example.com");
  userEvent.type(screen.getByPlaceholderText("Password"), "password");
  userEvent.click(screen.getByRole("button"));

  await waitFor(() => expect(screen.queryAllByRole("alert")).toHaveLength(0));
  await waitFor(() =>
    expect(
      screen.getByText("Welcome to Doublequote. We're delighted to have you.")
    ).toBeInTheDocument()
  );
});
