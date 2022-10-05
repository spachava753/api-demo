// src/users/usersController.ts
import {
    Body,
    Controller,
    Get,
    Path,
    Post,
    Put,
    Route,
    Response,
    SuccessResponse,
    Tags,
    Example,
    TsoaResponse,
    Res, Delete,
} from "tsoa";
import { ValidateError } from "tsoa";
import { User } from "./user";
import { UsersService, UserCreationParams, UserUpdateParams } from "./usersService";

interface ValidateErrorJSON {
    message: "Validation failed";
    details: { [name: string]: unknown };
}

@Route("/users")
@Tags("User")
export class UsersController extends Controller {

    /**
     * Retrieves the details of an existing user.
     * Supply the unique user ID from either and receive corresponding user details.
     * @param userId The user's identifier
     * @example userId "5"
     * @example userId "2"
     */
    @Example<User>({
        id: 111, //"52907745-7672-470e-a803-a2f8feb52944",
        name: "John",
        role: "Admin",
    })
    @Get("{userId}")
    public async getUser(
        @Path() userId: number, @Res() notFoundResponse: TsoaResponse<404, { reason: string }>
    ): Promise<User> {
        const cUser = new UsersService().get(userId);

        if (!cUser) {
            return notFoundResponse(404, { reason: "User Doesn't Exist. Please provide a valid id" });
        }

        return cUser;
    }

    /**
     * Retrieves the details of an All Users.
     */
    @Get()
    public async getUsers(): Promise<User[]> {
        return new UsersService().getAll();
    }

    @Delete("{userId}")
    public async deleteUser(
        @Path() userId: number, @Res() notFoundResponse: TsoaResponse<404, { reason: string }>
    ): Promise<string> {
        const id = new UsersService().delete(userId);
        if (id == -1) {
            return notFoundResponse(404, { reason: "User Doesn't Exist. Please provide a valid id" });
        }
        return "Success";
    }

    //@Response<ValidateErrorJSON>(422, "Validation Failed")
    @Post()
    @SuccessResponse("201", "Created") // Custom success response
    @Response<ValidateErrorJSON>(422, "Validation Failed", {
        message: "Validation failed",
        details: {
            requestBody: {
                message: "id is an excess property and therefore not allowed",
                value: "52907745-7672-470e-a803-a2f8feb52944",
            },
        },
    })
    public async createUser(
        @Body() requestBody: UserCreationParams
    ): Promise<void> {

        if (requestBody.name.startsWith("test") == true) {
            this.setStatus(422);
            let error: ValidateError = { status: 423, name: "verror", fields: {}, message: "Validation Failed" };
            throw error //ValidateError()
        }

        this.setStatus(201); // set return status 201
        new UsersService().create(requestBody);
        return;
    }


    //@Response<ValidateErrorJSON>(422, "Validation Failed")
    @Put("{userId}")
    @SuccessResponse("201", "Created") // Custom success response
    @Response<ValidateErrorJSON>(422, "Validation Failed", {
        message: "Validation failed",
        details: {
            requestBody: {
                message: "id is an excess property and therefore not allowed",
                value: "52907745-7672-470e-a803-a2f8feb52944",
            },
        },
    })
    public async updateUser(
        @Path() userId: number, @Body() requestBody: UserUpdateParams
    ): Promise<void> {

        if (requestBody.name.startsWith("test") == true) {
            this.setStatus(422);
            let error: ValidateError = { status: 423, name: "verror", fields: {}, message: "Validation Failed" };
            throw error //ValidateError()
        }

        this.setStatus(201); // set return status 201
        new UsersService().update(userId, requestBody);
        return;
    }
}
  