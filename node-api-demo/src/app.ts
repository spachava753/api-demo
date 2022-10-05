// src/app.ts

import methodOverride from "method-override";
import bodyParser from "body-parser";
import { ValidateError } from "tsoa";
import swaggerJson from "../build/swagger.json";
import * as swaggerUI from "swagger-ui-express";
import morgan from "morgan";

import express, {
  Response as ExResponse,
  Request as ExRequest,
  NextFunction,
  json, urlencoded,
} from "express";

import { RegisterRoutes } from "../build/routes";

export const app = express();

// Use body parser to read sent json payloads

app.use(
  urlencoded({
    extended: true,
  })
);
app.use(json());

app.use(morgan("tiny"));
app.use(bodyParser.urlencoded({ extended: true }));
app.use(bodyParser.json());
app.use(methodOverride());
app.use(["/openapi", "/docs", "/swagger"], swaggerUI.serve, swaggerUI.setup(swaggerJson));

RegisterRoutes(app);

app.use(function errorHandler(
  err: unknown,
  req: ExRequest,
  res: ExResponse,
  next: NextFunction
): ExResponse | void {
  if (err instanceof ValidateError) {
    console.log("Caught Validation Error for:",req.path, err.fields);
    return res.status(422).json({
      message: "Validation Failed",
      details: err?.fields,
    });
  }
  if (err instanceof Error) {
    return res.status(500).json({
      message: "Internal Server Error",
    });
  }

  next();
});

app.use(function notFoundHandler(_req, res: ExResponse) {
  res.status(404).send({
    message: "API endpoint not found",
  });
});