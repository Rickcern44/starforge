export class User {
    constructor(public readonly id: string,
                public email: string,
                public name: string,
                public roles: string[]) {}

    hasRole(role: string): boolean {
        return this.roles.includes(role);
    }
}