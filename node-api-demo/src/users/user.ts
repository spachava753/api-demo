/**
 * User objects allow you to associate actions performed
 * in the system with the user that performed them.
 * The User object contains common information across
 * every user in the system with role.
 */
export interface User {
    id: number;

     /**
      * The Role of the user 
      */
    role: string;
      /**
      * The Name of the user 
      */
    name: string;

  }