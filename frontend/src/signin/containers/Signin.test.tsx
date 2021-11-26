import "@testing-library/react";
import { render } from "../../testUtils";
import React from "react";
import { screen, waitFor } from "@testing-library/react";
import userEvent from "@testing-library/user-event";
import { BrowserRouter, MemoryRouter } from "react-router-dom";
import { Signin } from "./Signin";

const mockedNavigate = jest.fn();

jest.mock("react-router-dom", () => ({
  ...jest.requireActual("react-router-dom"),
  useNavigate: () => mockedNavigate,
}));

test("renders without crashing", async () => {
  render(
    <BrowserRouter>
      <Signin />
    </BrowserRouter>
  );

  expect(screen.getByText("Welcome back")).toBeInTheDocument();
  expect(screen.getByPlaceholderText("Email")).toBeInTheDocument();
  expect(screen.getByPlaceholderText("Password")).toBeInTheDocument();
  expect(screen.getByText("Sign in")).toBeInTheDocument();
});

test("requires email and password", async () => {
  render(
    <BrowserRouter>
      <Signin />
    </BrowserRouter>
  );
  userEvent.click(screen.getByRole("button"));

  await waitFor(() => {
    const errors = screen.queryAllByRole("alert");

    expect(errors).toHaveLength(2);

    expect(errors[0]).toHaveTextContent("An email is required.");
    expect(errors[1]).toHaveTextContent("A password is required.");
  });
});

test("submits signin info and redirects", async () => {
  render(
    <MemoryRouter initialEntries={[{ state: { path: "test" } }]}>
      <Signin />
    </MemoryRouter>
  );
  userEvent.type(screen.getByPlaceholderText("Email"), "test@example.com");
  userEvent.type(screen.getByPlaceholderText("Password"), "password");
  userEvent.click(screen.getByRole("button"));

  await waitFor(() => expect(mockedNavigate).toHaveBeenCalledWith("test"));
});

test("submits signin info and redirects to home", async () => {
  render(
    <MemoryRouter>
      <Signin />
    </MemoryRouter>
  );
  userEvent.type(screen.getByPlaceholderText("Email"), "test@example.com");
  userEvent.type(screen.getByPlaceholderText("Password"), "password");
  userEvent.click(screen.getByRole("button"));

  await waitFor(() => expect(mockedNavigate).toHaveBeenCalledWith("/"));
});
