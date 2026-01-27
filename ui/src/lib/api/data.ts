import type {UpcomingGame} from "$lib/types/gameTypes";


export const getUpcomingGames = () => {
    return upcomingGames
}

const upcomingGames: UpcomingGame[] = [
    {
        game: {
            id: 'g1',
            leagueId: 'l1',
            startTime: new Date().toISOString(),
            location: 'Berea Rec Center • Court 2',
            costInCents: 500,
            isCanceled: false,
            createdAt: new Date().toISOString()
        },
        league: {
            id: 'l1',
            name: 'Sunday Morning Run',
            isActive: true,
            createdAt: new Date().toISOString()
        },
        attendance: {
            gameId: 'g1',
            userId: 'u1',
            checkedIn: false,
            createdAt: new Date().toISOString()
        },
        payment: {
            id: 'p1',
            gameId: 'g1',
            userId: 'u1',
            amountCents: 700,
            method: 'venmo',
            status: 'pending',
            paidAt: new Date().toISOString(),
            confirmedBy: 'admin1'
        }
    },
    {
        game: {
            id: 'g2',
            leagueId: 'l1',
            startTime: new Date().toISOString(),
            location: 'Berea Rec Center • Court 1',
            costInCents: 500,
            isCanceled: false,
            createdAt: new Date().toISOString()
        },
        league: {
            id: 'l1',
            name: 'Sunday Morning Run',
            isActive: true,
            createdAt: new Date().toISOString()
        },
        attendance: {
            gameId: 'g1',
            userId: 'u1',
            checkedIn: false,
            createdAt: new Date().toISOString()
        },
        payment: {
            id: 'p1',
            gameId: 'g1',
            userId: 'u1',
            amountCents: 9000,
            method: 'venmo',
            status: "completed",
            paidAt: new Date().toISOString(),
            confirmedBy: 'admin1'
        }
    },
    {
        game: {
            id: 'g2',
            leagueId: 'l1',
            startTime: new Date("2026/2/12").toISOString(),
            location: 'Edgewater Courts • Court 1',
            costInCents: 0,
            isCanceled: false,
            createdAt: new Date().toISOString()
        },
        league: {
            id: 'l1',
            name: 'Sunday Morning Run',
            isActive: true,
            createdAt: new Date().toISOString()
        }
    }
];