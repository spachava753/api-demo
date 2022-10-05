// src/users/usersService.ts
import { User } from "./user";

// A post request should not contain an id.
export type UserCreationParams = Pick<User, "name" | "role">;
export type UserUpdateParams = Pick<User, "name" | "role">;

const users: User[] = [
    {
        id: 1,
        name: "Jane1 Doe",
        role: "Admin"
    },
    {
        id: 2,
        name: "Jane2 Doe",
        role: "Admin"
    },
    {
        id: 3,
        name: "Jane3 Doe",
        role: "Admin"
    },
]


export class UsersService {
    public get(id: number): User {

        const cUser: User = users.find((user) => {
            return user.id === id
        })!;

        return cUser;

    }

    public delete(id: number): number {
        const cUserIdx: number = users.indexOf(this.get(id))!;
        if (cUserIdx != -1) {
            users.splice(cUserIdx, 1)
            return id
        } else {
            return -1
        }
    }

    public getAll(): User[] {
        return users;
    }

    public create(userCreationParams: UserCreationParams): User {
        const user: User = {
            id: Math.floor(Math.random() * 10000), // Random
            ...userCreationParams,
        };
        users.push(user)
        return user;
    }

    public update(id: number, userUpdateParams: UserUpdateParams): User {
        const cUser: User = users.find((user) => {
            return user.id === id
        })!;
        cUser.name = userUpdateParams.name
        cUser.role = userUpdateParams.role
        return cUser;
    }
}
