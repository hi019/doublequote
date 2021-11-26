import { rest } from "msw";

const baseUrl = process.env.API_URL;

export const handlers = [
  rest.post(`${baseUrl}pub/register`, (req, res, ctx) => {
    return res(
      ctx.json({
        data: {
          require_email_verification: false,
        },
      })
    );
  }),

  rest.post(`${baseUrl}pub/login`, (req, res, ctx) => {
    return res(ctx.json({}));
  }),
];
