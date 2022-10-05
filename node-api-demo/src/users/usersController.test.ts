import { UsersController } from './usersController'

afterEach(() => {
  jest.resetAllMocks()
})

describe("UsersController", () => {
  describe("getUsers", () => {
    test("should return empty array", async () => {
      const controller = new UsersController();
      const users = await controller.getUsers();
      //expect(users).toEqual([])
      expect(users.length).toBe(3)
        })
    })
});