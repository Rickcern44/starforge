import type {Game} from "$lib/types/Game";

export type League = {
    ID: string
    Name: string
    CreatedAt: Date
    UpdatedAt: Date
    DeletedAt: Date
    IsActive: boolean
    Games: Game[]
}